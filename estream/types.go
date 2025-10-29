// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package estream

//go:generate stringer -type=blockID -trimprefix=block

type blockID int8

const (
	blockPlainKey blockID = iota + 1
	blockEncryptedKey
	blockEncStream
	blockPlainStream
	blockDatablock
	blockEOS
	blockEOF
	blockError
)

type checksumType uint8

//go:generate stringer -type=checksumType -trimprefix=checksumType

const (
	checksumTypeNone checksumType = iota
	checksumTypeXxhash

	checksumTypeUnknown
)

func (c checksumType) valid() bool {
	return c >= checksumTypeNone && c < checksumTypeUnknown
}
