#!/bin/bash
# Use rofi to find a file and then open it in acme

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
rofi -show 'open in acme' -modi 'open in acme':$CURRENT_DIR/acme-rofi-file-select.sh
