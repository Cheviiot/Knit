#! /usr/bin/env bash
# Fix WebKit cursor issue on Linux
export WEBKIT_DISABLE_COMPOSITING_MODE=1
# Use system cursor theme
export XCURSOR_PATH="${XCURSOR_PATH:-/usr/share/icons}"
