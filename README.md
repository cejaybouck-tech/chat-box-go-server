# Chat Box Go Server

A minimal, real-time chat server in Go using WebSockets. Connect with React or Vue clients — great for learning concurrency, WebSockets, and full‑stack integration.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev)
[![WebSockets](https://img.shields.io/badge/WebSockets-gorilla%2Fwebsocket-blue)](https://github.com/gorilla/websocket)

---

## Features

- Real-time chat over WebSockets
- Concurrent clients using goroutines
- Simple in-memory state with on-disk user registry (`users.json`)
- Works with any frontend (React, Vue, etc.)

## Tech Stack

| Part     | Tech                                      |
|----------|-------------------------------------------|
| Backend  | Go + `gorilla/websocket`                   |
| Frontend | React + Vite, Vue 3 + Vite                 |
| State    | In-memory; users persisted to `users.json` |

---

## Getting Started

### 1) Clone & run

```bash
git clone https://github.com/cejaybouck-tech/chat-box-go-server.git
cd chat-box-go-server/server
go run .
```

### 2) Configure environment (optional)

Create a `.env` file alongside the server binary:

```dotenv
APP_ORIGIN="http://localhost:5173"  # Allowed client origin (CORS/WebSocket upgrade)
PORT="8080"                          # HTTP port the server listens on
```

The server will create `users.json` on first run if it does not exist. If a user is missing, the server may create it on-demand.

### 3) WebSocket endpoint

- URL: `ws://<host>:<PORT>/chat`

---

## WebSocket Protocol

All messages are JSON. Below are the commonly used payloads.

### Authenticate (client → server)

```jsonc
{
  "type": "authenticate",
  "username": "alice",
  "password": "secret"
}
```

### Auth response (server → client)

```jsonc
{
  "type": "auth_response",
  "success": false,
  "message": "invalid credentials"
}
```

### Users list (server → client)

```jsonc
{
  "type": "users",
  "users": ["alice", "bob"]
}
```

### Presence events (server → clients)

```jsonc
{ "type": "user_joined", "username": "alice", "message": "" }
{ "type": "user_left",   "username": "alice", "message": "" }
```

### Chat message

```jsonc
{ "type": "message", "username": "alice", "message": "Hello" }
```

---

## Clients

- React client: https://github.com/cejaybouck-tech/chat-box-client-react
- Vue client: https://github.com/cejaybouck-tech/chat-box-client-vue

---

## Build

```bash
go build -o chat-box-server
./chat-box-server
```

## License

MIT
