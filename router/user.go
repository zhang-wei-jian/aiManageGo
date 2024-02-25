package router

import (
	"github.com/gin-gonic/gin"
	models "github.com/zhang-wei-jian/aiManageGo/modles"
)

func User(r *gin.Engine) {

	loginRouter := r.Group("/user")
	{
		loginRouter.POST("login", func(c *gin.Context) {
			//username := c.PostForm("username")
			//password := c.PostForm("password")

			var body models.User
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(400, gin.H{
					"error": err,
				})
				return
			}
			c.JSON(200, gin.H{
				"username": body.Username,
				"password": body.Password,
				"aa":       "aaa",
			})
		})
	}
}
