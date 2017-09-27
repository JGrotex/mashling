/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/mashling/cli/cli"
	"github.com/TIBCOSoftware/mashling/lib/util"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build",
	Short:     "Build mashling gateway from mashling.json",
	Long:      "Build mashling gateway from gateway description file - mashling.json",
}

func init() {
	CommandRegistry.RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option *cli.OptionInfo
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdBuild) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
}

// Exec implementation of cli.Command.Exec
func (c *cmdBuild) Exec(args []string) error {

	//Return, if any additanal arguments are passed
	if len(args) != 0 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given. \n\n")
		cmdUsage(c)
	}

	//check whether current directory contains valid mashling gateway project.
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Not able read current directory. \n\n")
		return err
	}
	var gatewayFile = path.Join(currentDir, util.Gateway_Definition_File_Name)
	var bytes []byte
	if b64GatewayJSON := os.Getenv("MASHLING_CONFIG"); b64GatewayJSON != "" {
		fmt.Fprintf(os.Stderr, "Environment variable MASHLING_CONFIG exists, using those contents to overwrite %s\n\n", util.Gateway_Definition_File_Name)
		bytes, err = base64.StdEncoding.DecodeString(b64GatewayJSON)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot read contents of existing MASHLING_CONFIG environment variable: %s\n\n", err.Error())
			os.Exit(1)
		}
		err = ioutil.WriteFile(gatewayFile, bytes, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot write contents of existing MASHLING_CONFIG environment variable to %s: %s\n\n", util.Gateway_Definition_File_Name, err.Error())
			os.Exit(1)
		}
	}
	if !fgutil.FileExists(gatewayFile) {
		fmt.Fprintf(os.Stderr, "Error: Invalid gateway project, didn't find "+gatewayFile+"\n\n")
		return err
	}

	// load gateway descriptor
	gatewayJSON, err := fgutil.LoadLocalFile(gatewayFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Error while loading gateway descriptor file "+gatewayFile+"\n\n")
		return err
	}

	isValidJson := false

	isValidJson, err = IsValidateGateway(gatewayJSON)

	if !isValidJson {
		fmt.Print("Mashling build aborted \n")
		return err
	}

	return BuildMashling(currentDir, gatewayJSON)
}
