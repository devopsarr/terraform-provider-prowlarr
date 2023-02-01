package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerProxyDataSourceName = "indexer_proxy"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerProxyDataSource{}

func NewIndexerProxyDataSource() datasource.DataSource {
	return &IndexerProxyDataSource{}
}

// IndexerProxyDataSource defines the indexer_proxy implementation.
type IndexerProxyDataSource struct {
	client *prowlarr.APIClient
}

func (i *IndexerProxyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxyDataSourceName
}

func (i *IndexerProxyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->Single [Indexer Proxy](../resources/indexer_proxy).",
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
				Required:            true,
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
	}
}

func (i *IndexerProxyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		i.client = client
	}
}

func (i *IndexerProxyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *IndexerProxy

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get indexerProxy current value
	response, _, err := i.client.IndexerProxyApi.ListIndexerProxy(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxyDataSourceName, err))

		return
	}

	indexerProxy, err := findIndexerProxy(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", indexerProxyDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxyDataSourceName)
	data.write(ctx, indexerProxy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findIndexerProxy(name string, indexerProxys []*prowlarr.IndexerProxyResource) (*prowlarr.IndexerProxyResource, error) {
	for _, i := range indexerProxys {
		if i.GetName() == name {
			return i, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(indexerProxyDataSourceName, "name", name)
}
