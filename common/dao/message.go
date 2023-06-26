package dao

import (
	"github.com/go-redis/redis"
	"go-im/common/conf/dbConf"
	cache "go-im/common/dao/pb"
	"go-im/common/model"
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

// SaveMsgSend mysql中保存发送的消息
func (dao *MessageDao) SaveMsgSend(send *model.MsgSend) error {
	return dao.db.Model(&model.MsgSend{}).Save(send).Error
}

// SaveMsgRecv mysql中保存未读消息列表
func (dao *MessageDao) SaveMsgRecv(recv *model.MsgReceive) error {
	return dao.db.Model(&model.MsgReceive{}).Save(recv).Error
}
