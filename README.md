# go-cli-chat

A simple command-line chat application written in Go using WebSockets. This project includes both a server and a client that communicate over WebSocket connections. Multiple clients can connect and chat with each other in real time through their terminals.

---

## 📦 Features

- Real-time messaging via WebSockets
- Multiple clients supported concurrently
- Messages are broadcasted to all connected clients (except the sender)
- Simple CLI interface for input and display
- Built with Go’s goroutines, channels, and mutex for concurrency

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/go-cli-chat.git
cd go-cli-chat
```
