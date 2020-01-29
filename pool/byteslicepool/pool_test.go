package byteslicepool

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetrics_String(t *testing.T) {
	m := Metrics{
		gets: 1,
		puts: 2,
		hits: 3,
		size: 5,
	}
	assert.Equal(t, fmt.Sprintf("Byte Slice Pool Metrics: \n\tGets: %d\n\tPuts: %d\n\tHits: %d\n\tSize: %d\n", m.gets, m.puts, m.hits, m.size), m.String())
}

func TestNew(t *testing.T) {
	var _ Pool = (*pool)(nil)

	p := New()
	assert.EqualValues(t, Metrics{}, p.Metrics())

	bs := p.Get(10)
	assert.NotNil(t, bs)
	assert.Equal(t, 10, len(bs))
	assert.EqualValues(t, Metrics{
		gets: 1,
		puts: 0,
		hits: 0,
		size: 0,
	}, p.Metrics())

	p.Put(bs)
	assert.EqualValues(t, Metrics{
		gets: 1,
		puts: 1,
		hits: 0,
		size: 16,
	}, p.Metrics())

	p.Put(bs)
	assert.EqualValues(t, Metrics{
		gets: 1,
		puts: 2,
		hits: 0,
		size: 32,
	}, p.Metrics())

	bs = p.Get(10)
	assert.NotNil(t, bs)
	assert.Equal(t, 10, len(bs))
	assert.EqualValues(t, Metrics{
		gets: 2,
		puts: 2,
		hits: 1,
		size: 16,
	}, p.Metrics())
}
