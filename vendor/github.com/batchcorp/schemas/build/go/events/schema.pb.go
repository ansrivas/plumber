// Code generated by protoc-gen-go. DO NOT EDIT.
// source: schema.proto

package events

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Schema_Type int32

const (
	Schema_UNKNOWN  Schema_Type = 0
	Schema_PLAIN    Schema_Type = 1
	Schema_JSON     Schema_Type = 2
	Schema_PROTOBUF Schema_Type = 3
)

var Schema_Type_name = map[int32]string{
	0: "UNKNOWN",
	1: "PLAIN",
	2: "JSON",
	3: "PROTOBUF",
}

var Schema_Type_value = map[string]int32{
	"UNKNOWN":  0,
	"PLAIN":    1,
	"JSON":     2,
	"PROTOBUF": 3,
}

func (x Schema_Type) String() string {
	return proto.EnumName(Schema_Type_name, int32(x))
}

func (Schema_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1c5fb4d8cc22d66a, []int{0, 0}
}

type Schema_UpdateType int32

const (
	Schema_INITIAL  Schema_UpdateType = 0
	Schema_EXISTING Schema_UpdateType = 1
)

var Schema_UpdateType_name = map[int32]string{
	0: "INITIAL",
	1: "EXISTING",
}

var Schema_UpdateType_value = map[string]int32{
	"INITIAL":  0,
	"EXISTING": 1,
}

func (x Schema_UpdateType) String() string {
	return proto.EnumName(Schema_UpdateType_name, int32(x))
}

func (Schema_UpdateType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1c5fb4d8cc22d66a, []int{0, 1}
}

type Schema struct {
	// The collector will ONLY fill out the 'id' for incoming messages - it is
	// the responsibility of downstream consumers to lookup the corresponding
	// schema configuration by 'id'.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Indicates the "data type" - what format are the collectors expecting to
	// receive the events in?
	Type Schema_Type       `protobuf:"varint,2,opt,name=type,proto3,enum=events.Schema_Type" json:"type,omitempty"`
	Raw  map[string][]byte `protobuf:"bytes,3,rep,name=raw,proto3" json:"raw,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Only used when Type == PROTOBUF
	ProtobufMessageName       string `protobuf:"bytes,4,opt,name=protobuf_message_name,json=protobufMessageName,proto3" json:"protobuf_message_name,omitempty"`
	ProtobufFileDescriptorSet []byte `protobuf:"bytes,5,opt,name=protobuf_file_descriptor_set,json=protobufFileDescriptorSet,proto3" json:"protobuf_file_descriptor_set,omitempty"`
	// The following fields are used by the schema-manager to facilitate schema updates
	UpdateType          Schema_UpdateType `protobuf:"varint,6,opt,name=update_type,json=updateType,proto3,enum=events.Schema_UpdateType" json:"update_type,omitempty"`
	UpdateCollectToken  string            `protobuf:"bytes,7,opt,name=update_collect_token,json=updateCollectToken,proto3" json:"update_collect_token,omitempty"`
	UpdateParquetSchema []byte            `protobuf:"bytes,8,opt,name=update_parquet_schema,json=updateParquetSchema,proto3" json:"update_parquet_schema,omitempty"`
	UpdateSqlSchema     []byte            `protobuf:"bytes,9,opt,name=update_sql_schema,json=updateSqlSchema,proto3" json:"update_sql_schema,omitempty"`
	UpdateFingerprint   string            `protobuf:"bytes,10,opt,name=update_fingerprint,json=updateFingerprint,proto3" json:"update_fingerprint,omitempty"`
	UpdateCollectId     string            `protobuf:"bytes,11,opt,name=update_collect_id,json=updateCollectId,proto3" json:"update_collect_id,omitempty"`
	// Schema version is used to create unique collect-update-* topics which in
	// turn allow the writer to write data using correct schema when there are
	// more than 1 schema updates in-flight. Talk to MG or DS.
	//
	// This field is incremented by the collectors on a schema update.
	SchemaVersion int64 `protobuf:"varint,12,opt,name=schema_version,json=schemaVersion,proto3" json:"schema_version,omitempty"`
	// The manifest message payload, as JSON, to infer the schema from.
	// Only used when Message.Type == GENERATE_SCHEMA
	SourcePayload        []byte   `protobuf:"bytes,13,opt,name=source_payload,json=sourcePayload,proto3" json:"source_payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Schema) Reset()         { *m = Schema{} }
func (m *Schema) String() string { return proto.CompactTextString(m) }
func (*Schema) ProtoMessage()    {}
func (*Schema) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c5fb4d8cc22d66a, []int{0}
}

func (m *Schema) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Schema.Unmarshal(m, b)
}
func (m *Schema) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Schema.Marshal(b, m, deterministic)
}
func (m *Schema) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Schema.Merge(m, src)
}
func (m *Schema) XXX_Size() int {
	return xxx_messageInfo_Schema.Size(m)
}
func (m *Schema) XXX_DiscardUnknown() {
	xxx_messageInfo_Schema.DiscardUnknown(m)
}

var xxx_messageInfo_Schema proto.InternalMessageInfo

func (m *Schema) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Schema) GetType() Schema_Type {
	if m != nil {
		return m.Type
	}
	return Schema_UNKNOWN
}

func (m *Schema) GetRaw() map[string][]byte {
	if m != nil {
		return m.Raw
	}
	return nil
}

func (m *Schema) GetProtobufMessageName() string {
	if m != nil {
		return m.ProtobufMessageName
	}
	return ""
}

func (m *Schema) GetProtobufFileDescriptorSet() []byte {
	if m != nil {
		return m.ProtobufFileDescriptorSet
	}
	return nil
}

func (m *Schema) GetUpdateType() Schema_UpdateType {
	if m != nil {
		return m.UpdateType
	}
	return Schema_INITIAL
}

func (m *Schema) GetUpdateCollectToken() string {
	if m != nil {
		return m.UpdateCollectToken
	}
	return ""
}

func (m *Schema) GetUpdateParquetSchema() []byte {
	if m != nil {
		return m.UpdateParquetSchema
	}
	return nil
}

func (m *Schema) GetUpdateSqlSchema() []byte {
	if m != nil {
		return m.UpdateSqlSchema
	}
	return nil
}

func (m *Schema) GetUpdateFingerprint() string {
	if m != nil {
		return m.UpdateFingerprint
	}
	return ""
}

func (m *Schema) GetUpdateCollectId() string {
	if m != nil {
		return m.UpdateCollectId
	}
	return ""
}

func (m *Schema) GetSchemaVersion() int64 {
	if m != nil {
		return m.SchemaVersion
	}
	return 0
}

func (m *Schema) GetSourcePayload() []byte {
	if m != nil {
		return m.SourcePayload
	}
	return nil
}

func init() {
	proto.RegisterEnum("events.Schema_Type", Schema_Type_name, Schema_Type_value)
	proto.RegisterEnum("events.Schema_UpdateType", Schema_UpdateType_name, Schema_UpdateType_value)
	proto.RegisterType((*Schema)(nil), "events.Schema")
	proto.RegisterMapType((map[string][]byte)(nil), "events.Schema.RawEntry")
}

func init() { proto.RegisterFile("schema.proto", fileDescriptor_1c5fb4d8cc22d66a) }

var fileDescriptor_1c5fb4d8cc22d66a = []byte{
	// 526 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0xe1, 0x6f, 0xd2, 0x40,
	0x18, 0xc6, 0x57, 0x0a, 0x8c, 0xbd, 0x30, 0x56, 0x8f, 0x19, 0x6f, 0xc6, 0x0f, 0x84, 0xc4, 0x0c,
	0x8d, 0x16, 0x83, 0xc9, 0x62, 0xf6, 0xc5, 0x6c, 0x3a, 0x4c, 0x75, 0x16, 0x52, 0x40, 0x8d, 0x5f,
	0x9a, 0xa3, 0x3d, 0xe0, 0xb2, 0xd2, 0x2b, 0xd7, 0x2b, 0x0b, 0x7f, 0xb6, 0xff, 0x81, 0xe9, 0x5d,
	0x61, 0xba, 0x6f, 0xed, 0xfb, 0xfc, 0xde, 0x7b, 0x9f, 0xe7, 0xee, 0x85, 0x46, 0x1a, 0x2c, 0xe9,
	0x8a, 0xd8, 0x89, 0xe0, 0x92, 0xa3, 0x2a, 0xdd, 0xd0, 0x58, 0xa6, 0x9d, 0x3f, 0x15, 0xa8, 0x8e,
	0x95, 0x80, 0x9a, 0x50, 0x62, 0x21, 0x36, 0xda, 0x46, 0xf7, 0xc8, 0x2b, 0xb1, 0x10, 0x9d, 0x43,
	0x59, 0x6e, 0x13, 0x8a, 0x4b, 0x6d, 0xa3, 0xdb, 0xec, 0xb7, 0x6c, 0xdd, 0x61, 0x6b, 0xda, 0x9e,
	0x6c, 0x13, 0xea, 0x29, 0x00, 0xbd, 0x02, 0x53, 0x90, 0x7b, 0x6c, 0xb6, 0xcd, 0x6e, 0xbd, 0xff,
	0xec, 0x11, 0xe7, 0x91, 0xfb, 0x9b, 0x58, 0x8a, 0xad, 0x97, 0x33, 0xa8, 0x0f, 0x4f, 0xd5, 0xfc,
	0x59, 0x36, 0xf7, 0x57, 0x34, 0x4d, 0xc9, 0x82, 0xfa, 0x31, 0x59, 0x51, 0x5c, 0x56, 0x63, 0x5b,
	0x3b, 0xf1, 0xbb, 0xd6, 0x5c, 0xb2, 0xa2, 0xe8, 0x23, 0xbc, 0xd8, 0xf7, 0xcc, 0x59, 0x44, 0xfd,
	0x90, 0xa6, 0x81, 0x60, 0x89, 0xe4, 0xc2, 0x4f, 0xa9, 0xc4, 0x95, 0xb6, 0xd1, 0x6d, 0x78, 0x67,
	0x3b, 0x66, 0xc0, 0x22, 0xfa, 0x79, 0x4f, 0x8c, 0xa9, 0x44, 0x97, 0x50, 0xcf, 0x92, 0x90, 0x48,
	0xea, 0xab, 0x3c, 0x55, 0x95, 0xe7, 0xec, 0x91, 0xcf, 0xa9, 0x22, 0x54, 0x2a, 0xc8, 0xf6, 0xdf,
	0xe8, 0x1d, 0x9c, 0x16, 0xbd, 0x01, 0x8f, 0x22, 0x1a, 0x48, 0x5f, 0xf2, 0x3b, 0x1a, 0xe3, 0x43,
	0xe5, 0x17, 0x69, 0xed, 0x93, 0x96, 0x26, 0xb9, 0x92, 0x47, 0x2c, 0x3a, 0x12, 0x22, 0xd6, 0x19,
	0x95, 0xbe, 0xbe, 0x78, 0x5c, 0x53, 0x3e, 0x5b, 0x5a, 0x1c, 0x69, 0xad, 0xb8, 0xfa, 0xd7, 0xf0,
	0xa4, 0xe8, 0x49, 0xd7, 0xd1, 0x8e, 0x3f, 0x52, 0xfc, 0x89, 0x16, 0xc6, 0xeb, 0xa8, 0x60, 0xdf,
	0x42, 0x31, 0xd5, 0x9f, 0xb3, 0x78, 0x41, 0x45, 0x22, 0x58, 0x2c, 0x31, 0x28, 0x3f, 0xc5, 0x29,
	0x83, 0x07, 0xe1, 0x9f, 0xa3, 0x77, 0x01, 0x58, 0x88, 0xeb, 0x8a, 0x3e, 0xf9, 0xcf, 0xbd, 0x13,
	0xa2, 0x97, 0xd0, 0xd4, 0xb3, 0xfd, 0x0d, 0x15, 0x29, 0xe3, 0x31, 0x6e, 0xb4, 0x8d, 0xae, 0xe9,
	0x1d, 0xeb, 0xea, 0x0f, 0x5d, 0x54, 0x18, 0xcf, 0x44, 0x90, 0x27, 0xdc, 0x46, 0x9c, 0x84, 0xf8,
	0x58, 0x59, 0x3d, 0xd6, 0xd5, 0x91, 0x2e, 0x3e, 0xbf, 0x80, 0xda, 0xee, 0xf1, 0x91, 0x05, 0xe6,
	0x1d, 0xdd, 0x16, 0xcb, 0x95, 0x7f, 0xa2, 0x53, 0xa8, 0x6c, 0x48, 0x94, 0xe9, 0xf5, 0x6a, 0x78,
	0xfa, 0xe7, 0xb2, 0xf4, 0xc1, 0xe8, 0x5c, 0x40, 0x59, 0x5d, 0x7d, 0x1d, 0x0e, 0xa7, 0xee, 0x37,
	0x77, 0xf8, 0xd3, 0xb5, 0x0e, 0xd0, 0x11, 0x54, 0x46, 0xb7, 0x57, 0x8e, 0x6b, 0x19, 0xa8, 0x06,
	0xe5, 0xaf, 0xe3, 0xa1, 0x6b, 0x95, 0x50, 0x03, 0x6a, 0x23, 0x6f, 0x38, 0x19, 0x5e, 0x4f, 0x07,
	0x96, 0xd9, 0x39, 0x07, 0x78, 0x78, 0xc4, 0xbc, 0xdb, 0x71, 0x9d, 0x89, 0x73, 0x75, 0x6b, 0x1d,
	0xe4, 0xe0, 0xcd, 0x2f, 0x67, 0x3c, 0x71, 0xdc, 0x2f, 0x96, 0x71, 0x6d, 0xff, 0x7e, 0xb3, 0x60,
	0x72, 0x99, 0xcd, 0xec, 0x80, 0xaf, 0x7a, 0x33, 0x22, 0x83, 0x65, 0xc0, 0x45, 0xd2, 0xd3, 0x29,
	0xd3, 0xde, 0x2c, 0x63, 0x51, 0xd8, 0x5b, 0xf0, 0x9e, 0xde, 0x90, 0x59, 0x55, 0xad, 0xd6, 0xfb,
	0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x21, 0x9e, 0xc3, 0xad, 0x42, 0x03, 0x00, 0x00,
}
