package monitor

import (
	"io"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
)

func EncodePasswd(str string) string {
	t := md5.New()
	io.WriteString(t, str)
	return hex.EncodeToString(t.Sum(nil))
}

type Node struct {
	Name 		string 	`json:"name"`		
	Connects	int 	`json:"connects"`
	Bind 		string 	`json:"bind"`
	Port 		string 	`json:"port"`
	Hash 		string 	`json:"hash"`
}

type MonitorInfo struct {
	NodeInfo 	Node 		`json:"nodeinfo,omitempty"`
	Packages 	int 		`json:"packages"`
	NextNode	[]Node 		`json:"nextnode,omitempty"`
	Genesis		bool 		`json:"genesis"`	
}

func (this *MonitorInfo)Encode() (string,error) {
	this.Hash()
	
	msgstr, err := json.Marshal(this)
	if nil != err {
		return "",err
	}

	return string(msgstr),nil
}

func (this *MonitorInfo)Decode(data string) (err error) {
	err = json.Unmarshal([]byte(data),this)
	if nil == err {
		return 
	}

	this.Hash()

	return 
}

func (this *MonitorInfo)Hash()  {
	for k,v := range this.NextNode {
		this.NextNode[k].Hash = EncodePasswd(v.Bind + v.Port)
	}
	
	this.NodeInfo.Hash = EncodePasswd(this.NodeInfo.Bind + this.NodeInfo.Port)

	return 
}

var NodeInfos map[string]MonitorInfo = make(map[string]MonitorInfo)