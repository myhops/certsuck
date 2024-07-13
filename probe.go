package certsuck

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"slices"
	"text/template"
)

// Probe collects certs.
type Probe struct {
	insecure bool
	rootCAs *x509.CertPool
}

// Chains contains the result of a connection probe to a server.
type Chains struct {
	// Verified contains the verified chains.
	// In most cases the entry 0 contains the peer certificate chain
	// and entry 1 the certificates including the root and the server.
	Verified [][]*x509.Certificate
	// Peer contains the peer certificate chain.
	Peer []*x509.Certificate
	// Longest contains the longest found chain.
	Longest []*x509.Certificate
	// LongestName contains the name of the chain that was the longest.
	// It is either Peer or Verified with its index.
	LongestName string
}

// Option is an Option for the New function.
type Option func(probe *Probe)

// WithInsecure allows unverified chains.
// The default for this option is true.
func WithInsecure(insecure ...bool) Option {
	return func(p *Probe) {
		p.insecure = true
		if len(insecure) > 0 {
			p.insecure = insecure[0]
		}
	}
}

// New returns a new Probe with the given option set.
func New(opts ...Option) *Probe {
	probe := &Probe{}
	for _, opt := range opts {
		opt(probe)
	}
	return probe
}

func WithRootCAs(rootCAs *x509.CertPool) Option {
	return func(p *Probe) {
		p.rootCAs = rootCAs
	}
}

// CollectCerts collect the certs from the server at hostPort and returns the found chains.
func (p *Probe) CollectCerts(hostPort string, opts ...Option) (*Chains, error) {
	return p.collectChains(hostPort)
}

func (p *Probe) collectChains(hostPort string) (*Chains, error) {
	var res = &Chains{}

	verifyConnectionCallback := func(state tls.ConnectionState) error {
		res.Verified = slices.Clone(state.VerifiedChains)
		res.Peer = slices.Clone(state.PeerCertificates)
		return nil
	}

	// Create a config for the callback.
	tlsCfg := tls.Config{
		InsecureSkipVerify: p.insecure,
		VerifyConnection:   verifyConnectionCallback,
		RootCAs: p.rootCAs,
	}
	// Connect to the host
	conn, err := tls.Dial("tcp", hostPort, &tlsCfg)
	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()

	res.setLongest()

	return res, nil
}

func (c *Chains) setLongest() {
	const (
		verified = "Verified"
		peer     = "Peer"
	)
	for i, cc := range c.Verified {
		if len(cc) > len(c.Longest) {
			c.LongestName = fmt.Sprintf("%s %d", verified, i)
			c.Longest = cc
		}
	}
	if len(c.Peer) > len(c.Longest) {
		c.LongestName = peer
		c.Longest = c.Peer
	}
}

// String returns a formatted string.
func (c *Chains) String() string {
	// return toString(c)
	return toStringTemplate(c)
}

func toStringTemplate(c *Chains) string {
	s, err := c.FormatTemplate(stringTemplate)
	if err != nil {
		return err.Error()
	}
	return s
}

// DefaultTemplate is the template that Chains.String uses.
// It is also shown as an example for modification.
const DefaultTemplate = `
{{- range $i, $item := .Verified -}}
{{- println "Verified" $i -}}
	{{- range $j, $val  := $item -}}
		{{- printf "%2d  " $j -}}Subject: {{ .Subject -}}{{- println -}}
		{{- print "    " -}}Issuer:  {{ .Issuer }}{{- println -}}
	{{- end }}
{{- end -}}

{{- if len .Peer | lt 0 -}}
{{- println "Peer" -}}
	{{- range $j, $val  := .Peer -}}
		{{- printf "%2d  " $j -}}Subject: {{ .Subject -}}{{- println -}}
		{{- print "    " -}}Issuer:  {{ .Issuer }}{{- println -}}
	{{- end }}
{{- end -}}
`

// The template used by String.
// You can override it with SetStringTemplate.
var stringTemplate = DefaultTemplate

// SetStringTemplate sets the template used by Chains.String.
// When this function returns an error, Chain.Strings will not work.
func SetStringTemplate(tpl string) error {
	t := template.New("string")
	if _, err := t.Parse(tpl); err != nil {
		return fmt.Errorf("parsing StringTemplate failed: %w", err)
	}
	return nil
}

// MustSetStringTemplate sets the template used by Chains.String.
// It panics if parsing the template fails.
func MustSetStringTemplate(tpl string) {
	if err := SetStringTemplate(tpl); err != nil {
		panic(err.Error())
	}
}

// FormatTemplate produces a string with the chain information.
// String uses this method with the DefaultTemplate for tpl.
func (c *Chains) FormatTemplate(tpl string) (string, error) {
	t := template.New("string")
	t, err := t.Parse(tpl)
	if err != nil {
		return "", err
	}
	w := bytes.Buffer{}
	if err := t.Execute(&w, c); err != nil {
		return "", err
	}
	return w.String(), nil
}
