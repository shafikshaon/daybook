# Docker Deployment - Local Machine

Deploy Daybook with Docker while using your host machine's PostgreSQL and Redis.

## 📋 Prerequisites

1. **Docker Desktop** installed and running
2. **PostgreSQL** running on host machine (port 5432)
3. **Redis** running on host machine (port 6379) - optional

## 🚀 Quick Start

### 1. Start Host Services

Make sure PostgreSQL and Redis are running on your host:

```bash
# Check PostgreSQL
psql -U postgres -h localhost -p 5432

# Check Redis
redis-cli ping
```

### 2. Configure Database

Create the database:
```bash
psql -U postgres -h localhost
CREATE DATABASE daybook;
\q
```

### 3. Deploy with Docker

```bash
# Just run this!
./docker-deploy.sh
```

That's it! Your app is running.

## 🌐 Access Your Application

- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

## 📊 Manage Containers

```bash
# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# View all logs
docker-compose logs -f

# Stop containers
docker-compose down

# Restart containers
docker-compose restart

# Rebuild and restart
docker-compose up -d --build

# View running containers
docker ps
```

## ⚙️ Configuration

### Database Connection

Edit `docker-compose.yml`:
```yaml
environment:
  - DB_HOST=host.docker.internal  # Connects to host PostgreSQL
  - DB_PORT=5432
  - DB_NAME=daybook
  - DB_USER=postgres
  - DB_PASSWORD=123456  # Change this!
```

### Frontend API URL

Edit `frontend/.env.docker`:
```
VITE_API_URL=http://localhost:8080/api/v1
```

### Change Ports

Edit `docker-compose.yml`:
```yaml
ports:
  - "8080:8080"  # Backend: host:container
  - "3000:80"    # Frontend: host:container
```

## 🔧 How It Works

```
┌─────────────────────────────────────────────┐
│           Your Host Machine                 │
│                                             │
│  PostgreSQL (5432) ◄─────┐                │
│  Redis (6379)      ◄─────┤                │
│                           │                 │
│  ┌────────────────────────┴──────────────┐ │
│  │     Docker Containers                 │ │
│  │                                       │ │
│  │  ┌──────────────┐  ┌───────────────┐ │ │
│  │  │   Backend    │  │   Frontend    │ │ │
│  │  │   :8080      │  │   :3000       │ │ │
│  │  └──────────────┘  └───────────────┘ │ │
│  └───────────────────────────────────────┘ │
│                                             │
│  Browser: http://localhost:3000             │
└─────────────────────────────────────────────┘
```

**Key Points:**
- Backend and Frontend run in Docker containers
- PostgreSQL runs on your host machine
- Redis runs on your host machine
- Containers connect to host via `host.docker.internal`

## 🐛 Troubleshooting

### Backend Can't Connect to PostgreSQL

```bash
# Check if PostgreSQL is listening on all interfaces
psql -U postgres -h localhost

# If not, edit postgresql.conf:
# listen_addresses = '*'

# Edit pg_hba.conf to allow connections from Docker:
# host    all    all    172.17.0.0/16    md5

# Restart PostgreSQL
```

### Containers Not Starting

```bash
# Check logs
docker-compose logs

# Rebuild images
docker-compose build --no-cache

# Remove old containers
docker-compose down -v
docker-compose up -d
```

### Port Already in Use

```bash
# Find what's using the port
lsof -i :8080
lsof -i :3000

# Kill the process or change port in docker-compose.yml
```

### Frontend Can't Reach Backend

Make sure backend is accessible:
```bash
# Test from host
curl http://localhost:8080/health

# If that works but frontend doesn't, rebuild frontend:
docker-compose up -d --build frontend
```

## 📦 Manual Commands

If you prefer manual control:

```bash
# Build images
docker-compose build

# Start containers
docker-compose up -d

# Stop containers
docker-compose down

# View logs
docker-compose logs -f

# Restart a service
docker-compose restart backend

# Rebuild a service
docker-compose up -d --build backend
```

## 🔄 Update Application

When you change code:

```bash
# Backend changes
docker-compose up -d --build backend

# Frontend changes
docker-compose up -d --build frontend

# Both
docker-compose up -d --build
```

## 🗑️ Clean Up

```bash
# Stop and remove containers
docker-compose down

# Remove containers and volumes
docker-compose down -v

# Remove images
docker rmi daybook-backend daybook-frontend
```

## 📝 Environment Variables

### Backend (.env in container)
- `SERVER_HOST`: 0.0.0.0
- `SERVER_PORT`: 8080
- `DB_HOST`: host.docker.internal
- `DB_PASSWORD`: Change in docker-compose.yml
- `JWT_SECRET`: Change in docker-compose.yml

### Frontend (.env.docker)
- `VITE_API_URL`: http://localhost:8080/api/v1

## 🎯 Production Tips

For production deployment:

1. **Use Docker Secrets** for passwords
2. **Use environment variables** from a `.env` file
3. **Set up HTTPS** with a reverse proxy
4. **Use volumes** for persistent data
5. **Set resource limits** in docker-compose.yml
6. **Enable health checks**
7. **Set up log rotation**

## 💡 Why This Setup?

**Pros:**
- ✅ Keep your existing PostgreSQL and Redis
- ✅ Easy to develop and test
- ✅ Containers are isolated
- ✅ No data loss if containers restart
- ✅ Fast rebuilds

**Cons:**
- ⚠️ Requires host services to be running
- ⚠️ Not fully portable (needs host setup)

For a fully containerized setup, see `KUBERNETES.md` or use Docker Compose with PostgreSQL/Redis containers.
