package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/storage/sql"
)

var (
	_getUserStatInfo  = "select id,user_id,follow_count,follower_count,status from `%s` where user_id=?;"
	_getUserStatInfos = "select id,user_id,follow_count,follower_count,status from `%s` where user_id in (%s);"
)

func getUserStatTableName() string {
	return "user_stat"
}

// IncrFollowCount 增加关注数
func (d *repository) IncrFollowCount(ctx context.Context, db *gorm.DB, userID uint64, step int) error {
	err := db.Exec("insert into user_stat set user_id=?, follow_count=1, created_at=? on duplicate key update "+
		"follow_count=follow_count+?, updated_at=?",
		userID, time.Now(), step, time.Now()).Error
	if err != nil {
		return errors.Wrap(err, "[user_stat_repo] incr user follow count")
	}
	return nil
}

// IncrFollowerCount 增加粉丝数
func (d *repository) IncrFollowerCount(ctx context.Context, db *gorm.DB, userID uint64, step int) error {
	err := db.Exec("insert into user_stat set user_id=?, follower_count=1, created_at=? on duplicate key update "+
		"follower_count=follower_count+?, updated_at=?",
		userID, time.Now(), step, time.Now()).Error
	if err != nil {
		return errors.Wrap(err, "[user_stat_repo] incr user follower count")
	}
	return nil
}

// GetUserStatByID 获取用户统计数据
func (d *repository) GetUserStatByID(ctx context.Context, userID uint64) (res *model.UserStatModel, err error) {
	res = &model.UserStatModel{}
	_sql := fmt.Sprintf(_getUserStatInfo, getUserStatTableName())
	row := d.db.QueryRow(ctx, _sql, userID)
	err = row.Scan(&res.ID, &res.UserID, &res.FollowCount, &res.FollowerCount, &res.Status)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("[dao.GetUserStatByID] row scan err, sql: %s, err: %v", _sql, err)
		return nil, errors.Wrap(err, "[dao.user_stat] get user stat err")
	}
	return
}

// GetUserStatByIDs 批量获取用户统计数据
func (d *repository) GetUserStatByIDs(ctx context.Context, userID []uint64) (map[uint64]*model.UserStatModel, error) {
	if len(userID) == 0 {
		return nil, nil
	}
	userStats := make([]*model.UserStatModel, 0)
	res := make(map[uint64]*model.UserStatModel)

	var userIDsStr []string
	for _, v := range userID {
		userIDsStr = append(userIDsStr, cast.ToString(v))
	}

	_sql := fmt.Sprintf(_getUserStatInfos, getUserStatTableName(), strings.Join(userIDsStr, ","))
	rows, err := d.db.Query(ctx, _sql)
	if err != nil {
		log.Errorf("d.orm.Query(%v), err: %v", _sql, err)
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		r := &model.UserStatModel{}
		if err = rows.Scan(&r.ID, &r.UserID, &r.FollowCount, &r.FollowerCount, &r.Status); err != nil {
			log.Errorf("rows.Load() err: %v", err)
			continue
		}
		if r.ID != 0 {
			userStats = append(userStats, r)
		}
	}

	for _, v := range userStats {
		res[v.UserID] = v
	}

	return res, nil
}
