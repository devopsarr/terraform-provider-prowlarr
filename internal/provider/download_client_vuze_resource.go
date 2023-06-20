package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	downloadClientVuzeResourceName   = "download_client_vuze"
	downloadClientVuzeImplementation = "Vuze"
	downloadClientVuzeConfigContract = "TransmissionSettings"
	downloadClientVuzeProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientVuzeResource{}
	_ resource.ResourceWithImportState = &DownloadClientVuzeResource{}
)

func NewDownloadClientVuzeResource() resource.Resource {
	return &DownloadClientVuzeResource{}
}

// DownloadClientVuzeResource defines the download client implementation.
type DownloadClientVuzeResource struct {
	client *prowlarr.APIClient
}

// DownloadClientVuze describes the download client data model.
type DownloadClientVuze struct {
	Tags         types.Set    `tfsdk:"tags"`
	Categories   types.Set    `tfsdk:"categories"`
	Name         types.String `tfsdk:"name"`
	Host         types.String `tfsdk:"host"`
	URLBase      types.String `tfsdk:"url_base"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	Category     types.String `tfsdk:"category"`
	Directory    types.String `tfsdk:"directory"`
	ItemPriority types.Int64  `tfsdk:"item_priority"`
	Priority     types.Int64  `tfsdk:"priority"`
	Port         types.Int64  `tfsdk:"port"`
	ID           types.Int64  `tfsdk:"id"`
	AddPaused    types.Bool   `tfsdk:"add_paused"`
	UseSsl       types.Bool   `tfsdk:"use_ssl"`
	Enable       types.Bool   `tfsdk:"enable"`
}

func (d DownloadClientVuze) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:           d.Tags,
		Categories:     d.Categories,
		Name:           d.Name,
		Host:           d.Host,
		URLBase:        d.URLBase,
		Username:       d.Username,
		Password:       d.Password,
		Category:       d.Category,
		Directory:      d.Directory,
		ItemPriority:   d.ItemPriority,
		Priority:       d.Priority,
		Port:           d.Port,
		ID:             d.ID,
		AddPaused:      d.AddPaused,
		UseSsl:         d.UseSsl,
		Enable:         d.Enable,
		Implementation: types.StringValue(downloadClientVuzeImplementation),
		ConfigContract: types.StringValue(downloadClientVuzeConfigContract),
		Protocol:       types.StringValue(downloadClientVuzeProtocol),
	}
}

func (d *DownloadClientVuze) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Categories = client.Categories
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.Username = client.Username
	d.Password = client.Password
	d.Category = client.Category
	d.Directory = client.Directory
	d.ItemPriority = client.ItemPriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.AddPaused = client.AddPaused
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
}

func (r *DownloadClientVuzeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientVuzeResourceName
}

func (r *DownloadClientVuzeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Vuze resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients) and [Vuze](https://wiki.servarr.com/prowlarr/supported#vuze).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"categories": schema.SetNestedAttribute{
				MarkdownDescription: "List of mapped categories.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: DownloadClientResource{}.getClientCategorySchema().Attributes,
				},
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"add_paused": schema.BoolAttribute{
				MarkdownDescription: "Add paused flag.",
				Optional:            true,
				Computed:            true,
			},
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"item_priority": schema.Int64Attribute{
				MarkdownDescription: "Older Movie priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
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
			"category": schema.StringAttribute{
				MarkdownDescription: "Category.",
				Optional:            true,
				Computed:            true,
			},
			"directory": schema.StringAttribute{
				MarkdownDescription: "Directory.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientVuzeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientVuzeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientVuze

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientVuze
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientVuzeResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientVuzeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientVuzeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientVuze

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientVuze current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientVuzeResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientVuzeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientVuzeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientVuze

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientVuze
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientVuzeResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientVuzeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientVuzeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientVuze

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientVuze current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientVuzeResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientVuzeResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientVuzeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientVuzeResourceName+": "+req.ID)
}

func (d *DownloadClientVuze) write(ctx context.Context, downloadClient *prowlarr.DownloadClientResource) {
	genericDownloadClient := DownloadClient{}
	genericDownloadClient.write(ctx, downloadClient)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientVuze) read(ctx context.Context) *prowlarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx)
}
