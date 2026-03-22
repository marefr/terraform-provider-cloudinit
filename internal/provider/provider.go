package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &CloudinitProvider{}

type CloudinitProvider struct {
	version string
}

func (p *CloudinitProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudinit"
	resp.Version = p.version
}

func (p *CloudinitProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for generating cloud-init NoCloud ISO images.",
		MarkdownDescription: `
The cloudinit provider is used to generate cloud-init ISO images, suitable for use with the
[NoCloud data source](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html) and the
[drive with labeled filesystem](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html#source-2-drive-with-labeled-filesystem)
configuration source.

This provider requires no configuration. For information on the resources it provides, see the navigation bar.
`,
	}
}

func (p *CloudinitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Nothing to configure.
}

func (p *CloudinitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *CloudinitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCloudInitISOResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CloudinitProvider{
			version: version,
		}
	}
}
