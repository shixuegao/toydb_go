// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: proto/message.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AddrType int32

const (
	AddrType_PEERS  AddrType = 0
	AddrType_PEER   AddrType = 1
	AddrType_LOCAL  AddrType = 2
	AddrType_CLIENT AddrType = 3
)

// Enum value maps for AddrType.
var (
	AddrType_name = map[int32]string{
		0: "PEERS",
		1: "PEER",
		2: "LOCAL",
		3: "CLIENT",
	}
	AddrType_value = map[string]int32{
		"PEERS":  0,
		"PEER":   1,
		"LOCAL":  2,
		"CLIENT": 3,
	}
)

func (x AddrType) Enum() *AddrType {
	p := new(AddrType)
	*p = x
	return p
}

func (x AddrType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AddrType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_message_proto_enumTypes[0].Descriptor()
}

func (AddrType) Type() protoreflect.EnumType {
	return &file_proto_message_proto_enumTypes[0]
}

func (x AddrType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AddrType.Descriptor instead.
func (AddrType) EnumDescriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{0}
}

type Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddrType AddrType `protobuf:"varint,1,opt,name=addrType,proto3,enum=proto.AddrType" json:"addrType,omitempty"`
	Peer     string   `protobuf:"bytes,2,opt,name=peer,proto3" json:"peer,omitempty"`
}

func (x *Address) Reset() {
	*x = Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{0}
}

func (x *Address) GetAddrType() AddrType {
	if x != nil {
		return x.AddrType
	}
	return AddrType_PEERS
}

func (x *Address) GetPeer() string {
	if x != nil {
		return x.Peer
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType EventType  `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
	Term      uint64     `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	From      *Address   `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	To        *Address   `protobuf:"bytes,4,opt,name=to,proto3" json:"to,omitempty"`
	Event     *anypb.Any `protobuf:"bytes,5,opt,name=event,proto3" json:"event,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

func (x *Message) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *Message) GetFrom() *Address {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *Message) GetTo() *Address {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *Message) GetEvent() *anypb.Any {
	if x != nil {
		return x.Event
	}
	return nil
}

var File_proto_message_proto protoreflect.FileDescriptor

var file_proto_message_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x07, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2b, 0x0a, 0x08, 0x61, 0x64, 0x64, 0x72, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x41, 0x64, 0x64, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x61, 0x64, 0x64, 0x72, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x65, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x70, 0x65, 0x65, 0x72, 0x22, 0xbe, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x2f, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x22, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x1e, 0x0a, 0x02, 0x74,
	0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x2a, 0x0a, 0x05, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79,
	0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2a, 0x36, 0x0a, 0x08, 0x41, 0x64, 0x64, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x45, 0x45, 0x52, 0x53, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x50, 0x45, 0x45, 0x52, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x4f, 0x43, 0x41,
	0x4c, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x4c, 0x49, 0x45, 0x4e, 0x54, 0x10, 0x03, 0x42,
	0x08, 0x5a, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_message_proto_rawDescOnce sync.Once
	file_proto_message_proto_rawDescData = file_proto_message_proto_rawDesc
)

func file_proto_message_proto_rawDescGZIP() []byte {
	file_proto_message_proto_rawDescOnce.Do(func() {
		file_proto_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_message_proto_rawDescData)
	})
	return file_proto_message_proto_rawDescData
}

var file_proto_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_message_proto_goTypes = []interface{}{
	(AddrType)(0),     // 0: proto.AddrType
	(*Address)(nil),   // 1: proto.Address
	(*Message)(nil),   // 2: proto.Message
	(EventType)(0),    // 3: proto.EventType
	(*anypb.Any)(nil), // 4: google.protobuf.Any
}
var file_proto_message_proto_depIdxs = []int32{
	0, // 0: proto.Address.addrType:type_name -> proto.AddrType
	3, // 1: proto.Message.event_type:type_name -> proto.EventType
	1, // 2: proto.Message.from:type_name -> proto.Address
	1, // 3: proto.Message.to:type_name -> proto.Address
	4, // 4: proto.Message.event:type_name -> google.protobuf.Any
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_message_proto_init() }
func file_proto_message_proto_init() {
	if File_proto_message_proto != nil {
		return
	}
	file_proto_event_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Address); i {
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
		file_proto_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			RawDescriptor: file_proto_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_message_proto_goTypes,
		DependencyIndexes: file_proto_message_proto_depIdxs,
		EnumInfos:         file_proto_message_proto_enumTypes,
		MessageInfos:      file_proto_message_proto_msgTypes,
	}.Build()
	File_proto_message_proto = out.File
	file_proto_message_proto_rawDesc = nil
	file_proto_message_proto_goTypes = nil
	file_proto_message_proto_depIdxs = nil
}
