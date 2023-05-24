package dao

import (
	"github.com/go-redis/redis"
	"go-im/common/conf/dbConf"
	"go-im/common/model"
	"gorm.io/gorm"
)

type UserDao struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewUserDao() *UserDao {
	return &UserDao{db: dbConf.NewDBClient(), cache: dbConf.NewRedisClient()}
}

func (dao *UserDao) SaveUser(user *model.User) error {
	return dao.db.Create(user).Error
}

func (dao *UserDao) IsExist(account string) (bool, error) {
	var count int64
	err := dao.db.Model(&model.User{}).Select("account = ?", account).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 1, err
}

func (dao *UserDao) SaveLoginStatus(account, token string, nickName string) error {
	return dao.cache.HMSet(UserLoginInfo+account, map[string]interface{}{token: token, account: account, nickName: nickName}).Err()
}

func (dao *UserDao) GetLoginStatus(account string) (map[string]string, error) {
	return dao.cache.HGetAll(UserLoginInfo + account).Result()
}
