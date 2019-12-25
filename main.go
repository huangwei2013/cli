package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/rancher/cli/api"
)

var serverCtx = api.ServerContext{
	StartAt: time.Now().Unix(),
}

func init() {
	g.Cfg().SetFileName("config/config.yaml")

	// init Server Context
	ip, err := api.ExternalIP()
	if err != nil {
		fmt.Println(err)
	}
	serverCtx.HostIP = string(ip)
	serverCtx.Port = g.Cfg().GetInt("base.port")
	serverCtx.Server = g.Server()
	serverCtx.RancherHost = g.Cfg().GetString("base.rancherHost")
	serverCtx.UserConfigs = make(map[string]*api.UserConfig)

	// init routers
	serverCtx.Server.BindHandler("/login", serverCtx.Login)
	serverCtx.Server.BindHandler("/cluster/list", serverCtx.ClusterLs)
	serverCtx.Server.BindHandler("/project/list", serverCtx.ProjectLs)
	serverCtx.Server.BindHandler("/project/get", serverCtx.ProjectGetByID)
	serverCtx.Server.BindHandler("/service/list", serverCtx.ServiceLs)
	serverCtx.Server.BindHandler("/pipeline/list", serverCtx.PipelineLs)
	serverCtx.Server.BindHandler("/pipeline/get", serverCtx.PipelineGetByID)
	serverCtx.Server.BindHandler("/pipeline/create", serverCtx.PipelineCreate)
	serverCtx.Server.BindHandler("/pipelineexecution/list", serverCtx.PipelineExecutionLs)
	serverCtx.Server.BindHandler("/pipelineexecution/get", serverCtx.PipelineExecutionGetByID)
	serverCtx.Server.BindHandler("/pipelineexecution/create", serverCtx.PipelineExecutionCreate)
	serverCtx.Server.BindHandler("/workload/list", serverCtx.WorkloadLs)
	serverCtx.Server.BindHandler("/deployment/list", serverCtx.DeploymentLs)
	serverCtx.Server.BindHandler("/deployment/get", serverCtx.DeploymentGetByID)
	serverCtx.Server.BindHandler("/deployment/create", serverCtx.DeploymentCreate)
	serverCtx.Server.BindHandler("/sourceCodeRepository/list", serverCtx.SourceCodeRepositoryLs)
	serverCtx.Server.BindHandler("/sourceCodeCredential/list", serverCtx.SourceCodeCredentialLs)
	serverCtx.Server.BindHandler("/sourceCodeCredential/get", serverCtx.SourceCodeCredentialGetByID)
	serverCtx.Server.BindHandler("/sourceCodeProvider/list", serverCtx.SourceCodeProviderLs)
	serverCtx.Server.BindHandler("/sourceCodeProviderConfig/list", serverCtx.SourceCodeProviderConfigLs)

}


func main(){
	fmt.Println("Server Staring...")
	serverCtx.Server.SetPort(serverCtx.Port)
	serverCtx.Server.Run()
}