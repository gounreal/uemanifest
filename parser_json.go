package uemanifest

import (
	"errors"
	"io"
)

func ParseJson(r io.Reader) (*Manifest, error) {
	return nil, errors.New("json manifests are not supported") // + yet*
}
