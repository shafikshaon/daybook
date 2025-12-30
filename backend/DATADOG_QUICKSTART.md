# Datadog APM Quick Start Guide

Get Datadog monitoring running on your Linux machine in 5 minutes.

## Step 1: Get Your API Key (1 minute)

1. Sign up at https://www.datadoghq.com/ (free trial available)
2. Go to https://app.datadoghq.com/organization-settings/api-keys
3. Copy your API key

## Step 2: Install Datadog Agent (2 minutes)

Replace `YOUR_API_KEY` with your actual API key:

```bash
DD_API_KEY=f3ba6830d00c5b5d38a424a1da87affb \
DD_SITE="datadoghq.com" \
DD_APM_ENABLED=true \
bash -c "$(curl -L https://s3.amazonaws.com/dd-agent/scripts/install_script.sh)"
```

Wait for installation to complete, then start the agent:

```bash
sudo systemctl start datadog-agent
sudo systemctl enable datadog-agent
```

Verify it's running:

```bash
sudo datadog-agent status | grep "Agent (.*) is running"
```

You should see: `Agent (v7.x.x) is running`

## Step 3: Enable APM in Your Backend (1 minute)

Edit your `.env` file:

```bash
nano .env
```

Change these values:

```bash
DD_ENABLED=true              # Changed from false
DD_SERVICE=daybook-backend
DD_ENV=production           # Or development/staging
DD_AGENT_HOST=localhost
DD_TRACE_AGENT_PORT=8126
```

Save and exit (Ctrl+X, Y, Enter).

## Step 4: Start Your Application (30 seconds)

```bash
./daybook-backend
```

Look for these lines in the logs:

```
Datadog APM tracer initialized for service: daybook-backend (env: production) ✓
Datadog GORM tracing enabled ✓
Datadog Redis tracing enabled ✓
Datadog APM middleware configured ✓
```

## Step 5: View Traces in Datadog (30 seconds)

1. Make some API requests:
   ```bash
   curl http://localhost:8080/health
   curl http://localhost:8080/api/v1/transactions
   ```

2. Go to https://app.datadoghq.com/apm/traces

3. Select **Service**: `daybook-backend`

4. You should see traces appearing within 1-2 minutes!

## What You'll See

### Service Overview
- Request rate
- Error rate
- Latency (P50, P75, P95, P99)
- Throughput

### Individual Traces
Every request shows:
- Total duration
- HTTP method and endpoint
- Database queries executed
- Redis commands (if any)
- Status code
- User ID (if authenticated)

### Database Insights
- Slowest queries
- Most frequent queries
- Query execution time
- Connection pool stats

### Service Map
Visual representation of:
- Your API service
- PostgreSQL database
- Redis cache
- External dependencies

## Troubleshooting

### No traces appearing?

**Check agent status:**
```bash
sudo datadog-agent status
```

Look for:
```
APM Agent
=========
  Status: Running
  Pid: xxxxx
  Receiver: localhost:8126
```

**Check application logs:**
```bash
# Should see these messages
grep -i "datadog" /path/to/your/logs
```

**Test agent connectivity:**
```bash
telnet localhost 8126
```

### Agent not starting?

**Check logs:**
```bash
sudo journalctl -u datadog-agent -n 50
```

**Verify API key:**
```bash
sudo datadog-agent config | grep api_key
```

## Next Steps

1. **Read the full guides:**
   - [DATADOG_SETUP.md](DATADOG_SETUP.md) - Complete Linux setup
   - [MONITORING_GUIDE.md](MONITORING_GUIDE.md) - What gets tracked

2. **Set up dashboards:**
   - Go to https://app.datadoghq.com/dashboard/lists
   - Create custom dashboards for your KPIs

3. **Configure alerts:**
   - Go to https://app.datadoghq.com/monitors/create
   - Set up monitors for:
     - Error rate > 1%
     - P95 latency > 1000ms
     - Database slow queries > 500ms

4. **Use custom metrics:**
   - See [MONITORING_GUIDE.md](MONITORING_GUIDE.md) section 4
   - Track business events like user registrations, transactions

## Cost Information

- **Free Tier**: 150GB ingested logs/month, 5 hosts
- **Pro Tier**: $15/host/month (APM included)
- **Enterprise**: Custom pricing

For most small deployments, the free tier is sufficient for initial monitoring.

## Support

- Datadog Docs: https://docs.datadoghq.com/
- APM Guide: https://docs.datadoghq.com/tracing/
- Support: https://help.datadoghq.com/

---

**That's it! Your application is now fully monitored. 🎉**
