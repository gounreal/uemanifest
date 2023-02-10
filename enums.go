package uemanifest

// Declares flags for manifest headers which specify storage types.
type EFileMetaFlags uint8

const (
	EFileMetaFlagsNone           EFileMetaFlags = 0
	EFileMetaFlagsReadOnly       EFileMetaFlags = 1      // Flag for readonly file.
	EFileMetaFlagsCompressed     EFileMetaFlags = 1 << 1 // Flag for natively compressed.
	EFileMetaFlagsUnixExecutable EFileMetaFlags = 1 << 2 // Flag for unix executable.
)

// A flags enum for manifest headers which specify storage types.
type EManifestStorageFlags uint8

const (
	EManifestStorageFlagsNone       EManifestStorageFlags = 0      // Stored as raw data.
	EManifestStorageFlagsCompressed EManifestStorageFlags = 1      // Flag for compressed data.
	EManifestStorageFlagsEncrypted  EManifestStorageFlags = 1 << 1 // Flag for encrypted. If also compressed, decrypt first. Encryption will ruin compressibility.
)

// Enum which describes the FManifestMeta data version.
type EManifestMetaVersion uint8

const (
	EManifestMetaVersionOriginal EManifestMetaVersion = iota
	EManifestMetaVersionSerialisesBuildId

	// Always after the latest version, signifies the latest version plus 1 to allow initialization simplicity.
	EManifestMetaVersionLatestPlusOne
	EManifestMetaVersionLatest EManifestMetaVersion = (EManifestMetaVersionLatestPlusOne - 1)
)

// An enum type to describe supported features of a certain manifest.
type EFeatureLevel int32

const (
	EFeatureLevelOriginal                                     EFeatureLevel = iota // The original version.
	EFeatureLevelCustomFields                                                      // Support for custom fields.
	EFeatureLevelStartStoringVersion                                               // Started storing the version number.
	EFeatureLevelDataFileRenames                                                   // Made after data files where renamed to include the hash value, these chunks now go to ChunksV2.
	EFeatureLevelStoresIfChunkOrFileData                                           // Manifest stores whether build was constructed with chunk or file data.
	EFeatureLevelStoresDataGroupNumbers                                            // Manifest stores group number for each chunk/file data for reference so that external readers don't need to know how to calculate them.
	EFeatureLevelChunkCompressionSupport                                           // Added support for chunk compression, these chunks now go to ChunksV3. NB: Not File Data Compression yet.
	EFeatureLevelStoresPrerequisitesInfo                                           // Manifest stores product prerequisites info.
	EFeatureLevelStoresChunkFileSizes                                              // Manifest stores chunk download sizes.
	EFeatureLevelStoredAsCompressedUClass                                          // Manifest can optionally be stored using UObject serialization and compressed.
	EFeatureLevelUNUSED_0                                                          // These two features were removed and never used.
	EFeatureLevelUNUSED_1                                                          // These two features were removed and never used.
	EFeatureLevelStoresChunkDataShaHashes                                          // Manifest stores chunk data SHA1 hash to use in place of data compare, for faster generation.
	EFeatureLevelStoresPrerequisiteIds                                             // Manifest stores Prerequisite Ids.
	EFeatureLevelStoredAsBinaryData                                                // The first minimal binary format was added. UObject classes will no longer be saved out when binary selected.
	EFeatureLevelVariableSizeChunksWithoutWindowSizeChunkInfo                      // Temporary level where manifest can reference chunks with dynamic window size, but did not serialize them. Chunks from here onwards are stored in ChunksV4.
	EFeatureLevelVariableSizeChunks                                                // Manifest can reference chunks with dynamic window size, and also serializes them.
	EFeatureLevelUsesRuntimeGeneratedBuildId                                       // Manifest uses a build id generated from its metadata.
	EFeatureLevelUsesBuildTimeGeneratedBuildId                                     // Manifest uses a build id generated unique at build time, and stored in manifest.

	EFeatureLevelLatestPlusOne                                                  // !! Always after the latest version entry, signifies the latest version plus 1 to allow the following Latest alias.
	EFeatureLevelLatest              = (EFeatureLevelLatestPlusOne - 1)         // An alias for the actual latest version value.
	EFeatureLevelLatestNoChunks      = EFeatureLevelStoresChunkFileSizes        // An alias to provide the latest version of a manifest supported by file data (nochunks).
	EFeatureLevelLatestJson          = EFeatureLevelStoresPrerequisiteIds       // An alias to provide the latest version of a manifest supported by a json serialized format.
	EFeatureLevelFirstOptimisedDelta = EFeatureLevelUsesRuntimeGeneratedBuildId // An alias to provide the first available version of optimised delta manifest saving.

	EFeatureLevelStoresUniqueBuildId = EFeatureLevelUsesRuntimeGeneratedBuildId // More aliases, but this time for values that have been renamed

	// JSON manifests were stored with a version of 255 during a certain CL range due to a bug.
	// We will treat this as being StoresChunkFileSizes in code.
	EFeatureLevelBrokenJsonVersion EFeatureLevel = 255
	EFeatureLevelInvalid           EFeatureLevel = -1 // This is for UObject default, so that we always serialize it.
)

// func (e EFeatureLevel) C()
