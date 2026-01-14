package yeet

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Config struct {
	APIKey string
	Host   string
	Client *http.Client
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("YEET_API_KEY", nil),
				Description: "API key for authenticating with the yeet API.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.yeet.cx",
				Description: "Host URL for the yeet API.",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip TLS verification (Useful for localhost development).",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"yeet_host": resourceHost(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	host := d.Get("host").(string)
	insecure := d.Get("insecure").(bool)

	var diags diag.Diagnostics

	client := &http.Client{}
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	config := &Config{
		APIKey: apiKey,
		Host:   host,
		Client: client,
	}

	return config, diags
}
