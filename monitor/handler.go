package monitor

import (
	"strings"
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SendPostRequestHttp(url string, param string) (response string,err error) {
	postReq, err := http.NewRequest("POST",url,strings.NewReader(param))
    if err != nil {
        return 
    }

    postReq.Header.Set("Content-Type", "application/json;encoding=utf-8")

    client := &http.Client{}
    resp, err := client.Do(postReq)
    if err != nil {
        return 
	} 
	
	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil {
		return "",err
	}

	defer resp.Body.Close()

	response = string(body)

	return 
}

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