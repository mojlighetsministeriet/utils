package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
)

// FetchPublicKey will fetch a PEM encoded §public RSA key over HTTP and return it's struct
func FetchPublicKey(url string) (publicKey *rsa.PublicKey, err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New("Failed to fetch public key")
		return
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	block, _ := pem.Decode(body)
	if block == nil {
		err = errors.New("Unable to decode pem")
		return
	}

	key, _ := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey = key.(*rsa.PublicKey)

	return
}
