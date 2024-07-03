package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
)

var (
	ErrToFewArguments = errors.New("to few arguments")
)

type options struct {
	hostPort string
	showOut  bool
	noRoot   bool
	noServer bool
}

func getOptions(args []string) (*options, error) {
	if len(args) < 1 {
		return nil, ErrToFewArguments
	}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	opts := &options{}
	fs.StringVar(&opts.hostPort, "host", "", "Hostname plus port")
	fs.BoolVar(&opts.showOut, "out", false, "Show pem output [false]")
	fs.BoolVar(&opts.noRoot, "no-root", false, "Do not show the root cert in pem output [false]")
	fs.BoolVar(&opts.noServer, "no-server", false, "Do not show the server cert in pem output [false]")

	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}
	return opts, nil
}

func usage(args []string) {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	opts := &options{}
	fs.StringVar(&opts.hostPort, "host", "", "Hostname plus port")
	fs.BoolVar(&opts.showOut, "out", false, "Show pem output [false]")
	fs.BoolVar(&opts.noRoot, "no-root", false, "Do not show the root cert in pem output [false]")
	fs.BoolVar(&opts.noServer, "no-server", false, "Do not show the server cert in pem output [false]")

	fs.Usage()
}

func showChains(w io.Writer, chains [][]*x509.Certificate) {
	// Show the certs.
	for i, chain := range chains {
		fmt.Printf("Chain %d\n", i)
		for _, crt := range chain {
			fmt.Fprintf(w, "  Subject: %s\n  Issuer:    %s\n", crt.Subject, crt.Issuer)
		}
	}
}

type showPemsOptions struct {
	noRoot   bool
	noServer bool
}

func showPems(w io.Writer, chain []*x509.Certificate, opts ...showPemsOptions) error {
	var noRoot bool
	var noServer bool

	if len(opts) > 0 {
		noRoot = opts[0].noRoot
		noServer = opts[0].noServer
	}
	var pb = pem.Block{
		Type: "CERTIFICATE",
	}
	for i, crt := range chain {
		isServer := i == 0
		if isServer && noServer {
			continue
		}
		isRoot := isRoot(crt)
		if noRoot && isRoot {
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

func run(opts *options) error {
	var chains = [][]*x509.Certificate{}
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
		return err
	}
	defer conn.Close()
	showChains(os.Stdout, chains)

	if !opts.showOut {
		return nil
	}

	var longest []*x509.Certificate
	for _, c := range chains {
		if len(c) > len(longest) {
			longest = c
		}
	}
	fmt.Fprintln(os.Stdout)
	showPems(os.Stdout, longest, showPemsOptions{noRoot: opts.noRoot, noServer: opts.noServer})
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
