// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

//msgp:tag json
//go:generate msgp -d clearomitted -d "timezone utc" -file $GOFILE

// CatalogDataFile contains information about an output file from a catalog job run.
type CatalogDataFile struct {
	Key         string `json:"key"`
	Size        uint64 `json:"size"`
	MD5Checksum string `json:"MD5Checksum"`
}

// CatalogManifestVersion represents the version of a catalog manifest.
type CatalogManifestVersion string

// Valid values for CatalogManifestVersion.
const (
	// We use AWS S3's manifest file version here as we are following the same
	// format at least initially.
	CatalogManifestVersion1 CatalogManifestVersion = "2016-11-30"
)

// CatalogManifest represents the manifest of a catalog job's result.
type CatalogManifest struct {
	SourceBucket      string                 `json:"sourceBucket"`
	DestinationBucket string                 `json:"destinationBucket"`
	Version           CatalogManifestVersion `json:"version"`
	CreationTimestamp string                 `json:"creationTimestamp"`
	FileFormat        string                 `json:"fileFormat"`
	FileSchema        string                 `json:"fileSchema"`
	Files             []CatalogDataFile      `json:"files"`
}
