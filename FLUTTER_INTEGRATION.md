# Flutter App Integration Guide

This guide shows how to update your Flutter app to use the new Go log ingestion server.

## 1. Update Analytics Service Configuration

In your `lib/services/analytics_service.dart`, update the server URL and add the API key:

```dart
class AnalyticsService {
  // Update the server URL to point to your Go server
  static const String LOG_SERVER_URL = 'https://your-domain.com/api/v1';
  // OR for local development:
  // static const String LOG_SERVER_URL = 'http://localhost:8080/api/v1';
  
  // Add API key configuration
  static const String API_KEY = 'your-api-key-here';
  
  // ... rest of your existing code
}
```

## 2. Update HTTP Headers

Modify your HTTP requests to include the API key header:

```dart
// In _sendSingleEvent method
final response = await http.post(
  Uri.parse('$LOG_SERVER_URL/ingest'),
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': API_KEY,  // Add this line
  },
  body: jsonEncode(event),
).timeout(const Duration(seconds: 5));

// In _flushBuffer method  
final response = await http.post(
  Uri.parse('$LOG_SERVER_URL/batch-ingest'),
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': API_KEY,  // Add this line
  },
  body: jsonEncode({'logs': batch}),
).timeout(const Duration(seconds: 10));

// In _syncOfflineLogs method
final response = await http.post(
  Uri.parse('$LOG_SERVER_URL/batch-ingest'),
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': API_KEY,  // Add this line
  },
  body: jsonEncode({'logs': batch}),
).timeout(const Duration(seconds: 15));
```

## 3. Environment-Based Configuration

For better security, use environment-based configuration:

### Option A: Using flutter_dotenv

1. Add to `pubspec.yaml`:
```yaml
dependencies:
  flutter_dotenv: ^5.1.0
```

2. Create `.env` file in your Flutter project root:
```env
LOG_SERVER_URL=https://your-domain.com/api/v1
LOG_SERVER_API_KEY=your-api-key-here
```

3. Update `analytics_service.dart`:
```dart
import 'package:flutter_dotenv/flutter_dotenv.dart';

class AnalyticsService {
  static String get LOG_SERVER_URL => 
    dotenv.env['LOG_SERVER_URL'] ?? 'http://localhost:8080/api/v1';
  
  static String get API_KEY => 
    dotenv.env['LOG_SERVER_API_KEY'] ?? '';
  
  // ... rest of code
}
```

4. Load environment in `main.dart`:
```dart
Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await dotenv.load(fileName: ".env");
  runApp(MyApp());
}
```

### Option B: Using Build Configurations

1. Create different configuration files:
   - `lib/config/dev_config.dart`
   - `lib/config/prod_config.dart`

2. Development config (`dev_config.dart`):
```dart
class Config {
  static const String LOG_SERVER_URL = 'http://localhost:8080/api/v1';
  static const String API_KEY = 'habit-tracker-key-dev';
}
```

3. Production config (`prod_config.dart`):
```dart
class Config {
  static const String LOG_SERVER_URL = 'https://your-domain.com/api/v1';
  static const String API_KEY = 'your-production-api-key';
}
```

4. Use conditional imports in `analytics_service.dart`:
```dart
import 'config/dev_config.dart' if (dart.library.io) 'config/prod_config.dart';

class AnalyticsService {
  static String get LOG_SERVER_URL => Config.LOG_SERVER_URL;
  static String get API_KEY => Config.API_KEY;
  // ... rest of code
}
```

## 4. API Key Security

### Generate Your API Keys

Run the setup script to generate secure API keys:
```bash
cd Change66-Log-Server
./scripts/setup.sh
```

This will generate two API keys:
- Development key: For testing and development
- Production key: For your live app

### Store Keys Securely

**Never commit API keys to version control!**

1. Add `.env` to your `.gitignore`
2. Use CI/CD environment variables for production
3. Consider using Flutter's secure storage for runtime key management

### Key Rotation

The Go server supports automatic key rotation. Update your Flutter app configuration when keys are rotated.

## 5. Error Handling Updates

Update error handling to work with the new server responses:

```dart
Future<void> _sendSingleEvent(Map<String, dynamic> event) async {
  try {
    final response = await http.post(
      Uri.parse('$LOG_SERVER_URL/ingest'),
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': API_KEY,
      },
      body: jsonEncode(event),
    ).timeout(const Duration(seconds: 5));

    if (response.statusCode == 201) {
      // Success - event ingested
      if (kDebugMode) {
        final responseData = jsonDecode(response.body);
        print('✅ Event ingested: ${responseData['data']['event_id']}');
      }
    } else if (response.statusCode == 401) {
      // Invalid API key
      if (kDebugMode) {
        print('❌ Invalid API key - check configuration');
      }
      await _storeOffline(event);
    } else if (response.statusCode == 429) {
      // Rate limited
      if (kDebugMode) {
        print('⚠️ Rate limited - will retry later');
      }
      await _storeOffline(event);
    } else {
      // Other error
      if (kDebugMode) {
        print('❌ Server error ${response.statusCode}: ${response.body}');
      }
      await _storeOffline(event);
    }
  } catch (e) {
    await _storeOffline(event);
    if (kDebugMode) {
      print('Failed to send single event: $e');
    }
  }
}
```

## 6. Testing Integration

### Local Testing

1. Start the Go server locally:
```bash
cd Change66-Log-Server
make run
```

2. Update your Flutter app to use `http://localhost:8080/api/v1`

3. Test the integration:
```bash
cd Change66-Log-Server
./scripts/test-api.sh
```

### Production Testing

1. Deploy your Go server to your production environment
2. Update Flutter app configuration with production URL and API key
3. Monitor logs to ensure events are being received

## 7. Monitoring Integration

Add monitoring to your Flutter app:

```dart
// Add this method to AnalyticsService
Future<Map<String, dynamic>> getServerStatus() async {
  try {
    final response = await http.get(
      Uri.parse('$LOG_SERVER_URL/status'),
      headers: {'X-API-Key': API_KEY},
    ).timeout(const Duration(seconds: 5));

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    }
  } catch (e) {
    if (kDebugMode) {
      print('Failed to get server status: $e');
    }
  }
  return {};
}

// Use this for debugging/admin screens
Future<void> checkServerHealth() async {
  final status = await getServerStatus();
  if (status.isNotEmpty) {
    print('📊 Server Status: ${status['data']['service']['uptime']}');
    print('📈 Total Logs: ${status['data']['metrics']['total_logs']}');
  }
}
```

## 8. Performance Considerations

### Batch Size Optimization

The Go server supports up to 1000 events per batch. Update your Flutter app:

```dart
class AnalyticsService {
  static const int MAX_BATCH_SIZE = 1000; // Increased from 50
  // ... rest of code
}
```

### Connection Pooling

The Go server handles connection pooling automatically. No changes needed in Flutter.

### Retry Logic

Implement exponential backoff for failed requests:

```dart
Future<void> _sendWithRetry(Function sendFunction, int maxRetries) async {
  int retries = 0;
  while (retries < maxRetries) {
    try {
      await sendFunction();
      return; // Success
    } catch (e) {
      retries++;
      if (retries < maxRetries) {
        await Future.delayed(Duration(seconds: math.pow(2, retries).toInt()));
      }
    }
  }
}
```

## 9. Migration Checklist

- [ ] Update server URL configuration
- [ ] Add API key configuration
- [ ] Update HTTP headers in all requests
- [ ] Update error handling for new response codes
- [ ] Test with local Go server
- [ ] Deploy Go server to production
- [ ] Update production configuration
- [ ] Monitor logs and metrics
- [ ] Update CI/CD pipelines if needed
- [ ] Document new configuration for team

## 10. Troubleshooting

### Common Issues

1. **401 Unauthorized**: Check API key configuration
2. **429 Too Many Requests**: Implement retry logic with backoff
3. **Connection Refused**: Verify server URL and network connectivity
4. **SSL/TLS Errors**: Ensure proper HTTPS configuration

### Debug Mode

Add debug logging to track requests:

```dart
if (kDebugMode) {
  print('📤 Sending to: $LOG_SERVER_URL');
  print('🔑 Using API key: ${API_KEY.substring(0, 8)}...');
  print('📊 Event type: ${event['event_type']}');
}
```

### Server Logs

Monitor Go server logs for debugging:
```bash
docker-compose logs -f log-server
```

## Support

If you encounter issues:
1. Check the Go server logs
2. Verify API key configuration  
3. Test with the provided test script
4. Check network connectivity
5. Review this integration guide
