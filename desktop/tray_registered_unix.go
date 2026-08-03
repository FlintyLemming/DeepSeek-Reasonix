//go:build !windows && !darwin

package main

import (
	"errors"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

const trayRegistrationPollInterval = 50 * time.Millisecond

// trayIconRegistered confirms the StatusNotifierWatcher actually accepted this
// process's tray item before a window close may hide to the tray. systray's
// ready callback fires ahead of registration on Linux (and regardless of it),
// so the trayReady flag alone never proves the icon is visible; hiding the
// window without a registered icon strands the app with no way to reopen or
// quit it. When registration cannot be confirmed, closing quits.
func trayIconRegistered(timeout time.Duration) bool {
	return awaitTrayRegistration(trayItemRegistered, timeout, trayRegistrationPollInterval)
}

// awaitTrayRegistration polls check until it succeeds or timeout elapses.
// Registration completes asynchronously right after the tray starts, so a
// close racing the tray startup gets a short grace period instead of a false
// negative.
func awaitTrayRegistration(check func() bool, timeout, interval time.Duration) bool {
	if check() {
		return true
	}
	if timeout <= 0 || interval <= 0 {
		return false
	}
	deadline := time.Now().Add(timeout)
	for {
		time.Sleep(interval)
		if check() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
	}
}

func trayItemRegistered() bool {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false
	}
	// SessionBus returns a shared singleton connection. Do not close it here:
	// other desktop integrations in this process may still be using it.
	obj := conn.Object("org.freedesktop.DBus", "/org/freedesktop/DBus")
	var owner string
	if err := obj.Call("org.freedesktop.DBus.GetNameOwner", 0, "org.kde.StatusNotifierWatcher").Store(&owner); err != nil || owner == "" {
		return false
	}
	watcher := conn.Object("org.kde.StatusNotifierWatcher", "/StatusNotifierWatcher")
	variant, err := watcher.GetProperty("org.kde.StatusNotifierWatcher.RegisteredStatusNotifierItems")
	if err != nil {
		// A watcher that vanished mid-check took the icon with it: closing must
		// quit. Any other read failure (e.g. a watcher that does not expose the
		// registered-items list) keeps the historical hide-to-tray behavior.
		return !dbusNameOwnerLost(err)
	}
	items, ok := variant.Value().([]string)
	if !ok {
		return true
	}
	// Watchers record the registering connection either as its unique name or
	// as the well-known StatusNotifierItem name systray requests. The shared
	// session connection owns both, so any owned name in the list confirms us.
	return trayRegistrationListed(items, conn.Names()...)
}

func dbusNameOwnerLost(err error) bool {
	var dbusErr dbus.Error
	if errors.As(err, &dbusErr) {
		switch dbusErr.Name {
		case "org.freedesktop.DBus.Error.NameHasNoOwner",
			"org.freedesktop.DBus.Error.ServiceUnknown":
			return true
		}
	}
	return false
}

// trayRegistrationListed reports whether any registered item belongs to one of
// the given bus names. Watchers record items either as the bare bus name
// ("org.kde.StatusNotifierItem-1234-1" on Plasma) or as "busName@objectPath"
// (the ayatana/GNOME-style watcher), so both shapes are matched.
func trayRegistrationListed(items []string, names ...string) bool {
	for _, item := range items {
		item = strings.TrimSpace(item)
		for _, name := range names {
			if name == "" {
				continue
			}
			if item == name || strings.HasPrefix(item, name+"@") {
				return true
			}
		}
	}
	return false
}
