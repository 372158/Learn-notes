package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"shortlink/internal/model"
	"shortlink/pkg/cache"
	"shortlink/pkg/utils"
)

type ShortLinkService struct {
	DB *gorm.DB
}

func NewShortLinkService(db *gorm.DB) *ShortLinkService {
	return &ShortLinkService{DB: db}
}

// CreateShortLink 创建短链接
func (s *ShortLinkService) CreateShortLink(originalURL string, expiredAt *time.Time) (*model.ShortLink, error) {
	// 1. 检查是否已存在该 URL 的短链接（可选，为了节省空间）
	var existing model.ShortLink
	if err := s.DB.Where("original_url = ?", originalURL).First(&existing).Error; err == nil {
		return &existing, nil // 已存在则直接返回
	}

	// 2. 生成短码（带碰撞重试机制）
	var shortCode string
	for i := 0; i < 5; i++ { // 最多尝试 5 次
		// 用 URL 加上随机种子生成哈希，降低碰撞概率
		seed := ""
		if i > 0 {
			seed = fmt.Sprintf("%d", rand.Int63())
		}
		shortCode = utils.GenerateShortCode(originalURL+seed, 6)

		// 检查短码是否已被占用
		var count int64
		if err := s.DB.Model(&model.ShortLink{}).Where("short_code = ?", shortCode).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("查询短码失败: %w", err)
		}
		if count == 0 {
			break
		}
		if i == 4 {
			return nil, errors.New("生成短码失败，请重试")
		}
	}

	// 3. 创建记录
	shortLink := &model.ShortLink{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		ClickCount:  0,
		ExpiredAt:   expiredAt,
	}

	if err := s.DB.Create(shortLink).Error; err != nil {
		return nil, err
	}

	return shortLink, nil
}

// GetOriginalURL 根据短码获取原始 URL（带缓存）
func (s *ShortLinkService) GetOriginalURL(shortCode string) (string, error) {
	cacheKey := "shortlink:" + shortCode

	// 1. 先查 Redis 缓存
	cached, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
		return cached, nil
	}

	// 2. 缓存未命中，查数据库
	var shortLink model.ShortLink
	if err := s.DB.Where("short_code = ?", shortCode).First(&shortLink).Error; err != nil {
		return "", errors.New("短链接不存在或已过期")
	}

	// 3. 检查是否过期
	if shortLink.ExpiredAt != nil && shortLink.ExpiredAt.Before(time.Now()) {
		return "", errors.New("短链接已过期")
	}

	// 4. 异步增加点击次数（不影响响应速度）
	go s.incrementClickCount(shortCode)

	// 5. 写入缓存（设置 1 小时过期）
	cache.Rdb.Set(cache.Ctx, cacheKey, shortLink.OriginalURL, time.Hour)

	return shortLink.OriginalURL, nil
}

// incrementClickCount 增加点击次数（异步）
func (s *ShortLinkService) incrementClickCount(shortCode string) {
	s.DB.Model(&model.ShortLink{}).Where("short_code = ?", shortCode).
		UpdateColumn("click_count", gorm.Expr("click_count + 1"))
}