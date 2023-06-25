package dao

import (
	"github.com/go-redis/redis"
	"go-im/common/conf/dbConf"
	cache "go-im/common/dao/pb"
	"gorm.io/gorm"
	"sync"
)

type MessageDao struct {
	cache *redis.Client
	db    *gorm.DB
}

var messageDao *MessageDao

var messageO sync.Once

func initMessage() {
	messageDao = &MessageDao{
		cache: dbConf.NewRedisClient(),
		db:    dbConf.NewDBClient(),
	}
}

func NewMessageDao() *MessageDao {
	messageO.Do(initMessage)
	return messageDao
}

func (dao *MessageDao) SavePrivateTextMsg(msg *cache.PrivateMsg) error {
	return dao.cache.SAdd(PrivateMessageInbox+msg.To, msg).Err()
}
