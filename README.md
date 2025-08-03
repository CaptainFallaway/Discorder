# Discorder

## Build

```bash
go mod download
go build -o discorder .
```

## Usage

```bash
Usage: ./discorder <token> <action> [args...]
   or: DISCORD_TOKEN=your_token ./program <action> [args...]
Available actions: rels, dms, gdm, rdm, msgs
```

## Examples

```bash
# List all relationships
./discorder your_token rels
# List all direct message channels
./discorder your_token dms
# Create or get a DM channel with a user
./discorder your_token gdm <user_id>
# Remove a DM channel with a user
./discorder your_token rdm <user_id>
# Get messages from a channel
./discorder your_token msgs <channel_id>
```
