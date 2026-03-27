package core

import (
	"errors"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var (
	ErrInvalidTOTP = errors.New("invalid TOTP code")
	ErrInvalidKey  = errors.New("invalid TOTP key")
)

type TOTPManager struct {
	issuer    string
	digits    otp.Digits
	algorithm otp.Algorithm
}

func NewTOTPManager(issuer string) *TOTPManager {
	return &TOTPManager{
		issuer:    issuer,
		digits:    otp.DigitsSix,
		algorithm: otp.AlgorithmSHA1,
	}
}

func (t *TOTPManager) GenerateSecret(accountName string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      t.issuer,
		AccountName: accountName,
		Digits:      t.digits,
		Algorithm:   t.algorithm,
		SecretSize:  20,
	})
	if err != nil {
		return "", ErrInvalidKey
	}
	return key.Secret(), nil
}

func (t *TOTPManager) ValidateCode(secret, code string) error {
	valid, err := totp.Validate(code, secret)
	if err != nil {
		return ErrInvalidTOTP
	}
	if !valid {
		return ErrInvalidTOTP
	}
	return nil
}

func (t *TOTPManager) GenerateProvisioningURI(secret, accountName string) string {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      t.issuer,
		AccountName: accountName,
		Digits:      t.digits,
		Algorithm:   t.algorithm,
		Secret:      []byte(secret),
	})
	return key.URL()
}

func (t *TOTPManager) GetTimeRemaining() int {
	period := 30
	now := time.Now().Unix()
	remaining := period - int(now%int64(period))
	return remaining
}
