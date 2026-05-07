#!/bin/bash
set -euo pipefail

FACTORIO_IMG={$1:"factorio:stable-2.0.76"}
# Run factorio server
docker run -it -p 127.0.0.1:34197:34197/udp -p 127.0.0.1:27015:27015/tcp -v "$(pwd)/factorio-volume:/factorio" factoriotools/$FACTORIO_IMG
