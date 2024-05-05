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
	downloadClientTorrentBlackholeResourceName   = "download_client_torrent_blackhole"
	downloadClientTorrentBlackholeImplementation = "TorrentBlackhole"
	downloadClientTorrentBlackholeConfigContract = "TorrentBlackholeSettings"
	downloadClientTorrentBlackholeProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientTorrentBlackholeResource{}
	_ resource.ResourceWithImportState = &DownloadClientTorrentBlackholeResource{}
)

func NewDownloadClientTorrentBlackholeResource() resource.Resource {
	return &DownloadClientTorrentBlackholeResource{}
}

// DownloadClientTorrentBlackholeResource defines the download client implementation.
type DownloadClientTorrentBlackholeResource struct {
	client *prowlarr.APIClient
}

// DownloadClientTorrentBlackhole describes the download client data model.
type DownloadClientTorrentBlackhole struct {
	Tags                types.Set    `tfsdk:"tags"`
	Categories          types.Set    `tfsdk:"categories"`
	Name                types.String `tfsdk:"name"`
	TorrentFolder       types.String `tfsdk:"torrent_folder"`
	MagnetFileExtension types.String `tfsdk:"magnet_file_extension"`
	Priority            types.Int64  `tfsdk:"priority"`
	ID                  types.Int64  `tfsdk:"id"`
	Enable              types.Bool   `tfsdk:"enable"`
	SaveMagnetFiles     types.Bool   `tfsdk:"save_magnet_files"`
}

func (d DownloadClientTorrentBlackhole) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                d.Tags,
		Categories:          d.Categories,
		Name:                d.Name,
		TorrentFolder:       d.TorrentFolder,
		MagnetFileExtension: d.MagnetFileExtension,
		Priority:            d.Priority,
		ID:                  d.ID,
		Enable:              d.Enable,
		SaveMagnetFiles:     d.SaveMagnetFiles,
		Implementation:      types.StringValue(downloadClientTorrentBlackholeImplementation),
		ConfigContract:      types.StringValue(downloadClientTorrentBlackholeConfigContract),
		Protocol:            types.StringValue(downloadClientTorrentBlackholeProtocol),
	}
}

func (d *DownloadClientTorrentBlackhole) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Categories = client.Categories
	d.Name = client.Name
	d.TorrentFolder = client.TorrentFolder
	d.MagnetFileExtension = client.MagnetFileExtension
	d.Priority = client.Priority
	d.ID = client.ID
	d.Enable = client.Enable
	d.SaveMagnetFiles = client.SaveMagnetFiles
}

func (r *DownloadClientTorrentBlackholeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientTorrentBlackholeResourceName
}

func (r *DownloadClientTorrentBlackholeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->\nDownload Client Torrent Blackhole resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients) and [TorrentBlackhole](https://wiki.servarr.com/prowlarr/supported#torrentblackhole).",
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
			"save_magnet_files": schema.BoolAttribute{
				MarkdownDescription: "Save magnet files flag.",
				Optional:            true,
				Computed:            true,
			},
			"torrent_folder": schema.StringAttribute{
				MarkdownDescription: "Torrent folder.",
				Required:            true,
			},
			"magnet_file_extension": schema.StringAttribute{
				MarkdownDescription: "Magnet file extension.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientTorrentBlackholeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientTorrentBlackholeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientTorrentBlackhole

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientTorrentBlackhole
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientTorrentBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientTorrentBlackholeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientTorrentBlackholeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientTorrentBlackhole

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientTorrentBlackhole current value
	response, _, err := r.client.DownloadClientAPI.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientTorrentBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientTorrentBlackholeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientTorrentBlackholeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientTorrentBlackhole

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientTorrentBlackhole
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientTorrentBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientTorrentBlackholeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientTorrentBlackholeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientTorrentBlackhole current value
	_, err := r.client.DownloadClientAPI.DeleteDownloadClient(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientTorrentBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientTorrentBlackholeResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientTorrentBlackholeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientTorrentBlackholeResourceName+": "+req.ID)
}

func (d *DownloadClientTorrentBlackhole) write(ctx context.Context, downloadClient *prowlarr.DownloadClientResource, diags *diag.Diagnostics) {
	genericDownloadClient := d.toDownloadClient()
	genericDownloadClient.write(ctx, downloadClient, diags)
	d.fromDownloadClient(genericDownloadClient)
}

func (d *DownloadClientTorrentBlackhole) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx, diags)
}
