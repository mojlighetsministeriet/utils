package jwt

import (
	"crypto/rsa"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/mojlighetsministeriet/identity-provider/entity"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	josejwt "github.com/SermoDigital/jose/jwt"
)

// Generate a new JWT token from a user
func Generate(issuer string, privateKey *rsa.PrivateKey, account entity.Account) (serializedToken []byte, err error) {
	claims := jws.Claims{}

	account.BeforeSave()

	claims.SetExpiration(time.Now().Add(time.Duration(60*20) * time.Second))
	claims.SetSubject(account.ID)
	claims.SetIssuer(issuer)
	claims.Set("email", account.Email)
	claims.Set("roles", account.RolesSerialized)

	token := jws.NewJWT(claims, crypto.SigningMethodRS256)

	serializedToken, err = token.Serialize(privateKey)

	return
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
