#!/bin/sh
# RPM %post scriptlet: refresh the freedesktop caches the package feeds.
# Runs as root during install/upgrade; every step no-ops safely when the
# tool or directory is absent so the transaction never fails on them.
#
# Touch the .desktop file so GNOME Shell / Gio.AppInfo notice the Exec=
# change after upgrades from pre-v1.20 packages that used reasonix-guard.
if [ -f /usr/share/applications/reasonix.desktop ]; then
	touch /usr/share/applications/reasonix.desktop >/dev/null 2>&1 || true
fi
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
