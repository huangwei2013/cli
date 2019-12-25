package api

import (

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

func (serverCtx *ServerContext) PodLs (r *ghttp.Request){

	rules := map[string]string{
		"token": "required",
		"projectId": "required",
	}
	msgs := map[string]interface{}{
		"token": "token 不能为空",
		"projectId": "projectId 不能为空",
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

	listOpts := BaseListOpts()
	if reqJson.GetString("nodeId") != "" {
		listOpts.Filters["NodeID"] = reqJson.GetString("nodeId")
	}
	if reqJson.GetString("namespaceId") != "" {
		listOpts.Filters["NamespaceId"] = reqJson.GetString("namespaceId")
	}

	collection, err := mc.ProjectClient.Pod.List(listOpts)
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  collection.Data)
}