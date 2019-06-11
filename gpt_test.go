package main

import (
	"github.com/rekby/gpt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader_GetGPTHeader(t *testing.T) {
	f, _ := os.Open("testdata/river.gpt")
	defer f.Close()

	var dummy [0x200]byte
	_, err := f.Read(dummy[:])
	require.NoError(t, err)

	table, err := gpt.ReadTable(f, 512)
	require.NoError(t, err)

	assert.Equal(t, gpt.Header{
		Signature: [8]byte{'E', 'F', 'I', ' ', 'P', 'A', 'R', 'T'},
		Revision:  0x00010000,
		Size:      0x5c,

		CRC:            0xf188a072,
		Reserved:       0,
		HeaderStartLBA: 1,

		DiskGUID: [16]uint8{0x32, 0x1b, 0x10, 0x98, 0xe2, 0xbb, 0xf2, 0x4b, 0xa0, 0x6e, 0x2b, 0xb3, 0x3d, 0x0, 0xc, 0x20},

		FirstUsableLBA: 0x22,
		LastUsableLBA:  0,

		PartitionsTableStartLBA: 2,
		PartitionsArrLen:        64,
		PartitionEntrySize:      128,
		PartitionsCRC:           0x10e3285f,
		TrailingBytes:           make([]byte, 420),
	}, table.Header)
}

func TestReader_GetPartition(t *testing.T) {
	f, _ := os.Open("testdata/river.gpt")
	defer f.Close()

	var dummy [0x200]byte
	_, err := f.Read(dummy[:])
	require.NoError(t, err)

	table, err := gpt.ReadTable(f, 512)
	require.NoError(t, err)

	assert.Equal(t, gpt.Partition{
		Type:          [16]byte{0x2c, 0xba, 0xa0, 0xde, 0xdd, 0xcb, 0x05, 0x48, 0xb4, 0xf9, 0xf4, 0x28, 0x25, 0x1c, 0x3e, 0x98},
		Id:            [16]byte{0xfd, 0x4a, 0xad, 0xac, 0x29, 0xb7, 0x1a, 0x1e, 0xb2, 0x15, 0xc0, 0xf5, 0x7d, 0xf2, 0x93, 0x82},
		FirstLBA:      256,
		LastLBA:       1279,
		Flags:         [8]byte{0x68, 0, 0, 0, 0, 0, 0, 0x10},
		PartNameUTF16: [72]byte{'s', 0, 'b', 0, 'l', 0, '1', 0, '_', 0, 'a'},
		TrailingBytes: make([]byte, 0),
	}, table.Partitions[0])
}

func TestNew(t *testing.T) {
	f, _ := os.Open("testdata/river.gpt")
	defer f.Close()

	var dummy [0x200]byte
	_, err := f.Read(dummy[:])
	require.NoError(t, err)

	_, err = gpt.ReadTable(f, 512)
	require.NoError(t, err)
}

func TestReader_GetGPTHeader_running(t *testing.T) {
	f, _ := os.Open("testdata/river.running.gpt")
	defer f.Close()

	var dummy [0x200]byte
	_, err := f.Read(dummy[:])
	require.NoError(t, err)

	table, err := gpt.ReadTable(f, 512)
	require.NoError(t, err)

	assert.Equal(t, gpt.Header{
		Signature: [8]byte{'E', 'F', 'I', ' ', 'P', 'A', 'R', 'T'},
		Revision:  0x00010000,
		Size:      0x5c,

		CRC:            0xc235ea73,
		Reserved:       0,
		HeaderStartLBA: 1,

		DiskGUID: [16]uint8{0x32, 0x1b, 0x10, 0x98, 0xe2, 0xbb, 0xf2, 0x4b, 0xa0, 0x6e, 0x2b, 0xb3, 0x3d, 0x0, 0xc, 0x20},

		FirstUsableLBA:     0x22,
		LastUsableLBA:      0x0747bfde,
		HeaderCopyStartLBA: 0x0747bfff,

		PartitionsTableStartLBA: 2,
		PartitionsArrLen:        64,
		PartitionEntrySize:      128,
		PartitionsCRC:           0xff8b115d,
		TrailingBytes:           make([]byte, 420),
	}, table.Header)
}

func TestReader_GetPartition_running(t *testing.T) {
	f, _ := os.Open("testdata/river.running.gpt")
	defer f.Close()

	var dummy [0x200]byte
	_, err := f.Read(dummy[:])
	require.NoError(t, err)

	table, err := gpt.ReadTable(f, 512)
	require.NoError(t, err)

	assert.Equal(t, gpt.Partition{
		Type:          [16]byte{0x2c, 0xba, 0xa0, 0xde, 0xdd, 0xcb, 0x05, 0x48, 0xb4, 0xf9, 0xf4, 0x28, 0x25, 0x1c, 0x3e, 0x98},
		Id:            [16]byte{0xfd, 0x4a, 0xad, 0xac, 0x29, 0xb7, 0x1a, 0x1e, 0xb2, 0x15, 0xc0, 0xf5, 0x7d, 0xf2, 0x93, 0x82},
		FirstLBA:      256,
		LastLBA:       1279,
		Flags:         [8]byte{0x68, 0, 0, 0, 0, 0, 0, 0x10},
		PartNameUTF16: [72]byte{'s', 0, 'b', 0, 'l', 0, '1', 0, '_', 0, 'a'},
		TrailingBytes: make([]byte, 0),
	}, table.Partitions[0])
}
