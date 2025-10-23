package helper

import (
	"context"
	"report-service/pkg/constants"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func ToBsonM(v interface{}) bson.M {
	if m, ok := v.(bson.M); ok {
		return m
	}
	if m, ok := v.(map[string]interface{}); ok {
		return bson.M(m)
	}
	return bson.M{}
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func GetLatestTimeStr(updatedAt, managerUpdatedAt string) string {
	if updatedAt == "" && managerUpdatedAt == "" {
		return ""
	}

	layouts := []string{
		"2006-01-02T15:04:05.000Z07:00",
		"2006-01-02T15:04:05.000000",
		"2006-01-02T15:04:05.000",
	}

	parseTime := func(s string) time.Time {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, s); err == nil {
				return t
			}
		}
		return time.Time{}
	}

	t1 := parseTime(updatedAt)
	t2 := parseTime(managerUpdatedAt)

	if t1.IsZero() && t2.IsZero() {
		return ""
	}

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	var latest time.Time
	if t1.After(t2) {
		latest = t1
	} else {
		latest = t2
	}

	// Format về dạng: "YYYY-MM-DD HH:mm:ss"
	return latest.In(loc).Format("2006-01-02 15:04:05")
}
