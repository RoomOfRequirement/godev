package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUint8(t *testing.T) {
	b := byte('H')
	assert.Equal(t, uint8(0x48), ByteToUint8(b))
	assert.Equal(t, b, Uin8ToByte(0x48))
}

func TestBytesToUint(t *testing.T) {
	b := []byte("Hello World")
	t.Logf("%#04x", b)
	// BE
	assert.Equal(t, uint16(0x4865), BytesToUint16BE(b))
	assert.Equal(t, uint32(0x48656c6c), BytesToUint32BE(b))
	assert.Equal(t, uint64(0x48656c6c6f20576f), BytesToUint64BE(b))
	// LE
	assert.Equal(t, uint16(0x6548), BytesToUint16LE(b))
	assert.Equal(t, uint32(0x6c6c6548), BytesToUint32LE(b))
	assert.Equal(t, uint64(0x6f57206f6c6c6548), BytesToUint64LE(b))
}

func TestUintToBytes(t *testing.T) {
	b := []byte("Hello World")
	t.Logf("%#04x", b)
	// BE
	assert.Equal(t, []byte("He"), Uint16ToBytesBE(BytesToUint16BE(b)))
	assert.Equal(t, []byte("Hell"), Uint32ToBytesBE(BytesToUint32BE(b)))
	assert.Equal(t, []byte("Hello Wo"), Uint64ToBytesBE(BytesToUint64BE(b)))
	// LE
	assert.Equal(t, []byte("He"), Uint16ToBytesLE(BytesToUint16LE(b)))
	assert.Equal(t, []byte("Hell"), Uint32ToBytesLE(BytesToUint32LE(b)))
	assert.Equal(t, []byte("Hello Wo"), Uint64ToBytesLE(BytesToUint64LE(b)))
}
