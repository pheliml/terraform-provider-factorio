#!/bin/bash
set -euo pipefail

CLIENT_MOD_DIR=$1

# Copy the mod to your client mods folder
cp -r ../../mod/terraform-crud-api "$HOME/$CLIENT_MOD_DIR"
# Create a folder to store the Factorio server data
mkdir factorio-volume
# Copy the factorio mod to the mods directory
mkdir factorio-volume/mods
cp -r ../../mod/terraform-crud-api factorio-volume/mods
# Configure the rcon pw
mkdir factorio-volume/config
echo "SOMEPASSWORD" > factorio-volume/config/rconpw

chmod +x ./run.sh
./run.sh