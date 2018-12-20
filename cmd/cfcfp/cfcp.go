package main

import (
	"log"
	"os"

	"github.com/elankath/cftool/pkg/cfcp"

	"github.com/jawher/mow.cli"
)

func init() {
	log.SetFlags(0)
}

func main() {
	app := cli.App("cfcp",
		`Copy files from Local FS to CF Container FS.  
Requires that the Cloud Foundry Client be installed and available in the path and
that you have logged into your relevant CF API endpoint (cloud controller) via 'cf login'

Notes:
* This tool does NOT support wildcards currently.
* This tool does NOT create TRG as dir in container FS if it doesn't exist!
	(feature will be added later when time permits)
* Follow https://docs.cloudfoundry.org/cf-cli/install-go-cli.html for installing the cf CLI`)

	app.Spec = "-a=<cfAppName> [-x=<excludePath>] [-s=<simulated>] [SRC] TRG"

	var (
		appName   = app.StringOpt("a", "", "CF Application Name")
		excludes  = app.StringsOpt("x", nil, "Exclude Path from Copy. Can be any partial path segment. Can be repeated")
		simulated = app.BoolOpt("s", false, "Simulates a copy by printing out what will be copied, but does not really copy")
		source    = app.StringArg("SRC", ".", "Source Path (file/dir) to Copy. (Default: current dir). All files in a Dir Path are copied")
		target    = app.StringArg("TRG", "", "Target Path in CF container. Created if not existing")
	)

	// Specify the action to execute when the app is invoked correctly
	app.Action = func() {
		config := cfcp.CopierConfig{
			AppName: *appName,
			Simulated: *simulated,
		}
		copier, err := cfcp.NewCopier(config)
		if err != nil {
			log.Fatalf("**** Failed to create copier: %s\n", err)
		}
		err = copier.Copy(*source, *target, *excludes)
		if err != nil {
			log.Fatalf("**** Copy failed: %s\n", err)
		}
	}

	// Invoke the app passing in os.Args
	app.Run(os.Args)
}
