package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload_Integration(t *testing.T) {
	provider := NewWebshareProxyProvider()

	// Act
	proxies, err := provider.Download()

	assert.Nil(t, err, "Download should not return an error")

	assert.NotEmpty(t, proxies, "The downloaded proxy list should not be empty")

	assert.Equal(t, 10, len(proxies))

	// Verify the format of a few proxies. This checks if the parsing logic is correct.
	for _, p := range proxies {
		assert.NotEmpty(t, p.Host, "Proxy host should not be empty")
		assert.NotEmpty(t, p.Port, "Proxy port should not be empty")
		assert.NotEmpty(t, p.User, "Proxy user should not be empty")
		assert.NotEmpty(t, p.Password, "Proxy password should not be empty")
	}

}
