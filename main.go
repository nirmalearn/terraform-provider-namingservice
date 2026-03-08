package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"namingservice_name": resourceName(),
		},
	}
}

func resourceName() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNameCreate,
		ReadContext:   resourceNameRead,
		DeleteContext: resourceNameDelete,
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Full HTTP endpoint for name generation, for example http://localhost:8000/generate-name.",
			},
			"app": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application short code.",
			},
			"env": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Environment short code.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region short code.",
			},
			"csp": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud service provider short code.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource type short code.",
			},
			"generated_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Generated resource name returned by the API.",
			},
		},
	}
}

func resourceNameCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	payload := map[string]string{
		"app":           d.Get("app").(string),
		"env":           d.Get("env").(string),
		"region":        d.Get("region").(string),
		"csp":           d.Get("csp").(string),
		"resource_type": d.Get("resource_type").(string),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal request body: %w", err))
	}

	apiURL := d.Get("api_url").(string)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(body))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create request: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to call naming service: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read naming service response: %w", err))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.FromErr(fmt.Errorf("naming service returned status %d: %s", resp.StatusCode, string(respBody)))
	}

	var out struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode naming service response: %w", err))
	}
	if out.Name == "" {
		return diag.FromErr(fmt.Errorf("naming service response did not include a non-empty name"))
	}

	if err := d.Set("generated_name", out.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set generated_name: %w", err))
	}

	d.SetId(out.Name)
	return nil
}

func resourceNameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNameDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
