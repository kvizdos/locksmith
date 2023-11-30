package fingerprint

import (
	"encoding/json"
	"testing"
)

func TestFingerprintJSONParseDoesNotReceiveIPAddressOrQuickHash(t *testing.T) {
	input := `{"screen":"483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e","timezone":"America/New_York","hardwareConcurrency":10,"deviceMemory":8,"canvas":"7676cbb516459b747df10789fc5bc96c57eebdb0c34d30319aaddabd4441f06d","lang":"en-US","webgl":"ee28f782c28905031900aaf78ed484ee9f2ac58e6def0a9d2a3bc7e87c6e4303","touch":"d18b034841f555937a0f971ad47da6514c68b9c71ea8f75eb9fd29290dd4e5f3","battery":true,"platform":"macOS","userAgent":"e463666c37a4a0156adaf7ef8e685f340f69e0adb04d9df79bad3cc9ff23cfd0","windowSize":"20af4c47a93f572daa7a2ee166100693141df25307b20a128cd10e0c74e9ba32","dnt":false,"devices":"e6dc3b1fb56af9a32897959f9d02be8eb018a8209e06963f82e80436b0357227","ipAddress":"192.158.23.11"}`

	var parsedFingerprint Fingerprint
	json.Unmarshal([]byte(input), &parsedFingerprint)

	if parsedFingerprint.IPAddress != "" {
		t.Errorf("got IP Address.. bad!!!")
	}

	input = `{"screen":"483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e","timezone":"America/New_York","hardwareConcurrency":10,"deviceMemory":8,"canvas":"7676cbb516459b747df10789fc5bc96c57eebdb0c34d30319aaddabd4441f06d","lang":"en-US","webgl":"ee28f782c28905031900aaf78ed484ee9f2ac58e6def0a9d2a3bc7e87c6e4303","touch":"d18b034841f555937a0f971ad47da6514c68b9c71ea8f75eb9fd29290dd4e5f3","battery":true,"platform":"macOS","userAgent":"e463666c37a4a0156adaf7ef8e685f340f69e0adb04d9df79bad3cc9ff23cfd0","windowSize":"20af4c47a93f572daa7a2ee166100693141df25307b20a128cd10e0c74e9ba32","dnt":false,"devices":"e6dc3b1fb56af9a32897959f9d02be8eb018a8209e06963f82e80436b0357227","quickHash":"192.158.23.11"}`

	json.Unmarshal([]byte(input), &parsedFingerprint)

	if parsedFingerprint.QuickHash != "" {
		t.Errorf("got quick hash.. bad!!!")
	}
}
