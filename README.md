# Notification Service

A Go-based notification service that consumes messages from Kafka, stores them in Supabase, and sends notifications via email (SendGrid) and Telegram.

## Features

- Subscribes to Kafka topic "notifications" and processes incoming messages
- Stores notifications in a Supabase database table
- Sends email notifications using SendGrid
- Sends Telegram notifications using the Telegram Bot API

## Prerequisites

- Go 1.21 or higher
- Kafka server
- Supabase account and project
- SendGrid account and API key (for email notifications)
- Telegram Bot API token (for Telegram notifications)

## Installation

1. Clone this repository:

```
git clone https://github.com/yourusername/notification_service.git
cd notification_service
```

2. Install dependencies:

```
go mod download
```

## Configuration

Create a `.env` file in the root directory with the following variables:

```
# Kafka configuration
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
KAFKA_TOPIC=notifications

# Supabase configuration
SUPABASE_URL=https://your-supabase-project.supabase.co
SUPABASE_API_KEY=your-supabase-api-key
SUPABASE_NOTIFICATIONS_TABLE=notifications

# SendGrid configuration
SENDGRID_API_KEY=your-sendgrid-api-key
SENDGRID_FROM_EMAIL=your-sender-email@example.com
SENDGRID_FROM_NAME=Notification Service

# Telegram configuration
TELEGRAM_BOT_TOKEN=your-telegram-bot-token
```

## Supabase Setup

1. Create a new Supabase project
2. Create a `notifications` table with the following schema:

```sql
CREATE TABLE notifications (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  user_id VARCHAR NOT NULL,
  type VARCHAR NOT NULL,
  channel VARCHAR NOT NULL,
  subject VARCHAR,
  content TEXT NOT NULL,
  status VARCHAR NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  sent_at TIMESTAMP WITH TIME ZONE,
  metadata JSONB
);
```

## Running the Service

### Locally

Build and run the application:

```
go build -o notification_service ./cmd/api
./notification_service
```

Or simply:

```
go run ./cmd/api
```

### Using Docker

Build and start the service using Docker Compose:

```
docker-compose up -d
```

This will start:
- Zookeeper
- Kafka
- The notification service

To stop the services:

```
docker-compose down
```

## Testing

### Using the Test Script

A test script is provided in the `scripts` directory to send test messages to Kafka:

```
go run ./scripts/send_test_message.go -channel your-email@example.com
```

Options:
- `-type`: Type of notification (`email` or `telegram`, default: `email`)
- `-channel`: Email address or Telegram chat ID (required)
- `-subject`: Notification subject (default: "Test Notification")
- `-content`: Notification content (default: "This is a test notification from Kafka.")
- `-user`: User ID (default: "user-123")
- `-bootstrap-servers`: Kafka bootstrap servers (default: value from .env or "localhost:9092")
- `-topic`: Kafka topic (default: value from .env or "notifications")

Example for sending a Telegram notification:

```
go run ./scripts/send_test_message.go -type telegram -channel 123456789 -content "Test message from Go!"
```

### Using Kafka Producer Tools

To test the service, you can send a message to the Kafka topic using a Kafka producer tool or a test client. Example message:

```json
{
  "user_id": "user-123",
  "type": "email",
  "channel": "user@example.com",
  "subject": "Test Notification",
  "content": "This is a test notification from Kafka."
}
```

## Kafka Message Format

The service expects Kafka messages in the following JSON format:

```json
{
  "user_id": "user-123",
  "type": "email", // or "telegram"
  "channel": "user@example.com", // or Telegram chat ID
  "subject": "Notification Subject",
  "content": "This is the notification content.",
  "metadata": {
    // optional additional data
  }
}
```

## License

MIT 