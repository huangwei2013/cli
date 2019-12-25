package api

import (
	"fmt"
	"time"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"github.com/rancher/cli/utils"
)

type CACertResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (serverCtx *ServerContext) Login(r *ghttp.Request){
	rules := map[string]string{
		"userName": "required",
		"password": "required",
	}
	msgs := map[string]interface{}{
		"userName": "userName 账号不能为空",
		"password": "password 密码不能为空",
	}

	reqJson := getReqJson(r)
	if err := gvalid.CheckMap(reqJson.ToMap(), rules, msgs); err != nil {
		sendRsp(r, 0, err.String())
	}

	token, err := login(serverCtx, reqJson.GetString("userName"), reqJson.GetString("password"))
	if nil != err {
		sendRsp(r,401, err.Error(), map[string]interface{}{})
	}

	// build-up context & config by token
	rancherServerConfig := &ServerConfig{}
	rancherServerConfig.URL = GetRancherUri(serverCtx.RancherHost)

	auth := utils.SplitOnColon(token)
	rancherServerConfig.AccessKey = auth[0]
	rancherServerConfig.SecretKey = auth[1]
	rancherServerConfig.TokenKey = token
	rancherServerConfig.Project = ""

	// something extra, making compatible with Rancher/cli process, to get cert
	var mc *MasterClient
	mc, err = GetMasterClient(reqJson.GetString("projectId"), rancherServerConfig)
	if nil != err {
		sendRsp(r,500, err.Error())
	}

	serverCtx.UserConfigs[token] = &UserConfig{
		UserID:		0,
		UserName:   reqJson.GetString("userName"),
		Password:   reqJson.GetString("password"),
		Token:		token,
		LoginAt:	time.Now().Unix(),
		RancherServerConfig:rancherServerConfig,
		MC:			mc,
	}

	sendRsp(r,0, "OK", map[string]interface{}{"token":token})
}

func login(serverCtx *ServerContext, username string, password string ) (string, error){
	c := ghttp.NewClient()
	token := ""

	url := fmt.Sprintf("%s/%s", GetRancherUri(serverCtx.RancherHost), "v3-public/localProviders/local?action=login" )
	data := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"description\":\"UI Session\",\"responseType\":\"cookie\",\"ttl\":57600000,\"labels\":{\"ui-session\":\"true\"}}", username, password)
	if response, err := c.Post(url, data); err != nil {
		return "", err
	} else {
		defer response.Close()
		token = response.Header.Get("Set-Cookie")[7:73]
	}
	return token, nil
}
