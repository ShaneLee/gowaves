// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transactions_api.proto

package generated

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TransactionStatus_Status int32

const (
	TransactionStatus_NOT_EXISTS  TransactionStatus_Status = 0
	TransactionStatus_UNCONFIRMED TransactionStatus_Status = 1
	TransactionStatus_CONFIRMED   TransactionStatus_Status = 2
)

var TransactionStatus_Status_name = map[int32]string{
	0: "NOT_EXISTS",
	1: "UNCONFIRMED",
	2: "CONFIRMED",
}

var TransactionStatus_Status_value = map[string]int32{
	"NOT_EXISTS":  0,
	"UNCONFIRMED": 1,
	"CONFIRMED":   2,
}

func (x TransactionStatus_Status) String() string {
	return proto.EnumName(TransactionStatus_Status_name, int32(x))
}

func (TransactionStatus_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_121a662cf7c9700a, []int{0, 0}
}

type TransactionStatus struct {
	Id                   []byte                   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status               TransactionStatus_Status `protobuf:"varint,2,opt,name=status,proto3,enum=waves.node.grpc.TransactionStatus_Status" json:"status,omitempty"`
	Height               int64                    `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *TransactionStatus) Reset()         { *m = TransactionStatus{} }
func (m *TransactionStatus) String() string { return proto.CompactTextString(m) }
func (*TransactionStatus) ProtoMessage()    {}
func (*TransactionStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_121a662cf7c9700a, []int{0}
}

func (m *TransactionStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionStatus.Unmarshal(m, b)
}
func (m *TransactionStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionStatus.Marshal(b, m, deterministic)
}
func (m *TransactionStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionStatus.Merge(m, src)
}
func (m *TransactionStatus) XXX_Size() int {
	return xxx_messageInfo_TransactionStatus.Size(m)
}
func (m *TransactionStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionStatus.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionStatus proto.InternalMessageInfo

func (m *TransactionStatus) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *TransactionStatus) GetStatus() TransactionStatus_Status {
	if m != nil {
		return m.Status
	}
	return TransactionStatus_NOT_EXISTS
}

func (m *TransactionStatus) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type TransactionResponse struct {
	Id                   []byte             `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Height               int64              `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	Transaction          *SignedTransaction `protobuf:"bytes,3,opt,name=transaction,proto3" json:"transaction,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *TransactionResponse) Reset()         { *m = TransactionResponse{} }
func (m *TransactionResponse) String() string { return proto.CompactTextString(m) }
func (*TransactionResponse) ProtoMessage()    {}
func (*TransactionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_121a662cf7c9700a, []int{1}
}

func (m *TransactionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionResponse.Unmarshal(m, b)
}
func (m *TransactionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionResponse.Marshal(b, m, deterministic)
}
func (m *TransactionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionResponse.Merge(m, src)
}
func (m *TransactionResponse) XXX_Size() int {
	return xxx_messageInfo_TransactionResponse.Size(m)
}
func (m *TransactionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionResponse proto.InternalMessageInfo

func (m *TransactionResponse) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *TransactionResponse) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *TransactionResponse) GetTransaction() *SignedTransaction {
	if m != nil {
		return m.Transaction
	}
	return nil
}

type TransactionsRequest struct {
	Sender               []byte     `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Recipient            *Recipient `protobuf:"bytes,2,opt,name=recipient,proto3" json:"recipient,omitempty"`
	TransactionIds       [][]byte   `protobuf:"bytes,3,rep,name=transaction_ids,json=transactionIds,proto3" json:"transaction_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *TransactionsRequest) Reset()         { *m = TransactionsRequest{} }
func (m *TransactionsRequest) String() string { return proto.CompactTextString(m) }
func (*TransactionsRequest) ProtoMessage()    {}
func (*TransactionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_121a662cf7c9700a, []int{2}
}

func (m *TransactionsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionsRequest.Unmarshal(m, b)
}
func (m *TransactionsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionsRequest.Marshal(b, m, deterministic)
}
func (m *TransactionsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionsRequest.Merge(m, src)
}
func (m *TransactionsRequest) XXX_Size() int {
	return xxx_messageInfo_TransactionsRequest.Size(m)
}
func (m *TransactionsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionsRequest proto.InternalMessageInfo

func (m *TransactionsRequest) GetSender() []byte {
	if m != nil {
		return m.Sender
	}
	return nil
}

func (m *TransactionsRequest) GetRecipient() *Recipient {
	if m != nil {
		return m.Recipient
	}
	return nil
}

func (m *TransactionsRequest) GetTransactionIds() [][]byte {
	if m != nil {
		return m.TransactionIds
	}
	return nil
}

type TransactionsByIdRequest struct {
	TransactionIds       [][]byte `protobuf:"bytes,3,rep,name=transaction_ids,json=transactionIds,proto3" json:"transaction_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TransactionsByIdRequest) Reset()         { *m = TransactionsByIdRequest{} }
func (m *TransactionsByIdRequest) String() string { return proto.CompactTextString(m) }
func (*TransactionsByIdRequest) ProtoMessage()    {}
func (*TransactionsByIdRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_121a662cf7c9700a, []int{3}
}

func (m *TransactionsByIdRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionsByIdRequest.Unmarshal(m, b)
}
func (m *TransactionsByIdRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionsByIdRequest.Marshal(b, m, deterministic)
}
func (m *TransactionsByIdRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionsByIdRequest.Merge(m, src)
}
func (m *TransactionsByIdRequest) XXX_Size() int {
	return xxx_messageInfo_TransactionsByIdRequest.Size(m)
}
func (m *TransactionsByIdRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionsByIdRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionsByIdRequest proto.InternalMessageInfo

func (m *TransactionsByIdRequest) GetTransactionIds() [][]byte {
	if m != nil {
		return m.TransactionIds
	}
	return nil
}

type CalculateFeeResponse struct {
	AssetId              []byte   `protobuf:"bytes,1,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty"`
	Amount               uint64   `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CalculateFeeResponse) Reset()         { *m = CalculateFeeResponse{} }
func (m *CalculateFeeResponse) String() string { return proto.CompactTextString(m) }
func (*CalculateFeeResponse) ProtoMessage()    {}
func (*CalculateFeeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_121a662cf7c9700a, []int{4}
}

func (m *CalculateFeeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CalculateFeeResponse.Unmarshal(m, b)
}
func (m *CalculateFeeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CalculateFeeResponse.Marshal(b, m, deterministic)
}
func (m *CalculateFeeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CalculateFeeResponse.Merge(m, src)
}
func (m *CalculateFeeResponse) XXX_Size() int {
	return xxx_messageInfo_CalculateFeeResponse.Size(m)
}
func (m *CalculateFeeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CalculateFeeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CalculateFeeResponse proto.InternalMessageInfo

func (m *CalculateFeeResponse) GetAssetId() []byte {
	if m != nil {
		return m.AssetId
	}
	return nil
}

func (m *CalculateFeeResponse) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type SignRequest struct {
	Transaction          *Transaction `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	SignerPublicKey      []byte       `protobuf:"bytes,2,opt,name=signer_public_key,json=signerPublicKey,proto3" json:"signer_public_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *SignRequest) Reset()         { *m = SignRequest{} }
func (m *SignRequest) String() string { return proto.CompactTextString(m) }
func (*SignRequest) ProtoMessage()    {}
func (*SignRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_121a662cf7c9700a, []int{5}
}

func (m *SignRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignRequest.Unmarshal(m, b)
}
func (m *SignRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignRequest.Marshal(b, m, deterministic)
}
func (m *SignRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignRequest.Merge(m, src)
}
func (m *SignRequest) XXX_Size() int {
	return xxx_messageInfo_SignRequest.Size(m)
}
func (m *SignRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SignRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SignRequest proto.InternalMessageInfo

func (m *SignRequest) GetTransaction() *Transaction {
	if m != nil {
		return m.Transaction
	}
	return nil
}

func (m *SignRequest) GetSignerPublicKey() []byte {
	if m != nil {
		return m.SignerPublicKey
	}
	return nil
}

func init() {
	proto.RegisterEnum("waves.node.grpc.TransactionStatus_Status", TransactionStatus_Status_name, TransactionStatus_Status_value)
	proto.RegisterType((*TransactionStatus)(nil), "waves.node.grpc.TransactionStatus")
	proto.RegisterType((*TransactionResponse)(nil), "waves.node.grpc.TransactionResponse")
	proto.RegisterType((*TransactionsRequest)(nil), "waves.node.grpc.TransactionsRequest")
	proto.RegisterType((*TransactionsByIdRequest)(nil), "waves.node.grpc.TransactionsByIdRequest")
	proto.RegisterType((*CalculateFeeResponse)(nil), "waves.node.grpc.CalculateFeeResponse")
	proto.RegisterType((*SignRequest)(nil), "waves.node.grpc.SignRequest")
}

func init() { proto.RegisterFile("transactions_api.proto", fileDescriptor_121a662cf7c9700a) }

var fileDescriptor_121a662cf7c9700a = []byte{
	// 598 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xc5, 0x49, 0x15, 0xc8, 0xa4, 0xc4, 0xed, 0x82, 0x4a, 0x6a, 0x71, 0x88, 0x2c, 0x24, 0x02,
	0x07, 0xab, 0x0a, 0x1c, 0x80, 0x03, 0xa8, 0x29, 0x6d, 0x65, 0x21, 0x52, 0xe4, 0xa4, 0x02, 0x21,
	0x81, 0xb5, 0xb5, 0xa7, 0xe9, 0xaa, 0x89, 0xd7, 0xec, 0xae, 0x8b, 0xf2, 0x03, 0x88, 0x6f, 0xe1,
	0x0b, 0xf8, 0x3c, 0xe4, 0xcd, 0x26, 0xd9, 0x52, 0xa5, 0xf4, 0xc0, 0xc9, 0x9a, 0xd9, 0xb7, 0xf3,
	0xde, 0xbc, 0xd9, 0x31, 0x6c, 0x29, 0x41, 0x33, 0x49, 0x13, 0xc5, 0x78, 0x26, 0x63, 0x9a, 0xb3,
	0x20, 0x17, 0x5c, 0x71, 0xe2, 0x7e, 0xa7, 0x17, 0x28, 0x83, 0x8c, 0xa7, 0x18, 0x8c, 0x44, 0x9e,
	0x78, 0xae, 0xc0, 0x84, 0xe5, 0x0c, 0x33, 0x35, 0x43, 0x78, 0x9b, 0xd6, 0x4d, 0x93, 0xf2, 0x58,
	0x76, 0xc1, 0xcf, 0x31, 0x96, 0x89, 0x60, 0xb9, 0x8a, 0x05, 0xca, 0x62, 0x6c, 0xe0, 0xfe, 0x6f,
	0x07, 0x36, 0x87, 0xcb, 0x1b, 0x03, 0x45, 0x55, 0x21, 0x49, 0x13, 0x2a, 0x2c, 0x6d, 0x39, 0x6d,
	0xa7, 0xb3, 0x1e, 0x55, 0x58, 0x4a, 0x76, 0xa1, 0x26, 0xf5, 0x49, 0xab, 0xd2, 0x76, 0x3a, 0xcd,
	0xee, 0x93, 0xe0, 0x2f, 0x1d, 0xc1, 0x95, 0x1a, 0xc1, 0xec, 0x13, 0x99, 0x8b, 0x64, 0x0b, 0x6a,
	0x67, 0xc8, 0x46, 0x67, 0xaa, 0x55, 0x6d, 0x3b, 0x9d, 0x6a, 0x64, 0x22, 0xff, 0x05, 0xd4, 0x16,
	0xa4, 0xd0, 0x3f, 0x1a, 0xc6, 0xfb, 0x9f, 0xc2, 0xc1, 0x70, 0xb0, 0x71, 0x8b, 0xb8, 0xd0, 0x38,
	0xee, 0xef, 0x1d, 0xf5, 0x0f, 0xc2, 0xe8, 0xfd, 0xfe, 0xdb, 0x0d, 0x87, 0xdc, 0x85, 0xfa, 0x32,
	0xac, 0xf8, 0x53, 0xb8, 0x67, 0xb1, 0x46, 0x28, 0x73, 0x9e, 0x49, 0xbc, 0xa2, 0x7d, 0x49, 0x5c,
	0xb1, 0x89, 0xc9, 0x2b, 0x68, 0x58, 0x56, 0x69, 0x55, 0x8d, 0x6e, 0xcb, 0x34, 0x36, 0x60, 0xa3,
	0x0c, 0x53, 0xbb, 0xbc, 0x0d, 0xf6, 0x7f, 0x38, 0x97, 0xb8, 0x65, 0x84, 0xdf, 0x0a, 0x94, 0xaa,
	0xe4, 0x92, 0x98, 0xa5, 0x28, 0x0c, 0xbf, 0x89, 0x48, 0x00, 0xf5, 0xc5, 0x9c, 0xb4, 0x8c, 0x46,
	0x77, 0xc3, 0x30, 0x45, 0xf3, 0x7c, 0xb4, 0x84, 0x90, 0xc7, 0xe0, 0x5a, 0x74, 0x31, 0x4b, 0x65,
	0xab, 0xda, 0xae, 0x76, 0xd6, 0xa3, 0xa6, 0x95, 0x0e, 0x53, 0xe9, 0xf7, 0xe0, 0x81, 0xad, 0xa3,
	0x37, 0x0d, 0xd3, 0xb9, 0x96, 0x1b, 0xd7, 0x08, 0xe1, 0xfe, 0x1e, 0x1d, 0x27, 0xc5, 0x98, 0x2a,
	0x3c, 0x40, 0x5c, 0x18, 0xb9, 0x0d, 0x77, 0xa8, 0x94, 0xa8, 0xe2, 0x85, 0x9d, 0xb7, 0x75, 0x1c,
	0x6a, 0x4f, 0xe9, 0x84, 0x17, 0xa6, 0x99, 0xb5, 0xc8, 0x44, 0x3e, 0x87, 0x46, 0xe9, 0xdc, 0x5c,
	0xc2, 0xf3, 0xcb, 0x16, 0x3b, 0xba, 0x71, 0x62, 0x1a, 0x5f, 0x65, 0x2e, 0x79, 0x0a, 0x9b, 0xb2,
	0xb4, 0x5f, 0xc4, 0x79, 0x71, 0x32, 0x66, 0x49, 0x7c, 0x8e, 0x53, 0xcd, 0xb3, 0x1e, 0xb9, 0xb3,
	0x83, 0x0f, 0x3a, 0xff, 0x0e, 0xa7, 0xdd, 0x9f, 0x6b, 0xe0, 0xda, 0x06, 0xec, 0xe6, 0x8c, 0xc4,
	0xe0, 0x1e, 0xa2, 0xb2, 0xb3, 0xe4, 0xd1, 0x75, 0xef, 0x75, 0x3e, 0x3d, 0xef, 0x5a, 0xd4, 0xdc,
	0x96, 0x1d, 0x87, 0x0c, 0x35, 0x41, 0xf9, 0x6a, 0x71, 0xef, 0x8c, 0x66, 0x23, 0xbc, 0x29, 0xc1,
	0xb6, 0x41, 0x85, 0x7a, 0x1f, 0x07, 0x7a, 0x1d, 0x23, 0xbd, 0x8d, 0x3b, 0x0e, 0xf9, 0x02, 0x0d,
	0x53, 0xb5, 0x90, 0x28, 0x49, 0xe7, 0xda, 0x8a, 0xd6, 0xa0, 0x3d, 0xff, 0xdf, 0xcb, 0xb8, 0xe3,
	0x90, 0xaf, 0xd0, 0x3c, 0x44, 0x75, 0x9c, 0x25, 0x3c, 0x3b, 0x65, 0x62, 0x82, 0xe9, 0x7f, 0x36,
	0xe5, 0x35, 0xac, 0x95, 0xa3, 0x27, 0x0f, 0xaf, 0xe0, 0xad, 0x17, 0xe1, 0xad, 0xdc, 0x2f, 0xf2,
	0x06, 0xea, 0x3d, 0xc1, 0x69, 0x9a, 0x50, 0xa9, 0xc8, 0x4a, 0xd8, 0xea, 0x02, 0xbd, 0x97, 0xe0,
	0x25, 0x7c, 0x32, 0x3b, 0xce, 0xc7, 0x54, 0x9d, 0x72, 0x31, 0x09, 0xca, 0x3f, 0x67, 0x29, 0xe2,
	0x73, 0x7d, 0x84, 0x19, 0x0a, 0xaa, 0x30, 0xfd, 0x55, 0x71, 0x3f, 0xea, 0x12, 0xfd, 0x52, 0xe1,
	0xa1, 0xc8, 0x93, 0x93, 0x9a, 0xfe, 0x17, 0x3e, 0xfb, 0x13, 0x00, 0x00, 0xff, 0xff, 0x54, 0x74,
	0x37, 0x55, 0x76, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TransactionsApiClient is the client API for TransactionsApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TransactionsApiClient interface {
	GetTransactions(ctx context.Context, in *TransactionsRequest, opts ...grpc.CallOption) (TransactionsApi_GetTransactionsClient, error)
	GetStateChanges(ctx context.Context, in *TransactionsRequest, opts ...grpc.CallOption) (TransactionsApi_GetStateChangesClient, error)
	GetStatuses(ctx context.Context, in *TransactionsByIdRequest, opts ...grpc.CallOption) (TransactionsApi_GetStatusesClient, error)
	GetUnconfirmed(ctx context.Context, in *TransactionsRequest, opts ...grpc.CallOption) (TransactionsApi_GetUnconfirmedClient, error)
	Sign(ctx context.Context, in *SignRequest, opts ...grpc.CallOption) (*SignedTransaction, error)
	Broadcast(ctx context.Context, in *SignedTransaction, opts ...grpc.CallOption) (*SignedTransaction, error)
}

type transactionsApiClient struct {
	cc *grpc.ClientConn
}

func NewTransactionsApiClient(cc *grpc.ClientConn) TransactionsApiClient {
	return &transactionsApiClient{cc}
}

func (c *transactionsApiClient) GetTransactions(ctx context.Context, in *TransactionsRequest, opts ...grpc.CallOption) (TransactionsApi_GetTransactionsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TransactionsApi_serviceDesc.Streams[0], "/waves.node.grpc.TransactionsApi/GetTransactions", opts...)
	if err != nil {
		return nil, err
	}
	x := &transactionsApiGetTransactionsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TransactionsApi_GetTransactionsClient interface {
	Recv() (*TransactionResponse, error)
	grpc.ClientStream
}

type transactionsApiGetTransactionsClient struct {
	grpc.ClientStream
}

func (x *transactionsApiGetTransactionsClient) Recv() (*TransactionResponse, error) {
	m := new(TransactionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *transactionsApiClient) GetStateChanges(ctx context.Context, in *TransactionsRequest, opts ...grpc.CallOption) (TransactionsApi_GetStateChangesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TransactionsApi_serviceDesc.Streams[1], "/waves.node.grpc.TransactionsApi/GetStateChanges", opts...)
	if err != nil {
		return nil, err
	}
	x := &transactionsApiGetStateChangesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TransactionsApi_GetStateChangesClient interface {
	Recv() (*InvokeScriptResult, error)
	grpc.ClientStream
}

type transactionsApiGetStateChangesClient struct {
	grpc.ClientStream
}

func (x *transactionsApiGetStateChangesClient) Recv() (*InvokeScriptResult, error) {
	m := new(InvokeScriptResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *transactionsApiClient) GetStatuses(ctx context.Context, in *TransactionsByIdRequest, opts ...grpc.CallOption) (TransactionsApi_GetStatusesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TransactionsApi_serviceDesc.Streams[2], "/waves.node.grpc.TransactionsApi/GetStatuses", opts...)
	if err != nil {
		return nil, err
	}
	x := &transactionsApiGetStatusesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TransactionsApi_GetStatusesClient interface {
	Recv() (*TransactionStatus, error)
	grpc.ClientStream
}

type transactionsApiGetStatusesClient struct {
	grpc.ClientStream
}

func (x *transactionsApiGetStatusesClient) Recv() (*TransactionStatus, error) {
	m := new(TransactionStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *transactionsApiClient) GetUnconfirmed(ctx context.Context, in *TransactionsRequest, opts ...grpc.CallOption) (TransactionsApi_GetUnconfirmedClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TransactionsApi_serviceDesc.Streams[3], "/waves.node.grpc.TransactionsApi/GetUnconfirmed", opts...)
	if err != nil {
		return nil, err
	}
	x := &transactionsApiGetUnconfirmedClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TransactionsApi_GetUnconfirmedClient interface {
	Recv() (*TransactionResponse, error)
	grpc.ClientStream
}

type transactionsApiGetUnconfirmedClient struct {
	grpc.ClientStream
}

func (x *transactionsApiGetUnconfirmedClient) Recv() (*TransactionResponse, error) {
	m := new(TransactionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *transactionsApiClient) Sign(ctx context.Context, in *SignRequest, opts ...grpc.CallOption) (*SignedTransaction, error) {
	out := new(SignedTransaction)
	err := c.cc.Invoke(ctx, "/waves.node.grpc.TransactionsApi/Sign", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionsApiClient) Broadcast(ctx context.Context, in *SignedTransaction, opts ...grpc.CallOption) (*SignedTransaction, error) {
	out := new(SignedTransaction)
	err := c.cc.Invoke(ctx, "/waves.node.grpc.TransactionsApi/Broadcast", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionsApiServer is the server API for TransactionsApi service.
type TransactionsApiServer interface {
	GetTransactions(*TransactionsRequest, TransactionsApi_GetTransactionsServer) error
	GetStateChanges(*TransactionsRequest, TransactionsApi_GetStateChangesServer) error
	GetStatuses(*TransactionsByIdRequest, TransactionsApi_GetStatusesServer) error
	GetUnconfirmed(*TransactionsRequest, TransactionsApi_GetUnconfirmedServer) error
	Sign(context.Context, *SignRequest) (*SignedTransaction, error)
	Broadcast(context.Context, *SignedTransaction) (*SignedTransaction, error)
}

// UnimplementedTransactionsApiServer can be embedded to have forward compatible implementations.
type UnimplementedTransactionsApiServer struct {
}

func (*UnimplementedTransactionsApiServer) GetTransactions(req *TransactionsRequest, srv TransactionsApi_GetTransactionsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetTransactions not implemented")
}
func (*UnimplementedTransactionsApiServer) GetStateChanges(req *TransactionsRequest, srv TransactionsApi_GetStateChangesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStateChanges not implemented")
}
func (*UnimplementedTransactionsApiServer) GetStatuses(req *TransactionsByIdRequest, srv TransactionsApi_GetStatusesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStatuses not implemented")
}
func (*UnimplementedTransactionsApiServer) GetUnconfirmed(req *TransactionsRequest, srv TransactionsApi_GetUnconfirmedServer) error {
	return status.Errorf(codes.Unimplemented, "method GetUnconfirmed not implemented")
}
func (*UnimplementedTransactionsApiServer) Sign(ctx context.Context, req *SignRequest) (*SignedTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sign not implemented")
}
func (*UnimplementedTransactionsApiServer) Broadcast(ctx context.Context, req *SignedTransaction) (*SignedTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Broadcast not implemented")
}

func RegisterTransactionsApiServer(s *grpc.Server, srv TransactionsApiServer) {
	s.RegisterService(&_TransactionsApi_serviceDesc, srv)
}

func _TransactionsApi_GetTransactions_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TransactionsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TransactionsApiServer).GetTransactions(m, &transactionsApiGetTransactionsServer{stream})
}

type TransactionsApi_GetTransactionsServer interface {
	Send(*TransactionResponse) error
	grpc.ServerStream
}

type transactionsApiGetTransactionsServer struct {
	grpc.ServerStream
}

func (x *transactionsApiGetTransactionsServer) Send(m *TransactionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _TransactionsApi_GetStateChanges_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TransactionsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TransactionsApiServer).GetStateChanges(m, &transactionsApiGetStateChangesServer{stream})
}

type TransactionsApi_GetStateChangesServer interface {
	Send(*InvokeScriptResult) error
	grpc.ServerStream
}

type transactionsApiGetStateChangesServer struct {
	grpc.ServerStream
}

func (x *transactionsApiGetStateChangesServer) Send(m *InvokeScriptResult) error {
	return x.ServerStream.SendMsg(m)
}

func _TransactionsApi_GetStatuses_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TransactionsByIdRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TransactionsApiServer).GetStatuses(m, &transactionsApiGetStatusesServer{stream})
}

type TransactionsApi_GetStatusesServer interface {
	Send(*TransactionStatus) error
	grpc.ServerStream
}

type transactionsApiGetStatusesServer struct {
	grpc.ServerStream
}

func (x *transactionsApiGetStatusesServer) Send(m *TransactionStatus) error {
	return x.ServerStream.SendMsg(m)
}

func _TransactionsApi_GetUnconfirmed_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TransactionsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TransactionsApiServer).GetUnconfirmed(m, &transactionsApiGetUnconfirmedServer{stream})
}

type TransactionsApi_GetUnconfirmedServer interface {
	Send(*TransactionResponse) error
	grpc.ServerStream
}

type transactionsApiGetUnconfirmedServer struct {
	grpc.ServerStream
}

func (x *transactionsApiGetUnconfirmedServer) Send(m *TransactionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _TransactionsApi_Sign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionsApiServer).Sign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/waves.node.grpc.TransactionsApi/Sign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionsApiServer).Sign(ctx, req.(*SignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionsApi_Broadcast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignedTransaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionsApiServer).Broadcast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/waves.node.grpc.TransactionsApi/Broadcast",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionsApiServer).Broadcast(ctx, req.(*SignedTransaction))
	}
	return interceptor(ctx, in, info, handler)
}

var _TransactionsApi_serviceDesc = grpc.ServiceDesc{
	ServiceName: "waves.node.grpc.TransactionsApi",
	HandlerType: (*TransactionsApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sign",
			Handler:    _TransactionsApi_Sign_Handler,
		},
		{
			MethodName: "Broadcast",
			Handler:    _TransactionsApi_Broadcast_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetTransactions",
			Handler:       _TransactionsApi_GetTransactions_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetStateChanges",
			Handler:       _TransactionsApi_GetStateChanges_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetStatuses",
			Handler:       _TransactionsApi_GetStatuses_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetUnconfirmed",
			Handler:       _TransactionsApi_GetUnconfirmed_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "transactions_api.proto",
}
