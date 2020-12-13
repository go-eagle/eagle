// 重复提交模型封装
package pkg

import (
	"fmt"
	"time"

	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/redis"
	"github.com/1024casts/snake/pkg/utils"
)

type CRepeat struct {
	cRepeatClient redis.CheckRepeat
}

func NewCRepeat() *CRepeat {
	return &CRepeat{
		cRepeatClient: redis.NewCheckRepeat(redis.RedisClient),
	}
}

func (c *CRepeat) getKey(userID int64, check string) string {
	key, err := utils.Md5(fmt.Sprintf("%d:%s", userID, check))
	if err != nil {
		log.Warnf("md5 string err: %v", err)
	}
	return key
}

func (c *CRepeat) Set(userID int64, check string, value interface{}, expiration time.Duration) error {
	return c.cRepeatClient.Set(c.getKey(userID, check), value, expiration)
}

func (c *CRepeat) SetNX(userID int64, check string, value interface{}, expiration time.Duration) (bool, error) {
	return c.cRepeatClient.SetNX(c.getKey(userID, check), value, expiration)
}

func (c *CRepeat) Get(userID int64, check string) (interface{}, error) {
	return c.cRepeatClient.Get(c.getKey(userID, check))
}

func (c *CRepeat) Del(userID int64, check string) int64 {
	return c.cRepeatClient.Del(c.getKey(userID, check))
}
