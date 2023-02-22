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
	notificationDiscordResourceName   = "notification_discord"
	notificationDiscordImplementation = "Discord"
	notificationDiscordConfigContract = "DiscordSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationDiscordResource{}
	_ resource.ResourceWithImportState = &NotificationDiscordResource{}
)

func NewNotificationDiscordResource() resource.Resource {
	return &NotificationDiscordResource{}
}

// NotificationDiscordResource defines the notification implementation.
type NotificationDiscordResource struct {
	client *prowlarr.APIClient
}

// NotificationDiscord describes the notification data model.
type NotificationDiscord struct {
	Tags                  types.Set    `tfsdk:"tags"`
	ImportFields          types.Set    `tfsdk:"import_fields"`
	GrabFields            types.Set    `tfsdk:"grab_fields"`
	WebHookURL            types.String `tfsdk:"web_hook_url"`
	Name                  types.String `tfsdk:"name"`
	Username              types.String `tfsdk:"username"`
	Avatar                types.String `tfsdk:"avatar"`
	Author                types.String `tfsdk:"author"`
	ID                    types.Int64  `tfsdk:"id"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
}

func (n NotificationDiscord) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		ImportFields:          n.ImportFields,
		GrabFields:            n.GrabFields,
		WebHookURL:            n.WebHookURL,
		Avatar:                n.Avatar,
		Username:              n.Username,
		Author:                n.Author,
		Name:                  n.Name,
		ID:                    n.ID,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		ConfigContract:        types.StringValue(notificationDiscordConfigContract),
		Implementation:        types.StringValue(notificationDiscordImplementation),
	}
}

func (n *NotificationDiscord) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.GrabFields = notification.GrabFields
	n.ImportFields = notification.ImportFields
	n.WebHookURL = notification.WebHookURL
	n.Avatar = notification.Avatar
	n.Username = notification.Username
	n.Author = notification.Author
	n.Name = notification.Name
	n.ID = notification.ID
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
}

func (r *NotificationDiscordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationDiscordResourceName
}

func (r *NotificationDiscordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Discord resource.\nFor more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect) and [Discord](https://wiki.servarr.com/prowlarr/supported#discord).",
		Attributes: map[string]schema.Attribute{
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationDiscord name.",
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
			"web_hook_url": schema.StringAttribute{
				MarkdownDescription: "Web hook URL.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"avatar": schema.StringAttribute{
				MarkdownDescription: "Avatar.",
				Optional:            true,
				Computed:            true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "Author.",
				Optional:            true,
				Computed:            true,
			},
			"grab_fields": schema.SetAttribute{
				MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"import_fields": schema.SetAttribute{
				MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *NotificationDiscordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationDiscordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationDiscord

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationDiscord
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationDiscordResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationDiscordResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationDiscordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationDiscord

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationDiscord current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationDiscordResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationDiscordResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationDiscordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationDiscord

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationDiscord
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationDiscordResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationDiscordResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationDiscordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationDiscord

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationDiscord current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationDiscordResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationDiscordResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationDiscordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationDiscordResourceName+": "+req.ID)
}

func (n *NotificationDiscord) write(ctx context.Context, notification *prowlarr.NotificationResource) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification)
	n.fromNotification(genericNotification)
}

func (n *NotificationDiscord) read(ctx context.Context) *prowlarr.NotificationResource {
	return n.toNotification().read(ctx)
}
