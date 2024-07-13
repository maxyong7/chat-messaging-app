package v1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func getUserIDFromContext(c *gin.Context) (string, error) {
	userID, ok := c.Get("user_id")

	// userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	userIDStr, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("user ID is invalid type")
	}
	return userIDStr, nil
}

func queryParamInt(c *gin.Context, name string, defaultvalue int) int {
	param := c.Param(name)
	result, err := strconv.Atoi(param)
	if err != nil {
		return defaultvalue
	}
	return result
}

func queryParamCursor(c *gin.Context) (time.Time, error) {
	cursor := c.Param("cursor")
	if cursor == "" {
		return time.Now(), nil
	}

	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, err
	}

	var cur time.Time
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return time.Time{}, err
	}
	return cur, nil
}

func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		tokenString = c.GetHeader("Sec-Websocket-Protocol")
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString = strings.Replace(tokenString, "access_token, ", "", 1)
	} else {
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// c.Request.SetPathValue("user_id", userID)
	c.Set("user_id", userID)
	c.Next()
}

func createToken(userUuid string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userUuid
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}
