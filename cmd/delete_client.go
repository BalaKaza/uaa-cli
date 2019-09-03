package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func DeleteClientValidations(cfg config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) == 0 {
		return cli.MissingArgumentError("client_id")
	}
	return nil
}

func DeleteClientCmd(api *uaa.API, clientId string) error {
	_, err := api.DeleteClient(clientId)
	if err != nil {
		return err
	}

	log.Infof("Successfully deleted client %v.", utils.Emphasize(clientId))
	return nil
}

var deleteClientCmd = &cobra.Command{
	Use:   "delete-client CLIENT_ID",
	Short: "Delete a client registration",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(DeleteClientValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		api := NewApiFromSavedConfig()
		cli.NotifyErrorsWithRetry(DeleteClientCmd(api, args[0]), log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(deleteClientCmd)
	deleteClientCmd.Annotations = make(map[string]string)
	deleteClientCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	deleteClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to delete the client")
}
