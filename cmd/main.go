package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	_ "net/http/pprof"

	"github.com/esrrhs/gohome/common"
	"github.com/esrrhs/gohome/geoip"
	"github.com/esrrhs/gohome/loggo"
	"github.com/esrrhs/pingtunnel"
	"github.com/esrrhs/pingtunnel/http2socks5"
	"github.com/esrrhs/pingtunnel/iptables"
	"github.com/esrrhs/pingtunnel/option"
	"github.com/esrrhs/pingtunnel/yellowsocks"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	defer common.CrashLog()

	option.ParsFlag()

	if *option.Model != "client" && *option.Model != "server" {
		flag.Usage()
		return
	}
	if *option.Model == "client" {
		if len(*option.ListenSocks5) == 0 || len(*option.Server) == 0 {
			flag.Usage()
			return
		}
		if *option.OpenSock5 == 0 && len(*option.Target) == 0 {
			flag.Usage()
			return
		}
		if *option.OpenSock5 != 0 {
			*option.Tcpmode = 1
			go http2socks5.Run()
			go yellowsocks.Run()
		}
	}
	if *option.TcpmodeMaxwin*10 > pingtunnel.FRAME_MAX_ID {
		fmt.Println("set tcp win to big, max = " + strconv.Itoa(pingtunnel.FRAME_MAX_ID/10))
		return
	}

	level := loggo.LEVEL_INFO
	if loggo.NameToLevel(*option.Loglevel) >= 0 {
		level = loggo.NameToLevel(*option.Loglevel)
	}
	loggo.Ini(loggo.Config{
		Level:     level,
		Prefix:    "pingtunnel",
		MaxDay:    1,
		NoLogFile: *option.Nolog > 0,
		NoPrint:   *option.Noprint > 0,
	})
	loggo.Info("start...")
	loggo.Info("key %d", *option.Key)

	if *option.Model == "server" {
		s, err := pingtunnel.NewServer(*option.Key, *option.Maxconn, *option.MaxProcessThread, *option.MaxProcessBuffer, *option.Conntt)
		if err != nil {
			loggo.Error("ERROR: %s", err.Error())
			return
		}
		loggo.Info("Server start")
		err = s.Run()
		if err != nil {
			loggo.Error("Run ERROR: %s", err.Error())
			return
		}
	} else if *option.Model == "client" {

		loggo.Info("type %s", *option.Model)
		loggo.Info("listen %s", *option.ListenSocks5)
		loggo.Info("server %s", *option.Server)
		loggo.Info("target %s", *option.Target)

		if *option.Tcpmode == 0 {
			*option.TcpmodeBuffersize = 0
			*option.TcpmodeMaxwin = 0
			*option.TcpmodeResendTimems = 0
			*option.TcpmodeCompress = 0
			*option.TcpmodeStatus = 0
		}

		if len(*option.S5filter) > 0 {
			err := geoip.Load(*option.S5ftfile)
			if err != nil {
				loggo.Error("Load Sock5 ip file ERROR: %s", err.Error())
				return
			}
		}
		filter := func(addr string) bool {
			if len(*option.S5filter) <= 0 {
				return true
			}

			taddr, err := net.ResolveTCPAddr("tcp", addr)
			if err != nil {
				return false
			}

			ret, err := geoip.GetCountryIsoCode(taddr.IP.String())
			if err != nil {
				return false
			}
			if len(ret) <= 0 {
				return false
			}
			return ret != *option.S5filter
		}

		c, err := pingtunnel.NewClient(*option.ListenSocks5, *option.Server, *option.Target, *option.Timeout, *option.Key,
			*option.Tcpmode, *option.TcpmodeBuffersize, *option.TcpmodeMaxwin, *option.TcpmodeResendTimems, *option.TcpmodeCompress,
			*option.TcpmodeStatus, *option.OpenSock5, *option.Maxconn, &filter)
		if err != nil {
			loggo.Error("ERROR: %s", err.Error())
			return
		}
		loggo.Info("Client Listen %s (%s) Server %s (%s) TargetPort %s:", c.Addr(), c.IPAddr(),
			c.ServerAddr(), c.ServerIPAddr(), c.TargetAddr())
		err = c.Run()
		go iptables.Run()
		if err != nil {
			loggo.Error("Run ERROR: %s", err.Error())
			return
		}
	} else {
		return
	}

	if *option.Profile > 0 {
		go http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*option.Profile), nil)
	}

	for {
		time.Sleep(time.Hour)
	}
}
