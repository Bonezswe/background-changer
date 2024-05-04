package main

import (
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

const (
	spiGetDesktopWallpaper = 0x0073
	spiSetDesktopWallpaper = 0x0014
	uiParam                = 0x0000
	spifUpdateINIFile      = 0x01
	spifSendChange         = 0x02
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

func Get() (string, error) {
	var filename [256]uint16
	systemParametersInfo.Call(
		uintptr(spiGetDesktopWallpaper),
		uintptr(cap(filename)),
		uintptr(unsafe.Pointer(&filename[0])),
		uintptr(0),
	)

	return strings.Trim(string(utf16.Decode(filename[:])), "\x00"), nil
}

func SetFromFile(filename string) error {
	filenameUTF16, err := syscall.UTF16PtrFromString(filename)

	if err != nil {
		return err
	}

	systemParametersInfo.Call(
		uintptr(spiSetDesktopWallpaper),
		uintptr(uiParam),
		uintptr(unsafe.Pointer(filenameUTF16)),
		uintptr(spifUpdateINIFile|spifSendChange),
	)

	return nil
}

func SetMode(mode Mode) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, "Control Panel\\Desktop", registry.SET_VALUE)

	if err != nil {
		return err
	}

	defer key.Close()

	var tile string

	if mode == Tile {
		tile = "1"
	} else {
		tile = "0"
	}

	err = key.SetStringValue("TileWallpaer", tile)

	if err != nil {
		return err
	}

	var style string

	switch mode {
	case Center, Tile:
		style = "0"
	case Fit:
		style = "6"
	case Stretch:
		style = "2"
	case Span:
		style = "22"
	case Crop:
		style = "10"
	default:
		panic("Invalid Wallpaper mode")
	}

	err = key.SetStringValue("WallpaperStyle", style)

	if err != nil {
		return err
	}

	path, err := Get()

	if err != nil {
		return err
	}

	return SetFromFile(path)
}

func GetScreenResolution() string {
	getSystemMetrics := user32.NewProc("GetSystemMetrics")

	const (
		SM_CXSCREEN = 0
		SM_CYSCREEN = 1
	)

	screenWidth, _, _ := getSystemMetrics.Call(uintptr(SM_CXSCREEN))
	screenHeight, _, _ := getSystemMetrics.Call(uintptr(SM_CYSCREEN))

	width := int(screenWidth)
	height := int(screenHeight)

	res := GetResolutionName(width, height)

	return res
}
