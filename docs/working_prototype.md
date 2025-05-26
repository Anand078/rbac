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
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    
    "github.com/yourusername/rbac-system/internal/database"
    "github.com/yourusername/rbac-system/internal/models"
    "github.com/yourusername/rbac-system/pkg/utils"
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
    rows, err := h.db.Query(query, search, limit