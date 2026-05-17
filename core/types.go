package core

type FileKind string

const (
	CodeKind   FileKind = "code"
	TextKind   FileKind = "text"
	JSONKind   FileKind = "json"
	BinaryKind FileKind = "binary"
	UnkownKind FileKind = "unkown"
)
