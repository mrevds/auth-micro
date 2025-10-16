package utils

import (
    "auth-micro/internal/auth/config"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type Claims struct {
    UserID string `json:"user_id"`
    Type   string `json:"type"`
    jwt.RegisteredClaims
}

// JWTManager управляет JWT токенами
type JWTManager struct {
    cfg *config.Config
}

func NewJWTManager(cfg *config.Config) *JWTManager {
    return &JWTManager{cfg: cfg}
}

func (j *JWTManager) GenerateToken(userID string) (string, error) {
    claims := Claims{
        UserID: userID,
        Type:   "access",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JWT.AccessTokenDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.cfg.JWT.SecretKey))
}

func (j *JWTManager) GenerateRefreshToken(userID string) (string, error) {
    claims := Claims{
        UserID: userID,
        Type:   "refresh",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JWT.RefreshTokenDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.cfg.JWT.SecretKey))
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(j.cfg.JWT.SecretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrTokenInvalidClaims
}

// Утилиты для паролей (не зависят от конфига)
func CheckPasswordHash(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}