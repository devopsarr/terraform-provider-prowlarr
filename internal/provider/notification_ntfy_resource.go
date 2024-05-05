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
	notificationNtfyResourceName   = "notification_ntfy"
	notificationNtfyImplementation = "Ntfy"
	notificationNtfyConfigContract = "NtfySettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationNtfyResource{}
	_ resource.ResourceWithImportState = &NotificationNtfyResource{}
)

func NewNotificationNtfyResource() resource.Resource {
	return &NotificationNtfyResource{}
}

// NotificationNtfyResource defines the notification implementation.
type NotificationNtfyResource struct {
	client *prowlarr.APIClient
}

// NotificationNtfy describes the notification data model.
type NotificationNtfy struct {
	Tags                  types.Set    `tfsdk:"tags"`
	FieldTags             types.Set    `tfsdk:"field_tags"`
	Topics                types.Set    `tfsdk:"topics"`
	ClickURL              types.String `tfsdk:"click_url"`
	ServerURL             types.String `tfsdk:"server_url"`
	Username              types.String `tfsdk:"username"`
	Name                  types.String `tfsdk:"name"`
	Password              types.String `tfsdk:"password"`
	AccessToken           types.String `tfsdk:"access_token"`
	Priority              types.Int64  `tfsdk:"priority"`
	ID                    types.Int64  `tfsdk:"id"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	IncludeManualGrabs    types.Bool   `tfsdk:"include_manual_grabs"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnHealthRestored      types.Bool   `tfsdk:"on_health_restored"`
}

func (n NotificationNtfy) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		FieldTags:             n.FieldTags,
		Topics:                n.Topics,
		ServerURL:             n.ServerURL,
		ClickURL:              n.ClickURL,
		Username:              n.Username,
		Password:              n.Password,
		AccessToken:           n.AccessToken,
		Name:                  n.Name,
		ItemPriority:          n.Priority,
		ID:                    n.ID,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		IncludeManualGrabs:    n.IncludeManualGrabs,
		OnGrab:                n.OnGrab,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		OnHealthRestored:      n.OnHealthRestored,
		ConfigContract:        types.StringValue(notificationNtfyConfigContract),
		Implementation:        types.StringValue(notificationNtfyImplementation),
	}
}

func (n *NotificationNtfy) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.FieldTags = notification.FieldTags
	n.Topics = notification.Topics
	n.ServerURL = notification.ServerURL
	n.ClickURL = notification.ClickURL
	n.Username = notification.Username
	n.Password = notification.Password
	n.AccessToken = notification.AccessToken
	n.Name = notification.Name
	n.Priority = notification.ItemPriority
	n.ID = notification.ID
	n.IncludeManualGrabs = notification.IncludeManualGrabs
	n.OnGrab = notification.OnGrab
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnHealthRestored = notification.OnHealthRestored
}

func (r *NotificationNtfyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationNtfyResourceName
}

func (r *NotificationNtfyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->\nNotification Ntfy resource.\nFor more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect) and [Ntfy](https://wiki.servarr.com/prowlarr/supported#ntfy).",
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
				MarkdownDescription: "NotificationNtfy name.",
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
				MarkdownDescription: "Priority. `1` Min, `2` Low, `3` Default, `4` High, `5` Max.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2, 3, 4, 5),
				},
			},
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Server URL.",
				Optional:            true,
				Computed:            true,
			},
			"click_url": schema.StringAttribute{
				MarkdownDescription: "Click URL.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"topics": schema.SetAttribute{
				MarkdownDescription: "Topics.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"field_tags": schema.SetAttribute{
				MarkdownDescription: "Tags and emojis.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationNtfyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationNtfyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationNtfy

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationNtfy
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationNtfyResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationNtfyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationNtfyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationNtfy

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationNtfy current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationNtfyResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationNtfyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationNtfyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationNtfy

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationNtfy
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationNtfyResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationNtfyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationNtfyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationNtfy current value
	_, err := r.client.NotificationAPI.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationNtfyResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationNtfyResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationNtfyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationNtfyResourceName+": "+req.ID)
}

func (n *NotificationNtfy) write(ctx context.Context, notification *prowlarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationNtfy) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
