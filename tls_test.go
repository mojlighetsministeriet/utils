package utils_test

import (
	"crypto/tls"
	"testing"

	"github.com/mojlighetsministeriet/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetCACertificatesTLSConfig(test *testing.T) {
	config, err := utils.GetCACertificatesTLSConfig()
	assert.NoError(test, err)
	assert.Equal(test, true, len(config.RootCAs.Subjects()) > 10)
}

func TestGetCACertificatesTLSConfigFromFilename(test *testing.T) {
	config, err := utils.GetTLSConfigFromFilename("/etc/ssl/certs/ca-certificates.crt")
	assert.NoError(test, err)
	assert.Equal(test, true, len(config.RootCAs.Subjects()) > 10)
}

func TestFailGetCACertificatesTLSConfigFromFilenameWithBadFilename(test *testing.T) {
	var expectedOutput *tls.Config
	config, err := utils.GetTLSConfigFromFilename("/badfilename.crt")
	assert.Error(test, err)
	assert.Equal(test, expectedOutput, config)
}
