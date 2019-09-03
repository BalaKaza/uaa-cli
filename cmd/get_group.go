package cmd

import (
	"errors"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func GetGroupCmd(api *uaa.API, printer cli.Printer, name, attributes string) error {
	group, err := api.GetGroupByName(name, attributes)
	if err != nil {
		return err
	}

	return printer.Print(group)
}

func GetGroupValidations(cfg config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("The positional argument GROUPNAME must be specified.")
	}
	return nil
}

var getGroupCmd = &cobra.Command{
	Use:   "get-group GROUPNAME",
	Short: "Look up a group by group name",
	PreRun: func(cmd *cobra.Command, args []string) {
		cli.NotifyValidationErrors(GetGroupValidations(GetSavedConfig(), args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := GetGroupCmd(GetAPIFromSavedTokenInContext(), cli.NewJsonPrinter(log), args[0], attributes)
		cli.NotifyErrorsWithRetry(err, log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(getGroupCmd)
	getGroupCmd.Annotations = make(map[string]string)
	getGroupCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"

	getGroupCmd.Flags().StringVarP(&attributes, "attributes", "a", "", `include only these comma-separated user attributes to improve query performance`)
	getGroupCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to get the group")
}
