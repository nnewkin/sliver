# Sliver Web UI

**合规与免责声明**：本项目仅面向获得合法书面授权的网络安全从业者、红队渗透测试人员，用于授权范围内的安全评估、攻防演练与技术研究。严禁将本项目用于任何未授权的非法入侵、攻击、窃密、破坏等行为，任何违反国家及地区法律法规的使用行为，全部法律后果由使用者自行承担，开发者与内容输出方不承担任何直接或间接的法律责任。使用本项目前，你必须确保已获得目标系统与网络的书面合法授权，严格遵守当地网络安全相关法律法规。

---

Web-based management interface for [Sliver](https://github.com/BishopFox/sliver) C2 framework.

## Features

- **100% Sliver API Coverage**: All Sliver native gRPC APIs are exposed via Web UI
- **Real-time Communication**: WebSocket-based event streaming for sessions and beacons
- **Telegram Bot Integration**: Full-featured Telegram bot for remote operations
- **RBAC Permissions**: Role-based access control (Super Admin, Operator, Auditor)
- **TOTP 2FA**: Two-factor authentication for enhanced security
- **Audit Logging**: Complete operation logging and tracking

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Sliver Web UI                             │
├─────────────────────────────────────────────────────────────┤
│  Frontend (React + Ant Design)     │  Backend (Go + Gin)   │
│  - Dashboard                       │  - REST API           │
│  - Session Management              │  - WebSocket Server   │
│  - Beacon Control                  │  - Telegram Bot       │
│  - Listener Management             │  - gRPC Client        │
│  - Implant Generation              │                       │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Sliver Server                             │
│                    (gRPC :50050)                             │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

### Using Docker Compose

1. Clone the repository
2. Generate a secure JWT secret:
   ```bash
   export JWT_SECRET=$(openssl rand -base64 32)
   ```
3. Start services:
   ```bash
   docker-compose up -d
   ```
4. Access the Web UI at `http://localhost:8080`
5. Login with default credentials (create admin on first login)

### Manual Build

#### Backend

```bash
cd sliver-web
go mod download
go build -o sliver-web ./cmd/server
./sliver-web --config sliver-web.conf
```

#### Frontend

```bash
cd web
npm install
npm run dev
```

## Configuration

Copy `sliver-web.conf.example` to `sliver-web.conf` and configure:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  jwt_secret: "CHANGE_ME"
  jwt_expiry_hours: 24

sliver:
  grpc_addr: "localhost:50050"
  certs_dir: "./certs"

database:
  path: "./data/sliver-web.db"

telegram:
  enabled: false
  bot_token: ""
  mode: "polling"
```

## Telegram Bot Setup

1. Create a bot via [@BotFather](https://t.me/BotFather)
2. Get your bot token
3. Configure in Settings > Telegram Bot
4. Add your Telegram User ID to the whitelist

### Bot Commands

| Command | Description |
|---------|-------------|
| `/start` | Start the bot |
| `/sessions` | List active sessions |
| `/beacons` | List all beacons |
| `/exec <session_id> <cmd>` | Execute command |
| `/alerts` | View alert history |

## API Documentation

### Authentication

```
POST /api/auth/login
POST /api/auth/2fa/verify
POST /api/auth/logout
GET  /api/auth/profile
```

### Sessions

```
GET  /api/sessions
GET  /api/sessions/:id
DELETE /api/sessions/:id
POST /api/sessions/:id/exec
GET  /api/sessions/:id/processes
```

### Beacons

```
GET  /api/beacons
GET  /api/beacons/:id
GET  /api/beacons/:id/tasks
POST /api/beacons/:id/task
```

### Listeners

```
GET  /api/listeners
POST /api/listeners/mtls
POST /api/listeners/dns
POST /api/listeners/http
DELETE /api/listeners/:id
```

### Generate

```
POST /api/generate
GET  /api/generate/download/:name
```

## Security

- JWT-based authentication
- TOTP 2FA support
- RBAC permission model
- Complete audit logging
- Sensitive data encryption

## License

See [LICENSE](LICENSE) file.
