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

const tagsDetailsDataSourceName = "tags_details"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TagsDetailsDataSource{}

func NewTagsDetailsDataSource() datasource.DataSource {
	return &TagsDetailsDataSource{}
}

// TagsDetailsDataSource defines the tags details implementation.
type TagsDetailsDataSource struct {
	client *prowlarr.APIClient
}

// Tags describes the tags data model.
type TagsDetails struct {
	Tags types.Set    `tfsdk:"tags"`
	ID   types.String `tfsdk:"id"`
}

func (d *TagsDetailsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + tagsDetailsDataSourceName
}

func (d *TagsDetailsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Tag -->[Tag](../resources/tag) list with their associated resources.",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tag list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Tags ID.",
							Computed:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "Tags label.",
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
				},
			},
		},
	}
}

func (d *TagsDetailsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *TagsDetailsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get tagss current value
	response, _, err := d.client.TagDetailsApi.ListTagDetail(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, tagsDetailsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+tagsDetailsDataSourceName)
	// Map response body to resource schema attribute
	tags := make([]TagDetails, len(response))
	for i, t := range response {
		tags[i].write(ctx, t, &resp.Diagnostics)
	}

	tagList, diags := types.SetValueFrom(ctx, TagDetails{}.getType(), tags)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, TagsDetails{Tags: tagList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
