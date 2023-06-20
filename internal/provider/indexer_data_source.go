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

const indexerDataSourceName = "indexer"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerDataSource{}

func NewIndexerDataSource() datasource.DataSource {
	return &IndexerDataSource{}
}

// IndexerDataSource defines the indexer implementation.
type IndexerDataSource struct {
	client *prowlarr.APIClient
}

func (d *IndexerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerDataSourceName
}

func (d *IndexerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Indexers -->Single [Indexer](../resources/indexer).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable RSS flag.",
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Computed:            true,
			},
			"app_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Application profile ID.",
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Indexer configuration template.",
				Computed:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Indexer implementation name.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Indexer name.",
				Required:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"language": schema.StringAttribute{
				MarkdownDescription: "Language.",
				Computed:            true,
			},
			"privacy": schema.StringAttribute{
				MarkdownDescription: "Privacy.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Indexer ID.",
				Computed:            true,
			},
			"fields": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Set of configuration fields.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Field name.",
							Computed:            true,
						},
						"text_value": schema.StringAttribute{
							MarkdownDescription: "Text value.",
							Computed:            true,
						},
						"sensitive_value": schema.StringAttribute{
							MarkdownDescription: "Sensitive string value.",
							Computed:            true,
							Sensitive:           true,
						},
						"number_value": schema.NumberAttribute{
							MarkdownDescription: "Number value.",
							Computed:            true,
						},
						"bool_value": schema.BoolAttribute{
							MarkdownDescription: "Bool value.",
							Computed:            true,
						},
						"set_value": schema.SetAttribute{
							MarkdownDescription: "Set value.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *IndexerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *IndexerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var indexer *Indexer

	resp.Diagnostics.Append(req.Config.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get indexers current value
	response, _, err := d.client.IndexerApi.ListIndexer(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerDataSourceName, err))

		return
	}

	value, err := findIndexer(indexer.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", indexerDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerDataSourceName)
	indexer.write(ctx, value, &resp.Diagnostics)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func findIndexer(name string, indexers []*prowlarr.IndexerResource) (*prowlarr.IndexerResource, error) {
	for _, t := range indexers {
		if t.GetName() == name {
			return t, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(indexerDataSourceName, "name", name)
}
