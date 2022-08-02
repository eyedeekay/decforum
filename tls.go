package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"i2pgit.org/idk/reseed-tools/reseed"
	//"github.com/MDrollette/i2p-tools/reseed"
)

func checkOrNewTLSCert(tlsHost string, tlsCert, tlsKey *string) error {
	//tlsCert := tlsHost + ".crt"
	//tlsKey := tlsHost + ".pem"
	_, certErr := os.Stat(*tlsCert)
	_, keyErr := os.Stat(*tlsKey)
	if certErr != nil || keyErr != nil {
		if certErr != nil {
			fmt.Printf("Unable to read TLS certificate '%s'\n", *tlsCert)
		}
		if keyErr != nil {
			fmt.Printf("Unable to read TLS key '%s'\n", *tlsKey)
		}

		if err := createTLSCertificate(tlsHost); nil != err {
			return err
		}

		*tlsCert = tlsHost + ".crt"
		*tlsKey = tlsHost + ".pem"
	}

	return nil
}

func createTLSCertificate(host string) error {
	fmt.Println("Generating TLS keys. This may take a minute...")
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	tlsCert, err := reseed.NewTLSCertificate(host, priv)
	if nil != err {
		return err
	}

	// save the TLS certificate
	certOut, err := os.Create(host + ".crt")
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", host+".crt", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: tlsCert})
	certOut.Close()
	fmt.Printf("\tTLS certificate saved to: %s\n", host+".crt")

	// save the TLS private key
	privFile := host + ".pem"
	keyOut, err := os.OpenFile(privFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %v", privFile, err)
	}
	secp384r1, err := asn1.Marshal(asn1.ObjectIdentifier{1, 3, 132, 0, 34}) // http://www.ietf.org/rfc/rfc5480.txt
	pem.Encode(keyOut, &pem.Block{Type: "EC PARAMETERS", Bytes: secp384r1})
	ecder, err := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: ecder})
	pem.Encode(keyOut, &pem.Block{Type: "CERTIFICATE", Bytes: tlsCert})

	keyOut.Close()
	fmt.Printf("\tTLS private key saved to: %s\n", privFile)

	// CRL
	crlFile := host + ".crl"
	crlOut, err := os.OpenFile(crlFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", crlFile, err)
	}
	crlcert, err := x509.ParseCertificate(tlsCert)
	if err != nil {
		return fmt.Errorf("Certificate with unknown critical extension was not parsed: %s", err)
	}

	now := time.Now()
	revokedCerts := []pkix.RevokedCertificate{
		{
			SerialNumber:   crlcert.SerialNumber,
			RevocationTime: now,
		},
	}

	crlBytes, err := crlcert.CreateCRL(rand.Reader, priv, revokedCerts, now, now)
	if err != nil {
		return fmt.Errorf("error creating CRL: %s", err)
	}
	_, err = x509.ParseDERCRL(crlBytes)
	if err != nil {
		return fmt.Errorf("error reparsing CRL: %s", err)
	}
	pem.Encode(crlOut, &pem.Block{Type: "X509 CRL", Bytes: crlBytes})
	crlOut.Close()
	fmt.Printf("\tTLS CRL saved to: %s\n", crlFile)

	return nil
}
