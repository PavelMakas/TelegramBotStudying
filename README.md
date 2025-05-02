# Telegram Bot with Message Reversal and Story Generation

This is a Telegram bot written in Go that provides two main functionalities:

1. **Message Reversal**: The bot replies to any incoming message with its reversed version.
   - Example: "hello" → "olleh"

2. **Story Generation**: Using the `/story` command, the bot generates a short story (max 400 characters) in a specified style using OpenAI's API.
   - Example: `/story sci-fi` will generate a science fiction story

## Setup

1. Clone the repository
2. Create a `.env` file in the root directory with the following variables:
   ```
   TELEGRAM_BOT_TOKEN=your_telegram_bot_token
   OPENAI_API_KEY=your_openai_api_key
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the bot:
   ```bash
   go run main.go
   ```

## Project Structure

- `main.go` - Main application entry point
- `config/` - Configuration handling
- `handlers/` - Bot command and message handlers
- `utils/` - Utility functions including string reversal
- `tests/` - Unit tests

## Testing

Run the tests using:
```bash
go test ./...
``` 