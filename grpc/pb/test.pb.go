// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: proto/test.proto

package pb

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

type CPU struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumberCores   uint32  `protobuf:"varint,1,opt,name=number_cores,json=numberCores,proto3" json:"number_cores,omitempty"`
	NumberThreads uint32  `protobuf:"varint,2,opt,name=number_threads,json=numberThreads,proto3" json:"number_threads,omitempty"`
	MaxHzG        float64 `protobuf:"fixed64,3,opt,name=max_hz_g,json=maxHzG,proto3" json:"max_hz_g,omitempty"`
}

func (x *CPU) Reset() {
	*x = CPU{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CPU) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CPU) ProtoMessage() {}

func (x *CPU) ProtoReflect() protoreflect.Message {
	mi := &file_proto_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CPU.ProtoReflect.Descriptor instead.
func (*CPU) Descriptor() ([]byte, []int) {
	return file_proto_test_proto_rawDescGZIP(), []int{0}
}

func (x *CPU) GetNumberCores() uint32 {
	if x != nil {
		return x.NumberCores
	}
	return 0
}

func (x *CPU) GetNumberThreads() uint32 {
	if x != nil {
		return x.NumberThreads
	}
	return 0
}

func (x *CPU) GetMaxHzG() float64 {
	if x != nil {
		return x.MaxHzG
	}
	return 0
}

var File_proto_test_proto protoreflect.FileDescriptor

var file_proto_test_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x69, 0x0a, 0x03, 0x43, 0x50, 0x55,
	0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x72, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f,
	0x72, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x74, 0x68,
	0x72, 0x65, 0x61, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x73, 0x12, 0x18, 0x0a, 0x08, 0x6d, 0x61,
	0x78, 0x5f, 0x68, 0x7a, 0x5f, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x6d, 0x61,
	0x78, 0x48, 0x7a, 0x47, 0x32, 0x2f, 0x0a, 0x0c, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x0a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x50, 0x55, 0x1a, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x43, 0x50, 0x55, 0x42, 0x05, 0x5a, 0x03, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_test_proto_rawDescOnce sync.Once
	file_proto_test_proto_rawDescData = file_proto_test_proto_rawDesc
)

func file_proto_test_proto_rawDescGZIP() []byte {
	file_proto_test_proto_rawDescOnce.Do(func() {
		file_proto_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_test_proto_rawDescData)
	})
	return file_proto_test_proto_rawDescData
}

var file_proto_test_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_test_proto_goTypes = []interface{}{
	(*CPU)(nil), // 0: proto.CPU
}
var file_proto_test_proto_depIdxs = []int32{
	0, // 0: proto.HelloService.Hello:input_type -> proto.CPU
	0, // 1: proto.HelloService.Hello:output_type -> proto.CPU
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_test_proto_init() }
func file_proto_test_proto_init() {
	if File_proto_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CPU); i {
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
			RawDescriptor: file_proto_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_test_proto_goTypes,
		DependencyIndexes: file_proto_test_proto_depIdxs,
		MessageInfos:      file_proto_test_proto_msgTypes,
	}.Build()
	File_proto_test_proto = out.File
	file_proto_test_proto_rawDesc = nil
	file_proto_test_proto_goTypes = nil
	file_proto_test_proto_depIdxs = nil
}
