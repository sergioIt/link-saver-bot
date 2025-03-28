# Link Saver Bot

A Telegram bot that helps users save and manage their links.

## Setup

1. Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

2. Get a Telegram Bot Token:
   - Message [@BotFather](https://t.me/botfather) on Telegram
   - Create a new bot using the `/newbot` command
   - Copy the token provided by BotFather

3. Run the bot using one of these methods:

   Using exec form with environment variable (recommended):
   ```bash
   docker run -d --name link-saver-bot -v $(pwd)/user_data:/app/user_data link-saver-bot ./main -telegram-token="your_token_here"
   ```

   Using .env file:
   ```bash
   docker run -d --name link-saver-bot --env-file .env -v $(pwd)/user_data:/app/user_data link-saver-bot
   ```

   Using environment variable:
   ```bash
   docker run -d --name link-saver-bot -e TELEGRAM_TOKEN=your_token_here -v $(pwd)/user_data:/app/user_data link-saver-bot
   ```

   Using Docker Compose:
   ```bash
   TELEGRAM_TOKEN=your_token_here docker compose up -d
   ```
   Note: Docker Compose will use the name 'link-saver-bot' as specified in docker-compose.yml

## Development

Requirements:
- Go 1.22 or later
- Docker (optional)

Running locally:
```bash
go run main.go -telegram-token="your_token_here"
```

⚠️ Security Notes:
- Be careful when using command-line arguments as they might be visible in process listings (ps aux)
- For production, prefer using environment files or Docker secrets

## Container Management

Stop the container:
```bash
docker stop link-saver-bot
```

Remove the container:
```bash
docker rm link-saver-bot
```

View logs:
```bash
docker logs link-saver-bot
```

Follow logs:
```bash
docker logs -f link-saver-bot
```
