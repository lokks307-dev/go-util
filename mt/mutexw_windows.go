package mt

import (
	"syscall"
	"unsafe"
)

const (
	ERROR_ALREADY_EXISTS = 183
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutex  = kernel32.NewProc("CreateMutexW")
	procReleaseMutex = kernel32.NewProc("ReleaseMutex")
)

func CreateMutex(name string) (uintptr, error) {
	var ret uintptr
	var err error
	var ptrName *uint16

	ptrName, err = syscall.UTF16PtrFromString(name)
	if err != nil {
		return ret, err
	}

	ret, _, err = procCreateMutex.Call(
		0,
		0,
		uintptr(unsafe.Pointer(ptrName)),
	)
	switch int(err.(syscall.Errno)) {
	case 0:
		return ret, nil
	default:
		return ret, err
	}
}

func ReleaseMutex(hmutex uintptr) bool {
	ret, _, _ := procReleaseMutex.Call(hmutex)
	return ret != 0
}

func CheckMutex(mname string) (uintptr, bool) {
	id, err := CreateMutex(mname)

	if err != nil {
		if int(err.(syscall.Errno)) == ERROR_ALREADY_EXISTS {
			return 0, false
		}
	}

	return id, true
}
