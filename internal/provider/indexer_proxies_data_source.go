package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerProxiesDataSourceName = "indexer_proxies"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerProxiesDataSource{}

func NewIndexerProxiesDataSource() datasource.DataSource {
	return &IndexerProxiesDataSource{}
}

// IndexerProxiesDataSource defines the indexer proxies implementation.
type IndexerProxiesDataSource struct {
	client *prowlarr.APIClient
}

// IndexerProxies describes the indexer proxies data model.
type IndexerProxies struct {
	IndexerProxies types.Set    `tfsdk:"indexer_proxies"`
	ID             types.String `tfsdk:"id"`
}

func (d *IndexerProxiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxiesDataSourceName
}

func (d *IndexerProxiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->List all available [Indexer Proxies](../resources/indexer_proxy).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"indexer_proxies": schema.SetNestedAttribute{
				MarkdownDescription: "Indexer Client list..",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"config_contract": schema.StringAttribute{
							MarkdownDescription: "IndexerProxy configuration template.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "IndexerProxy implementation name.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Indexer Proxy name.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Indexer Proxy ID.",
							Computed:            true,
						},
						// Field values
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port.",
							Computed:            true,
						},
						"request_timeout": schema.Int64Attribute{
							MarkdownDescription: "Request timeout.",
							Computed:            true,
						},
						"host": schema.StringAttribute{
							MarkdownDescription: "host.",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "Username.",
							Computed:            true,
						},
						"password": schema.StringAttribute{
							MarkdownDescription: "Password.",
							Computed:            true,
							Sensitive:           true,
						},
					},
				},
			},
		},
	}
}

func (d *IndexerProxiesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *IndexerProxiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *IndexerProxies

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get indexer proxies current value
	response, _, err := d.client.IndexerProxyApi.ListIndexerProxy(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxiesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxiesDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]IndexerProxy, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.IndexerProxies.Type(context.Background()), &data.IndexerProxies)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
