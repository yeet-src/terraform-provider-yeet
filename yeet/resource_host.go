package yeet

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// pruneKeySizeBytes is the size of the random prune key.
const pruneKeySizeBytes = 32

func resourceHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		DeleteContext: resourceHostDelete,
		Schema: map[string]*schema.Schema{
			"prune_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				Description: "The host's prune key. If not provided, a random key will be generated.",
			},
		},
	}
}

func resourceHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	var pruneKey string

	// Check if prune_key was provided by the user
	if v, ok := d.GetOk("prune_key"); ok {
		pruneKey = v.(string)
	} else {
		// Generate a random pruneKey if not provided
		pruneKeyBytes := make([]byte, pruneKeySizeBytes)
		if _, err := rand.Read(pruneKeyBytes); err != nil {
			return diag.FromErr(fmt.Errorf("failed to generate random pruneKey: %w", err))
		}
		pruneKey = hex.EncodeToString(pruneKeyBytes)
	}

	// Store the prune key in state
	if err := d.Set("prune_key", pruneKey); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pruneKey)

	return diags
}

func resourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Nothing to read from the API, just ensure the key is still in state
	if d.Id() == "" {
		return diag.Errorf("resource ID is empty")
	}

	return diags
}

func resourceHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)

	pruneKey := d.Get("prune_key").(string)

	// Call the prune endpoint with the api_key and the generated prune key
	url := fmt.Sprintf("%s/hosts/prune", config.Host)

	payload := map[string]string{
		"prune_key": pruneKey,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal request payload: %w", err))
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create request: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))

	resp, err := config.Client.Do(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to call prune endpoint: %w", err))
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return diag.FromErr(fmt.Errorf("prune endpoint returned error status %d: %s", resp.StatusCode, string(body)))
	}

	d.SetId("")
	return diags
}
