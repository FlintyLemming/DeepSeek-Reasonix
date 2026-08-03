//go:build !windows

package main

import _ "embed"

// Linux tray icon: a dedicated 256x256 8-bit PNG. The StatusNotifier protocol
// ships the icon as a raw ARGB32 pixmap over D-Bus, so the 1024x1024 app icon
// would cost a ~4MiB property per registration and scale down to the same
// ~22px tray slot anyway; some hosts reject or fail to render such oversized
// pixmaps, leaving the backgrounded app with no way back. The tray never runs
// on macOS (see tray_supported_darwin.go), so this asset is Linux-only in
// practice. Regenerate from build/appicon.png only via area averaging — the
// asset test pins the size, bit depth, and Reasonix blue background.
//
//go:embed build/linux/trayicon.png
var trayIconBytes []byte
