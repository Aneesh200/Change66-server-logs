#!/bin/bash

##############################################
# View Notification Logs - Helper Script
##############################################

BASE_URL="http://logs.biopeak.authify.tech/api/v1"
API_KEY="habit-tracker-key-dev"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

##############################################
# 1. All Notification Logs
##############################################

view_all_notifications() {
    print_header "📬 All Notification-Related Logs"
    
    curl -s -H "X-API-Key: $API_KEY" \
      "$BASE_URL/logs/recent?limit=500" | python3 -c "
import sys, json
data = json.load(sys.stdin)
notification_logs = [log for log in data['data'] if 'notification' in log['event_name'].lower() or 'push' in log['event_name'].lower()]
print(f\"Found {len(notification_logs)} notification logs\n\")
for log in notification_logs[:20]:  # Show first 20
    print(f\"ID: {log['id']} | {log['event_name']}\")
    print(f\"  Time: {log['timestamp']}\")
    print(f\"  User: {log['user_id']}\")
    if 'properties' in log:
        if 'success' in log['properties']:
            status = '✅' if log['properties']['success'] else '❌'
            print(f\"  Status: {status} {log['properties'].get('success')}\")
        if 'error' in log['properties']:
            print(f\"  Error: {log['properties']['error'][:100]}\")
    print()
"
}

##############################################
# 2. FCM Token Generation
##############################################

view_fcm_token_logs() {
    print_header "🔑 FCM Token Generation Logs"
    
    curl -s -H "X-API-Key: $API_KEY" \
      "$BASE_URL/logs/recent?limit=100" | python3 -c "
import sys, json
data = json.load(sys.stdin)
fcm_logs = [log for log in data['data'] if 'fcm_token' in log['event_name'].lower()]
print(f\"Found {len(fcm_logs)} FCM token logs\n\")
for log in fcm_logs:
    print(f\"ID: {log['id']} | {log['event_name']}\")
    print(f\"  Time: {log['timestamp']}\")
    if 'properties' in log:
        if 'success' in log['properties']:
            status = '✅' if log['properties']['success'] else '❌'
            print(f\"  Status: {status}\")
        if 'generation_duration_ms' in log['properties']:
            print(f\"  Duration: {log['properties']['generation_duration_ms']}ms\")
        if 'token_length' in log['properties']:
            print(f\"  Token Length: {log['properties']['token_length']}\")
    print()
"
}

##############################################
# 3. Token Store Failures
##############################################

view_token_store_failures() {
    print_header "❌ Token Store Failures"
    
    curl -s -H "X-API-Key: $API_KEY" \
      "$BASE_URL/logs/recent?limit=100" | python3 -c "
import sys, json
data = json.load(sys.stdin)
failure_logs = [log for log in data['data'] if 'token_store_failed' in log['event_name'].lower()]
print(f\"Found {len(failure_logs)} token store failure logs\n\")
for log in failure_logs:
    print(f\"ID: {log['id']} | {log['timestamp']}\")
    if 'properties' in log:
        if 'dashboard_url' in log['properties']:
            print(f\"  Trying to reach: {log['properties']['dashboard_url']}\")
        if 'error' in log['properties']:
            error = log['properties']['error']
            # Extract key error info
            if 'Connection refused' in error:
                print(f\"  Issue: Dashboard not running (Connection refused)\")
            elif 'No route to host' in error:
                print(f\"  Issue: Network unreachable (No route to host)\")
            else:
                print(f\"  Error: {error[:150]}\")
    print()
"
}

##############################################
# 4. Habit Completion Notifications
##############################################

view_habit_notifications() {
    print_header "✅ Habit Completion Notifications"
    
    curl -s -H "X-API-Key: $API_KEY" \
      "$BASE_URL/logs/recent?limit=100" | python3 -c "
import sys, json
data = json.load(sys.stdin)
habit_logs = [log for log in data['data'] if 'habit_completion' in log['event_name'].lower()]
print(f\"Found {len(habit_logs)} habit completion notification logs\n\")
for log in habit_logs:
    print(f\"ID: {log['id']} | {log['event_name']}\")
    print(f\"  Time: {log['timestamp']}\")
    if 'properties' in log:
        if 'habit_title' in log['properties']:
            print(f\"  Habit: {log['properties']['habit_title']}\")
        if 'current_streak' in log['properties']:
            print(f\"  Streak: {log['properties']['current_streak']}\")
        if 'processing_duration_ms' in log['properties']:
            print(f\"  Processing: {log['properties']['processing_duration_ms']}ms\")
        if 'success' in log['properties']:
            status = '✅' if log['properties']['success'] else '❌'
            print(f\"  Status: {status}\")
    print()
"
}

##############################################
# 5. Summary Statistics
##############################################

view_notification_stats() {
    print_header "📊 Notification Statistics"
    
    curl -s -H "X-API-Key: $API_KEY" \
      "$BASE_URL/logs/recent?limit=500" | python3 -c "
import sys, json
from collections import Counter

data = json.load(sys.stdin)
notification_logs = [log for log in data['data'] if 'notification' in log['event_name'].lower() or 'push' in log['event_name'].lower()]

print(f\"Total Notification Logs: {len(notification_logs)}\n\")

# Event type breakdown
event_types = Counter([log['event_name'] for log in notification_logs])
print(\"Event Types:\")
for event_type, count in event_types.most_common():
    print(f\"  {event_type}: {count}\")

print()

# Success rate
success_logs = [log for log in notification_logs if log.get('properties', {}).get('success') == True]
failed_logs = [log for log in notification_logs if log.get('properties', {}).get('success') == False]
print(f\"Success Rate: {len(success_logs)}/{len(success_logs) + len(failed_logs)} ({len(success_logs)*100/(len(success_logs)+len(failed_logs)) if (len(success_logs)+len(failed_logs)) > 0 else 0:.1f}%)\")

print()

# Recent errors
print(\"Recent Errors:\")
error_logs = [log for log in notification_logs if 'error' in log.get('properties', {})][:5]
for log in error_logs:
    error = log['properties']['error']
    if 'Connection refused' in error:
        print(f\"  • Dashboard not running (localhost)\")
    elif 'No route to host' in error:
        print(f\"  • Network unreachable (192.168.x.x)\")
    else:
        print(f\"  • {error[:80]}...\")
"
}

##############################################
# Main Menu
##############################################

case "${1:-all}" in
    all)
        view_all_notifications
        ;;
    fcm)
        view_fcm_token_logs
        ;;
    failures)
        view_token_store_failures
        ;;
    habits)
        view_habit_notifications
        ;;
    stats)
        view_notification_stats
        ;;
    *)
        echo "Usage: $0 [all|fcm|failures|habits|stats]"
        echo ""
        echo "Commands:"
        echo "  all       - View all notification logs (default)"
        echo "  fcm       - View FCM token generation logs"
        echo "  failures  - View token store failures"
        echo "  habits    - View habit completion notifications"
        echo "  stats     - View notification statistics"
        echo ""
        echo "Examples:"
        echo "  $0 all"
        echo "  $0 fcm"
        echo "  $0 stats"
        exit 1
        ;;
esac

