package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	notificationGotifyResourceName   = "notification_gotify"
	notificationGotifyImplementation = "Gotify"
	notificationGotifyConfigContract = "GotifySettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationGotifyResource{}
	_ resource.ResourceWithImportState = &NotificationGotifyResource{}
)

func NewNotificationGotifyResource() resource.Resource {
	return &NotificationGotifyResource{}
}

// NotificationGotifyResource defines the notification implementation.
type NotificationGotifyResource struct {
	client *prowlarr.APIClient
}

// NotificationGotify describes the notification data model.
type NotificationGotify struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Server                types.String `tfsdk:"server"`
	Name                  types.String `tfsdk:"name"`
	AppToken              types.String `tfsdk:"app_token"`
	Priority              types.Int64  `tfsdk:"priority"`
	ID                    types.Int64  `tfsdk:"id"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	IncludeManualGrabs    types.Bool   `tfsdk:"include_manual_grabs"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnHealthRestored      types.Bool   `tfsdk:"on_health_restored"`
}

func (n NotificationGotify) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		Server:                n.Server,
		AppToken:              n.AppToken,
		ItemPriority:          n.Priority,
		Name:                  n.Name,
		ID:                    n.ID,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		IncludeManualGrabs:    n.IncludeManualGrabs,
		OnGrab:                n.OnGrab,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		OnHealthRestored:      n.OnHealthRestored,
		ConfigContract:        types.StringValue(notificationGotifyConfigContract),
		Implementation:        types.StringValue(notificationGotifyImplementation),
	}
}

func (n *NotificationGotify) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Server = notification.Server
	n.AppToken = notification.AppToken
	n.Priority = notification.ItemPriority
	n.Name = notification.Name
	n.ID = notification.ID
	n.IncludeManualGrabs = notification.IncludeManualGrabs
	n.OnGrab = notification.OnGrab
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnHealthRestored = notification.OnHealthRestored
}

func (r *NotificationGotifyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationGotifyResourceName
}

func (r *NotificationGotifyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Gotify resource.\nFor more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect) and [Gotify](https://wiki.servarr.com/prowlarr/supported#gotify).",
		Attributes: map[string]schema.Attribute{
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_health_restored": schema.BoolAttribute{
				MarkdownDescription: "On health restored flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On release grab flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_manual_grabs": schema.BoolAttribute{
				MarkdownDescription: "Include manual grab flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationGotify name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority. `0` Min, `2` Low, `5` Normal, `8` High.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 2, 5, 8),
				},
			},
			"server": schema.StringAttribute{
				MarkdownDescription: "Server.",
				Required:            true,
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "App token.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationGotifyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationGotifyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationGotify

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationGotify
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationGotifyResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationGotifyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationGotifyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationGotify

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationGotify current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationGotifyResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationGotifyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationGotifyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationGotify

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationGotify
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationGotifyResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationGotifyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationGotifyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationGotify current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationGotifyResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationGotifyResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationGotifyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationGotifyResourceName+": "+req.ID)
}

func (n *NotificationGotify) write(ctx context.Context, notification *prowlarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationGotify) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
