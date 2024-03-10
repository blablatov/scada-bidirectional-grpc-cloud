// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.4
// source: grpc-cloud.proto

//option go_package = "./;cloudservice"; //dir of create proto-file

package __

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

type RequestIO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Sensors     []string `protobuf:"bytes,2,rep,name=sensors,proto3" json:"sensors,omitempty"`
	Description string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Measurement float64  `protobuf:"fixed64,4,opt,name=measurement,proto3" json:"measurement,omitempty"`
	Destination string   `protobuf:"bytes,5,opt,name=destination,proto3" json:"destination,omitempty"`
	Sensor      string   `protobuf:"bytes,6,opt,name=sensor,proto3" json:"sensor,omitempty"`
}

func (x *RequestIO) Reset() {
	*x = RequestIO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_cloud_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestIO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestIO) ProtoMessage() {}

func (x *RequestIO) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_cloud_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestIO.ProtoReflect.Descriptor instead.
func (*RequestIO) Descriptor() ([]byte, []int) {
	return file_grpc_cloud_proto_rawDescGZIP(), []int{0}
}

func (x *RequestIO) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RequestIO) GetSensors() []string {
	if x != nil {
		return x.Sensors
	}
	return nil
}

func (x *RequestIO) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *RequestIO) GetMeasurement() float64 {
	if x != nil {
		return x.Measurement
	}
	return 0
}

func (x *RequestIO) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

func (x *RequestIO) GetSensor() string {
	if x != nil {
		return x.Sensor
	}
	return ""
}

type StatusIO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status string       `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	IOList []*RequestIO `protobuf:"bytes,3,rep,name=IOList,proto3" json:"IOList,omitempty"`
}

func (x *StatusIO) Reset() {
	*x = StatusIO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_cloud_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusIO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusIO) ProtoMessage() {}

func (x *StatusIO) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_cloud_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusIO.ProtoReflect.Descriptor instead.
func (*StatusIO) Descriptor() ([]byte, []int) {
	return file_grpc_cloud_proto_rawDescGZIP(), []int{1}
}

func (x *StatusIO) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *StatusIO) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *StatusIO) GetIOList() []*RequestIO {
	if x != nil {
		return x.IOList
	}
	return nil
}

// Номера и имена зарезервированных полей сообщений. Don't use this
type Res struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Res) Reset() {
	*x = Res{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_cloud_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Res) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Res) ProtoMessage() {}

func (x *Res) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_cloud_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Res.ProtoReflect.Descriptor instead.
func (*Res) Descriptor() ([]byte, []int) {
	return file_grpc_cloud_proto_rawDescGZIP(), []int{2}
}

var File_grpc_cloud_proto protoreflect.FileDescriptor

var file_grpc_cloud_proto_rawDesc = []byte{
	0x0a, 0x10, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0c, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x22, 0xb3, 0x01, 0x0a, 0x09, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x4f, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x07, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x6d, 0x65,
	0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x0b, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x22, 0x63, 0x0a, 0x08, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x49, 0x4f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x06, 0x49, 0x4f,
	0x4c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x49, 0x4f, 0x52, 0x06, 0x49, 0x4f, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x43, 0x0a, 0x03, 0x52,
	0x65, 0x73, 0x4a, 0x04, 0x08, 0x07, 0x10, 0x08, 0x4a, 0x04, 0x08, 0x08, 0x10, 0x09, 0x4a, 0x04,
	0x08, 0x09, 0x10, 0x11, 0x4a, 0x08, 0x08, 0x63, 0x10, 0x80, 0x80, 0x80, 0x80, 0x02, 0x52, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x0a, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x32, 0x54, 0x0a, 0x0d, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x12, 0x43, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x43, 0x6c, 0x6f, 0x75,
	0x64, 0x12, 0x17, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x4f, 0x1a, 0x16, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x49, 0x4f, 0x28, 0x01, 0x30, 0x01, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_cloud_proto_rawDescOnce sync.Once
	file_grpc_cloud_proto_rawDescData = file_grpc_cloud_proto_rawDesc
)

func file_grpc_cloud_proto_rawDescGZIP() []byte {
	file_grpc_cloud_proto_rawDescOnce.Do(func() {
		file_grpc_cloud_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_cloud_proto_rawDescData)
	})
	return file_grpc_cloud_proto_rawDescData
}

var file_grpc_cloud_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_grpc_cloud_proto_goTypes = []interface{}{
	(*RequestIO)(nil), // 0: cloudservice.RequestIO
	(*StatusIO)(nil),  // 1: cloudservice.StatusIO
	(*Res)(nil),       // 2: cloudservice.Res
}
var file_grpc_cloud_proto_depIdxs = []int32{
	0, // 0: cloudservice.StatusIO.IOList:type_name -> cloudservice.RequestIO
	0, // 1: cloudservice.CloudExchange.processCloud:input_type -> cloudservice.RequestIO
	1, // 2: cloudservice.CloudExchange.processCloud:output_type -> cloudservice.StatusIO
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_grpc_cloud_proto_init() }
func file_grpc_cloud_proto_init() {
	if File_grpc_cloud_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_cloud_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestIO); i {
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
		file_grpc_cloud_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusIO); i {
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
		file_grpc_cloud_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Res); i {
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
			RawDescriptor: file_grpc_cloud_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_cloud_proto_goTypes,
		DependencyIndexes: file_grpc_cloud_proto_depIdxs,
		MessageInfos:      file_grpc_cloud_proto_msgTypes,
	}.Build()
	File_grpc_cloud_proto = out.File
	file_grpc_cloud_proto_rawDesc = nil
	file_grpc_cloud_proto_goTypes = nil
	file_grpc_cloud_proto_depIdxs = nil
}
