package dmsofts

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

const PROCESS_ALL_ACCESS = 0x1F0FFF

var (
	kernel32               = syscall.MustLoadDLL("kernel32.dll")
	procOpenProcess        = kernel32.MustFindProc("OpenProcess")
	procReadProcessMemory  = kernel32.MustFindProc("ReadProcessMemory")
	procWriteProcessMemory = kernel32.MustFindProc("WriteProcessMemory")
	VirtualProtectEx       = kernel32.MustFindProc("VirtualProtectEx")
)

//	获取根据进程名获取进程pid
//	@param openName：进程名
func ProcessName(openName string) (pid int32) {
	pids, _ := process.Pids()
	for _, pid := range pids {
		pn, _ := process.NewProcess(pid)
		pName, _ := pn.Name()
		if openName == pName {
			return pn.Pid
		}
	}
	return -1
}

//	打开进程 返回进程句柄
//	@param pid：进程pid
func OpenProcess(pid int32) uintptr {
	handle, _, _ := procOpenProcess.Call(uintptr(PROCESS_ALL_ACCESS), uintptr(1), uintptr(pid))
	return handle
}

//	读内存
//	handle 进程句柄
//	adress 内存地址  因为可能涉及到64位的程序 所有这里用 uint64
func ReadProcessMemoryInt64(handle uintptr, adress uint64) (Data uint64) {

	var (
		data   uint64
		length uint32
	)
	ret, _, e := procReadProcessMemory.Call(uintptr(handle),
		uintptr(adress),
		uintptr(unsafe.Pointer(&data)),
		4, uintptr(unsafe.Pointer(&length))) // read 2 bytes
	if ret == 0 {
		fmt.Println("  Error:", e)
	}
	return data
}

//	读内存
//	handle 进程句柄
//	adress 内存地址  因为可能涉及到64位的程序 所有这里用 uint64
func ReadProcessMemoryFloat32(handle uintptr, adress uint64) (Data float32) {

	var (
		data   float32
		length uint32
	)
	ret, _, e := procReadProcessMemory.Call(uintptr(handle),
		uintptr(adress),
		uintptr(unsafe.Pointer(&data)),
		4, uintptr(unsafe.Pointer(&length))) // read 2 bytes
	if ret == 0 {
		fmt.Println("  Error:", e)
	}
	return data
}

//	写内存小数型
//	@param handle: 进程句柄
//	@param Address: 内存地址  因为可能涉及到64位的程序 所有这里用 uint64
//	@param lpBuffer: 改写数值
//	@param nSize: 字节数
func WriteProcessMemoryFloat32(handle uintptr, address uint64, lpBuffer float32, nSize uintptr) (int, bool) {
	var nBytesWritten int

	ret, _, _ := procWriteProcessMemory.Call(
		handle,
		uintptr(address),
		uintptr(unsafe.Pointer(&lpBuffer)),
		nSize,
		uintptr(unsafe.Pointer(&nBytesWritten)),
	)
	return nBytesWritten, ret != 0
}

//	写内存整数型
//	@param handle: 进程句柄
//	@param Address: 内存地址  因为可能涉及到64位的程序 所有这里用 uint64
//	@param lpBuffer: 改写数值
//	@param nSize: 字节数
func WriteProcessMemoryInt64(handle uintptr, address uint64, lpBuffer int64, nSize uintptr) (int, bool) {
	var nBytesWritten int

	ret, _, _ := procWriteProcessMemory.Call(
		handle,
		uintptr(address),
		uintptr(unsafe.Pointer(&lpBuffer)),
		nSize,
		uintptr(unsafe.Pointer(&nBytesWritten)),
	)
	return nBytesWritten, ret != 0
}

//  改变内存保护属性 改变为可读可写
//	@param handle: 进程句柄
//  @param Address: 内存地址  因为可能涉及到64位的程序 所有这里用 uint64
//  说明： 调用这个方法系统会返回64位的指针 所有说32位系统是无法使用该API的 不然就会程序卡死
func VirtualProtectExReadWrite(handle uintptr, address uint64, name uintptr, a int) bool {

	//handle uintptr, 进程句柄
	//address uint64, 内存地址
	//dwSize uint, 改变区域大小
	//flNewProtect uint, 改变保护属性
	//lpflOldProtect uint, 改变后返回的指针

	ret, _, _ := VirtualProtectEx.Call(
		handle,
		uintptr(address),
		1000,
		name,
		uintptr(unsafe.Pointer(&a)),
	)
	return ret != 0
}

//  改变内存保护属性 改变为可读
//	@param handle: 进程句柄
//  @param Address: 内存地址  因为可能涉及到64位的程序 所有这里用 uint64
//  说明： 调用这个方法系统会返回64位的指针 所有说32位系统是无法使用该API的 不然就会程序卡死
func VirtualProtectExRead(handle uintptr, address uint64) (uintptr, bool) {
	var a int
	//handle uintptr, 进程句柄
	//address uint64, 内存地址
	//dwSize uint, 改变区域大小
	//flNewProtect uint, 改变保护属性
	//lpflOldProtect uint, 改变后返回的指针
	ret, _, _ := VirtualProtectEx.Call(
		handle,
		uintptr(address),
		4,
		windows.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&a)),
	)
	return uintptr(unsafe.Pointer(&a)), ret != 0
}
