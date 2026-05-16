#!/bin/bash
set -euo pipefail

# Enable user to modify server mount
sudo chown -R $USER:$USER factorio-volume/

# --- CONFIGURATION ---
CLIENT_MOD_DIR="$HOME/.factorio/mods"
SRC_MOD_DIR="../../../mod/terraform-crud-api"
VOLUME_DIR="factorio-volume"

echo "Reading mod metadata..."
if [ ! -f "$SRC_MOD_DIR/info.json" ]; then
    echo "Error: Cannot find $SRC_MOD_DIR/info.json"
    exit 1
fi

# Automatically extract name and version from info.json
MOD_NAME=$(grep -o '"name": *"[^"]*"' "$SRC_MOD_DIR/info.json" | grep -o '"[^"]*"$' | tr -d '"')
MOD_VERSION=$(grep -o '"version": *"[^"]*"' "$SRC_MOD_DIR/info.json" | grep -o '"[^"]*"$' | tr -d '"')

# Define the zip filename and the expected internal directory structure
ZIP_NAME="${MOD_NAME}_${MOD_VERSION}.zip"
INTERNAL_DIR="${MOD_NAME}_${MOD_VERSION}"

echo "Mod detected: $MOD_NAME (v$MOD_VERSION)"

# Create a temporary staging area to build the zip archive cleanly
TMP_DIR=$(mktemp -d)
mkdir -p "$TMP_DIR/$INTERNAL_DIR"
cp -r "$SRC_MOD_DIR/." "$TMP_DIR/$INTERNAL_DIR/"

echo "Packaging mod into zip archive..."
# Compress from within the temp directory to maintain the correct root path structure
(cd "$TMP_DIR" && zip -rq "$ZIP_NAME" "$INTERNAL_DIR")

# --- CLIENT DEPLOYMENT ---
echo "Deploying to Factorio Client..."
mkdir -p "$CLIENT_MOD_DIR"

# Remove any old directories or zip files matching this mod to prevent conflicts
rm -rf "$CLIENT_MOD_DIR/${MOD_NAME}"_*
cp "$TMP_DIR/$ZIP_NAME" "$CLIENT_MOD_DIR/"

# --- SERVER DEPLOYMENT ---
echo "Preparing Server Volume..."
mkdir -p "$VOLUME_DIR/mods"
mkdir -p "$VOLUME_DIR/config"

# Clear old server assets and copy the freshly generated archive
sudo rm -rf "$VOLUME_DIR/mods/${MOD_NAME}"_*
sudo cp "$TMP_DIR/$ZIP_NAME" "$VOLUME_DIR/mods/"

# Configure RCON password
sudo echo "SOMEPASSWORD" > "$VOLUME_DIR/config/rconpw"

# --- SERVER SETTINGS DEPLOYMENT (OFFLINE / IP MODE) ---
echo "Deploying server-settings.json to volume..."
cat << 'EOF' > "$VOLUME_DIR/config/server-settings.json"
{
  "name": "Terraform Dev Server",
  "description": "Private local development server via direct IP connection.",
  "tags": ["dev", "local"],
  "max_players": 0,

  "visibility": {
    "public": false,
    "lan": true
  },

  "username": "",
  "password": "",
  "token": "",
  "game_password": "",

  "require_user_verification": false,

  "max_upload_in_kilobytes_per_second": 0,
  "max_upload_slots": 5,
  "minimum_latency_in_ticks": 0,
  "max_heartbeats_per_second": 60,

  "ignore_player_limit_for_returning_players": false,
  "allow_commands": "admins-only",
  "autosave_interval": 10,
  "autosave_slots": 5,
  "afk_autokick_interval": 0,
  "auto_pause": true,
  "only_admins_can_pause_the_game": true,
  "autosave_only_on_server": true,
  "non_blocking_saving": false
}
EOF

# --- CLEANUP ---
rm -rf "$TMP_DIR"

echo "Starting run.sh..."
chmod +x ./run.sh
./run.sh