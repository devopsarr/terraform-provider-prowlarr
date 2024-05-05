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
	notificationJoinResourceName   = "notification_join"
	notificationJoinImplementation = "Join"
	notificationJoinConfigContract = "JoinSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationJoinResource{}
	_ resource.ResourceWithImportState = &NotificationJoinResource{}
)

func NewNotificationJoinResource() resource.Resource {
	return &NotificationJoinResource{}
}

// NotificationJoinResource defines the notification implementation.
type NotificationJoinResource struct {
	client *prowlarr.APIClient
}

// NotificationJoin describes the notification data model.
type NotificationJoin struct {
	Tags                  types.Set    `tfsdk:"tags"`
	DeviceNames           types.String `tfsdk:"device_names"`
	Name                  types.String `tfsdk:"name"`
	APIKey                types.String `tfsdk:"api_key"`
	Priority              types.Int64  `tfsdk:"priority"`
	ID                    types.Int64  `tfsdk:"id"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	IncludeManualGrabs    types.Bool   `tfsdk:"include_manual_grabs"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnHealthRestored      types.Bool   `tfsdk:"on_health_restored"`
}

func (n NotificationJoin) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		DeviceNames:           n.DeviceNames,
		APIKey:                n.APIKey,
		ItemPriority:          n.Priority,
		Name:                  n.Name,
		ID:                    n.ID,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		IncludeManualGrabs:    n.IncludeManualGrabs,
		OnGrab:                n.OnGrab,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		OnHealthRestored:      n.OnHealthRestored,
		ConfigContract:        types.StringValue(notificationJoinConfigContract),
		Implementation:        types.StringValue(notificationJoinImplementation),
	}
}

func (n *NotificationJoin) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.DeviceNames = notification.DeviceNames
	n.APIKey = notification.APIKey
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

func (r *NotificationJoinResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationJoinResourceName
}

func (r *NotificationJoinResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->\nNotification Join resource.\nFor more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect) and [Join](https://wiki.servarr.com/prowlarr/supported#join).",
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
				MarkdownDescription: "NotificationJoin name.",
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
				MarkdownDescription: "Priority. `-2` Silent, `-1` Quiet, `0` Normal, `1` High, `2` Emergency.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-2, -1, 0, 1, 2),
				},
			},
			"device_names": schema.StringAttribute{
				MarkdownDescription: "Device names. Comma separated list.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationJoinResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationJoinResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationJoin

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationJoin
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationJoinResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationJoinResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationJoinResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationJoin

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationJoin current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationJoinResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationJoinResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationJoinResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationJoin

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationJoin
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationJoinResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationJoinResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationJoinResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationJoin current value
	_, err := r.client.NotificationAPI.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationJoinResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationJoinResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationJoinResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationJoinResourceName+": "+req.ID)
}

func (n *NotificationJoin) write(ctx context.Context, notification *prowlarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationJoin) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
