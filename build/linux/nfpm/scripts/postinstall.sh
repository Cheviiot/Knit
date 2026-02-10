#!/bin/sh

# Update icon cache for proper icon display
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
  echo "Updating icon cache..."
  gtk-update-icon-cache -f -t /usr/share/icons/hicolor 2>/dev/null || true
fi

# Update desktop database for .desktop file changes
# This makes the application appear in application menus and registers its capabilities.
if command -v update-desktop-database >/dev/null 2>&1; then
  echo "Updating desktop database..."
  update-desktop-database -q /usr/share/applications 2>/dev/null || true
fi

# Update MIME database for custom URL schemes (x-scheme-handler)
# This ensures the system knows how to handle magnet links.
if command -v update-mime-database >/dev/null 2>&1; then
  echo "Updating MIME database..."
  update-mime-database /usr/share/mime 2>/dev/null || true
fi

exit 0
