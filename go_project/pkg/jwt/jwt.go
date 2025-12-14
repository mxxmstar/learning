package jwt_manager

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig JWT相关配置
type JWTConfig struct {
	SecretKey []byte `mapstructure:"jwt_secret_key"` // 用于签名的密钥
	Issuer    string `mapstructure:"jwt_issuer"`     // 签发者（可选）
	Expire    int    `mapstructure:"jwt_expire"`
}

type CustomClaims struct {
	UserID   uint64 `json:"user_id"`
	DeviceID string `json:"device_id,omitempty"`
	jwt.RegisteredClaims
}

type JWT struct {
	config *JWTConfig
}

func NewJWT(secretKey []byte, issuer string, expire int) *JWT {
	return &JWT{
		config: &JWTConfig{
			SecretKey: secretKey,
			Issuer:    issuer,
			Expire:    expire,
		},
	}
}

func (j *JWT) GenerateToken(userID uint64, deviceID string) (string, error) {
	if userID == 0 {
		return "", errors.New("user id is null")
	}

	claims := CustomClaims{
		UserID:   userID,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.config.Expire) * time.Second)),
			Subject:   strconv.FormatUint(userID, 10),
		},
	}

	// 创建并签名JWT 这里是HS256算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.config.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is null")
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.config.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
