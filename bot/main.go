package main

import (
	"crypto/tls"
	"fmt"
	//"net"
	"github.com/Acey9/spring/common"
	"os"
	"time"
)

type Bot struct {
	beatQueue chan int
	cmdQueue  chan string
}

func (bot *Bot) sendHeartbeat() error {
	fmt.Println("TODO")
	return nil
}

func (bot *Bot) Exe() error {
	//TODO
	fmt.Println("exe")
	return nil
}

func main() {
	args := os.Args[1:]
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	//conn, err := net.DialTimeout("tcp", args[0], time.Second*3)
	conn, err := tls.Dial("tcp", args[0], conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	sleep := time.Millisecond * time.Duration(5000)

	login := common.Pack(common.LOGIN, []byte("1.0 KC bid"))
	fmt.Printf("login:% X\n", login)
	heartbeat := common.Pack(common.HEARTBEAT, []byte("\x00"))
	fmt.Println("heartbeat.len:", len(heartbeat))

	err = common.WritePacket(conn, login)
	if err != nil {
		fmt.Println(1, err)
		return
	}

	err = common.WritePacket(conn, heartbeat)
	if err != nil {
		fmt.Println(2, err)
		return
	}

	for {
		conn.SetDeadline(time.Now().Add(10 * time.Second))
		pkt, err := common.ReadPacket(conn)
		if err != nil {
			fmt.Println(3, err)
			return
		}

		if pkt.Type == common.HEARTBEAT {
			err = common.WritePacket(conn, heartbeat)
			fmt.Println("heartbeat.")
			if err != nil {
				fmt.Println(4, err)
				return
			}
			time.Sleep(sleep)
		} else if pkt.Type == common.COMMAND {
			//TODO
			fmt.Println("command:", pkt.Body)
		} else if pkt.Type == common.STATUS {
			//TODO
			fmt.Println("upload status.")
		} else {
			continue
		}
	}

}
