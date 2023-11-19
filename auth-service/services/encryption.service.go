package services

import (
	"auth-service/errors"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type EncryptionService struct {
	SecretKey string
}

type TokenData struct {
	UserID         string
	ExpirationDate time.Time
}

func (e EncryptionService) GenerateToken(userID string) (string, *errors.ErrorStruct) {
	expirationDate := time.Now().Add(time.Hour)

	data := TokenData{
		UserID:         userID,
		ExpirationDate: expirationDate,
	}

	tokenJSON, err := json.Marshal(data)

	if err != nil {
		return "", errors.NewError(err.Error(), 500)
	}
	h := hmac.New(sha256.New, []byte(e.SecretKey))
	h.Write(tokenJSON)
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	token := fmt.Sprintf("%s.%s", base64.URLEncoding.EncodeToString(tokenJSON), signature)
	return token, nil

}

func (e EncryptionService) ValidateToken(token string) (string, *errors.ErrorStruct) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", errors.NewError("Invalid token format", 500)
	}

	encodedData, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", errors.NewError(err.Error(), 500)
	}

	h := hmac.New(sha256.New, []byte(e.SecretKey))
	h.Write(encodedData)
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(parts[1]), []byte(expectedSignature)) {
		return "", errors.NewError("Invalid token signature", 401)
	}

	var data TokenData
	err = json.Unmarshal(encodedData, &data)
	if err != nil {
		return "", errors.NewError(err.Error(), 500)
	}

	if time.Now().After(data.ExpirationDate) {
		return "", errors.NewError("Token has expired", 401)
	}

	return data.UserID, nil
}
