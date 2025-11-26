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

const indexersDataSourceName = "indexers"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexersDataSource{}

func NewIndexersDataSource() datasource.DataSource {
	return &IndexersDataSource{}
}

// IndexersDataSource defines the indexers implementation.
type IndexersDataSource struct {
	client *prowlarr.APIClient
	auth   context.Context
}

// Indexers describes the indexers data model.
type Indexers struct {
	Indexers types.Set    `tfsdk:"indexers"`
	ID       types.String `tfsdk:"id"`
}

func (d *IndexersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexersDataSourceName
}

func (d *IndexersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Indexers -->\nList all available [Indexers](../resources/indexer).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"indexers": schema.SetNestedAttribute{
				MarkdownDescription: "Indexer list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							MarkdownDescription: "Enable RSS flag.",
							Computed:            true,
						},
						"redirect": schema.BoolAttribute{
							MarkdownDescription: "Redirect download request from client to indexer instead of proxying via Prowlarr.",
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
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *IndexersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *IndexersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Indexers

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get indexers current value
	response, _, err := d.client.IndexerAPI.ListIndexer(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexersDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexersDataSourceName)
	// Map response body to resource schema attribute
	indexers := make([]Indexer, len(response))
	for i, t := range response {
		indexers[i].write(ctx, &t, &resp.Diagnostics)
	}

	tfsdk.ValueFrom(ctx, indexers, data.Indexers.Type(ctx), &data.Indexers)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
