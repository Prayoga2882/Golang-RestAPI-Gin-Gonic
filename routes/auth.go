package routes

import (
	"fmt"
	"main/config"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/danilopolani/gocialite.v1/structs"
)

var JWT_SECRET = "SECRET_POWER"

func CheckToken(c *gin.Context){
	c.JSON(200, gin.H{
		"Status" : "Succesfully",
	})
}

func RedirectHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")


	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GITHUB"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GITHUB"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		"google": {
			"clientID":     os.Getenv("CLIENT_ID_GOOGLE"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GOOGLE"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{},
		"google": []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

func CallbackHandler(c *gin.Context) {
	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	// Handle callback and check for errors
	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	var newUser = getOrRegister(provider, (*structs.User)(user))

	var jwtToken = getToken(&newUser)

	c.JSON(200, gin.H{
		"Status" 	: "Succesfully",
		"Data" 	 	: newUser,
		"JWT"		: jwtToken,
	})
}

func getOrRegister(provider string, user *structs.User) models.User{
	var userData models.User

	config.DB.Where("provider = ? AND social_id = ?", provider, user.ID).First(&userData)

	if userData.ID == 0 {
		newUser := models.User {
			FullName	: user.FullName,
			Username	: user.Username,
			Email		: user.Email,
			Social_Id	: user.ID,
			Provider	: provider,
			Avatar		: user.Avatar,
		}
		config.DB.Create(&newUser)
		return newUser
	}
	return userData
}

func getToken(user *models.User) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp" : time.Now().AddDate(0, 0, 7).Unix(),
		"iat" : time.Now().Unix(),
	})

	
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}