package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	downloadClientUsenetDownloadStationResourceName   = "download_client_usenet_download_station"
	downloadClientUsenetDownloadStationImplementation = "UsenetDownloadStation"
	downloadClientUsenetDownloadStationConfigContract = "DownloadStationSettings"
	downloadClientUsenetDownloadStationProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientUsenetDownloadStationResource{}
	_ resource.ResourceWithImportState = &DownloadClientUsenetDownloadStationResource{}
)

func NewDownloadClientUsenetDownloadStationResource() resource.Resource {
	return &DownloadClientUsenetDownloadStationResource{}
}

// DownloadClientUsenetDownloadStationResource defines the download client implementation.
type DownloadClientUsenetDownloadStationResource struct {
	client *prowlarr.APIClient
}

// DownloadClientUsenetDownloadStation describes the download client data model.
type DownloadClientUsenetDownloadStation struct {
	Tags        types.Set    `tfsdk:"tags"`
	Categories  types.Set    `tfsdk:"categories"`
	Name        types.String `tfsdk:"name"`
	Host        types.String `tfsdk:"host"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Category    types.String `tfsdk:"category"`
	TVDirectory types.String `tfsdk:"station_directory"`
	Priority    types.Int64  `tfsdk:"priority"`
	Port        types.Int64  `tfsdk:"port"`
	ID          types.Int64  `tfsdk:"id"`
	UseSsl      types.Bool   `tfsdk:"use_ssl"`
	Enable      types.Bool   `tfsdk:"enable"`
}

func (d DownloadClientUsenetDownloadStation) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:           d.Tags,
		Categories:     d.Categories,
		Name:           d.Name,
		Host:           d.Host,
		Username:       d.Username,
		Password:       d.Password,
		Category:       d.Category,
		TVDirectory:    d.TVDirectory,
		Priority:       d.Priority,
		Port:           d.Port,
		ID:             d.ID,
		UseSsl:         d.UseSsl,
		Enable:         d.Enable,
		Implementation: types.StringValue(downloadClientUsenetDownloadStationImplementation),
		ConfigContract: types.StringValue(downloadClientUsenetDownloadStationConfigContract),
		Protocol:       types.StringValue(downloadClientUsenetDownloadStationProtocol),
	}
}

func (d *DownloadClientUsenetDownloadStation) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Categories = client.Categories
	d.Name = client.Name
	d.Host = client.Host
	d.Username = client.Username
	d.Password = client.Password
	d.Category = client.Category
	d.TVDirectory = client.TVDirectory
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
}

func (r *DownloadClientUsenetDownloadStationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientUsenetDownloadStationResourceName
}

func (r *DownloadClientUsenetDownloadStationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->\nDownload Client UsenetDownloadStation resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients) and [UsenetDownloadStation](https://wiki.servarr.com/prowlarr/supported#usenetdownloadstation).",
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
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
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
			"station_directory": schema.StringAttribute{
				MarkdownDescription: "Directory.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientUsenetDownloadStationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientUsenetDownloadStationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientUsenetDownloadStation

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientUsenetDownloadStation
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientUsenetDownloadStationResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientUsenetDownloadStationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientUsenetDownloadStationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientUsenetDownloadStation

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientUsenetDownloadStation current value
	response, _, err := r.client.DownloadClientAPI.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientUsenetDownloadStationResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientUsenetDownloadStationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientUsenetDownloadStationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientUsenetDownloadStation

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientUsenetDownloadStation
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientUsenetDownloadStationResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientUsenetDownloadStationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientUsenetDownloadStationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientUsenetDownloadStation current value
	_, err := r.client.DownloadClientAPI.DeleteDownloadClient(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientUsenetDownloadStationResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientUsenetDownloadStationResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientUsenetDownloadStationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientUsenetDownloadStationResourceName+": "+req.ID)
}

func (d *DownloadClientUsenetDownloadStation) write(ctx context.Context, downloadClient *prowlarr.DownloadClientResource, diags *diag.Diagnostics) {
	genericDownloadClient := d.toDownloadClient()
	genericDownloadClient.write(ctx, downloadClient, diags)
	d.fromDownloadClient(genericDownloadClient)
}

func (d *DownloadClientUsenetDownloadStation) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx, diags)
}
