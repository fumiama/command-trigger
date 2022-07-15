package packet

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fumiama/blake2b-simd"
	"github.com/stretchr/testify/assert"
)

func TestPacket(t *testing.T) {
	var key [blake2b.Size]byte
	_, _ = rand.Read(key[:])
	p := NewPacket(&key, 123456, "hello world")
	assert.Equal(t, true, p.IsValid(&key))
	assert.Equal(t, time.Now().Unix(), p.Unix())
	assert.Equal(t, "hello world", p.Text())
	p.echo = [8]byte{}
	assert.Equal(t, false, p.IsValid(&key))
}
