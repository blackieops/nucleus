package auth

import (
	"strconv"

	"com.blackieops.nucleus/data"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	DBContext *data.Context
}

func (a *AuthRouter) Mount(r *gin.RouterGroup) {
	r.GET("/login", func(c *gin.Context) {
		users := FindAllUsers(a.DBContext)
		c.HTML(200, "auth_login.html", gin.H{"users": users})
	})

	r.POST("/login", func(c *gin.Context) {
		userId, err := strconv.Atoi(c.PostForm("UserID"))
		if err != nil {
			c.AbortWithStatus(422)
			return
		}

		user := FindUser(a.DBContext, userId)
		session := sessions.Default(c)
		session.Set("CurrentUserID", user.ID)
		session.Save()

		if rt := session.Get("ReturnTo"); rt != nil {
			c.Redirect(302, rt.(string))
			c.Abort()
			return
		}

		c.Redirect(302, "/")
	})
}
