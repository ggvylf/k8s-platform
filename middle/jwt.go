package middle

import (
	"k8s-platform/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		// login接口不处理
		if len(c.Request.URL.String()) >= 10 && c.Request.URL.String()[0:10] == "/api/login" {
			c.Next()
		} else {

			token := c.Request.Header.Get("Authorization")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":  "缺少token",
					"data": nil,
				})
				c.Abort()
				return
			}

			// 解析token内容
			claims, err := utils.JWTToken.ParseToken(token)

			// 错误处理
			if err != nil {
				if err.Error() == "TokenExpired" {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg":  "token已过期",
						"data": nil,
					})
					c.Abort()
					return
				}
				if err.Error() == "TokenMalformed" {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg":  "token格式错误",
						"data": nil,
					})
					c.Abort()
					return
				}
				if err.Error() == "TokenNotValidYet" {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg":  "token时间未生效",
						"data": nil,
					})
					c.Abort()
					return
				}
			}

			// 把token添加到context
			c.Set("claims", claims)
			c.Next()
		}

	}

}
