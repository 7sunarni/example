package yellowsocks

import (
	"io"
	"net"

	"github.com/esrrhs/gohome/common"
	"github.com/esrrhs/gohome/loggo"
	"github.com/esrrhs/gohome/network"
	"github.com/esrrhs/pingtunnel/option"
)

func Run() {
	if *option.ListenTCP == "" {
		return
	}
	defer common.CrashLog()
	tcpaddr, err := net.ResolveTCPAddr("tcp", *option.ListenTCP)
	if err != nil {
		loggo.Error("listen fail %s", err)
		return
	}

	tcplistenConn, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		loggo.Error("Error listening for tcp packets: %s", err)
		return
	}
	loggo.Info("listen ok %s", tcpaddr.String())

	dstaddr, err := net.ResolveTCPAddr("tcp", *option.ListenSocks5)
	if err != nil {
		loggo.Error("target fail %s", err)
		return
	}
	loggo.Info("target %s", dstaddr.String())

	for {
		conn, err := tcplistenConn.AcceptTCP()
		if err != nil {
			loggo.Info("Error accept tcp %s", err)
			continue
		}

		go process(conn, dstaddr)
	}
}

func process(conn *net.TCPConn, socks5addr *net.TCPAddr) {

	defer common.CrashLog()

	loggo.Info("start conn from %s", conn.RemoteAddr())

	host, port, err := getOriginalDst(conn)

	loggo.Info("parse conn from %s -> %s:%d", conn.RemoteAddr(), host, port)

	socks5conn, err := net.DialTCP("tcp", nil, socks5addr)
	if err != nil {
		loggo.Info("dial socks5 conn fail %s %v", socks5addr, err)
		return
	}

	loggo.Info("dial socks5 conn ok %s -> %s:%d", conn.RemoteAddr(), host, port)

	err = network.Sock5Handshake(socks5conn, 0, "", "")
	if err != nil {
		loggo.Error("sock5Handshake fail %s", err)
		return
	}

	loggo.Info("Handshake socks5 conn ok %s -> %s:%d", conn.RemoteAddr(), host, port)

	err = network.Sock5SetRequest(socks5conn, host, port, 0)
	if err != nil {
		conn.Close()
		loggo.Error("sock5SetRequest fail %s", err)
		return
	}

	loggo.Info("SetRequest socks5 conn ok %s -> %s:%d", conn.RemoteAddr(), host, port)

	go transfer(conn, socks5conn, conn.RemoteAddr().String(), socks5conn.RemoteAddr().String())
	go transfer(socks5conn, conn, socks5conn.RemoteAddr().String(), conn.RemoteAddr().String())

	loggo.Info("process conn ok %s -> %s:%d", conn.RemoteAddr(), host, port)
}

func transfer(destination io.WriteCloser, source io.ReadCloser, dst string, src string) {
	defer common.CrashLog()
	defer destination.Close()
	defer source.Close()
	// loggo.Info("begin transfer from %s -> %s", src, dst)
	io.Copy(destination, source)
	// loggo.Info("end transfer from %s -> %s", src, dst)
}
