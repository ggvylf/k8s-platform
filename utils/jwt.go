package utils

import (
	"errors"
	"k8s-platform/config"

	"github.com/golang-jwt/jwt/v5"
)

type jwtToken struct{}

var JWTToken jwtToken

// 定义token内容
type CustomeCliams struct {
	jwt.RegisteredClaims
	UserName string `json:"username"`
	Password string `json:"password"`
}

// 解析token
func (*jwtToken) ParseToken(tokenstring string) (claims *CustomeCliams, err error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenstring, &CustomeCliams{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret), nil
	})

	if err != nil {
		return nil, errors.New("token解析失败")
	}

	// 判断验证信息,v5版本里已经没有*jwt.ValidationError类型了
	// if err != nil {
	// 	logger.Error("parse token failed ", err)
	// 	if ve, ok := err.(*jwt.ValidationError); ok {
	// 		// 位运算
	// 		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
	// 			return nil, errors.New("TokenMalformed")
	// 		} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
	// 			// Token is expired
	// 			return nil, errors.New("TokenExpired")
	// 		} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
	// 			return nil, errors.New("TokenNotValidYet")
	// 		} else {
	// 			return nil, errors.New("TokenInvalid")
	// 		}
	// 	}
	// }

	// 转换类型
	if claims, ok := token.Claims.(*CustomeCliams); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token解析失败")

}
