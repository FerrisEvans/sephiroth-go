package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"sephiroth-go/core"
	"sephiroth-go/model/req"
	"time"
)

type Jwt struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT() *Jwt {
	return &Jwt{
		[]byte(core.Config.Jwt.SigningKey),
	}
}

func (j *Jwt) CreateClaims(baseClaims req.BaseClaims) req.CustomClaims {
	bf, _ := ParseDuration(core.Config.Jwt.BufferTime)
	ep, _ := ParseDuration(core.Config.Jwt.ExpiresTime)
	claims := req.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},                   // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    core.Config.Jwt.Issuer,                    // 签名的发行者
		},
	}
	return claims
}

// CreateToken 创建一个token
func (j *Jwt) CreateToken(claims req.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *Jwt) CreateTokenByOldToken(oldToken string, claims req.CustomClaims) (string, error) {
	v, err, _ := core.ConcurrencyControl.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// ParseToken 解析 token
func (j *Jwt) ParseToken(tokenString string) (*req.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &req.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		//if ve, ok := err.(*jwt.ValidationError); ok {
		//	if ve.Errors&jwt.ValidationErrorMalformed != 0 {
		//		return nil, TokenMalformed
		//	} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
		//		// Token is expired
		//		return nil, TokenExpired
		//	} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
		//		return nil, TokenNotValidYet
		//	} else {
		//		return nil, TokenInvalid
		//	}
		//}
	}
	if token != nil {
		if claims, ok := token.Claims.(*req.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
