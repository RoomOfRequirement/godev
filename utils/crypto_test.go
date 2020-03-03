package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	plaintext := []byte("this is a plaintext")
	key := []byte("abcd1234efgh7890") // 16 bytes
	cipherStr, err := Encrypt(plaintext, key)
	assert.NoError(t, err)
	assert.NotEmpty(t, cipherStr)
	t.Log(cipherStr)

	// invalid key
	key = []byte("")
	cipherStr, err = Encrypt(plaintext, key)
	assert.Error(t, err)
	assert.Empty(t, cipherStr)
}

func TestDecrypt(t *testing.T) {
	cipherStr := "Jti0Wp7MxK20TJwN00oCeA1lqd9wwPt6RyiGYBIWG3zWfFbV8/yv4RaFZF/MlnXO"
	key := []byte("abcd1234efgh7890")
	plainStr, err := Decrypt(cipherStr, key)
	assert.NoError(t, err)
	assert.Equal(t, "this is a plaintext", plainStr)

	// invalid key
	key = []byte("")
	plainStr, err = Decrypt(cipherStr, key)
	assert.Error(t, err)
	assert.Empty(t, plainStr)

	// wrong key
	key = []byte("efgh7890abcd1234")
	plainStr, err = Decrypt(cipherStr, key)
	assert.Error(t, err)
	assert.Empty(t, plainStr)

	// wrong cipherStr length, shorter than block size
	cipherStr = "Jti0Wp7MxK20"
	key = []byte("abcd1234efgh7890")
	plainStr, err = Decrypt(cipherStr, key)
	assert.Error(t, err)
	assert.Empty(t, plainStr)
	// wrong cipherStr length, illegal base64 data
	cipherStr = "Jti0Wp7MxK20TJwN00oCeA1lqd9wwPt6RyiGYBIWG3zWfFbV8/yv4RaFZF/"
	key = []byte("abcd1234efgh7890")
	plainStr, err = Decrypt(cipherStr, key)
	assert.Error(t, err)
	assert.Empty(t, plainStr)
	// wrong cipherStr length, not a multiple of the block size
	cipherStr = "Jti0Wp7MxK20TJwN00oCeA1lqd9wwPt6RyiGYBIWG3zWfFbV"
	key = []byte("abcd1234efgh7890")
	plainStr, err = Decrypt(cipherStr, key)
	assert.Error(t, err)
	assert.Empty(t, plainStr)
}

// Benchmark_CBCEncrypt-4   	  487633	      2117 ns/op	     896 B/op	      12 allocs/op
func Benchmark_CBCEncrypt(b *testing.B) {
	plainText := []byte(GenerateRandomString(32))
	key := []byte(GenerateRandomString(16))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Encrypt(plainText, key)
	}
}

// Benchmark_CBCDecrypt-4   	 1551403	       684 ns/op	     672 B/op	       9 allocs/op
func Benchmark_CBCDecrypt(b *testing.B) {
	cipherText := "Jti0Wp7MxK20TJwN00oCeA1lqd9wwPt6RyiGYBIWG3zWfFbV8/yv4RaFZF/MlnXO"
	key := []byte("abcd1234efgh7890")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decrypt(cipherText, key)
	}
}
