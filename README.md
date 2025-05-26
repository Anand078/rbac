# rbac

A Go-based application implementing Role-Based Access Control (RBAC) with authentication.

## Features

- User Registration and Login
- JWT-based Authentication
- Role Management (Create, List, Assign/Remove to Users)
- Permission Management (Create, List, Grant/Revoke to Roles)
- Endpoint Authorization based on roles and permissions
- Integration with Supabase

## Technologies Used

- Go (version 1.24.3)
- Gin (Web Framework)
- JWT (for authentication)
- Supabase (Database and Auth integration)
- PostgreSQL (Database)

## Prerequisites

- Go 1.24.3 or later
- Docker (for containerized deployment)
- A PostgreSQL database
- Supabase project with database and authentication enabled

## Setup

1. Clone the repository:

```bash
git clone https://github.com/Anand078/rbac.git
cd rbac
```

2. Create a `.env` file in the root directory with the following environment variables:

```
DATABASE_URL=your_database_connection_string
SUPABASE_URL=your_supabase_project_url
SUPABASE_SERVICE_KEY=your_supabase_service_key
JWT_SECRET=your_jwt_secret_key
PORT=8080
```

Replace the placeholder values with your actual database and Supabase credentials.

3. Ensure your PostgreSQL database schema is set up for RBAC. (Details on specific tables might be needed, refer to internal/database for schema details if necessary).

## Running the Application

### Locally

Make sure you have Go installed and the environment variables set in your `.env` file.

```bash
go run cmd/main.go
```

The application should start on the port specified in the `.env` file (default is 8080).

### With Docker

Make sure you have Docker installed and the `.env` file created.

1. Build the Docker image:

```bash
docker build -t rbac-system .
```

2. Run the Docker container:

```bash
docker run -p 8080:8080 --env-file .env rbac-system
```

The application will be accessible at `http://localhost:8080`.

## API Endpoints

Here are some of the main API endpoints:

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login a user and get a JWT token
- `POST /api/roles/create` - Create a new role (Admin only)
- `GET /api/roles` - Get all roles (Authenticated users)
- `GET /api/users/:userID/roles` - Get roles for a specific user (Authenticated users)
- `POST /api/users/assign-role` - Assign a role to a user (Admin only)
- `DELETE /api/users/:userID/roles/:roleID` - Remove a role from a user (Admin only)
- `POST /api/permissions/create` - Create a new permission (Admin only)
- `GET /api/permissions` - Get all permissions (Authenticated users)
- `GET /api/roles/:roleID/permissions` - Get permissions for a specific role (Authenticated users)
- `POST /api/permissions/grant` - Grant a permission to a role (Admin only)
- `DELETE /api/roles/:roleID/permissions/:permissionID` - Revoke a permission from a role (Admin only)
- `GET /health` - Health check endpoint

(Note: Specific request/response bodies and detailed authorization rules for each endpoint would require deeper code inspection or documentation. This list is based on the routes defined in `cmd/main.go`.)