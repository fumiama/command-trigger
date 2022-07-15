package main

import (
	"flag"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
	"unsafe"

	"github.com/fumiama/blake2b-simd"
	"github.com/fumiama/command-trigger/cmd"
	"github.com/fumiama/command-trigger/math"
	"github.com/fumiama/command-trigger/packet"
	base14 "github.com/fumiama/go-base16384"
	"github.com/sirupsen/logrus"
)

func main() {
	h := flag.Bool("h", false, "display this help")
	i := flag.Int64("i", 10, "min execute interval")
	d := flag.Int64("d", 5, "max valid time diff")
	k := flag.String("k", "", "64 bytes hmac key in base16384 format")
	t := flag.String("t", "", "send/recv trigger to this addr:port")
	e := flag.String("e", "", "execute this command on triggered")
	m := flag.String("m", "", "send additional message")
	g := flag.Bool("g", false, "generate a random key and exit")
	D := flag.Bool("D", false, "debug-level log output")
	flag.Parse()
	if *h {
		flag.Usage()
		os.Exit(0)
	}
	var key [blake2b.Size]byte
	if *g {
		_, err := rand.Read(key[:])
		if err != nil {
			panic(err)
		}
		logrus.Infoln("New HMAC key:", base14.EncodeToString(key[:]))
		os.Exit(0)
	}
	if *k == "" {
		panic("hmac key must not be empty")
	}
	if *t == "" && *e == "" {
		panic("ether -t or -t&-e should be selected")
	}
	if *t == "" && *e != "" {
		panic("-e requires -t to be selected")
	}
	dekey := base14.DecodeFromString(*k)
	if dekey == nil {
		panic("invalid hmac key")
	}
	copy(key[:], dekey)
	if *D {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if *e != "" {
		cmds := strings.Split(*e, " ")
		logrus.Errorln(listen(*t, &key, *i, *d, func(conn net.PacketConn, remo net.Addr, echo uint64, text string) {
			logrus.Infoln(remo, "triggered with message:", text)
			var c *exec.Cmd
			if len(cmds) > 1 {
				c = exec.Command(cmds[0], cmds[1:]...)
			} else {
				c = exec.Command(cmds[0])
			}
			logrus.Debugln("exec cmd:", c.String())
			data, err := c.CombinedOutput()
			msg := string(data)
			if msg == "" {
				msg = "ok"
			}
			if err != nil {
				msg = err.Error()
				if msg == "" {
					msg = "error"
				}
			}
			if len(msg) > packet.TextSize {
				msg = msg[:packet.TextSize]
			}
			logrus.Debugln("get result:", msg)
			p := packet.NewPacket(&key, echo, msg)
			_, err = conn.WriteTo(p.Bytes(), remo)
			if err != nil {
				panic(err)
			}
		}))
	} else {
		remo, err := net.ResolveUDPAddr("udp", *t)
		if err != nil {
			panic(err)
		}
		tg := cmd.NewTrigger(&key, *d)
		err = tg.Trigger(remo, *m)
		if err != nil {
			logrus.Errorln(err)
		}
		os.Exit(0)
	}
}

func listen(addr string, key *[64]byte, minExecuteInterval, maxValidTimeDiff int64, do func(conn net.PacketConn, remo net.Addr, echo uint64, text string)) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	t := packet.Packet{}
	r := (*[packet.PacketSize]byte)(unsafe.Pointer(&t))
	prevdotime := int64(0)
	for {
		n, remo, err := conn.ReadFrom(r[:])
		if err != nil {
			return err
		}
		if n != len(r) || !t.IsValid(key) {
			logrus.Debugln("invalid trigger packet from", remo)
			continue
		}
		newt := time.Now().Unix()
		if math.Abs64(newt-t.Unix()) > maxValidTimeDiff || math.Abs64(newt-prevdotime) < minExecuteInterval {
			logrus.Debugln("timeout trigger packet from", remo)
			continue
		}
		go do(conn, remo, t.Echo(), t.Text())
		prevdotime = newt
	}
}
