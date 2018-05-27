package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/mojlighetsministeriet/utils/jwt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type Account struct {
	ID    string
	Email string
	Roles []string
}

func (account *Account) GetID() string {
	return account.ID
}

func (account *Account) GetEmail() string {
	return account.Email
}

func (account *Account) GetRolesSerialized() string {
	return strings.Join(account.Roles, ",")
}

func TestGenerateAndParseIfValid(test *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(test, err)

	account := Account{
		ID:    uuid.Must(uuid.NewV4()).String(),
		Email: "tech+testing@mojlighetsministerietest.se",
		Roles: []string{"user"},
	}

	accessToken, err := jwt.Generate("test-service", privateKey, &account)
	assert.NoError(test, err)

	parsedToken, err := jwt.ParseIfValid(&privateKey.PublicKey, accessToken)
	assert.NoError(test, err)
	assert.Equal(test, account.ID, parsedToken.Claims().Get("sub").(string))
	assert.Equal(test, "tech+testing@mojlighetsministerietest.se", parsedToken.Claims().Get("email"))
	assert.Equal(test, "user", parsedToken.Claims().Get("roles"))
}

func TestFailParseIfValidWithBadToken(test *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(test, err)

	parsedToken, err := jwt.ParseIfValid(&privateKey.PublicKey, []byte{1, 2, 3, 4})
	assert.Error(test, err)
	assert.Equal(test, nil, parsedToken)
}

func TestFailParseIfValidWithBadPublicKey(test *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(test, err)

	account := Account{
		ID:    uuid.Must(uuid.NewV4()).String(),
		Email: "tech+testing@mojlighetsministerietest.se",
		Roles: []string{"user"},
	}

	accessToken, err := jwt.Generate("test-service", privateKey, &account)
	assert.NoError(test, err)

	wrongPrivateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(test, err)

	parsedToken, err := jwt.ParseIfValid(&wrongPrivateKey.PublicKey, accessToken)
	assert.Error(test, err)
	assert.Equal(test, false, parsedToken.Claims().Has("email"))
}

func TestGetClaimsFromContextIfValid(test *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(test, err)

	account := Account{
		ID:    uuid.Must(uuid.NewV4()).String(),
		Email: "tech+testing@mojlighetsministerietest.se",
		Roles: []string{"user"},
	}

	accessToken, err := jwt.Generate("test-service", privateKey, &account)
	assert.NoError(test, err)

	router := echo.New()
	request := httptest.NewRequest(echo.GET, "/", nil)
	request.Header.Add("Authorization", "Bearer "+string(accessToken))
	recorder := httptest.NewRecorder()
	context := router.NewContext(request, recorder)

	claims, err := jwt.GetClaimsFromContextIfValid(&privateKey.PublicKey, context)
	assert.NoError(test, err)
	assert.Equal(test, "tech+testing@mojlighetsministerietest.se", claims.Get("email"))
}

func TestWithInvalidKeyFailToGetClaimsFromContextIfValid(test *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(test, err)

	account := Account{
		ID:    uuid.Must(uuid.NewV4()).String(),
		Email: "tech+testing@mojlighetsministerietest.se",
		Roles: []string{"user"},
	}

	accessToken, err := jwt.Generate("test-service", privateKey, &account)
	assert.NoError(test, err)

	router := echo.New()
	request := httptest.NewRequest(echo.GET, "/", nil)
	request.Header.Add("Authorization", "Bearer "+string(accessToken))
	recorder := httptest.NewRecorder()
	context := router.NewContext(request, recorder)

	wrongPrivateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(test, err)

	claims, err := jwt.GetClaimsFromContextIfValid(&wrongPrivateKey.PublicKey, context)
	assert.Error(test, err)
	assert.Equal(test, nil, claims.Get("email"))
}
