package verificationcodes

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrVerificationRateLimited       = errors.New("verification rate limited")
	ErrVerificationMethodUnsupported = errors.New("verification method unsupported")
	ErrVerificationExpired           = errors.New("verification expired")
	ErrVerificationCodeDoesNotBelong = errors.New("verification code does not belong to the user")
)

type EmailVerificationSender interface {
	SendVerificationEmail(ctx context.Context, userEmail string, userID string, token string) error
}

type PhoneVerificationSender interface {
	SendVerificationPhone(ctx context.Context, userPhone string, userID string, token string) error
}

type VerificationSender interface {
	EmailVerificationSender
}

type stubSender struct{}

func (s *stubSender) SendVerificationEmail(ctx context.Context, userEmail string, userID string, token string) error {
	fmt.Printf("EMAIL > %s (%s) -> %s\n", userEmail, userID, token)
	return nil
}

func (s *stubSender) SendVerificationPhone(ctx context.Context, userPhone string, userID string, token string) error {
	fmt.Printf("PHONE > %s (%s) -> %s\n", userPhone, userID, token)
	return nil
}

type Verifier interface {
	DeleteCode(ctx context.Context, method verifierMethod, userID string) error
	CheckCode(ctx context.Context, method verifierMethod, userID string, code string) (bool, error)
	SendVerification(ctx context.Context, lsu users.LocksmithUserInterface, method verifierMethod, forValue string) error
}

type verifier struct {
	db     database.DatabaseAccessor
	sender VerificationSender
}

func NewVerifier(db database.DatabaseAccessor, sender VerificationSender) Verifier {
	if sender == nil {
		sender = &stubSender{}
	}
	return &verifier{db: db, sender: sender}
}

var errNotFound = errors.New("verification code not found")

type verifierMethod string

const (
	VerifierMethod_EMAIL        verifierMethod = "email"
	VerifierMethod_PHONE_NUMBER verifierMethod = "phone"
)

type code struct {
	ForUserID string
	Method    verifierMethod
	ForValue  string
	Code      string
	ExpiresAt time.Time
}

func (e *verifier) getCode(ctx context.Context, method verifierMethod, userID string) (*code, error) {
	rawCode, found := e.db.FindOne("email_verifications", map[string]any{
		"user_id": userID,
		"method":  string(method),
	})
	if !found {
		return nil, errNotFound
	}

	c := code{
		ForUserID: userID,
		Method:    method,
		ExpiresAt: time.Unix(rawCode.(map[string]any)["expires_at"].(int64), 0),
		ForValue:  rawCode.(map[string]any)["value"].(string),
		Code:      rawCode.(map[string]any)["code"].(string),
	}

	_, err := e.db.UpdateOne("email_verifications", map[string]any{
		"_id": rawCode.(map[string]any)["_id"].(primitive.ObjectID),
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.SET: {
			"expires_at": time.Now().UTC().Add(1 * time.Hour).Unix(),
		},
	})

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (e *verifier) generateCode(ctx context.Context, method verifierMethod, userID string, forValue string) (*code, error) {
	c := code{
		ForUserID: userID,
		Method:    method,
		ForValue:  forValue,
		Code:      rand.Text(),
		ExpiresAt: time.Now().UTC().Add(1 * time.Hour),
	}

	_, err := e.db.InsertOne("email_verifications", map[string]any{
		"user_id":    c.ForUserID,
		"code":       c.Code,
		"method":     string(c.Method),
		"value":      c.ForValue,
		"expires_at": c.ExpiresAt.Unix(),
	})
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (e *verifier) DeleteCode(ctx context.Context, method verifierMethod, userID string) error {
	_, err := e.db.DeleteOne("email_verifications", map[string]any{
		"user_id": userID,
		"method":  string(method),
	})
	return err
}

func (e *verifier) CheckCode(ctx context.Context, method verifierMethod, userID string, code string) (bool, error) {
	c, err := e.getCode(ctx, method, userID)
	if err != nil {
		if errors.Is(err, errNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get code during checkCode: %w", err)
	}

	isGoodCode := c.Code == code

	if !isGoodCode {
		return false, nil
	}

	if time.Now().After(c.ExpiresAt) {
		return false, ErrVerificationExpired
	}

	if c.ForUserID != userID {
		return false, ErrVerificationCodeDoesNotBelong
	}

	return true, nil
}

func (e *verifier) SendVerification(ctx context.Context, lsu users.LocksmithUserInterface, method verifierMethod, forValue string) error {
	c, err := e.getCode(ctx, method, lsu.GetID())
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return fmt.Errorf("failed to get code during sendVerificationEmail: %w", err)
		}

		// code not found, generate a new one
		c, err = e.generateCode(ctx, method, lsu.GetID(), forValue)
		if err != nil {
			return fmt.Errorf("failed to generate email verification code: %w", err)
		}
	}

	var sendErr error
	switch method {
	case VerifierMethod_EMAIL:
		sendErr = e.sender.SendVerificationEmail(ctx, lsu.GetEmail(), lsu.GetID(), c.Code)
	case VerifierMethod_PHONE_NUMBER:
		asPhone, ok := e.sender.(PhoneVerificationSender)
		if !ok {
			return ErrVerificationMethodUnsupported
		}
		sendErr = asPhone.SendVerificationPhone(ctx, c.ForValue, lsu.GetID(), c.Code)
	default:
		return ErrVerificationMethodUnsupported
	}
	if sendErr != nil {
		if errors.Is(err, ErrVerificationRateLimited) {
			return ErrVerificationRateLimited // pass it through directly.
		}
		if errors.Is(err, ErrVerificationMethodUnsupported) {
			return ErrVerificationMethodUnsupported // pass it through directly.
		}
		return fmt.Errorf("failed to send verification email: %w", sendErr)
	}

	return nil
}
