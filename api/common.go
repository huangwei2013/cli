package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/rancher/norman/types"
)

type ServerContext struct {
	StartAt int64
	HostIP string
	Port int
	RancherHost string
	UserConfigs map[string]*UserConfig // key:value = token:UserConfig
	Server *ghttp.Server
}

// user config
type UserConfig struct{
	UserID int
	UserName string
	Password string
	Token string
	LoginAt int64
	RancherServerConfig *ServerConfig
	MC *MasterClient
}


//ServerConfig holds the config for each server the user has setup
type ServerConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	TokenKey  string `json:"tokenKey"`
	URL       string `json:"url"`
	Project   string `json:"project"`
	CACerts   string `json:"cacert"`
}

func getReqJson(r *ghttp.Request ) *gjson.Json {
	reqJson, err := gjson.DecodeToJson(r.GetRaw())
	if err != nil {
		sendRsp(r, 400, err.Error())
	}
	return reqJson
}

// 标准返回结果数据结构封装。
// 返回固定数据结构的JSON:
// code:  错误码(0:成功, 1:失败, >1:错误码);
// msg:  请求结果信息;
// data: 请求结果,根据不同接口返回结果的数据结构不同;
func sendRsp(r *ghttp.Request, code int, msg string, data ...interface{}){
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(g.Map{
		"code":  code,
		"msg":  msg,
		"data": responseData,
	})
	r.Exit()
}

func BaseListOpts() *types.ListOpts {
	return &types.ListOpts{
		Filters: map[string]interface{}{
			"limit": -2,
			"all":   true,
		},
	}
}

func GetRancherUri(rancherHost string) string{
	return fmt.Sprintf("https://%s", rancherHost)
}


func ParseClusterIDFromProjectID(projectID string) string {
	return strings.Split(projectID, ":")[0]
}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func GetMasterClient(projectId string, config *ServerConfig) (*MasterClient, error){
	//make a new one
	config.Project = projectId
	mc, err := NewMasterClient(config)
	if nil != err {
		if _, ok := err.(*url.Error); ok && strings.Contains(err.Error(), "certificate signed by unknown authority") {
			mc, err = getCertFromServer(config)
			if nil != err {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return mc, nil
}


func getCertFromServer(cf *ServerConfig) (*MasterClient, error) {
	req, err := http.NewRequest("GET", cf.URL+"/v3/settings/cacerts", nil)
	if nil != err {
		return nil, err
	}

	req.SetBasicAuth(cf.AccessKey, cf.SecretKey)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if nil != err {
		return nil, err
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return nil, err
	}

	var certReponse *CACertResponse
	err = json.Unmarshal(content, &certReponse)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse response from %s/v3/settings/cacerts\nError: %s\nResponse:\n%s", cf.URL, err, content)
	}

	cert, err := verifyCert([]byte(certReponse.Value))
	if nil != err {
		return nil, err
	}

	cf.CACerts = cert
	return NewMasterClient(cf)
}

func verifyCert(caCert []byte) (string, error) {
	// replace the escaped version of the line break
	caCert = bytes.Replace(caCert, []byte(`\n`), []byte("\n"), -1)
	block, _ := pem.Decode(caCert)

	if nil == block {
		return "", errors.New("No cert was found")
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if nil != err {
		return "", err
	}

	if !parsedCert.IsCA {
		return "", errors.New("CACerts is not valid")
	}
	return string(caCert), nil
}
