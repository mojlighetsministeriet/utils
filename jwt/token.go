package jwt

import (
	"crypto/rsa"
	"strings"
	"time"

	"github.com/labstack/echo"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	josejwt "github.com/SermoDigital/jose/jwt"
)

// Account describes an account used to generate a token
type Account interface {
	GetID() string
	GetEmail() string
	GetRolesSerialized() string
}

// GenerateWithCustomExpiration generates a new JWT token from an account with a custom expiration time
func GenerateWithCustomExpiration(issuer string, privateKey *rsa.PrivateKey, account Account, expiration time.Time) (serializedToken []byte, err error) {
	claims := jws.Claims{}

	claims.SetExpiration(expiration)
	claims.SetSubject(account.GetID())
	claims.SetIssuer(issuer)
	claims.Set("email", account.GetEmail())
	claims.Set("roles", account.GetRolesSerialized())

	token := jws.NewJWT(claims, crypto.SigningMethodRS256)

	serializedToken, err = token.Serialize(privateKey)

	return
}

// Generate a new JWT token from an account
func Generate(issuer string, privateKey *rsa.PrivateKey, account Account) ([]byte, error) {
	return GenerateWithCustomExpiration(issuer, privateKey, account, time.Now().Add(time.Duration(60*20)*time.Second))
}

// ParseIfValid return a parsed JWT token if it is valid
func ParseIfValid(publicKey *rsa.PublicKey, tokenData []byte) (token josejwt.JWT, err error) {
	token, err = jws.ParseJWT(tokenData)
	if err != nil {
		return
	}

	err = token.Validate(publicKey, crypto.SigningMethodRS256)
	if err != nil {
		claims := jws.Claims{}
		claims.SetExpiration(time.Now().Add(time.Duration(60*20) * time.Second))
		token = jws.NewJWT(claims, crypto.SigningMethodRS256)
	}

	return
}

// GetTokenFromContext will extract the token bytes from the HTTP request header connected to a echo.Context object
func GetTokenFromContext(context echo.Context) (result []byte) {
	token := context.Request().Header.Get("Authorization")
	token = strings.Replace(token, "Bearer", "", -1)
	token = strings.Trim(strings.Replace(token, "bearer", "", -1), " ")

	if len(token) > 20 {
		result = []byte(token)
	}

	return
}

// GetClaimsFromContextIfValid validates the JWT token and fetches the claims from the JWT
func GetClaimsFromContextIfValid(publicKey *rsa.PublicKey, context echo.Context) (claims josejwt.Claims, err error) {
	token, err := ParseIfValid(publicKey, GetTokenFromContext(context))
	if err != nil {
		return
	}

	claims = token.Claims()

	return
}
