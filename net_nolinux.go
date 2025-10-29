//go:build !linux

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

// GetNetInfo returns information of the given network interface
// Not implemented for non-linux platforms
func GetNetInfo(addr string, iface string) NetInfo {
	return NetInfo{
		NodeCommon: NodeCommon{
			Addr:  addr,
			Error: "Not implemented for non-linux platforms",
		},
		Interface: iface,
	}
}
