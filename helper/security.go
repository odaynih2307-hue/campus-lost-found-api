package helper

import (
	"campus-lost-found-api/app/model"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
func RandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func SHA256Hex(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func IssueAccessToken(user model.User, secret, issuer string, minutes int) (string, error) {
	now := time.Now()
	c := jwt.MapClaims{"sub": user.ID, "username": user.Username, "role": user.Role, "iss": issuer, "iat": now.Unix(), "exp": now.Add(time.Duration(minutes) * time.Minute).Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(secret))
}
func ParseAccessToken(token, secret, issuer string) (model.AuthUser, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, Unauthorized("token tidak valid")
		}
		return []byte(secret), nil
	}, jwt.WithIssuer(issuer))
	if err != nil || !t.Valid {
		return model.AuthUser{}, Unauthorized("token tidak valid")
	}
	c, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return model.AuthUser{}, Unauthorized("token tidak valid")
	}
	sub, ok := c["sub"].(float64)
	if !ok {
		return model.AuthUser{}, Unauthorized("token tidak valid")
	}
	username, _ := c["username"].(string)
	role, _ := c["role"].(string)
	return model.AuthUser{UserID: int64(sub), Username: username, Role: role}, nil
}
