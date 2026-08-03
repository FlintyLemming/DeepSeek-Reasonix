//go:build darwin

package main

import "time"

// macOS never runs the tray (see tray_supported_darwin.go); background close
// hides the whole application and restores from the Dock, so this hook is
// unreachable and unconditionally satisfied.
func trayIconRegistered(time.Duration) bool { return true }
