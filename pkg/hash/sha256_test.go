package hash

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHashSHA256(t *testing.T) {
	tests := []struct {
		key      string
		data     []byte
		expected string
	}{
		{
			key:      "test-key",
			data:     []byte("test data"),
			expected: "6ba8fe37b46e711c155c1d65fc59c2826255bff64b7a2b385d5588577d472160",
		},
		{
			key:      "another-key",
			data:     []byte("some other data"),
			expected: "0b0ada949c490d79a152d2bacd4472eabbdc71494936dbafd75c35bc6c6af09e",
		},
		{
			key:      "",
			data:     []byte("data with empty key"),
			expected: "27dbe8e637a12fcedc552951be7bc4eeb9118c5796706c1be24e7a2a921be3f1",
		},
		{
			key:      "key",
			data:     []byte(""),
			expected: "2c70e12b7a0646f92279f427c7b38e7334d8e5389cff167a1dc30e73f826b683",
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			hash := GetHashSHA256(tt.key, tt.data)
			hashHex := hex.EncodeToString(hash)
			assert.Equal(t, tt.expected, hashHex)
		})
	}
}
