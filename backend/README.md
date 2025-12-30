# Daybook Backend API

A comprehensive personal finance management backend built with Go, Gin, GORM, PostgreSQL, and Redis.

## Features

- **User Authentication** - JWT-based authentication with signup, login, and profile management
- **Accounts** - Manage multiple accounts (cash, checking, savings, credit cards, brokerage)
- **Transactions** - Track income, expenses, and transfers with categories and tags
- **Credit Cards** - Manage credit cards, statements, payments, and rewards
- **Investments** - Portfolio management with stocks, bonds, ETFs, crypto, and dividend tracking
- **Bills** - Recurring bill tracking with payment reminders
- **Budgets** - Category-based budgets with progress tracking and alerts
- **Savings Goals** - Goal setting with contribution tracking and automated rules
- **Fixed Deposits** - FD management with interest calculations
- **Settings** - User preferences (currency, theme, notifications)
- **Reports & Analytics** - Transaction statistics and financial insights

## Tech Stack

- **Go 1.21+** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **PostgreSQL** - Primary database
- **Redis** - Caching layer
- **JWT** - Authentication tokens
- **Docker** - Containerization
- **Viper** - Configuration management
- **Datadog APM** - Application performance monitoring (optional)

## Project Structure

```
backend/
├── cmd/                    # Command-line tools
├── config/                 # Configuration management
├── database/               # Database initialization and connection
├── handlers/               # HTTP request handlers
├── logger/                 # Logging utilities
├── middleware/             # HTTP middleware (auth, CORS, etc.)
├── models/                 # Data models
├── repository/             # Data access layer
├── routes/                 # Route definitions
├── services/               # Business logic
├── storage/                # File storage
├── utilities/              # Helper functions
├── main.go                 # Application entry point
├── config.yaml             # Configuration file
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose configuration
├── Makefile                # Build automation
└── README.md               # This file
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Redis 7+ (optional, for caching)
- Docker & Docker Compose (for containerized deployment)

## Getting Started

### 1. Clone the Repository

```bash
cd backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `config.yaml` or set environment variables to configure:
- Database connection
- Redis connection
- JWT secret
- CORS settings
- Server port

### 4. Setup Database

Make sure PostgreSQL is running, then the application will automatically create tables on startup.

### 5. Run the Application

#### Option A: Local Development

```bash
# Run directly
go run main.go

# Or use make
make run
```

#### Option B: Docker Compose

```bash
# Start all services (PostgreSQL, Redis, Backend)
make docker-up

# View logs
make docker-logs

# Stop all services
make docker-down
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - Register new user
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/me` - Get current user profile
- `PUT /api/v1/auth/profile` - Update user profile
- `PUT /api/v1/auth/change-password` - Change password

### Accounts
- `GET /api/v1/accounts` - List all accounts
- `POST /api/v1/accounts` - Create account
- `GET /api/v1/accounts/:id` - Get account
- `PUT /api/v1/accounts/:id` - Update account
- `DELETE /api/v1/accounts/:id` - Delete account
- `PATCH /api/v1/accounts/:id/balance` - Update balance

### Transactions
- `GET /api/v1/transactions` - List transactions (with filters)
- `POST /api/v1/transactions` - Create transaction
- `POST /api/v1/transactions/bulk` - Bulk import
- `GET /api/v1/transactions/stats` - Get statistics
- `GET /api/v1/transactions/:id` - Get transaction
- `PUT /api/v1/transactions/:id` - Update transaction
- `DELETE /api/v1/transactions/:id` - Delete transaction

### Credit Cards
- `GET /api/v1/credit-cards` - List credit cards
- `POST /api/v1/credit-cards` - Create credit card
- `POST /api/v1/credit-cards/:id/payment` - Record payment
- `GET /api/v1/credit-cards/:id/statements` - Get statements
- `GET /api/v1/rewards` - List rewards
- `POST /api/v1/rewards` - Record reward

### Investments
- `GET /api/v1/investments` - List investments
- `POST /api/v1/investments` - Create investment
- `POST /api/v1/investments/:id/buy` - Buy shares
- `POST /api/v1/investments/:id/sell` - Sell shares
- `GET /api/v1/portfolios` - List portfolios
- `POST /api/v1/portfolios` - Create portfolio
- `GET /api/v1/dividends` - List dividends
- `POST /api/v1/dividends` - Record dividend

### Bills
- `GET /api/v1/bills` - List bills
- `POST /api/v1/bills` - Create bill
- `POST /api/v1/bills/:id/pay` - Mark as paid
- `GET /api/v1/bill-payments` - Payment history

### Budgets
- `GET /api/v1/budgets` - List budgets
- `POST /api/v1/budgets` - Create budget
- `GET /api/v1/budgets/:id/progress` - Get progress

### Savings Goals
- `GET /api/v1/savings-goals` - List goals
- `POST /api/v1/savings-goals` - Create goal
- `POST /api/v1/savings-goals/:id/contribute` - Add contribution
- `POST /api/v1/savings-goals/:id/withdraw` - Withdraw
- `GET /api/v1/automated-rules` - List rules
- `POST /api/v1/automated-rules` - Create rule

### Fixed Deposits
- `GET /api/v1/fixed-deposits` - List FDs
- `POST /api/v1/fixed-deposits` - Create FD
- `POST /api/v1/fixed-deposits/:id/withdraw` - Withdraw FD

### Settings
- `GET /api/v1/settings` - Get settings
- `PUT /api/v1/settings` - Update settings

### Health Check
- `GET /health` - API health status

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

Get a token by calling the `/api/v1/auth/login` endpoint.

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Vet Code

```bash
make vet
```

## Docker Deployment

### Build Docker Image

```bash
make docker-build
```

### Start Services

```bash
make docker-up
```

This will start:
- PostgreSQL on port 5432
- Redis on port 6379
- Backend API on port 8080

## Environment Variables

Override configuration with environment variables:

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=daybook

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT Configuration
JWT_SECRET=your-secret-key

# Datadog APM (Optional)
DD_ENABLED=false
DD_SERVICE=daybook-backend
DD_ENV=development
DD_AGENT_HOST=localhost
DD_TRACE_AGENT_PORT=8126
```

## Database Migrations

Database schema is automatically created/updated on application startup using GORM Auto Migration.

## Security

- Passwords are hashed using bcrypt
- JWT tokens for stateless authentication
- CORS configured for frontend integration
- User-specific data isolation
- SQL injection protection via GORM

## Performance

- Redis caching for frequently accessed data
- Database indexes on foreign keys and frequently queried fields
- Connection pooling for database connections
- Soft deletes for data recovery

## Monitoring with Datadog APM

The application includes **comprehensive Datadog APM integration** for tracking everything - HTTP requests, database queries, Redis operations, and custom business metrics.

### Quick Start

For detailed setup instructions, see **[DATADOG_SETUP.md](DATADOG_SETUP.md)** (Linux installation guide)
For monitoring details, see **[MONITORING_GUIDE.md](MONITORING_GUIDE.md)** (what gets tracked)

### Prerequisites

- Datadog account (sign up at https://www.datadoghq.com/)
- Datadog Agent installed and running on your system or infrastructure

### What Gets Tracked

When enabled, Datadog tracks:
- ✅ All HTTP requests (endpoints, latency, status codes, errors)
- ✅ All database queries (SQL statements, execution time, rows affected)
- ✅ All Redis operations (commands, keys, cache hits/misses)
- ✅ Custom business metrics (user registrations, transactions, payments, etc.)
- ✅ Errors and exceptions with full stack traces
- ✅ Service dependencies and distributed traces

### Configuration

Enable Datadog monitoring by setting the following environment variables in your `.env` file:

```bash
# Enable Datadog APM
DD_ENABLED=true

# Service name (appears in Datadog UI)
DD_SERVICE=daybook-backend

# Environment name (e.g., development, staging, production)
DD_ENV=production

# Datadog Agent host
DD_AGENT_HOST=localhost

# Datadog Agent trace port
DD_TRACE_AGENT_PORT=8126
```

### What Gets Monitored

When Datadog is enabled, the following metrics are automatically collected:

- **HTTP Requests**: All API endpoints with request/response times, status codes, and errors
- **Database Queries**: GORM database operations with query performance
- **Service Performance**: Request latency, throughput, and error rates
- **Distributed Traces**: End-to-end request flow across services
- **Custom Tags**: User ID, request ID, and other contextual information

### Viewing Traces

1. Start the Datadog Agent on your system
2. Enable Datadog in your `.env` file (set `DD_ENABLED=true`)
3. Start the application
4. Make some API requests
5. View traces in your Datadog dashboard at https://app.datadoghq.com/apm/traces

### Production Deployment

For production environments:

```bash
DD_ENABLED=true
DD_SERVICE=daybook-backend
DD_ENV=production
DD_AGENT_HOST=datadog-agent-service  # Your Datadog Agent hostname
DD_TRACE_AGENT_PORT=8126
```

### Disabling Datadog

To disable Datadog monitoring, set `DD_ENABLED=false` in your `.env` file or remove the environment variable. The application will run normally without any tracing overhead.

## Error Handling

All API responses follow a standard format:

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please create an issue in the repository.
