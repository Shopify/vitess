# MySQL Shutdown Fix for vttestserver

## Problem Description

vttestserver was experiencing cases where mysqld wasn't shutdown properly when receiving SIGTERM signals. This led to:

- Orphaned MySQL processes
- Stale socket files
- Data corruption risks
- Container restart failures

## Root Causes

1. **Timeout Mismatch**: vttest layer timeout (60s) was shorter than MySQL shutdown timeout (300s)
2. **Poor Error Handling**: Shutdown failures were logged but not properly handled
3. **Missing Cleanup**: No cleanup of orphaned processes or stale files
4. **Fixed Timeouts**: No way to configure timeouts for different MySQL configurations

## Solutions Implemented

### 1. Extended and Configurable Timeouts

**Before:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
```

**After:**
```go
// Configurable timeout with sensible default
ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
```

### 2. Enhanced Error Handling

**New Features:**
- Force shutdown attempts if graceful shutdown fails
- Process cleanup and signal handling (SIGTERM → SIGKILL)
- Stale socket and PID file removal
- Better error reporting with timeout information

### 3. Configuration Options

**New CLI Flag:**
```bash
--mysql_shutdown_timeout=6m  # Default: 6 minutes
```

**Environment Variable:**
```bash
MYSQL_SHUTDOWN_TIMEOUT=6m
```

## Usage Instructions

### For Docker Deployments

1. **Update your docker-compose.yml or run command:**
```yaml
environment:
  - MYSQL_SHUTDOWN_TIMEOUT=6m  # Adjust based on your MySQL setup
```

2. **For large MySQL instances with big buffer pools:**
```yaml
environment:
  - MYSQL_SHUTDOWN_TIMEOUT=10m  # Longer timeout for large datasets
```

### For Direct vttestserver Usage

```bash
vttestserver \
  --mysql_shutdown_timeout=6m \
  --keyspaces=test_keyspace \
  --num_shards=2
```

## Troubleshooting

### 1. Use the Debug Script

```bash
# Run after shutdown issues
./tools/debug_mysql_shutdown.sh /path/to/mysql/data
```

### 2. Check MySQL Error Logs

```bash
# Look for shutdown-related messages
tail -f /vt/vtdataroot/vt_*/error.log | grep -i shutdown
```

### 3. Monitor Shutdown Process

```bash
# In another terminal during shutdown
watch 'ps aux | grep mysql'
```

### 4. Common Issues and Solutions

**Issue: "Context deadline exceeded"**
- **Solution**: Increase `--mysql_shutdown_timeout`
- **Example**: `--mysql_shutdown_timeout=10m`

**Issue: Stale socket files**
- **Solution**: The fix now automatically cleans these up
- **Manual**: `rm /vt/vtdataroot/vt_*/mysql.sock*`

**Issue: Orphaned MySQL processes**
- **Solution**: The fix now sends SIGTERM then SIGKILL
- **Manual**: `killall mysqld` or specific PID kill

**Issue: Large InnoDB buffer pool taking time**
- **Solution**: Use longer timeout or set `innodb_fast_shutdown=1`
- **Example**: `--mysql_shutdown_timeout=15m`

## Performance Tuning

### MySQL Configuration Recommendations

Add to your MySQL configuration to speed up shutdown:

```ini
# my.cnf additions for faster shutdown
innodb_fast_shutdown=1          # Skip full purge and buffer pool flush
innodb_flush_method=O_DIRECT    # Reduce OS cache interference
sync_binlog=0                   # Reduce binlog sync during shutdown
```

### Container Resource Limits

Ensure adequate resources for graceful shutdown:

```yaml
# docker-compose.yml
services:
  vttestserver:
    deploy:
      resources:
        limits:
          memory: 1G
          cpus: 2
    # Ensure sufficient time for shutdown
    stop_grace_period: 10m
```

## Testing the Fix

### 1. Test Normal Shutdown

```bash
# Start vttestserver
docker run -d --name vttestserver vitess/vttestserver

# Send SIGTERM and verify clean shutdown
docker stop vttestserver

# Check logs for clean shutdown
docker logs vttestserver | grep -i shutdown
```

### 2. Test Force Shutdown Scenario

```bash
# Simulate a hung MySQL by sending SIGSTOP to mysqld
docker exec vttestserver killall -STOP mysqld

# Try to stop container - should trigger force cleanup
timeout 30s docker stop vttestserver

# Verify cleanup occurred
./tools/debug_mysql_shutdown.sh
```

## Migration Guide

### Updating Existing Deployments

1. **Update vttestserver image** to version with the fix
2. **Add environment variable**:
   ```bash
   MYSQL_SHUTDOWN_TIMEOUT=6m
   ```
3. **Increase container stop_grace_period**:
   ```yaml
   stop_grace_period: 10m
   ```
4. **Test shutdown behavior** in staging environment

### Backward Compatibility

- Default timeout increased from 60s to 6m (safe increase)
- Existing deployments will automatically benefit from improved error handling
- No breaking changes to existing APIs

## Monitoring and Alerts

### Metrics to Monitor

1. **Shutdown Duration**: How long MySQL takes to shut down
2. **Force Shutdown Events**: When graceful shutdown fails
3. **Orphaned Processes**: Processes left behind after container stop
4. **Socket File Issues**: Stale socket files requiring cleanup

### Sample Monitoring Script

```bash
#!/bin/bash
# monitor_mysql_shutdown.sh

LOG_FILE="/var/log/mysql_shutdown_monitor.log"

check_mysql_shutdown() {
    MYSQLD_PIDS=$(pgrep mysqld || true)
    SOCKET_FILES=$(find /vt -name "*.sock" 2>/dev/null || true)
    
    if [ -n "$MYSQLD_PIDS" ] || [ -n "$SOCKET_FILES" ]; then
        echo "$(date): WARNING - MySQL shutdown incomplete" >> "$LOG_FILE"
        echo "PIDs: $MYSQLD_PIDS" >> "$LOG_FILE"
        echo "Sockets: $SOCKET_FILES" >> "$LOG_FILE"
        # Send alert to monitoring system
        curl -X POST "$ALERT_WEBHOOK" -d "MySQL shutdown incomplete"
    fi
}

check_mysql_shutdown
```

## Related Pull Requests

The fixes implemented are based on lessons learned from these Vitess improvements:

- [PR #14568](https://github.com/Shopify/vitess/pull/14568) - feat: Allow configurable MySQL shutdown timeout
- [PR #15575](https://github.com/Shopify/vitess/pull/15575) - Refactor: Adjust mysqlctld onterm_timeout default

These PRs addressed similar timeout issues in other Vitess components and provided the foundation for this vttestserver fix. 