package middleware

import (
	"net/http"

	helper "github.com/AnthuanGarcia/appWebSeguridad/src/helpers"

	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authentication(c *gin.Context) {

	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
		c.Abort()
		return
	}

	claims, err := helper.ValidateToken(clientToken)
	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	c.Set("username", claims.Username)
	c.Set("nombre", claims.Nombre)
	c.Set("ape_paterno", claims.Ape_Paterno)
	c.Set("ID", claims.ID)

	c.Next()

}
