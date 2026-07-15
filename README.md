# Pizza Tracker

Real-time pizza order tracking system built with Go and Gin.

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)

## Setup

1. Install dependencies:
```bash
go mod download
```

3. Create `.env` file:
```env
PORT="3033"
DB_PATH="./data/pizza.db"
SESSION_SECRET="your-secret-key-here"
```

## Run

```bash
go run cmd/main.go
```

The server will start at `http://localhost:3033`

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `PORT` | Server port | `8080` |
| `DB_PATH` | SQLite database path | (required) |
| `SESSION_SECRET` | Secret key for session encryption | (required) |

## Project Structure

```
cmd/                    # Application entry point
internal/
├── admin/              # Admin handlers and routes
├── app/                # App initialization
├── config/             # Configuration loading
├── database/           # Database setup
├── middleware/          # Auth middleware
├── order/              # Order handlers, models, repository
├── route/              # Route setup
├── session/            # Session store
├── shared/
│   ├── notification/   # Real-time notifications (SSE)
│   └── util/           # Shared utilities
└── user/               # User handlers, models, repository
templates/              # HTML templates
```
