package api

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	projectClient "github.com/rancher/types/client/project/v3"
)

func (serverCtx *ServerContext) PipelineLs(r *ghttp.Request){
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
	userConfig, ok := serverCtx.UserConfigs[token]
	if !ok {
		sendRsp(r,500, "Token : Invalid or  Nonexistent")
	}

	mc, err := GetMasterClient(reqJson.GetString("projectId"), userConfig.RancherServerConfig)
	if nil != err {
		sendRsp(r,500, err.Error())
	}

	listOpts := BaseListOpts()
	listOpts.Filters["system"] = false
	collection, err := mc.ProjectClient.Pipeline.List(listOpts)
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  collection.Data)
}


func (serverCtx *ServerContext) PipelineGetByID(r *ghttp.Request){
	rules := map[string]string{
		"token": "required",
		"projectId": "required",
		"pipelineId": "required",
	}
	msgs := map[string]interface{}{
		"token": "token 不能为空",
		"projectId": "projectId 不能为空",
		"pipelineId": "pipelineId 不能为空",
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

	pipeline, err := mc.ProjectClient.Pipeline.ByID(reqJson.GetString("pipelineId"))
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  pipeline)
}


func (serverCtx *ServerContext) PipelineCreate(r *ghttp.Request){
	rules := map[string]string{
		"token": "required",
		"projectId": "required",
		"pipelineName": "required",
		"repositoryURL": "required",
		"sourceCodeCredentialId": "required",
	}
	msgs := map[string]interface{}{
		"token": "token 不能为空",
		"projectId": "projectId 不能为空",
		"pipelineName": "pipelineName 不能为空",
		"repositoryURL": "repositoryURL 不能为空",
		"sourceCodeCredentialId": "sourceCodeCredentialId 不能为空",
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

	sourceCodeCredential, err := mc.ProjectClient.SourceCodeCredential.ByID(reqJson.GetString("sourceCodeCredentialId"))
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	newPipelineObj := &projectClient.Pipeline{
		ProjectID:      reqJson.GetString("projectId"),
		Name:        	reqJson.GetString("pipelineName"),
		RepositoryURL:	reqJson.GetString("repositoryURL"),
		SourceCodeCredentialID:	reqJson.GetString("sourceCodeCredentialId"),
		SourceCodeCredential:	sourceCodeCredential,
		TriggerWebhookPr		:false,
		TriggerWebhookPush		:true,
		TriggerWebhookTag		:false,
	}

	newPipeline, err := mc.ProjectClient.Pipeline.Create(newPipelineObj)
	if nil != err {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK", newPipeline)
}
