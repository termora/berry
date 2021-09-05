// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: rpc/members.proto

package rpc

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

type SendGuildMemberChunkData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Members []*Member `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty"`
}

func (x *SendGuildMemberChunkData) Reset() {
	*x = SendGuildMemberChunkData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_members_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendGuildMemberChunkData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendGuildMemberChunkData) ProtoMessage() {}

func (x *SendGuildMemberChunkData) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_members_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendGuildMemberChunkData.ProtoReflect.Descriptor instead.
func (*SendGuildMemberChunkData) Descriptor() ([]byte, []int) {
	return file_rpc_members_proto_rawDescGZIP(), []int{0}
}

func (x *SendGuildMemberChunkData) GetMembers() []*Member {
	if x != nil {
		return x.Members
	}
	return nil
}

type SendGuildMemberChunkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendGuildMemberChunkResponse) Reset() {
	*x = SendGuildMemberChunkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_members_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendGuildMemberChunkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendGuildMemberChunkResponse) ProtoMessage() {}

func (x *SendGuildMemberChunkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_members_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendGuildMemberChunkResponse.ProtoReflect.Descriptor instead.
func (*SendGuildMemberChunkResponse) Descriptor() ([]byte, []int) {
	return file_rpc_members_proto_rawDescGZIP(), []int{1}
}

type UpdateGuildMemberData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Member *Member `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
}

func (x *UpdateGuildMemberData) Reset() {
	*x = UpdateGuildMemberData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_members_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateGuildMemberData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateGuildMemberData) ProtoMessage() {}

func (x *UpdateGuildMemberData) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_members_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateGuildMemberData.ProtoReflect.Descriptor instead.
func (*UpdateGuildMemberData) Descriptor() ([]byte, []int) {
	return file_rpc_members_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateGuildMemberData) GetMember() *Member {
	if x != nil {
		return x.Member
	}
	return nil
}

type UpdateGuildMemberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkedSent bool `protobuf:"varint,1,opt,name=chunked_sent,json=chunkedSent,proto3" json:"chunked_sent,omitempty"`
}

func (x *UpdateGuildMemberResponse) Reset() {
	*x = UpdateGuildMemberResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_members_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateGuildMemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateGuildMemberResponse) ProtoMessage() {}

func (x *UpdateGuildMemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_members_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateGuildMemberResponse.ProtoReflect.Descriptor instead.
func (*UpdateGuildMemberResponse) Descriptor() ([]byte, []int) {
	return file_rpc_members_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateGuildMemberResponse) GetChunkedSent() bool {
	if x != nil {
		return x.ChunkedSent
	}
	return false
}

type RemoveGuildMemberData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *RemoveGuildMemberData) Reset() {
	*x = RemoveGuildMemberData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_members_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveGuildMemberData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveGuildMemberData) ProtoMessage() {}

func (x *RemoveGuildMemberData) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_members_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveGuildMemberData.ProtoReflect.Descriptor instead.
func (*RemoveGuildMemberData) Descriptor() ([]byte, []int) {
	return file_rpc_members_proto_rawDescGZIP(), []int{4}
}

func (x *RemoveGuildMemberData) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type RemoveGuildMemberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkedSent bool `protobuf:"varint,1,opt,name=chunked_sent,json=chunkedSent,proto3" json:"chunked_sent,omitempty"`
}

func (x *RemoveGuildMemberResponse) Reset() {
	*x = RemoveGuildMemberResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_members_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveGuildMemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveGuildMemberResponse) ProtoMessage() {}

func (x *RemoveGuildMemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_members_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveGuildMemberResponse.ProtoReflect.Descriptor instead.
func (*RemoveGuildMemberResponse) Descriptor() ([]byte, []int) {
	return file_rpc_members_proto_rawDescGZIP(), []int{5}
}

func (x *RemoveGuildMemberResponse) GetChunkedSent() bool {
	if x != nil {
		return x.ChunkedSent
	}
	return false
}

type Member struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId        uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Username      string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Discriminator string   `protobuf:"bytes,3,opt,name=discriminator,proto3" json:"discriminator,omitempty"`
	RoleIds       []uint64 `protobuf:"varint,4,rep,packed,name=role_ids,json=roleIds,proto3" json:"role_ids,omitempty"`
	Nickname      string   `protobuf:"bytes,5,opt,name=nickname,proto3" json:"nickname,omitempty"`
}

func (x *Member) Reset() {
	*x = Member{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_members_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Member) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Member) ProtoMessage() {}

func (x *Member) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_members_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Member.ProtoReflect.Descriptor instead.
func (*Member) Descriptor() ([]byte, []int) {
	return file_rpc_members_proto_rawDescGZIP(), []int{6}
}

func (x *Member) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Member) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Member) GetDiscriminator() string {
	if x != nil {
		return x.Discriminator
	}
	return ""
}

func (x *Member) GetRoleIds() []uint64 {
	if x != nil {
		return x.RoleIds
	}
	return nil
}

func (x *Member) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

var File_rpc_members_proto protoreflect.FileDescriptor

var file_rpc_members_proto_rawDesc = []byte{
	0x0a, 0x11, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3d, 0x0a, 0x18, 0x53, 0x65, 0x6e, 0x64, 0x47, 0x75, 0x69, 0x6c, 0x64,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x21, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x07, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x22, 0x1e, 0x0a, 0x1c, 0x53, 0x65, 0x6e, 0x64, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x38, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x75, 0x69, 0x6c,
	0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1f, 0x0a, 0x06, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x3e, 0x0a, 0x19,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x75,
	0x6e, 0x6b, 0x65, 0x64, 0x5f, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0b, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x65, 0x64, 0x53, 0x65, 0x6e, 0x74, 0x22, 0x30, 0x0a, 0x15,
	0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x3e,
	0x0a, 0x19, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x65, 0x64, 0x5f, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0b, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x65, 0x64, 0x53, 0x65, 0x6e, 0x74, 0x22, 0x9a,
	0x01, 0x0a, 0x06, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x24,
	0x0a, 0x0d, 0x64, 0x69, 0x73, 0x63, 0x72, 0x69, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x69, 0x73, 0x63, 0x72, 0x69, 0x6d, 0x69, 0x6e,
	0x61, 0x74, 0x6f, 0x72, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x04, 0x52, 0x07, 0x72, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x73, 0x12,
	0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0xf8, 0x01, 0x0a, 0x12,
	0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x50, 0x0a, 0x14, 0x53, 0x65, 0x6e, 0x64, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x19, 0x2e, 0x53, 0x65, 0x6e,
	0x64, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x68, 0x75, 0x6e,
	0x6b, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x1d, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x47, 0x75, 0x69, 0x6c,
	0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x75,
	0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x61, 0x74,
	0x61, 0x1a, 0x1a, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a,
	0x11, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x16, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x47, 0x75, 0x69, 0x6c, 0x64,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x1a, 0x2e, 0x52, 0x65, 0x6d,
	0x6f, 0x76, 0x65, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x72, 0x6d, 0x6f, 0x72, 0x61, 0x2f, 0x62, 0x65, 0x72,
	0x72, 0x79, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74,
	0x69, 0x63, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_members_proto_rawDescOnce sync.Once
	file_rpc_members_proto_rawDescData = file_rpc_members_proto_rawDesc
)

func file_rpc_members_proto_rawDescGZIP() []byte {
	file_rpc_members_proto_rawDescOnce.Do(func() {
		file_rpc_members_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_members_proto_rawDescData)
	})
	return file_rpc_members_proto_rawDescData
}

var file_rpc_members_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_rpc_members_proto_goTypes = []interface{}{
	(*SendGuildMemberChunkData)(nil),     // 0: SendGuildMemberChunkData
	(*SendGuildMemberChunkResponse)(nil), // 1: SendGuildMemberChunkResponse
	(*UpdateGuildMemberData)(nil),        // 2: UpdateGuildMemberData
	(*UpdateGuildMemberResponse)(nil),    // 3: UpdateGuildMemberResponse
	(*RemoveGuildMemberData)(nil),        // 4: RemoveGuildMemberData
	(*RemoveGuildMemberResponse)(nil),    // 5: RemoveGuildMemberResponse
	(*Member)(nil),                       // 6: Member
}
var file_rpc_members_proto_depIdxs = []int32{
	6, // 0: SendGuildMemberChunkData.members:type_name -> Member
	6, // 1: UpdateGuildMemberData.member:type_name -> Member
	0, // 2: GuildMemberService.SendGuildMemberChunk:input_type -> SendGuildMemberChunkData
	2, // 3: GuildMemberService.UpdateGuildMember:input_type -> UpdateGuildMemberData
	4, // 4: GuildMemberService.RemoveGuildMember:input_type -> RemoveGuildMemberData
	1, // 5: GuildMemberService.SendGuildMemberChunk:output_type -> SendGuildMemberChunkResponse
	3, // 6: GuildMemberService.UpdateGuildMember:output_type -> UpdateGuildMemberResponse
	5, // 7: GuildMemberService.RemoveGuildMember:output_type -> RemoveGuildMemberResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rpc_members_proto_init() }
func file_rpc_members_proto_init() {
	if File_rpc_members_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_members_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendGuildMemberChunkData); i {
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
		file_rpc_members_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendGuildMemberChunkResponse); i {
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
		file_rpc_members_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateGuildMemberData); i {
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
		file_rpc_members_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateGuildMemberResponse); i {
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
		file_rpc_members_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveGuildMemberData); i {
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
		file_rpc_members_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveGuildMemberResponse); i {
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
		file_rpc_members_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Member); i {
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
			RawDescriptor: file_rpc_members_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_members_proto_goTypes,
		DependencyIndexes: file_rpc_members_proto_depIdxs,
		MessageInfos:      file_rpc_members_proto_msgTypes,
	}.Build()
	File_rpc_members_proto = out.File
	file_rpc_members_proto_rawDesc = nil
	file_rpc_members_proto_goTypes = nil
	file_rpc_members_proto_depIdxs = nil
}