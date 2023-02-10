package uemanifest

import "github.com/google/uuid"

type Manifest struct {
	ManifestMeta     FManifestMeta
	ChunkDataList    []FChunkInfo      // The list of chunks.
	FileManifestList []FFileManifest   // The list of files.
	CustomFields     map[string]string // The map of field name to field data.
}

type FManifestMeta struct {
	FeatureLevel  EFeatureLevel // The feature level support this build was created with, regardless of the serialised format.
	IsFileData    bool          // Whether this is a legacy 'nochunks' build.
	AppID         uint32        // The app id provided at generation.
	AppName       string        // The app name string provided at generation.
	BuildVersion  string        // The build version string provided at generation.
	LaunchExe     string        // The file in this manifest designated the application executable of the build.
	LaunchCommand string        // The command line required when launching the application executable.
	PrereqIDs     []string      // The set of prerequisite ids for dependencies that this build's prerequisite installer will apply.
	PrereqName    string        // A display string for the prerequisite provided at generation.
	PrereqPath    string        // The file in this manifest designated the launch executable of the prerequisite installer.
	PrereqArgs    string        // The command line required when launching the prerequisite installer.
	BuildID       string        // A unique build id generated at original chunking time to identify an exact build.
}

type FChunkInfo struct {
	Guid        uuid.UUID // The GUID for this data.
	Hash        uint64    // The FRollingHash hashed value for this chunk data.
	ShaHash     [20]byte  // The FSHA hashed value for this chunk data.
	GroupNumber uint8     // The group number this chunk divides into.
	WindowSize  uint32    // The window size for this chunk.
	FileSize    int64     // The file download size for this chunk.
}

type FFileManifest struct {
	Filename      string         // The build relative filename.
	SymlinkTarget string         // Whether this is a symlink to another file.
	FileHash      [20]byte       // The file SHA1.
	FileMetaFlags EFileMetaFlags // The flags for this file.
	InstallTags   []string       // The install tags for this file.
	ChunkParts    []FChunkPart   // The list of chunk parts to stitch.
	FileSize      uint64         // The size of this file.
}

type FChunkPart struct {
	Guid   uuid.UUID // The GUID of the chunk containing this part.
	Offset uint32    // The offset of the first byte into the chunk.
	Size   uint32    // The size of this part.
}
