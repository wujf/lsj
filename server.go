package main

import(
    "log"
    "net"
    "time"
    "github.com/pkg/taptun"
)

/*
type ServerContext struct {
    conn        *net.UDPConn
    tun         *taptun.Tun
}

func upLink(src int, dst *net.UDPConn) {
    buf := make([1500]byte)
}

func downLink(src *net.UDPConn, dst int) {
    buf := make([1500]byte)
    n, client, err := s.conn.ReadFromUDP(buf)
    log.Println("ReadFromUDP: client:", client, ", ", n, " bytes")
}

func (s *ServerContext) Serv() error {
    s.tun = tunOpen()
    defer tunClose(s.tun)

    go upLink(s.tun, s.conn)
    go downLink(s.conn, s.tun)
}
*/

var clientAddr *net.UDPAddr

func up(tun *taptun.Tun, conn *net.UDPConn) {
    buf := make([]byte, 1500)
    for {
        n, client, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Println("read UDP error: ", err)
            break
        }
        clientAddr = client
        log.Printf("read %d bytes from UDP\n", n)
        if n, err = tun.Write(buf[:n]); err != nil {
            log.Println("send to TUN error:", err)
        }
    }
}

func down(tun *taptun.Tun, conn *net.UDPConn) {
    buf := make([]byte, 1500)
    for {
        size, err := tun.Read(buf)
        if err != nil {
            log.Println("read error: ", err)
            break
        }

        log.Printf("read %d bytes from TUN\n", size)
        if size, err = conn.WriteTo(buf[:size], clientAddr); err!= nil {
            log.Println("send to UDP error:", err)
        }
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

    addr, _ := net.ResolveUDPAddr("udp", ":8883")
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Fatal("Listen UDP error: ", err)
        return
    }
    defer conn.Close()

    go up(tun, conn)
    go down(tun, conn)

    for {
        time.Sleep(time.Second * 3)
        log.Println("Server running ...")
    }
}
