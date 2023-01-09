package provider

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"

	"github.com/devopsarr/terraform-provider-prowlarr/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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

const (
	downloadClientTransmissionResourceName   = "download_client_transmission"
	downloadClientTransmissionImplementation = "Transmission"
	downloadClientTransmissionConfigContract = "TransmissionSettings"
	downloadClientTransmissionProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientTransmissionResource{}
	_ resource.ResourceWithImportState = &DownloadClientTransmissionResource{}
)

func NewDownloadClientTransmissionResource() resource.Resource {
	return &DownloadClientTransmissionResource{}
}

// DownloadClientTransmissionResource defines the download client implementation.
type DownloadClientTransmissionResource struct {
	client *prowlarr.APIClient
}

// DownloadClientTransmission describes the download client data model.
type DownloadClientTransmission struct {
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

func (d DownloadClientTransmission) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:         d.Tags,
		Categories:   d.Categories,
		Name:         d.Name,
		Host:         d.Host,
		URLBase:      d.URLBase,
		Username:     d.Username,
		Password:     d.Password,
		Category:     d.Category,
		Directory:    d.Directory,
		ItemPriority: d.ItemPriority,
		Priority:     d.Priority,
		Port:         d.Port,
		ID:           d.ID,
		AddPaused:    d.AddPaused,
		UseSsl:       d.UseSsl,
		Enable:       d.Enable,
	}
}

func (d *DownloadClientTransmission) fromDownloadClient(client *DownloadClient) {
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

func (r *DownloadClientTransmissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientTransmissionResourceName
}

func (r *DownloadClientTransmissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Transmission resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients) and [Transmission](https://wiki.servarr.com/prowlarr/supported#transmission).",
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
				Optional:            true,
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
				MarkdownDescription: "Priority. `0` Last, `1` First.",
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

func (r *DownloadClientTransmissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*prowlarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *prowlarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DownloadClientTransmissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientTransmission

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientTransmission
	request := client.read(ctx)

	response, aaa, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		test, _ := io.ReadAll(aaa.Body)
		resp.Diagnostics.AddError(tools.ClientError, string(test))
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", downloadClientTransmissionResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientTransmissionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientTransmissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientTransmission

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientTransmission current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientTransmissionResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientTransmissionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientTransmissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientTransmission

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientTransmission
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", downloadClientTransmissionResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientTransmissionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientTransmissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientTransmission

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientTransmission current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientTransmissionResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientTransmissionResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientTransmissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+downloadClientTransmissionResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (d *DownloadClientTransmission) write(ctx context.Context, downloadClient *prowlarr.DownloadClientResource) {
	genericDownloadClient := DownloadClient{
		Enable:     types.BoolValue(downloadClient.GetEnable()),
		Priority:   types.Int64Value(int64(downloadClient.GetPriority())),
		ID:         types.Int64Value(int64(downloadClient.GetId())),
		Name:       types.StringValue(downloadClient.GetName()),
		Tags:       types.SetValueMust(types.Int64Type, nil),
		Categories: types.SetValueMust(DownloadClientResource{}.getClientCategorySchema().Type(), nil),
	}

	tfsdk.ValueFrom(ctx, downloadClient.Tags, genericDownloadClient.Tags.Type(ctx), &genericDownloadClient.Tags)
	genericDownloadClient.writeFields(ctx, downloadClient.Fields)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientTransmission) read(ctx context.Context) *prowlarr.DownloadClientResource {
	tags := make([]*int32, len(d.Tags.Elements()))
	categories := make([]*prowlarr.DownloadClientCategory, 0)

	tfsdk.ValueAs(ctx, d.Categories, &categories)
	tfsdk.ValueAs(ctx, d.Tags, &tags)

	client := prowlarr.NewDownloadClientResource()
	client.SetEnable(d.Enable.ValueBool())
	client.SetPriority(int32(d.Priority.ValueInt64()))
	client.SetId(int32(d.ID.ValueInt64()))
	client.SetConfigContract(downloadClientTransmissionConfigContract)
	client.SetImplementation(downloadClientTransmissionImplementation)
	client.SetName(d.Name.ValueString())
	client.SetProtocol(downloadClientTransmissionProtocol)
	client.SetTags(tags)
	client.SetCategories(categories)
	client.SetFields(d.toDownloadClient().readFields(ctx))

	return client
}
