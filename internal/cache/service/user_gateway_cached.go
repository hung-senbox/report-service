package service

import (
	"context"
	"fmt"
	"time"

	"report-service/internal/cache"
	"report-service/internal/gateway"
	gw_response "report-service/internal/gateway/dto/response"
)

// CachedUserGateway wraps UserGateway with Redis caching.
type CachedUserGateway struct {
	inner gateway.UserGateway
	cache cache.Cache
	ttl   time.Duration
}

func NewCachedUserGateway(inner gateway.UserGateway, cache cache.Cache) gateway.UserGateway {
	return &CachedUserGateway{
		inner: inner,
		cache: cache,
		ttl:   30 * time.Minute,
	}
}

// ==============================
// === Example: GetUserInfo ===
// ==============================
func (g *CachedUserGateway) GetUserInfo(ctx context.Context, userID string) (*gw_response.UserInfo, error) {
	cacheKey := fmt.Sprintf("user:info:%s", userID)

	var cached gw_response.UserInfo
	if err := g.cache.Get(ctx, cacheKey, &cached); err == nil && cached.ID != "" {
		return &cached, nil
	}

	user, err := g.inner.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	_ = g.cache.Set(ctx, cacheKey, user, int(g.ttl.Seconds()))
	return user, nil
}

// ==============================
// === GetStudentInfo ===
// ==============================
func (g *CachedUserGateway) GetStudentInfo(ctx context.Context, studentID string) (*gw_response.StudentResponse, error) {
	cacheKey := fmt.Sprintf("student:%s", studentID)

	var cached gw_response.StudentResponse
	if err := g.cache.Get(ctx, cacheKey, &cached); err == nil && cached.ID != "" {
		return &cached, nil
	}

	student, err := g.inner.GetStudentInfo(ctx, studentID)
	if err != nil {
		return nil, err
	}

	_ = g.cache.Set(ctx, cacheKey, student, int(g.ttl.Seconds()))
	return student, nil
}

// ==============================
// === GetTeachersByUser ===
// ==============================
func (g *CachedUserGateway) GetTeachersByUser(ctx context.Context, userID string) ([]*gw_response.TeacherResponse, error) {
	cacheKey := fmt.Sprintf("user:%s:teachers", userID)

	var cached []*gw_response.TeacherResponse
	if err := g.cache.Get(ctx, cacheKey, &cached); err == nil && len(cached) > 0 {
		return cached, nil
	}

	teachers, err := g.inner.GetTeachersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	_ = g.cache.Set(ctx, cacheKey, teachers, int(g.ttl.Seconds()))
	return teachers, nil
}

// ==============================
// === GetTeacherByUserAndOrg ===
// ==============================
func (g *CachedUserGateway) GetTeacherByUserAndOrganization(ctx context.Context, userID, orgID string) (*gw_response.TeacherResponse, error) {
	cacheKey := fmt.Sprintf("teacher:user:%s:org:%s", userID, orgID)

	var cached gw_response.TeacherResponse
	if err := g.cache.Get(ctx, cacheKey, &cached); err == nil && cached.ID != "" {
		return &cached, nil
	}

	teacher, err := g.inner.GetTeacherByUserAndOrganization(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	_ = g.cache.Set(ctx, cacheKey, teacher, int(g.ttl.Seconds()))
	return teacher, nil
}

// ==============================
// === Invalidate cache ===
// ==============================
func (g *CachedUserGateway) InvalidateUserCache(ctx context.Context, userID string) error {
	patterns := []string{
		fmt.Sprintf("user:info:%s", userID),
		fmt.Sprintf("user:%s:teachers", userID),
	}
	for _, key := range patterns {
		_ = g.cache.Delete(ctx, key)
	}
	return nil
}

// === Các hàm còn lại không cache ===

func (g *CachedUserGateway) GetCurrentUser(ctx context.Context) (*gw_response.CurrentUser, error) {
	return g.inner.GetCurrentUser(ctx)
}

func (g *CachedUserGateway) GetUserByTeacher(ctx context.Context, teacherID string) (*gw_response.CurrentUser, error) {
	return g.inner.GetUserByTeacher(ctx, teacherID)
}

func (g *CachedUserGateway) GetTeacherInfo(ctx context.Context, userID string, organizationID string) (*gw_response.TeacherResponse, error) {
	return g.inner.GetTeacherInfo(ctx, userID, organizationID)
}
