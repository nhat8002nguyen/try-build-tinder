# Spark API Documentation

Complete API reference for the Spark dating application backend.

## Base URL

```
http://localhost:8080
```

## Authentication

Most endpoints require JWT authentication. Include the access token in the Authorization header:

```
Authorization: Bearer <access_token>
```

Tokens are obtained through registration or login and should be refreshed before expiration using the refresh token.

---

## Table of Contents

1. [Authentication](#authentication-endpoints)
2. [Users](#users)
3. [Discovery](#discovery)
4. [Swipes](#swipes)
5. [Matches](#matches)
6. [Messages](#messages)
7. [Notifications](#notifications)
8. [WebSocket](#websocket)

---

## Authentication Endpoints

### Register

Create a new user account.

**Endpoint:** `POST /api/auth/register`

**Authentication:** None required

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "SecurePassword123",
  "name": "John Doe"
}
```

**Validation:**
- `email`: Valid email format (required)
- `password`: Minimum 8 characters (required)
- `name`: Required

**Success Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john.doe@example.com",
      "name": "John Doe",
      "gender": "",
      "birthdate": null,
      "bio": "",
      "location": {
        "latitude": 0,
        "longitude": 0
      },
      "preferences": {
        "min_age": 0,
        "max_age": 0,
        "max_distance": 0,
        "gender_preference": []
      },
      "is_verified": false,
      "is_active": true,
      "last_active_at": "2026-01-22T10:00:00Z",
      "created_at": "2026-01-22T10:00:00Z",
      "updated_at": "2026-01-22T10:00:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 86400
    }
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input (duplicate email, weak password, etc.)
- `500 Internal Server Error`: Server error

---

### Login

Authenticate with existing credentials.

**Endpoint:** `POST /api/auth/login`

**Authentication:** None required

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "SecurePassword123"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john.doe@example.com",
      "name": "John Doe",
      ...
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 86400
    }
  }
}
```

**Error Responses:**
- `401 Unauthorized`: Invalid credentials
- `400 Bad Request`: Missing required fields

---

### Refresh Token

Obtain a new access token using a refresh token.

**Endpoint:** `POST /api/auth/refresh`

**Authentication:** None required

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400
  }
}
```

**Error Responses:**
- `401 Unauthorized`: Invalid or expired refresh token

---

### Get Current User

Get the authenticated user's profile.

**Endpoint:** `GET /api/auth/me`

**Authentication:** Required (Bearer token)

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john.doe@example.com",
    "name": "John Doe",
    "gender": "male",
    "birthdate": "1995-05-15T00:00:00Z",
    "bio": "Love hiking and coffee!",
    "location": {
      "latitude": 37.7749,
      "longitude": -122.4194
    },
    "preferences": {
      "min_age": 22,
      "max_age": 35,
      "max_distance": 50,
      "gender_preference": ["female"]
    },
    "is_verified": false,
    "is_active": true,
    "last_active_at": "2026-01-22T10:00:00Z",
    "created_at": "2026-01-22T09:00:00Z",
    "updated_at": "2026-01-22T10:00:00Z",
    "photos": [
      {
        "id": "photo-uuid-1",
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
        "photo_url": "/uploads/550e8400-e29b-41d4-a716-446655440000/photo1.jpg",
        "display_order": 0,
        "is_approved": true,
        "created_at": "2026-01-22T09:30:00Z"
      }
    ]
  }
}
```

**Error Responses:**
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: User not found

---

### OAuth Authentication

Authenticate using OAuth providers (Google, Facebook).

**Endpoint:** `GET /api/auth/oauth/:provider`

**Parameters:**
- `provider`: `google` or `facebook`

**Flow:**
1. Redirect user to this endpoint
2. User authenticates with OAuth provider
3. OAuth provider redirects to callback
4. Frontend receives tokens via query parameters

**OAuth Callback Endpoint:** `GET /api/auth/oauth/:provider/callback`

The callback redirects to: `http://localhost:3000/auth/callback?access_token=...&refresh_token=...`

---

## Users

### Get User by ID

Retrieve a user's public profile.

**Endpoint:** `GET /api/users/:id`

**Authentication:** Required

**Parameters:**
- `id`: User UUID

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "gender": "male",
    "birthdate": "1995-05-15T00:00:00Z",
    "bio": "Love hiking and coffee!",
    "photos": [
      {
        "id": "photo-uuid-1",
        "photo_url": "/uploads/550e8400-e29b-41d4-a716-446655440000/photo1.jpg",
        "display_order": 0
      }
    ],
    "is_verified": false
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid user ID format
- `404 Not Found`: User not found

---

### Update Profile

Update the current user's profile.

**Endpoint:** `PUT /api/users/me`

**Authentication:** Required

**Request Body:**
```json
{
  "name": "John Doe",
  "gender": "male",
  "birthdate": "1995-05-15T00:00:00Z",
  "bio": "Love hiking and coffee!",
  "preferences": {
    "min_age": 22,
    "max_age": 35,
    "max_distance": 50,
    "gender_preference": ["female"]
  }
}
```

**Field Validations:**
- `name`: Optional string
- `gender`: One of: `male`, `female`, `other`
- `birthdate`: ISO 8601 date string
- `bio`: Optional string
- `preferences.min_age`: Integer >= 18
- `preferences.max_age`: Integer <= 100
- `preferences.max_distance`: Integer (kilometers)
- `preferences.gender_preference`: Array of gender values

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john.doe@example.com",
    "name": "John Doe",
    ...
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input
- `500 Internal Server Error`: Update failed

---

### Upload Photo

Upload a profile photo.

**Endpoint:** `POST /api/users/me/photos`

**Authentication:** Required

**Request:** `multipart/form-data`
- `photo`: Image file (JPEG, PNG, GIF, WebP)

**Constraints:**
- Maximum 6 photos per user
- Maximum file size: 5MB
- Supported formats: JPEG, PNG, GIF, WebP

**Success Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "photo-uuid-1",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "photo_url": "/uploads/550e8400-e29b-41d4-a716-446655440000/photo1.jpg",
    "display_order": 0,
    "is_approved": true,
    "created_at": "2026-01-22T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid file, exceeded photo limit, or file too large
- `500 Internal Server Error`: Upload failed

---

### Delete Photo

Delete a profile photo.

**Endpoint:** `DELETE /api/users/me/photos/:photoId`

**Authentication:** Required

**Parameters:**
- `photoId`: Photo UUID

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Photo deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid photo ID
- `500 Internal Server Error`: Deletion failed

---

### Update Location

Update user's current location.

**Endpoint:** `PUT /api/users/me/location`

**Authentication:** Required

**Request Body:**
```json
{
  "latitude": 37.7749,
  "longitude": -122.4194
}
```

**Validation:**
- `latitude`: -90 to 90
- `longitude`: -180 to 180

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Location updated successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid coordinates
- `500 Internal Server Error`: Update failed

---

## Discovery

### Get Potential Matches

Get a list of potential matches based on user preferences and location.

**Endpoint:** `GET /api/discover`

**Authentication:** Required

**Query Parameters:**
- `limit` (optional): Number of results, max 50 (default: system default)
- `min_age` (optional): Minimum age filter (>= 18)
- `max_age` (optional): Maximum age filter (<= 100)
- `max_distance` (optional): Maximum distance in kilometers
- `gender_preference` (optional): Comma-separated genders (e.g., `male,female`)

**Example Request:**
```
GET /api/discover?limit=10&min_age=22&max_age=35&max_distance=50&gender_preference=female
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "user-uuid-1",
      "name": "Jane Smith",
      "age": 28,
      "gender": "female",
      "bio": "Adventure seeker and foodie",
      "distance": 5.2,
      "photos": [
        {
          "id": "photo-uuid-1",
          "photo_url": "/uploads/user-uuid-1/photo1.jpg",
          "display_order": 0
        }
      ],
      "is_verified": true
    },
    ...
  ]
}
```

**Algorithm:**
The discovery system uses a weighted scoring algorithm:
- Activity score: Recent activity is prioritized
- Distance score: Closer users rank higher
- Profile score: Complete profiles (photos, bio, verification) rank higher

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch matches

---

## Swipes

### Create Swipe

Record a swipe (like or dislike) on another user.

**Endpoint:** `POST /api/swipes`

**Authentication:** Required

**Request Body:**
```json
{
  "target_id": "550e8400-e29b-41d4-a716-446655440000",
  "direction": "like"
}
```

**Field Validations:**
- `target_id`: Valid user UUID (required)
- `direction`: Must be `like` or `dislike` (required)

**Success Response (201 Created):**

**Case 1: No Match**
```json
{
  "success": true,
  "data": {
    "is_match": false
  }
}
```

**Case 2: Mutual Match**
```json
{
  "success": true,
  "data": {
    "is_match": true,
    "match": {
      "id": "match-uuid-1",
      "user1_id": "current-user-uuid",
      "user2_id": "550e8400-e29b-41d4-a716-446655440000",
      "matched_at": "2026-01-22T10:00:00Z",
      "last_message_at": null
    }
  }
}
```

**Notes:**
- A match occurs when both users swipe right (like) on each other
- Users cannot swipe on themselves
- Duplicate swipes are prevented by the system
- When a match occurs, both users receive a WebSocket notification

**Error Responses:**
- `400 Bad Request`: Invalid target ID or trying to swipe on yourself
- `500 Internal Server Error`: Failed to record swipe

---

## Matches

### Get All Matches

Retrieve all matches for the current user.

**Endpoint:** `GET /api/matches`

**Authentication:** Required

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "match-uuid-1",
      "other_user": {
        "id": "user-uuid-1",
        "name": "Jane Smith",
        "age": 28,
        "gender": "female",
        "bio": "Adventure seeker",
        "photos": [
          {
            "id": "photo-uuid-1",
            "photo_url": "/uploads/user-uuid-1/photo1.jpg",
            "display_order": 0
          }
        ]
      },
      "matched_at": "2026-01-22T10:00:00Z",
      "last_message_at": "2026-01-22T11:30:00Z"
    },
    ...
  ]
}
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch matches

---

### Get Match Details

Get details of a specific match.

**Endpoint:** `GET /api/matches/:id`

**Authentication:** Required

**Parameters:**
- `id`: Match UUID

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "match-uuid-1",
    "other_user": {
      "id": "user-uuid-1",
      "name": "Jane Smith",
      "age": 28,
      "gender": "female",
      "bio": "Adventure seeker",
      "photos": [...]
    },
    "matched_at": "2026-01-22T10:00:00Z",
    "last_message_at": "2026-01-22T11:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid match ID format
- `403 Forbidden`: User is not part of this match
- `404 Not Found`: Match not found

---

## Messages

### Get Messages

Retrieve messages from a specific match.

**Endpoint:** `GET /api/matches/:id/messages`

**Authentication:** Required

**Parameters:**
- `id`: Match UUID

**Query Parameters:**
- `limit` (optional): Number of messages, max 100 (default: 50)
- `offset` (optional): Offset for pagination (default: 0)

**Example Request:**
```
GET /api/matches/match-uuid-1/messages?limit=50&offset=0
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "message-uuid-1",
      "match_id": "match-uuid-1",
      "sender_id": "user-uuid-1",
      "content": "Hey! How are you?",
      "is_read": true,
      "created_at": "2026-01-22T11:30:00Z"
    },
    {
      "id": "message-uuid-2",
      "match_id": "match-uuid-1",
      "sender_id": "current-user-uuid",
      "content": "I'm great, thanks! How about you?",
      "is_read": true,
      "created_at": "2026-01-22T11:32:00Z"
    },
    ...
  ]
}
```

**Notes:**
- Messages are returned in chronological order (oldest first)
- Calling this endpoint automatically marks all messages as read for the current user
- Use pagination for long conversations

**Error Responses:**
- `400 Bad Request`: Invalid match ID
- `403 Forbidden`: User is not part of this match
- `500 Internal Server Error`: Failed to fetch messages

---

### Send Message

Send a message to a match.

**Endpoint:** `POST /api/matches/:id/messages`

**Authentication:** Required

**Parameters:**
- `id`: Match UUID

**Request Body:**
```json
{
  "content": "Hey! How are you?"
}
```

**Validation:**
- `content`: Required, 1-1000 characters

**Success Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "message-uuid-1",
    "match_id": "match-uuid-1",
    "sender_id": "current-user-uuid",
    "content": "Hey! How are you?",
    "is_read": false,
    "created_at": "2026-01-22T11:30:00Z"
  }
}
```

**Notes:**
- The recipient receives a real-time notification via WebSocket
- The match's `last_message_at` timestamp is updated
- Message content is stored as plain text

**Error Responses:**
- `400 Bad Request`: Invalid match ID or message content
- `403 Forbidden`: User is not part of this match
- `500 Internal Server Error`: Failed to send message

---

## Notifications

### Get Notifications

Retrieve user notifications.

**Endpoint:** `GET /api/notifications`

**Authentication:** Required

**Query Parameters:**
- `limit` (optional): Number of notifications, max 50 (default: 20)
- `offset` (optional): Offset for pagination (default: 0)

**Example Request:**
```
GET /api/notifications?limit=20&offset=0
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "notifications": [
      {
        "id": "notif-uuid-1",
        "user_id": "current-user-uuid",
        "type": "match",
        "title": "New Match!",
        "message": "You matched with Jane Smith",
        "data": {
          "match_id": "match-uuid-1",
          "user_id": "user-uuid-1"
        },
        "is_read": false,
        "created_at": "2026-01-22T10:00:00Z"
      },
      {
        "id": "notif-uuid-2",
        "user_id": "current-user-uuid",
        "type": "message",
        "title": "New Message",
        "message": "Jane Smith sent you a message",
        "data": {
          "match_id": "match-uuid-1",
          "message_id": "message-uuid-1"
        },
        "is_read": true,
        "created_at": "2026-01-22T11:30:00Z"
      }
    ],
    "unread_count": 1
  }
}
```

**Notification Types:**
- `match`: New match notification
- `message`: New message notification
- `like`: Someone liked you (premium feature)

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch notifications

---

### Mark Notification as Read

Mark a specific notification as read.

**Endpoint:** `PUT /api/notifications/:id/read`

**Authentication:** Required

**Parameters:**
- `id`: Notification UUID

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Notification marked as read"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid notification ID
- `500 Internal Server Error`: Failed to update notification

---

## WebSocket

### WebSocket Connection

Establish a real-time WebSocket connection for instant notifications.

**Endpoint:** `WS /ws`

**Authentication:** Required (via query parameter)

**Connection URL:**
```
ws://localhost:8080/ws?token=<access_token>
```

**Message Types:**

**1. Match Notification**
```json
{
  "type": "match",
  "payload": {
    "match_id": "match-uuid-1",
    "other_user_id": "user-uuid-1"
  }
}
```

**2. New Message Notification**
```json
{
  "type": "message",
  "payload": {
    "message": {
      "id": "message-uuid-1",
      "match_id": "match-uuid-1",
      "sender_id": "user-uuid-1",
      "content": "Hey! How are you?",
      "is_read": false,
      "created_at": "2026-01-22T11:30:00Z"
    },
    "match_id": "match-uuid-1"
  }
}
```

**Client Implementation Example (JavaScript):**
```javascript
const ws = new WebSocket(`ws://localhost:8080/ws?token=${accessToken}`);

ws.onopen = () => {
  console.log('WebSocket connected');
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'match') {
    console.log('New match!', data.payload);
    // Show match notification
  } else if (data.type === 'message') {
    console.log('New message!', data.payload);
    // Update chat UI
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket disconnected');
  // Implement reconnection logic
};
```

**Notes:**
- WebSocket connections are maintained per user
- Multiple connections from the same user are supported (e.g., multiple devices)
- Connections automatically close when the token expires
- Implement exponential backoff for reconnection attempts

---

## Error Response Format

All error responses follow this structure:

```json
{
  "success": false,
  "error": "Error message description"
}
```

**Common HTTP Status Codes:**
- `200 OK`: Request succeeded
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters or body
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Authenticated but not authorized for this resource
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **General endpoints**: 100 requests per minute per IP
- **Authentication endpoints**: 10 requests per minute per IP
- **Message sending**: 30 messages per minute per user

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642857600
```

---

## Testing with Postman

1. **Import the Collection:**
   - Download `Spark_API.postman_collection.json`
   - Import into Postman: File → Import

2. **Set Variables:**
   - `baseUrl`: `http://localhost:8080` (default)
   - `accessToken`: Auto-populated after login/register
   - `refreshToken`: Auto-populated after login/register
   - `userId`: Auto-populated after login/register

3. **Workflow:**
   - Start with "Register" or "Login" to get tokens
   - Tokens are automatically saved and used in subsequent requests
   - Use "Get Current User" to verify authentication
   - Explore other endpoints with authentication

4. **Tips:**
   - The collection includes test scripts that auto-save tokens
   - Update profile and location before testing discovery
   - Upload photos for better testing experience
   - Use multiple accounts to test swipes and matches

---

## Development Notes

- **CORS**: Configured to allow `http://localhost:3000` (frontend)
- **File Uploads**: Stored locally in `./uploads` directory (configurable to S3)
- **Database**: PostgreSQL with automatic migrations on startup
- **Cache**: Redis for session management and rate limiting
- **Logging**: Request logging via middleware

---

## Support

For issues or questions:
- Check the main [README.md](./README.md)
- Review [TESTING.md](./TESTING.md) for testing guidelines
- Create an issue in the repository
