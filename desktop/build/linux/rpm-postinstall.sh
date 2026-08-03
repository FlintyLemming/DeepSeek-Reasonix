#!/bin/sh
# RPM %post scriptlet: refresh the freedesktop caches the package feeds.
# Runs as root during install/upgrade; every step no-ops safely when the
# tool or directory is absent so the transaction never fails on them.
if [ -d /usr/share/icons/hicolor ]; then
	touch --no-create /usr/share/icons/hicolor >/dev/null 2>&1 || true
fi
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
	gtk-update-icon-cache -qf /usr/share/icons/hicolor >/dev/null 2>&1 || true
fi
if command -v update-desktop-database >/dev/null 2>&1; then
	update-desktop-database -q /usr/share/applications >/dev/null 2>&1 || true
fi
exit 0
