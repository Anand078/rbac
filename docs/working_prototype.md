# Working Prototype

## Overview

This section demonstrates a complete working example of the RBAC system integrated with a Course Management feature. The prototype includes:

1. Database setup for courses
2. Complete handler implementation
3. Testing scenarios with cURL commands
4. Expected responses

## 1. Database Setup

First, create the courses table in your Supabase database:

```sql
-- Create courses table
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    teacher_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for teacher lookup
CREATE INDEX idx_courses_teacher_id ON courses(teacher_id);

-- Create a sample course for testing
INSERT INTO courses (name, description, teacher_id) 
VALUES (
    'Introduction to Programming',
    'Learn the basics of programming with Python',
    (SELECT id FROM users WHERE email = 'teacher@example.com' LIMIT 1)
);
```

## 2. Course Handler Implementation

### internal/models/course.go

```go
package models

import (
    "time"
    "github.com/google/uuid"
)

type Course struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    TeacherID   uuid.UUID `json:"teacher_id" db:"teacher_id"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCourseRequest struct {
    Name        string `json:"name" binding:"required,min=3,max=255"`
    Description string `json:"description" binding:"max=1000"`
}

type UpdateCourseRequest struct {
    Name        string `json:"name" binding:"omitempty,min=3,max=255"`
    Description string `json:"description" binding:"omitempty,max=1000"`
}

type CourseListResponse struct {
    Courses []Course `json:"courses"`
    Total   int      `json:"total"`
    Page    int      `json:"page"`
    Limit   int      `json:"limit"`
}
```

### internal/handlers/courses.go

```go
package handlers

import (
    "database/sql"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "github.com/Anand078/rbac/internal/database"
    "github.com/Anand078/rbac/internal/models"
    "github.com/Anand078/rbac/pkg/utils"
)

type CourseHandler struct {
    db *database.DB
}

func NewCourseHandler(db *database.DB) *CourseHandler {
    return &CourseHandler{db: db}
}

// ListCourses - GET /api/courses
// Requires: course:read permission
func (h *CourseHandler) ListCourses(c *gin.Context) {
    // Get pagination parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    search := c.Query("search")

    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 20
    }

    offset := (page - 1) * limit

    // Build query
    query := `
        SELECT id, name, description, teacher_id, created_at, updated_at
        FROM courses
        WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

    // Get total count
    var total int
    countQuery := `
        SELECT COUNT(*) FROM courses
        WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
    `
    err := h.db.QueryRow(countQuery, search).Scan(&total)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to count courses")
        return
    }

    // Get courses
    rows, err := h.db.Query(query, search, limit, offset)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch courses")
        return
    }
    defer rows.Close()

    courses := []models.Course{}
    for rows.Next() {
        var course models.Course
        var teacherID sql.NullString // Use sql.NullString to handle potential NULL teacher_id
        err := rows.Scan(&course.ID, &course.Name, &course.Description, &teacherID, &course.CreatedAt, &course.UpdatedAt)
        if err != nil {
            log.Printf("Error scanning course row: %v", err) // Log the error but continue
            continue
        }
        if teacherID.Valid {
            course.TeacherID, _ = uuid.Parse(teacherID.String) // Parse only if not NULL
        } else {
            course.TeacherID = uuid.Nil // Set to Nil UUID if NULL
        }
        courses = append(courses, course)
    }

    if err = rows.Err(); err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Error iterating over courses")
        return
    }

    c.JSON(http.StatusOK, models.CourseListResponse{
        Courses: courses,
        Total:   total,
        Page:    page,
        Limit:   limit,
    })
}

// CreateCourse - POST /api/courses
// Requires: course:create permission
func (h *CourseHandler) CreateCourse(c *gin.Context) {
    var req models.CreateCourseRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }

    // Get the authenticated user's ID from the context (set by auth middleware)
    userID, exists := c.Get("user_id")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    // Ensure the user ID is a valid UUID
    teacherID, ok := userID.(string)
    if !ok {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format in context")
        return
    }
    teacherUUID, err := uuid.Parse(teacherID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format")
        return
    }

    course := models.Course{
        ID:          uuid.New(),
        Name:        req.Name,
        Description: req.Description,
        TeacherID:   teacherUUID,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    query := `
        INSERT INTO courses (id, name, description, teacher_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    err = h.db.QueryRow(query,
        course.ID,
        course.Name,
        course.Description,
        course.TeacherID,
        course.CreatedAt,
        course.UpdatedAt,
    ).Scan(&course.ID)

    if err != nil {
        log.Printf("Error creating course: %v", err)
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create course")
        return
    }

    c.JSON(http.StatusCreated, course)
}
```

## 3. Testing Scenarios

### List Courses

```bash
curl -X GET "http://localhost:8080/api/courses?page=1&limit=10" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### Create Course

```bash
curl -X POST "http://localhost:8080/api/courses" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Advanced Algorithms", "description": "Deep dive into algorithms"}'
```

## 4. Expected Responses

### List Courses (200 OK)

```json
{
  "courses": [
    {
      "id": "c1b2e3f4-5678-1234-9abc-def012345678",
      "name": "Introduction to Programming",
      "description": "Learn the basics of programming with Python",
      "teacher_id": "a1b2c3d4-5678-1234-9abc-def012345678",
      "created_at": "2024-05-26T10:00:00Z",
      "updated_at": "2024-05-26T10:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

### Create Course (201 Created)

```json
{
  "id": "d2e3f4a5-6789-1234-9abc-def012345678",
  "name": "Advanced Algorithms",
  "description": "Deep dive into algorithms",
  "teacher_id": "a1b2c3d4-5678-1234-9abc-def012345678",
  "created_at": "2024-05-26T10:05:00Z",
  "updated_at": "2024-05-26T10:05:00Z"
}
```

### Error Example (400 Bad Request)

```json
{
  "error": "Invalid request body"
}
```

---

This prototype demonstrates how to integrate RBAC with a course management feature, including database schema, handler logic, and example API usage.