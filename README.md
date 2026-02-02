# Spark - Tinder Clone

A full-stack dating application built with Go, PostgreSQL, and React. Designed to be scalable and reliable for 10,000+ active users.

## Features

- **User Authentication**: Email/password and OAuth (Google, Facebook)
- **User Profiles**: Photos, bio, preferences, location
- **Discovery**: Geolocation-based matching with smart ranking algorithm
- **Swipe Mechanism**: Like/dislike with match detection
- **Real-time Messaging**: WebSocket-based chat with polling fallback
- **Notifications**: Match and message notifications

## Tech Stack

### Backend
- **Go** with Gin framework
- **PostgreSQL** with GORM
- **Redis** for caching and sessions
- **JWT** for authentication
- **Gorilla WebSocket** for real-time communication

### Frontend
- **React 18** with TypeScript
- **Vite** for build tooling
- **TailwindCSS** for styling
- **Framer Motion** for animations
- **React Query** for server state
- **Zustand** for client state

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 20+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repo-url>
   cd try-build-tinder
   ```

2. **Start infrastructure services**
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

3. **Setup backend**
   ```bash
   cd backend
   
   # Copy environment template and configure it
   cp .env.example .env
   # Edit .env and set your JWT_SECRET at minimum
   
   # Install dependencies and run
   go mod tidy
   go run cmd/server/main.go
   ```

4. **Setup frontend**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

5. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080/api

### Environment Variables

The backend includes a `.env.example` template file. Copy it to `.env` and configure:

```bash
cd backend
cp .env.example .env
```

**Required Configuration:**
- `JWT_SECRET`: Change to a secure random string (required for authentication)
- `DATABASE_URL`: Default works with docker-compose.dev.yml
- `REDIS_URL`: Default works with docker-compose.dev.yml

**Optional Configuration:**
- OAuth credentials (Google/Facebook) - only needed if using OAuth login
- S3 storage - only needed if using cloud storage instead of local

Example `.env` contents:

```env
# Server
SERVER_PORT=8080
ENVIRONMENT=development

# Database (works with docker-compose.dev.yml)
DATABASE_URL=postgres://postgres:postgres@localhost:5432/tinder_clone?sslmode=disable

# Redis (works with docker-compose.dev.yml)
REDIS_URL=redis://localhost:6379

# JWT - CHANGE THIS!
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRE_HOURS=24

# OAuth (optional - leave empty to skip OAuth)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/oauth/google/callback

FACEBOOK_CLIENT_ID=
FACEBOOK_CLIENT_SECRET=
FACEBOOK_REDIRECT_URL=http://localhost:8080/api/auth/oauth/facebook/callback

# Storage
STORAGE_TYPE=local
LOCAL_STORAGE_DIR=./uploads
```

### Production Deployment

```bash
# Build and run with Docker Compose
docker-compose up -d --build

# Access at http://localhost
```

## Testing

### Backend Tests

The backend has comprehensive unit tests for services, utilities, and middleware.

```bash
cd backend

# Run all tests
go test ./... -v

# Run specific package
go test ./internal/utils/... -v

# Run with coverage
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Current Test Coverage:**
- ✅ **Utils (validators): ALL TESTS PASSING** (4 test suites, 13 test cases)

### Frontend Tests

The frontend uses Vitest for component and integration testing.

```bash
cd frontend

# Run tests in watch mode
npm run test

# Run tests once
npm run test:run

# Run with UI
npm run test:ui
```

**Current Test Coverage:**
- ✅ Components (Landing): Test infrastructure ready
- ✅ Store (auth): Test infrastructure ready

See [TESTING.md](./TESTING.md) for detailed testing documentation.

## API Documentation

For complete API documentation with request/response examples, authentication details, and error handling:

📖 **[View Full API Documentation](./API_DOCUMENTATION.md)**

🧪 **[Download Postman Collection](./Spark_API.postman_collection.json)**

### Quick Reference

**Authentication:**
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login with email/password
- `POST /api/auth/refresh` - Refresh JWT token
- `GET /api/auth/oauth/:provider` - OAuth login (Google, Facebook)
- `GET /api/auth/me` - Get current user

**Users:**
- `GET /api/users/:id` - Get user profile
- `PUT /api/users/me` - Update profile
- `POST /api/users/me/photos` - Upload photo (max 6 photos)
- `DELETE /api/users/me/photos/:photoId` - Delete photo
- `PUT /api/users/me/location` - Update location

**Discovery:**
- `GET /api/discover` - Get potential matches (with filters)

**Swipes:**
- `POST /api/swipes` - Record swipe (like/dislike)

**Matches:**
- `GET /api/matches` - Get all matches
- `GET /api/matches/:id` - Get match details

**Messages:**
- `GET /api/matches/:id/messages` - Get messages
- `POST /api/matches/:id/messages` - Send message
- `WS /ws` - WebSocket connection for real-time updates

**Notifications:**
- `GET /api/notifications` - Get notifications
- `PUT /api/notifications/:id/read` - Mark as read

### Testing with Postman

1. Import `Spark_API.postman_collection.json` into Postman
2. Start with "Register" or "Login" to get authentication tokens
3. Tokens are automatically saved to collection variables
4. All authenticated endpoints will use the saved token

## Architecture

```
┌─────────────────┐
│  React Frontend │ (SPA with WebSocket client)
└────────┬────────┘
         │ HTTP/WebSocket
┌────────▼────────────────────────────────────┐
│  Golang API Server                          │
│  - REST API (Gin framework)                 │
│  - WebSocket Hub (Gorilla WebSocket)        │
│  - Auth Service (JWT + OAuth)               │
│  - Matching Engine                          │
│  - Message Service                          │
└────────┬────────────────────────────────────┘
         │
    ┌────┴────┬──────────┬──────────┐
    │         │          │          │
┌───▼───┐ ┌──▼───┐ ┌───▼───┐ ┌───▼───┐
│Postgres│ │Redis │ │  S3   │ │Queue  │
│  DB    │ │Cache │ │Photos │ │(Redis)│
└────────┘ └──────┘ └───────┘ └───────┘
```

## Scalability Features

- **Database**: Connection pooling, indexed queries, optimized for geospatial queries
- **Caching**: Redis for sessions, hot data, and rate limiting
- **WebSocket**: Hub pattern with connection registry for efficient broadcasting
- **Backend**: Stateless design allowing horizontal scaling
- **Frontend**: SPA with code splitting and lazy loading

## Matching Algorithm

The discovery system uses a scoring function:

```
score = w1 * activity_score + w2 * distance_score + w3 * profile_score

- activity_score: Based on last active time (recent = higher)
- distance_score: Inverse of distance (closer = higher)
- profile_score: Profile completeness (photos, bio, verification)
```

## Project Structure

```
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── config/          # Configuration
│   │   ├── database/        # Database connection
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # Auth middleware
│   │   ├── models/          # Data models
│   │   ├── services/        # Business logic
│   │   ├── utils/           # Utilities
│   │   └── websocket/       # WebSocket hub
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── contexts/        # React contexts
│   │   ├── pages/           # Page components
│   │   ├── services/        # API clients
│   │   ├── store/           # State management
│   │   └── types/           # TypeScript types
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml       # Production setup
├── docker-compose.dev.yml   # Development setup
└── README.md
```

## License

MIT
