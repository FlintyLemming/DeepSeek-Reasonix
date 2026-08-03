//go:build !windows && !darwin

package main

import (
	"testing"
	"time"
)

func TestAwaitTrayRegistrationAcceptsAnAlreadyRegisteredTray(t *testing.T) {
	calls := 0
	if !awaitTrayRegistration(func() bool { calls++; return true }, 0, 0) {
		t.Fatal("an already registered tray must be accepted immediately")
	}
	if calls != 1 {
		t.Fatalf("check ran %d times, want a single immediate check", calls)
	}
}

func TestAwaitTrayRegistrationPollsUntilRegistrationCompletes(t *testing.T) {
	calls := 0
	ok := awaitTrayRegistration(func() bool {
		calls++
		return calls >= 3
	}, time.Second, time.Millisecond)
	if !ok {
		t.Fatal("registration that completes during the grace period must be accepted")
	}
	if calls < 3 {
		t.Fatalf("check ran %d times, want polling until success", calls)
	}
}

func TestAwaitTrayRegistrationTimesOutWhenTheIconNeverRegisters(t *testing.T) {
	start := time.Now()
	if awaitTrayRegistration(func() bool { return false }, 30*time.Millisecond, 5*time.Millisecond) {
		t.Fatal("a tray that never registers must not be reported as registered")
	}
	if elapsed := time.Since(start); elapsed < 30*time.Millisecond {
		t.Fatalf("gave up after %v, want the full grace period", elapsed)
	}
}

func TestAwaitTrayRegistrationRejectsImmediatelyWithoutGracePeriod(t *testing.T) {
	calls := 0
	if awaitTrayRegistration(func() bool { calls++; return false }, 0, time.Millisecond) {
		t.Fatal("zero timeout must not accept an unregistered tray")
	}
	if calls != 1 {
		t.Fatalf("check ran %d times, want a single immediate check", calls)
	}
}

func TestTrayRegistrationListedMatchesUniqueAndWellKnownNames(t *testing.T) {
	items := []string{":1.42", "org.kde.StatusNotifierItem-999-1", "some-other-app"}
	if !trayRegistrationListed(items, ":1.42", "org.kde.StatusNotifierItem-1234-1") {
		t.Fatal("unique connection name must match")
	}
	if !trayRegistrationListed(items, ":1.99", "org.kde.StatusNotifierItem-999-1") {
		t.Fatal("well-known StatusNotifierItem name must match")
	}
	if trayRegistrationListed(items, ":1.98", "org.kde.StatusNotifierItem-1234-1") {
		t.Fatal("foreign items must not count as our registration")
	}
	if trayRegistrationListed(nil, ":1.42") {
		t.Fatal("an empty registration list must not match")
	}
	if trayRegistrationListed([]string{"", " "}, "", "") {
		t.Fatal("blank entries and names must never match")
	}
	if !trayRegistrationListed([]string{"  :1.42  "}, ":1.42") {
		t.Fatal("surrounding whitespace must be ignored")
	}
}

func TestTrayRegistrationListedMatchesNameAtPathEntries(t *testing.T) {
	// ayatana/GNOME-style watchers record "uniqueName@objectPath".
	items := []string{":1.715@/StatusNotifierItem", ":1.50@/org/ayatana/NotificationItem/app"}
	if !trayRegistrationListed(items, ":1.715", "org.kde.StatusNotifierItem-715-1") {
		t.Fatal("unique name must match its name@path entry")
	}
	if !trayRegistrationListed(items, ":1.50") {
		t.Fatal("unique name must match a name@path entry with a nested path")
	}
	if trayRegistrationListed(items, ":1.71") {
		t.Fatal("a shorter unique name must not prefix-match a longer one")
	}
	if trayRegistrationListed(items, ":1.716") {
		t.Fatal("an unregistered unique name must not match")
	}
}
