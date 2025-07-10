#!/bin/bash

# Debug script for MySQL shutdown issues in vttestserver
# Usage: ./debug_mysql_shutdown.sh [mysql_data_dir]

set -e

MYSQL_DATA_DIR=${1:-/vt/vtdataroot}
TABLET_DIR=$(find "$MYSQL_DATA_DIR" -name "vt_*" -type d | head -1)

echo "=== MySQL Shutdown Debug Information ==="
echo "Timestamp: $(date)"
echo "MySQL Data Dir: $MYSQL_DATA_DIR"
echo "Tablet Dir: $TABLET_DIR"
echo

# Check for MySQL processes
echo "=== Active MySQL Processes ==="
ps aux | grep mysql | grep -v grep || echo "No MySQL processes found"
echo

# Check for socket files
echo "=== Socket Files ==="
find "$MYSQL_DATA_DIR" -name "*.sock" -o -name "*.sock.lock" 2>/dev/null || echo "No socket files found"
echo

# Check for PID files
echo "=== PID Files ==="
find "$MYSQL_DATA_DIR" -name "*.pid" 2>/dev/null || echo "No PID files found"
echo

# Check MySQL error log
if [ -n "$TABLET_DIR" ] && [ -f "$TABLET_DIR/error.log" ]; then
    echo "=== MySQL Error Log (last 20 lines) ==="
    tail -20 "$TABLET_DIR/error.log"
    echo
fi

# Check for InnoDB shutdown status
if [ -n "$TABLET_DIR" ] && [ -f "$TABLET_DIR/error.log" ]; then
    echo "=== InnoDB Shutdown Events ==="
    grep -i "shutdown" "$TABLET_DIR/error.log" | tail -10 || echo "No shutdown events found"
    echo
fi

# Check disk space
echo "=== Disk Space ==="
df -h "$MYSQL_DATA_DIR"
echo

# Check memory usage
echo "=== Memory Usage ==="
free -h
echo

# Check for hanging connections
echo "=== Network Connections ==="
netstat -an | grep :3306 || echo "No MySQL connections found"
echo

echo "=== Debug Complete ===" 