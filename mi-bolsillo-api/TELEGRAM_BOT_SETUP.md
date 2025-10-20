# Telegram Bot Setup Guide

This guide explains how to set up and use the Telegram bot for Mi Bolsillo API.

## Features

The Telegram bot provides the following capabilities:

1. **Upload Bills**: Send a photo of a bill/receipt and the bot will automatically parse it and save it to your account
2. **List Bills**: View your recent bills with details
3. **Summary**: Get spending summaries for different time periods (this month, last month, last week, all time)

## Setup Instructions

### 1. Create a Telegram Bot

1. Open Telegram and search for [@BotFather](https://t.me/botfather)
2. Send `/newbot` command
3. Follow the instructions to choose a name and username for your bot
4. BotFather will provide you with a bot token (looks like `123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ`)
5. Save this token - you'll need it for the next step

### 2. Configure Environment Variables

Add the following to your `.env` file:

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
GROK_API_KEY=your_grok_api_key_here
```

The Grok API key is required for:
- Parsing bill images (extracting text and amounts from photos)
- Intent detection (understanding user queries)

### 3. Set Up Webhook

Once your server is running, you need to tell Telegram where to send updates.

**Option A: Using ngrok for local development**

```bash
# Start ngrok (in a separate terminal)
ngrok http 8080

# Copy the HTTPS URL (e.g., https://abc123.ngrok.io)

# Set the webhook using curl
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook?url=<YOUR_NGROK_URL>/telegram/webhook"

# Example:
curl -X POST "https://api.telegram.org/bot123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ/setWebhook?url=https://abc123.ngrok.io/telegram/webhook"
```

**Option B: Using your production server**

```bash
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook?url=<YOUR_SERVER_URL>/telegram/webhook"

# Example:
curl -X POST "https://api.telegram.org/bot123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ/setWebhook?url=https://api.mibolsillo.com/telegram/webhook"
```

**Verify webhook setup:**

```bash
curl "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getWebhookInfo"
```

### 4. Start the Server

```bash
go run cmd/main.go
```

### 5. Test the Bot

1. Open Telegram and search for your bot using its username
2. Send `/start` or any message to begin
3. The bot will welcome you and explain what it can do

## Usage Examples

### Upload a Bill

Simply send a photo of your bill or receipt to the bot. The bot will:
1. Process the image
2. Extract merchant name, items, amounts, and date
3. Save the bill to your account
4. Send you a confirmation with the parsed details

### List Your Bills

Send any of these messages:
- "Show my bills"
- "List my expenses"
- "What are my recent bills?"
- "Show me my last 10 bills"

### Get a Summary

Send any of these messages:
- "How much did I spend last month?"
- "Summary of this month"
- "Total expenses this month"
- "What did I spend last week?"

## User Mapping

**IMPORTANT**: The current implementation uses a simple in-memory mapping of Telegram user IDs to application user IDs. For production use, you should:

1. Implement a proper user registration system
2. Store user mappings in the database
3. Add authentication/authorization for Telegram users

Current behavior:
- First-time users are automatically assigned a user ID in format `tg_<telegram_user_id>`
- This mapping is stored in memory and will be lost on server restart
- Each Telegram user has their own isolated bill data

## API Endpoints

The Telegram webhook is exposed at:
```
POST /telegram/webhook
```

This endpoint:
- Is publicly accessible (no Clerk authentication required)
- Accepts Telegram update objects
- Processes updates asynchronously
- Returns immediately to satisfy Telegram's timeout requirements

## Architecture

```
Telegram → Webhook → Handler → Intent Detection → Services → Database
                              ↓
                        Grok API (for parsing & intent)
```

Components:
- **Telegram Client** (`internal/adapters/outbound/telegram/`): Handles communication with Telegram API
- **Intent Detection Service** (`internal/core/services/intent_detection_service.go`): Uses Grok to understand user queries
- **Telegram Webhook Handler** (`internal/adapters/inbound/handlers/telegram_webhook_handler.go`): Processes incoming messages
- **Grok Client** (`internal/adapters/outbound/grok/`): Parses bill images and detects intent

## Troubleshooting

### Bot doesn't respond

1. Check webhook is set correctly:
   ```bash
   curl "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getWebhookInfo"
   ```

2. Check server logs for errors

3. Verify server is accessible from the internet (Telegram needs to reach your webhook)

### Image parsing fails

1. Ensure image is clear and well-lit
2. Check that GROK_API_KEY is set correctly
3. Verify you have sufficient Grok API credits

### Intent detection not working

1. Check server logs to see detected intent
2. Verify GROK_API_KEY is valid
3. Try more explicit queries (e.g., "list my bills" instead of "bills")

## Security Considerations

For production deployment:

1. **Verify webhook requests**: Add Telegram webhook secret verification
2. **Rate limiting**: Implement rate limiting on the webhook endpoint
3. **User authentication**: Implement proper user registration and linking
4. **Input validation**: Add validation for all user inputs
5. **Database transactions**: Ensure atomic operations for bill creation

## Example Interaction

```
User: Hi
Bot: Welcome to Mi Bolsillo! 👋
     I can help you manage your bills and expenses...

User: [sends photo of receipt]
Bot: 📸 Processing your bill image...
     ✅ Bill saved successfully!
     🏪 Merchant: Walmart
     💰 Total: USD 45.67
     📅 Date: 2025-10-15
     📝 Items: 8

User: Show my bills
Bot: 📋 Your Recent Bills (showing 5 of 12)
     1. Walmart
        💰 USD 45.67 (PEN 171.38 / USD 45.67)
        📅 2025-10-15
        📝 8 items
     ...

User: How much did I spend this month?
Bot: 📊 Expense Summary - This Month
     💰 Total Spent
        PEN 1,234.56
        USD 329.88
     📋 Number of Bills: 15
     By Category (PEN):
        • Food: 567.89
        • Transportation: 234.56
        • Entertainment: 432.11
```
