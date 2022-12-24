package jwt

import (
	"bluebell/settings"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySecret = []byte("test.com")

type MyClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成token
func GenToken(userID uint64, username string) (string, error) {
	c := &MyClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(settings.AuthSettings.JWTExpire) * time.Hour).Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),                                                                 // 签发时间
			Issuer:    "my-project",                                                                      // 签发者
			Id:        "",                                                                                // 按需求选这个, 有些实现中, 会控制这个ID是不是在黑/白名单来判断是否还有效
			NotBefore: 0,                                                                                 // 生效起始时间
			Subject:   "",                                                                                // 主题                                           // 签发人
		},
	}
	// 使用指定的签名方法创建前名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定secret签名并获取完整的编码后字符串token
	return token.SignedString(mySecret)
}

// 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}

	// 校验token
	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("invalid token")
}
