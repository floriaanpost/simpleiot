// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.12
// source: serial.proto

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

type Serial struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject string   `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Points  []*Point `protobuf:"bytes,2,rep,name=points,proto3" json:"points,omitempty"`
}

func (x *Serial) Reset() {
	*x = Serial{}
	if protoimpl.UnsafeEnabled {
		mi := &file_serial_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Serial) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Serial) ProtoMessage() {}

func (x *Serial) ProtoReflect() protoreflect.Message {
	mi := &file_serial_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Serial.ProtoReflect.Descriptor instead.
func (*Serial) Descriptor() ([]byte, []int) {
	return file_serial_proto_rawDescGZIP(), []int{0}
}

func (x *Serial) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *Serial) GetPoints() []*Point {
	if x != nil {
		return x.Points
	}
	return nil
}

var File_serial_proto protoreflect.FileDescriptor

var file_serial_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x70, 0x62, 0x1a, 0x0b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x45, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x21, 0x0a, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x06,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x42, 0x0d, 0x5a, 0x0b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_serial_proto_rawDescOnce sync.Once
	file_serial_proto_rawDescData = file_serial_proto_rawDesc
)

func file_serial_proto_rawDescGZIP() []byte {
	file_serial_proto_rawDescOnce.Do(func() {
		file_serial_proto_rawDescData = protoimpl.X.CompressGZIP(file_serial_proto_rawDescData)
	})
	return file_serial_proto_rawDescData
}

var file_serial_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_serial_proto_goTypes = []interface{}{
	(*Serial)(nil), // 0: pb.Serial
	(*Point)(nil),  // 1: pb.Point
}
var file_serial_proto_depIdxs = []int32{
	1, // 0: pb.Serial.points:type_name -> pb.Point
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_serial_proto_init() }
func file_serial_proto_init() {
	if File_serial_proto != nil {
		return
	}
	file_point_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_serial_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Serial); i {
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
			RawDescriptor: file_serial_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_serial_proto_goTypes,
		DependencyIndexes: file_serial_proto_depIdxs,
		MessageInfos:      file_serial_proto_msgTypes,
	}.Build()
	File_serial_proto = out.File
	file_serial_proto_rawDesc = nil
	file_serial_proto_goTypes = nil
	file_serial_proto_depIdxs = nil
}
