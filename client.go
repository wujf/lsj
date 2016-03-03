package main

import(
    "log"
    "os"
    "net"
    "time"
    "github.com/pkg/taptun"
)

var serverAddr *net.UDPAddr

func down(tun *taptun.Tun, conn *net.UDPConn) {
    buf := make([]byte, 1400)
    for {
        size, err := conn.Read(buf)
        if err != nil {
            log.Println("read UDP: ", err)
            break
        }
        log.Println("read %d bytes from UDP\n", size)
        n, err := tun.Write(buf)
        log.Printf("send %d bytes to TUN\n", n)
    }
}

func up(tun *taptun.Tun, conn *net.UDPConn) {
    buf := make([]byte, 1500)
    for {
        size, err := tun.Read(buf)
        if err != nil {
            log.Println("read error: ", err)
            break
        }

        log.Printf("read %d bytes from TUN\n", size)
        n, err := conn.Write(buf, size)
        if err != nil {
            log.Fatal("wirte UDP server error:", err)
            break
        }

        log.Printf("sent %d bytes to UDP server\n", n)
    }
}

func main() {
    tun, err := taptun.OpenTun()
    if err != nil {
        log.Fatal("Open TUN device error: ", err)
        return
    }
    defer tun.Close()
    log.Println("Opened TUN device: ", tun.String())

    serverAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
    if err != nil {
        log.Fatal("Can't resolv server addr:", os.Args[1])
        return
    }

    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        log.Fatal("Dial UDP error: ", err, ", server:", serverAddr)
        return
    }
    defer conn.Close()

    go up(tun, conn)
    go down(tun, conn)

    for {
        time.Sleep(time.Second * 3)
        log.Println("Client running ...")
    }
}
