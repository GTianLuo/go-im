// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: message/message.proto

package message

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CmdType int32

const (
	CmdType_PrivateMsgCmd   CmdType = 0
	CmdType_GroupMsgCmd     CmdType = 1
	CmdType_HeartBeatCmd    CmdType = 2
	CmdType_AuthRequestCmd  CmdType = 3
	CmdType_MsgAckCmd       CmdType = 4
	CmdType_HeartBeatAckCmd CmdType = 5
	CmdType_AuthResponseCmd CmdType = 6
	CmdType_SystemCmd       CmdType = 7
)

// Enum value maps for CmdType.
var (
	CmdType_name = map[int32]string{
		0: "PrivateMsgCmd",
		1: "GroupMsgCmd",
		2: "HeartBeatCmd",
		3: "AuthRequestCmd",
		4: "MsgAckCmd",
		5: "HeartBeatAckCmd",
		6: "AuthResponseCmd",
		7: "SystemCmd",
	}
	CmdType_value = map[string]int32{
		"PrivateMsgCmd":   0,
		"GroupMsgCmd":     1,
		"HeartBeatCmd":    2,
		"AuthRequestCmd":  3,
		"MsgAckCmd":       4,
		"HeartBeatAckCmd": 5,
		"AuthResponseCmd": 6,
		"SystemCmd":       7,
	}
)

func (x CmdType) Enum() *CmdType {
	p := new(CmdType)
	*p = x
	return p
}

func (x CmdType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CmdType) Descriptor() protoreflect.EnumDescriptor {
	return file_message_message_proto_enumTypes[0].Descriptor()
}

func (CmdType) Type() protoreflect.EnumType {
	return &file_message_message_proto_enumTypes[0]
}

func (x CmdType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CmdType.Descriptor instead.
func (CmdType) EnumDescriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{0}
}

type MsgType int32

const (
	MsgType_TextMsg    MsgType = 0
	MsgType_PictureMsg MsgType = 1
)

// Enum value maps for MsgType.
var (
	MsgType_name = map[int32]string{
		0: "TextMsg",
		1: "PictureMsg",
	}
	MsgType_value = map[string]int32{
		"TextMsg":    0,
		"PictureMsg": 1,
	}
)

func (x MsgType) Enum() *MsgType {
	p := new(MsgType)
	*p = x
	return p
}

func (x MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_message_message_proto_enumTypes[1].Descriptor()
}

func (MsgType) Type() protoreflect.EnumType {
	return &file_message_message_proto_enumTypes[1]
}

func (x MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgType.Descriptor instead.
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{1}
}

// 顶层cmd pb结构
type Cmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      CmdType `protobuf:"varint,1,opt,name=Type,proto3,enum=message.CmdType" json:"Type,omitempty"`
	MsgId     string  `protobuf:"bytes,2,opt,name=MsgId,proto3" json:"MsgId,omitempty"` //该字段在上下游有不同的含义，作上游消息的时候MsgId是会话内的唯一Id，作下游消息的时候是全局唯一id
	From      string  `protobuf:"bytes,3,opt,name=From,proto3" json:"From,omitempty"`
	Timestamp int64   `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Payload   []byte  `protobuf:"bytes,5,opt,name=Payload,proto3" json:"Payload,omitempty"`
}

func (x *Cmd) Reset() {
	*x = Cmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cmd) ProtoMessage() {}

func (x *Cmd) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cmd.ProtoReflect.Descriptor instead.
func (*Cmd) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{0}
}

func (x *Cmd) GetType() CmdType {
	if x != nil {
		return x.Type
	}
	return CmdType_PrivateMsgCmd
}

func (x *Cmd) GetMsgId() string {
	if x != nil {
		return x.MsgId
	}
	return ""
}

func (x *Cmd) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Cmd) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Cmd) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

// 聊天
type ChatMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type MsgType `protobuf:"varint,1,opt,name=Type,proto3,enum=message.MsgType" json:"Type,omitempty"`
	To   string  `protobuf:"bytes,3,opt,name=To,proto3" json:"To,omitempty"`
	Data []byte  `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *ChatMsg) Reset() {
	*x = ChatMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMsg) ProtoMessage() {}

func (x *ChatMsg) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMsg.ProtoReflect.Descriptor instead.
func (*ChatMsg) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{1}
}

func (x *ChatMsg) GetType() MsgType {
	if x != nil {
		return x.Type
	}
	return MsgType_TextMsg
}

func (x *ChatMsg) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *ChatMsg) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

//鉴权请求
type AuthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"`
}

func (x *AuthRequest) Reset() {
	*x = AuthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthRequest) ProtoMessage() {}

func (x *AuthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthRequest.ProtoReflect.Descriptor instead.
func (*AuthRequest) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{2}
}

func (x *AuthRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

//鉴权响应
type AuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status int32  `protobuf:"varint,1,opt,name=Status,proto3" json:"Status,omitempty"` //1表示内部错误，0表示鉴权失败，1表示成功,2表示该网关连接达到上限，需重新连接
	ErrMsg string `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
}

func (x *AuthResponse) Reset() {
	*x = AuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthResponse) ProtoMessage() {}

func (x *AuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthResponse.ProtoReflect.Descriptor instead.
func (*AuthResponse) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{3}
}

func (x *AuthResponse) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *AuthResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

var File_message_message_proto protoreflect.FileDescriptor

var file_message_message_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x8d, 0x01, 0x0a, 0x03, 0x43, 0x6d, 0x64, 0x12, 0x24, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x43, 0x6d, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x4d, 0x73, 0x67, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4d,
	0x73, 0x67, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x22, 0x53, 0x0a, 0x07, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x24, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x54, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x54,
	0x6f, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x22, 0x23, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x3e, 0x0a, 0x0c, 0x41, 0x75,
	0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x2a, 0x9b, 0x01, 0x0a, 0x07, 0x43,
	0x6d, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x4d, 0x73, 0x67, 0x43, 0x6d, 0x64, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x4d, 0x73, 0x67, 0x43, 0x6d, 0x64, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x43, 0x6d, 0x64, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e,
	0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x6d, 0x64, 0x10, 0x03,
	0x12, 0x0d, 0x0a, 0x09, 0x4d, 0x73, 0x67, 0x41, 0x63, 0x6b, 0x43, 0x6d, 0x64, 0x10, 0x04, 0x12,
	0x13, 0x0a, 0x0f, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x41, 0x63, 0x6b, 0x43,
	0x6d, 0x64, 0x10, 0x05, 0x12, 0x13, 0x0a, 0x0f, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6d, 0x64, 0x10, 0x06, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x43, 0x6d, 0x64, 0x10, 0x07, 0x2a, 0x26, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x54, 0x65, 0x78, 0x74, 0x4d, 0x73, 0x67, 0x10, 0x00,
	0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65, 0x4d, 0x73, 0x67, 0x10, 0x01,
	0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x3b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_message_proto_rawDescOnce sync.Once
	file_message_message_proto_rawDescData = file_message_message_proto_rawDesc
)

func file_message_message_proto_rawDescGZIP() []byte {
	file_message_message_proto_rawDescOnce.Do(func() {
		file_message_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_message_proto_rawDescData)
	})
	return file_message_message_proto_rawDescData
}

var file_message_message_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_message_message_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_message_message_proto_goTypes = []interface{}{
	(CmdType)(0),         // 0: message.CmdType
	(MsgType)(0),         // 1: message.MsgType
	(*Cmd)(nil),          // 2: message.Cmd
	(*ChatMsg)(nil),      // 3: message.ChatMsg
	(*AuthRequest)(nil),  // 4: message.AuthRequest
	(*AuthResponse)(nil), // 5: message.AuthResponse
}
var file_message_message_proto_depIdxs = []int32{
	0, // 0: message.Cmd.Type:type_name -> message.CmdType
	1, // 1: message.ChatMsg.Type:type_name -> message.MsgType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_message_message_proto_init() }
func file_message_message_proto_init() {
	if File_message_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cmd); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_message_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_message_proto_goTypes,
		DependencyIndexes: file_message_message_proto_depIdxs,
		EnumInfos:         file_message_message_proto_enumTypes,
		MessageInfos:      file_message_message_proto_msgTypes,
	}.Build()
	File_message_message_proto = out.File
	file_message_message_proto_rawDesc = nil
	file_message_message_proto_goTypes = nil
	file_message_message_proto_depIdxs = nil
}
