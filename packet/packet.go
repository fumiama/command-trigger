// Package packet is the UDP packet to be sent
package packet

import (
	"bytes"
	"encoding/binary"
	"time"
	"unsafe"

	"github.com/fumiama/blake2b-simd"
)

const (
	PacketSize = 256
	TextSize   = PacketSize - 8 - 8 - 1 - blake2b.Size
)

type Packet struct {
	unix [8]byte
	echo [8]byte
	txtl uint8
	text [TextSize]byte
	hmac [blake2b.Size]byte
}

// NewPacket fills fields, then calculates hmac with key
func NewPacket(key *[64]byte, echo uint64, text string) (t Packet) {
	if len(text) > TextSize {
		panic("text too long")
	}
	binary.LittleEndian.PutUint64(t.unix[:], uint64(time.Now().Unix()))
	binary.LittleEndian.PutUint64(t.echo[:], echo)
	t.txtl = uint8(len(text))
	copy(t.text[:], text)
	h := blake2b.NewMAC(blake2b.Size, key[:])
	_, _ = h.Write((*[PacketSize - blake2b.Size]byte)(unsafe.Pointer(&t))[:])
	_ = h.Sum(t.hmac[:0])
	return
}

// IsValid checks the packet's validity
func (t *Packet) IsValid(key *[64]byte) bool {
	var myhash [blake2b.Size]byte
	h := blake2b.NewMAC(blake2b.Size, key[:])
	_, _ = h.Write((*[PacketSize - blake2b.Size]byte)(unsafe.Pointer(t))[:])
	return bytes.Equal(h.Sum(myhash[:0]), t.hmac[:])
}

// Unix wraps timestamp by second
func (t *Packet) Unix() int64 {
	return int64(binary.LittleEndian.Uint64(t.unix[:]))
}

// Echo uniquely marks a request
func (t *Packet) Echo() uint64 {
	return binary.LittleEndian.Uint64(t.echo[:])
}

// Text is the reply message
func (t *Packet) Text() string {
	return string(t.text[:t.txtl])
}

// Bytes is the packet bytes (no copy)
func (t *Packet) Bytes() []byte {
	return (*[PacketSize]byte)(unsafe.Pointer(t))[:]
}
