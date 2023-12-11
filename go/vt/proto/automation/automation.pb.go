//
//Copyright 2019 The Vitess Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// Protobuf data structures for the automation framework.

// Messages (e.g. Task) are used both for checkpoint data and API access
// (e.g. retrieving the current status of a pending cluster operation).

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.3
// source: automation.proto

package automation

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

type ClusterOperationState int32

const (
	ClusterOperationState_UNKNOWN_CLUSTER_OPERATION_STATE ClusterOperationState = 0
	ClusterOperationState_CLUSTER_OPERATION_NOT_STARTED   ClusterOperationState = 1
	ClusterOperationState_CLUSTER_OPERATION_RUNNING       ClusterOperationState = 2
	ClusterOperationState_CLUSTER_OPERATION_DONE          ClusterOperationState = 3
)

// Enum value maps for ClusterOperationState.
var (
	ClusterOperationState_name = map[int32]string{
		0: "UNKNOWN_CLUSTER_OPERATION_STATE",
		1: "CLUSTER_OPERATION_NOT_STARTED",
		2: "CLUSTER_OPERATION_RUNNING",
		3: "CLUSTER_OPERATION_DONE",
	}
	ClusterOperationState_value = map[string]int32{
		"UNKNOWN_CLUSTER_OPERATION_STATE": 0,
		"CLUSTER_OPERATION_NOT_STARTED":   1,
		"CLUSTER_OPERATION_RUNNING":       2,
		"CLUSTER_OPERATION_DONE":          3,
	}
)

func (x ClusterOperationState) Enum() *ClusterOperationState {
	p := new(ClusterOperationState)
	*p = x
	return p
}

func (x ClusterOperationState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ClusterOperationState) Descriptor() protoreflect.EnumDescriptor {
	return file_automation_proto_enumTypes[0].Descriptor()
}

func (ClusterOperationState) Type() protoreflect.EnumType {
	return &file_automation_proto_enumTypes[0]
}

func (x ClusterOperationState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ClusterOperationState.Descriptor instead.
func (ClusterOperationState) EnumDescriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{0}
}

type TaskState int32

const (
	TaskState_UNKNOWN_TASK_STATE TaskState = 0
	TaskState_NOT_STARTED        TaskState = 1
	TaskState_RUNNING            TaskState = 2
	TaskState_DONE               TaskState = 3
)

// Enum value maps for TaskState.
var (
	TaskState_name = map[int32]string{
		0: "UNKNOWN_TASK_STATE",
		1: "NOT_STARTED",
		2: "RUNNING",
		3: "DONE",
	}
	TaskState_value = map[string]int32{
		"UNKNOWN_TASK_STATE": 0,
		"NOT_STARTED":        1,
		"RUNNING":            2,
		"DONE":               3,
	}
)

func (x TaskState) Enum() *TaskState {
	p := new(TaskState)
	*p = x
	return p
}

func (x TaskState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskState) Descriptor() protoreflect.EnumDescriptor {
	return file_automation_proto_enumTypes[1].Descriptor()
}

func (TaskState) Type() protoreflect.EnumType {
	return &file_automation_proto_enumTypes[1]
}

func (x TaskState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskState.Descriptor instead.
func (TaskState) EnumDescriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{1}
}

type ClusterOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// TaskContainer are processed sequentially, one at a time.
	SerialTasks []*TaskContainer `protobuf:"bytes,2,rep,name=serial_tasks,json=serialTasks,proto3" json:"serial_tasks,omitempty"`
	// Cached value. This has to be re-evaluated e.g. after a checkpoint load because running tasks may have already finished.
	State ClusterOperationState `protobuf:"varint,3,opt,name=state,proto3,enum=automation.ClusterOperationState" json:"state,omitempty"`
	// Error of the first task which failed. Set after state advanced to CLUSTER_OPERATION_DONE. If empty, all tasks succeeded. Cached value, see state above.
	Error string `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *ClusterOperation) Reset() {
	*x = ClusterOperation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterOperation) ProtoMessage() {}

func (x *ClusterOperation) ProtoReflect() protoreflect.Message {
	mi := &file_automation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterOperation.ProtoReflect.Descriptor instead.
func (*ClusterOperation) Descriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{0}
}

func (x *ClusterOperation) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ClusterOperation) GetSerialTasks() []*TaskContainer {
	if x != nil {
		return x.SerialTasks
	}
	return nil
}

func (x *ClusterOperation) GetState() ClusterOperationState {
	if x != nil {
		return x.State
	}
	return ClusterOperationState_UNKNOWN_CLUSTER_OPERATION_STATE
}

func (x *ClusterOperation) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// TaskContainer holds one or more task which may be executed in parallel.
// "concurrency", if > 0, limits the amount of concurrently executed tasks.
type TaskContainer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParallelTasks []*Task `protobuf:"bytes,1,rep,name=parallel_tasks,json=parallelTasks,proto3" json:"parallel_tasks,omitempty"`
	Concurrency   int32   `protobuf:"varint,2,opt,name=concurrency,proto3" json:"concurrency,omitempty"`
}

func (x *TaskContainer) Reset() {
	*x = TaskContainer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskContainer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskContainer) ProtoMessage() {}

func (x *TaskContainer) ProtoReflect() protoreflect.Message {
	mi := &file_automation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskContainer.ProtoReflect.Descriptor instead.
func (*TaskContainer) Descriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{1}
}

func (x *TaskContainer) GetParallelTasks() []*Task {
	if x != nil {
		return x.ParallelTasks
	}
	return nil
}

func (x *TaskContainer) GetConcurrency() int32 {
	if x != nil {
		return x.Concurrency
	}
	return 0
}

// Task represents a specific task which should be automatically executed.
type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Task specification.
	Name       string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Parameters map[string]string `protobuf:"bytes,2,rep,name=parameters,proto3" json:"parameters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Runtime data.
	Id    string    `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	State TaskState `protobuf:"varint,4,opt,name=state,proto3,enum=automation.TaskState" json:"state,omitempty"`
	// Set after state advanced to DONE.
	Output string `protobuf:"bytes,5,opt,name=output,proto3" json:"output,omitempty"`
	// Set after state advanced to DONE. If empty, the task did succeed.
	Error string `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_automation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{2}
}

func (x *Task) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Task) GetParameters() map[string]string {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *Task) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Task) GetState() TaskState {
	if x != nil {
		return x.State
	}
	return TaskState_UNKNOWN_TASK_STATE
}

func (x *Task) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

func (x *Task) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type EnqueueClusterOperationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Parameters map[string]string `protobuf:"bytes,2,rep,name=parameters,proto3" json:"parameters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *EnqueueClusterOperationRequest) Reset() {
	*x = EnqueueClusterOperationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnqueueClusterOperationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnqueueClusterOperationRequest) ProtoMessage() {}

func (x *EnqueueClusterOperationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_automation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnqueueClusterOperationRequest.ProtoReflect.Descriptor instead.
func (*EnqueueClusterOperationRequest) Descriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{3}
}

func (x *EnqueueClusterOperationRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EnqueueClusterOperationRequest) GetParameters() map[string]string {
	if x != nil {
		return x.Parameters
	}
	return nil
}

type EnqueueClusterOperationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *EnqueueClusterOperationResponse) Reset() {
	*x = EnqueueClusterOperationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnqueueClusterOperationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnqueueClusterOperationResponse) ProtoMessage() {}

func (x *EnqueueClusterOperationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_automation_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnqueueClusterOperationResponse.ProtoReflect.Descriptor instead.
func (*EnqueueClusterOperationResponse) Descriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{4}
}

func (x *EnqueueClusterOperationResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetClusterOperationStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetClusterOperationStateRequest) Reset() {
	*x = GetClusterOperationStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClusterOperationStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClusterOperationStateRequest) ProtoMessage() {}

func (x *GetClusterOperationStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_automation_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClusterOperationStateRequest.ProtoReflect.Descriptor instead.
func (*GetClusterOperationStateRequest) Descriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{5}
}

func (x *GetClusterOperationStateRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetClusterOperationStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State ClusterOperationState `protobuf:"varint,1,opt,name=state,proto3,enum=automation.ClusterOperationState" json:"state,omitempty"`
}

func (x *GetClusterOperationStateResponse) Reset() {
	*x = GetClusterOperationStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClusterOperationStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClusterOperationStateResponse) ProtoMessage() {}

func (x *GetClusterOperationStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_automation_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClusterOperationStateResponse.ProtoReflect.Descriptor instead.
func (*GetClusterOperationStateResponse) Descriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{6}
}

func (x *GetClusterOperationStateResponse) GetState() ClusterOperationState {
	if x != nil {
		return x.State
	}
	return ClusterOperationState_UNKNOWN_CLUSTER_OPERATION_STATE
}

type GetClusterOperationDetailsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetClusterOperationDetailsRequest) Reset() {
	*x = GetClusterOperationDetailsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClusterOperationDetailsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClusterOperationDetailsRequest) ProtoMessage() {}

func (x *GetClusterOperationDetailsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_automation_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClusterOperationDetailsRequest.ProtoReflect.Descriptor instead.
func (*GetClusterOperationDetailsRequest) Descriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{7}
}

func (x *GetClusterOperationDetailsRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetClusterOperationDetailsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Full snapshot of the execution e.g. including output of each task.
	ClusterOp *ClusterOperation `protobuf:"bytes,2,opt,name=cluster_op,json=clusterOp,proto3" json:"cluster_op,omitempty"`
}

func (x *GetClusterOperationDetailsResponse) Reset() {
	*x = GetClusterOperationDetailsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClusterOperationDetailsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClusterOperationDetailsResponse) ProtoMessage() {}

func (x *GetClusterOperationDetailsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_automation_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClusterOperationDetailsResponse.ProtoReflect.Descriptor instead.
func (*GetClusterOperationDetailsResponse) Descriptor() ([]byte, []int) {
	return file_automation_proto_rawDescGZIP(), []int{8}
}

func (x *GetClusterOperationDetailsResponse) GetClusterOp() *ClusterOperation {
	if x != nil {
		return x.ClusterOp
	}
	return nil
}

var File_automation_proto protoreflect.FileDescriptor

var file_automation_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xaf,
	0x01, 0x0a, 0x10, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x3c, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x74, 0x61,
	0x73, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x75, 0x74, 0x6f,
	0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x43, 0x6f, 0x6e, 0x74, 0x61,
	0x69, 0x6e, 0x65, 0x72, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x54, 0x61, 0x73, 0x6b,
	0x73, 0x12, 0x37, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x21, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x22, 0x6a, 0x0a, 0x0d, 0x54, 0x61, 0x73, 0x6b, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x12, 0x37, 0x0a, 0x0e, 0x70, 0x61, 0x72, 0x61, 0x6c, 0x6c, 0x65, 0x6c, 0x5f, 0x74, 0x61,
	0x73, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x75, 0x74, 0x6f,
	0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x0d, 0x70, 0x61, 0x72,
	0x61, 0x6c, 0x6c, 0x65, 0x6c, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f,
	0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0b, 0x63, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x22, 0x86, 0x02, 0x0a,
	0x04, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x40, 0x0a, 0x0a, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e,
	0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x2e,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2b, 0x0a, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x61, 0x75, 0x74,
	0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x1a, 0x3d, 0x0a, 0x0f, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xcf, 0x01, 0x0a, 0x1e, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x5a, 0x0a, 0x0a,
	0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x3a, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x45, 0x6e,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x1a, 0x3d, 0x0a, 0x0f, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x31, 0x0a, 0x1f, 0x45, 0x6e, 0x71, 0x75, 0x65,
	0x75, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x31, 0x0a, 0x1f, 0x47, 0x65,
	0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x5b, 0x0a,
	0x20, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x37, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x21, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x33, 0x0a, 0x21, 0x47, 0x65,
	0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x61, 0x0a, 0x22, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x0a, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x5f, 0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x75, 0x74, 0x6f,
	0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x4f, 0x70, 0x2a, 0x9a, 0x01, 0x0a, 0x15, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x23, 0x0a, 0x1f,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f,
	0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x10,
	0x00, 0x12, 0x21, 0x0a, 0x1d, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x4f, 0x50, 0x45,
	0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x52, 0x54,
	0x45, 0x44, 0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f,
	0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e,
	0x47, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x4f,
	0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x4f, 0x4e, 0x45, 0x10, 0x03, 0x2a,
	0x4b, 0x0a, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x12,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x52,
	0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47,
	0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x4f, 0x4e, 0x45, 0x10, 0x03, 0x42, 0x29, 0x5a, 0x27,
	0x76, 0x69, 0x74, 0x65, 0x73, 0x73, 0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x69, 0x74, 0x65, 0x73, 0x73,
	0x2f, 0x67, 0x6f, 0x2f, 0x76, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x74,
	0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_automation_proto_rawDescOnce sync.Once
	file_automation_proto_rawDescData = file_automation_proto_rawDesc
)

func file_automation_proto_rawDescGZIP() []byte {
	file_automation_proto_rawDescOnce.Do(func() {
		file_automation_proto_rawDescData = protoimpl.X.CompressGZIP(file_automation_proto_rawDescData)
	})
	return file_automation_proto_rawDescData
}

var file_automation_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_automation_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_automation_proto_goTypes = []interface{}{
	(ClusterOperationState)(0),                 // 0: automation.ClusterOperationState
	(TaskState)(0),                             // 1: automation.TaskState
	(*ClusterOperation)(nil),                   // 2: automation.ClusterOperation
	(*TaskContainer)(nil),                      // 3: automation.TaskContainer
	(*Task)(nil),                               // 4: automation.Task
	(*EnqueueClusterOperationRequest)(nil),     // 5: automation.EnqueueClusterOperationRequest
	(*EnqueueClusterOperationResponse)(nil),    // 6: automation.EnqueueClusterOperationResponse
	(*GetClusterOperationStateRequest)(nil),    // 7: automation.GetClusterOperationStateRequest
	(*GetClusterOperationStateResponse)(nil),   // 8: automation.GetClusterOperationStateResponse
	(*GetClusterOperationDetailsRequest)(nil),  // 9: automation.GetClusterOperationDetailsRequest
	(*GetClusterOperationDetailsResponse)(nil), // 10: automation.GetClusterOperationDetailsResponse
	nil, // 11: automation.Task.ParametersEntry
	nil, // 12: automation.EnqueueClusterOperationRequest.ParametersEntry
}
var file_automation_proto_depIdxs = []int32{
	3,  // 0: automation.ClusterOperation.serial_tasks:type_name -> automation.TaskContainer
	0,  // 1: automation.ClusterOperation.state:type_name -> automation.ClusterOperationState
	4,  // 2: automation.TaskContainer.parallel_tasks:type_name -> automation.Task
	11, // 3: automation.Task.parameters:type_name -> automation.Task.ParametersEntry
	1,  // 4: automation.Task.state:type_name -> automation.TaskState
	12, // 5: automation.EnqueueClusterOperationRequest.parameters:type_name -> automation.EnqueueClusterOperationRequest.ParametersEntry
	0,  // 6: automation.GetClusterOperationStateResponse.state:type_name -> automation.ClusterOperationState
	2,  // 7: automation.GetClusterOperationDetailsResponse.cluster_op:type_name -> automation.ClusterOperation
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_automation_proto_init() }
func file_automation_proto_init() {
	if File_automation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_automation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterOperation); i {
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
		file_automation_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskContainer); i {
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
		file_automation_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
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
		file_automation_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnqueueClusterOperationRequest); i {
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
		file_automation_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnqueueClusterOperationResponse); i {
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
		file_automation_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClusterOperationStateRequest); i {
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
		file_automation_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClusterOperationStateResponse); i {
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
		file_automation_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClusterOperationDetailsRequest); i {
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
		file_automation_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClusterOperationDetailsResponse); i {
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
			RawDescriptor: file_automation_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_automation_proto_goTypes,
		DependencyIndexes: file_automation_proto_depIdxs,
		EnumInfos:         file_automation_proto_enumTypes,
		MessageInfos:      file_automation_proto_msgTypes,
	}.Build()
	File_automation_proto = out.File
	file_automation_proto_rawDesc = nil
	file_automation_proto_goTypes = nil
	file_automation_proto_depIdxs = nil
}
