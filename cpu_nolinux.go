//go:build !linux

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import "errors"

func getCPUFreqStats() (stats []CPUFreqStats, err error) {
	return nil, errors.New("not implemented for non-linux platforms")
}
