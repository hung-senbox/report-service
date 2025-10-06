package helper

import (
	"context"
	"report-service/pkg/constants"
	"strconv"
	"strings"
)

func ParseAppLanguage(header string, defaultVal uint) uint {
	header = strings.TrimSpace(strings.Trim(header, "\""))
	if val, err := strconv.Atoi(header); err == nil {
		return uint(val)
	}
	return defaultVal
}

func GetHeaders(ctx context.Context) map[string]string {
	headers := make(map[string]string)

	if lang, ok := ctx.Value(constants.AppLanguage).(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	return headers
}

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(constants.UserID).(string); ok {
		return userID
	}
	return ""
}
