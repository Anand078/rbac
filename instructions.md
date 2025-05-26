# RBAC or ACL Module Implementation Guide

## Overview

This document outlines the design and implementation of a Role-Based Access Control (RBAC) or Access Control List (ACL) module for an edtech platform using Golang and Supabase.

### Scope

The RBAC or ACL system will be used to manage access to various parts of the platform and interactions between microservices. This includes:

- Users (students, teachers)
- Microservice-to-microservice permissions (e.g., Service A can only perform operations X, Y, Z on Service B)

### Assumptions

- Users are defined as either students or teachers.
- The system supports role-based access where each user is assigned roles that determine their allowed actions.
- Microservices will interact with each other using the RBAC or ACL system for enforcing permissions.

## Technical Design

### Database Schema Design

#### Entity-Relationship Diagram (ERD)

The following tables and relationships are defined for the RBAC system:

1.  **Users Table**:
    - `id` (Primary Key, Integer)
    - `name` (String)
    - `email` (String, Unique)
    - `role_id` (Foreign Key referencing Roles Table)

2.  **Roles Table**:
    - `id` (Primary Key, Integer)
    - `name` (String, e.g., `student`, `teacher`, `admin`)

3.  **Permissions Table**:
    - `id` (Primary Key, Integer)
    - `name` (String, e.g., `create_course`, `view_grades`)

4.  **Role_Permissions Table** (Many-to-many relationship between Roles and Permissions):
    - `role_id` (Foreign Key referencing Roles Table)
    - `permission_id` (Foreign Key referencing Permissions Table)

5.  **User_Roles Table** (Many-to-many relationship between Users and Roles):
    - `user_id` (Foreign Key referencing Users Table)
    - `role_id` (Foreign Key referencing Roles Table)

#### Example ERD

```plaintext
Users --< User_Roles >-- Roles --< Role_Permissions >-- Permissions
```

### API Design

The following APIs are designed to manage roles, permissions, and user access:

1.  **POST /api/roles/create**
    - Description: Create a new role.
    - Request Body:

    ```json
    {
      "name": "teacher"
    }
    ```

    - Response:

    ```json
    {
      "id": 1,
      "name": "teacher"
    }
    ```

2.  **GET /api/users/{userID}/roles**
    - Description: Retrieve roles assigned to a user.
    - Response:

    ```json
    [
      {
        "id": 1,
        "name": "teacher"
      }
    ]
    ```

3.  **POST /api/permissions/grant**
    - Description: Grant permissions to a role.
    - Request Body:

    ```json
    {
      "role_id": 1,
      "permission_id": 2
    }
    ```

    - Response:

    ```json
    {
      "status": "success",
      "message": "Permission granted"
    }
    ```

4.  **DELETE /api/permissions/revoke**
    - Description: Revoke permissions from a role.
    - Request Body:

    ```json
    {
      "role_id": 1,
      "permission_id": 2
    }
    ```

    - Response:

    ```json
    {
      "status": "success",
      "message": "Permission revoked"
    }
    ```

### Sequence Diagrams

The following sequence diagrams illustrate the flow of authorization checks and role/permission management.

#### Authorization Flow

1.  **User Request for Resource**:
    - User makes a request to access a resource.
    - The backend system checks the user's roles and associated permissions to determine if access is allowed.
    - If access is granted, the requested resource is provided. If denied, an error message is returned.

```plaintext
+--------+     +--------------+     +---------------+     +------------------+
|  User  | --> | Backend API  | --> | Authorization  | --> | Resource Service |
+--------+     +--------------+     +---------------+     +------------------+
```

## Implementation

### Working Prototype

The following Go code snippets demonstrate how the RBAC or ACL system can be implemented.

#### Example Code: Assigning a Role to a User

```go
package main

import (
    "fmt"
    "github.com/supabase/supabase-go"
)

func AssignRole(userID int, roleID int) error {
    // Connect to Supabase DB
    db := supabase.NewClient("https://your-project-url", "your-api-key")

    // Update user's role in the database
    _, err := db.Table("user_roles").Insert(map[string]interface{}{
        "user_id": userID,
        "role_id": roleID,
    })

    if err != nil {
        return fmt.Errorf("failed to assign role: %v", err)
    }
    return nil
}
```

#### API Integration Example

Here is an example of using the role assignment API to assign roles to a user:

```bash
curl -X POST https://your-api-url/api/roles/create \
  -H "Content-Type: application/json" \
  -d '{"name": "teacher"}'
```

## Documentation

The following documentation provides guidance on how to use the RBAC or ACL system.

### API Usage Guide

1.  **Creating Roles**:

    Use the POST /api/roles/create endpoint to create new roles like `student`, `teacher`, or `admin`.

2.  **Assigning Roles to Users**:

    The system allows for users to be assigned specific roles using the POST /api/users/{userID}/roles endpoint.

3.  **Managing Permissions**:

    Permissions can be granted or revoked from roles using the POST /api/permissions/grant and DELETE /api/permissions/revoke endpoints.

### Developer Integration

Microservices can check roles and permissions using the authorization middleware.

Ensure that each service checks the permission before performing any action based on the role of the requesting user.

### Deployment Instructions

1.  **Set up Supabase**: Create a new project on Supabase and set up your tables as outlined in the database schema.
2.  **Deploy Backend Service**: Deploy the Go backend service to your preferred cloud provider (e.g., AWS, Google Cloud).
3.  **Testing**: Use Postman or cURL to test the API endpoints. Ensure that roles and permissions are functioning as expected.

