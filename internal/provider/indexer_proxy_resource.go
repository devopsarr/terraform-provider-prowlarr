package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerProxyResourceName = "indexer_proxy"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerProxyResource{}
	_ resource.ResourceWithImportState = &IndexerProxyResource{}
)

var indexerProxyFields = helpers.Fields{
	Ints:    []string{"port", "requestTimeout"},
	Strings: []string{"host", "password", "username"},
}

func NewIndexerProxyResource() resource.Resource {
	return &IndexerProxyResource{}
}

// IndexerProxyResource defines the indexer proxy implementation.
type IndexerProxyResource struct {
	client *prowlarr.APIClient
	auth   context.Context
}

// IndexerProxy describes the indexer proxy data model.
type IndexerProxy struct {
	Tags           types.Set    `tfsdk:"tags"`
	Name           types.String `tfsdk:"name"`
	ConfigContract types.String `tfsdk:"config_contract"`
	Implementation types.String `tfsdk:"implementation"`
	Host           types.String `tfsdk:"host"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Port           types.Int64  `tfsdk:"port"`
	RequestTimeout types.Int64  `tfsdk:"request_timeout"`
	ID             types.Int64  `tfsdk:"id"`
}

func (i IndexerProxy) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"tags":            types.SetType{}.WithElementType(types.Int64Type),
			"name":            types.StringType,
			"config_contract": types.StringType,
			"implementation":  types.StringType,
			"host":            types.StringType,
			"username":        types.StringType,
			"password":        types.StringType,
			"port":            types.Int64Type,
			"request_timeout": types.Int64Type,
			"id":              types.Int64Type,
		})
}

func (r *IndexerProxyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxyResourceName
}

func (r *IndexerProxyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->\nGeneric Indexer Proxy resource. When possible use a specific resource instead.\nFor more information refer to [Indexer Proxy](https://wiki.servarr.com/prowlarr/settings#indexer-proxies).",
		Attributes: map[string]schema.Attribute{
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "IndexerProxy configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "IndexerProxy implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Indexer Proxy name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Indexer Proxy ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"request_timeout": schema.Int64Attribute{
				MarkdownDescription: "Request timeout.",
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
		},
	}
}

func (r *IndexerProxyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *IndexerProxyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var proxy *IndexerProxy

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerProxy
	request := proxy.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerProxyAPI.CreateIndexerProxy(r.auth).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerProxyResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerProxyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state IndexerProxy

	state.writeSensitive(proxy)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerProxyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var proxy *IndexerProxy

	resp.Diagnostics.Append(req.State.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerProxy current value
	response, _, err := r.client.IndexerProxyAPI.GetIndexerProxyById(r.auth, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxyResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state IndexerProxy

	state.writeSensitive(proxy)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerProxyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var proxy *IndexerProxy

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerProxy
	request := proxy.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerProxyAPI.UpdateIndexerProxy(r.auth, strconv.Itoa(int(request.GetId()))).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerProxyResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerProxyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state IndexerProxy

	state.writeSensitive(proxy)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerProxyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerProxy current value
	_, err := r.client.IndexerProxyAPI.DeleteIndexerProxy(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerProxyResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerProxyResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerProxyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerProxyResourceName+": "+req.ID)
}

func (i *IndexerProxy) write(ctx context.Context, indexerProxy *prowlarr.IndexerProxyResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	i.ID = types.Int64Value(int64(indexerProxy.GetId()))
	i.ConfigContract = types.StringValue(indexerProxy.GetConfigContract())
	i.Implementation = types.StringValue(indexerProxy.GetImplementation())
	i.Name = types.StringValue(indexerProxy.GetName())
	i.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, indexerProxy.Tags)
	diags.Append(localDiag...)
	helpers.WriteFields(ctx, i, indexerProxy.GetFields(), indexerProxyFields)
}

func (i *IndexerProxy) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.IndexerProxyResource {
	proxy := prowlarr.NewIndexerProxyResource()
	proxy.SetId(int32(i.ID.ValueInt64()))
	proxy.SetConfigContract(i.ConfigContract.ValueString())
	proxy.SetImplementation(i.Implementation.ValueString())
	proxy.SetName(i.Name.ValueString())
	diags.Append(i.Tags.ElementsAs(ctx, &proxy.Tags, true)...)
	proxy.SetFields(helpers.ReadFields(ctx, i, indexerProxyFields))

	return proxy
}

// writeSensitive copy sensitive data from another resource.
func (i *IndexerProxy) writeSensitive(proxy *IndexerProxy) {
	if !proxy.Password.IsUnknown() {
		i.Password = proxy.Password
	}
}
