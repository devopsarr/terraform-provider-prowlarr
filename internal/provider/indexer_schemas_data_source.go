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

const indexerSchemasDataSourceName = "indexer_schemas"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerSchemasDataSource{}

func NewIndexerSchemasDataSource() datasource.DataSource {
	return &IndexerSchemasDataSource{}
}

// IndexerSchemasDataSource defines the indexers implementation.
type IndexerSchemasDataSource struct {
	client *prowlarr.APIClient
}

// IndexerSchemas describes the indexers data model.
type IndexerSchemas struct {
	IndexerSchemas types.List   `tfsdk:"indexer_schemas"`
	ID             types.String `tfsdk:"id"`
}

func (d *IndexerSchemasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerSchemasDataSourceName
}

func (d *IndexerSchemasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Indexers -->List all available [Indexer Schemas](../data-sources/indexer_schema).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"indexer_schemas": schema.ListAttribute{
				MarkdownDescription: "Indexer name list.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *IndexerSchemasDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *IndexerSchemasDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexers current value
	response, _, err := d.client.IndexerApi.ListIndexerSchema(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerSchemasDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexersDataSourceName)
	// Map response body to resource schema attribute
	indexers := make([]string, len(response))
	for i, t := range response {
		indexers[i] = t.GetName()
	}

	indexerList, diags := types.ListValueFrom(ctx, types.StringType, indexers)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, IndexerSchemas{IndexerSchemas: indexerList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
