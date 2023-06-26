package messaging

import (
	"errors"
	"github.com/panjf2000/ants/v2"
	"go-im/common/conf/middlewareConf"
	"go-im/common/dao"
	"go-im/common/log"
	"go-im/common/message"
	"go-im/common/model"
	"google.golang.org/protobuf/proto"

	"go-im/common/conf/serviceConf"
	"go-im/common/mq"
)

type MsgHandler struct {
	mqConsumer *mq.MQConsumer // 消费mq消息
	mqWorker   *mq.MQWorker
	//msgS       *msgStorage.MsgStorage
	workPool *ants.Pool
	done     chan struct{}
}

func initMsgHandle() error {
	msgHandler := &MsgHandler{
		done: make(chan struct{}, 1),
	}
	pool, err := ants.NewPool(1000)
	if err != nil {
		return err
	}
	msgHandler.workPool = pool
	ch, err := middlewareConf.GetMqChannel()
	if err != nil {
		return err
	}
	consumer, err := mq.NewConsumer(ch, serviceConf.GetMessagingMqQueueName(), msgHandler.workPool)
	if err != nil {
		return errors.New("failed init Manager: " + err.Error())
	}
	msgHandler.mqConsumer = consumer

	mqWorker, err := mq.NewWorker(ch, serviceConf.GetMessagingMqXName(), "")
	if err != nil {
		return err
	}
	msgHandler.mqWorker = mqWorker
	go msgHandler.handleMsg()
	return nil
}

func (msgH *MsgHandler) handleMsg() {
	msgH.mqConsumer.ConsumeMsg(msgH.done, func(msg []byte) bool {
		cmd := &message.Cmd{}
		if err := proto.Unmarshal(msg, cmd); err != nil {
			log.Error(err)
			return false
		}
		switch cmd.Type {
		case message.CmdType_PrivateMsgCmd:
			if err := msgH.HandlePrivateMsg(msg, cmd); err != nil {
				log.Error(err)
				return false
			}
		case message.CmdType_GroupMsgCmd:
		}
		return true
	})
}

func (msgH *MsgHandler) HandlePrivateMsg(cmdBytes []byte, cmd *message.Cmd) error {

	//msgDao := dao.NewMessageDao()
	msg := &message.ChatMsg{}
	if err := proto.Unmarshal(cmd.Payload, msg); err != nil {
		return errors.New("HandlePrivateMsg: failed to unmarshal: " + err.Error())
	}
	//判断用户是否在线
	gatewayStatusDao := dao.NewGatewayStatus()
	isOnline, deviceId, err := gatewayStatusDao.UserIsOnline(msg.To)
	if err != nil {
		return errors.New("HandlePrivateMsg:" + err.Error())
	}
	if isOnline {
		// TODO
		log.Error("======+====")
		// 用户在线，转发消息
		if err := msgH.mqWorker.PublishMsg(cmdBytes, serviceConf.GetMessagingMqRoutingKey(deviceId)); err != nil {
			// 消息push失败
			return errors.New("HandlePrivateMsg: failed to push msg to mq :" + err.Error())
		}
		if err := msgH.storagePrivateOnlineMsg(cmd, msg); err != nil {
			return errors.New("HandlePrivateMsg:" + err.Error())
		}
	}
	//TODO 处理离线消息
	return nil
}

func (msgH *MsgHandler) storagePrivateOnlineMsg(cmd *message.Cmd, msg *message.ChatMsg) error {

	messageDao := dao.NewMessageDao()
	// 持久化存储发送消息记录
	if err := messageDao.SaveMsgSend(&model.MsgSend{
		MsgFrom:    cmd.From,
		MsgTo:      msg.To,
		MsgSeq:     cmd.MsgId,
		MsgContent: msg.Data,
		SendTime:   cmd.Timestamp,
		CmdType:    cmd.Type,
		MsgType:    msg.Type,
	}); err != nil {
		return err
	}
	// 持久化存储消息接收表
	if err := messageDao.SaveMsgRecv(&model.MsgReceive{
		MsgFrom: cmd.From,
		MsgTo:   msg.To,
		MsgSeq:  cmd.MsgId,
	}); err != nil {
		return err
	}
	return nil
}
