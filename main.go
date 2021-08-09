package main

import (
	"github.com/chriskuchin/terraform-provider-looker/pkg/looker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: looker.Provider,
	})
}
