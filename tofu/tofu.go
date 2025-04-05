package tofu

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const knownHostsFile = ".gemini-known-hosts"

func HandleTOFU(rawCerts [][]byte, hostname string) error {
	cert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		return err
	}

	fingerprint := certFingerprint(cert)
	known, err := readKnownFingerprint(hostname)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if known == "" {
		//fmt.Println("[TOFU] First time seeing", hostname, "- trusting cert:", fingerprint)
		return saveFingerprint(hostname, fingerprint)
	}

	if known != fingerprint {
		return fmt.Errorf("[TOFU] Certificate mismatch for %s! Expected %s, got %s", hostname, known, fingerprint)
	}

	return nil
}

func certFingerprint(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(sum[:])
}

func knownHostsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Can't find home directory:", err)
	}
	return filepath.Join(home, knownHostsFile)
}

func readKnownFingerprint(hostname string) (string, error) {
	file := knownHostsPath()
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	lines := strings.SplitSeq(string(data), "\n")
	for line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[0] == hostname {
			return parts[1], nil
		}
	}
	return "", nil
}

func saveFingerprint(hostname, fingerprint string) error {
	file := knownHostsPath()
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s %s\n", hostname, fingerprint)
	return err
}
