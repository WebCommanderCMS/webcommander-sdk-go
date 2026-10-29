//go:build windows

package client

import "syscall"

// consoleIsUnicode reports whether the console output code page is 65001; the legacy OEM and ANSI pages have no U+2500.
func consoleIsUnicode() bool {
	const cpUTF8 = 65001
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleOutputCP")
	if err := proc.Find(); err != nil {
		return false
	}
	codePage, _, _ := proc.Call()
	return uint32(codePage) == cpUTF8 //#nosec G115 -- a Windows code page identifier is bounded by the API
}
