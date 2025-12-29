# Datadog APM Setup Guide for Linux

This guide will help you set up comprehensive Datadog APM monitoring for the Daybook backend on Linux.

## Prerequisites

- Linux-based machine (Ubuntu, Debian, CentOS, RHEL, etc.)
- Sudo access
- Datadog account (sign up at https://www.datadoghq.com/)

## Step 1: Get Your Datadog API Key

1. Sign up for Datadog at https://www.datadoghq.com/
2. Log in to your Datadog account
3. Navigate to: **Organization Settings → API Keys**
   - URL: https://app.datadoghq.com/organization-settings/api-keys
4. Copy your API key (you'll need it for the next step)

## Step 2: Install Datadog Agent on Linux

### Option A: One-Line Installation (Recommended)

Replace `YOUR_API_KEY` with your actual Datadog API key:

```bash
DD_API_KEY=YOUR_API_KEY \
DD_SITE="datadoghq.com" \
DD_APM_ENABLED=true \
bash -c "$(curl -L https://s3.amazonaws.com/dd-agent/scripts/install_script.sh)"
```

### Option B: Manual Installation

#### Ubuntu/Debian:

```bash
# Add Datadog repository
sudo sh -c "echo 'deb [signed-by=/usr/share/keyrings/datadog-archive-keyring.gpg] https://apt.datadoghq.com/ stable 7' > /etc/apt/sources.list.d/datadog.list"

# Import Datadog GPG key
sudo touch /usr/share/keyrings/datadog-archive-keyring.gpg
sudo chmod a+r /usr/share/keyrings/datadog-archive-keyring.gpg
curl https://keys.datadoghq.com/DATADOG_APM_SIGNING_KEYS.public | sudo gpg --no-default-keyring --keyring /usr/share/keyrings/datadog-archive-keyring.gpg --import --batch

# Update and install
sudo apt-get update
sudo apt-get install datadog-agent
```

#### CentOS/RHEL:

```bash
# Add Datadog repository
cat <<EOF | sudo tee /etc/yum.repos.d/datadog.repo
[datadog]
name=Datadog, Inc.
baseurl=https://yum.datadoghq.com/stable/7/x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://keys.datadoghq.com/DATADOG_RPM_KEY_CURRENT.public
       https://keys.datadoghq.com/DATADOG_RPM_KEY_B01082D3.public
       https://keys.datadoghq.com/DATADOG_RPM_KEY_FD4BF915.public
EOF

# Install
sudo yum install datadog-agent
```

## Step 3: Configure the Datadog Agent

Edit the agent configuration file:

```bash
sudo nano /etc/datadog-agent/datadog.yaml
```

Update the following settings:

```yaml
# Your Datadog API key
api_key: YOUR_API_KEY

# Datadog site (use datadoghq.com for US, datadoghq.eu for EU)
site: datadoghq.com

# Enable APM
apm_config:
  enabled: true

  # Allow traces from non-local traffic (important for Docker/containers)
  apm_non_local_traffic: true

  # Port for receiving traces
  receiver_port: 8126

  # Environment tag
  env: production

# Enable process monitoring
process_config:
  enabled: true

# Enable logs (optional but recommended)
logs_enabled: true
logs_config:
  container_collect_all: true
```

Save and exit (Ctrl+X, then Y, then Enter).

## Step 4: Start the Datadog Agent

```bash
# Start the agent
sudo systemctl start datadog-agent

# Enable auto-start on boot
sudo systemctl enable datadog-agent

# Check status
sudo systemctl status datadog-agent

# View agent logs
sudo tail -f /var/log/datadog/agent.log
```

## Step 5: Verify Agent Installation

```bash
# Check agent status
sudo datadog-agent status

# You should see:
# - Agent running
# - APM Agent running on port 8126
# - Connection to Datadog successful
```

## Step 6: Configure Your Daybook Backend

Update your `.env` file:

```bash
# Enable Datadog APM
DD_ENABLED=true

# Service name (will appear in Datadog)
DD_SERVICE=daybook-backend

# Environment (production, staging, development)
DD_ENV=production

# Datadog Agent host (localhost if agent is on same machine)
DD_AGENT_HOST=localhost

# Datadog Agent trace port
DD_TRACE_AGENT_PORT=8126
```

## Step 7: Start Your Application

```bash
cd /path/to/daybook/backend
./daybook-backend
```

You should see in the logs:
```
Initializing Datadog APM tracer...
Datadog APM tracer initialized for service: daybook-backend (env: production)
Datadog APM middleware configured
```

## Step 8: Generate Some Traffic

Make some API requests to your backend:

```bash
# Health check
curl http://localhost:8080/health

# Login (example)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Create some transactions, check accounts, etc.
```

## Step 9: View Traces in Datadog

1. Go to https://app.datadoghq.com/apm/traces
2. Select **Service**: `daybook-backend`
3. Select **Env**: `production` (or whatever you set in DD_ENV)
4. You should see traces appearing within 1-2 minutes

## What Gets Tracked

With the enhanced integration, Datadog will track:

### ✅ HTTP Requests
- All API endpoints
- Request/response times
- Status codes (200, 404, 500, etc.)
- HTTP methods (GET, POST, PUT, DELETE)
- Route patterns

### ✅ Database Operations
- All SQL queries via GORM
- Query execution time
- Table names
- Query types (SELECT, INSERT, UPDATE, DELETE)
- Database errors and slow queries

### ✅ Redis Operations
- All Redis commands
- Command execution time
- Cache hits/misses
- Connection pool metrics

### ✅ Custom Business Metrics
- User registrations
- Transaction creations
- Account operations
- Goal contributions
- Budget tracking

### ✅ Errors & Exceptions
- Stack traces
- Error types
- Error rates
- Failed requests

## Datadog Dashboard Features

Once data is flowing, you can use:

1. **APM → Services**: Overview of all services
2. **APM → Traces**: Individual request traces
3. **APM → Service Map**: Visual service dependencies
4. **Dashboards**: Create custom dashboards
5. **Monitors**: Set up alerts for errors, slow queries, etc.

## Useful Datadog Agent Commands

```bash
# Check agent status
sudo datadog-agent status

# Restart agent
sudo systemctl restart datadog-agent

# Stop agent
sudo systemctl stop datadog-agent

# View logs
sudo tail -f /var/log/datadog/agent.log

# Check APM stats
sudo datadog-agent status | grep -A 20 "APM Agent"

# Test configuration
sudo datadog-agent configcheck
```

## Troubleshooting

### Agent Not Starting

```bash
# Check logs
sudo journalctl -u datadog-agent -n 50 --no-pager

# Verify configuration
sudo datadog-agent configcheck
```

### No Traces Appearing in Datadog

1. **Check agent is receiving traces:**
   ```bash
   sudo datadog-agent status | grep -A 20 "APM Agent"
   ```

2. **Check application logs** for Datadog initialization messages

3. **Verify network connectivity:**
   ```bash
   telnet localhost 8126
   ```

4. **Check firewall rules** (if agent is on different machine)

### High Memory Usage

If the agent uses too much memory, adjust in `/etc/datadog-agent/datadog.yaml`:

```yaml
apm_config:
  max_memory: 500000000  # 500MB
  max_cpu_percent: 50
```

## Production Best Practices

1. **Set appropriate environment tags:**
   ```bash
   DD_ENV=production
   DD_SERVICE=daybook-backend
   DD_VERSION=1.0.0  # Application version
   ```

2. **Use service tags for better filtering:**
   - Add custom tags in your application code
   - Tag by region, datacenter, instance type, etc.

3. **Set up monitors and alerts:**
   - High error rates
   - Slow database queries
   - Memory/CPU spikes
   - Service downtime

4. **Sample traces in high-traffic environments:**
   - Configure sampling rate in agent config to reduce costs
   - Keep 100% error traces, sample successful requests

5. **Secure your API key:**
   - Store in environment variables or secrets manager
   - Rotate keys periodically
   - Use different keys for different environments

## Cost Optimization

- **Free Tier**: 150GB ingested logs/month, 5 hosts
- **Pro Tier**: $15/host/month (billed annually)
- **Use sampling** for high-traffic applications
- **Filter noisy endpoints** (like health checks)

## Additional Resources

- Datadog Documentation: https://docs.datadoghq.com/
- APM Guide: https://docs.datadoghq.com/tracing/
- Go Tracer: https://docs.datadoghq.com/tracing/setup_overview/setup/go/
- Support: https://www.datadoghq.com/support/

## Support

For issues or questions about this setup, check:
1. Datadog Agent logs: `/var/log/datadog/agent.log`
2. Application logs: Check your backend logs
3. Datadog support: https://help.datadoghq.com/
