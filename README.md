# Discorder

This is a CLI tool to list relationships and DM channels, create or remove DM channels, view your guilds and their channels, and export messages from channels you can access with your user token.

This tool is intended for personal use and does not support bot accounts. It uses the Discord API directly, so it may break if Discord changes their API.

## Features

- List all relationships (friends, blocked users, etc.)
- List all direct message channels (dms)
- Create a DM channel with a user
- Remove a DM channel (either a user or a group DM)
- List all guilds the user belongs to
- List channels in a guild
- Get all messages from a channel (pipe to a file or pager)

## Build

```bash
go mod download
go build -o discorder ./cmd/discorder
```

Or install:

```bash
go install ./cmd/discorder
```

## Auth token

Provide the Discord token in one of these ways:

1. First argument when running the program
2. Environment variable `DISCORD_TOKEN`
3. `.env` file containing `DISCORD_TOKEN=...` in the working directory

## Usage

```bash
Usage: ./discorder <token> <action> [args...]
   or: DISCORD_TOKEN=your_token ./discorder <action> [args...]
Available actions: relationships, dms, create-dm, remove-dm, guilds, guild-channels, messages
```

## Examples

```bash
# List all relationships (friends, blocked users, etc.)
./discorder your_token relationships
# List all direct message channels
./discorder your_token dms
# Create a DM channel with a user
./discorder your_token create-dm <user_id>
# Remove a DM channel (either a user or a group DM)
./discorder your_token remove-dm <channel_id>
# List guilds you belong to
./discorder your_token guilds
# List channels in a guild
./discorder your_token guild-channels <guild_id>
# Get all messages from a channel (recommended to pipe out to a file)
./discorder your_token messages <channel_id>
```
