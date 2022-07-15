package cmd

import (
	"math/rand"
	"net"
	"time"
	"unsafe"

	"github.com/fumiama/command-trigger/math"
	"github.com/fumiama/command-trigger/packet"
	"github.com/sirupsen/logrus"
)

type Trigger struct {
	key              *[64]byte
	maxValidTimeDiff int64
}

func NewTrigger(key *[64]byte, maxValidTimeDiff int64) (t Trigger) {
	t.key = key
	t.maxValidTimeDiff = maxValidTimeDiff
	return
}

func (t *Trigger) Trigger(remo net.Addr, text string) error {
	addr, err := net.ResolveUDPAddr(remo.Network(), remo.String())
	if err != nil {
		return err
	}
	conn, err := net.DialUDP(remo.Network(), nil, addr)
	if err != nil {
		return err
	}
	echo := rand.Uint64()
	p := packet.NewPacket(t.key, echo, text)
	logrus.Debugln("send trigger to", remo, ":", text)
	_, err = conn.Write(p.Bytes())
	if err != nil {
		return err
	}
	defer conn.Close()
	r := (*[packet.PacketSize]byte)(unsafe.Pointer(&p))
	n, err := conn.Read(r[:])
	if err != nil {
		return err
	}
	if n != len(r) || !p.IsValid(t.key) || p.Echo() != echo {
		logrus.Debugln("invalid reply packet from", remo)
		return nil
	}
	newt := time.Now().Unix()
	if math.Abs64(newt-p.Unix()) > t.maxValidTimeDiff {
		logrus.Debug("timeout reply packet from", remo)
		return nil
	}
	logrus.Infoln(remo, "reply:", p.Text())
	return nil
}
