package hibp

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HIBPEnforcement string

const (
	STRICT HIBPEnforcement = "strict"
	LOOSE  HIBPEnforcement = "loose"
)

type HIBPSettings struct {
	Enabled     bool
	AppName     string
	Enforcement HIBPEnforcement
	HTTPClient  *http.Client
}

// A packge used for HIBP integration
func sha1Hash(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	return fmt.Sprintf("%X", hasher.Sum(nil))
}

// It is critical that this function is async
// as HIBP can take 1-2 seconds to respond.
// Run this at the beginning of login/registration
// and only listen for the result when required.
func CheckPassword(userAgent string, password string, isTaken chan bool, httpClient *http.Client) {
	hashedPass := sha1Hash(password)

	url := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", hashedPass[:5])
	method := "GET"

	client := httpClient

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		isTaken <- false
		return
	}
	req.Header.Add("user-agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		isTaken <- false
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		isTaken <- false
		return
	}

	isTaken <- checkResults(hashedPass, body)
}

func checkResults(checkHash string, hibpResults []byte) bool {
	checkHash = checkHash[5:]
	scanner := bufio.NewScanner(bytes.NewReader(hibpResults))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue // Skip lines that don't have exactly two parts
		}

		hash, countStr := parts[0], parts[1]
		if countStr == "0" {
			continue
		}

		if strings.EqualFold(hash, checkHash) {
			return true
		}
	}

	return false
}
