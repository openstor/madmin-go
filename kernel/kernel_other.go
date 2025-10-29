// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

//go:build !linux

package kernel

// VersionFromRelease only implemented on Linux.
func VersionFromRelease(_ string) (uint32, error) {
	return 0, nil
}

// Version only implemented on Linux.
func Version(_, _, _ int) uint32 {
	return 0
}

// CurrentRelease only implemented on Linux.
func CurrentRelease() (string, error) {
	return "", nil
}

// CurrentVersion only implemented on Linux.
func CurrentVersion() (uint32, error) {
	return 0, nil
}
