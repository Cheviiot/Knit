#!/bin/bash
# Knit launcher wrapper
# Fix WebKit cursor issue
export WEBKIT_DISABLE_COMPOSITING_MODE=1
export XCURSOR_PATH="${XCURSOR_PATH:-/usr/share/icons}"

exec /usr/lib/knit/knit-bin "$@"
