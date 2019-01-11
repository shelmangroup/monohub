// Code generated by protoc-gen-go. DO NOT EDIT.
// source: commit.proto

package monohub

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Parent struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Sha                  string   `protobuf:"bytes,2,opt,name=sha" json:"sha,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Parent) Reset()         { *m = Parent{} }
func (m *Parent) String() string { return proto.CompactTextString(m) }
func (*Parent) ProtoMessage()    {}
func (*Parent) Descriptor() ([]byte, []int) {
	return fileDescriptor_commit_b0b64943c0218a98, []int{0}
}
func (m *Parent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Parent.Unmarshal(m, b)
}
func (m *Parent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Parent.Marshal(b, m, deterministic)
}
func (dst *Parent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Parent.Merge(dst, src)
}
func (m *Parent) XXX_Size() int {
	return xxx_messageInfo_Parent.Size(m)
}
func (m *Parent) XXX_DiscardUnknown() {
	xxx_messageInfo_Parent.DiscardUnknown(m)
}

var xxx_messageInfo_Parent proto.InternalMessageInfo

func (m *Parent) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Parent) GetSha() string {
	if m != nil {
		return m.Sha
	}
	return ""
}

type Stats struct {
	Additions            int64    `protobuf:"varint,1,opt,name=additions" json:"additions,omitempty"`
	Deletions            int64    `protobuf:"varint,2,opt,name=deletions" json:"deletions,omitempty"`
	Total                int64    `protobuf:"varint,3,opt,name=total" json:"total,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stats) Reset()         { *m = Stats{} }
func (m *Stats) String() string { return proto.CompactTextString(m) }
func (*Stats) ProtoMessage()    {}
func (*Stats) Descriptor() ([]byte, []int) {
	return fileDescriptor_commit_b0b64943c0218a98, []int{1}
}
func (m *Stats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stats.Unmarshal(m, b)
}
func (m *Stats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stats.Marshal(b, m, deterministic)
}
func (dst *Stats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stats.Merge(dst, src)
}
func (m *Stats) XXX_Size() int {
	return xxx_messageInfo_Stats.Size(m)
}
func (m *Stats) XXX_DiscardUnknown() {
	xxx_messageInfo_Stats.DiscardUnknown(m)
}

var xxx_messageInfo_Stats proto.InternalMessageInfo

func (m *Stats) GetAdditions() int64 {
	if m != nil {
		return m.Additions
	}
	return 0
}

func (m *Stats) GetDeletions() int64 {
	if m != nil {
		return m.Deletions
	}
	return 0
}

func (m *Stats) GetTotal() int64 {
	if m != nil {
		return m.Total
	}
	return 0
}

type Author struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Login                string   `protobuf:"bytes,2,opt,name=login" json:"login,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	Email                string   `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Author) Reset()         { *m = Author{} }
func (m *Author) String() string { return proto.CompactTextString(m) }
func (*Author) ProtoMessage()    {}
func (*Author) Descriptor() ([]byte, []int) {
	return fileDescriptor_commit_b0b64943c0218a98, []int{2}
}
func (m *Author) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Author.Unmarshal(m, b)
}
func (m *Author) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Author.Marshal(b, m, deterministic)
}
func (dst *Author) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Author.Merge(dst, src)
}
func (m *Author) XXX_Size() int {
	return xxx_messageInfo_Author.Size(m)
}
func (m *Author) XXX_DiscardUnknown() {
	xxx_messageInfo_Author.DiscardUnknown(m)
}

var xxx_messageInfo_Author proto.InternalMessageInfo

func (m *Author) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Author) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *Author) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Author) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type Tree struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Sha                  string   `protobuf:"bytes,2,opt,name=sha" json:"sha,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tree) Reset()         { *m = Tree{} }
func (m *Tree) String() string { return proto.CompactTextString(m) }
func (*Tree) ProtoMessage()    {}
func (*Tree) Descriptor() ([]byte, []int) {
	return fileDescriptor_commit_b0b64943c0218a98, []int{3}
}
func (m *Tree) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tree.Unmarshal(m, b)
}
func (m *Tree) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tree.Marshal(b, m, deterministic)
}
func (dst *Tree) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tree.Merge(dst, src)
}
func (m *Tree) XXX_Size() int {
	return xxx_messageInfo_Tree.Size(m)
}
func (m *Tree) XXX_DiscardUnknown() {
	xxx_messageInfo_Tree.DiscardUnknown(m)
}

var xxx_messageInfo_Tree proto.InternalMessageInfo

func (m *Tree) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Tree) GetSha() string {
	if m != nil {
		return m.Sha
	}
	return ""
}

type Commit struct {
	Committer            *Author  `protobuf:"bytes,2,opt,name=committer" json:"committer,omitempty"`
	Message              string   `protobuf:"bytes,3,opt,name=message" json:"message,omitempty"`
	Tree                 *Tree    `protobuf:"bytes,4,opt,name=tree" json:"tree,omitempty"`
	Author               *Author  `protobuf:"bytes,5,opt,name=author" json:"author,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Commit) Reset()         { *m = Commit{} }
func (m *Commit) String() string { return proto.CompactTextString(m) }
func (*Commit) ProtoMessage()    {}
func (*Commit) Descriptor() ([]byte, []int) {
	return fileDescriptor_commit_b0b64943c0218a98, []int{4}
}
func (m *Commit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Commit.Unmarshal(m, b)
}
func (m *Commit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Commit.Marshal(b, m, deterministic)
}
func (dst *Commit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Commit.Merge(dst, src)
}
func (m *Commit) XXX_Size() int {
	return xxx_messageInfo_Commit.Size(m)
}
func (m *Commit) XXX_DiscardUnknown() {
	xxx_messageInfo_Commit.DiscardUnknown(m)
}

var xxx_messageInfo_Commit proto.InternalMessageInfo

func (m *Commit) GetCommitter() *Author {
	if m != nil {
		return m.Committer
	}
	return nil
}

func (m *Commit) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Commit) GetTree() *Tree {
	if m != nil {
		return m.Tree
	}
	return nil
}

func (m *Commit) GetAuthor() *Author {
	if m != nil {
		return m.Author
	}
	return nil
}

type File struct {
	Filename             string   `protobuf:"bytes,1,opt,name=filename" json:"filename,omitempty"`
	Additions            int64    `protobuf:"varint,2,opt,name=additions" json:"additions,omitempty"`
	Changes              int64    `protobuf:"varint,3,opt,name=changes" json:"changes,omitempty"`
	Deletions            int64    `protobuf:"varint,4,opt,name=deletions" json:"deletions,omitempty"`
	Status               string   `protobuf:"bytes,5,opt,name=status" json:"status,omitempty"`
	BlobUrl              string   `protobuf:"bytes,6,opt,name=blob_url,json=blobUrl" json:"blob_url,omitempty"`
	Patch                string   `protobuf:"bytes,7,opt,name=patch" json:"patch,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *File) Reset()         { *m = File{} }
func (m *File) String() string { return proto.CompactTextString(m) }
func (*File) ProtoMessage()    {}
func (*File) Descriptor() ([]byte, []int) {
	return fileDescriptor_commit_b0b64943c0218a98, []int{5}
}
func (m *File) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_File.Unmarshal(m, b)
}
func (m *File) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_File.Marshal(b, m, deterministic)
}
func (dst *File) XXX_Merge(src proto.Message) {
	xxx_messageInfo_File.Merge(dst, src)
}
func (m *File) XXX_Size() int {
	return xxx_messageInfo_File.Size(m)
}
func (m *File) XXX_DiscardUnknown() {
	xxx_messageInfo_File.DiscardUnknown(m)
}

var xxx_messageInfo_File proto.InternalMessageInfo

func (m *File) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *File) GetAdditions() int64 {
	if m != nil {
		return m.Additions
	}
	return 0
}

func (m *File) GetChanges() int64 {
	if m != nil {
		return m.Changes
	}
	return 0
}

func (m *File) GetDeletions() int64 {
	if m != nil {
		return m.Deletions
	}
	return 0
}

func (m *File) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *File) GetBlobUrl() string {
	if m != nil {
		return m.BlobUrl
	}
	return ""
}

func (m *File) GetPatch() string {
	if m != nil {
		return m.Patch
	}
	return ""
}

type CommitRequest struct {
	Sha                  string   `protobuf:"bytes,1,opt,name=sha" json:"sha,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommitRequest) Reset()         { *m = CommitRequest{} }
func (m *CommitRequest) String() string { return proto.CompactTextString(m) }
func (*CommitRequest) ProtoMessage()    {}
func (*CommitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_commit_b0b64943c0218a98, []int{6}
}
func (m *CommitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommitRequest.Unmarshal(m, b)
}
func (m *CommitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommitRequest.Marshal(b, m, deterministic)
}
func (dst *CommitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitRequest.Merge(dst, src)
}
func (m *CommitRequest) XXX_Size() int {
	return xxx_messageInfo_CommitRequest.Size(m)
}
func (m *CommitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CommitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CommitRequest proto.InternalMessageInfo

func (m *CommitRequest) GetSha() string {
	if m != nil {
		return m.Sha
	}
	return ""
}

type CommitResponse struct {
	Sha                  string    `protobuf:"bytes,1,opt,name=sha" json:"sha,omitempty"`
	Author               *Author   `protobuf:"bytes,2,opt,name=author" json:"author,omitempty"`
	Files                []*File   `protobuf:"bytes,3,rep,name=files" json:"files,omitempty"`
	Stats                *Stats    `protobuf:"bytes,4,opt,name=stats" json:"stats,omitempty"`
	Parents              []*Parent `protobuf:"bytes,5,rep,name=parents" json:"parents,omitempty"`
	Commit               *Commit   `protobuf:"bytes,6,opt,name=commit" json:"commit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CommitResponse) Reset()         { *m = CommitResponse{} }
func (m *CommitResponse) String() string { return proto.CompactTextString(m) }
func (*CommitResponse) ProtoMessage()    {}
func (*CommitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_commit_b0b64943c0218a98, []int{7}
}
func (m *CommitResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommitResponse.Unmarshal(m, b)
}
func (m *CommitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommitResponse.Marshal(b, m, deterministic)
}
func (dst *CommitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitResponse.Merge(dst, src)
}
func (m *CommitResponse) XXX_Size() int {
	return xxx_messageInfo_CommitResponse.Size(m)
}
func (m *CommitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CommitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CommitResponse proto.InternalMessageInfo

func (m *CommitResponse) GetSha() string {
	if m != nil {
		return m.Sha
	}
	return ""
}

func (m *CommitResponse) GetAuthor() *Author {
	if m != nil {
		return m.Author
	}
	return nil
}

func (m *CommitResponse) GetFiles() []*File {
	if m != nil {
		return m.Files
	}
	return nil
}

func (m *CommitResponse) GetStats() *Stats {
	if m != nil {
		return m.Stats
	}
	return nil
}

func (m *CommitResponse) GetParents() []*Parent {
	if m != nil {
		return m.Parents
	}
	return nil
}

func (m *CommitResponse) GetCommit() *Commit {
	if m != nil {
		return m.Commit
	}
	return nil
}

func init() {
	proto.RegisterType((*Parent)(nil), "monohub.Parent")
	proto.RegisterType((*Stats)(nil), "monohub.Stats")
	proto.RegisterType((*Author)(nil), "monohub.Author")
	proto.RegisterType((*Tree)(nil), "monohub.Tree")
	proto.RegisterType((*Commit)(nil), "monohub.Commit")
	proto.RegisterType((*File)(nil), "monohub.File")
	proto.RegisterType((*CommitRequest)(nil), "monohub.CommitRequest")
	proto.RegisterType((*CommitResponse)(nil), "monohub.CommitResponse")
}

func init() { proto.RegisterFile("commit.proto", fileDescriptor_commit_b0b64943c0218a98) }

var fileDescriptor_commit_b0b64943c0218a98 = []byte{
	// 447 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x4b, 0x8e, 0xd4, 0x30,
	0x10, 0x55, 0xfe, 0x93, 0x1a, 0xa6, 0x41, 0x16, 0x42, 0x01, 0xb1, 0x60, 0x02, 0x12, 0x1f, 0x41,
	0x2f, 0x86, 0x13, 0x20, 0x24, 0xd6, 0xc8, 0x80, 0x04, 0x2b, 0xe4, 0xee, 0x14, 0x1d, 0x4b, 0x4e,
	0xdc, 0xc4, 0xce, 0x71, 0x38, 0x0c, 0xd7, 0xe1, 0x14, 0xc8, 0x55, 0x4e, 0xf7, 0xf4, 0x80, 0x10,
	0xbb, 0xbc, 0x7a, 0x2f, 0xe5, 0xaa, 0xf7, 0x6c, 0xb8, 0xb5, 0xb5, 0xc3, 0xa0, 0xfd, 0x7a, 0x3f,
	0x59, 0x6f, 0x45, 0x35, 0xd8, 0xd1, 0xf6, 0xf3, 0xa6, 0x7d, 0x09, 0xe5, 0x7b, 0x35, 0xe1, 0xe8,
	0xc5, 0x1d, 0xc8, 0xe6, 0xc9, 0x34, 0xc9, 0xa3, 0xe4, 0x59, 0x2d, 0xc3, 0x67, 0xa8, 0xb8, 0x5e,
	0x35, 0x29, 0x57, 0x5c, 0xaf, 0xda, 0x2f, 0x50, 0x7c, 0xf0, 0xca, 0x3b, 0xf1, 0x10, 0x6a, 0xd5,
	0x75, 0xda, 0x6b, 0x3b, 0x3a, 0xfa, 0x25, 0x93, 0xc7, 0x42, 0x60, 0x3b, 0x34, 0xc8, 0x6c, 0xca,
	0xec, 0xa1, 0x20, 0xee, 0x42, 0xe1, 0xad, 0x57, 0xa6, 0xc9, 0x88, 0x61, 0xd0, 0x7e, 0x86, 0xf2,
	0xcd, 0xec, 0x7b, 0x3b, 0x89, 0x15, 0xa4, 0xba, 0x8b, 0x4d, 0x53, 0xdd, 0x05, 0xbd, 0xb1, 0x3b,
	0x3d, 0xc6, 0x41, 0x18, 0x08, 0x01, 0xf9, 0xa8, 0x06, 0xa4, 0x26, 0xb5, 0xa4, 0xef, 0xa0, 0xc4,
	0x41, 0x69, 0xd3, 0xe4, 0xac, 0x24, 0xd0, 0xbe, 0x80, 0xfc, 0xe3, 0x84, 0xf8, 0x5f, 0x0b, 0xfe,
	0x48, 0xa0, 0x7c, 0x4b, 0x46, 0x89, 0x57, 0x50, 0xb3, 0x65, 0x1e, 0x27, 0x92, 0x9c, 0x5f, 0xdd,
	0x5e, 0x47, 0xdb, 0xd6, 0x3c, 0xaa, 0x3c, 0x2a, 0x44, 0x03, 0xd5, 0x80, 0xce, 0xa9, 0xdd, 0x32,
	0xd2, 0x02, 0xc5, 0x25, 0xe4, 0x7e, 0x42, 0xa4, 0xa1, 0xce, 0xaf, 0x2e, 0x0e, 0x3d, 0xc2, 0x50,
	0x92, 0x28, 0xf1, 0x14, 0x4a, 0x45, 0x1d, 0x9b, 0xe2, 0xef, 0x07, 0x45, 0xba, 0xfd, 0x99, 0x40,
	0xfe, 0x4e, 0x1b, 0x14, 0x0f, 0xe0, 0xec, 0x9b, 0x36, 0x48, 0x16, 0xf0, 0x46, 0x07, 0x7c, 0x1a,
	0x4e, 0x7a, 0x33, 0x9c, 0x06, 0xaa, 0x6d, 0xaf, 0xc6, 0x1d, 0xba, 0x18, 0xc0, 0x02, 0x4f, 0x63,
	0xcb, 0x6f, 0xc6, 0x76, 0x0f, 0x4a, 0xe7, 0x95, 0x9f, 0x1d, 0xcd, 0x58, 0xcb, 0x88, 0xc4, 0x7d,
	0x38, 0xdb, 0x18, 0xbb, 0xf9, 0x1a, 0xbc, 0x2d, 0x79, 0xf3, 0x80, 0x3f, 0x4d, 0x26, 0xe4, 0xb1,
	0x57, 0x7e, 0xdb, 0x37, 0x15, 0xe7, 0x41, 0xa0, 0xbd, 0x84, 0x0b, 0xb6, 0x58, 0xe2, 0xf7, 0x19,
	0x9d, 0x5f, 0x62, 0x48, 0x8e, 0x31, 0xfc, 0x4a, 0x60, 0xb5, 0x68, 0xdc, 0xde, 0x8e, 0x0e, 0xff,
	0x14, 0x5d, 0x33, 0x2d, 0xfd, 0xa7, 0x69, 0xe2, 0x31, 0x14, 0xc1, 0x9b, 0xb0, 0x6f, 0x76, 0x92,
	0x40, 0x70, 0x52, 0x32, 0x27, 0x9e, 0x40, 0x11, 0x16, 0x72, 0x31, 0xa6, 0xd5, 0x41, 0x44, 0x17,
	0x5e, 0x32, 0x29, 0x9e, 0x43, 0xb5, 0xa7, 0xe7, 0x12, 0x5c, 0xc8, 0x4e, 0x0e, 0xe5, 0x67, 0x24,
	0x17, 0x3e, 0x8c, 0xc7, 0xb7, 0x83, 0x5c, 0xb9, 0xae, 0x8c, 0x9b, 0x45, 0x7a, 0x53, 0xd2, 0x93,
	0x7c, 0xfd, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xb2, 0xc3, 0xb8, 0x16, 0xa2, 0x03, 0x00, 0x00,
}
