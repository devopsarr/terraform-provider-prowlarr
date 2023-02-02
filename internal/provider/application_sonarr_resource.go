package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	applicationSonarrResourceName   = "application_sonarr"
	applicationSonarrImplementation = "Sonarr"
	applicationSonarrConfigContract = "SonarrSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ApplicationSonarrResource{}
	_ resource.ResourceWithImportState = &ApplicationSonarrResource{}
)

func NewApplicationSonarrResource() resource.Resource {
	return &ApplicationSonarrResource{}
}

// ApplicationSonarrResource defines the application implementation.
type ApplicationSonarrResource struct {
	client *prowlarr.APIClient
}

// ApplicationSonarr describes the application data model.
type ApplicationSonarr struct {
	SyncCategories      types.Set    `tfsdk:"sync_categories"`
	AnimeSyncCategories types.Set    `tfsdk:"anime_sync_categories"`
	Tags                types.Set    `tfsdk:"tags"`
	Name                types.String `tfsdk:"name"`
	SyncLevel           types.String `tfsdk:"sync_level"`
	ProwlarrURL         types.String `tfsdk:"prowlarr_url"`
	BaseURL             types.String `tfsdk:"base_url"`
	APIKey              types.String `tfsdk:"api_key"`
	ID                  types.Int64  `tfsdk:"id"`
}

func (n ApplicationSonarr) toApplication() *Application {
	return &Application{
		SyncCategories:      n.SyncCategories,
		AnimeSyncCategories: n.AnimeSyncCategories,
		Tags:                n.Tags,
		Name:                n.Name,
		SyncLevel:           n.SyncLevel,
		ProwlarrURL:         n.ProwlarrURL,
		BaseURL:             n.BaseURL,
		APIKey:              n.APIKey,
		ID:                  n.ID,
	}
}

func (n *ApplicationSonarr) fromApplication(application *Application) {
	n.Tags = application.Tags
	n.Name = application.Name
	n.ID = application.ID
	n.SyncLevel = application.SyncLevel
	n.SyncCategories = application.SyncCategories
	n.AnimeSyncCategories = application.AnimeSyncCategories
	n.BaseURL = application.BaseURL
	n.ProwlarrURL = application.ProwlarrURL
	n.APIKey = application.APIKey
}

func (r *ApplicationSonarrResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + applicationSonarrResourceName
}

func (r *ApplicationSonarrResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Applications -->Application Sonarr resource.\nFor more information refer to [Application](https://wiki.servarr.com/prowlarr/settings#applications) and [Sonarr](https://wiki.servarr.com/prowlarr/supported#sonarr).",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Application name.",
				Required:            true,
			},
			"sync_level": schema.StringAttribute{
				MarkdownDescription: "Sync level.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("addOnly", "disabled", "fullSync"),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Application ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Required:            true,
			},
			"prowlarr_url": schema.StringAttribute{
				MarkdownDescription: "Prowlarr URL.",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"sync_categories": schema.SetAttribute{
				MarkdownDescription: "Sync categories.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"anime_sync_categories": schema.SetAttribute{
				MarkdownDescription: "Anime sync categories.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *ApplicationSonarrResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ApplicationSonarrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var application *ApplicationSonarr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ApplicationSonarr
	request := application.read(ctx)

	response, _, err := r.client.ApplicationApi.CreateApplications(ctx).ApplicationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, applicationSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+applicationSonarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	application.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &application)...)
}

func (r *ApplicationSonarrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var application *ApplicationSonarr

	resp.Diagnostics.Append(req.State.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ApplicationSonarr current value
	response, _, err := r.client.ApplicationApi.GetApplicationsById(ctx, int32(application.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, applicationSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+applicationSonarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	application.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &application)...)
}

func (r *ApplicationSonarrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var application *ApplicationSonarr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ApplicationSonarr
	request := application.read(ctx)

	response, _, err := r.client.ApplicationApi.UpdateApplications(ctx, strconv.Itoa(int(request.GetId()))).ApplicationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, applicationSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+applicationSonarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	application.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &application)...)
}

func (r *ApplicationSonarrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var application *ApplicationSonarr

	resp.Diagnostics.Append(req.State.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ApplicationSonarr current value
	_, err := r.client.ApplicationApi.DeleteApplications(ctx, int32(application.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, applicationSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+applicationSonarrResourceName+": "+strconv.Itoa(int(application.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ApplicationSonarrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+applicationSonarrResourceName+": "+req.ID)
}

func (n *ApplicationSonarr) write(ctx context.Context, application *prowlarr.ApplicationResource) {
	genericApplication := Application{
		SyncLevel: types.StringValue(string(application.GetSyncLevel())),
		ID:        types.Int64Value(int64(application.GetId())),
		Name:      types.StringValue(application.GetName()),
		Tags:      types.SetValueMust(types.Int64Type, nil),
	}
	tfsdk.ValueFrom(ctx, application.Tags, genericApplication.Tags.Type(ctx), &genericApplication.Tags)
	genericApplication.writeFields(ctx, application.GetFields())
	n.fromApplication(&genericApplication)
}

func (n *ApplicationSonarr) read(ctx context.Context) *prowlarr.ApplicationResource {
	tags := make([]*int32, len(n.Tags.Elements()))
	tfsdk.ValueAs(ctx, n.Tags, &tags)

	application := prowlarr.NewApplicationResource()
	application.SetSyncLevel(prowlarr.ApplicationSyncLevel(n.SyncLevel.ValueString()))
	application.SetId(int32(n.ID.ValueInt64()))
	application.SetName(n.Name.ValueString())
	application.SetConfigContract(applicationSonarrConfigContract)
	application.SetImplementation(applicationSonarrImplementation)
	application.SetTags(tags)
	application.SetFields(n.toApplication().readFields(ctx))

	return application
}
