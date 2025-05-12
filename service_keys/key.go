package service_keys

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceKeyCredentials struct {
	ClientID string `json:"client_id"`
	SecretID string `json:"secret_id"`
	Secret   string `json:"client_secret"`
}

type ServiceKey struct {
	ID           primitive.ObjectID `json:"id"`
	ClientID     string             `json:"client_id"`
	Name         string             `json:"name"`
	Scopes       []string           `json:"scopes"`
	Secrets      []ServiceSecret    `json:"-"`
	AllowListIPs []string           `json:"ips"`
	CreatedAt    time.Time          `json:"created_at"`
	CreatedBy    string             `json:"created_by"`
	RevokedAt    *time.Time         `json:"revoked_at,omitzero"`
	RevokedBy    string             `json:"revoked_by,omitempty"`
}

func (k ServiceKey) ToMap() map[string]any {
	secrets := []map[string]any{}
	for _, secret := range k.Secrets {
		secrets = append(secrets, secret.ToMap())
	}

	return map[string]any{
		"_id":        k.ID,
		"client_id":  k.ClientID,
		"name":       k.Name,
		"scopes":     k.Scopes,
		"secrets":    secrets,
		"ips":        k.AllowListIPs,
		"created":    k.CreatedAt.Unix(),
		"created_by": k.CreatedBy,
		"revoked_at": (func() any {
			if k.RevokedAt != nil {
				return k.RevokedAt.Unix()
			}
			return nil
		})(),
		"revoked_by": (func() any {
			if k.RevokedBy != "" {
				return k.RevokedBy
			}
			return nil
		})(),
	}
}

func ServiceKeyFromMap(rawInput any) ServiceKey {
	input := rawInput.(map[string]interface{})
	rawSecrets := input["secrets"].([]map[string]any)
	secrets := []ServiceSecret{}
	for _, raw := range rawSecrets {
		secret := ServiceSecret{}.FromMap(raw)
		secrets = append(secrets, secret)
	}

	rawScopes := input["scopes"].([]string)
	scopes := []string{}
	for _, s := range rawScopes {
		scopes = append(scopes, s)
	}

	revokedBy := ""
	var revokedAt *time.Time
	if v, ok := input["revoked_by"].(string); ok {
		revokedBy = v
	}

	if v, ok := input["revoked_at"].(int64); ok {
		ra := time.Unix(v, 0)
		revokedAt = &ra
	}

	ips := []string{}
	if ips, ok := input["ips"].([]interface{}); ok {
		for _, ip := range ips {
			ips = append(ips, ip)
		}
	}

	return ServiceKey{
		ID:           input["_id"].(primitive.ObjectID),
		ClientID:     input["client_id"].(string),
		Name:         input["name"].(string),
		CreatedAt:    time.Unix(input["created"].(int64), 0),
		CreatedBy:    input["created_by"].(string),
		AllowListIPs: ips,
		Scopes:       scopes,
		Secrets:      secrets,
		RevokedAt:    revokedAt,
		RevokedBy:    revokedBy,
	}
}

type ServiceSecret struct {
	ID         string     `json:"id"`
	SecretHash string     `json:"secret_hash"`
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  string     `json:"created_by"`
	LastUsedAt *time.Time `json:"last_used_at,omitzero"`
	RevokedAt  *time.Time `json:"revoked_at,omitzero"`
	RevokedBy  string     `json:"revoked_by,omitempty"`
}

func (s ServiceSecret) ToMap() map[string]any {
	return map[string]any{
		"id":          s.ID,
		"secret_hash": s.SecretHash,
		"created_at":  s.CreatedAt.Unix(),
		"created_by":  s.CreatedBy,
		"last_used_at": func() any {
			if s.LastUsedAt != nil && !s.LastUsedAt.IsZero() {
				return s.LastUsedAt.Unix()
			}
			return nil
		}(),
		"revoked_at": func() any {
			if s.RevokedAt != nil {
				return s.RevokedAt.Unix()
			}
			return nil
		}(),
		"revoked_by": func() any {
			if s.RevokedBy != "" {
				return s.RevokedBy
			}
			return nil
		}(),
	}
}

func (s ServiceSecret) FromMap(input map[string]any) ServiceSecret {
	var revoked *time.Time
	if input["revoked_at"] != nil {
		rv := time.Unix(input["revoked_at"].(int64), 0)
		revoked = &rv
	}

	var lastUsed *time.Time
	if input["last_used_at"] != nil {
		lu := time.Unix(input["last_used_at"].(int64), 0)
		lastUsed = &lu
	}

	var revokedBy string
	if v, ok := input["revoked_by"].(string); ok {
		revokedBy = v
	}

	return ServiceSecret{
		ID:         input["id"].(string),
		SecretHash: input["secret_hash"].(string),
		CreatedBy:  input["created_by"].(string),
		CreatedAt:  time.Unix(input["created_at"].(int64), 0),
		LastUsedAt: lastUsed,
		RevokedAt:  revoked,
		RevokedBy:  revokedBy,
	}
}
