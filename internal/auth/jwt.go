package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	priv *rsa.PrivateKey
	pub  *rsa.PublicKey
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTService() (*JWTService, error) {
	accessMin := envInt("ACCESS_TOKEN_TTL_MIN", 15)
	refreshH  := envInt("REFRESH_TOKEN_TTL_HOURS", 24*30)

	priv, pub, err := loadOrGenerateKeys()
	if err != nil {
		return nil, err
	}
	return &JWTService{
		priv:       priv,
		pub:        pub,
		accessTTL:  time.Duration(accessMin) * time.Minute,
		refreshTTL: time.Duration(refreshH) * time.Hour,
	}, nil
}

func envInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil { return i }
	}
	return def
}

func loadOrGenerateKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privB64 := os.Getenv("JWT_PRIVATE_KEY_BASE64")
	pubB64  := os.Getenv("JWT_PUBLIC_KEY_BASE64")
	if privB64 != "" && pubB64 != "" {
		privPEM, _ := base64.StdEncoding.DecodeString(privB64)
		pubPEM,  _ := base64.StdEncoding.DecodeString(pubB64)
		priv, err := parseRSAPrivateKeyFromPEM(privPEM)
		if err != nil { return nil, nil, err }
		pub,  err := parseRSAPublicKeyFromPEM(pubPEM)
		if err != nil { return nil, nil, err }
		return priv, pub, nil
	}
	// Dev fallback: generate ephemeral
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil { return nil, nil, err }
	return priv, &priv.PublicKey, nil
}

func parseRSAPrivateKeyFromPEM(p []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(p)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private pem")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func parseRSAPublicKeyFromPEM(p []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(p)
	if block == nil || (block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY") {
		return nil, errors.New("invalid public pem")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil { return nil, err }
	k, ok := pub.(*rsa.PublicKey)
	if !ok { return nil, errors.New("not rsa public key") }
	return k, nil
}

func (j *JWTService) SignAccess(userID string, roles []string) (string, string, error) {
	now := time.Now().UTC()
	jti := uuid.NewString()
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": jti,
		"roles": roles,
		"iat": now.Unix(),
		"exp": now.Add(j.accessTTL).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	s, err := tok.SignedString(j.priv)
	return s, jti, err
}

func (j *JWTService) SignRefresh(userID string) (string, string, time.Time, error) {
	now := time.Now().UTC()
	jti := uuid.NewString()
	exp := now.Add(j.refreshTTL)
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": jti,
		"typ": "refresh",
		"iat": now.Unix(),
		"exp": exp.Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	s, err := tok.SignedString(j.priv)
	return s, jti, exp, err
}

func (j *JWTService) Parse(token string) (*jwt.Token, jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return j.pub, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))
	if err != nil { return nil, nil, err }
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid { return nil, nil, errors.New("invalid claims") }
	return parsed, claims, nil
}

// AccessTTL returns the access token TTL duration
func (j *JWTService) AccessTTL() time.Duration {
	return j.accessTTL
}

