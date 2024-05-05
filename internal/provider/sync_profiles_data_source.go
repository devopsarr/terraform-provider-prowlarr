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

const syncProfilesDataSourceName = "sync_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SyncProfilesDataSource{}

func NewSyncProfilesDataSource() datasource.DataSource {
	return &SyncProfilesDataSource{}
}

// SyncProfilesDataSource defines the sync profiles implementation.
type SyncProfilesDataSource struct {
	client *prowlarr.APIClient
}

// SyncProfiles describes the sync profiles data model.
type SyncProfiles struct {
	SyncProfiles types.Set    `tfsdk:"sync_profiles"`
	ID           types.String `tfsdk:"id"`
}

func (d *SyncProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + syncProfilesDataSourceName
}

func (d *SyncProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Applications -->\nList all available [Sync Profiles](../resources/sync_profile).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"sync_profiles": schema.SetNestedAttribute{
				MarkdownDescription: "Sync Profile list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
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
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SyncProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *SyncProfilesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get sync profiles current value
	response, _, err := d.client.AppProfileAPI.ListAppProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, syncProfilesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+syncProfilesDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]SyncProfile, len(response))
	for i, p := range response {
		profiles[i].write(&p)
	}

	profileList, diags := types.SetValueFrom(ctx, SyncProfile{}.getType(), profiles)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, SyncProfiles{SyncProfiles: profileList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
