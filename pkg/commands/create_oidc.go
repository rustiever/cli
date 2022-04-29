package commands

import (
	"fmt"

	commonv3 "github.com/RafayLabs/rcloud-base/proto/types/commonpb/v3"
	systemv3 "github.com/RafayLabs/rcloud-base/proto/types/systempb/v3"
	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/oidc"
	"github.com/spf13/cobra"
)

const (
	ClientIDFlag    = "clientid"
	CallbackUrlFlag = "callback-url"
	ScopesFlag      = "scopes"
)

type CreateOIDCProviderOptions struct {
	clientId    string
	callbackUrl string
	scopes      []string
	logger      log.Logger
	config      *config.Config
}

func NewCreateOIDCProviderOptions(logger log.Logger, config *config.Config) *CreateOIDCProviderOptions {
	o := new(CreateOIDCProviderOptions)
	o.logger = logger
	o.config = config
	return o
}

func (c *CreateOIDCProviderOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.PersistentFlags()
	flagSet.StringVar(&c.clientId, ClientIDFlag, "",
		"Client Id generated during for Oauth provider registration.")
	flagSet.StringVar(&c.callbackUrl, CallbackUrlFlag, "",
		"Callback URL to be configured during Oauth Registration.")
	flagSet.StringSliceVar(&c.scopes, ScopesFlag, nil,
		"Scopes that are required from OIDC Provider.")
}

func (o *CreateOIDCProviderOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MinimumNArgs(1)(cmd, args)

}

func (o *CreateOIDCProviderOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())

	name := args[0]
	flagSet := cmd.Flags()
	err := fmt.Errorf("flags not triggered")

	if flagSet.Changed(ClientIDFlag) && flagSet.Changed(CallbackUrlFlag) && flagSet.Changed(ScopesFlag) {
		oidcProvider := &systemv3.OIDCProvider{
			Kind: "OIDCProvider",
			Metadata: &commonv3.Metadata{
				Name:         name,
				Organization: config.GetConfig().Organization,
			},
			Spec: &systemv3.OIDCProviderSpec{
				ProviderName: name,
				ClientId:     o.clientId,
				CallbackUrl:  o.callbackUrl,
				Scopes:       o.scopes,
				Predefined:   false,
			},
		}

		err := oidc.CreateOIDCProvider(oidcProvider)
		if err != nil {
			return fmt.Errorf("failed to create oidc provider, cause: %s", err.Error())
		}
	} else {
		return err
	}
	return nil
}