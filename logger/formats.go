package logger

var logs = map[LogFormats]LogFormat{
	INVALID_METHOD: {
		Regex:      "'(\\d+.\\d+.\\d+.\\d+)' used an incorrect method on '(\\w+)' expected '(\\w+)' got '(\\w+)'",
		RegexOrder: []string{"srcip", "path", "extra_data", "extra_data"},
		FmtPattern: "'%s' used an incorrect method on '%s' expected '%s' got '%s'",
	},
	BAD_REQUEST: {
		Regex:      "'(\\d+.\\d+.\\d+.\\d+)' used a bad request on '(\\w+)'",
		RegexOrder: []string{"srcip", "url"},
		FmtPattern: "'%s' used a bad request on '%s'",
	},
	LOGIN: {
		Regex:      "User '(\\w+)' logged in from destination '(\\d+.\\d+.\\d+.\\d+)'",
		RegexOrder: []string{"dstuser", "srcip"},
		FmtPattern: "User '%s' logged in from destination '%s'",
	},
	LOGIN_FAIL_BAD_PASSWORD: {
		Regex:      "User '(\\w+)' presented a bad password from destination '(\\d+.\\d+.\\d+.\\d+)'",
		RegexOrder: []string{"dstuser", "srcip"},
		FmtPattern: "User '%s' presented a bad password from destination '%s'",
	},
	LOGIN_INVALID_USERNAME: {
		Regex:      "'(\\d+\\.\\d+\\.\\d+\\.\\d+)' attempted to login to an invalid username '(\\w+)'",
		RegexOrder: []string{"srcip", "dstuser"},
		FmtPattern: "'%s' attempted to login to an invalid username '%s'",
	},
	INVITE_CODE_MALFORMED: {
		Regex:      "'(\\d+.\\d+.\\d+.\\d+)' used an invalid invite code '(\\w+)'",
		RegexOrder: []string{"srcip", "extra_data"},
		FmtPattern: "'%s' used an invalid invite code '%s'",
	},
	INVITE_CODE_FAKE_VIEW: {
		Regex:      "'(\\d+.\\d+.\\d+.\\d+)' used a fake or expired invite code to try and view regitration '(\\w+)'",
		RegexOrder: []string{"srcip", "extra_data"},
		FmtPattern: "'%s' used a fake or expired invite code to try and view registration '%s'",
	},
	INVITE_CODE_FAKE: {
		Regex:      "'(\\d+.\\d+.\\d+.\\d+)' used a fake invite code '(\\w+)'",
		RegexOrder: []string{"srcip", "extra_data"},
		FmtPattern: "'%s' used a fake invite code '%s'",
	},
	INVITE_CODE_INCORRECT_EMAIL: {
		Regex:      "'(\\d+.\\d+.\\d+.\\d+)' used an incorrect email on attached invite '(\\w+)'",
		RegexOrder: []string{"srcip", "extra_data"},
		FmtPattern: "'%s' used an incorrect email on attached invite '%s'",
	},
	INVITE_CODE_USED: {
		Regex:      "'(\\d+.\\d+.\\d+.\\d+)' registered '(\\w+)' using invite code '(\\w+)'",
		RegexOrder: []string{"srcip", "dstuser", "extra_data"},
		FmtPattern: "'%s' registered '%s' using invite code '%s'",
	},
	REGISTRATION_SUCCESS: {
		Regex:      "'(\\d+\\.\\d+\\.\\d+\\.\\d+)' registered '([^']+)'",
		RegexOrder: []string{"srcip", "dstuser"},
		FmtPattern: "'%s' registered '%s'",
	},
}

type LogFormats string

const (
	LOGIN                       LogFormats = "login"
	LOGIN_FAIL_BAD_PASSWORD     LogFormats = "login_bad_password"
	LOGIN_INVALID_USERNAME      LogFormats = "login_invalid_username"
	INVALID_METHOD              LogFormats = "invalid_method"
	BAD_REQUEST                 LogFormats = "bad_request"
	INVITE_CODE_FAKE_VIEW       LogFormats = "malformed_invite_code_view"
	INVITE_CODE_MALFORMED       LogFormats = "malformed_invite_code"
	INVITE_CODE_FAKE            LogFormats = "fake_invite_code"
	INVITE_CODE_INCORRECT_EMAIL LogFormats = "incorect_email_invite_code"
	INVITE_CODE_USED            LogFormats = "used_invite_code"
	REGISTRATION_SUCCESS        LogFormats = "registration_success"
)
