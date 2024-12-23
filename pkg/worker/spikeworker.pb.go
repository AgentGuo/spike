// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.20.3
// source: pkg/worker/spikeworker.proto

package worker

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

type CallWorkerFunctionReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload   string `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	RequestId string `protobuf:"bytes,2,opt,name=requestId,proto3" json:"requestId,omitempty"`
}

func (x *CallWorkerFunctionReq) Reset() {
	*x = CallWorkerFunctionReq{}
	mi := &file_pkg_worker_spikeworker_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CallWorkerFunctionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallWorkerFunctionReq) ProtoMessage() {}

func (x *CallWorkerFunctionReq) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_worker_spikeworker_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallWorkerFunctionReq.ProtoReflect.Descriptor instead.
func (*CallWorkerFunctionReq) Descriptor() ([]byte, []int) {
	return file_pkg_worker_spikeworker_proto_rawDescGZIP(), []int{0}
}

func (x *CallWorkerFunctionReq) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

func (x *CallWorkerFunctionReq) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

type CallWorkerFunctionResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload   string `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	RequestId string `protobuf:"bytes,2,opt,name=requestId,proto3" json:"requestId,omitempty"`
}

func (x *CallWorkerFunctionResp) Reset() {
	*x = CallWorkerFunctionResp{}
	mi := &file_pkg_worker_spikeworker_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CallWorkerFunctionResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallWorkerFunctionResp) ProtoMessage() {}

func (x *CallWorkerFunctionResp) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_worker_spikeworker_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallWorkerFunctionResp.ProtoReflect.Descriptor instead.
func (*CallWorkerFunctionResp) Descriptor() ([]byte, []int) {
	return file_pkg_worker_spikeworker_proto_rawDescGZIP(), []int{1}
}

func (x *CallWorkerFunctionResp) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

func (x *CallWorkerFunctionResp) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

var File_pkg_worker_spikeworker_proto protoreflect.FileDescriptor

var file_pkg_worker_spikeworker_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x6b, 0x67, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2f, 0x73, 0x70, 0x69,
	0x6b, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4f,
	0x0a, 0x15, 0x43, 0x61, 0x6c, 0x6c, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x46, 0x75, 0x6e, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x22,
	0x50, 0x0a, 0x16, 0x43, 0x61, 0x6c, 0x6c, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x46, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49,
	0x64, 0x32, 0x5b, 0x0a, 0x12, 0x53, 0x70, 0x69, 0x6b, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x12, 0x43, 0x61, 0x6c, 0x6c, 0x57,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e,
	0x43, 0x61, 0x6c, 0x6c, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x57, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x42, 0x55,
	0x0a, 0x20, 0x69, 0x6f, 0x2e, 0x70, 0x69, 0x78, 0x65, 0x6c, 0x73, 0x64, 0x62, 0x2e, 0x70, 0x69,
	0x78, 0x65, 0x6c, 0x73, 0x2e, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x68, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x72, 0x42, 0x0b, 0x53, 0x70, 0x69, 0x6b, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5a,
	0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x67, 0x65, 0x6e,
	0x74, 0x47, 0x75, 0x6f, 0x2f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_worker_spikeworker_proto_rawDescOnce sync.Once
	file_pkg_worker_spikeworker_proto_rawDescData = file_pkg_worker_spikeworker_proto_rawDesc
)

func file_pkg_worker_spikeworker_proto_rawDescGZIP() []byte {
	file_pkg_worker_spikeworker_proto_rawDescOnce.Do(func() {
		file_pkg_worker_spikeworker_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_worker_spikeworker_proto_rawDescData)
	})
	return file_pkg_worker_spikeworker_proto_rawDescData
}

var file_pkg_worker_spikeworker_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_worker_spikeworker_proto_goTypes = []any{
	(*CallWorkerFunctionReq)(nil),  // 0: CallWorkerFunctionReq
	(*CallWorkerFunctionResp)(nil), // 1: CallWorkerFunctionResp
}
var file_pkg_worker_spikeworker_proto_depIdxs = []int32{
	0, // 0: SpikeWorkerService.CallWorkerFunction:input_type -> CallWorkerFunctionReq
	1, // 1: SpikeWorkerService.CallWorkerFunction:output_type -> CallWorkerFunctionResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_worker_spikeworker_proto_init() }
func file_pkg_worker_spikeworker_proto_init() {
	if File_pkg_worker_spikeworker_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_worker_spikeworker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_worker_spikeworker_proto_goTypes,
		DependencyIndexes: file_pkg_worker_spikeworker_proto_depIdxs,
		MessageInfos:      file_pkg_worker_spikeworker_proto_msgTypes,
	}.Build()
	File_pkg_worker_spikeworker_proto = out.File
	file_pkg_worker_spikeworker_proto_rawDesc = nil
	file_pkg_worker_spikeworker_proto_goTypes = nil
	file_pkg_worker_spikeworker_proto_depIdxs = nil
}
