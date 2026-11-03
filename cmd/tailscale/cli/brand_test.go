// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestBrandText(t *testing.T) {
	tests := []struct{ in, want string }{
		{"tailscale up [flags]", "thelo up [flags]"},
		{`see "tailscale set --help" for valid settings`, `see "thelo set --help" for valid settings`},
		{"use `tailscale set --exit-node=`", "use `thelo set --exit-node=`"},
		{"Use 'sudo tailscale up'.", "Use 'sudo thelo up'."},
		{"(tailscale status)", "(thelo status)"},
		// the product name in running text
		{"Tailscale is starting. Please wait.", "Thelo is starting. Please wait."},
		{"Tailscale is stopped.", "Thelo is stopped."},
		{"Tailscale SSH", "Thelo SSH"},
		{"Tailscale DNS: enabled.", "Thelo DNS: enabled."},
		{"could not reach the Tailscale coordination server", "could not reach the Thelo coordination server"},
		{"Connect to Tailscale, logging in if needed", "Connect to Thelo, logging in if needed"},
		{"Run 'tailscale up' to start Tailscale.", "Run 'thelo up' to start Thelo."},
		// tailnet and Funnel
		{"Share a directory with your tailnet", "Share a directory with your thelonet"},
		{"Tailnet Lock is ENABLED.", "Thelonet Lock is ENABLED."},
		{"MagicDNS: enabled tailnet-wide (suffix = x)", "MagicDNS: enabled thelonet-wide (suffix = x)"},
		{"switch between tailnets", "switch between thelonets"},
		{"Turn on/off Funnel service", "Turn on/off Halo service"},
		{"Tailscale Funnel is not enabled", "Halo is not enabled"},
		{"# Funnel on:", "# Halo on:"},
		{"Show current serve/funnel status", "Show current serve/halo status"},
		{"tailscale funnel status", "thelo halo status"},
		{"use `tailscale funnel`", "use `thelo halo`"},
		{"tailscale funnelx", "thelo funnelx"},
		// left alone: identifiers, paths, URLs, the daemon, the credit and the company
		{"tailscaled --state=x", "tailscaled --state=x"},
		{"https://pkgs.tailscale.com/stable", "https://pkgs.tailscale.com/stable"},
		{"tailscale.com/cmd/tailscale", "tailscale.com/cmd/tailscale"},
		{"TailscaleIPs", "TailscaleIPs"},
		{"CurrentTailnet", "CurrentTailnet"},
		{"web.ts-tailnet.ts.net", "web.ts-tailnet.ts.net"},
		{"tailnet.ts.net", "tailnet.ts.net"},
		{"AllowFunnel", "AllowFunnel"},
		{"/Applications/Tailscale.app", "/Applications/Tailscale.app"},
		{brandCredit, brandCredit},
		{"Copyright (c) Tailscale Inc & contributors", "Copyright (c) Tailscale Inc & contributors"},
	}
	for _, tt := range tests {
		if got := brandText(tt.in); got != tt.want {
			t.Errorf("brandText(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestBrandErrKeepsTheChain(t *testing.T) {
	base := errors.New("final argument to 'tailscale file cp' must end in colon")
	err := brandErr(base)
	if !strings.Contains(err.Error(), "'thelo file cp'") {
		t.Errorf("message not branded: %q", err)
	}
	if !errors.Is(err, base) {
		t.Error("errors.Is no longer finds the original error")
	}
	if plain := errors.New("no command here"); brandErr(plain) != plain {
		t.Error("an error without command references should be returned as is")
	}
	if brandErr(nil) != nil {
		t.Error("nil must stay nil")
	}
}

func TestRootHelpIsBranded(t *testing.T) {
	root := newRootCmd()
	help := usageFunc(root)
	for _, want := range []string{"thelo [flags] <subcommand>", brandProduct, brandCredit} {
		if !strings.Contains(help, want) {
			t.Errorf("root help missing %q:\n%s", want, help)
		}
	}
	for _, sub := range root.Subcommands {
		if sub.Name == "up" {
			if h := usageFunc(sub); !strings.Contains(h, "thelo up [flags]") || strings.Contains(h, "tailscale up") {
				t.Errorf("up help not branded:\n%s", h)
			}
		}
	}
}

func TestBrandWriter(t *testing.T) {
	var buf strings.Builder
	w := branded(&buf)
	if branded(w) != w {
		t.Error("branded must not wrap twice")
	}
	if unbranded(w) != &buf {
		t.Error("unbranded must return the underlying writer")
	}
	fmt.Fprintln(w, "Tailscale is starting. Please wait.")
	fmt.Fprintln(w, `{"Health": ["Tailscale is starting. Please wait."]}`)
	got := buf.String()
	want := "Thelo is starting. Please wait.\n" + `{"Health": ["Tailscale is starting. Please wait."]}` + "\n"
	if got != want {
		t.Errorf("got %q, want %q (JSON passes through unchanged)", got, want)
	}
}

func TestHaloCommand(t *testing.T) {
	var names []string
	for _, sub := range newRootCmd().Subcommands {
		names = append(names, sub.Name)
	}
	got := strings.Join(names, " ")
	if !strings.Contains(" "+got+" ", " halo ") || strings.Contains(" "+got+" ", " funnel ") {
		t.Errorf("want a halo subcommand and no funnel: %s", got)
	}
}

func TestLicensesReproducesUpstream(t *testing.T) {
	l := brandLicenses()
	for _, want := range []string{"Copyright (c) 2020 Tailscale Inc & contributors.", "THIS SOFTWARE IS PROVIDED", "/licenses/tailscale.md", "/PATENTS"} {
		if !strings.Contains(l, want) {
			t.Errorf("licenses missing %q", want)
		}
	}
}
