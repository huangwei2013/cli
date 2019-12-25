package api

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	utils "github.com/rancher/cli/utils"
	projectClient "github.com/rancher/types/client/project/v3"
)

func (serverCtx *ServerContext) PipelineExecutionLs(r *ghttp.Request){

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
	if reqJson.GetString("pipelineId") != "" {
		listOpts.Filters["pipelineId"] = reqJson.GetString("pipelineId")
	}

	collection, err := mc.ProjectClient.PipelineExecution.List(listOpts)
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  collection.Data)
}


func (serverCtx *ServerContext) PipelineExecutionGetByID(r *ghttp.Request){
	rules := map[string]string{
		"token": "required",
		"projectId": "required",
		"pipelineExecutionId": "required",
	}
	msgs := map[string]interface{}{
		"token": "token 不能为空",
		"pipelineExecutionId": "pipelineExecutionId 不能为空",
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

	pipelineExecution, err := mc.ProjectClient.PipelineExecution.ByID(reqJson.GetString("pipelineExecutionId"))
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  pipelineExecution)
}


func (serverCtx *ServerContext) PipelineExecutionCreate(r *ghttp.Request){
	rules := map[string]string{
		"token": "required",
		"projectId": "required",
		"pipelineId": "required",
		"pipelineExecutionName": "required",
	}
	msgs := map[string]interface{}{
		"token": "token 不能为空",
		"projectId": "projectId 不能为空",
		"pipelineId": "pipelineId 不能为空",
		"pipelineExecutionName": "pipelineExecutionName 不能为空",
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

	pipelineConfig, _ := GetPipelineConfig()
	toCreate := &projectClient.PipelineExecution{
		TriggeredBy: "user",
		PipelineConfig: pipelineConfig,
		PipelineID:    reqJson.GetString("pipelineId"),
	}
	toCreate.Name = fmt.Sprintf("%s-%d", pipeline.Name, pipeline.NextRun)
	toCreate.Labels = map[string]string{utils.PipelineFinishLabel: ""}
	toCreate.Run = pipeline.NextRun

	toCreate.State = utils.StateWaiting
	toCreate.Started = time.Now().Format(time.RFC3339)
	toCreate.Conditions = nil

	for i := 0; i < len(toCreate.Stages); i++ {
		stage := &toCreate.Stages[i]
		stage.State = utils.StateWaiting
		stage.Started = ""
		stage.Ended = ""
		for j := 0; j < len(stage.Steps); j++ {
			step := &stage.Steps[j]
			step.State = utils.StateWaiting
			step.Started = ""
			step.Ended = ""
		}
	}

	newPipelineExecution, err := mc.ProjectClient.PipelineExecution.Create(toCreate)
	if nil != err {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK", newPipelineExecution)
}


func GetPipelineConfig()(*projectClient.PipelineConfig, error){

	path := "./source/.rancher-pipeline.yml"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pipelineConfig, err := utils.PipelineConfigFromYaml(content)
	if err != nil {
		return nil, err
	}

	return pipelineConfig, nil
}