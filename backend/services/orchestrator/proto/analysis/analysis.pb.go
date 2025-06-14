// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: proto/analysis/analysis.proto

package analysis

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VideoToAnalyzeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Bucket        string                 `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	ObjectKey     string                 `protobuf:"bytes,2,opt,name=object_key,json=objectKey,proto3" json:"object_key,omitempty"`
	ExerciseName  string                 `protobuf:"bytes,3,opt,name=exercise_name,json=exerciseName,proto3" json:"exercise_name,omitempty"`
	Auth0Id       string                 `protobuf:"bytes,4,opt,name=auth0_id,json=auth0Id,proto3" json:"auth0_id,omitempty"`
	VideoId       int64                  `protobuf:"varint,5,opt,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VideoToAnalyzeRequest) Reset() {
	*x = VideoToAnalyzeRequest{}
	mi := &file_proto_analysis_analysis_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VideoToAnalyzeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoToAnalyzeRequest) ProtoMessage() {}

func (x *VideoToAnalyzeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_analysis_analysis_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoToAnalyzeRequest.ProtoReflect.Descriptor instead.
func (*VideoToAnalyzeRequest) Descriptor() ([]byte, []int) {
	return file_proto_analysis_analysis_proto_rawDescGZIP(), []int{0}
}

func (x *VideoToAnalyzeRequest) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

func (x *VideoToAnalyzeRequest) GetObjectKey() string {
	if x != nil {
		return x.ObjectKey
	}
	return ""
}

func (x *VideoToAnalyzeRequest) GetExerciseName() string {
	if x != nil {
		return x.ExerciseName
	}
	return ""
}

func (x *VideoToAnalyzeRequest) GetAuth0Id() string {
	if x != nil {
		return x.Auth0Id
	}
	return ""
}

func (x *VideoToAnalyzeRequest) GetVideoId() int64 {
	if x != nil {
		return x.VideoId
	}
	return 0
}

type VideoToAnalyzeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VideoToAnalyzeResponse) Reset() {
	*x = VideoToAnalyzeResponse{}
	mi := &file_proto_analysis_analysis_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VideoToAnalyzeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoToAnalyzeResponse) ProtoMessage() {}

func (x *VideoToAnalyzeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_analysis_analysis_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoToAnalyzeResponse.ProtoReflect.Descriptor instead.
func (*VideoToAnalyzeResponse) Descriptor() ([]byte, []int) {
	return file_proto_analysis_analysis_proto_rawDescGZIP(), []int{1}
}

func (x *VideoToAnalyzeResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *VideoToAnalyzeResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_analysis_analysis_proto protoreflect.FileDescriptor

var file_proto_analysis_analysis_proto_rawDesc = string([]byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73,
	0x2f, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x22, 0xa9, 0x01, 0x0a, 0x15, 0x56, 0x69,
	0x64, 0x65, 0x6f, 0x54, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x78,
	0x65, 0x72, 0x63, 0x69, 0x73, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x65, 0x78, 0x65, 0x72, 0x63, 0x69, 0x73, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x61, 0x75, 0x74, 0x68, 0x30, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x30, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x49, 0x64, 0x22, 0x4c, 0x0a, 0x16, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x54, 0x6f,
	0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x32, 0x61, 0x0a, 0x0c, 0x4f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x12, 0x51, 0x0a, 0x0c, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x56, 0x69,
	0x64, 0x65, 0x6f, 0x12, 0x1f, 0x2e, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x2e, 0x56,
	0x69, 0x64, 0x65, 0x6f, 0x54, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x2e,
	0x56, 0x69, 0x64, 0x65, 0x6f, 0x54, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x12, 0x5a, 0x10, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x3b, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
})

var (
	file_proto_analysis_analysis_proto_rawDescOnce sync.Once
	file_proto_analysis_analysis_proto_rawDescData []byte
)

func file_proto_analysis_analysis_proto_rawDescGZIP() []byte {
	file_proto_analysis_analysis_proto_rawDescOnce.Do(func() {
		file_proto_analysis_analysis_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_analysis_analysis_proto_rawDesc), len(file_proto_analysis_analysis_proto_rawDesc)))
	})
	return file_proto_analysis_analysis_proto_rawDescData
}

var file_proto_analysis_analysis_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_analysis_analysis_proto_goTypes = []any{
	(*VideoToAnalyzeRequest)(nil),  // 0: analysis.VideoToAnalyzeRequest
	(*VideoToAnalyzeResponse)(nil), // 1: analysis.VideoToAnalyzeResponse
}
var file_proto_analysis_analysis_proto_depIdxs = []int32{
	0, // 0: analysis.Orchestrator.AnalyzeVideo:input_type -> analysis.VideoToAnalyzeRequest
	1, // 1: analysis.Orchestrator.AnalyzeVideo:output_type -> analysis.VideoToAnalyzeResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_analysis_analysis_proto_init() }
func file_proto_analysis_analysis_proto_init() {
	if File_proto_analysis_analysis_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_analysis_analysis_proto_rawDesc), len(file_proto_analysis_analysis_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_analysis_analysis_proto_goTypes,
		DependencyIndexes: file_proto_analysis_analysis_proto_depIdxs,
		MessageInfos:      file_proto_analysis_analysis_proto_msgTypes,
	}.Build()
	File_proto_analysis_analysis_proto = out.File
	file_proto_analysis_analysis_proto_goTypes = nil
	file_proto_analysis_analysis_proto_depIdxs = nil
}
