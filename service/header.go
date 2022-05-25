/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package service

import (
	"dragonfly/controllers/api"
	"github.com/gin-gonic/gin"
)

func Header() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		for name, values := range c.Request.Header {
			// Loop over all values for the name.
			for _, value := range values {
				if name == "Authorization" {
					token = value
				}
			}
		}
		api.Token = token

		c.Next()
	}
}
