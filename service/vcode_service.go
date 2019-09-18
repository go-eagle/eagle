package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/1024casts/snake/pkg/redis"
	"github.com/lexkong/log"
	"github.com/pkg/errors"
)

// 验证码服务，主要提供生成验证码和获取验证码

// 直接初始化，可以避免在使用时再实例化
var VCodeService = NewVCodeService()

const (
	verifyCodeRedisKey = "app:login:vcode:%s" // 验证码key
	maxDurationTime    = 10 * time.Minute     // 验证码有效期
)

// 校验码服务，生成校验码和获得校验码
type vcodeService struct {
}

func NewVCodeService() *vcodeService {
	return &vcodeService{}
}

// 生成校验码
func (srv *vcodeService) GenLoginVCode(phone string) (string, error) {
	// step1: 生成随机数
	vcode := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	// step2: 写入到redis里
	// 使用set, key使用前缀+手机号 缓存10分钟）
	key := fmt.Sprintf(verifyCodeRedisKey, phone)
	err := redis.Client.Set(key, vcode, maxDurationTime).Err()
	if err != nil {
		log.Warnf("[vcode_service] redis set err, %v", err)
		return "", errors.New("redis set err")
	}

	return vcode, nil
}

// 生成校验码
func (srv *vcodeService) GetLoginVCode(phone string) (string, error) {
	// 直接从redis里获取
	key := fmt.Sprintf(verifyCodeRedisKey, phone)
	vcode, err := redis.Client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		log.Warnf("[vcode_service] redis get err, %v", err)
		return "", errors.New("redis get err")
	}

	return vcode, nil
}
