package main

import (
	"fmt"
	"net"
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
	conn, err := net.DialTimeout("tcp", "127.0.0.1:18888", time.Second*3)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	sleep := time.Millisecond * time.Duration(5000)

	login := Pack(LOGIN, "1.0 KC bid")
	fmt.Printf("login:% X\n", login)
	heartbeat := Pack(HEARTBEAT, "\x00")
	fmt.Println("heartbeat.len:", len(heartbeat))

	err = WritePacket(conn, login)
	if err != nil {
		return
	}

	err = WritePacket(conn, heartbeat)
	if err != nil {
		return
	}

	for {
		conn.SetDeadline(time.Now().Add(10 * time.Second))
		pkt, err := ReadPacket(conn)
		if err != nil {
			return
		}

		if pkt.Type == HEARTBEAT {
			err = WritePacket(conn, heartbeat)
			fmt.Println("heartbeat.")
			if err != nil {
				return
			}
			time.Sleep(sleep)
		} else if pkt.Type == COMMAND {
			//TODO
			fmt.Println("command:", pkt.Body)
		} else if pkt.Type == STATUS {
			//TODO
			fmt.Println("upload status.")
		} else {
			continue
		}
	}

}
