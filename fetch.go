package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func Fetch(host, path string) (status, body string, err error) {
ifRedirect:
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return handleTOFU(rawCerts, host)
		},
	}

	fmt.Printf("Connecting to gemini://%s/%s\r\n", host, path)
	addr := net.JoinHostPort(host, "1965")
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return "", "", fmt.Errorf("TLS connection failed: %v", err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "gemini://%s/%s\r\n", host, path)

	reader := bufio.NewReader(conn)
	status, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read response:", err)
	}

	status_no, err := strconv.Atoi(strings.Fields(status)[0])
	if err != nil {
		return "", "", fmt.Errorf("weird status err: %v", err)
	}
	// TODO integrate this function with browser to update history properly
	if status_no == 30 || status_no == 31 {
		new_url := strings.Fields(status)[1]
		new_url = strings.TrimPrefix(new_url, "gemini://")
		host, path = getHostPath(new_url)
		fmt.Printf("Redirect: gemini://%s/%s\r\n", host, path)
		goto ifRedirect
	}
	if status_no != 20 {
		return "", "", fmt.Errorf("status was not 20 but was %d, status: %s", status_no, status)
	}

	var b strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		b.WriteString(line)
	}

	return status, b.String(), nil
}
