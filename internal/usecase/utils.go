package usecase

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

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
