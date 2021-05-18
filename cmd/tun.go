package cmd

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func Mktun() {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	ifr := make([]byte, 18)
	copy(ifr, []byte("tun0"))

	ifr[16] = 0x01
	ifr[17] = 0x10


	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), uintptr(0x400454ca), uintptr(unsafe.Pointer(&ifr[0])))
	if errno != 0 {
		panic(errno)
	}
	for {
		buf := make([]byte, 2048)
		read, err := file.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println(buf[:read])
	}
}