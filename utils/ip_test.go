package utils

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestIsPrivateIPv4(t *testing.T) {
	assert.False(t, IsPrivateIPv4(nil))
	assert.False(t, IsPrivateIPv4([]byte{1, 1, 1, 1}))
	assert.True(t, IsPrivateIPv4([]byte{192, 168, 1, 1}))
}

func TestGetPublicIP(t *testing.T) {
	ip, err := GetPublicIP()
	assert.NoError(t, err)
	assert.NotNil(t, ip)
}

func TestGetLocalIP(t *testing.T) {
	_, err := GetLocalIP(true)
	assert.NoError(t, err)

	_, err = GetLocalIP(false)
	assert.NoError(t, err)
}

func TestIsIPBetween(t *testing.T) {
	is, err := IsIPBetween([]byte{192, 168, 0, 0}, []byte{192, 168, 255, 255}, []byte{192, 168, 2, 3})
	assert.NoError(t, err)
	assert.True(t, is)

	is, err = IsIPBetween([]byte{192, 168, 0, 0}, []byte{192, 168, 1, 1}, []byte{192, 168, 2, 3})
	assert.NoError(t, err)
	assert.False(t, is)

	// err
	is, err = IsIPBetween(nil, []byte{192, 168, 255, 255}, []byte{192, 168, 2, 3})
	assert.Error(t, err)
	assert.False(t, is)

	is, err = IsIPBetween([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, []byte{192, 168, 255, 255}, []byte{192, 168, 2, 3})
	assert.Error(t, err)
	assert.False(t, is)
}

func TestIPv4AtoN(t *testing.T) {
	assert.Equal(t, int64(1)<<24|int64(1)<<16|int64(1)<<8|int64(1), IPv4AtoN([]byte{1, 1, 1, 1}))
	// invalid
	assert.Equal(t, int64(0), IPv4AtoN(nil))
}

func TestIPv4NtoA(t *testing.T) {
	assert.Equal(t, net.IP{1, 1, 1, 1}, IPv4NtoA(int64(1)<<24|int64(1)<<16|int64(1)<<8|int64(1)))
	assert.Equal(t, net.IP{0, 0, 0, 0}, IPv4NtoA(0))
}

func TestIPv6AtoN(t *testing.T) {
	// nil
	assert.Equal(t, "", IPv6AtoN(nil))
	// v4
	assert.Equal(t, "01010101", IPv6AtoN(net.IP{1, 1, 1, 1}))
	// v6
	assert.Equal(t, "01010101010101010101010101010101", IPv6AtoN(net.IP{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
}

func TestIPv6NtoA(t *testing.T) {
	// v4
	assert.Equal(t, net.IP{1, 1, 1, 1}, IPv6NtoA("01010101"))
	// v6
	assert.Equal(t, net.IP{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, IPv6NtoA("01010101010101010101010101010101"))
	// nil
	assert.Nil(t, IPv6NtoA("010101010101"))
	assert.Nil(t, IPv6NtoA("hi"))
}
