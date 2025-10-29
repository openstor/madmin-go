//go:build linux

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"fmt"

	"github.com/safchain/ethtool"
)

// GetNetInfo returns information of the given network interface
func GetNetInfo(addr string, iface string) (ni NetInfo) {
	ni.Addr = addr
	ni.Interface = iface

	ethHandle, err := ethtool.NewEthtool()
	if err != nil {
		ni.Error = err.Error()
		return ni
	}
	defer ethHandle.Close()

	di, err := ethHandle.DriverInfo(ni.Interface)
	if err != nil {
		ni.Error = fmt.Sprintf("Error getting driver info for %s: %s", ni.Interface, err.Error())
		return ni
	}

	ni.Driver = di.Driver
	ni.FirmwareVersion = di.FwVersion

	ring, err := ethHandle.GetRing(ni.Interface)
	if err != nil {
		ni.Error = fmt.Sprintf("Error getting ring parameters for %s: %s", ni.Interface, err.Error())
		return ni
	}

	ni.Settings = &NetSettings{
		RxMaxPending: ring.RxMaxPending,
		TxMaxPending: ring.TxMaxPending,
		RxPending:    ring.RxPending,
		TxPending:    ring.TxPending,
	}

	channels, err := ethHandle.GetChannels(iface)
	if err != nil {
		ni.Error = fmt.Sprintf("Error getting channels for %s: %s", ni.Interface, err.Error())
	}
	ni.Settings.CombinedCount = channels.CombinedCount
	ni.Settings.MaxCombined = channels.MaxCombined

	return ni
}
