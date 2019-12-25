package api

import (
	"fmt"
	"time"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	projectClient "github.com/rancher/types/client/project/v3"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (serverCtx *ServerContext) DeploymentLs(r *ghttp.Request){

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

	userConfig, ok :=  serverCtx.UserConfigs[reqJson.GetString("token")]
	if !ok {
		sendRsp(r,500, "Token : Invalid or  Nonexistent")
	}

	mc, err := GetMasterClient(reqJson.GetString("projectId"), userConfig.RancherServerConfig)
	if nil != err {
		sendRsp(r,500, err.Error())
	}

	collection, err := mc.ProjectClient.Deployment.List(BaseListOpts())
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  collection.Data)
}

func (serverCtx *ServerContext) DeploymentGetByID(r *ghttp.Request){
	rules := map[string]string{
		"token": "required",
		"projectId": "required",
		"deploymentId": "required",
	}
	msgs := map[string]interface{}{
		"token": "token 不能为空",
		"projectId": "projectId 不能为空",
		"deploymentId": "deploymentId 不能为空",
	}

	reqJson := getReqJson(r)
	if err := gvalid.CheckMap(reqJson.ToMap(), rules, msgs); err != nil {
		sendRsp(r, 0, err.String())
	}

	userConfig, ok :=  serverCtx.UserConfigs[reqJson.GetString("token")]
	if !ok {
		sendRsp(r,500, "Token : Invalid or  Nonexistent")
	}

	mc, err := GetMasterClient(reqJson.GetString("projectId"), userConfig.RancherServerConfig)
	if nil != err {
		sendRsp(r,500, err.Error())
	}

	deployment, err := mc.ProjectClient.Deployment.ByID(reqJson.GetString("deploymentId"))
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  deployment)
}

func (serverCtx *ServerContext) DeploymentCreate(r *ghttp.Request){

	rules := map[string]string{
		"token": "required",
		"projectId": "required",
		"image": "required",
		"deploymentName": "required",
		"namespaceId": "required",
		"containerName": "required",
	}
	msgs := map[string]interface{}{
		"token": "token 不能为空",
		"projectId": "projectId 不能为空",
		"image": "image 不能为空",
		"deploymentName": "deploymentName 不能为空",
		"namespaceId": "namespaceId 不能为空",
		"containerName": "containerName 不能为空",
	}

	reqJson := getReqJson(r)
	if err := gvalid.CheckMap(reqJson.ToMap(), rules, msgs); err != nil {
		sendRsp(r, 0, err.String())
	}

	userConfig, ok :=  serverCtx.UserConfigs[reqJson.GetString("token")]
	if !ok {
		sendRsp(r,500, "Token : Invalid or  Nonexistent")
	}

	deploymentConfig := &projectClient.DeploymentConfig{
		MinReadySeconds: 0,
		Strategy: "RollingUpdate",
		MaxSurge: intstr.FromInt(1),
		MaxUnavailable: intstr.FromInt(0),
	}

	containers := []projectClient.Container{
		projectClient.Container{
			Image: reqJson.GetString("image"),
			Name: fmt.Sprintf("%v-%v",reqJson.GetString("containerName"), time.Now().UnixNano()),
		},
	}

	newDeployment := &projectClient.Deployment{
		Name:        reqJson.GetString("deploymentName"),
		RestartPolicy: "Always",
		NamespaceId: reqJson.GetString("namespaceId"),
		DeploymentConfig: deploymentConfig,
		Containers:containers,
		ProjectID:  reqJson.GetString("projectId"),
	}

	mc, err := GetMasterClient(reqJson.GetString("projectId"), userConfig.RancherServerConfig)
	if nil != err {
		sendRsp(r,500, err.Error())
	}
	deployment, err := mc.ProjectClient.Deployment.Create(newDeployment)
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK",  deployment)
}


func (serverCtx *ServerContext) DeploymentDelete(r *ghttp.Request){

	rules := map[string]string{
		"token": "required",
		"projectId": "required",
		"deploymentId": "required",
	}
	msgs := map[string]interface{}{
		"token": "token 不能为空",
		"projectId": "projectId 不能为空",
		"deploymentId": "deploymentId 不能为空",
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

	deployment, err := mc.ProjectClient.Deployment.ByID(reqJson.GetString("deploymentId"))
	if err != nil {
		sendRsp(r,500, err.Error())
	}

	err = mc.ProjectClient.Deployment.Delete(deployment)
	if nil != err {
		sendRsp(r,500, err.Error())
	}

	sendRsp(r,0, "OK")
}