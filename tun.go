package main

import (
	"log"
	"syscall"
	"unsafe"
)

const (
	tunDevice    = "/dev/net/tun"
	IFF_TUN      = 0x0001
	IFF_TAP      = 020002
	TUNSETNOCSUM = 0x400454c8
	TUNSETIFF    = 0x400454ca
)

/*
type sockaddr struct {
    family uint16
    port uint16
    addr [20]byte
}
*/

type ifreq struct {
	ifname [16]byte
	flags  int16
}

func tunOpen() (fd int) {
	fd, err := syscall.Open(tunDevice, syscall.O_RDWR, 0)
	if err != nil {
            log.Fatal("open tun device error")
            return -1
	}

	log.Printf("open tun success: fd=%d\n", fd)
	var ifr ifreq
	ifr.flags = IFF_TUN
	cmd := TUNSETIFF
        if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(cmd), uintptr(unsafe.Pointer(&ifr))); err != 0 {
            log.Fatal("ioctl error: ", err)
            return -1
	}

        log.Println("ioctl TUNSETIFF success!")
        return fd
}

func tunClose(fd int) {
    syscall.Close(fd)
}

func main() {
    tunfd := tunOpen()
    defer tunClose(tunfd)

    //ifconf
}
