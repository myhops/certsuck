package certsuck

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func ParseCerts(data []byte)([]*x509.Certificate, error) {
	var res []*x509.Certificate
	const certType = "CERTIFICATE"
	for pb, rest := pem.Decode(data); pb != nil; pb, rest = pem.Decode(rest) {
		if pb.Type != certType {
			continue
		}
		crt, err := x509.ParseCertificate(pb.Bytes)
		if err != nil {
			continue
		}
		res = append(res, crt)
	}
	return res, nil
}

func ReadCertsFile(filename string) ([]*x509.Certificate, error) {
	// Load file.
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading cert file: %w", err)
	}

	return ParseCerts(data)
}

