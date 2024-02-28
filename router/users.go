package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	models "github.com/zhang-wei-jian/aiManageGo/modles"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type UserClaims struct {
	Uid        int
	UserAgent  string
	Username   string
	GrantScope string
	jwt.RegisteredClaims
}

// 签名密钥
const sign_key = "95osj3fUD7fo0mlYdDbncXz4VD2igvf0"

// 随机字符串
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(str_len int) string {
	rand_bytes := make([]rune, str_len)
	for i := range rand_bytes {
		rand_bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(rand_bytes)
}

func User(r *gin.Engine) {

	loginRouter := r.Group("/users")
	loginRouter.Use(JWTMiddleware())
	{
		loginRouter.POST("login", func(c *gin.Context) {

			//获取post的JSON
			var body models.User
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(400, gin.H{
					"error": " GetJSONErr",
				})
				return
			}
			tokenString, err := generateTokenUsingHs256(c)
			if err != nil {
				c.JSON(500, "jwt系统错误")
				return
			}

			if err != nil {
				panic(err)
			}
			c.Header("Authorization", "Bearer "+tokenString)
			c.JSON(200, gin.H{
				"username": body.Username,
				"password": body.Password,
				"jwt":      tokenString,
			})
		})

		loginRouter.GET("hi", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hi",
			})
		})
	}
}

func generateTokenUsingHs256(c *gin.Context) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	claims := UserClaims{
		Uid:        854654,
		UserAgent:  c.GetHeader("User-Agent"),
		Username:   "Tom",
		GrantScope: "read_user_info",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   // 签发者
			Subject:   "人类",                                            // 签发对象
			Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
			ID:        randStr(10),                                     // wt ID, 类似于盐值
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(sign_key))
	return tokenString, err
}

func parseTokenHs256(token_string string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(token_string, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sign_key), nil //返回签名密钥
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("claim invalid")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}

	return claims, nil
}

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验
		//if ctx.Request.URL.Path == "/users/login" || ctx.Request.URL.Path == "/v1/chat/completions" {
		if ctx.Request.URL.Path == "/users/login" {
			return
		}

		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			fmt.Println("no login")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//segs := strings.SplitN(tokenHeader, " ", 2)
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 格式不对，有人瞎搞
			fmt.Println("header - jwt - nil")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		//claims := &UserClaims{}
		claims := &UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			// token 不对，有人搞你
			fmt.Println("header - jwt - error")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 按照道理来说，是不可能走到这一步的
			fmt.Println("header - jwt - aaaa")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		if claims.ExpiresAt.Time.Before(now) {
			// 过期了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			fmt.Println("header - jwt - date")
			return
		}
		if claims.UserAgent != ctx.GetHeader("User-Agent") {
			// user agent 不相等
			ctx.AbortWithStatus(http.StatusUnauthorized)
			fmt.Println("header - jwt - aaaaac11a")
			fmt.Println(claims.UserAgent, "and", ctx.GetHeader("User-Agent"))

			return
		}
		// 为了演示，假设十秒钟刷新一次
		if claims.ExpiresAt.Time.Sub(now) < time.Second*50 {
			// 刷新
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
			if err != nil {
				// 因为刷新这个事情，并不是一定要做的，所以这里可以考虑打印日志
				// 暂时这样打印
				log.Println(err)
				return
			}
			ctx.Header("x-jwt-token", tokenStr)
			log.Println("刷新了 token")
		}
	}

}
