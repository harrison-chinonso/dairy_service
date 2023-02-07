package helper

import (
	"dairy_service/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strconv"
	"strings"
	"time"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

// GenerateJWT function takes a user model and generates a JWT containing the user’s id (id),
// the time at which the token was issued (iat), and the expiry date of the token (eat). Using the JWT_PRIVATE_KEY
// environment variable, a signed JWT is returned as a string.
func GenerateJWT(user model.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

// The ValidateJWT function ensures that the incoming request contains a valid token in the request header. This function will be used by the middleware to ensure that only authenticated requests are allowed past the middleware.
func ValidateJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

// The CurrentUser function will be used to get the user associated with the provided JWT by retrieving the id key from the
// parsed JWT and retrieve the corresponding user from the database.
func CurrentUser(context *gin.Context) (model.User, error) {
	err := ValidateJWT(context)
	if err != nil {
		return model.User{}, err
	}
	token, _ := getToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	user, err := model.FindUserById(userId)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// The getToken function uses the returned token string to parse the JWT, using the private key specified in .env.local.
func getToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

// The getTokenFromRequest function retrieves the bearer token from the request. Bearer tokens come in the format bearer <JWT>,
// hence the retrieved string is split and the JWT string is returned.
func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
