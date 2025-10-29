//go:build !linux

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package cgroup

import "errors"

// GetMemoryLimit - Not implemented in non-linux platforms
func GetMemoryLimit(_ int) (limit uint64, err error) {
	return limit, errors.New("not implemented for non-linux platforms")
}
