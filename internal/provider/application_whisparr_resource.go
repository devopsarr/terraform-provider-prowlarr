package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	applicationWhisparrResourceName   = "application_whisparr"
	applicationWhisparrImplementation = "Whisparr"
	applicationWhisparrConfigContract = "WhisparrSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ApplicationWhisparrResource{}
	_ resource.ResourceWithImportState = &ApplicationWhisparrResource{}
)

func NewApplicationWhisparrResource() resource.Resource {
	return &ApplicationWhisparrResource{}
}

// ApplicationWhisparrResource defines the application implementation.
type ApplicationWhisparrResource struct {
	client *prowlarr.APIClient
}

// ApplicationWhisparr describes the application data model.
type ApplicationWhisparr struct {
	SyncCategories types.Set    `tfsdk:"sync_categories"`
	Tags           types.Set    `tfsdk:"tags"`
	Name           types.String `tfsdk:"name"`
	SyncLevel      types.String `tfsdk:"sync_level"`
	ProwlarrURL    types.String `tfsdk:"prowlarr_url"`
	BaseURL        types.String `tfsdk:"base_url"`
	APIKey         types.String `tfsdk:"api_key"`
	ID             types.Int64  `tfsdk:"id"`
}

func (a ApplicationWhisparr) toApplication() *Application {
	return &Application{
		SyncCategories: a.SyncCategories,
		Tags:           a.Tags,
		Name:           a.Name,
		SyncLevel:      a.SyncLevel,
		ProwlarrURL:    a.ProwlarrURL,
		BaseURL:        a.BaseURL,
		APIKey:         a.APIKey,
		ID:             a.ID,
		ConfigContract: types.StringValue(applicationWhisparrConfigContract),
		Implementation: types.StringValue(applicationWhisparrImplementation),
	}
}

func (a *ApplicationWhisparr) fromApplication(application *Application) {
	a.Tags = application.Tags
	a.Name = application.Name
	a.ID = application.ID
	a.SyncLevel = application.SyncLevel
	a.SyncCategories = application.SyncCategories
	a.BaseURL = application.BaseURL
	a.ProwlarrURL = application.ProwlarrURL
	a.APIKey = application.APIKey
}

func (r *ApplicationWhisparrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + applicationWhisparrResourceName
}

func (r *ApplicationWhisparrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Applications -->\nApplication Whisparr resource.\nFor more information refer to [Application](https://wiki.servarr.com/prowlarr/settings#applications) and [Whisparr](https://wiki.servarr.com/prowlarr/supported#whisparr).",
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

func (r *ApplicationWhisparrResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ApplicationWhisparrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var application *ApplicationWhisparr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ApplicationWhisparr
	request := application.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ApplicationAPI.CreateApplications(ctx).ApplicationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, applicationWhisparrResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+applicationWhisparrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	application.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &application)...)
}

func (r *ApplicationWhisparrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var application *ApplicationWhisparr

	resp.Diagnostics.Append(req.State.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ApplicationWhisparr current value
	response, _, err := r.client.ApplicationAPI.GetApplicationsById(ctx, int32(application.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, applicationWhisparrResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+applicationWhisparrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	application.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &application)...)
}

func (r *ApplicationWhisparrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var application *ApplicationWhisparr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ApplicationWhisparr
	request := application.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ApplicationAPI.UpdateApplications(ctx, strconv.Itoa(int(request.GetId()))).ApplicationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, applicationWhisparrResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+applicationWhisparrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	application.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &application)...)
}

func (r *ApplicationWhisparrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ApplicationWhisparr current value
	_, err := r.client.ApplicationAPI.DeleteApplications(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, applicationWhisparrResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+applicationWhisparrResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ApplicationWhisparrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+applicationWhisparrResourceName+": "+req.ID)
}

func (a *ApplicationWhisparr) write(ctx context.Context, application *prowlarr.ApplicationResource, diags *diag.Diagnostics) {
	genericApplication := a.toApplication()
	genericApplication.write(ctx, application, diags)
	a.fromApplication(genericApplication)
}

func (a *ApplicationWhisparr) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.ApplicationResource {
	return a.toApplication().read(ctx, diags)
}
