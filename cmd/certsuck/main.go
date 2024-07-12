package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/myhops/certsuck/probe"
)

var (
	ErrMissingHost = errors.New("missing host")
)

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
	ow := os.Stdout

	if opts.showOpts {
		showOptions(ow, opts)
	}

	if len(opts.hostPort) == 0 {
		return ErrMissingHost
	}

	chains, err := probe.New(probe.WithInsecure(opts.insecure)).CollectCerts(opts.hostPort)
	if err != nil {
		return err
	}

	fmt.Fprint(ow, chains.String())

	// No out required.
	if !opts.showOut && !opts.derOut {
		return nil
	}

	if opts.showOut {
		fmt.Fprintln(ow)
		if err := showPems(ow, chains.Longest, opts); err != nil {
			return fmt.Errorf("error showing PEMs: %w", err)
		}
	}

	if opts.derOut {
		if err := writeDerFiles(chains.Longest, opts); err != nil {
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
		// usage(os.Args)
		os.Exit(2)
	}
	os.Exit(0)
}
