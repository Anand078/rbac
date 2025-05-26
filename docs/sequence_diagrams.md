# Sequence Diagrams

## 1. User Registration Flow

This diagram shows the complete user registration process including password hashing and role assignment.

```mermaid
sequenceDiagram
    participant User
    participant API as API Gateway
    participant AuthHandler as Auth Handler
    participant AuthService as Auth Service
    participant Bcrypt
    participant Database
    participant Response as Response Handler

    User->>API: POST /api/auth/register
    Note right of User: {email, name, password, role_id}
    
    API->>AuthHandler: Register(request)
    AuthHandler->>AuthHandler: Validate request body
    
    alt Invalid Request
        AuthHandler-->>User: 400 Bad Request
    end
    
    AuthHandler->>AuthService: Register(userData)
    AuthService->>Bcrypt: GenerateFromPassword(password)
    Bcrypt-->>AuthService: hashedPassword
    
    AuthService->>AuthService: Generate UUID for user
    
    AuthService->>Database: BEGIN TRANSACTION
    AuthService->>Database: INSERT INTO users (id, email, name, password_hash)
    
    alt Email Already Exists
        Database-->>AuthService: Constraint violation
        AuthService->>Database: ROLLBACK
        AuthService-->>AuthHandler: Error: Email exists
        AuthHandler-->>User: 409 Conflict
    end
    
    Database-->>AuthService: User created with timestamps
    
    opt Role ID Provided
        AuthService->>AuthService: Parse role_id to UUID
        AuthService->>Database: INSERT INTO user_roles (user_id, role_id)
        Database-->>AuthService: Role assigned
    end
    
    AuthService->>Database: COMMIT TRANSACTION
    Database-->>AuthService: Transaction complete
    
    AuthService-->>AuthHandler: User object (without password)
    AuthHandler->>Response: SuccessResponse(201, user)
    Response-->>User: 201 Created + User data
```

---

## 2. User Login Flow

This diagram illustrates the authentication process and JWT token generation.

```mermaid
sequenceDiagram
    participant User
    participant API as API Gateway
    participant AuthHandler as Auth Handler
    participant AuthService as Auth Service
    participant Database
    participant Bcrypt
    participant JWT as JWT Service
    participant Response as Response Handler

    User->>API: POST /api/auth/login
    Note right of User: {email, password}
    
    API->>AuthHandler: Login(request)
    AuthHandler->>AuthHandler: Validate request body
    
    alt Invalid Request
        AuthHandler-->>User: 400 Bad Request
    end
    
    AuthHandler->>AuthService: Login(credentials)
    AuthService->>Database: SELECT user WHERE email = ?
    
    alt User Not Found
        Database-->>AuthService: No rows
        AuthService-->>AuthHandler: Invalid credentials
        AuthHandler-->>User: 401 Unauthorized
    end
    
    Database-->>AuthService: User data with password_hash
    
    AuthService->>Bcrypt: CompareHashAndPassword(hash, password)
    
    alt Password Mismatch
        Bcrypt-->>AuthService: Error
        AuthService-->>AuthHandler: Invalid credentials
        AuthHandler-->>User: 401 Unauthorized
    end
    
    Bcrypt-->>AuthService: Password valid
    
    AuthService->>Database: SELECT roles for user_id
    Database-->>AuthService: User roles
    
    AuthService->>AuthService: Attach roles to user object
    
    AuthService->>JWT: GenerateToken(user_id, email, exp)
    Note right of JWT: Token expires in 24 hours
    JWT-->>AuthService: Signed JWT token
    
    AuthService-->>AuthHandler: LoginResponse{token, user}
    AuthHandler->>Response: SuccessResponse(200, data)
    Response-->>User: 200 OK + Token + User with roles
```

---

## 3. Authorization Check Flow

This diagram shows how the system validates permissions for protected endpoints.

```mermaid
sequenceDiagram
    participant User
    participant API as API Gateway
    participant AuthMiddleware as Auth Middleware
    participant JWT as JWT Parser
    participant Context as Request Context
    participant RBACService as RBAC Service
    participant Database
    participant Handler as Course Handler
    participant Response as Response Handler

    User->>API: GET /api/courses
    Note right of User: Authorization: Bearer <token>
    
    API->>AuthMiddleware: Authenticate()
    
    AuthMiddleware->>AuthMiddleware: Extract token from header
    
    alt No Authorization Header
        AuthMiddleware-->>User: 401 Unauthorized
    end
    
    AuthMiddleware->>JWT: Parse(token, secret)
    
    alt Invalid Token
        JWT-->>AuthMiddleware: Parse error
        AuthMiddleware-->>User: 401 Unauthorized
    end
    
    JWT-->>AuthMiddleware: Claims{user_id, email, exp}
    
    AuthMiddleware->>AuthMiddleware: Check expiration
    
    alt Token Expired
        AuthMiddleware-->>User: 401 Unauthorized
    end
    
    AuthMiddleware->>Context: Set("user_id", user_id)
    AuthMiddleware->>Context: Set("email", email)
    
    AuthMiddleware->>AuthMiddleware: Authorize("course", "read")
    AuthMiddleware->>RBACService: HasPermission(user_id, "course", "read")
    
    RBACService->>Database: Query user permissions
    Note right of Database: JOIN users, user_roles,<br/>role_permissions, permissions
    
    Database-->>RBACService: Permission check result
    
    alt No Permission
        RBACService-->>AuthMiddleware: false
        AuthMiddleware-->>User: 403 Forbidden
    end
    
    RBACService-->>AuthMiddleware: true
    AuthMiddleware->>API: Next() - Continue to handler
    
    API->>Handler: ListCourses(context)
    Handler->>Context: Get("user_id")
    Context-->>Handler: user_id
    
    Handler->>Database: SELECT courses
    Database-->>Handler: Course list
    
    Handler->>Response: SuccessResponse(200, courses)
    Response-->>User: 200 OK + Course data