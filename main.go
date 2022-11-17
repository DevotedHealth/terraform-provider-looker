package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hirosassa/terraform-provider-looker/pkg/looker"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: looker.Provider,
	})
}
