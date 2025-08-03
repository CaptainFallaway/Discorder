# Discorder

Discorder is a command-line tool designed to retrieve all messages from a user's DM channels. It also includes features for managing direct messages (DMs) and relationships, such as viewing your friends list, creating new / getting DM channels from a user id, and then also deleting existing DMs ones.

The project is very script-like since I made this in a few hours and I could not care less atm.

## Build

```bash
go mod download
go build -o discorder .
```

## Notes

There are three ways to provide the Discord token:

1. As the first argument when running the program.
2. As an environment variable `DISCORD_TOKEN`.
3. As an environment variable `DISCORD_TOKEN` in a `.env` file in the same directory as the program

## Usage

```bash
Usage: ./discorder <token> <action> [args...]
   or: DISCORD_TOKEN=your_token ./discorder <action> [args...]
Available actions: rels, dms, gdm, rdm, msgs
```

## Examples

```bash
# List all relationships (friends, blocked users, etc.)
./discorder your_token rels
# List all direct message channels
./discorder your_token dms
# Create or get a DM channel with a user
./discorder your_token gdm <user_id>
# Remove a DM channel (either a user or a group DM)
./discorder your_token rdm <channel_id>
# Get all messages from a channel (recommended to pipe out to a file)
./discorder your_token msgs <channel_id>
```
