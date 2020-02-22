package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewSnowflake(t *testing.T) {
	// invalid start time
	sf, err := NewSnowflake(SnowFlakeOptions{
		StartTime: time.Now().Add(time.Hour),
		MachineID: nil,
	})
	assert.Nil(t, sf)
	assert.Error(t, err)

	// invalid machine id func
	sf, err = NewSnowflake(SnowFlakeOptions{
		StartTime: time.Time{},
		MachineID: nil,
	})
	assert.Nil(t, sf)
	assert.Error(t, err)

	// invalid machine id
	sf, err = NewSnowflake(SnowFlakeOptions{
		StartTime: time.Time{},
		MachineID: func() uint16 {
			return MaxMachineID + 1
		},
	})
	assert.Nil(t, sf)
	assert.Error(t, err)

	// default start time
	sf, err = NewSnowflake(SnowFlakeOptions{
		StartTime: time.Time{},
		MachineID: func() uint16 {
			return MaxMachineID - 1
		},
	})
	assert.NotNil(t, sf)
	assert.NoError(t, err)
	assert.Equal(t, int64(DefaultStartTime), sf.startTime)
	assert.Equal(t, uint16(MaxMachineID-1), sf.machineID)

	// custom start time
	st := time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC)
	sf, err = NewSnowflake(SnowFlakeOptions{
		StartTime: st,
		MachineID: func() uint16 {
			return 0
		},
	})
	assert.NotNil(t, sf)
	assert.NoError(t, err)
	assert.Equal(t, timeToUint64(st), sf.startTime)
	assert.Equal(t, uint16(0), sf.machineID)
}

func TestSnowFlake_NextUID(t *testing.T) {
	sf, err := NewSnowflake(SnowFlakeOptions{
		StartTime: time.Time{},
		MachineID: func() uint16 {
			return MaxMachineID - 1
		},
	})
	assert.NoError(t, err)
	uid, err := sf.NextUID()
	assert.NoError(t, err)
	assert.LessOrEqual(t, uint64(7471237110495232), uid)

	// max sequence
	sf.sequence = MaxSequenceBits
	uid, err = sf.NextUID()
	assert.NoError(t, err)

	// err
	sf.startTime -= int64(time.Duration(70) * time.Duration(365*24) * time.Hour) // 70 years
	uid, err = sf.NextUID()
	assert.Error(t, err)
	assert.Equal(t, uint64(0), uid)
}
