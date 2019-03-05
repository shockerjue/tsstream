package controller

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

var ERROR_MSG map[int]string = map[int]string {
		200:    "Success",
		
		402:	"Access denied",
}

func getErrorMsg(code int) string {
        if "" == ERROR_MSG[code] {
                return "except"
        }

        return ERROR_MSG[code]
}

func handleErr(c *gin.Context, code int) {
        c.JSON(http.StatusOK, gin.H{"code" : code, "msg" : getErrorMsg(code)})
        c.Abort()
}

func handleFunc(c *gin.Context, data interface{}) {
        c.JSON(http.StatusOK,gin.H{"code": 200, "msg" :"Success","data": data})
}