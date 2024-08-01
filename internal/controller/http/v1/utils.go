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

func getUserUUIDFromContext(c *gin.Context) (string, error) {
	userID, ok := c.Get("user_uuid")

	// userID, ok := ctx.Value("user_uuid").(string)
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
	cursor, ok := c.GetQuery("cursor")
	if !ok || cursor == "" {
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

	userID, ok := claims["user_uuid"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// c.Request.SetPathValue("user_uuid", userID)
	c.Set("user_uuid", userID)
	c.Next()
}

func createToken(userUuid string, expirationHour time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_uuid"] = userUuid
	claims["exp"] = time.Now().Add(time.Hour * expirationHour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func encodeCursor(cursor *time.Time) string {
	if cursor == nil {
		return ""
	}
	if cursor.IsZero() {
		return ""
	}
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}
	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor
}
