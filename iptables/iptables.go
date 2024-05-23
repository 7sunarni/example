package iptables

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func Run() {
	ticker := time.NewTicker(time.Hour)
	run()
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
	c := exec.Command("sh", "-c", "iptables-save", "|", "grep", "YELLOWSOCKS")
	d := bytes.NewBuffer(make([]byte, 0))
	c.Stdout = d
	if err := c.Run(); err != nil {
		fmt.Println(err)
		return
	}
	if strings.Contains(d.String(), "YELLOWSOCKS") {
		return
	}
	for _, command := range commands {
		exec.Command("iptables", strings.Split(command, " ")...).Run()
	}
}
