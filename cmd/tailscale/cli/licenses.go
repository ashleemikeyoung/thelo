// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"context"

	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
)

var licensesCmd = &ffcli.Command{
	Name:       "licenses",
	ShortUsage: "tailscale licenses",
	ShortHelp:  "Get open source license information",
	LongHelp:   "Get open source license information",
	Exec:       runLicenses,
}

func runLicenses(ctx context.Context, args []string) error {
	// upstream points at tailscale.com/licenses/<platform>, the lists for
	// its own apps (licenses.LicensesURL); this client prints its own
	fmt.Fprint(unbranded(Stdout), brandLicenses())
	return nil
}
