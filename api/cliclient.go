package api

import (
	"fmt"
	"strings"

	errorsPkg "github.com/pkg/errors"
	"github.com/rancher/norman/clientbase"
	ntypes "github.com/rancher/norman/types"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
	"golang.org/x/sync/errgroup"
)

type MasterClient struct {
	ClusterClient    *clusterClient.Client
	ManagementClient *managementClient.Client
	ProjectClient    *projectClient.Client
	UserConfig       *ServerConfig
}

// NewMasterClient returns a new MasterClient with Cluster, Management and Project
// clients populated
func NewMasterClient(config *ServerConfig) (*MasterClient, error) {
	mc := &MasterClient{
		UserConfig: config,
	}

	var g errgroup.Group
	g.Go(mc.newManagementClient)
	g.Go(mc.newClusterClient)
	g.Go(mc.newProjectClient)

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return mc, nil
}

// NewManagementClient returns a new MasterClient with only the Management client
func NewManagementClient(config *ServerConfig) (*MasterClient, error) {
	mc := &MasterClient{
		UserConfig: config,
	}

	err := mc.newManagementClient()
	if err != nil {
		return nil, err
	}
	return mc, nil
}

// NewClusterClient returns a new MasterClient with only the Cluster client
func NewClusterClient(config *ServerConfig) (*MasterClient, error) {

	mc := &MasterClient{
		UserConfig: config,
	}

	err := mc.newClusterClient()
	if err != nil {
		return nil, err
	}
	return mc, nil
}

// NewProjectClient returns a new MasterClient with only the Project client
func NewProjectClient(config *ServerConfig) (*MasterClient, error) {

	mc := &MasterClient{
		UserConfig: config,
	}

	err := mc.newProjectClient()
	if err != nil {
		return nil, err
	}
	return mc, nil
}

func (mc *MasterClient) newManagementClient() error {
	options := createClientOpts(mc.UserConfig)

	// Setup the management client
	mClient, err := managementClient.NewClient(options)
	if err != nil {
		return err
	}
	mc.ManagementClient = mClient

	return nil
}

func (c ServerConfig) FocusedCluster() string {
	return strings.Split(c.Project, ":")[0]
}

func (mc *MasterClient) newClusterClient() error {
	options := createClientOpts(mc.UserConfig)
	options.URL = options.URL + "/clusters/" + mc.UserConfig.FocusedCluster()

	// Setup the project client
	cc, err := clusterClient.NewClient(options)
	if err != nil {
		if clientbase.IsNotFound(err) {
			err = errorsPkg.WithMessage(err, "Current cluster not available, try running `rancher context switch`. Error")
		}
		return err
	}
	mc.ClusterClient = cc

	return nil
}

func (mc *MasterClient) newProjectClient() error {
	options := createClientOpts(mc.UserConfig)
	options.URL = options.URL + "/projects/" + mc.UserConfig.Project

	// Setup the project client
	pc, err := projectClient.NewClient(options)
	if err != nil {
		if clientbase.IsNotFound(err) {
			err = errorsPkg.WithMessage(err, "Current project not available, try running `rancher context switch`. Error")
		}
		return err
	}
	mc.ProjectClient = pc

	return nil
}

func (mc *MasterClient) ByID(resource *ntypes.Resource, respObject interface{}) error {
	if _, ok := mc.ManagementClient.APIBaseClient.Types[resource.Type]; ok {
		return mc.ManagementClient.ByID(resource.Type, resource.ID, &respObject)
	} else if _, ok := mc.ProjectClient.APIBaseClient.Types[resource.Type]; ok {
		return mc.ProjectClient.ByID(resource.Type, resource.ID, &respObject)
	} else if _, ok := mc.ClusterClient.APIBaseClient.Types[resource.Type]; ok {
		return mc.ClusterClient.ByID(resource.Type, resource.ID, &respObject)
	}
	return fmt.Errorf("MasterClient - unknown resource type %v", resource.Type)
}

func createClientOpts(config *ServerConfig) *clientbase.ClientOpts {
	serverURL := config.URL

	if !strings.HasSuffix(serverURL, "/v3") {
		serverURL = config.URL + "/v3"
	}

	options := &clientbase.ClientOpts{
		URL:       serverURL,
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
		CACerts:   config.CACerts,
	}
	return options
}
