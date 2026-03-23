# ppx

A CLI tool for generating Go web project scaffolding with modular architecture.

## Installation

```bash
go install github.com/lpphub/ppx@latest
```

Or build from source:

```bash
git clone https://github.com/lpphub/ppx.git
cd ppx
go build -o ppx .
```

## Usage

### Create a new project

```bash
ppx new myproject
```

Or with a custom module name:

```bash
ppx new myproject --module github.com/user/myproject
```

### Create a new module

```bash
cd myproject
ppx module product
```

## Generated Project Structure

```
myproject/
├── config/
│   └── config.yml            # Configuration file
├── modules/
│   ├── core/                 # Module interface
│   │   └── module.go
│   ├── auth/                 # Authentication module
│   │   ├── module.go
│   │   ├── dto.go
│   │   ├── handler.go
│   │   └── service.go
│   ├── user/                 # User module
│   │   ├── module.go
│   │   ├── model.go
│   │   ├── dto.go
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repo.go
│   └── post/                 # Demo CRUD module
│       ├── module.go
│       ├── model.go
│       ├── dto.go
│       ├── handler.go
│       ├── service.go
│       └── repo.go
├── infra/
│   ├── init.go               # Infrastructure init
│   ├── config.go             # Config loader
│   ├── dbs.go                # Database connection
│   └── jwt/
│       └── jwt.go            # JWT utilities
├── server/
│   ├── app.go                # HTTP server
│   ├── helper/
│   │   └── response.go       # Response helpers
│   └── middleware/
│       ├── auth.go           # JWT authentication
│       └── cors.go           # CORS middleware
├── shared/
│   ├── consts/               # Constants
│   ├── contracts/            # Module contracts
│   ├── errs/                 # Error definitions
│   ├── pagination/           # Pagination utilities
│   └── strutils/             # String utilities
├── main.go                   # Application entry point
├── go.mod
├── Makefile
└── Dockerfile
```

## Features

- **Modular Architecture**: Clean separation of concerns with modules
- **Built-in Auth**: JWT-based authentication ready to use
- **Demo CRUD**: Post module as a complete CRUD example
- **Clean Code**: Repository, Service, Handler layers
- **Configuration**: YAML-based configuration
- **Middleware**: CORS, JWT authentication, logging
- **Docker Ready**: Includes Dockerfile for deployment

## API Endpoints

### Authentication

| Method | Endpoint          | Description      | Auth |
|--------|-------------------|------------------|------|
| POST   | /auth/register    | Register user    | No   |
| POST   | /auth/login       | Login user       | No   |
| POST   | /auth/refresh     | Refresh token    | No   |

### Posts

| Method | Endpoint     | Description    | Auth |
|--------|--------------|----------------|------|
| GET    | /posts       | List posts     | No   |
| GET    | /posts/:id   | Get post       | No   |
| POST   | /posts       | Create post    | Yes  |
| PUT    | /posts/:id   | Update post    | Yes  |
| DELETE | /posts/:id   | Delete post    | Yes  |

## Quick Start

```bash
# Create a new project
ppx new myproject

# Navigate to project
cd myproject

# Update configuration
# Edit config/config.yml with your database credentials

# Install dependencies
go mod tidy

# Run the server
go run .
```

## License

MIT
