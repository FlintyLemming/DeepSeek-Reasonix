//go:build windows

package main

import "time"

// The Windows tray icon lives in the notification area owned by the systray
// Win32 message loop that startDesktopTray runs; once the tray started, the
// icon is reachable and no extra confirmation is needed.
func trayIconRegistered(time.Duration) bool { return true }
