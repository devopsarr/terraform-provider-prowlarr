package provider

import (
	"context"
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

const (
	notificationTelegramResourceName   = "notification_telegram"
	notificationTelegramImplementation = "Telegram"
	notificationTelegramConfigContract = "TelegramSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationTelegramResource{}
	_ resource.ResourceWithImportState = &NotificationTelegramResource{}
)

func NewNotificationTelegramResource() resource.Resource {
	return &NotificationTelegramResource{}
}

// NotificationTelegramResource defines the notification implementation.
type NotificationTelegramResource struct {
	client *prowlarr.APIClient
}

// NotificationTelegram describes the notification data model.
type NotificationTelegram struct {
	Tags                  types.Set    `tfsdk:"tags"`
	ChatID                types.String `tfsdk:"chat_id"`
	TopicID               types.String `tfsdk:"topic_id"`
	Name                  types.String `tfsdk:"name"`
	BotToken              types.String `tfsdk:"bot_token"`
	ID                    types.Int64  `tfsdk:"id"`
	SendSilently          types.Bool   `tfsdk:"send_silently"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	IncludeManualGrabs    types.Bool   `tfsdk:"include_manual_grabs"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnHealthRestored      types.Bool   `tfsdk:"on_health_restored"`
}

func (n NotificationTelegram) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		ChatID:                n.ChatID,
		TopicID:               n.TopicID,
		BotToken:              n.BotToken,
		SendSilently:          n.SendSilently,
		Name:                  n.Name,
		ID:                    n.ID,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		IncludeManualGrabs:    n.IncludeManualGrabs,
		OnGrab:                n.OnGrab,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		OnHealthRestored:      n.OnHealthRestored,
		ConfigContract:        types.StringValue(notificationTelegramConfigContract),
		Implementation:        types.StringValue(notificationTelegramImplementation),
	}
}

func (n *NotificationTelegram) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.ChatID = notification.ChatID
	n.TopicID = notification.TopicID
	n.BotToken = notification.BotToken
	n.SendSilently = notification.SendSilently
	n.Name = notification.Name
	n.ID = notification.ID
	n.IncludeManualGrabs = notification.IncludeManualGrabs
	n.OnGrab = notification.OnGrab
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnHealthRestored = notification.OnHealthRestored
}

func (r *NotificationTelegramResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationTelegramResourceName
}

func (r *NotificationTelegramResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Telegram resource.\nFor more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect) and [Telegram](https://wiki.servarr.com/prowlarr/supported#telegram).",
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
				MarkdownDescription: "NotificationTelegram name.",
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
			"send_silently": schema.BoolAttribute{
				MarkdownDescription: "Send silently flag.",
				Optional:            true,
				Computed:            true,
			},
			"chat_id": schema.StringAttribute{
				MarkdownDescription: "Chat ID.",
				Required:            true,
			},
			"topic_id": schema.StringAttribute{
				MarkdownDescription: "Topic ID.",
				Optional:            true,
				Computed:            true,
			},
			"bot_token": schema.StringAttribute{
				MarkdownDescription: "Bot token.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationTelegramResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationTelegramResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationTelegram
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationTelegramResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTelegramResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationTelegram current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationTelegramResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTelegramResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationTelegram
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationTelegramResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTelegramResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationTelegram current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationTelegramResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationTelegramResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationTelegramResourceName+": "+req.ID)
}

func (n *NotificationTelegram) write(ctx context.Context, notification *prowlarr.NotificationResource) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification)
	n.fromNotification(genericNotification)
}

func (n *NotificationTelegram) read(ctx context.Context) *prowlarr.NotificationResource {
	return n.toNotification().read(ctx)
}
