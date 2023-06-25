package dao

import (
	"github.com/go-redis/redis"
	"go-im/common/conf/dbConf"
	"go-im/common/model"
	"gorm.io/gorm"
	"sync"
)

type UserDao struct {
	db    *gorm.DB
	cache *redis.Client
}

var userO sync.Once
var userDao *UserDao

func initUserDao() {
	userDao = &UserDao{db: dbConf.NewDBClient(), cache: dbConf.NewRedisClient()}
}

func NewUserDao() *UserDao {
	userO.Do(initUserDao)
	return userDao
}

func (dao *UserDao) SaveUser(user *model.User) error {
	return dao.db.Create(user).Error
}

func (dao *UserDao) UserIsExist(account string) (bool, error) {
	var count int64
	user := &model.User{}
	err := dao.db.Model(user).Where("account = ?", account).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 1, err
}

func (dao *UserDao) GetUserByAccount(account string) (*model.User, error) {
	user := &model.User{}
	err := dao.db.Where("account = ?", account).Find(user).Error
	return user, err
}

func (dao *UserDao) SaveLoginStatus(account, token string, nickName string) error {
	err := dao.cache.HMSet(UserLoginInfo+account,
		map[string]interface{}{"token": token, "nickName": nickName}).Err()
	if err != nil {
		return err
	}
	return dao.cache.Expire(UserLoginInfo+account, UserLoginInfoTTL).Err()
}

func (dao *UserDao) GetLoginStatus(account string) (map[string]string, error) {
	return dao.cache.HGetAll(UserLoginInfo + account).Result()
}
