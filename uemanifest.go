package uemanifest

import (
	"encoding/binary"
	"errors"
	"io"
)

const BinaryMagic = 0x44BEC00C

func Parse(r io.ReadSeeker) (*Manifest, error) {
	// Read the first 4 bytes from the reader
	var first4Bytes [4]byte
	if _, err := r.Read(first4Bytes[:]); err != nil {
		return nil, err
	}

	// If there are any '{' in the first 4 bytes, then it should be a JSON manifest.
	if isJsonManifest([4]byte(first4Bytes)) {
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		return ParseJson(r)
	}

	// If the first 4 bytes matches the binary's magic header, then should be a Binary manifest.
	if binary.LittleEndian.Uint32(first4Bytes[:]) == BinaryMagic {
		return ParseBinary(r)
	}

	return nil, errors.New("invalid manifest")
}

func isJsonManifest(firstBytes [4]byte) bool {
	for _, b := range firstBytes {
		if b == '{' {
			return true
		}
	}

	return false
}
