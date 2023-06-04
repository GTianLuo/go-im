package dao

import (
	"github.com/go-redis/redis"
	"go-im/common/conf/dbConf"
	"strconv"
)

var gatewayStatusDao *GatewayStatusDao

type GatewayStatusDao struct {
	cache *redis.Client
}

func initGatewayStatusDao() {
	gatewayStatusDao = &GatewayStatusDao{
		cache: dbConf.NewRedisClient(),
	}
}

func NewGatewayStatus() *GatewayStatusDao {
	o.Do(initGatewayStatusDao)
	return gatewayStatusDao
}

// SaveConnStatus redis中保存gateway(user)的长连接状态
func (dao *GatewayStatusDao) SaveConnStatus(deviceId int, account string) error {
	return dao.cache.SAdd(GateWayConnsStatus+strconv.Itoa(deviceId), account).Err()
}

// DelConnStatus  redis中删除gateway(user)的长连接状态
func (dao *GatewayStatusDao) DelConnStatus(deviceId int, account string) error {
	return dao.cache.SRem(GateWayConnsStatus+strconv.Itoa(deviceId), account).Err()
}

// GetGlobalMessageId 获取全局的message id
func (dao *GatewayStatusDao) GetGlobalMessageId() (int64, error) {
	return dao.cache.Incr(MessageGlobalId).Result()
}
