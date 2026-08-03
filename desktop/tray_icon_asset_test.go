package main

import (
	"image"
	"image/png"
	"os"
	"testing"
)

// TestLinuxTrayIconIsSizedEightBit pins the StatusNotifier pixmap source: a
// 256x256 8-bit RGBA PNG. The icon travels over D-Bus as raw ARGB32, so a
// 16-bit or oversized source balloons the payload and some hosts refuse to
// render it — the exact failure that strands a backgrounded app with no tray
// icon.
func TestLinuxTrayIconIsSizedEightBit(t *testing.T) {
	f, err := os.Open("build/linux/trayicon.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := img.Bounds().Dx(), 256; got != want {
		t.Fatalf("tray icon width = %d, want %d", got, want)
	}
	if got, want := img.Bounds().Dy(), 256; got != want {
		t.Fatalf("tray icon height = %d, want %d", got, want)
	}
	// Go's PNG decoder maps 8-bit RGBA to NRGBA/RGBA and 16-bit to the *64
	// types, so the concrete type witnesses the bit depth.
	switch img.(type) {
	case *image.NRGBA, *image.RGBA, *image.Paletted, *image.Gray:
	default:
		t.Fatalf("tray icon must be an 8-bit PNG, got %T", img)
	}

	assertFullCanvasRoundedIcon(t, img, 256)
}
