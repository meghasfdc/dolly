package metrics

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"time"

	"github.com/go-phorce/pkg/certutil"
	"github.com/go-phorce/pkg/metrics/tags"
)

// PublishCertExpirationInDays publish cert expiration time in Days
func PublishCertExpirationInDays(c *x509.Certificate, typ string) float32 {
	metricKey := []string{"cert", "expiry", "days",
		tags.Separator,
		"CN", c.Subject.CommonName,
		"Serial", c.SerialNumber.String(),
		"SKI", certutil.GetSubjectKeyID(c),
		"type", typ,
	}
	expiresIn := c.NotAfter.Sub(time.Now().UTC())
	expiresInDays := float32(expiresIn) / float32(time.Hour*24)
	SetGauge(metricKey, expiresInDays)
	return expiresInDays
}

// PublishCRLExpirationInDays publish CRL expiration time in Days
func PublishCRLExpirationInDays(c *pkix.CertificateList, issuerID string, issuer *x509.Certificate) float32 {
	PublishCertExpirationInDays(issuer, "issuer")

	metricKey := []string{"crl", "expiry", "days",
		tags.Separator,
		"CN", issuer.Subject.CommonName,
		"Serial", issuer.SerialNumber.String(),
		"SKI", certutil.GetSubjectKeyID(issuer),
		"IssuerID", issuerID,
	}
	expiresIn := c.TBSCertList.NextUpdate.Sub(time.Now().UTC())
	expiresInDays := float32(expiresIn) / float32(time.Hour*24)
	SetGauge(metricKey, expiresInDays)
	return expiresInDays
}