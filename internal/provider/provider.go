package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider              = &AG5Provider{}
	_ provider.ProviderWithFunctions = &AG5Provider{}
)

// AG5Provider provides deterministic, local-only helper functions for AG5
// Terraform configurations.
type AG5Provider struct {
	version string
}

func (p *AG5Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ag5"
	resp.Version = p.version
}

func (p *AG5Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *AG5Provider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (p *AG5Provider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *AG5Provider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

func (p *AG5Provider) Functions(context.Context) []func() function.Function {
	return []func() function.Function{
		NewShortenFunction,
	}
}

// New returns a factory for an AG5 provider at the supplied version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AG5Provider{version: version}
	}
}
