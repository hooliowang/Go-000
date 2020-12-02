package api

import (
	"errortest/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetUser(c *gin.Context) {
	code := 0
	msg := "OK"
	data := make(map[string]interface{})
	for {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			code = 500
			msg = "id param is invalid"
			break
		}

		user, err := service.GetUser(id)
		if err != nil {
			fmt.Println("err: ", err)
			fmt.Printf("original error: %T %v\n", errors.Cause(err), errors.Cause(err))
			fmt.Printf("stack trace:\n%+v\n", err)

			code = 500
			msg = err.Error()

			break
		}

		data["user"] = user

		break
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
