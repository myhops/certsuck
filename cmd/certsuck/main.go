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

func showChain(w io.Writer, chain []*x509.Certificate, indent string) {
	for i, crt := range chain {
		fmt.Fprintf(w, "%s%1d Subject: %s\n    Issuer:  %s\n", indent, i, crt.Subject, crt.Issuer)
	}
}

func showChains(w io.Writer, chains map[string][]*x509.Certificate) {
	// Show the certs.
	for k, chain := range chains {
		fmt.Printf("%s\n", k)
		showChain(w, chain, "  ")
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

func collectChains(opts *options) (map[string][]*x509.Certificate, error) {
	var verifiedChains = [][]*x509.Certificate{}
	var peerCerts = []*x509.Certificate{}

	verifyConnectionCallback := func(state tls.ConnectionState) error {
		verifiedChains = slices.Clone( state.VerifiedChains)
		peerCerts = slices.Clone(state.PeerCertificates)
		return nil
	}

	// Create a config for the callback.
	tlsCfg := tls.Config{
		InsecureSkipVerify: opts.insecure,
		VerifyConnection: verifyConnectionCallback,
	}
	// Connect to the host
	conn, err := tls.Dial("tcp", opts.hostPort, &tlsCfg)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := make(map[string][]*x509.Certificate)

	for i, chain := range verifiedChains {
		res[fmt.Sprintf("Verified chain %2d", i)] = chain
	}
	res["Peer chain"] = peerCerts
	return res, nil
}

func getLongestFromMap(chains map[string][]*x509.Certificate) string {
	var l int
	var longestKey string
	for k, chain := range chains {
		if len(chain) > l {
			longestKey = k
			l = len(chain)
		}
	}
	return longestKey
}

func run(opts *options) error {
	ow := os.Stdout

	if opts.showOpts {
		showOptions(ow, opts)
	}

	if len(opts.hostPort) == 0 {
		return ErrMissingHost
	}

	chains, err := collectChains(opts)
	if err != nil {
		return err
	}

	showChains(ow, chains)

	// No out required.
	if !opts.showOut && !opts.derOut {
		return nil
	}

	ln := getLongestFromMap(chains)

	if opts.showOut {
		fmt.Fprintln(ow)
		if err := showPems(ow, chains[ln], opts); err != nil {
			return fmt.Errorf("error showing PEMs: %w", err)
		}
	}

	if opts.derOut {
		if err := writeDerFiles(chains[ln], opts); err != nil {
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
