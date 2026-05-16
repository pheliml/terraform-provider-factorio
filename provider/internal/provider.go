package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-factorio/client"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"rcon_host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FACTORIO_RCON_HOST", nil),
				Description: "The hostname and port for RCON on the Factorio server",
			},
			"rcon_pw": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("FACTORIO_RCON_PW", nil),
				Description: "The password for RCON on the Factorio server",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"factorio_entity": resourceEntity(),
			"factorio_hello":  resourceHello(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"factorio_players": dataSourcePlayers(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	rconHost := d.Get("rcon_host").(string)
	rconPw := d.Get("rcon_pw").(string)
	if rconHost == "" {
		return nil, diag.Errorf("rcon_host was empty")
	}
	if rconPw == "" {
		return nil, diag.Errorf("rcon_pw was empty")
	}
	client, err := client.NewFactorioClient(rconHost, rconPw)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return client, nil
}
