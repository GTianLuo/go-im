package messaging

import (
	"errors"
	"go-im/common/conf/middlewareConf"
	"go-im/common/dao"
	cache "go-im/common/dao/pb"
	"go-im/common/log"
	"go-im/common/message"
	"google.golang.org/protobuf/proto"

	"github.com/panjf2000/ants/v2"

	"go-im/common/conf/serviceConf"
	"go-im/common/mq"
)

type MsgHandler struct {
	Consumer *mq.MQConsumer // 消费mq消息
	workPool *ants.Pool     //协程池
	done     chan struct{}
}

func initMsgHandle() error {
	msgHandler := &MsgHandler{
		done: make(chan struct{}, 1),
	}

	ch, err := middlewareConf.GetMqChannel()
	if err != nil {
		return err
	}
	consumer, err := mq.NewConsumer(ch, serviceConf.GetGatewayMqQueueName())
	if err != nil {
		return errors.New("failed init Manager: " + err.Error())
	}
	msgHandler.Consumer = consumer
	if msgHandler.workPool, err = ants.NewPool(serviceConf.GetGateWayWorkPoolNum()); err != nil {
		return errors.New("failed init Manager: " + err.Error())
	}
	go msgHandler.handleMsg()
	return nil
}

func (msgH *MsgHandler) handleMsg() {
	msgH.Consumer.ConsumeMsg(msgH.done, func(msg []byte) bool {
		cmd := &message.Cmd{}
		if err := proto.Unmarshal(msg, cmd); err != nil {
			log.Error(err)
			return false
		}
		switch cmd.Type {
		case message.CmdType_PrivateMsgCmd:
			if err := msgH.HandlePrivateMsg(cmd); err != nil {
				log.Error(err)
				return false
			}
		case message.CmdType_GroupMsgCmd:
		}
		return true
	})
}

func (msgH *MsgHandler) HandlePrivateMsg(cmd *message.Cmd) error {

	msgDao := dao.NewMessageDao()
	msg := &message.PrivateMsg{}
	if err := proto.Unmarshal(cmd.Payload, msg); err != nil {
		return err
	}
	switch msg.Type {
	case message.MsgType_TextMsg:
		err := msgDao.SavePrivateTextMsg(&cache.PrivateMsg{
			Timestamp: cmd.Timestamp,
			MsgId:     cmd.MsgId,
			From:      cmd.From,
			To:        msg.To,
			Content:   string(msg.Data),
		})
		return err
	default:
		return errors.New("invalid message.pb type")
	}
}
