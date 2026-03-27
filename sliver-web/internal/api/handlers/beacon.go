package handlers

import (
	"context"
	"net/http"

	"github.com/MonkeyCode/sliver-web/internal/db/repository"
	"github.com/MonkeyCode/sliver-web/internal/rpc/services"
	"github.com/gin-gonic/gin"
)

type BeaconHandler struct {
	beaconService *services.BeaconService
	auditRepo     *repository.AuditRepository
}

func NewBeaconHandler(beaconService *services.BeaconService, auditRepo *repository.AuditRepository) *BeaconHandler {
	return &BeaconHandler{
		beaconService: beaconService,
		auditRepo:     auditRepo,
	}
}

func (h *BeaconHandler) ListBeacons(c *gin.Context) {
	ctx := context.Background()
	beacons, err := h.beaconService.GetAllBeacons(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to get beacons",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": beacons,
	})
}

func (h *BeaconHandler) GetBeacon(c *gin.Context) {
	beaconID := c.Param("id")
	ctx := context.Background()

	beacon, err := h.beaconService.GetBeacon(ctx, beaconID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    10005,
			"message": "beacon not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": beacon,
	})
}

func (h *BeaconHandler) GetBeaconTasks(c *gin.Context) {
	beaconID := c.Param("id")
	ctx := context.Background()

	tasks, err := h.beaconService.GetBeaconTasks(ctx, beaconID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to get beacon tasks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": tasks,
	})
}

func (h *BeaconHandler) ExecuteBeaconTask(c *gin.Context) {
	beaconID := c.Param("id")

	var req struct {
		Command string `json:"command" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	taskID, err := h.beaconService.ExecuteCommand(ctx, beaconID, []string{req.Command})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to execute beacon task",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "beacon_task", "beacon", beaconID, req.Command, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"task_id": taskID,
		},
	})
}

func (h *BeaconHandler) CancelBeaconTask(c *gin.Context) {
	beaconID := c.Param("id")
	taskID := c.Param("tid")
	ctx := context.Background()

	if err := h.beaconService.CancelTask(ctx, beaconID, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to cancel beacon task",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "cancel_beacon_task", "beacon", beaconID, taskID, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "task cancelled successfully",
	})
}

func (h *BeaconHandler) GetTaskOutput(c *gin.Context) {
	beaconID := c.Param("id")
	taskID := c.Param("tid")
	ctx := context.Background()

	output, err := h.beaconService.GetTaskOutput(ctx, beaconID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to get task output",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"output": string(output),
		},
	})
}
