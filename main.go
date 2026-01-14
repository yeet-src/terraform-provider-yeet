package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/yeet-src/terraform-provider-yeet/yeet"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: yeet.Provider,
	})
}
