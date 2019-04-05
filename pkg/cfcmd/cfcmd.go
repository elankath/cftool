package cfcmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"net/http"

	"github.com/elankath/cftool/pkg/cmd"
	"github.com/pkg/errors"
)

// Client is a CF Client and must be obtained via NewClient
type Client struct {
	Target *Target
	Info   *Info
	http   *http.Client
}

// // Config represents the configuration info
// type Config struct {
// 	Token       string
// 	APIEndpoint string
// }

// Info represents the CF Info at /v2/info
type Info struct {
	Name                        string `json:"name"`
	Build                       string `json:"build"`
	Support                     string `json:"support"`
	Version                     int    `json:"version"`
	Description                 string `json:"description"`
	AuthorizationEndpoint       string `json:"authorization_endpoint"`
	TokenEndpoint               string `json:"token_endpoint"`
	MinCLIVersion               string `json:"min_cli_version"`
	MinCLIRecommendedCLIVersion string `json:"min_recommended_cli_version"`
	AppSSHEndpoint              string `json:"app_ssh_endpoint"`
	AppSSHHostKeyFingerprint    string `json:"app_ssh_host_key_fingerprint"`
	AppSSHOAuthClient           string `json:"app_ssh_oauth_client"`
	DopplerLoggingEndpoint      string `json:"doppler_logging_endpoint"`
	APIVersion                  string `json:"api_version"`
	OSBAPIVersion               string `json:"osbapi_version"`
	BITSEndpoint                string `json:"bits_endpoint"`
}
// "name": "",
// "build": "",
// "support": "",
// "version": 0,
// "description": "Cloud Foundry at SAP Cloud Platform",
// "authorization_endpoint": "https://login.cf.sap.hana.ondemand.com",
// "token_endpoint": "https://uaa.cf.sap.hana.ondemand.com",
// "min_cli_version": null,
// "min_recommended_cli_version": null,
// "app_ssh_endpoint": "ssh.cf.sap.hana.ondemand.com:2222",
// "app_ssh_host_key_fingerprint": "af:82:92:98:0b:80:c4:14:3e:0a:9b:c3:c8:4b:ae:21",
// "app_ssh_oauth_client": "ssh-proxy",
// "doppler_logging_endpoint": "wss://doppler.cf.sap.hana.ondemand.com:443",
// "api_version": "2.128.0",
// "osbapi_version": "2.14",
// "bits_endpoint": "https://bits.cf.sap.hana.ondemand.com"

// Target represents output of CF target
type Target struct {
	APIEndpoint string
	APIVersion  string
	User        string
	Org         string
	Space       string
}

func (t Target) String() string {
	return fmt.Sprintf("{APIEndpoint:%s, APIVersion:%s, User:%s, Org: %s, Space: %s}", t.APIEndpoint, t.APIVersion, t.User, t.Org, t.Space)
}

func GetTarget() (*Target, error) {
	cmdArgs := []string{"target"}
	cmdOut, err := cmd.Exec("cf", cmdArgs, false)
	if err != nil {
		return nil, err
	}
	if strings.Contains(cmdOut, "Not logged") {
		return nil, errors.New("Must be logged into cf")
	}
	return parseTargetString(cmdOut)
}

func GetInfo() (info *Info, err error) {
	cmdArgs := []string{"curl", "/v2/info"}
	cmdOut, err := cmd.Exec("cf", cmdArgs, false)
	if err != nil {
		return nil, err
	}
	return unmarshalInfo(cmdOut)
}

func GetAppGUID(appName string) (string, error) {
	cmdArgs := []string{"app", appName, "--guid"}
	guid, err := cmd.Exec("cf", cmdArgs, true)
	if err != nil {
		return "", errors.Wrapf(err, "Cannot obtain guid of app '%s'", appName)
	}
	return strings.TrimSpace(guid), nil
}

func GenSSHCode() (code string, err error) {
	cmdArgs := []string{"ssh-code"}
	out, err := cmd.Exec("cf", cmdArgs, true)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func getOAuthTokenViaCli() (string, error) {
	cmdArgs := []string{"app", "oauth-token"}
	token, err := cmd.Exec("cf", cmdArgs, true)
	if err != nil {
		return "", errors.Wrapf(err, "Cannot obtain OAuth token via cf cli")
	}
	return strings.TrimSpace(token), nil
}

func unmarshalInfo(json string) (*Info, error) {
	var info Info
	err := decodeJSON(json, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func decodeJSON(strJSON string, v interface{}) error {
	return json.Unmarshal([]byte(strJSON), v)
}

func parseTargetString(targetString string) (*Target, error) {
	matches := targetParseRe.FindStringSubmatch(targetString)
	if matches == nil {
		return nil, errors.New("Can't parse target string: " + targetString)
	}

	if len(matches) < 5 {
		return nil, errors.New("In-sufficient matches parse target string: " + targetString)
	}
	target := &Target{
		APIEndpoint: matches[1],
		APIVersion:  matches[2],
		User:        matches[3],
		Org:         matches[4],
		Space:       matches[5],
	}
	return target, nil
}

var targetParseRe = regexp.MustCompile(`api endpoint:\s*(.*?)\napi version:\s*(.*?)\nuser:\s*(.*?)\norg:\s*(.*?)\nspace:\s*(.*?)`)
