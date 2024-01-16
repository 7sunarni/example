package iptables

import (
	"os/exec"
	"strings"
	"time"
)

func Run() {
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ticker.C:
			run()
		}
	}
}

var commands = []string{
	"-t nat -N YELLOWSOCKS",
	"-t nat -A PREROUTING -p tcp -j YELLOWSOCKS",
	"-t nat -A YELLOWSOCKS -d 0.0.0.0/8 -j RETURN",
	"-t nat -A YELLOWSOCKS -d 10.0.0.0/8 -j RETURN",
	"-t nat -A YELLOWSOCKS -d 127.0.0.0/8 -j RETURN",
	"-t nat -A YELLOWSOCKS -d 169.254.0.0/16 -j RETURN",
	"-t nat -A YELLOWSOCKS -d 172.16.0.0/16 -j RETURN",
	"-t nat -A YELLOWSOCKS -d 192.168.0.0/16 -j RETURN",
	"-t nat -A YELLOWSOCKS -d 224.0.0.0/4 -j RETURN",
	"-t nat -A YELLOWSOCKS -d 240.0.0.0/4 -j RETURN",
	"-t nat -A YELLOWSOCKS -p tcp -j REDIRECT --to-ports 4466",
}

func run() {
	for _, command := range commands {
		exec.Command("iptables", strings.Split(command, " ")...).Run()
	}
}
