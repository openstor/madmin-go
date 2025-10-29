//go:build linux

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
)

func getCPUFreqStats() (stats []CPUFreqStats, err error) {
	// Attempt to read CPU stats for governor on CPU0
	// which is enough indicating atleast the system
	// has one CPU.
	cpuName := "cpu" + strconv.Itoa(0)

	governorPath := filepath.Join(
		"/sys/devices/system/cpu",
		cpuName,
		"cpufreq",
		"scaling_governor",
	)

	content, err1 := os.ReadFile(governorPath)
	if err1 != nil {
		err = err1
		return stats, err
	}

	stats = append(stats, CPUFreqStats{
		Name:     cpuName,
		Governor: string(bytes.TrimSpace(content)),
	})

	return stats, nil
}
