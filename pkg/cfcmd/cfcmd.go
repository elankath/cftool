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
	Description           string `json:"description"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	AppSSHEndpoint        string `json:"app_ssh_endpoint"`
	APIVersion            string `json:"api_version"`
}

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

// NewClient returns an intiaized CF Client that self-intiaizes by leveraging
// the cf-cli to discover the API endoint and session bearer
// token and . Requires that some user has already logged via `cf login`
// func NewClient() (c *Client, err error) {
// 	c.Target, err = getTarget()
// 	if err != nil {
// 		return nil, err
// 	}
// 	c.Info, err = GetInfo()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return c, nil
// }

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
		return
	}
	err = decodeJSON(cmdOut, info)
	return
}

func GetCFAppGUID(appName string) (string, error) {
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

func decodeJSON(strJSON string, v interface{}) error {
	return json.NewDecoder(strings.NewReader(strJSON)).Decode(v)
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
