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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	applicationRadarrResourceName   = "application_radarr"
	applicationRadarrImplementation = "Radarr"
	applicationRadarrConfigContract = "RadarrSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ApplicationRadarrResource{}
	_ resource.ResourceWithImportState = &ApplicationRadarrResource{}
)

func NewApplicationRadarrResource() resource.Resource {
	return &ApplicationRadarrResource{}
}

// ApplicationRadarrResource defines the application implementation.
type ApplicationRadarrResource struct {
	client *prowlarr.APIClient
}

// ApplicationRadarr describes the application data model.
type ApplicationRadarr struct {
	SyncCategories types.Set    `tfsdk:"sync_categories"`
	Tags           types.Set    `tfsdk:"tags"`
	Name           types.String `tfsdk:"name"`
	SyncLevel      types.String `tfsdk:"sync_level"`
	ProwlarrURL    types.String `tfsdk:"prowlarr_url"`
	BaseURL        types.String `tfsdk:"base_url"`
	APIKey         types.String `tfsdk:"api_key"`
	ID             types.Int64  `tfsdk:"id"`
}

func (a ApplicationRadarr) toApplication() *Application {
	return &Application{
		SyncCategories: a.SyncCategories,
		Tags:           a.Tags,
		Name:           a.Name,
		SyncLevel:      a.SyncLevel,
		ProwlarrURL:    a.ProwlarrURL,
		BaseURL:        a.BaseURL,
		APIKey:         a.APIKey,
		ID:             a.ID,
		ConfigContract: types.StringValue(applicationRadarrConfigContract),
		Implementation: types.StringValue(applicationRadarrImplementation),
	}
}

func (a *ApplicationRadarr) fromApplication(application *Application) {
	a.Tags = application.Tags
	a.Name = application.Name
	a.ID = application.ID
	a.SyncLevel = application.SyncLevel
	a.SyncCategories = application.SyncCategories
	a.BaseURL = application.BaseURL
	a.ProwlarrURL = application.ProwlarrURL
	a.APIKey = application.APIKey
}

func (r *ApplicationRadarrResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + applicationRadarrResourceName
}

func (r *ApplicationRadarrResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Applications -->Application Radarr resource.\nFor more information refer to [Application](https://wiki.servarr.com/prowlarr/settings#applications) and [Radarr](https://wiki.servarr.com/prowlarr/supported#radarr).",
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
		},
	}
}

func (r *ApplicationRadarrResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ApplicationRadarrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var application *ApplicationRadarr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ApplicationRadarr
	request := application.read(ctx)

	response, _, err := r.client.ApplicationApi.CreateApplications(ctx).ApplicationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, applicationRadarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+applicationRadarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	application.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &application)...)
}

func (r *ApplicationRadarrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var application *ApplicationRadarr

	resp.Diagnostics.Append(req.State.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ApplicationRadarr current value
	response, _, err := r.client.ApplicationApi.GetApplicationsById(ctx, int32(application.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, applicationRadarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+applicationRadarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	application.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &application)...)
}

func (r *ApplicationRadarrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var application *ApplicationRadarr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ApplicationRadarr
	request := application.read(ctx)

	response, _, err := r.client.ApplicationApi.UpdateApplications(ctx, strconv.Itoa(int(request.GetId()))).ApplicationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, applicationRadarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+applicationRadarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	application.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &application)...)
}

func (r *ApplicationRadarrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ApplicationRadarr current value
	_, err := r.client.ApplicationApi.DeleteApplications(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, applicationRadarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+applicationRadarrResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ApplicationRadarrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+applicationRadarrResourceName+": "+req.ID)
}

func (a *ApplicationRadarr) write(ctx context.Context, application *prowlarr.ApplicationResource) {
	genericApplication := a.toApplication()
	genericApplication.write(ctx, application)
	a.fromApplication(genericApplication)
}

func (a *ApplicationRadarr) read(ctx context.Context) *prowlarr.ApplicationResource {
	return a.toApplication().read(ctx)
}
