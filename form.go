package maventa

import (
	"io"
	"net/url"
)

type Form interface {
	Values() url.Values
	Files() map[string]File
}

type File struct {
	Filename string
	Reader   io.Reader
}
