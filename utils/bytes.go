package utils

import "encoding/binary"

/* network -> BigEndian, host -> LittleEndian (`ntoh`, `hton`)
 * uint8  -> byte      (byte is alias for uint8)
 * uint16 -> [2]byte
 * uint32 -> [4]byte
 * uint64 -> [8]byte
 */

// ByteToUint8 ...
//	useless: byte == uint8
func ByteToUint8(b byte) uint8 {
	return b
}

// Uin8ToByte ...
//	useless: uint8 == byte
func Uin8ToByte(v uint8) byte {
	return v
}

// BytesToUint16BE ...
// or use binary.Read(buf, order, data)
func BytesToUint16BE(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

// BytesToUint32BE ...
func BytesToUint32BE(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// BytesToUint64BE ...
func BytesToUint64BE(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Uint16ToBytesBE ...
// or use binary.Write(buf, order, data)
func Uint16ToBytesBE(v uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	return b
}

// Uint32ToBytesBE ...
func Uint32ToBytesBE(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}

// Uint64ToBytesBE ...
func Uint64ToBytesBE(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// BytesToUint16LE ...
func BytesToUint16LE(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

// BytesToUint32LE ...
func BytesToUint32LE(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

// BytesToUint64LE ...
func BytesToUint64LE(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

// Uint16ToBytesLE ...
func Uint16ToBytesLE(v uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	return b
}

// Uint32ToBytesLE ...
func Uint32ToBytesLE(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

// Uint64ToBytesLE ...
func Uint64ToBytesLE(v uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	return b
}
