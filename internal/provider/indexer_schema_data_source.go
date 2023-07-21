package provider

import (
	"context"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerSchemaDataSourceName = "indexer_schema"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerSchemaDataSource{}

func NewIndexerSchemaDataSource() datasource.DataSource {
	return &IndexerSchemaDataSource{}
}

// IndexerSchemaDataSource defines the indexer schema implementation.
type IndexerSchemaDataSource struct {
	client *prowlarr.APIClient
}

// IndexerSchema describes the indexer data model.
type IndexerSchema struct {
	IndexerURLs    types.Set    `tfsdk:"indexer_urls"`
	LegacyURLs     types.Set    `tfsdk:"legacy_urls"`
	Fields         types.Set    `tfsdk:"fields"`
	ConfigContract types.String `tfsdk:"config_contract"`
	Implementation types.String `tfsdk:"implementation"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Encoding       types.String `tfsdk:"encoding"`
	Protocol       types.String `tfsdk:"protocol"`
	Language       types.String `tfsdk:"language"`
	Privacy        types.String `tfsdk:"privacy"`
	ID             types.Int64  `tfsdk:"id"`
}

// SchemaField is part of IndexerSchema.
type SchemaField struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

func (d *IndexerSchemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerSchemaDataSourceName
}

func (d *IndexerSchemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer schema definition.",
		Attributes: map[string]schema.Attribute{
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
			"description": schema.StringAttribute{
				MarkdownDescription: "Indexer description.",
				Computed:            true,
			},
			"encoding": schema.StringAttribute{
				MarkdownDescription: "Indexer encoding.",
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Computed:            true,
			},
			"indexer_urls": schema.SetAttribute{
				MarkdownDescription: "List of available URLs.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"legacy_urls": schema.SetAttribute{
				MarkdownDescription: "List of legacy URLs.",
				Computed:            true,
				ElementType:         types.StringType,
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
				MarkdownDescription: "Schema ID.",
				Computed:            true,
			},
			"fields": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Set of configuration fields.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: d.getFieldSchema().Attributes,
				},
			},
		},
	}
}

func (d IndexerSchemaDataSource) getFieldSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Field name.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Field description.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Field type.",
				Computed:            true,
			},
		},
	}
}

func (d *IndexerSchemaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *IndexerSchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *IndexerSchema

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get indexers current value
	response, _, err := d.client.IndexerApi.ListIndexerSchema(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerSchemaDataSourceName, err))

		return
	}

	data.find(ctx, data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+indexerSchemaDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (i *IndexerSchema) find(ctx context.Context, name string, schemas []*prowlarr.IndexerResource, diags *diag.Diagnostics) {
	for id, indexer := range schemas {
		if indexer.GetName() == name {
			i.ID = types.Int64Value(int64(id))
			i.write(ctx, indexer, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(indexerSchemaDataSourceName, "name", name))
}

func (i *IndexerSchema) write(ctx context.Context, indexer *prowlarr.IndexerResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	i.ConfigContract = types.StringValue(indexer.GetConfigContract())
	i.Implementation = types.StringValue(indexer.GetImplementation())
	i.Name = types.StringValue(indexer.GetName())
	i.Description = types.StringValue(indexer.GetDescription())
	i.Protocol = types.StringValue(string(indexer.GetProtocol()))
	i.Encoding = types.StringValue(indexer.GetEncoding())
	i.Language = types.StringValue(indexer.GetLanguage())
	i.Privacy = types.StringValue(string(indexer.GetPrivacy()))

	fields := make([]SchemaField, len(indexer.GetFields()))
	for n, f := range indexer.GetFields() {
		fields[n].write(f)
	}

	i.Fields, tempDiag = types.SetValueFrom(ctx, IndexerSchemaDataSource{}.getFieldSchema().Type(), fields)
	diags.Append(tempDiag...)
	i.IndexerURLs, tempDiag = types.SetValueFrom(ctx, types.StringType, indexer.GetIndexerUrls())
	diags.Append(tempDiag...)
	i.LegacyURLs, tempDiag = types.SetValueFrom(ctx, types.StringType, indexer.GetLegacyUrls())
	diags.Append(tempDiag...)
}

func (f *SchemaField) write(field *prowlarr.Field) {
	f.Name = types.StringValue(field.GetName())
	f.Description = types.StringValue(field.GetHelpText())
	f.Type = types.StringValue(field.GetType())
}
