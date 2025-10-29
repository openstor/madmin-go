// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

//go:build (linux && 386) || (linux && amd64) || (linux && arm64) || (linux && loong64) || (linux && mips64) || (linux && mips64le) || (linux && mips)

package kernel

func utsnameStr(in []int8) string {
	out := make([]byte, 0, len(in))
	for i := 0; i < len(in); i++ {
		if in[i] == 0x00 {
			break
		}
		out = append(out, byte(in[i]))
	}
	return string(out)
}
