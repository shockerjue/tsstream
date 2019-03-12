package monitor

import (
	"net/http"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)


func PostMonitorFunc(c *gin.Context) {
	var monitor MonitorInfo
	err := c.ShouldBindJSON(&monitor)
	if nil != err {
		log.Error("monitor.Decode error:",err)

		return 
	}

	monitor.Hash()
	if monitor.Genesis {
		NodeInfos["genesis"] = monitor
	}else {
		NodeInfos[monitor.NodeInfo.Hash] = monitor
	}

	return 
}

func GetMonitorFunc(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"code": 200, "msg" :"Success","data": NodeInfos})

	return 
}