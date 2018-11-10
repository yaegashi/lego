package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"github.com/xenolf/lego/certificate"
	"github.com/xenolf/lego/log"
)

func checkFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0700)
	}
	return nil
}

func saveCertRes(certRes *certificate.Resource, c *cli.Context) {
	var domainName string

	// Check filename cli parameter
	if c.GlobalString("filename") == "" {
		// Make sure no funny chars are in the cert names (like wildcards ;))
		domainName = strings.Replace(certRes.Domain, "*", "_", -1)
	} else {
		domainName = c.GlobalString("filename")
	}

	// We store the certificate, private key and metadata in different files
	// as web servers would not be able to work with a combined file.
	certOut := filepath.Join(CertPath(c), domainName+".crt")

	err := checkFolder(filepath.Dir(certOut))
	if err != nil {
		log.Fatalf("Could not check/create path: %v", err)
	}

	err = ioutil.WriteFile(certOut, certRes.Certificate, 0600)
	if err != nil {
		log.Fatalf("Unable to save Certificate for domain %s\n\t%v", certRes.Domain, err)
	}

	issuerOut := filepath.Join(CertPath(c), domainName+".issuer.crt")

	if certRes.IssuerCertificate != nil {
		err = ioutil.WriteFile(issuerOut, certRes.IssuerCertificate, 0600)
		if err != nil {
			log.Fatalf("Unable to save IssuerCertificate for domain %s\n\t%v", certRes.Domain, err)
		}
	}

	if certRes.PrivateKey != nil {
		privOut := filepath.Join(CertPath(c), domainName+".key")

		// if we were given a CSR, we don't know the private key
		err = ioutil.WriteFile(privOut, certRes.PrivateKey, 0600)
		if err != nil {
			log.Fatalf("Unable to save PrivateKey for domain %s\n\t%v", certRes.Domain, err)
		}

		if c.GlobalBool("pem") {
			pemOut := filepath.Join(CertPath(c), domainName+".pem")
			err = ioutil.WriteFile(pemOut, bytes.Join([][]byte{certRes.Certificate, certRes.PrivateKey}, nil), 0600)
			if err != nil {
				log.Fatalf("Unable to save Certificate and PrivateKey in .pem for domain %s\n\t%v", certRes.Domain, err)
			}
		}

	} else if c.GlobalBool("pem") {
		// we don't have the private key; can't write the .pem file
		log.Fatalf("Unable to save pem without private key for domain %s\n\t%v; are you using a CSR?", certRes.Domain, err)
	}

	jsonBytes, err := json.MarshalIndent(certRes, "", "\t")
	if err != nil {
		log.Fatalf("Unable to marshal CertResource for domain %s\n\t%v", certRes.Domain, err)
	}

	metaOut := filepath.Join(CertPath(c), domainName+".json")
	err = ioutil.WriteFile(metaOut, jsonBytes, 0600)
	if err != nil {
		log.Fatalf("Unable to save CertResource for domain %s\n\t%v", certRes.Domain, err)
	}
}

// CertPath gets the path for certificates.
func CertPath(c *cli.Context) string {
	return filepath.Join(c.GlobalString("path"), "certificates")
}
