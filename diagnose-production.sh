#!/bin/bash

# Production Backend Diagnostics Script
# Run this on your production server to diagnose API issues

echo "======================================"
echo "Daybook Backend Diagnostics"
echo "======================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if backend process is running
echo "1. Checking Backend Process..."
if ps aux | grep -v grep | grep daybook-backend > /dev/null; then
    echo -e "${GREEN}✓ Backend process is running${NC}"
    ps aux | grep -v grep | grep daybook-backend
else
    echo -e "${RED}✗ Backend process is NOT running${NC}"
    echo "  → Start the backend service"
fi
echo ""

# Check what's listening on port 8080
echo "2. Checking Port 8080..."
if sudo lsof -i :8080 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Port 8080 is in use${NC}"
    sudo lsof -i :8080
else
    echo -e "${RED}✗ Nothing listening on port 8080${NC}"
    echo "  → Backend is not running or using different port"
fi
echo ""

# Test backend locally
echo "3. Testing Backend Locally (http://localhost:8080/health)..."
LOCAL_TEST=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health 2>/dev/null)
if [ "$LOCAL_TEST" = "200" ]; then
    echo -e "${GREEN}✓ Backend responds locally (HTTP $LOCAL_TEST)${NC}"
    curl -s http://localhost:8080/health | head -5
else
    echo -e "${RED}✗ Backend not responding locally (HTTP $LOCAL_TEST)${NC}"
    echo "  → Backend is not running properly"
fi
echo ""

# Check Nginx status
echo "4. Checking Nginx Status..."
if systemctl is-active --quiet nginx; then
    echo -e "${GREEN}✓ Nginx is running${NC}"
else
    echo -e "${RED}✗ Nginx is not running${NC}"
    echo "  → Start nginx: sudo systemctl start nginx"
fi
echo ""

# Test Nginx configuration
echo "5. Testing Nginx Configuration..."
if sudo nginx -t > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Nginx configuration is valid${NC}"
else
    echo -e "${RED}✗ Nginx configuration has errors${NC}"
    sudo nginx -t 2>&1
fi
echo ""

# Check if nginx config for API exists
echo "6. Looking for Nginx API Configuration..."
if [ -f /etc/nginx/sites-available/api.daybook.shafik.xyz ]; then
    echo -e "${GREEN}✓ Found: /etc/nginx/sites-available/api.daybook.shafik.xyz${NC}"
    echo "  Checking if enabled..."
    if [ -L /etc/nginx/sites-enabled/api.daybook.shafik.xyz ]; then
        echo -e "  ${GREEN}✓ Configuration is enabled${NC}"
    else
        echo -e "  ${RED}✗ Configuration is not enabled${NC}"
        echo "    → Run: sudo ln -s /etc/nginx/sites-available/api.daybook.shafik.xyz /etc/nginx/sites-enabled/"
    fi
elif [ -f /etc/nginx/conf.d/api.daybook.shafik.xyz.conf ]; then
    echo -e "${GREEN}✓ Found: /etc/nginx/conf.d/api.daybook.shafik.xyz.conf${NC}"
else
    echo -e "${RED}✗ No nginx configuration found for api.daybook.shafik.xyz${NC}"
    echo "  → Create nginx configuration (see PRODUCTION_DEPLOYMENT_GUIDE.md)"
fi
echo ""

# Test public API endpoint
echo "7. Testing Public API (https://api.daybook.shafik.xyz/health)..."
PUBLIC_TEST=$(curl -s -o /dev/null -w "%{http_code}" https://api.daybook.shafik.xyz/health 2>/dev/null)
if [ "$PUBLIC_TEST" = "200" ]; then
    echo -e "${GREEN}✓ Public API responds (HTTP $PUBLIC_TEST)${NC}"
    curl -s https://api.daybook.shafik.xyz/health | head -5
else
    echo -e "${RED}✗ Public API not responding correctly (HTTP $PUBLIC_TEST)${NC}"
    if [ "$PUBLIC_TEST" = "000" ]; then
        echo "  → DNS or connectivity issue"
    elif [ "$PUBLIC_TEST" = "404" ]; then
        echo "  → Nginx routing issue or backend not running"
    elif [ "$PUBLIC_TEST" = "502" ]; then
        echo "  → Backend not reachable from nginx"
    elif [ "$PUBLIC_TEST" = "503" ]; then
        echo "  → Backend temporarily unavailable"
    fi
fi
echo ""

# Check SSL certificate
echo "8. Checking SSL Certificate..."
if [ -f /etc/letsencrypt/live/api.daybook.shafik.xyz/fullchain.pem ]; then
    echo -e "${GREEN}✓ SSL certificate found${NC}"
    echo "  Expiry:"
    sudo openssl x509 -enddate -noout -in /etc/letsencrypt/live/api.daybook.shafik.xyz/fullchain.pem
else
    echo -e "${YELLOW}! SSL certificate not found at expected location${NC}"
    echo "  → Install certificate: sudo certbot --nginx -d api.daybook.shafik.xyz"
fi
echo ""

# Check firewall
echo "9. Checking Firewall..."
if command -v ufw > /dev/null; then
    if sudo ufw status | grep -q "Status: active"; then
        echo -e "${YELLOW}! UFW firewall is active${NC}"
        echo "  Checking if HTTP/HTTPS are allowed..."
        if sudo ufw status | grep -q "443"; then
            echo -e "  ${GREEN}✓ Port 443 (HTTPS) is allowed${NC}"
        else
            echo -e "  ${RED}✗ Port 443 (HTTPS) is not allowed${NC}"
            echo "    → Run: sudo ufw allow 443/tcp"
        fi
        if sudo ufw status | grep -q "80"; then
            echo -e "  ${GREEN}✓ Port 80 (HTTP) is allowed${NC}"
        else
            echo -e "  ${YELLOW}! Port 80 (HTTP) is not allowed${NC}"
            echo "    → Run: sudo ufw allow 80/tcp"
        fi
    else
        echo -e "${GREEN}✓ UFW firewall is inactive${NC}"
    fi
else
    echo "  UFW not found, checking iptables..."
    if sudo iptables -L -n | grep -q "443"; then
        echo -e "  ${GREEN}✓ Port 443 appears to be allowed${NC}"
    else
        echo -e "  ${YELLOW}! Cannot determine firewall status${NC}"
    fi
fi
echo ""

# Check backend .env file
echo "10. Checking Backend Configuration..."
if [ -f ./backend/.env ]; then
    echo -e "${GREEN}✓ Found backend/.env${NC}"
    echo "  Critical settings:"
    grep -E "^SERVER_PORT=" ./backend/.env || echo "  ${RED}✗ SERVER_PORT not set${NC}"
    grep -E "^SERVER_MODE=" ./backend/.env || echo "  ${YELLOW}! SERVER_MODE not set${NC}"
    grep -E "^CORS_ALLOWED_ORIGINS=" ./backend/.env || echo "  ${RED}✗ CORS_ALLOWED_ORIGINS not set${NC}"
elif [ -f ./.env ]; then
    echo -e "${GREEN}✓ Found .env${NC}"
    echo "  Critical settings:"
    grep -E "^SERVER_PORT=" ./.env || echo "  ${RED}✗ SERVER_PORT not set${NC}"
    grep -E "^SERVER_MODE=" ./.env || echo "  ${YELLOW}! SERVER_MODE not set${NC}"
    grep -E "^CORS_ALLOWED_ORIGINS=" ./.env || echo "  ${RED}✗ CORS_ALLOWED_ORIGINS not set${NC}"
else
    echo -e "${RED}✗ Backend .env file not found${NC}"
    echo "  → Create .env file with production settings"
fi
echo ""

# Check recent nginx error logs
echo "11. Recent Nginx Error Logs (last 5 lines)..."
if [ -f /var/log/nginx/error.log ]; then
    sudo tail -5 /var/log/nginx/error.log
else
    echo -e "${YELLOW}! No nginx error log found${NC}"
fi
echo ""

# Summary
echo "======================================"
echo "Diagnosis Summary"
echo "======================================"
echo ""

ISSUES=0

if ! ps aux | grep -v grep | grep daybook-backend > /dev/null; then
    echo -e "${RED}✗ Backend is not running${NC}"
    ((ISSUES++))
fi

if [ "$LOCAL_TEST" != "200" ]; then
    echo -e "${RED}✗ Backend not responding locally${NC}"
    ((ISSUES++))
fi

if ! systemctl is-active --quiet nginx; then
    echo -e "${RED}✗ Nginx is not running${NC}"
    ((ISSUES++))
fi

if [ "$PUBLIC_TEST" != "200" ]; then
    echo -e "${RED}✗ Public API not accessible (HTTP $PUBLIC_TEST)${NC}"
    ((ISSUES++))
fi

if [ $ISSUES -eq 0 ]; then
    echo -e "${GREEN}✓ No critical issues found!${NC}"
    echo ""
    echo "If you're still experiencing problems:"
    echo "  1. Check browser console for errors"
    echo "  2. Verify CORS settings in backend .env"
    echo "  3. Clear browser cache and retry"
else
    echo -e "${RED}Found $ISSUES critical issue(s)${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Review the issues above"
    echo "  2. See PRODUCTION_DEPLOYMENT_GUIDE.md for solutions"
    echo "  3. Re-run this script after making changes"
fi

echo ""
echo "======================================"
