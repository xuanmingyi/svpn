package main

import (
	"os"
	"syscall"
	"unsafe"
)

func MakeTunDevice(name string) (file *os.File, err error) {
	file, err = os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	ifr := make([]byte, 18)

	copy(ifr, []byte(name))

	ifr[16] = 0x01
	ifr[17] = 0x10

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(file.Fd()), uintptr(0x400454ca),
		uintptr(unsafe.Pointer(&ifr[0])))
	if errno != 0 {
		panic(errno)
	}

	return file, nil
}
