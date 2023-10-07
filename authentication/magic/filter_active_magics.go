package magic

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// Useful if you've already queried the user and just need a list
// of active magics back. Used in: the end of Validate() + ValidateSession()
func FilterActive(activeMagics chan MagicAuthentications, magics MagicAuthentications, manuallyExpireTokenID ...string) {
	doExpireCode := ""

	if len(manuallyExpireTokenID) > 0 {
		debased, _ := base64.StdEncoding.DecodeString(manuallyExpireTokenID[0])
		tokenParts := strings.Split(string(debased), ":")
		hasher := sha256.New()
		hasher.Write([]byte(tokenParts[1]))
		hashedToken := fmt.Sprintf("%x", hasher.Sum(nil))

		doExpireCode = hashedToken
	}

	output := MagicAuthentications{}
	for _, m := range magics {
		if m.ExpiresAt <= time.Now().UTC().Unix() || m.Code == doExpireCode {
			continue
		}

		output = append(output, m)
	}

	activeMagics <- output
}
