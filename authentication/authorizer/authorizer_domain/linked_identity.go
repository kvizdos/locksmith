package authorizer_domain

import "time"

type LinkedIdentity struct {
	Provider string    `json:"provider"`
	Subject  string    `json:"subject"`
	UserID   string    `json:"user_id"`
	Issuer   string    `json:"issuer"`
	LinkedAt time.Time `json:"linked_at"`
}

func (l LinkedIdentity) ToMap() map[string]any {
	return map[string]any{
		"provider":  l.Provider,
		"subject":   l.Subject,
		"user_id":   l.UserID,
		"linked_at": l.LinkedAt.Unix(),
		"issuer":    l.Issuer,
	}
}

func LinkedIdentityFromMap(m map[string]any) LinkedIdentity {
	return LinkedIdentity{
		Provider: m["provider"].(string),
		Issuer:   m["issuer"].(string),
		Subject:  m["subject"].(string),
		UserID:   m["user_id"].(string),
		LinkedAt: time.Unix(m["linked_at"].(int64), 0),
	}
}
