package util

import (
	"file/internal/biz"
	"file/internal/conf"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// CertDir returns the absolute path to the cert directory.
func CertDir() string {
	_, filename, _, _ := runtime.Caller(0)
	// internal/data → server root → cert/
	root := filepath.Join(filepath.Dir(filename), "..", "..")
	return filepath.Join(root, "cert")
}

// NewPublicPEM reads the RSA public key PEM from cert/<kid>-public.pem.
func NewPublicPEM(c *conf.Data) (biz.PublicPEM, error) {
	filename := fmt.Sprintf("%s-public.pem", c.Auth.Kid)
	data, err := os.ReadFile(filepath.Join(CertDir(), filename))
	return biz.PublicPEM(data), err
}
