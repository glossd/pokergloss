// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: mqsurvival/api.proto

package mqsurvival

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type TicketGift struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ToUserId string `protobuf:"bytes,1,opt,name=to_user_id,json=toUserId,proto3" json:"to_user_id,omitempty"`
	Tickets  int64  `protobuf:"varint,2,opt,name=tickets,proto3" json:"tickets,omitempty"`
}

func (x *TicketGift) Reset() {
	*x = TicketGift{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mqsurvival_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TicketGift) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TicketGift) ProtoMessage() {}

func (x *TicketGift) ProtoReflect() protoreflect.Message {
	mi := &file_mqsurvival_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TicketGift.ProtoReflect.Descriptor instead.
func (*TicketGift) Descriptor() ([]byte, []int) {
	return file_mqsurvival_api_proto_rawDescGZIP(), []int{0}
}

func (x *TicketGift) GetToUserId() string {
	if x != nil {
		return x.ToUserId
	}
	return ""
}

func (x *TicketGift) GetTickets() int64 {
	if x != nil {
		return x.Tickets
	}
	return 0
}

var File_mqsurvival_api_proto protoreflect.FileDescriptor

var file_mqsurvival_api_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6d, 0x71, 0x73, 0x75, 0x72, 0x76, 0x69, 0x76, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6d, 0x71, 0x73, 0x75, 0x72, 0x76, 0x69, 0x76,
	0x61, 0x6c, 0x22, 0x44, 0x0a, 0x0a, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x47, 0x69, 0x66, 0x74,
	0x12, 0x1c, 0x0a, 0x0a, 0x74, 0x6f, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x6f, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mqsurvival_api_proto_rawDescOnce sync.Once
	file_mqsurvival_api_proto_rawDescData = file_mqsurvival_api_proto_rawDesc
)

func file_mqsurvival_api_proto_rawDescGZIP() []byte {
	file_mqsurvival_api_proto_rawDescOnce.Do(func() {
		file_mqsurvival_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_mqsurvival_api_proto_rawDescData)
	})
	return file_mqsurvival_api_proto_rawDescData
}

var file_mqsurvival_api_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_mqsurvival_api_proto_goTypes = []interface{}{
	(*TicketGift)(nil), // 0: mqsurvival.TicketGift
}
var file_mqsurvival_api_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mqsurvival_api_proto_init() }
func file_mqsurvival_api_proto_init() {
	if File_mqsurvival_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mqsurvival_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TicketGift); i {
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
			RawDescriptor: file_mqsurvival_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mqsurvival_api_proto_goTypes,
		DependencyIndexes: file_mqsurvival_api_proto_depIdxs,
		MessageInfos:      file_mqsurvival_api_proto_msgTypes,
	}.Build()
	File_mqsurvival_api_proto = out.File
	file_mqsurvival_api_proto_rawDesc = nil
	file_mqsurvival_api_proto_goTypes = nil
	file_mqsurvival_api_proto_depIdxs = nil
}
