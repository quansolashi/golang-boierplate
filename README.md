# golang-boierplate

A modern Golang backend boilerplate featuring GraphQL, REST, MySQL, Redis, RabbitMQ, and Docker support. This project is designed for rapid backend development with best practices, modular structure, and production-ready configurations.

## Features

- **GraphQL & REST API**: Flexible API design with gqlgen and Gin.
- **Database**: MySQL integration using GORM and Ent.
- **Caching**: Redis support for fast data access.
- **Messaging**: RabbitMQ for asynchronous processing.
- **Swagger Docs**: Auto-generated API documentation.
- **Dockerized**: Easy local development and deployment.
- **Hot Reload**: Live code reloading with [air](https://github.com/cosmtrek/air).
- **Linting & Testing**: Integrated with golangci-lint and Go test.

---

## Project Structure

```
.
├── backend/         # Main backend application
│   ├── cmd/         # Application entrypoint
│   ├── internal/    # Application logic
│   ├── pkg/         # Shared packages (config, db, etc.)
│   ├── ent/         # Ent ORM files
│   ├── graph/       # GraphQL schema and resolvers
│   ├── Makefile     # Development commands
│   ├── Dockerfile   # Production Dockerfile
│   └── ...
├── db/              # Database/Redis/RabbitMQ data
├── infra/           # Infrastructure configs (nginx, certs)
├── docker-compose.yml
└── README.md
```

---

## Prerequisites

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- (For local dev) [Go 1.23+](https://golang.org/dl/), [Air](https://github.com/cosmtrek/air), [Make](https://www.gnu.org/software/make/)

---

## Environment Variables

The backend uses environment variables for configuration. A template file is provided at `backend/.env.example`.

1. Copy `.env.example` to `.env` in the `backend/` directory:

   ```sh
   cp backend/.env.example backend/.env
   ```

2. Edit `backend/.env` and fill in the required values for your environment.

---

## Running the Backend

### 1. Using Docker Compose (Recommended)

This will start MySQL, Redis, RabbitMQ, Nginx, and the backend server.

```sh
docker-compose up --build
```

- The backend will be available at `http://localhost:3002`
- Nginx will proxy at `http://localhost` (port 80)
- MySQL: `localhost:3306`, Redis: `localhost:6379`, RabbitMQ: `localhost:5672`

### 2. Local Development (Hot Reload)

1. **Install dependencies:**

   ```sh
   cd backend
   make setup
   ```

2. **Start MySQL, Redis, RabbitMQ (via Docker Compose):**

   ```sh
   docker-compose up mysql redis rabbitmq
   ```

3. **Run the backend with hot reload:**

   ```sh
   make dev
   ```

   This uses [air](https://github.com/cosmtrek/air) for live reloading.

---

## Useful Makefile Commands

- `make build` – Build the backend binary
- `make dev` – Run with hot reload (requires Air)
- `make test` – Run tests
- `make lint` – Run linter
- `make gqlgen` – Regenerate GraphQL code
- `make entgen` – Regenerate Ent ORM code

---

## Database Migration

To run database migrations:

```sh
cd backend
go run hack/migration/main.go
```

Or use your preferred migration tool.

---

## API Documentation

- Swagger UI is auto-generated and available at `/docs/web` (if enabled in the server).

---

## License

MIT

---

Let me know if you want to add more details or usage examples!
