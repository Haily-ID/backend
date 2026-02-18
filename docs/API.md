# API Endpoints Documentation

## Authentication

### Register

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123456789",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user"
  }
}
```

## User Endpoints (Protected)

### Get Profile

```http
GET /api/v1/users/me
Authorization: Bearer {token}
```

### Update Profile

```http
PUT /api/v1/users/me
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "John Updated"
}
```

### Get User by ID

```http
GET /api/v1/users/:id
Authorization: Bearer {token}
```

### List Users

```http
GET /api/v1/users?offset=0&limit=10
Authorization: Bearer {token}
```

## Company Endpoints (Protected)

### Create Company

```http
POST /api/v1/companies
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "My Company",
  "description": "Company description"
}
```

### Get My Companies (Owner)

```http
GET /api/v1/companies/mine
Authorization: Bearer {token}
```

### Get Joined Companies (Member)

```http
GET /api/v1/companies/joined
Authorization: Bearer {token}
```

### Get Company by ID

```http
GET /api/v1/companies/:id
Authorization: Bearer {token}
```

### Update Company

```http
PUT /api/v1/companies/:id
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Updated Company Name",
  "description": "Updated description"
}
```

### Delete Company

```http
DELETE /api/v1/companies/:id
Authorization: Bearer {token}
```

### Join Company

```http
POST /api/v1/companies/:id/join
Authorization: Bearer {token}
```

### Leave Company

```http
DELETE /api/v1/companies/:id/leave
Authorization: Bearer {token}
```

### Get Company Members

```http
GET /api/v1/companies/:id/members
Authorization: Bearer {token}
```

### List All Companies

```http
GET /api/v1/companies?offset=0&limit=10
Authorization: Bearer {token}
```

## Health Check

```http
GET /health
```

Response:

```json
{
  "status": "healthy",
  "database": "up",
  "redis": "up",
  "version": "1.0.0"
}
```

## Response Format

All responses follow an unnormalized format:

### Success Response

```json
{
  "instructions": [
    {
      "type": "user",
      "ids": ["123456789"]
    }
  ],
  "entities": {
    "user": {
      "123456789": {
        "id": "123456789",
        "email": "user@example.com",
        "name": "John Doe",
        "role": "user",
        "created_at": 1708291200,
        "updated_at": 1708291200
      }
    }
  }
}
```

### Error Response

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "validation failed on field 'Email': required",
    "field": "Email"
  }
}
```

## Error Codes

- `VALIDATION_ERROR` - Input validation failed
- `NOT_FOUND` - Resource not found
- `UNAUTHORIZED` - Authentication required
- `FORBIDDEN` - Insufficient permissions
- `CONFLICT` - Resource already exists
- `INTERNAL_SERVER_ERROR` - Server error
- `BAD_REQUEST` - Invalid request
