package main

import (
	"bytes"
	"encoding/binary"
	"net"
	//"fmt"
)

const (
	HEARTBEAT = iota
	COMMAND
	STATUS
	LOGIN
)

type Pkt struct {
	Len  uint16
	Type uint8
	Body string
}

func (pkt *Pkt) pack() []byte {
	_body := []byte(pkt.Body)
	_len := len(_body) + 3

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint16(_len))
	binary.Write(buf, binary.BigEndian, pkt.Type)
	binary.Write(buf, binary.BigEndian, _body)

	return buf.Bytes()
}

func Pack(_type uint8, body string) []byte {
	p := Pkt{0, _type, body}
	return p.pack()
}

func Unpack(buf []byte) *Pkt {
	_len := binary.BigEndian.Uint16(buf[0:2])
	_type := uint8(buf[2])
	body := string(buf[3:])
	return &Pkt{_len, _type, body}
}

func WritePacket(conn net.Conn, buf []byte) error {
	bufLen := len(buf)
	bufPos := 0
	for {
		n, err := conn.Write(buf[bufPos : bufPos+bufLen])
		if err != nil {
			return err
		}
		bufLen -= n
		if bufLen <= 0 {
			break
		}

		bufPos += n
	}
	return nil
}

func ReadPacket(conn net.Conn) (*Pkt, error) {

	headLen := uint16(3)
	head := make([]byte, headLen)
	n, err := conn.Read(head)
	if err != nil || uint16(n) != headLen {
		return nil, err
	}

	pktLen := binary.BigEndian.Uint16(head[0:2])

	buf := make([]byte, pktLen)
	bufPos := headLen
	buf[0] = head[0]
	buf[1] = head[1]
	buf[2] = head[2]

	last_len := pktLen - headLen
	for {
		n, err := conn.Read(buf[bufPos : bufPos+last_len])
		if err != nil {
			return nil, err
		}
		bufPos += uint16(n)
		last_len -= uint16(n)
		if last_len <= 0 {
			break
		}
	}
	return Unpack(buf), nil
}

/*
func main() {
	res := Pack(HEARTBEAT, "\x43")
	fmt.Printf("% X\n", res)
	up := Unpack(res)
	fmt.Println(up)

	res = Pack(HEARTBEAT, "Spring is great")
	fmt.Printf("% X\n", res)
	up = Unpack(res)
	fmt.Println(up)
}
*/
