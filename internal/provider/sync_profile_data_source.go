package provider

import (
	"context"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const syncProfileDataSourceName = "sync_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SyncProfileDataSource{}

func NewSyncProfileDataSource() datasource.DataSource {
	return &SyncProfileDataSource{}
}

// SyncProfileDataSource defines the sync_profile implementation.
type SyncProfileDataSource struct {
	client *prowlarr.APIClient
}

func (d *SyncProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + syncProfileDataSourceName
}

func (d *SyncProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Applications -->Single [Sync Profile](../resources/sync_profile).",
		Attributes: map[string]schema.Attribute{
			"enable_rss": schema.BoolAttribute{
				MarkdownDescription: "Enable RSS flag.",
				Computed:            true,
			},
			"enable_interactive_search": schema.BoolAttribute{
				MarkdownDescription: "Enable interactive search flag.",
				Computed:            true,
			},
			"enable_automatic_search": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic search flag.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Sync Profile ID.",
				Computed:            true,
			},
			"minimum_seeders": schema.Int64Attribute{
				MarkdownDescription: "Minimum seeders.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name.",
				Required:            true,
			},
		},
	}
}

func (d *SyncProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *SyncProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *SyncProfile

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get syncProfile current value
	response, _, err := d.client.AppProfileAPI.ListAppProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, syncProfileDataSourceName, err))

		return
	}

	data.find(data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+syncProfileDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *SyncProfile) find(name string, syncProfiles []prowlarr.AppProfileResource, diags *diag.Diagnostics) {
	for _, profile := range syncProfiles {
		if profile.GetName() == name {
			p.write(&profile)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(syncProfileDataSourceName, "name", name))
}
