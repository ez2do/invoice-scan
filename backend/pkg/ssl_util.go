package pkg

import (
	"crypto/x509"
	"encoding/pem"
	"go.uber.org/multierr"
)

func ParseX509CertsPem(pemCerts []byte) ([]*x509.Certificate, error) {
	var merr error
	var certs []*x509.Certificate
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			merr = multierr.Append(merr, err)
			continue
		}

		certs = append(certs, cert)
	}

	return certs, merr
}
