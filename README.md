# Multi-Window Media Sequencer

A full-stack media sequencing application that manages multiple display windows, continuously plays each window's configured playlist, supports dynamic playlist updates, and provides synchronized playback across all windows.

## Live Deployment

- Frontend: `<https://media-sequencer-frontend.onrender.com>`
- Backend: <https://multi-window-media-sequencer.onrender.com>
- Health check: <https://multi-window-media-sequencer.onrender.com/health>

Replace `<https://media-sequencer-frontend.onrender.com>` with the deployed React frontend URL.

## Features

- Multiple independent media display windows
- Continuous playlist playback
- Image, video, and blank media support
- Per-window ordered playlists with configurable durations
- Dynamic playlist additions and removals
- Persistent PostgreSQL storage
- Global synchronized playback across all windows
- Configurable synchronization duration
- Automatic return to each window's normal playlist after synchronization
- WebSocket notifications for playlist updates and synchronization events
- Dockerized local development
- Production deployment on Render

## Tech Stack

### Frontend

- React
- Vite
- CSS
- WebSocket API

### Backend

- Go 1.25
- Gin
- GORM
- PostgreSQL
- Gorilla WebSocket

### Infrastructure

- Docker
- Docker Compose
- Nginx
- Render

## Architecture

```text
React Frontend
      |
      | REST API / WebSocket
      v
Go Backend (Gin)
      |
     GORM
      |
      v
PostgreSQL
```

The frontend renders media based on playback state returned by the backend. The backend persists playlists, media, playback-cycle state, and synchronization state.

## Project Structure

```text
Multi-Window-Media-Sequencer/
├── backend/
│   ├── cmd/server/
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── service/
│   │   └── websocket/
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── assets/
│   │   ├── components/
│   │   ├── services/
│   │   ├── App.jsx
│   │   ├── App.css
│   │   └── main.jsx
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
├── docker-compose.yml
└── README.md
```

## Local Setup

### Prerequisites

- Docker Desktop and Docker Compose
- Go 1.25+ (for running the backend without Docker)
- Node.js 22+ (for running the frontend without Docker)

### Run with Docker Compose

From the project root:

```bash
docker compose up --build
```

Application URLs:

```text
Frontend: http://localhost:5173
Backend:  http://localhost:8080
Health:   http://localhost:8080/health
```

PostgreSQL is exposed on host port `5433` for local development. The frontend container is served by Nginx on port `80` internally and is mapped to host port `5173`.

Stop the application:

```bash
docker compose down
```

Remove containers and the local database volume:

```bash
docker compose down -v
```

> Removing the volume deletes local PostgreSQL data.

### Backend Without Docker

Start PostgreSQL separately, then run:

```bash
cd backend
go mod download
go run ./cmd/server
```

Configure the required environment variables before starting the backend:

```dotenv
DB_HOST=localhost
DB_PORT=5432
DB_NAME=media_sequencer
DB_USER=postgres
DB_PASSWORD=postgres
FRONTEND_URL=http://localhost:5173
GIN_MODE=debug
```

Use `GIN_MODE=release` in production. The backend also supports a `PORT` environment variable and defaults to `8080`.

### Frontend Without Docker

```bash
cd frontend
npm install
npm run dev
```

When running the frontend separately, configure the API and WebSocket URLs if they are not served through the same origin:

```dotenv
VITE_API_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080/ws
```

The frontend defaults to `/api` and derives the WebSocket URL from the current browser host when these variables are not set. Do not commit `.env` files containing secrets or environment-specific configuration.

## Database

The backend uses PostgreSQL with GORM. Database tables are automatically migrated when the server starts.

Main entities:

```text
Window
Media
PlaylistItem
SyncEvent
PlaybackCycle
ActiveSync
```

The application seeds initial media and window playlist data when the database is empty.

## Playback Logic

Each window has an independent ordered playlist. Every playlist item contains a media ID, position, and duration.

The backend maintains a persistent playback-cycle start time and calculates the current media from:

```text
cycle start + current time + playlist + item durations
```

When the configured playlist ends, playback wraps back to the beginning and continues. The assignment treats each window's total play size as a five-hour cycle; configured playlist content repeats within that cycle rather than introducing an automatic blank period.

Blank playback occurs only when blank media is explicitly configured.

## Synchronization

Synchronization is global across all windows.

Request:

```http
POST /api/sync
Content-Type: application/json
```

Example body:

```json
{
  "mediaId": 2,
  "duration": 15
}
```

The backend:

1. Validates the requested media.
2. Creates a synchronization event.
3. Stores the active synchronization state.
4. Records a common UTC start timestamp.
5. Broadcasts a `SYNC_PLAY` WebSocket event.

During synchronization, every window uses the same media and calculates its playback offset from the common start timestamp. After the configured duration expires, the active synchronization state is cleared and every window returns to its own normal playlist. The original playlists are not modified.

## WebSocket Events

WebSocket endpoint:

```text
/ws
```

### `SYNC_PLAY`

Sent when synchronized playback starts:

```json
{
  "type": "SYNC_PLAY",
  "mediaId": 2,
  "startedAt": "2026-09-10T18:18:41Z",
  "duration": 15
}
```

### `PLAYLIST_UPDATED`

Sent when a playlist changes:

```json
{
  "type": "PLAYLIST_UPDATED",
  "windowId": 2
}
```

The frontend uses these events to refresh playback and playlist state without a full page reload.

## API

### Health

```http
GET /health
```

### Windows

```http
GET /api/windows
GET /api/windows/:id
GET /api/windows/:id/current
```

### Media

```http
GET /api/media
POST /api/media
```

Example media:

```json
{
  "name": "Image 1",
  "type": "image",
  "url": "https://example.com/image.jpg",
  "duration": 30
}
```

Supported media types include images, videos, and explicitly configured blank media.

### Playlists

```http
GET /api/windows/:id/playlist
POST /api/windows/:id/playlist
DELETE /api/windows/:id/playlist/:itemId
```

Example playlist item:

```json
{
  "mediaId": 1,
  "position": 1,
  "duration": 30
}
```

### Synchronization

```http
POST /api/sync
```

See the [Synchronization](#synchronization) section for the request body.

## Seed Data

When the database is empty, the application creates:

### Media

- Image 1
- Image 2
- Video 1

### Windows

- Window 1
- Window 2
- Window 3

Each window starts with a different playlist configuration so independent playback and global synchronization can be demonstrated immediately.

## Deployment

The production deployment uses:

```text
Render
├── React Frontend
├── Go Backend
└── PostgreSQL
```

Example backend environment variables:

```dotenv
DB_HOST=<render-internal-host>
DB_PORT=5432
DB_NAME=media_sequencer
DB_USER=<render-db-user>
DB_PASSWORD=<render-db-password>
FRONTEND_URL=<deployed-frontend-url>
GIN_MODE=release
```

Example frontend environment variables:

```dotenv
VITE_API_URL=https://multi-window-media-sequencer.onrender.com/api
VITE_WS_URL=wss://multi-window-media-sequencer.onrender.com/ws
```

The production frontend uses HTTPS for REST requests and secure WebSockets (`wss`).

## Docker Details

### Backend

The backend uses a multi-stage Docker build:

```text
Go builder
  -> download dependencies
  -> build binary
  -> minimal Alpine runtime
```

### Frontend

The frontend uses:

```text
Node builder
  -> npm run build
  -> Nginx static server
```

Docker Compose runs the `frontend`, `backend`, and `postgres` services.

## Assumptions

- Media URLs point to externally accessible media resources.
- Media duration is configured as part of media or playlist configuration.
- Playback state is calculated from timestamps rather than requiring the backend to continuously push the next media item.
- Synchronization is global across all configured windows.
- A synchronization event temporarily overrides normal playback without modifying playlists.
- PostgreSQL is the source of truth for persistent application state.
- Blank playback is explicit and is not automatically inserted into unused cycle time.
- The browser is responsible for media rendering and playback between state updates.
- Videos are muted to improve reliable browser autoplay behavior.

## Verification Checklist

- Three windows load successfully.
- Each window follows its configured playlist.
- Media transitions occur continuously.
- Playlist items can be added and removed dynamically.
- Playlist changes are reflected without restarting the backend.
- Sync causes the same media to appear on all windows.
- All windows use the same synchronization start time.
- Sync expires after the configured duration.
- Each window resumes its own playlist after sync.
- Data persists after application restart.
- `/health` returns successfully.
- The production frontend communicates with the backend.
- The production WebSocket connection works.

## Author

**Aditya Sanskar Srivastav**

Final-year Computer Science and Engineering student.

GitHub: <https://github.com/Aditya7880900936>