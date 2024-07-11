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
	insecure  bool
	cacerts   string
}

const (
	fmtOptions       = "-host: %v, -out: %v, -root: %v, -no-server: %v, -insecure %v, -der-out: %v, -der-prefix: %s, -der-dir"
	fmtPrettyOptions = `  -host:       %v 
  -out:        %v
  -no-root:    %v
  -no-server:  %v
  -insecure:   %v
  -der-out:    %v
  -der-prefix: %s
  -der-dir:    %s`
)

func (opts *options) string(format string) string {
	out := strings.Builder{}
	fmt.Fprintf(&out, format,
		opts.hostPort, opts.showOut, opts.noRoot, opts.noServer, opts.insecure, opts.derOut, opts.derPrefix, opts.derDir)
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

	fs, opts := newFlagSetAndOptions(args[0])

	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}
	return opts, nil
}

func usage(args []string) {
	fs, _ := newFlagSetAndOptions(args[0])
	fs.Usage()
}

func newFlagSetAndOptions(name string) (*flag.FlagSet, *options)  {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	opts := &options{}
	fs.StringVar(&opts.hostPort, "host", "", "Hostname plus port")
	fs.StringVar(&opts.derPrefix, "der-prefix", "", "Prefix for the der files. Defaults to <host name>-")
	fs.StringVar(&opts.derDir, "der-dir", ".", "Path to write the der files to. Defaults to the current directory")
	fs.StringVar(&opts.cacerts, "cacerts", "", "Extra cacerts to use to verify the server, pem format")

	fs.BoolVar(&opts.derOut, "der-out", false, "Output der files. The names of the files is <host>-0x.der [false]")
	fs.BoolVar(&opts.showOut, "out", false, "Show pem output [false]")
	fs.BoolVar(&opts.noRoot, "no-root", false, "Omit the root cert in pem or der output [false]")
	fs.BoolVar(&opts.noServer, "no-server", false, "Omit the server cert in pem or der output [false]")
	fs.BoolVar(&opts.showOpts, "show-opts", false, "Show the options [false]")
	fs.BoolVar(&opts.insecure, "insecure", true, "Allow insecure certs [true]")
	return fs, opts
}