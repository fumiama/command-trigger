package cmd

import (
	"net"
	"time"

	"github.com/fumiama/command-trigger/packet"
	"github.com/sirupsen/logrus"
)

type Executer struct {
	key *[64]byte
	do  func(remo net.Addr, text string) string
}

func NewExecuter(key *[64]byte, waitsec time.Duration, do func(remo net.Addr, text string) string) (e Executer) {
	e.key = key
	e.do = do
	return
}

func (e *Executer) Execute(remo net.Addr, echo uint64, text string) error {
	conn, err := net.Dial(remo.Network(), remo.String())
	if err != nil {
		return err
	}
	defer conn.Close()
	p := packet.NewPacket(e.key, echo, e.do(remo, text))
	logrus.Debugln("write exec reply to", remo, ":", p.Text())
	_, err = conn.Write(p.Bytes())
	return err
}
