package register

import "encoding/json"

type registrationResponse struct {
	Error        string `json:"error,omitempty"`
	ConfirmEmail bool   `json:"confirmEmail,omitempty"`
	EmailBlocked bool   `json:"rejectEmail,omitempty"`
	DidYouMean   string `json:"didYouMean,omitempty"`
}

func (r registrationResponse) Marshal() []byte {
	js, _ := json.Marshal(r)
	return js
}

func (r *registrationResponse) Unmarshal(err []byte) {
	json.Unmarshal(err, r)
}
