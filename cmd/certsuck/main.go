package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	ErrMissingHost = errors.New("missing host")
)

func showChains(w io.Writer, chains [][]*x509.Certificate) {
	// Show the certs.
	for i, chain := range chains {
		fmt.Printf("Chain %d\n", i)
		for i, crt := range chain {
			fmt.Fprintf(w, "  %1d Subject: %s\n    Issuer:  %s\n", i, crt.Subject, crt.Issuer)
		}
	}
}

func getHostPart(hostPort string) string {
	parts := strings.Split(hostPort, ":")
	return parts[0]
}

func writeDerFiles(chain []*x509.Certificate, opts *options) error {
	// Create prefix
	prefix := getHostPart(opts.hostPort) + "-"
	if len(opts.derPrefix) > 0 {
		prefix = opts.derPrefix
	}
	// Build the file name using the prefix and path.
	prefix = filepath.Join(opts.derDir, prefix)

	for i, crt := range chain {
		if i == 0 && opts.noServer {
			continue
		}
		if isRoot(crt) && opts.noRoot {
			continue
		}
		name := fmt.Sprintf("%s%02d.der", prefix, i)
		if err := os.WriteFile(name, crt.Raw, 0777); err != nil {
			return fmt.Errorf("error writing %s", err)
		}
	}
	return nil
}

func showPems(w io.Writer, chain []*x509.Certificate, opts *options) error {
	var pb = pem.Block{
		Type: "CERTIFICATE",
	}
	for i, crt := range chain {
		isServer := i == 0
		if isServer && opts.noServer {
			continue
		}
		isRoot := isRoot(crt)
		if opts.noRoot && isRoot {
			continue
		}
		if isServer {
			fmt.Fprintln(w, "Server certificate")
		}
		if isRoot {
			fmt.Fprintln(w, "Root certificate")
		}
		fmt.Fprintf(w, "%d  Subject: %s\n   Issuer:    %s\n", i, crt.Subject, crt.Issuer)
		pb.Bytes = crt.Raw
		if err := pem.Encode(w, &pb); err != nil {
			return fmt.Errorf("error encoding %s", crt.Subject)
		}
	}
	return nil
}

func isRoot(crt *x509.Certificate) bool {
	return slices.Compare(crt.RawIssuer, crt.RawSubject) == 0
}

func showOptions(w io.Writer, opts *options) {
	fmt.Fprintf(w, "Options:\n%s\n", opts.prettyString())
}

func run(opts *options) error {
	var chains = [][]*x509.Certificate{}
	ow := os.Stdout

	if opts.showOpts {
		showOptions(ow, opts)
	}

	if len(opts.hostPort) == 0 {
		return ErrMissingHost
	}

	callback := func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		chains = slices.Clone(verifiedChains)
		return nil
	}

	// Create a config for the callback.
	tlsCfg := tls.Config{
		VerifyPeerCertificate: callback,
	}
	// Connect to the host
	conn, err := tls.Dial("tcp", opts.hostPort, &tlsCfg)
	if err != nil {
		showOptions(ow, opts)
		return err
	}
	defer conn.Close()
	showChains(ow, chains)

	// No out required.
	if !opts.showOut && !opts.derOut {
		return nil
	}

	var longest []*x509.Certificate
	for _, c := range chains {
		if len(c) > len(longest) {
			longest = c
		}
	}

	if opts.showOut {
		fmt.Fprintln(ow)
		if err := showPems(ow, longest, opts); err != nil {
			return fmt.Errorf("error showing PEMs: %w", err)
		}
	}

	if opts.derOut {
		if err := writeDerFiles(longest, opts); err != nil {
			return fmt.Errorf("error writing DER files: %w", err)
		}
	}
	return nil
}

func main() {
	opts, err := getOptions(os.Args)
	if err != nil {
		fmt.Printf("Options error: %s\n", err.Error())
		os.Exit(1)
	}
	if err := run(opts); err != nil {
		fmt.Printf("Run error: %s\n", err.Error())
		usage(os.Args)
		os.Exit(2)
	}
	os.Exit(0)
}
