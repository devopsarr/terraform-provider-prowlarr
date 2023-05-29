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

func (d *TagDetailsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + tagDetailsDataSourceName
}

func (d *TagDetailsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Tag -->Single [Tag](../resources/tag) with its associated resources.",
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
	var tag *TagDetails

	resp.Diagnostics.Append(req.Config.Get(ctx, &tag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tags current value
	response, _, err := d.client.TagDetailsApi.ListTagDetail(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, tagDetailsDataSourceName, err))

		return
	}

	value, err := findTagDetails(tag.Label.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", tagDetailsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+tagDetailsDataSourceName)
	tag.write(ctx, value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &tag)...)
}

func findTagDetails(label string, tags []*prowlarr.TagDetailsResource) (*prowlarr.TagDetailsResource, error) {
	for _, t := range tags {
		if t.GetLabel() == label {
			return t, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(tagDetailsDataSourceName, "label", label)
}

func (t *TagDetails) write(ctx context.Context, tag *prowlarr.TagDetailsResource) {
	t.ID = types.Int64Value(int64(tag.GetId()))
	t.Label = types.StringValue(tag.GetLabel())
	t.ApplicationIDs, _ = types.SetValueFrom(ctx, types.Int64Type, tag.GetApplicationIds())
	t.IndexerIDs, _ = types.SetValueFrom(ctx, types.Int64Type, tag.GetIndexerIds())
	t.IndexerProxyIDs, _ = types.SetValueFrom(ctx, types.Int64Type, tag.GetIndexerProxyIds())
	t.NotificationIDs, _ = types.SetValueFrom(ctx, types.Int64Type, tag.GetNotificationIds())
}
