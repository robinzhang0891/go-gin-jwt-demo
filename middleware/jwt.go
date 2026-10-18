package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT密钥，项目中建议放到配置文件
var secretKey = os.Getenv("JWT_SECRET")

func init() {
	if len(secretKey) < 32 {
		panic("JWT_SECRET must contain at least 32 bytes")
	}
}

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userId uint) (string, error) {
	claim := Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // 有效期2小时
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secretKey))
}

// JWTAuth JWT鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "未携带token"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "token格式错误"})
			c.Abort()
			return
		}

		tokenStr := parts[1] // validation explicitly require HS256 AND expiration added
		claim := &Claims{}
		token, err := jwt.ParseWithClaims(
			tokenStr,
			claim,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "token无效或已过期"})
			c.Abort()
			return
		}
		// 将用户ID存入上下文，后续接口使用
		c.Set("userId", claim.UserID)
		c.Next()
	}
}
