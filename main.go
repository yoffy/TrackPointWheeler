package main

import (
	//"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookExW")
	procGetMessage          = user32.NewProc("GetMessageW")
)

const (
	WH_MOUSE_LL    = 14
	WM_MOUSEMOVE   = 0x200
	WM_LBUTTONDOWN = 0x201
	WM_LBUTTONUP   = 0x202
	WM_RBUTTONDOWN = 0x204
	WM_RBUTTONUP   = 0x205
	WM_MBUTTONDOWN = 0x207
	WM_MBUTTONUP   = 0x208
	WM_MOUSEWHEEL  = 0x20A
	WM_MOUSEHWHEEL = 0x20E
	NULL           = 0
)

type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

type MSLLHOOKSTRUCT struct {
	Pt          POINT
	MouseData   DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo uintptr
}

type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

func HookMouse() {
	var hook HHOOK

	hook = SetWindowsHookEx(
		WH_MOUSE_LL,
		(HOOKPROC)(func(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
			if nCode == 0 && wparam == WM_MBUTTONDOWN {
				// never call CallNextHookEx because to ignore middle button
				return 1
			}
			return CallNextHookEx(hook, nCode, wparam, lparam)
		}), 0, 0)

	var msg MSG
	for GetMessage(&msg, 0, 0, 0) != 0 {
	}

	UnhookWindowsHookEx(hook)
	hook = 0
}

func main() {
	HookMouse()
}
