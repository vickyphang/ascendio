package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	ghOauth "golang.org/x/oauth2/github"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     ghOauth.Endpoint,
	}
	state          = "random"
	appID          = os.Getenv("APP_ID")
	privateKey     *rsa.PrivateKey
	privateKeyPath = os.Getenv("PRIVATE_KEY_PATH")
)

func init() {
	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		panic("failed to decode PEM block containing private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	privateKey = privKey
}

func main() {
	router := gin.Default()

	router.GET("/login", handleLogin)
	router.GET("/callback", handleCallback)

	router.Run(":8080")
}

func handleLogin(c *gin.Context) {
	url := oauthConf.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

func handleCallback(c *gin.Context) {
	code := c.Query("code")
	tok, err := oauthConf.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Now generate JWT for GitHub App
	jwtToken, err := generateJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tok.AccessToken, "jwt": jwtToken})
}

func generateJWT() (string, error) {
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iat": now,
		"exp": now + (10 * 60),
		"iss": appID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
