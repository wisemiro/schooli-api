package models

import "mime/multipart"

type FileUploadModel struct {
	ObjName     string
	FileBuf     multipart.File
	FileSize    int64
	ContentType string
	FolderName  string
}

type MultipleFileUploadModel struct {
	FileNames  []*multipart.FileHeader
	FolderName string
}
