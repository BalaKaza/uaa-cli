package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func SaveContext(ctx config.UaaContext, log *cli.Logger) {
	c := GetSavedConfig()
	c.AddContext(ctx)
	config.WriteConfig(c)
	log.Info("Access token added to active context.")
}

func addImplicitTokenToContext(clientId string, token oauth2.Token, log *cli.Logger) {
	ctx := config.UaaContext{
		GrantType: config.IMPLICIT,
		ClientId:  clientId,
		Token:     token,
	}

	SaveContext(ctx, log)
}

func ImplicitTokenArgumentValidation(cfg config.Config, args []string, port int) error {
	if err := cli.EnsureTargetInConfig(cfg); err != nil {
		return err
	}
	if len(args) < 1 {
		return cli.MissingArgumentError("client_id")
	}
	if port == 0 {
		return cli.MissingArgumentError("port")
	}
	return validateTokenFormatError(tokenFormat)
}

func ImplicitTokenCommandRun(doneRunning chan bool, clientId string, implicitImp cli.ClientImpersonator, log *cli.Logger) {
	implicitImp.Start()
	implicitImp.Authorize()
	tokenResponse := <-implicitImp.Done()
	addImplicitTokenToContext(clientId, tokenResponse, log)
	doneRunning <- true
}

var getImplicitToken = &cobra.Command{
	Use:   "get-implicit-token CLIENT_ID --port REDIRECT_URI_PORT",
	Short: "Obtain an access token using the implicit grant type",
	Long:  help.ImplicitGrant(),
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(ImplicitTokenArgumentValidation(cfg, args, port), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan bool)
		baseUrl := GetSavedConfig().GetActiveTarget().BaseUrl
		implicitImp := cli.NewImplicitClientImpersonator(args[0], baseUrl, tokenFormat, scope, port, log, open.Run)
		go ImplicitTokenCommandRun(done, args[0], implicitImp, GetLogger())
		<-done
	},
}

func init() {
	getImplicitToken.Flags().IntVarP(&port, "port", "", 0, "port on which to run local callback server")
	getImplicitToken.Flags().StringVarP(&scope, "scope", "", "openid", "comma-separated scopes to request in token")
	getImplicitToken.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
	getImplicitToken.Annotations = make(map[string]string)
	getImplicitToken.Annotations[TOKEN_CATEGORY] = "true"
	RootCmd.AddCommand(getImplicitToken)
}
