package tools

import (
	"errors"
	"time"

	jwt4 "github.com/golang-jwt/jwt/v4"
)

// MyClaims Create a struct that will be encoded to a JWT.
//  We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type MyClaims struct {
	Phone string `json:"phone"`
	jwt4.RegisteredClaims
}

var (
	MySecret         = "my_secret_key_1234567809"
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token.")
)

// GetToken : generate token by uniqueId,
// 这里传入的是手机号
func GetToken(phone string) (tokenStr string, err error) {
	claim := MyClaims{
		Phone: phone,
		RegisteredClaims: jwt4.RegisteredClaims{
			ExpiresAt: jwt4.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
			IssuedAt:  jwt4.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt4.NewNumericDate(time.Now()),                                       // 生效时间
		}}
	token := jwt4.NewWithClaims(jwt4.SigningMethodHS256, claim) // 使用HS256算法
	tokenStr, err = token.SignedString([]byte(MySecret))
	return tokenStr, err
}

func Secret() jwt4.Keyfunc {
	return func(token *jwt4.Token) (interface{}, error) {
		return []byte(MySecret), nil // 这是我的secret
	}
}

// ParseToken : parse token
func ParseToken(tokens string) (*MyClaims, error) {
	token, err := jwt4.ParseWithClaims(tokens, &MyClaims{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt4.ValidationError); ok {
			if ve.Errors&jwt4.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt4.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt4.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// RefreshToken : refresh token
func RefreshToken(tokenStr string) (string, error) {
	jwt4.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt4.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt4.Token) (interface{}, error) {
		return []byte(MySecret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		jwt4.TimeFunc = time.Now
		claims.RegisteredClaims.ExpiresAt.Time = time.Now().Add(3 * time.Hour)
		return GetToken(claims.Phone)
	}
	return "", TokenInvalid
}
