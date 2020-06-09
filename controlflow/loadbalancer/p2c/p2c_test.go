package p2c

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewP2C(t *testing.T) {
	p := NewP2C()
	assert.Equal(t, 0, len(p.nodes))
	assert.Equal(t, 0, len(p.nodesSet))
}

func TestP2C_AddNode(t *testing.T) {
	p := NewP2C()
	p.AddNode("0.0.0.0", 0)
	assert.Equal(t, 1, len(p.nodes))
	assert.Equal(t, 1, len(p.nodesSet))

	p.AddNode("0.0.0.0", 10)
	assert.Equal(t, 1, len(p.nodes))
	assert.Equal(t, 1, len(p.nodesSet))
	l, err := p.GetLoad("0.0.0.0")
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), l)
}

func TestP2C_DeleteNode(t *testing.T) {
	p := NewP2C()
	p.AddNode("0.0.0.0", 0)
	assert.Equal(t, 1, len(p.nodes))
	assert.Equal(t, 1, len(p.nodesSet))
	p.DeleteNode("0.0.0.1")
	assert.Equal(t, 1, len(p.nodes))
	assert.Equal(t, 1, len(p.nodesSet))
	p.DeleteNode("0.0.0.0")
	assert.Equal(t, 0, len(p.nodes))
	assert.Equal(t, 0, len(p.nodesSet))
}

func TestP2C_Get(t *testing.T) {
	p := NewP2C()
	n, err := p.Get("")
	assert.Error(t, ErrNoNodes, err)
	assert.Equal(t, "", n)
	p.AddNode("0.0.0.0", 0)
	p.AddNode("0.0.0.1", 10)
	n, err = p.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.1", n)
	_, err = p.Get("")
	assert.NoError(t, err)
}

func TestP2C_UpdateLoad(t *testing.T) {
	p := NewP2C()
	p.AddNode("0.0.0.0", 0)
	p.AddNode("0.0.0.1", 10)
	err := p.IncrLoad("0.0.0.0")
	assert.NoError(t, err)
	err = p.DecrLoad("0.0.0.1")
	assert.NoError(t, err)
	load, err := p.GetLoad("0.0.0.0")
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), load)
	load, err = p.GetLoad("0.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, uint64(9), load)

	err = p.IncrLoad("0.0.0.2")
	assert.Error(t, ErrNodeNotExist, err)
	err = p.DecrLoad("0.0.0.2")
	assert.Error(t, ErrNodeNotExist, err)
	load, err = p.GetLoad("0.0.0.2")
	assert.Error(t, ErrNodeNotExist, err)
	assert.Equal(t, uint64(0), load)

	err = p.UpdateLoad("0.0.0.0", 10)
	assert.NoError(t, err)
	load, err = p.GetLoad("0.0.0.0")
	assert.Equal(t, uint64(10), load)
	err = p.UpdateLoad("0.0.0.2", 10)
	assert.Error(t, ErrNodeNotExist, err)
}
