package main

import (
	"errors"
	"flag"
	"fmt"
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
	derDir    string
	showOpts  bool

	buildJks  bool
	storePass string
	keyPass   string
	keystore  string
}

const (
	fmtOptions       = "-host: %v, -out: %v, -root: %v, -no-server: %v, -der-out: %v, -der-prefix: %s, -der-dir"
	fmtPrettyOptions = `  -host:       %v 
  -out:        %v
  -no-root:    %v
  -no-server:  %v
  -der-out:    %v
  -der-prefix: %s
  -der-dir:    %s`
)

func (opts *options) string(format string) string {
	out := strings.Builder{}
	fmt.Fprintf(&out, format,
		opts.hostPort, opts.showOut, opts.noRoot, opts.noServer, opts.derOut, opts.derPrefix, opts.derDir)
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
	fs.StringVar(&opts.derDir, "der-dir", ".", "Path to write the der files to. Defaults to the current directory")
	fs.StringVar(&opts.keyPass, "key-pass", "changeme", "Password for the keys [changeme]")
	fs.StringVar(&opts.storePass, "store-pass", "changeme", "Password for the store [changeme]")
	fs.StringVar(&opts.keystore, "keystore", "truststore.jks", "Keystore name [truststore.jks]")
	fs.BoolVar(&opts.derOut, "der-out", false, "Output der files. The names of the files is <host>-0x.der [false]")
	fs.BoolVar(&opts.showOut, "out", false, "Show pem output [false]")
	fs.BoolVar(&opts.noRoot, "no-root", false, "Omit the root cert in pem or der output [false]")
	fs.BoolVar(&opts.noServer, "no-server", false, "Omit the server cert in pem or der output [false]")
	fs.BoolVar(&opts.showOpts, "show-opts", false, "Show the options [false]")
	fs.BoolVar(&opts.buildJks, "build-jks", false, "Build the jks file [false]")

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
	fs.StringVar(&opts.derDir, "der-dir", ".", "Path to write the der files to. Defaults to the current directory")
	fs.BoolVar(&opts.derOut, "der-out", false, "Output der files. The names of the files is <host>-0x.der [false]")
	fs.BoolVar(&opts.showOut, "out", false, "Show pem output [false]")
	fs.BoolVar(&opts.noRoot, "no-root", false, "Omit the root cert in pem or der output [false]")
	fs.BoolVar(&opts.noServer, "no-server", false, "Omit the server cert in pem or der output [false]")
	fs.BoolVar(&opts.showOpts, "show-opts", false, "Show the options [false]")
	fs.Usage()
}
