package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const syncProfileResourceName = "sync_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SyncProfileResource{}
	_ resource.ResourceWithImportState = &SyncProfileResource{}
)

func NewSyncProfileResource() resource.Resource {
	return &SyncProfileResource{}
}

// SyncProfileResource defines the sync profile implementation.
type SyncProfileResource struct {
	client *prowlarr.APIClient
}

// SyncProfile describes the sync profile data model.
type SyncProfile struct {
	Name                    types.String `tfsdk:"name"`
	ID                      types.Int64  `tfsdk:"id"`
	MinimumSeeders          types.Int64  `tfsdk:"minimum_seeders"`
	EnableRss               types.Bool   `tfsdk:"enable_rss"`
	EnableInteractiveSearch types.Bool   `tfsdk:"enable_interactive_search"`
	EnableAutomaticSearch   types.Bool   `tfsdk:"enable_automatic_search"`
}

func (r *SyncProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + syncProfileResourceName
}

func (r *SyncProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Applications -->Sync Profile resource.\nFor more information refer to [Sync Profiles](https://wiki.servarr.com/prowlarr/settings#sync-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"enable_rss": schema.BoolAttribute{
				MarkdownDescription: "Enable RSS flag.",
				Required:            true,
			},
			"enable_interactive_search": schema.BoolAttribute{
				MarkdownDescription: "Enable interactive search flag.",
				Required:            true,
			},
			"enable_automatic_search": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic search flag.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Sync Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"minimum_seeders": schema.Int64Attribute{
				MarkdownDescription: "Minimum seeders.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name.",
				Required:            true,
			},
		},
	}
}

func (r *SyncProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *SyncProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *SyncProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Sync Profile
	request := profile.read()

	response, _, err := r.client.AppProfileApi.CreateAppProfile(ctx).AppProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, syncProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created sync profile: "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *SyncProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *SyncProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get sync profile current value
	response, _, err := r.client.AppProfileApi.GetAppProfileById(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, syncProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+syncProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	profile.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *SyncProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *SyncProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update SyncProfile
	request := profile.read()

	response, _, err := r.client.AppProfileApi.UpdateAppProfile(ctx, fmt.Sprint(request.GetId())).AppProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, syncProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+syncProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *SyncProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete sync profile current value
	_, err := r.client.AppProfileApi.DeleteAppProfile(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, syncProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+syncProfileResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *SyncProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+syncProfileResourceName+": "+req.ID)
}

func (s *SyncProfile) read() *prowlarr.AppProfileResource {
	profile := *prowlarr.NewAppProfileResource()
	profile.SetName(s.Name.ValueString())
	profile.SetId(int32(s.ID.ValueInt64()))
	profile.SetMinimumSeeders(int32(s.MinimumSeeders.ValueInt64()))
	profile.SetEnableRss(s.EnableRss.ValueBool())
	profile.SetEnableInteractiveSearch(s.EnableInteractiveSearch.ValueBool())
	profile.SetEnableAutomaticSearch(s.EnableAutomaticSearch.ValueBool())

	return &profile
}

func (s *SyncProfile) write(profile *prowlarr.AppProfileResource) {
	s.ID = types.Int64Value(int64(profile.GetId()))
	s.Name = types.StringValue(profile.GetName())
	s.MinimumSeeders = types.Int64Value(int64(profile.GetMinimumSeeders()))
	s.EnableRss = types.BoolValue(profile.GetEnableRss())
	s.EnableInteractiveSearch = types.BoolValue(profile.GetEnableInteractiveSearch())
	s.EnableAutomaticSearch = types.BoolValue(profile.GetEnableAutomaticSearch())
}
