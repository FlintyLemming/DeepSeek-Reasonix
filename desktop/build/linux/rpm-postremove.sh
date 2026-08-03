#!/bin/sh
# RPM %postun scriptlet: refresh the freedesktop caches after removal or
# upgrade. Runs as root; $1 is the number of package instances left after
# the transaction, but the cache refresh is safe in every case.
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
