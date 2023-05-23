package option

import "flag"

var (
	Model               *string
	ListenSocks5        *string
	ListenTCP           *string
	ListenHTTP          *string
	Target              *string
	Server              *string
	Timeout             *int
	Key                 *int
	Tcpmode             *int
	TcpmodeBuffersize   *int
	TcpmodeMaxwin       *int
	TcpmodeResendTimems *int
	TcpmodeCompress     *int
	Nolog               *int
	Noprint             *int
	TcpmodeStatus       *int
	Loglevel            *string
	OpenSock5           *int
	Maxconn             *int
	MaxProcessThread    *int
	MaxProcessBuffer    *int
	Profile             *int
	Conntt              *int
	S5filter            *string
	S5ftfile            *string
)

func ParsFlag() {
	Model = flag.String("type", "", "client or server")
	ListenSocks5 = flag.String("l", "", "socks5 listen addr")
	ListenTCP = flag.String("ltcp", "", "tcp listen addr")
	ListenHTTP = flag.String("lhttp", "", "http listen addr")
	Target = flag.String("t", "", "target addr")
	Server = flag.String("s", "", "server addr")
	Timeout = flag.Int("timeout", 60, "conn timeout")
	Key = flag.Int("key", 0, "key")
	Tcpmode = flag.Int("tcp", 0, "tcp mode")
	TcpmodeBuffersize = flag.Int("tcp_bs", 1*1024*1024, "tcp mode buffer size")
	TcpmodeMaxwin = flag.Int("tcp_mw", 20000, "tcp mode max win")
	TcpmodeResendTimems = flag.Int("tcp_rst", 400, "tcp mode resend time ms")
	TcpmodeCompress = flag.Int("tcp_gz", 0, "tcp data compress")
	Nolog = flag.Int("nolog", 0, "write log file")
	Noprint = flag.Int("noprint", 0, "print stdout")
	TcpmodeStatus = flag.Int("tcp_stat", 0, "print tcp stat")
	Loglevel = flag.String("loglevel", "info", "log level")
	OpenSock5 = flag.Int("sock5", 0, "sock5 mode")
	Maxconn = flag.Int("maxconn", 0, "max num of connections")
	MaxProcessThread = flag.Int("maxprt", 100, "max process thread in server")
	MaxProcessBuffer = flag.Int("maxprb", 1000, "max process thread's buffer in server")
	Profile = flag.Int("profile", 0, "open profile")
	Conntt = flag.Int("conntt", 1000, "the connect call's timeout")
	S5filter = flag.String("s5filter", "", "sock5 filter")
	S5ftfile = flag.String("s5ftfile", "GeoLite2-Country.mmdb", "sock5 filter file")
	flag.Parse()
}
