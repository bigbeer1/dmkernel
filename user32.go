package dmkernel

import (
	"github.com/lxn/win"
	"syscall"
)

//根据窗口标题取窗口句柄
// @param 窗口标题

func FindWindow(name string) win.HWND {

	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(name))

	return hwnd
}
