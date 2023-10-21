package entities

import "io"

type ReadSeekerAt interface {
	io.ReadSeeker
	io.ReaderAt
}
type Attachment struct {
	File        ReadSeekerAt
	ContentType string
	Prefix      string
	Ext         string
}
