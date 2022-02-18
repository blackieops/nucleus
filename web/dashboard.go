package web

import (
	"github.com/gin-gonic/gin"
)

func (r *WebRouter) handleDashboardShow(c *gin.Context) {
	// Since we do not have a dashboard implemeneted yet, this just redirects
	// to the user profile settings page, since that's basically all you can do
	// at this point.
	c.Redirect(302, "/web/me")
}
