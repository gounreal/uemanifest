package uemanifest

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/gounreal/uereader"
)

type FManifestHeader struct {
	HeaderSize           uint32                // The size of this header.
	DataSizeUncompressed uint32                // The size of this data uncompressed.
	DataSizeCompressed   uint32                // The size of this data compressed.
	SHAHash              [20]byte              // The SHA1 hash for the manifest data that follows.
	StoredAs             EManifestStorageFlags // How the chunk data is stored.
	Version              EFeatureLevel         // The version of this header and manifest data format, driven by the feature level.
}

func ParseBinary(r io.ReadSeeker) (*Manifest, error) {
	ar := uereader.NewReader("main", r)

	header := uereader.SubReader(ar, "ManifestHeader", readHeader)
	ar.SetPos(int64(header.HeaderSize))
	if err := ar.Err(); err != nil {
		return nil, err
	}

	rawData := bytes.NewReader(ar.Bytes(int(header.DataSizeCompressed)))
	if header.StoredAs&EManifestStorageFlagsCompressed != 0 {
		newReader, err := zlib.NewReader(rawData)
		if err != nil {
			return nil, err
		}

		newData, err := io.ReadAll(newReader)
		if err != nil {
			return nil, err
		}

		rawData = bytes.NewReader(newData)
	}

	ar = uereader.NewReader("manifest", rawData)

	var manifest = &Manifest{
		ManifestMeta:     uereader.SubReader(ar, "ManifestMeta", readManifestMeta),
		ChunkDataList:    uereader.SubReader(ar, "ChunkDataList", readChunkDataList),
		FileManifestList: uereader.SubReader(ar, "FileManifestList", readFileManifestList),
		CustomFields:     uereader.SubReader(ar, "CustomFields", readCustomFields),
	}

	return manifest, ar.Err()
}

func readHeader(ar *uereader.Reader) (*FManifestHeader, error) {
	header := &FManifestHeader{
		HeaderSize:           ar.UInt32(),
		DataSizeUncompressed: ar.UInt32(),
		DataSizeCompressed:   ar.UInt32(),
		SHAHash:              ar.ShaHash(),
		StoredAs:             EManifestStorageFlags(ar.Byte()),
		Version:              EFeatureLevel(ar.Int32()),
	}

	return header, nil
}

func readManifestMeta(ar *uereader.Reader) (FManifestMeta, error) {
	startPos := ar.Pos()

	dataSize := ar.UInt32()
	dataVersion := EManifestMetaVersion(ar.UInt8())

	meta := FManifestMeta{
		FeatureLevel:  EFeatureLevel(ar.Int32()),
		IsFileData:    ar.Bool(),
		AppID:         ar.UInt32(),
		AppName:       ar.String(),
		BuildVersion:  ar.String(),
		LaunchExe:     ar.String(),
		LaunchCommand: ar.String(),
		PrereqIDs:     uereader.ReadSlice(ar, func(ar *uereader.Reader) (string, error) { return ar.String(), nil }),
		PrereqName:    ar.String(),
		PrereqPath:    ar.String(),
		PrereqArgs:    ar.String(),
	}

	if dataVersion >= EManifestMetaVersionSerialisesBuildId {
		meta.BuildID = ar.String()
	}

	ar.SetPos(startPos + int64(dataSize))

	return meta, nil
}

func readChunkDataList(ar *uereader.Reader) ([]FChunkInfo, error) {
	startPos := ar.Pos()

	dataSize := ar.UInt32()
	ar.UInt8() // dataVersionInt := EChunkDataListVersion(uint8)
	elementCount := ar.Int32()

	chunks := make([]FChunkInfo, elementCount)

	for idx := range chunks {
		chunks[idx].Guid = ar.BigEndianUUID()
	}
	for idx := range chunks {
		chunks[idx].Hash = ar.UInt64()
	}
	for idx := range chunks {
		chunks[idx].ShaHash = ar.ShaHash()
	}
	for idx := range chunks {
		chunks[idx].GroupNumber = ar.UInt8()
	}
	for idx := range chunks {
		chunks[idx].WindowSize = ar.UInt32()
	}
	for idx := range chunks {
		chunks[idx].FileSize = ar.Int64()
	}

	ar.SetPos(startPos + int64(dataSize))
	return chunks, nil
}

func readFileManifestList(ar *uereader.Reader) ([]FFileManifest, error) {
	startPos := ar.Pos()

	dataSize := ar.UInt32()
	ar.UInt8() // dataVersionInt := EChunkDataListVersion(uint8)
	elementCount := ar.Int32()

	files := make([]FFileManifest, elementCount)

	for idx := range files {
		files[idx].Filename = ar.String()
	}
	for idx := range files {
		files[idx].SymlinkTarget = ar.String()
	}
	for idx := range files {
		files[idx].FileHash = ar.ShaHash()
	}
	for idx := range files {
		files[idx].FileMetaFlags = EFileMetaFlags(ar.UInt8())
	}
	for idx := range files {
		files[idx].InstallTags = uereader.ReadSlice(ar, func(ar *uereader.Reader) (string, error) { return ar.String(), nil })
	}
	for idx := range files {
		files[idx].ChunkParts = uereader.ReadSlice(ar, func(ar *uereader.Reader) (FChunkPart, error) {
			cStartPos := ar.Pos()
			cDataSize := ar.UInt32()

			chunk := FChunkPart{
				Guid:   ar.BigEndianUUID(),
				Offset: ar.UInt32(),
				Size:   ar.UInt32(),
			}

			files[idx].FileSize += uint64(chunk.Size)
			ar.SetPos(cStartPos + int64(cDataSize))

			return chunk, nil
		})
	}

	ar.SetPos(startPos + int64(dataSize))
	return files, nil
}

func readCustomFields(ar *uereader.Reader) (map[string]string, error) {
	startPos := ar.Pos()

	dataSize := ar.UInt32()
	ar.UInt8() // dataVersionInt := EChunkDataListVersion(uint8)
	elementCount := ar.Int32()

	keys := uereader.ReadArray(ar, elementCount, func(ar *uereader.Reader) (string, error) { return ar.String(), nil })
	values := uereader.ReadArray(ar, elementCount, func(ar *uereader.Reader) (string, error) { return ar.String(), nil })

	fields := map[string]string{}
	for idx, key := range keys {
		fields[key] = values[idx]
	}

	ar.SetPos(startPos + int64(dataSize))
	return fields, nil
}
