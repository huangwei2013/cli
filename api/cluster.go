package api

import (

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

func (serverCtx *ServerContext) ClusterLs(r *ghttp.Request){

	rules := map[string]string{
		"token": "required",
	}
	msgs := map[string]interface{}{
		"token": "token不能为空",
	}

	reqJson := getReqJson(r)
	if err := gvalid.CheckMap(reqJson.ToMap(), rules, msgs); err != nil {
		sendRsp(r, 0, err.String())
	}

	token := reqJson.GetString("token")
	userConfig, ok :=  serverCtx.UserConfigs[token]
	if !ok {
		sendRsp(r,500, "Token : Invalid or  Nonexistent")
	}

	collection, err := userConfig.MC.ManagementClient.Cluster.List(BaseListOpts())
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  collection.Data)
}