package v1

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
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
