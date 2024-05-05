package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (d *IndexerProxiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxiesDataSourceName
}

func (d *IndexerProxiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->\nList all available [Indexer Proxies](../resources/indexer_proxy).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"indexer_proxies": schema.SetNestedAttribute{
				MarkdownDescription: "Indexer Client list.",
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

func (d *IndexerProxiesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer proxies current value
	response, _, err := d.client.IndexerProxyAPI.ListIndexerProxy(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxiesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxiesDataSourceName)
	// Map response body to resource schema attribute
	proxies := make([]IndexerProxy, len(response))
	for i, p := range response {
		proxies[i].write(ctx, &p, &resp.Diagnostics)
	}

	proxyList, diags := types.SetValueFrom(ctx, IndexerProxy{}.getType(), proxies)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, IndexerProxies{IndexerProxies: proxyList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
