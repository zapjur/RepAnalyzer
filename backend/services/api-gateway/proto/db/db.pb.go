// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: proto/db/db.proto

package db

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

type GetUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Auth0Id       string                 `protobuf:"bytes,1,opt,name=auth0_id,json=auth0Id,proto3" json:"auth0_id,omitempty"`
	Email         string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	mi := &file_proto_db_db_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_db_db_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserRequest.ProtoReflect.Descriptor instead.
func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_db_db_proto_rawDescGZIP(), []int{0}
}

func (x *GetUserRequest) GetAuth0Id() string {
	if x != nil {
		return x.Auth0Id
	}
	return ""
}

func (x *GetUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type GetUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserResponse) Reset() {
	*x = GetUserResponse{}
	mi := &file_proto_db_db_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserResponse) ProtoMessage() {}

func (x *GetUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_db_db_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserResponse.ProtoReflect.Descriptor instead.
func (*GetUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_db_db_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *GetUserResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type UploadVideoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	ExerciseName  string                 `protobuf:"bytes,2,opt,name=exercise_name,json=exerciseName,proto3" json:"exercise_name,omitempty"`
	Auth0Id       string                 `protobuf:"bytes,3,opt,name=auth0_id,json=auth0Id,proto3" json:"auth0_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadVideoRequest) Reset() {
	*x = UploadVideoRequest{}
	mi := &file_proto_db_db_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadVideoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadVideoRequest) ProtoMessage() {}

func (x *UploadVideoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_db_db_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadVideoRequest.ProtoReflect.Descriptor instead.
func (*UploadVideoRequest) Descriptor() ([]byte, []int) {
	return file_proto_db_db_proto_rawDescGZIP(), []int{2}
}

func (x *UploadVideoRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *UploadVideoRequest) GetExerciseName() string {
	if x != nil {
		return x.ExerciseName
	}
	return ""
}

func (x *UploadVideoRequest) GetAuth0Id() string {
	if x != nil {
		return x.Auth0Id
	}
	return ""
}

type UploadVideoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	VideoId       int64                  `protobuf:"varint,3,opt,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadVideoResponse) Reset() {
	*x = UploadVideoResponse{}
	mi := &file_proto_db_db_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadVideoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadVideoResponse) ProtoMessage() {}

func (x *UploadVideoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_db_db_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadVideoResponse.ProtoReflect.Descriptor instead.
func (*UploadVideoResponse) Descriptor() ([]byte, []int) {
	return file_proto_db_db_proto_rawDescGZIP(), []int{3}
}

func (x *UploadVideoResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UploadVideoResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *UploadVideoResponse) GetVideoId() int64 {
	if x != nil {
		return x.VideoId
	}
	return 0
}

type GetUserVideosByExerciseRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Auth0Id       string                 `protobuf:"bytes,1,opt,name=auth0_id,json=auth0Id,proto3" json:"auth0_id,omitempty"`
	ExerciseName  string                 `protobuf:"bytes,2,opt,name=exercise_name,json=exerciseName,proto3" json:"exercise_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserVideosByExerciseRequest) Reset() {
	*x = GetUserVideosByExerciseRequest{}
	mi := &file_proto_db_db_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserVideosByExerciseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserVideosByExerciseRequest) ProtoMessage() {}

func (x *GetUserVideosByExerciseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_db_db_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserVideosByExerciseRequest.ProtoReflect.Descriptor instead.
func (*GetUserVideosByExerciseRequest) Descriptor() ([]byte, []int) {
	return file_proto_db_db_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserVideosByExerciseRequest) GetAuth0Id() string {
	if x != nil {
		return x.Auth0Id
	}
	return ""
}

func (x *GetUserVideosByExerciseRequest) GetExerciseName() string {
	if x != nil {
		return x.ExerciseName
	}
	return ""
}

type VideoInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	ExerciseName  string                 `protobuf:"bytes,2,opt,name=exercise_name,json=exerciseName,proto3" json:"exercise_name,omitempty"`
	Auth0Id       string                 `protobuf:"bytes,3,opt,name=auth0_id,json=auth0Id,proto3" json:"auth0_id,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Id            int64                  `protobuf:"varint,5,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VideoInfo) Reset() {
	*x = VideoInfo{}
	mi := &file_proto_db_db_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VideoInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoInfo) ProtoMessage() {}

func (x *VideoInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_db_db_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoInfo.ProtoReflect.Descriptor instead.
func (*VideoInfo) Descriptor() ([]byte, []int) {
	return file_proto_db_db_proto_rawDescGZIP(), []int{5}
}

func (x *VideoInfo) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *VideoInfo) GetExerciseName() string {
	if x != nil {
		return x.ExerciseName
	}
	return ""
}

func (x *VideoInfo) GetAuth0Id() string {
	if x != nil {
		return x.Auth0Id
	}
	return ""
}

func (x *VideoInfo) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *VideoInfo) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type VideoAnalysisRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	VideoId       int64                  `protobuf:"varint,1,opt,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
	Type          string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	ResultUrl     string                 `protobuf:"bytes,3,opt,name=result_url,json=resultUrl,proto3" json:"result_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VideoAnalysisRequest) Reset() {
	*x = VideoAnalysisRequest{}
	mi := &file_proto_db_db_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VideoAnalysisRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoAnalysisRequest) ProtoMessage() {}

func (x *VideoAnalysisRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_db_db_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoAnalysisRequest.ProtoReflect.Descriptor instead.
func (*VideoAnalysisRequest) Descriptor() ([]byte, []int) {
	return file_proto_db_db_proto_rawDescGZIP(), []int{6}
}

func (x *VideoAnalysisRequest) GetVideoId() int64 {
	if x != nil {
		return x.VideoId
	}
	return 0
}

func (x *VideoAnalysisRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *VideoAnalysisRequest) GetResultUrl() string {
	if x != nil {
		return x.ResultUrl
	}
	return ""
}

type SaveAnalysisResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SaveAnalysisResponse) Reset() {
	*x = SaveAnalysisResponse{}
	mi := &file_proto_db_db_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SaveAnalysisResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveAnalysisResponse) ProtoMessage() {}

func (x *SaveAnalysisResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_db_db_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveAnalysisResponse.ProtoReflect.Descriptor instead.
func (*SaveAnalysisResponse) Descriptor() ([]byte, []int) {
	return file_proto_db_db_proto_rawDescGZIP(), []int{7}
}

func (x *SaveAnalysisResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SaveAnalysisResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetUserVideosByExerciseResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Videos        []*VideoInfo           `protobuf:"bytes,1,rep,name=videos,proto3" json:"videos,omitempty"`
	Success       bool                   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserVideosByExerciseResponse) Reset() {
	*x = GetUserVideosByExerciseResponse{}
	mi := &file_proto_db_db_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserVideosByExerciseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserVideosByExerciseResponse) ProtoMessage() {}

func (x *GetUserVideosByExerciseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_db_db_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserVideosByExerciseResponse.ProtoReflect.Descriptor instead.
func (*GetUserVideosByExerciseResponse) Descriptor() ([]byte, []int) {
	return file_proto_db_db_proto_rawDescGZIP(), []int{8}
}

func (x *GetUserVideosByExerciseResponse) GetVideos() []*VideoInfo {
	if x != nil {
		return x.Videos
	}
	return nil
}

func (x *GetUserVideosByExerciseResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *GetUserVideosByExerciseResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_db_db_proto protoreflect.FileDescriptor

var file_proto_db_db_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x62, 0x2f, 0x64, 0x62, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x02, 0x64, 0x62, 0x22, 0x41, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x75, 0x74,
	0x68, 0x30, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x75, 0x74,
	0x68, 0x30, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x45, 0x0a, 0x0f, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x22, 0x66, 0x0a, 0x12, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x56, 0x69, 0x64, 0x65, 0x6f,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x78, 0x65,
	0x72, 0x63, 0x69, 0x73, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x65, 0x78, 0x65, 0x72, 0x63, 0x69, 0x73, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19,
	0x0a, 0x08, 0x61, 0x75, 0x74, 0x68, 0x30, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x30, 0x49, 0x64, 0x22, 0x64, 0x0a, 0x13, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x22,
	0x60, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x73,
	0x42, 0x79, 0x45, 0x78, 0x65, 0x72, 0x63, 0x69, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x75, 0x74, 0x68, 0x30, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x30, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d,
	0x65, 0x78, 0x65, 0x72, 0x63, 0x69, 0x73, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x78, 0x65, 0x72, 0x63, 0x69, 0x73, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x22, 0x8c, 0x01, 0x0a, 0x09, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x78, 0x65, 0x72, 0x63, 0x69, 0x73, 0x65, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x78, 0x65, 0x72, 0x63, 0x69,
	0x73, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x75, 0x74, 0x68, 0x30, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x30, 0x49,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x64, 0x0a, 0x14, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x69, 0x64, 0x65,
	0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x69, 0x64, 0x65,
	0x6f, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x55, 0x72, 0x6c, 0x22, 0x4a, 0x0a, 0x14, 0x53, 0x61, 0x76, 0x65, 0x41, 0x6e,
	0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x7c, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x56, 0x69, 0x64,
	0x65, 0x6f, 0x73, 0x42, 0x79, 0x45, 0x78, 0x65, 0x72, 0x63, 0x69, 0x73, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x06, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x64, 0x62, 0x2e, 0x56, 0x69, 0x64, 0x65, 0x6f,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x06, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0xad, 0x02, 0x0a, 0x09, 0x44, 0x42, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32,
	0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x64, 0x62, 0x2e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e,
	0x64, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x44, 0x0a, 0x11, 0x53, 0x61, 0x76, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x65, 0x64, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x12, 0x16, 0x2e, 0x64, 0x62, 0x2e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x64, 0x62, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x56, 0x69, 0x64, 0x65, 0x6f,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x42, 0x79, 0x45, 0x78, 0x65, 0x72, 0x63,
	0x69, 0x73, 0x65, 0x12, 0x22, 0x2e, 0x64, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x56, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x42, 0x79, 0x45, 0x78, 0x65, 0x72, 0x63, 0x69, 0x73, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x64, 0x62, 0x2e, 0x47, 0x65, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x42, 0x79, 0x45, 0x78, 0x65, 0x72,
	0x63, 0x69, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x0c,
	0x53, 0x61, 0x76, 0x65, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x12, 0x18, 0x2e, 0x64,
	0x62, 0x2e, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x64, 0x62, 0x2e, 0x53, 0x61, 0x76, 0x65,
	0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x64, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_db_db_proto_rawDescOnce sync.Once
	file_proto_db_db_proto_rawDescData []byte
)

func file_proto_db_db_proto_rawDescGZIP() []byte {
	file_proto_db_db_proto_rawDescOnce.Do(func() {
		file_proto_db_db_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_db_db_proto_rawDesc), len(file_proto_db_db_proto_rawDesc)))
	})
	return file_proto_db_db_proto_rawDescData
}

var file_proto_db_db_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_db_db_proto_goTypes = []any{
	(*GetUserRequest)(nil),                  // 0: db.GetUserRequest
	(*GetUserResponse)(nil),                 // 1: db.GetUserResponse
	(*UploadVideoRequest)(nil),              // 2: db.UploadVideoRequest
	(*UploadVideoResponse)(nil),             // 3: db.UploadVideoResponse
	(*GetUserVideosByExerciseRequest)(nil),  // 4: db.GetUserVideosByExerciseRequest
	(*VideoInfo)(nil),                       // 5: db.VideoInfo
	(*VideoAnalysisRequest)(nil),            // 6: db.VideoAnalysisRequest
	(*SaveAnalysisResponse)(nil),            // 7: db.SaveAnalysisResponse
	(*GetUserVideosByExerciseResponse)(nil), // 8: db.GetUserVideosByExerciseResponse
}
var file_proto_db_db_proto_depIdxs = []int32{
	5, // 0: db.GetUserVideosByExerciseResponse.videos:type_name -> db.VideoInfo
	0, // 1: db.DBService.GetUser:input_type -> db.GetUserRequest
	2, // 2: db.DBService.SaveUploadedVideo:input_type -> db.UploadVideoRequest
	4, // 3: db.DBService.GetUserVideosByExercise:input_type -> db.GetUserVideosByExerciseRequest
	6, // 4: db.DBService.SaveAnalysis:input_type -> db.VideoAnalysisRequest
	1, // 5: db.DBService.GetUser:output_type -> db.GetUserResponse
	3, // 6: db.DBService.SaveUploadedVideo:output_type -> db.UploadVideoResponse
	8, // 7: db.DBService.GetUserVideosByExercise:output_type -> db.GetUserVideosByExerciseResponse
	7, // 8: db.DBService.SaveAnalysis:output_type -> db.SaveAnalysisResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_db_db_proto_init() }
func file_proto_db_db_proto_init() {
	if File_proto_db_db_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_db_db_proto_rawDesc), len(file_proto_db_db_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_db_db_proto_goTypes,
		DependencyIndexes: file_proto_db_db_proto_depIdxs,
		MessageInfos:      file_proto_db_db_proto_msgTypes,
	}.Build()
	File_proto_db_db_proto = out.File
	file_proto_db_db_proto_goTypes = nil
	file_proto_db_db_proto_depIdxs = nil
}
