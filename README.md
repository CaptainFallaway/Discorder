# Discorder

A small CLI tool to interact with Discord's API for personal account data: list relationships and DMs, create/delete DM channels, and dump messages.

This repository has been refactored to idiomatic Go layout:

- `cmd/discorder/` – main CLI entrypoint
- `internal/discord/` – Discord client and data types
- `internal/cli/` – CLI helpers: formatting, sorting, printing, message utilities

## Features

- List all relationships (friends, blocked users, etc.)
- List all direct message channels
- Create a DM channel with a user
- Remove a DM channel (either a user or a group DM)
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
