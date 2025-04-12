package browser

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/krbreyn/gemcat/data"
	"github.com/krbreyn/gemcat/tofu"
)

func FetchGemini(url *url.URL, doCache bool) (status, body string, err error) {

ifRedirect:
	if url.Scheme != "gemini" {
		return "", "", fmt.Errorf("only gemini connections are handled, got %s", url.String())
	}

	host, path := url.Host, url.Path

	tlsDialer := &tls.Dialer{
		NetDialer: &net.Dialer{
			Timeout: 7 * time.Second,
		},
		Config: &tls.Config{
			InsecureSkipVerify: true,
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				return tofu.HandleTOFU(rawCerts, host)
			},
		},
	}

	if doCache {
		isStale, err := data.IsCacheStale(url, time.Hour*24)
		if err != nil {
			return "", "", fmt.Errorf("cache error: %w\n", err)
		}

		if !isStale {
			content, err := data.LoadFromCache(url)
			if err != nil {
				return "", "", fmt.Errorf("cache error: %w\n", err)
			} else {
				// fmt.Println("cache hit")
				return "20 [cache hit]", string(content), nil
			}
		}
	}

	addr := net.JoinHostPort(host, "1965")
	conn, err := tlsDialer.Dial("tcp", addr)
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
		return status, "", fmt.Errorf("weird status err: %v", err)
	}

	// TODO integrate this function with browser to update history properly
	if status_no >= 30 && status_no <= 39 {
		url, err = url.Parse(strings.Fields(status)[1])
		if err != nil {
			return status, "", fmt.Errorf("redirect url parse error: %w", err)
		}
		// TODO use an OutputObject?
		// fmt.Printf("Redirect: %s\r\n", new_url.String())
		goto ifRedirect
	}

	if status_no < 20 && status_no > 29 {
		return status, "", fmt.Errorf("status was not 2x but was %d", status_no)
	}

	var b strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		b.WriteString(line)
	}

	content := b.String()

	err = data.CacheGemFile(url, []byte(content))
	if err != nil {
		return "", "", fmt.Errorf("cache err: %w", err)
	}

	return status, content, nil
}
