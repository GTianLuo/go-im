package dao

import (
	"errors"
	"github.com/go-redis/redis"
	"go-im/common/conf/dbConf"
	"go-im/common/log"
	"strconv"
	"sync"
)

var gatewayStatusDao *GatewayStatusDao

var gatewayO sync.Once

type GatewayStatusDao struct {
	cache *redis.Client
}

func initGatewayStatusDao() {
	gatewayStatusDao = &GatewayStatusDao{
		cache: dbConf.NewRedisClient(),
	}
}

func NewGatewayStatus() *GatewayStatusDao {
	gatewayO.Do(initGatewayStatusDao)
	return gatewayStatusDao
}

// SaveConnStatus redis中保存gateway(user)的长连接状态
func (dao *GatewayStatusDao) SaveConnStatus(deviceId string, account string) error {
	return dao.cache.Set(GateWayConnsStatus+account, deviceId, -1).Err()
}

// DelConnStatus  redis中删除gateway(user)的长连接状态
func (dao *GatewayStatusDao) DelConnStatus(account string) error {
	return dao.cache.Del(GateWayConnsStatus + account).Err()
}

// GetGlobalMessageId 获取全局的message id
func (dao *GatewayStatusDao) GetGlobalMessageId() (string, error) {
	id, err := dao.cache.Incr(MessageGlobalId).Result()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

// UserIsOnline 判断用户是否在线，若在线返回所在设备Id
func (dao *GatewayStatusDao) UserIsOnline(account string) (bool, string, error) {
	log.Info(account)
	result, err := dao.cache.Exists(GateWayConnsStatus + account).Result()
	if err != nil {
		return false, "", err
	}
	if result == 0 {
		// 用户不在线
		return false, "", nil
	}

	cmd := dao.cache.Get(GateWayConnsStatus + account)
	if err := cmd.Err(); err != nil {
		return false, "", errors.New("dao: UserIsOnline:" + err.Error())
	}
	rs := cmd.String()
	return true, rs, nil
}
