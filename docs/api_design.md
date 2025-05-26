# API Design

## Authentication
All protected endpoints require JWT Bearer token in Authorization header:
```
Authorization: Bearer <jwt-token>
```

---

## Authentication Endpoints

### Register User
**POST** `/api/auth/register`

Creates a new user account with optional role assignment.

**Request Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
    "email": "teacher@example.com",
    "name": "John Teacher",
    "password": "securepassword123",
    "role_id": "550e8400-e29b-41d4-a716-446655440000"  // Optional
}
```

**Response Codes:**
- `201 Created` - User successfully created
- `400 Bad Request` - Invalid input data
- `409 Conflict` - Email already exists

**Success Response (201):**
```json
{
    "success": true,
    "message": "User created successfully",
    "data": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "email": "teacher@example.com",
        "name": "John Teacher",
        "created_at": "2024-01-20T10:00:00Z",
        "updated_at": "2024-01-20T10:00:00Z"
    }
}
```

**Error Response (400):**
```json
{
    "success": false,
    "message": "An error occurred",
    "error": "password: must be at least 6 characters"
}
```

---

### Login
**POST** `/api/auth/login`

Authenticates user and returns JWT token.

**Request Body:**
```json
{
    "email": "teacher@example.com",
    "password": "securepassword123"
}
```

**Response Codes:**
- `200 OK` - Successfully authenticated
- `401 Unauthorized` - Invalid credentials
- `400 Bad Request` - Missing required fields

**Success Response (200):**
```json
{
    "success": true,
    "message": "Login successful",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzZTQ1NjctZTg5Yi0xMmQzLWE0NTYtNDI2NjE0MTc0MDAwIiwiZW1haWwiOiJ0ZWFjaGVyQGV4YW1wbGUuY29tIiwiZXhwIjoxNzA1ODM5NjAwfQ.abc123...",
        "user": {
            "id": "123e4567-e89b-12d3-a456-426614174000",
            "email": "teacher@example.com",
            "name": "John Teacher",
            "created_at": "2024-01-20T10:00:00Z",
            "updated_at": "2024-01-20T10:00:00Z",
            "roles": [
                {
                    "id": "550e8400-e29b-41d4-a716-446655440000",
                    "name": "teacher",
                    "description": "Teacher role with course management access",
                    "created_at": "2024-01-15T08:00:00Z"
                }
            ]
        }
    }
}
```

---

## Role Management Endpoints

### Create Role
**POST** `/api/roles/create`

Creates a new role. Requires admin privileges.

**Required Permission:** Admin role

**Request Headers:**
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "name": "teaching_assistant",
    "description": "Assistant role with limited course access"
}
```

**Response Codes:**
- `201 Created` - Role created successfully
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing/invalid token
- `403 Forbidden` - Insufficient permissions
- `409 Conflict` - Role name already exists

**Success Response (201):**
```json
{
    "success": true,
    "message": "Role created successfully",
    "data": {
        "id": "650e8400-e29b-41d4-a716-446655440001",
        "name": "teaching_assistant",
        "description": "Assistant role with limited course access",
        "created_at": "2024-01-20T11:00:00Z"
    }
}
```

---

### List All Roles
**GET** `/api/roles`

Retrieves all available roles in the system.

**Required Permission:** Authenticated user

**Response (200):**
```json
{
    "success": true,
    "message": "Roles retrieved successfully",
    "data": [
        {
            "id": "450e8400-e29b-41d4-a716-446655440001",
            "name": "student",
            "description": "Student role with basic access",
            "created_at": "2024-01-01T00:00:00Z"
        },
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "name": "teacher",
            "description": "Teacher role with course management access",
            "created_at": "2024-01-01T00:00:00Z"
        },
        {
            "id": "650e8400-e29b-41d4-a716-446655440002",
            "name": "admin",
            "description": "Administrator role with full access",
            "created_at": "2024-01-01T00:00:00Z"
        }
    ]
}
```

---

### Get User Roles
**GET** `/api/users/:userID/roles`

Retrieves all roles assigned to a specific user.

**Required Permission:** Authenticated user (can view own roles) or Admin

**URL Parameters:**
- `userID` - UUID of the user

**Response (200):**
```json
{
    "success": true,
    "message": "User roles retrieved successfully",
    "data": [
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "name": "teacher",
            "description": "Teacher role with course management access",
            "created_at": "2024-01-01T00:00:00Z"
        }
    ]
}
```

---

### Assign Role to User
**POST** `/api/users/assign-role`

Assigns a role to a user. Requires admin privileges.

**Required Permission:** Admin role

**Request Body:**
```json
{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "role_id": "650e8400-e29b-41d4-a716-446655440001"
}
```

**Response Codes:**
- `200 OK` - Role assigned successfully
- `400 Bad Request` - Invalid UUIDs
- `404 Not Found` - User or role not found
- `409 Conflict` - User already has this role

**Success Response (200):**
```json
{
    "success": true,
    "message": "Role assigned successfully",
    "data": null
}
```

---

### Remove Role from User
**DELETE** `/api/users/:userID/roles/:roleID`

Removes a role from a user. Requires admin privileges.

**Required Permission:** Admin role

**URL Parameters:**
- `userID` - UUID of the user
- `roleID` - UUID of the role to remove

**Response (200):**
```json
{
    "success": true,
    "message": "Role removed successfully",
    "data": null
}
```

---

## Permission Management Endpoints

### Create Permission
**POST** `/api/permissions/create`

Creates a new permission. Requires admin privileges.

**Required Permission:** Admin role

**Request Body:**
```json
{
    "name": "manage_assignments",
    "resource": "assignments",
    "action": "manage",
    "description": "Full control over assignments"
}
```

**Response (201):**
```json
{
    "success": true,
    "message": "Permission created successfully",
    "data": {
        "id": "750e8400-e29b-41d4-a716-446655440003",
        "name": "manage_assignments",
        "resource": "assignments",
        "action": "manage",
        "description": "Full control over assignments",
        "created_at": "2024-01-20T12:00:00Z"
    }
}
```

---

### List All Permissions
**GET** `/api/permissions`

Retrieves all available permissions in the system.

**Required Permission:** Authenticated user

**Response (200):**
```json
{
    "success": true,
    "message": "Permissions retrieved successfully",
    "data": [
        {
            "id": "850e8400-e29b-41d4-a716-446655440004",
            "name": "view_course",
            "resource": "course",
            "action": "read",
            "description": "View course details",
            "created_at": "2024-01-01T00:00:00Z"
        },
        {
            "id": "850e8400-e29b-41d4-a716-446655440005",
            "name": "create_course",
            "resource": "course",
            "action": "create",
            "description": "Create new courses",
            "created_at": "2024-01-01T00:00:00Z"
        }
    ]
}
```

---

### Get Role Permissions
**GET** `/api/roles/:roleID/permissions`

Retrieves all permissions assigned to a specific role.

**URL Parameters:**
- `roleID` - UUID of the role

**Response (200):**
```json
{
    "success": true,
    "message": "Role permissions retrieved successfully",
    "data": [
        {
            "id": "850e8400-e29b-41d4-a716-446655440004",
            "name": "view_course",
            "resource": "course",
            "action": "read",
            "description": "View course details",
            "created_at": "2024-01-01T00:00:00Z"
        },
        {
            "id": "850e8400-e29b-41d4-a716-446655440005",
            "name": "create_course",
            "resource": "course",
            "action": "create",
            "description": "Create new courses",
            "created_at": "2024-01-01T00:00:00Z"
        }
    ]
}
```

---

### Grant Permission to Role
**POST** `/api/permissions/grant`

Grants a permission to a role. Requires admin privileges.

**Required Permission:** Admin role

**Request Body:**
```json
{
    "role_id": "550e8400-e29b-41d4-a716-446655440000",
    "permission_id": "750e8400-e29b-41d4-a716-446655440003"
}
```

**Response (200):**
```json
{
    "success": true,
    "message": "Permission granted successfully",
    "data": null
}
```

---

### Revoke Permission from Role
**DELETE** `/api/roles/:roleID/permissions/:permissionID`

Revokes a permission from a role. Requires admin privileges.

**Required Permission:** Admin role

**URL Parameters:**
- `roleID` - UUID of the role
- `permissionID` - UUID of the permission to revoke

**Response (200):**
```json
{
    "success": true,
    "message": "Permission revoked successfully",
    "data": null
}
```

---

## Protected Resource Endpoints

### List Courses
**GET** `/api/courses`

Retrieves a list of courses. Requires course:read permission.

**Required Permission:** `course:read`

**Query Parameters:**
- `page` (optional) - Page number for pagination (default: 1)
- `limit` (optional) - Items per page (default: 20)
- `search` (optional) - Search term for course name

**Response (200):**
```json
{
    "success": true,
    "message": "Courses retrieved successfully",
    "data": [
        {
            "id": "950e8400-e29b-41d4-a716-446655440006",
            "name": "Introduction to Computer Science",
            "description": "Basic programming concepts and algorithms",
            "teacher_id": "123e4567-e89b-12d3-a456-426614174000",
            "created_at": "2024-01-20T13:00:00Z"
        }
    ]
}
```

---

### Create Course
**POST** `/api/courses`

Creates a new course. Requires course:create permission.

**Required Permission:** `course:create`

**Request Body:**
```json
{
    "name": "Advanced Mathematics",
    "description": "Advanced math course for senior students"
}
```

**Response (201):**
```json
{
    "success": true,
    "message": "Course created successfully",
    "data": {
        "id": "a50e8400-e29b-41d4-a716-446655440007",
        "name": "Advanced Mathematics",
        "description": "Advanced math course for senior students",
        "teacher_id": "123e4567-e89b-12d3-a456-426614174000",
        "created_at": "2024-01-20T14:00:00Z"
    }
}
```

---

### View Grades
**GET** `/api/grades`

Retrieves grades. Requires grades:read permission.

**Required Permission:** `grades:read`

**Response (200):**
```json
{
    "success": true,
    "message": "Grades retrieved successfully",
    "data": [
        {
            "student_id": "b50e8400-e29b-41d4-a716-446655440008",
            "course_id": "950e8400-e29b-41d4-a716-446655440006",
            "grade": "A",
            "points": 95.5,
            "submitted_at": "2024-01-15T16:00:00Z"
        }
    ]
}
```

---

## Error Response Format

All error responses follow this format:

```json
{
    "success": false,
    "message": "An error occurred",
    "error": "Detailed error message here"
}
```

### Common Error Codes

| Status Code | Description | Common Causes |
|-------------|-------------|---------------|
| 400 | Bad Request | Invalid input, missing required fields |
| 401 | Unauthorized | Missing or invalid authentication token |
| 403 | Forbidden | Insufficient permissions for the requested resource |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource already exists (e.g., duplicate email) |
| 422 | Unprocessable Entity | Valid request but unable to process |
| 500 | Internal Server Error | Server-side error |

---

## Rate Limiting

API endpoints have the following rate limits:
- Authentication endpoints: 5 requests per minute per IP
- Other endpoints: 100 requests per minute per user

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1705839600
```

---

## Pagination

List endpoints support pagination with these query parameters:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20, max: 100)

Paginated responses include metadata:
```json
{
    "success": true,
    "message": "Data retrieved successfully",
    "data": [...],
    "meta": {
        "page": 1,
        "limit": 20,
        "total": 100,
        "total_pages": 5
    }
}
```