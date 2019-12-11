package cfcp

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/elankath/cfutil/pkg/cfcmd"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// CopierConfig represents the configuration struct for the Copier and is needed for obtaining a Copier via NewCopier
type CopierConfig struct {
	AppName   string
	Simulated bool
}

// Copier copies files from the local FS to the CF container FS of an application instance.
type Copier struct {
	config       *CopierConfig
	appGUID      string
	numInstances int
	// target *cfcmd.Target
	info *cfcmd.Info
}

// PasswordGenerator is a function type that when invoked returns an authentication code (for SSH) or an error.
type PasswordGenerator func() (code string, err error)

// NewCopier initiazes a Copier so that its Copy method can be invoked.
func NewCopier(config CopierConfig) (*Copier, error) {
	c := &Copier{
		config: &config,
	}
	info, err := cfcmd.GetInfo()
	if err != nil {
		return nil, err
	}
	c.info = info
	fmt.Printf("**** SSH Endpoint: %s\n", c.info.AppSSHEndpoint)
	return c, nil
}

// Copy the files/dirs matching the source paths, minutes the excludes to the target dir
// on the container instance
func (c *Copier) Copy(source string, target string, excludes []string) error {
	if !strings.HasPrefix(target, "/") {
		fmt.Println("Target path for secure copy must be absolute. Prefixing /")
		target = "/" + target
	}
	fi, err := os.Stat(source)
	switch {
	case err != nil:
		return err
	case fi.IsDir():
		log.Fatalf("TODO: Not yet implemented dir copy") // TODO: Implement me
	default:
		err = c.copyFile(source, target, excludes)
	}
	if err != nil {
		return err
	}
	return nil
}

// cmdArgs := []string{"curl", "/v2/info"}
// v2Info, err := execCommand("cf", cmdArgs, false)
// matches := sshEndpointRe.FindStringSubmatch(v2Info)
// if matches == nil {
// 	return "", -1, errors.New("Can't find ssh endpoint in /v2/info")
// }
// sshEndpoint := matches[1]
// if sshEndpoint == "" {
// 	return "", -1, errors.New("Empty ssh endpoint in /v2/info")
// }
// splits := strings.Split(sshEndpoint, ":")
// var host string
// var port int
// host = splits[0]
// if len(splits) > 1 {
// 	port, err = strconv.Atoi(splits[1])
// } else {
// 	port = 22
// }
// if err != nil {
// 	return "", -1, err
// }
// return host, port, nil
// }

// func getCFAppNumInstances(appName string) (int, error) {
// 	cmdArgs := []string{"app", appName}
// 	appSummary, err := execCommand("cf", cmdArgs, false)
// 	if err != nil {
// 		return -1, err
// 	}
// 	count := 0
// 	for _, ln := range strings.Split(appSummary, "\n") {
// 		if strings.Index(ln, "#") == 0 {
// 			count++
// 		}
// 	}
// 	return count, nil
// }

func (c *Copier) copyFile(sourcePath string, targetPath string, excludes []string) error {
	excluded, match := isExcluded(sourcePath, excludes)
	if excluded {
		fmt.Printf("**** %s excluded for copy due to exclude match with: %s", sourcePath, match)
		return nil
	}
	code, err := cfcmd.GenSSHCode()
	if err != nil {
		return errors.Wrap(err, "(copyFile) Cannot generate SSH pass code")
	}
	return c.doCopy(sourcePath, targetPath, code)
}

func (c *Copier) doCopy(sourcePath string, targetPath string, passCode string) error {
	fmt.Printf("Copying %s to %s\n", sourcePath, targetPath)
	c.numInstances
	// cf:23162dc7-ca4d-4e77-b656-65cd6d16ba66/0
	clientConfig := &ssh.ClientConfig{}
	return nil
}

// func (c *Copier) doCopy(sourcePath string, targetPath string, passCode string) error {
// 	cmd := "scp"
// 	for i := 0; i < c.numInstances; i++ {
// 		cmdArgs := []string{"-P",
// 			strconv.Itoa(c.sshPort),
// 			fmt.Sprintf("-oUser=cf:%s/%d", c.AppGUID, i),
// 			sourcePath,
// 			c.sshHost + ":" + targetPath,
// 		}
// 		//Example: scp -P 2222 -oUser=cf:9775c4e8-1ea6-4080-ab45-059e7e640310/0 build.sh ssh.cf.sap.hana.ondemand.com:/tmp
// 		cmdOut, err := execCommandWithStringInput(cmd, cmdArgs, passCode, true)
// 		fmt.Print(cmdOut)
// 		if err != nil {
// 			return errors.Wrapf(err, "Failed copy of %s to %s with code %s", sourcePath, targetPath, passCode)
// 		}
// 	}
// 	return nil
// }

func isExcluded(path string, excludes []string) (excluded bool, match string) {
	for _, x := range excludes {
		if strings.Contains(path, x) {
			excluded, match = true, x
			return
		}
	}
	excluded, match = false, ""
	return
}

// func NewCodeGenerator1(done <-chan struct{}) PasswordGenerator {
// 	fmt.Println("**** Obtaining SSH Code...")
// 	chCodes := make(chan strResult, 4)
// 	go func() {
// 		for {
// 			cmdArgs := []string{"ssh-code"}
// 			out, err := execCommand("cf", cmdArgs, true)
// 			var result strResult
// 			if err != nil {
// 				result.err = err
// 			} else {
// 				result.value = string(out)
// 			}
// 			select {
// 			case chCodes <- result:
// 			case <-done:
// 				fmt.Println("*** Stopping Code Generation")
// 				return
// 			}
// 		}
// 	}()
// 	return func() (code string, err error) {
// 		result, ok := <-chCodes
// 		if !ok {
// 			return "", errors.New("PasswordGenerator terminated")
// 		}
// 		if result.err != nil {
// 			return "", err
// 		}
// 		return result.value, nil
// 	}
// }

type strResult struct {
	value string
	err   error
}

var sshEndpointRe = regexp.MustCompile(`"app_ssh_endpoint": "(.*?)"`)
