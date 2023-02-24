package model

import "io"

type File struct {
	Size        int64
	Name        string
	ContentType string
	File        io.Reader
}
