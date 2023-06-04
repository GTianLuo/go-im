package dao

import (
	"sync"
	"time"
)

var o sync.Once

const (
	UserLoginInfo    = "imSys:user:token:"
	UserLoginInfoTTL = time.Hour * 24 * 7 //一周未使用

	GateWayConnsStatus = "imSys:gateway:deviceId:"

	MessageGlobalId = "imSys:message:id:"
)
