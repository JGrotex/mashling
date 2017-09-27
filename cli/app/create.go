/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"flag"
	"fmt"
	"os"

	"encoding/base64"
	"encoding/json"
	"path"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/mashling/cli/cli"
	"github.com/TIBCOSoftware/mashling/lib/model"
	"github.com/TIBCOSoftware/mashling/lib/types"
)

var optCreate = &cli.OptionInfo{
	Name:      "create",
	UsageLine: "create AppName",
	Short:     "create a mashling gateway",
	Long: `Creates a mashling gateway.

Options:
    -f       specify the mashling.json to create gateway project from
 `,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option    *cli.OptionInfo
	fileName  string
	vendorDir string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdCreate) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.fileName), "f", "", "gateway app file")
}

// Exec implementation of cli.Command.Exec
func (c *cmdCreate) Exec(args []string) error {

	var gatewayJSON string
	var gatewayName string
	var err error
	var bytes []byte
	var mashling types.Microgateway

	if c.fileName != "" {

		if fgutil.IsRemote(c.fileName) {

			gatewayJSON, err = fgutil.LoadRemoteFile(c.fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", c.fileName, err.Error())
				os.Exit(2)
			}
		} else {
			gatewayJSON, err = fgutil.LoadLocalFile(c.fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", c.fileName, err.Error())
				os.Exit(2)
			}

			if len(args) != 0 {
				gatewayName = args[0]
			}
		}
	} else {
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "Error: Gateway name not specified\n\n")
			cmdUsage(c)
		}

		if len(args) != 1 {
			fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
			cmdUsage(c)
		}

		gatewayName = args[0]

		if b64GatewayJSON := os.Getenv("MASHLING_CONFIG"); b64GatewayJSON != "" {
			fmt.Fprintf(os.Stderr, "Environment variable MASHLING_CONFIG exists, using those contents...\n\n")
			bytes, err = base64.StdEncoding.DecodeString(b64GatewayJSON)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Cannot read contents of existing MASHLING_CONFIG environment variable: %s\n\n", err.Error())
				os.Exit(1)
			}
		} else {
			mashling, err = model.CreateMashlingSampleModel()
			if err != nil {
				return err
			}
			bytes, err = json.MarshalIndent(mashling, "", "\t")
			if err != nil {
				return err
			}
		}
		gatewayJSON = string(bytes)
	}

	currentDir, err := os.Getwd()

	if err != nil {
		return err
	}

	appDir := path.Join(currentDir, gatewayName)

	isValidJSON, invalidErr := IsValidateGateway(gatewayJSON)

	if !isValidJSON {
		fmt.Fprintf(os.Stderr, "Error: Provided JSON file is invalid: %s\n\n", invalidErr.Error())
		os.Exit(1)
	}

	return CreateMashling(SetupNewProjectEnv(), gatewayJSON, appDir, gatewayName, c.vendorDir)
}
