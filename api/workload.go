package api

import (

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

func (serverCtx *ServerContext) WorkloadLs(r *ghttp.Request){

	rules := map[string]string{
		"token": "required",
	}
	msgs := map[string]interface{}{
		"token": "token 不能为空",
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

	mc, err := GetMasterClient(reqJson.GetString("projectId"), userConfig.RancherServerConfig)
	if nil != err {
		sendRsp(r,500, err.Error())
	}

	collection, err := mc.ProjectClient.Workload.List(BaseListOpts())
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  collection.Data)
}