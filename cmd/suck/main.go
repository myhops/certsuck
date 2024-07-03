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
	"strings"
)

var (
	ErrToFewArguments = errors.New("to few arguments")
)

type options struct {
	hostPort  string
	showOut   bool
	noRoot    bool
	noServer  bool
	derOut    bool
	derPrefix string
	showOpts  bool
}

const (
	fmtOptions = "Hostport: %v, showOut: %v, noRoot: %v, noServer: %v, derOut: %v, derPrefix: %s"
	fmtPrettyOptions = `  Hostport: %v 
  showOut: %v
  noRoot: %v
  noServer: %v
  derOut: %v
  derPrefix: %s`
)

func (opts *options) string(format string) string {
	out := strings.Builder{}
	fmt.Fprintf(&out, format,
		opts.hostPort, opts.showOut, opts.noRoot, opts.noServer, opts.derOut, opts.derPrefix)
	return out.String()	
}

func (opts *options) String() string {
	return opts.string(fmtOptions)
}

func (opts *options) prettyString() string {
	return opts.string(fmtPrettyOptions)
}

func getOptions(args []string) (*options, error) {
	if len(args) < 1 {
		return nil, ErrToFewArguments
	}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	opts := &options{}
	fs.StringVar(&opts.hostPort, "host", "", "Hostname plus port")
	fs.StringVar(&opts.derPrefix, "der-prefix", "", "Prefix for the der files. Defaults to <host name>-")
	fs.BoolVar(&opts.derOut, "der-out", false, "Output der files. The names of the files is <host>-0x.der [false]")
	fs.BoolVar(&opts.showOut, "out", false, "Show pem output [false]")
	fs.BoolVar(&opts.noRoot, "no-root", false, "Do not show the root cert in pem output [false]")
	fs.BoolVar(&opts.noServer, "no-server", false, "Do not show the server cert in pem output [false]")
	fs.BoolVar(&opts.showOpts, "show-opts", false, "Show the options [false]")

	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}
	return opts, nil
}

func usage(args []string) {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	opts := &options{}
	fs.StringVar(&opts.hostPort, "host", "", "Hostname plus port")
	fs.StringVar(&opts.derPrefix, "der-prefix", "", "Prefix for the der files. Defaults to <host name>-")
	fs.BoolVar(&opts.derOut, "der-out", false, "Output der files. The names of the files is <host>-0x.der [false]")
	fs.BoolVar(&opts.showOut, "out", false, "Show pem output [false]")
	fs.BoolVar(&opts.noRoot, "no-root", false, "Do not show the root cert in pem output [false]")
	fs.BoolVar(&opts.noServer, "no-server", false, "Do not show the server cert in pem output [false]")
	fs.BoolVar(&opts.showOpts, "show-opts", false, "Show the options [false]")
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

	for i, crt := range chain {
		name := fmt.Sprintf("%s%02d.der", prefix, i)
		if err := os.WriteFile(name, crt.Raw, 0777); err != nil {
			return fmt.Errorf("error writing %s", err)
		}
	}
	return nil
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

func showOptions(w io.Writer, opts *options) {
	fmt.Fprintf(w, "Options:\n%s\n", opts.prettyString())
}

func run(opts *options) error {
	var chains = [][]*x509.Certificate{}
	ow := os.Stdout

	if opts.showOpts {
		showOptions(ow, opts)
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

	if opts.derOut {
		if err := writeDerFiles(longest, opts); err != nil {
			return err
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
