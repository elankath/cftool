package cfcmd

import (
	"flag"
	"testing"

	"github.com/elankath/logflag"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var targetString = `api endpoint:   https://api.cf.sap.hana.ondemand.com
api version:    2.123.0
user:           i034796
org:            Cloud_Integration_teammosaic
space:          dev`

var infoJSON = `{
   "name": "",
   "build": "",
   "support": "",
   "version": 0,
   "description": "Cloud Foundry at SAP Cloud Platform",
   "authorization_endpoint": "https://login.cf.sap.hana.ondemand.com",
   "token_endpoint": "https://uaa.cf.sap.hana.ondemand.com",
   "min_cli_version": null,
   "min_recommended_cli_version": null,
   "app_ssh_endpoint": "ssh.cf.sap.hana.ondemand.com:2222",
   "app_ssh_host_key_fingerprint": "af:82:92:98:0b:80:c4:14:3e:0a:9b:c3:c8:4b:ae:21",
   "app_ssh_oauth_client": "ssh-proxy",
   "doppler_logging_endpoint": "wss://doppler.cf.sap.hana.ondemand.com:443",
   "api_version": "2.123.0",
   "osbapi_version": "2.14"
}`

func init() {
	flag.Parse()
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logflag.Parse()
}

func TestDecodeJSON(t *testing.T) {
	var info *Info
	err := decodeJSON(infoJSON, info)
	if assert.NoError(t, err, "Could not parse JSON:", infoJSON) && assert.NotNil(t, info) {
		log.Infof("Parsed Info: %v\n", info)
	} else {
		log.Errorln("TestDecodeJSON failed")
	}
}

func TestParseTargetString(t *testing.T) {
	target, err := parseTargetString(targetString)
	if assert.NoError(t, err, "Could not parse:", targetString) && assert.NotNil(t, target) {
		log.Infof("Parsed Target: %s\n", target)
	} else {
		log.Errorln("TestParseTargetString failed")
	}
}
