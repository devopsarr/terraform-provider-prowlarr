package provider

import (
	"context"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const tagDetailsDataSourceName = "tag_details"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TagDetailsDataSource{}

func NewTagDetailsDataSource() datasource.DataSource {
	return &TagDetailsDataSource{}
}

// TagDetailsDataSource defines the tag details implementation.
type TagDetailsDataSource struct {
	client *prowlarr.APIClient
}

// Tag describes the tag data model.
type TagDetails struct {
	NotificationIDs types.Set    `tfsdk:"notification_ids"`
	IndexerIDs      types.Set    `tfsdk:"indexer_ids"`
	IndexerProxyIDs types.Set    `tfsdk:"indexer_proxy_ids"`
	ApplicationIDs  types.Set    `tfsdk:"application_ids"`
	Label           types.String `tfsdk:"label"`
	ID              types.Int64  `tfsdk:"id"`
}

func (t TagDetails) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"notification_ids":  types.SetType{}.WithElementType(types.Int64Type),
			"indexer_ids":       types.SetType{}.WithElementType(types.Int64Type),
			"indexer_proxy_ids": types.SetType{}.WithElementType(types.Int64Type),
			"application_ids":   types.SetType{}.WithElementType(types.Int64Type),
			"label":             types.StringType,
			"id":                types.Int64Type,
		})
}

func (d *TagDetailsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + tagDetailsDataSourceName
}

func (d *TagDetailsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Tag -->\nSingle [Tag](../resources/tag) with its associated resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Tag ID.",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Tag label.",
				Required:            true,
			},
			"notification_ids": schema.SetAttribute{
				MarkdownDescription: "List of associated notifications.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"indexer_ids": schema.SetAttribute{
				MarkdownDescription: "List of associated indexers.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"indexer_proxy_ids": schema.SetAttribute{
				MarkdownDescription: "List of associated indexer proxies.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"application_ids": schema.SetAttribute{
				MarkdownDescription: "List of associated applications.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (d *TagDetailsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *TagDetailsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *TagDetails

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tags current value
	response, _, err := d.client.TagDetailsAPI.ListTagDetail(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, tagDetailsDataSourceName, err))

		return
	}

	data.find(ctx, data.Label.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+tagDetailsDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (t *TagDetails) find(ctx context.Context, label string, tags []prowlarr.TagDetailsResource, diags *diag.Diagnostics) {
	for _, tag := range tags {
		if tag.GetLabel() == label {
			t.write(ctx, &tag, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(tagDetailsDataSourceName, "label", label))
}

func (t *TagDetails) write(ctx context.Context, tag *prowlarr.TagDetailsResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	t.ID = types.Int64Value(int64(tag.GetId()))
	t.Label = types.StringValue(tag.GetLabel())
	t.ApplicationIDs, tempDiag = types.SetValueFrom(ctx, types.Int64Type, tag.GetApplicationIds())
	diags.Append(tempDiag...)
	t.IndexerIDs, tempDiag = types.SetValueFrom(ctx, types.Int64Type, tag.GetIndexerIds())
	diags.Append(tempDiag...)
	t.IndexerProxyIDs, tempDiag = types.SetValueFrom(ctx, types.Int64Type, tag.GetIndexerProxyIds())
	diags.Append(tempDiag...)
	t.NotificationIDs, tempDiag = types.SetValueFrom(ctx, types.Int64Type, tag.GetNotificationIds())
	diags.Append(tempDiag...)
}
