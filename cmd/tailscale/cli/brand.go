// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"io"
	"strings"

	"tailscale.com/version"
)

func init() {
	// thelo halo, not thelo funnel
	i := infoMap[funnel]
	i.Name = brandFunnelCommand
	infoMap[funnel] = i
}

// Thelo branding (Qavats Ranch). This client is built from Tailscale's open-source
// code and joins the Thelo network (a thelonet) run by Headscale.
// Only what people read is renamed: the command name, help and messages.
// Protocols, file paths, the daemon's service label and API names keep
// their upstream names so that upstream releases merge cleanly.
const (
	// brandProduct is the product name shown to people.
	brandProduct = "Thelo"
	// brandShort is the product name in running text ("Thelo is starting").
	brandShort = "Thelo"
	// brandCommand is the name of this CLI.
	brandCommand = "thelo"
	// brandNet is what upstream calls a tailnet.
	brandNet = "thelonet"
	// brandFunnel is what upstream calls Funnel (Tailscale Funnel).
	brandFunnel = "Halo"
	// brandFunnelCommand is the subcommand upstream calls funnel.
	brandFunnelCommand = "halo"
	// brandCredit names the upstream project, as its BSD-3-Clause licence asks.
	brandCredit = "Built on Tailscale (https://github.com/tailscale/tailscale), BSD-3-Clause."
	// brandCompany is the company that ships Thelo.
	brandCompany = "Qavats Ranch"
)

// brandText rewrites what people read: command references, so that
// `tailscale up` reads `thelo up`, and the product's words, so that
// "Tailscale is starting" reads "Thelo is starting" (see [brandProse]). A reference is "tailscale " at the start
// of the text or after whitespace, a quote, a backtick or "(". That leaves
// alone "Tailscale" (the company or protocol), "tailscaled", and hostnames
// such as "pkgs.tailscale.com". (No regexp: tailscaled's dependency test
// keeps it out of the binary.)
func brandText(s string) string {
	return brandProse(brandCommands(s))
}

// brandCommands rewrites command references: "tailscale <subcommand>".
func brandCommands(s string) string {
	const ref = "tailscale "
	if !strings.Contains(s, ref) {
		return s
	}
	var b strings.Builder
	for {
		i := strings.Index(s, ref)
		if i < 0 {
			b.WriteString(s)
			return b.String()
		}
		b.WriteString(s[:i])
		if i == 0 || strings.ContainsRune(" \t\n'\"`(", rune(s[i-1])) {
			b.WriteString(brandCommand + " ")
			if rest := s[i+len(ref):]; strings.HasPrefix(rest, "funnel") && isProductMention("", rest[len("funnel"):]) {
				b.WriteString(brandFunnelCommand)
				s = rest[len("funnel"):]
				continue
			}
		} else {
			b.WriteString(ref)
		}
		s = s[i+len(ref):]
	}
}

// brandErr applies [brandText] to an error's message for display.
func brandErr(err error) error {
	if err == nil {
		return nil
	}
	msg := brandText(err.Error())
	if msg == err.Error() {
		return err
	}
	return &brandedError{msg: msg, err: err}
}

type brandedError struct {
	msg string
	err error
}

func (e *brandedError) Error() string { return e.msg }
func (e *brandedError) Unwrap() error { return e.err }

// brandProse replaces the product's words in running text: "Tailscale" with
// [brandShort], "tailnet" with [brandNet] (capitalised and plural too), and
// "Funnel" or "Tailscale Funnel" with [brandFunnel] ("serve/funnel" too). It leaves alone the
// words inside identifiers, paths and hostnames ("TailscaleIPs",
// "Tailscale.app", "tailscale.com", "CurrentTailnet", "ts-tailnet.ts.net"),
// the upstream credit ("Built on Tailscale") and the company ("Tailscale Inc").
func brandProse(s string) string {
	s = replaceWord(s, "Tailscale Funnel", brandFunnel)
	s = replaceWord(s, "Tailscale", brandShort)
	s = replaceWord(s, "Funnel", brandFunnel)
	s = replaceWord(s, "serve/funnel", "serve/"+brandFunnelCommand)
	s = replaceWord(s, "tailnet", brandNet)
	return replaceWord(s, "Tailnet", strings.ToUpper(brandNet[:1])+brandNet[1:])
}

// replaceWord replaces name with repl where [isProductMention] says it is
// the word on its own (a plural "s" may follow).
func replaceWord(s, name, repl string) string {
	if !strings.Contains(s, name) {
		return s
	}
	var b strings.Builder
	for {
		i := strings.Index(s, name)
		if i < 0 {
			b.WriteString(s)
			return b.String()
		}
		before, after := b.String()+s[:i], s[i+len(name):]
		b.WriteString(s[:i])
		if isProductMention(before, after) || strings.HasPrefix(after, "s") && isProductMention(before, after[1:]) {
			b.WriteString(repl)
		} else {
			b.WriteString(name)
		}
		s = after
	}
}

func isProductMention(before, after string) bool {
	if before != "" {
		if c := before[len(before)-1]; isWordByte(c) || strings.IndexByte("./-_@", c) >= 0 {
			return false
		}
	}
	if after != "" {
		c := after[0]
		if isWordByte(c) {
			return false
		}
		if c == '.' && len(after) > 1 && isWordByte(after[1]) {
			return false // Tailscale.app, Tailscale.exe
		}
	}
	if strings.HasSuffix(before, "Built on ") || strings.HasPrefix(after, " Inc") || strings.HasPrefix(after, ", Inc") {
		return false
	}
	return true
}

func isWordByte(c byte) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// brandWriter applies [brandText] to what the CLI prints for people. JSON
// documents (a write that starts with '{' or '[') pass through unchanged, so
// scripts reading --json output see upstream's values.
type brandWriter struct{ w io.Writer }

func (bw brandWriter) Write(p []byte) (int, error) {
	if t := strings.TrimLeft(string(p), " \t\r\n"); strings.HasPrefix(t, "{") || strings.HasPrefix(t, "[") {
		return bw.w.Write(p)
	}
	if _, err := io.WriteString(bw.w, brandText(string(p))); err != nil {
		return 0, err
	}
	return len(p), nil
}

// branded wraps w in a [brandWriter] once.
func branded(w io.Writer) io.Writer {
	if _, ok := w.(brandWriter); ok {
		return w
	}
	return brandWriter{w}
}

// unbranded returns the writer under a [brandWriter].
func unbranded(w io.Writer) io.Writer {
	if bw, ok := w.(brandWriter); ok {
		return bw.w
	}
	return w
}

// brandLicenses is what `thelo licenses` prints: where the code comes from,
// upstream's licence in full (BSD-3-Clause asks binary redistributions to
// reproduce it), and where the dependencies' licences are listed.
func brandLicenses() string {
	tag := "v" + version.Short()
	return brandProduct + ` is built from Tailscale's open-source client,
https://github.com/tailscale/tailscale (release ` + tag + `), and is
distributed under the same licence, reproduced below. Tailscale Inc. also
grants a patent licence for its code:
https://github.com/tailscale/tailscale/blob/` + tag + `/PATENTS

The open-source packages it is built with, and their licences, are listed at
https://github.com/tailscale/tailscale/blob/` + tag + `/licenses/tailscale.md

` + upstreamLicense
}

// upstreamLicense is Tailscale's LICENSE file, unchanged.
const upstreamLicense = `BSD 3-Clause License

Copyright (c) 2020 Tailscale Inc & contributors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`
