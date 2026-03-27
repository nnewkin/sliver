package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MonkeyCode/sliver-web/internal/api"
	"github.com/MonkeyCode/sliver-web/internal/bot"
	"github.com/MonkeyCode/sliver-web/internal/core"
	"github.com/MonkeyCode/sliver-web/internal/db"
	"github.com/MonkeyCode/sliver-web/internal/rpc"
	"github.com/MonkeyCode/sliver-web/internal/ws"
	"github.com/sirupsen/logrus"
)

var (
	configPath = flag.String("config", "sliver-web.conf", "Path to configuration file")
	debug      = flag.Bool("debug", false, "Enable debug logging")
)

func main() {
	flag.Parse()

	log := logrus.New()
	if *debug {
		log.SetLevel(logrus.DebugLevel)
	}

	config, err := core.LoadConfig(*configPath)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	if err := os.MkdirAll(config.Database.Path, 0755); err != nil {
		log.WithError(err).Fatal("Failed to create database directory")
	}

	database, err := db.NewDatabase(config.Database.Path)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}
	defer database.Close()

	sliverClient, err := rpc.NewSliverClient(config.Sliver.GRPCAddr, config.Sliver.CertsDir)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to Sliver server")
	}

	wsServer := ws.NewServer()
	go wsServer.Start()

	botServer := bot.NewBotServer(config, database, sliverClient, wsServer)
	if config.Telegram.Enabled {
		if err := botServer.Start(); err != nil {
			log.WithError(err).Warn("Failed to start Telegram bot")
		}
	}

	apiServer := api.NewServer(config, database, sliverClient, wsServer, botServer)
	go apiServer.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	botServer.Stop()
	wsServer.Stop()
	apiServer.Stop()
	fmt.Println("Server stopped")
}
