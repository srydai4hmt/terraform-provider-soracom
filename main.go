package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/nigamishibumi/terraform-provider-soracom/soracom"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: soracom.Provider,
	})
}
