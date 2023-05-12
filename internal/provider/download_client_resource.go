package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const downloadClientResourceName = "download_client"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientResource{}
	_ resource.ResourceWithImportState = &DownloadClientResource{}
)

var downloadClientFields = helpers.Fields{
	Bools:                  []string{"addPaused", "useSsl", "startOnAdd", "sequentialOrder", "addStopped", "saveMagnetFiles", "readOnly"},
	Ints:                   []string{"port", "itemPriority", "initialState", "intialState"},
	IntsExceptions:         []string{"priority"},
	Strings:                []string{"host", "apiKey", "urlBase", "rpcPath", "secretToken", "password", "username", "tvImportedCategory", "directory", "destinationDirectory", "destination", "category", "nzbFolder", "strmFolder", "torrentFolder", "magnetFileExtension", "apiUrl", "appId", "appToken"},
	StringSlices:           []string{"fieldTags", "postImTags"},
	StringSlicesExceptions: []string{"tags"},
	IntSlices:              []string{"additionalTags"},
}

func NewDownloadClientResource() resource.Resource {
	return &DownloadClientResource{}
}

// DownloadClientResource defines the download client implementation.
type DownloadClientResource struct {
	client *prowlarr.APIClient
}

// DownloadClient describes the download client data model.
type DownloadClient struct {
	Tags                 types.Set    `tfsdk:"tags"`
	PostImTags           types.Set    `tfsdk:"post_im_tags"`
	FieldTags            types.Set    `tfsdk:"field_tags"`
	AdditionalTags       types.Set    `tfsdk:"additional_tags"`
	Categories           types.Set    `tfsdk:"categories"`
	NzbFolder            types.String `tfsdk:"nzb_folder"`
	Category             types.String `tfsdk:"category"`
	Implementation       types.String `tfsdk:"implementation"`
	Name                 types.String `tfsdk:"name"`
	Protocol             types.String `tfsdk:"protocol"`
	MagnetFileExtension  types.String `tfsdk:"magnet_file_extension"`
	TorrentFolder        types.String `tfsdk:"torrent_folder"`
	StrmFolder           types.String `tfsdk:"strm_folder"`
	Host                 types.String `tfsdk:"host"`
	ConfigContract       types.String `tfsdk:"config_contract"`
	Destination          types.String `tfsdk:"destination"`
	Directory            types.String `tfsdk:"directory"`
	Username             types.String `tfsdk:"username"`
	TvImportedCategory   types.String `tfsdk:"tv_imported_category"`
	Password             types.String `tfsdk:"password"`
	SecretToken          types.String `tfsdk:"secret_token"`
	RPCPath              types.String `tfsdk:"rpc_path"`
	URLBase              types.String `tfsdk:"url_base"`
	APIKey               types.String `tfsdk:"api_key"`
	APIURL               types.String `tfsdk:"api_url"`
	AppID                types.String `tfsdk:"app_id"`
	AppToken             types.String `tfsdk:"app_token"`
	DestinationDirectory types.String `tfsdk:"destination_directory"`
	ItemPriority         types.Int64  `tfsdk:"item_priority"`
	IntialState          types.Int64  `tfsdk:"intial_state"`
	InitialState         types.Int64  `tfsdk:"initial_state"`
	Priority             types.Int64  `tfsdk:"priority"`
	Port                 types.Int64  `tfsdk:"port"`
	ID                   types.Int64  `tfsdk:"id"`
	AddStopped           types.Bool   `tfsdk:"add_stopped"`
	SaveMagnetFiles      types.Bool   `tfsdk:"save_magnet_files"`
	ReadOnly             types.Bool   `tfsdk:"read_only"`
	SequentialOrder      types.Bool   `tfsdk:"sequential_order"`
	StartOnAdd           types.Bool   `tfsdk:"start_on_add"`
	UseSsl               types.Bool   `tfsdk:"use_ssl"`
	AddPaused            types.Bool   `tfsdk:"add_paused"`
	Enable               types.Bool   `tfsdk:"enable"`
}

// ClientCategory is part of DownloadClient.
type ClientCategory struct {
	Categories types.Set    `tfsdk:"categories"`
	Name       types.String `tfsdk:"name"`
}

func (r *DownloadClientResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientResourceName
}

func (r *DownloadClientResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Generic Download Client resource. When possible use a specific resource instead.\nFor more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients).",
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
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "DownloadClient configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "DownloadClient implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("usenet", "torrent"),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"categories": schema.SetNestedAttribute{
				MarkdownDescription: "List of mapped categories.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getClientCategorySchema().Attributes,
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
			"start_on_add": schema.BoolAttribute{
				MarkdownDescription: "Start on add flag.",
				Optional:            true,
				Computed:            true,
			},
			"sequential_order": schema.BoolAttribute{
				MarkdownDescription: "Sequential order flag.",
				Optional:            true,
				Computed:            true,
			},
			"add_stopped": schema.BoolAttribute{
				MarkdownDescription: "Add stopped flag.",
				Optional:            true,
				Computed:            true,
			},
			"save_magnet_files": schema.BoolAttribute{
				MarkdownDescription: "Save magnet files flag.",
				Optional:            true,
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Read only flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"item_priority": schema.Int64Attribute{
				MarkdownDescription: "Priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"initial_state": schema.Int64Attribute{
				MarkdownDescription: "Initial state. `0` Start, `1` ForceStart, `2` Pause.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"intial_state": schema.Int64Attribute{
				MarkdownDescription: "Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.",
				Optional:            true,
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
			},
			"rpc_path": schema.StringAttribute{
				MarkdownDescription: "RPC path.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "API URL.",
				Optional:            true,
				Computed:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "App ID.",
				Optional:            true,
				Computed:            true,
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "App Token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"secret_token": schema.StringAttribute{
				MarkdownDescription: "Secret token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
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
			"tv_imported_category": schema.StringAttribute{
				MarkdownDescription: "TV imported category.",
				Optional:            true,
				Computed:            true,
			},
			"directory": schema.StringAttribute{
				MarkdownDescription: "Directory.",
				Optional:            true,
				Computed:            true,
			},
			"destination_directory": schema.StringAttribute{
				MarkdownDescription: "Movie directory.",
				Optional:            true,
				Computed:            true,
			},
			"destination": schema.StringAttribute{
				MarkdownDescription: "Destination.",
				Optional:            true,
				Computed:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Category.",
				Optional:            true,
				Computed:            true,
			},
			"nzb_folder": schema.StringAttribute{
				MarkdownDescription: "NZB folder.",
				Optional:            true,
				Computed:            true,
			},
			"strm_folder": schema.StringAttribute{
				MarkdownDescription: "STRM folder.",
				Optional:            true,
				Computed:            true,
			},
			"torrent_folder": schema.StringAttribute{
				MarkdownDescription: "Torrent folder.",
				Optional:            true,
				Computed:            true,
			},
			"magnet_file_extension": schema.StringAttribute{
				MarkdownDescription: "Magnet file extension.",
				Optional:            true,
				Computed:            true,
			},
			"additional_tags": schema.SetAttribute{
				MarkdownDescription: "Additional tags, `0` TitleSlug, `1` Quality, `2` Language, `3` ReleaseGroup, `4` Year, `5` Indexer, `6` Network.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"field_tags": schema.SetAttribute{
				MarkdownDescription: "Field tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"post_im_tags": schema.SetAttribute{
				MarkdownDescription: "Post import tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r DownloadClientResource) getClientCategorySchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of client category.",
				Optional:            true,
				Computed:            true,
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "List of categories.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *DownloadClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClient

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClient
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state DownloadClient

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DownloadClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClient

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClient current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state DownloadClient

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DownloadClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClient

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClient
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state DownloadClient

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DownloadClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClient

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClient current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientResourceName+": "+req.ID)
}

func (d *DownloadClient) write(ctx context.Context, downloadClient *prowlarr.DownloadClientResource) {
	d.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, downloadClient.GetTags())
	d.Enable = types.BoolValue(downloadClient.GetEnable())
	d.Priority = types.Int64Value(int64(downloadClient.GetPriority()))
	d.ID = types.Int64Value(int64(downloadClient.GetId()))
	d.ConfigContract = types.StringValue(downloadClient.GetConfigContract())
	d.Implementation = types.StringValue(downloadClient.GetImplementation())
	d.Name = types.StringValue(downloadClient.GetName())
	d.Protocol = types.StringValue(string(downloadClient.GetProtocol()))
	d.AdditionalTags = types.SetValueMust(types.Int64Type, nil)
	d.FieldTags = types.SetValueMust(types.StringType, nil)
	d.PostImTags = types.SetValueMust(types.StringType, nil)
	d.Categories = types.SetValueMust(DownloadClientResource{}.getClientCategorySchema().Type(), nil)

	categories := make([]ClientCategory, len(downloadClient.GetCategories()))
	for i, c := range downloadClient.GetCategories() {
		categories[i].write(ctx, c)
	}

	tfsdk.ValueFrom(ctx, categories, d.Categories.Type(ctx), &d.Categories)
	helpers.WriteFields(ctx, d, downloadClient.GetFields(), downloadClientFields)
}

func (c *ClientCategory) write(ctx context.Context, category *prowlarr.DownloadClientCategory) {
	c.Name = types.StringValue(category.GetClientCategory())
	c.Categories = types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, category.Categories, c.Categories.Type(ctx), &c.Categories)
}

func (d *DownloadClient) read(ctx context.Context) *prowlarr.DownloadClientResource {
	tags := make([]*int32, len(d.Tags.Elements()))
	tfsdk.ValueAs(ctx, d.Tags, &tags)

	categories := make([]*ClientCategory, len(d.Categories.Elements()))
	tfsdk.ValueAs(ctx, d.Categories, &categories)

	clientCategories := make([]*prowlarr.DownloadClientCategory, len(d.Categories.Elements()))
	for n, c := range categories {
		clientCategories[n] = c.read(ctx)
	}

	client := prowlarr.NewDownloadClientResource()
	client.SetEnable(d.Enable.ValueBool())
	client.SetPriority(int32(d.Priority.ValueInt64()))
	client.SetId(int32(d.ID.ValueInt64()))
	client.SetConfigContract(d.ConfigContract.ValueString())
	client.SetImplementation(d.Implementation.ValueString())
	client.SetName(d.Name.ValueString())
	client.SetProtocol(prowlarr.DownloadProtocol(d.Protocol.ValueString()))
	client.SetTags(tags)
	client.SetFields(helpers.ReadFields(ctx, d, downloadClientFields))
	client.SetCategories(clientCategories)

	return client
}

func (c *ClientCategory) read(ctx context.Context) *prowlarr.DownloadClientCategory {
	categories := make([]*int32, len(c.Categories.Elements()))
	tfsdk.ValueAs(ctx, c.Categories, &categories)

	category := prowlarr.NewDownloadClientCategory()
	category.SetCategories(categories)
	category.SetClientCategory(c.Name.ValueString())

	return category
}
