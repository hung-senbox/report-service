package cache

import "context"

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
}

// Helper tạo key nhất quán
func UserCacheKey(userID string) string {
	return "user:" + userID
}

func StudentCacheKey(studenID string) string {
	return "student:" + studenID
}

func TeacherCacheKey(teacherID string) string {
	return "teacher:" + teacherID
}

func StaffCacheKey(staffID string) string {
	return "staff:" + staffID
}
