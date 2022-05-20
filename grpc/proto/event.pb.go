// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: proto/event.proto

package proto

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

//事件类型
type EventType int32

const (
	EventType_HEARTBEAT         EventType = 0
	EventType_CONFIRM_LEADER    EventType = 1
	EventType_SOLICIT_VOTE      EventType = 2
	EventType_GRANT_VOTE        EventType = 3
	EventType_REPLICATE_ENTRIES EventType = 4
	EventType_ACCEPT_ENTRIES    EventType = 5
	EventType_REJECT_ENTRIES    EventType = 6
	EventType_CLIENT_REQUEST    EventType = 7
	EventType_CLIENT_RESPONSE   EventType = 8
)

// Enum value maps for EventType.
var (
	EventType_name = map[int32]string{
		0: "HEARTBEAT",
		1: "CONFIRM_LEADER",
		2: "SOLICIT_VOTE",
		3: "GRANT_VOTE",
		4: "REPLICATE_ENTRIES",
		5: "ACCEPT_ENTRIES",
		6: "REJECT_ENTRIES",
		7: "CLIENT_REQUEST",
		8: "CLIENT_RESPONSE",
	}
	EventType_value = map[string]int32{
		"HEARTBEAT":         0,
		"CONFIRM_LEADER":    1,
		"SOLICIT_VOTE":      2,
		"GRANT_VOTE":        3,
		"REPLICATE_ENTRIES": 4,
		"ACCEPT_ENTRIES":    5,
		"REJECT_ENTRIES":    6,
		"CLIENT_REQUEST":    7,
		"CLIENT_RESPONSE":   8,
	}
)

func (x EventType) Enum() *EventType {
	p := new(EventType)
	*p = x
	return p
}

func (x EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_event_proto_enumTypes[0].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_proto_event_proto_enumTypes[0]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{0}
}

//请求类型
type ReqType int32

const (
	ReqType_QUERY  ReqType = 0
	ReqType_MUTATE ReqType = 1
	ReqType_STATUS ReqType = 2
)

// Enum value maps for ReqType.
var (
	ReqType_name = map[int32]string{
		0: "QUERY",
		1: "MUTATE",
		2: "STATUS",
	}
	ReqType_value = map[string]int32{
		"QUERY":  0,
		"MUTATE": 1,
		"STATUS": 2,
	}
)

func (x ReqType) Enum() *ReqType {
	p := new(ReqType)
	*p = x
	return p
}

func (x ReqType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ReqType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_event_proto_enumTypes[1].Descriptor()
}

func (ReqType) Type() protoreflect.EnumType {
	return &file_proto_event_proto_enumTypes[1]
}

func (x ReqType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ReqType.Descriptor instead.
func (ReqType) EnumDescriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{1}
}

type HeartbeatEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType   EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
	CommitIndex uint64    `protobuf:"varint,2,opt,name=commit_index,json=commitIndex,proto3" json:"commit_index,omitempty"`
	CommitTerm  uint64    `protobuf:"varint,3,opt,name=commit_term,json=commitTerm,proto3" json:"commit_term,omitempty"`
}

func (x *HeartbeatEvent) Reset() {
	*x = HeartbeatEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatEvent) ProtoMessage() {}

func (x *HeartbeatEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatEvent.ProtoReflect.Descriptor instead.
func (*HeartbeatEvent) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{0}
}

func (x *HeartbeatEvent) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

func (x *HeartbeatEvent) GetCommitIndex() uint64 {
	if x != nil {
		return x.CommitIndex
	}
	return 0
}

func (x *HeartbeatEvent) GetCommitTerm() uint64 {
	if x != nil {
		return x.CommitTerm
	}
	return 0
}

type ConfirmleaderEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType    EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
	CommitIndex  uint64    `protobuf:"varint,2,opt,name=commit_index,json=commitIndex,proto3" json:"commit_index,omitempty"`
	HasCommitted bool      `protobuf:"varint,3,opt,name=has_committed,json=hasCommitted,proto3" json:"has_committed,omitempty"`
}

func (x *ConfirmleaderEvent) Reset() {
	*x = ConfirmleaderEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfirmleaderEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfirmleaderEvent) ProtoMessage() {}

func (x *ConfirmleaderEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfirmleaderEvent.ProtoReflect.Descriptor instead.
func (*ConfirmleaderEvent) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{1}
}

func (x *ConfirmleaderEvent) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

func (x *ConfirmleaderEvent) GetCommitIndex() uint64 {
	if x != nil {
		return x.CommitIndex
	}
	return 0
}

func (x *ConfirmleaderEvent) GetHasCommitted() bool {
	if x != nil {
		return x.HasCommitted
	}
	return false
}

type SolicitVoteEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
	LastIndex uint64    `protobuf:"varint,2,opt,name=last_index,json=lastIndex,proto3" json:"last_index,omitempty"`
	LastTerm  uint64    `protobuf:"varint,3,opt,name=last_term,json=lastTerm,proto3" json:"last_term,omitempty"`
}

func (x *SolicitVoteEvent) Reset() {
	*x = SolicitVoteEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SolicitVoteEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SolicitVoteEvent) ProtoMessage() {}

func (x *SolicitVoteEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SolicitVoteEvent.ProtoReflect.Descriptor instead.
func (*SolicitVoteEvent) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{2}
}

func (x *SolicitVoteEvent) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

func (x *SolicitVoteEvent) GetLastIndex() uint64 {
	if x != nil {
		return x.LastIndex
	}
	return 0
}

func (x *SolicitVoteEvent) GetLastTerm() uint64 {
	if x != nil {
		return x.LastTerm
	}
	return 0
}

type GrantVoteEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
}

func (x *GrantVoteEvent) Reset() {
	*x = GrantVoteEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrantVoteEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrantVoteEvent) ProtoMessage() {}

func (x *GrantVoteEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrantVoteEvent.ProtoReflect.Descriptor instead.
func (*GrantVoteEvent) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{3}
}

func (x *GrantVoteEvent) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

type ReplicateEntriesEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
	BaseIndex uint64    `protobuf:"varint,2,opt,name=base_index,json=baseIndex,proto3" json:"base_index,omitempty"`
	BaseTerm  uint64    `protobuf:"varint,3,opt,name=base_term,json=baseTerm,proto3" json:"base_term,omitempty"`
	Entries   []*Entry  `protobuf:"bytes,4,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *ReplicateEntriesEvent) Reset() {
	*x = ReplicateEntriesEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplicateEntriesEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplicateEntriesEvent) ProtoMessage() {}

func (x *ReplicateEntriesEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplicateEntriesEvent.ProtoReflect.Descriptor instead.
func (*ReplicateEntriesEvent) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{4}
}

func (x *ReplicateEntriesEvent) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

func (x *ReplicateEntriesEvent) GetBaseIndex() uint64 {
	if x != nil {
		return x.BaseIndex
	}
	return 0
}

func (x *ReplicateEntriesEvent) GetBaseTerm() uint64 {
	if x != nil {
		return x.BaseTerm
	}
	return 0
}

func (x *ReplicateEntriesEvent) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type RejectEntriesEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
}

func (x *RejectEntriesEvent) Reset() {
	*x = RejectEntriesEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RejectEntriesEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RejectEntriesEvent) ProtoMessage() {}

func (x *RejectEntriesEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RejectEntriesEvent.ProtoReflect.Descriptor instead.
func (*RejectEntriesEvent) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{5}
}

func (x *RejectEntriesEvent) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

type AcceptEntriesEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
	LastIndex uint64    `protobuf:"varint,2,opt,name=last_index,json=lastIndex,proto3" json:"last_index,omitempty"`
}

func (x *AcceptEntriesEvent) Reset() {
	*x = AcceptEntriesEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptEntriesEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptEntriesEvent) ProtoMessage() {}

func (x *AcceptEntriesEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptEntriesEvent.ProtoReflect.Descriptor instead.
func (*AcceptEntriesEvent) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{6}
}

func (x *AcceptEntriesEvent) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

func (x *AcceptEntriesEvent) GetLastIndex() uint64 {
	if x != nil {
		return x.LastIndex
	}
	return 0
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Server        string            `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
	Leader        string            `protobuf:"bytes,2,opt,name=leader,proto3" json:"leader,omitempty"`
	Term          uint64            `protobuf:"varint,3,opt,name=term,proto3" json:"term,omitempty"`
	CommitIndex   uint64            `protobuf:"varint,4,opt,name=commit_index,json=commitIndex,proto3" json:"commit_index,omitempty"`
	ApplyIndex    uint64            `protobuf:"varint,5,opt,name=apply_index,json=applyIndex,proto3" json:"apply_index,omitempty"`
	Storage       string            `protobuf:"bytes,6,opt,name=storage,proto3" json:"storage,omitempty"`
	StorageSize   uint64            `protobuf:"varint,7,opt,name=storage_size,json=storageSize,proto3" json:"storage_size,omitempty"`
	NodeLastIndex map[string]uint64 `protobuf:"bytes,8,rep,name=node_last_index,json=nodeLastIndex,proto3" json:"node_last_index,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{7}
}

func (x *Status) GetServer() string {
	if x != nil {
		return x.Server
	}
	return ""
}

func (x *Status) GetLeader() string {
	if x != nil {
		return x.Leader
	}
	return ""
}

func (x *Status) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *Status) GetCommitIndex() uint64 {
	if x != nil {
		return x.CommitIndex
	}
	return 0
}

func (x *Status) GetApplyIndex() uint64 {
	if x != nil {
		return x.ApplyIndex
	}
	return 0
}

func (x *Status) GetStorage() string {
	if x != nil {
		return x.Storage
	}
	return ""
}

func (x *Status) GetStorageSize() uint64 {
	if x != nil {
		return x.StorageSize
	}
	return 0
}

func (x *Status) GetNodeLastIndex() map[string]uint64 {
	if x != nil {
		return x.NodeLastIndex
	}
	return nil
}

type ClientRequestEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType   EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
	Id          string    `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	RequestType ReqType   `protobuf:"varint,3,opt,name=request_type,json=requestType,proto3,enum=proto.ReqType" json:"request_type,omitempty"`
	Command     []byte    `protobuf:"bytes,4,opt,name=command,proto3" json:"command,omitempty"`
}

func (x *ClientRequestEvent) Reset() {
	*x = ClientRequestEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientRequestEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientRequestEvent) ProtoMessage() {}

func (x *ClientRequestEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientRequestEvent.ProtoReflect.Descriptor instead.
func (*ClientRequestEvent) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{8}
}

func (x *ClientRequestEvent) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

func (x *ClientRequestEvent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ClientRequestEvent) GetRequestType() ReqType {
	if x != nil {
		return x.RequestType
	}
	return ReqType_QUERY
}

func (x *ClientRequestEvent) GetCommand() []byte {
	if x != nil {
		return x.Command
	}
	return nil
}

type ClientResponseEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType   EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=proto.EventType" json:"event_type,omitempty"`
	Id          string    `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	RequestType ReqType   `protobuf:"varint,3,opt,name=request_type,json=requestType,proto3,enum=proto.ReqType" json:"request_type,omitempty"`
	Result      []byte    `protobuf:"bytes,4,opt,name=result,proto3" json:"result,omitempty"` //当request_type = QUERY、MUTATE的时
	Status      *Status   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"` //当request_type = STATUS时
}

func (x *ClientResponseEvent) Reset() {
	*x = ClientResponseEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientResponseEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientResponseEvent) ProtoMessage() {}

func (x *ClientResponseEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientResponseEvent.ProtoReflect.Descriptor instead.
func (*ClientResponseEvent) Descriptor() ([]byte, []int) {
	return file_proto_event_proto_rawDescGZIP(), []int{9}
}

func (x *ClientResponseEvent) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_HEARTBEAT
}

func (x *ClientResponseEvent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ClientResponseEvent) GetRequestType() ReqType {
	if x != nil {
		return x.RequestType
	}
	return ReqType_QUERY
}

func (x *ClientResponseEvent) GetResult() []byte {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *ClientResponseEvent) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

var File_proto_event_proto protoreflect.FileDescriptor

var file_proto_event_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85, 0x01,
	0x0a, 0x0e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x2f, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x74,
	0x65, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x8d, 0x01, 0x0a, 0x12, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72,
	0x6d, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x0a,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x23, 0x0a, 0x0d, 0x68, 0x61, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x74, 0x65,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x68, 0x61, 0x73, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x74, 0x65, 0x64, 0x22, 0x7f, 0x0a, 0x10, 0x53, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x74,
	0x56, 0x6f, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x0a, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09,
	0x6c, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6c, 0x61,
	0x73, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x41, 0x0a, 0x0e, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x56,
	0x6f, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x22, 0xac, 0x01, 0x0a, 0x15, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x62, 0x61, 0x73, 0x65, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x74, 0x65, 0x72, 0x6d,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x62, 0x61, 0x73, 0x65, 0x54, 0x65, 0x72, 0x6d,
	0x12, 0x26, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x45, 0x0a, 0x12, 0x52, 0x65, 0x6a, 0x65,
	0x63, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x2f,
	0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x22,
	0x64, 0x0a, 0x12, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x6c, 0x61, 0x73, 0x74,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x22, 0xd9, 0x02, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04,
	0x74, 0x65, 0x72, 0x6d, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x70, 0x70, 0x6c, 0x79,
	0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x61, 0x70,
	0x70, 0x6c, 0x79, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x48, 0x0a, 0x0f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x4c, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0d, 0x6e, 0x6f, 0x64, 0x65, 0x4c, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x1a,
	0x40, 0x0a, 0x12, 0x4e, 0x6f, 0x64, 0x65, 0x4c, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0xa2, 0x01, 0x0a, 0x12, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x31, 0x0a, 0x0c, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x71, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0xc8, 0x01, 0x0a, 0x13, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x2f,
	0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x31, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65,
	0x71, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x25, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x2a, 0xb8, 0x01, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0d, 0x0a, 0x09, 0x48, 0x45, 0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x10, 0x00, 0x12, 0x12,
	0x0a, 0x0e, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x52, 0x4d, 0x5f, 0x4c, 0x45, 0x41, 0x44, 0x45, 0x52,
	0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x4f, 0x4c, 0x49, 0x43, 0x49, 0x54, 0x5f, 0x56, 0x4f,
	0x54, 0x45, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x47, 0x52, 0x41, 0x4e, 0x54, 0x5f, 0x56, 0x4f,
	0x54, 0x45, 0x10, 0x03, 0x12, 0x15, 0x0a, 0x11, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54,
	0x45, 0x5f, 0x45, 0x4e, 0x54, 0x52, 0x49, 0x45, 0x53, 0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x41,
	0x43, 0x43, 0x45, 0x50, 0x54, 0x5f, 0x45, 0x4e, 0x54, 0x52, 0x49, 0x45, 0x53, 0x10, 0x05, 0x12,
	0x12, 0x0a, 0x0e, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x5f, 0x45, 0x4e, 0x54, 0x52, 0x49, 0x45,
	0x53, 0x10, 0x06, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x4c, 0x49, 0x45, 0x4e, 0x54, 0x5f, 0x52, 0x45,
	0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x07, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x4c, 0x49, 0x45, 0x4e,
	0x54, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0x08, 0x2a, 0x2c, 0x0a, 0x07,
	0x52, 0x65, 0x71, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x51, 0x55, 0x45, 0x52, 0x59,
	0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x55, 0x54, 0x41, 0x54, 0x45, 0x10, 0x01, 0x12, 0x0a,
	0x0a, 0x06, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x10, 0x02, 0x42, 0x08, 0x5a, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_event_proto_rawDescOnce sync.Once
	file_proto_event_proto_rawDescData = file_proto_event_proto_rawDesc
)

func file_proto_event_proto_rawDescGZIP() []byte {
	file_proto_event_proto_rawDescOnce.Do(func() {
		file_proto_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_event_proto_rawDescData)
	})
	return file_proto_event_proto_rawDescData
}

var file_proto_event_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_proto_event_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_event_proto_goTypes = []interface{}{
	(EventType)(0),                // 0: proto.EventType
	(ReqType)(0),                  // 1: proto.ReqType
	(*HeartbeatEvent)(nil),        // 2: proto.HeartbeatEvent
	(*ConfirmleaderEvent)(nil),    // 3: proto.ConfirmleaderEvent
	(*SolicitVoteEvent)(nil),      // 4: proto.SolicitVoteEvent
	(*GrantVoteEvent)(nil),        // 5: proto.GrantVoteEvent
	(*ReplicateEntriesEvent)(nil), // 6: proto.ReplicateEntriesEvent
	(*RejectEntriesEvent)(nil),    // 7: proto.RejectEntriesEvent
	(*AcceptEntriesEvent)(nil),    // 8: proto.AcceptEntriesEvent
	(*Status)(nil),                // 9: proto.Status
	(*ClientRequestEvent)(nil),    // 10: proto.ClientRequestEvent
	(*ClientResponseEvent)(nil),   // 11: proto.ClientResponseEvent
	nil,                           // 12: proto.Status.NodeLastIndexEntry
	(*Entry)(nil),                 // 13: proto.Entry
}
var file_proto_event_proto_depIdxs = []int32{
	0,  // 0: proto.HeartbeatEvent.event_type:type_name -> proto.EventType
	0,  // 1: proto.ConfirmleaderEvent.event_type:type_name -> proto.EventType
	0,  // 2: proto.SolicitVoteEvent.event_type:type_name -> proto.EventType
	0,  // 3: proto.GrantVoteEvent.event_type:type_name -> proto.EventType
	0,  // 4: proto.ReplicateEntriesEvent.event_type:type_name -> proto.EventType
	13, // 5: proto.ReplicateEntriesEvent.entries:type_name -> proto.Entry
	0,  // 6: proto.RejectEntriesEvent.event_type:type_name -> proto.EventType
	0,  // 7: proto.AcceptEntriesEvent.event_type:type_name -> proto.EventType
	12, // 8: proto.Status.node_last_index:type_name -> proto.Status.NodeLastIndexEntry
	0,  // 9: proto.ClientRequestEvent.event_type:type_name -> proto.EventType
	1,  // 10: proto.ClientRequestEvent.request_type:type_name -> proto.ReqType
	0,  // 11: proto.ClientResponseEvent.event_type:type_name -> proto.EventType
	1,  // 12: proto.ClientResponseEvent.request_type:type_name -> proto.ReqType
	9,  // 13: proto.ClientResponseEvent.status:type_name -> proto.Status
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_proto_event_proto_init() }
func file_proto_event_proto_init() {
	if File_proto_event_proto != nil {
		return
	}
	file_proto_entry_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartbeatEvent); i {
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
		file_proto_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfirmleaderEvent); i {
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
		file_proto_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SolicitVoteEvent); i {
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
		file_proto_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrantVoteEvent); i {
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
		file_proto_event_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplicateEntriesEvent); i {
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
		file_proto_event_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RejectEntriesEvent); i {
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
		file_proto_event_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptEntriesEvent); i {
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
		file_proto_event_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
		file_proto_event_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientRequestEvent); i {
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
		file_proto_event_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientResponseEvent); i {
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
			RawDescriptor: file_proto_event_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_event_proto_goTypes,
		DependencyIndexes: file_proto_event_proto_depIdxs,
		EnumInfos:         file_proto_event_proto_enumTypes,
		MessageInfos:      file_proto_event_proto_msgTypes,
	}.Build()
	File_proto_event_proto = out.File
	file_proto_event_proto_rawDesc = nil
	file_proto_event_proto_goTypes = nil
	file_proto_event_proto_depIdxs = nil
}