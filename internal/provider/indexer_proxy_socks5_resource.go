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
	indexerProxySocks5ResourceName   = "indexer_proxy_socks5"
	indexerProxySocks5Implementation = "Socks5"
	indexerProxySocks5ConfigContract = "Socks5Settings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerProxySocks5Resource{}
	_ resource.ResourceWithImportState = &IndexerProxySocks5Resource{}
)

func NewIndexerProxySocks5Resource() resource.Resource {
	return &IndexerProxySocks5Resource{}
}

// IndexerProxySocks5Resource defines the indexer proxy implementation.
type IndexerProxySocks5Resource struct {
	client *prowlarr.APIClient
	auth   context.Context
}

// IndexerProxySocks5 describes the indexer proxy data model.
type IndexerProxySocks5 struct {
	Tags     types.Set    `tfsdk:"tags"`
	Name     types.String `tfsdk:"name"`
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Port     types.Int64  `tfsdk:"port"`
	ID       types.Int64  `tfsdk:"id"`
}

func (i IndexerProxySocks5) toIndexerProxy() *IndexerProxy {
	return &IndexerProxy{
		Tags:           i.Tags,
		Name:           i.Name,
		Host:           i.Host,
		Username:       i.Username,
		Password:       i.Password,
		Port:           i.Port,
		ID:             i.ID,
		ConfigContract: types.StringValue(indexerProxySocks5ConfigContract),
		Implementation: types.StringValue(indexerProxySocks5Implementation),
	}
}

func (i *IndexerProxySocks5) fromIndexerProxy(proxy *IndexerProxy) {
	i.Tags = proxy.Tags
	i.Name = proxy.Name
	i.Host = proxy.Host
	i.Username = proxy.Username
	i.Password = proxy.Password
	i.Port = proxy.Port
	i.ID = proxy.ID
}

func (r *IndexerProxySocks5Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxySocks5ResourceName
}

func (r *IndexerProxySocks5Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->\nIndexer Proxy Socks5 resource.\nFor more information refer to [Indexer Proxy](https://wiki.servarr.com/prowlarr/settings#indexer-proxies) and [Socks5](https://wiki.servarr.com/prowlarr/supported#socks5).",
		Attributes: map[string]schema.Attribute{
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
				Required:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *IndexerProxySocks5Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *IndexerProxySocks5Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var proxy *IndexerProxySocks5

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerProxySocks5
	request := proxy.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerProxyAPI.CreateIndexerProxy(r.auth).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerProxySocks5ResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerProxySocks5ResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxySocks5Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var proxy IndexerProxySocks5

	resp.Diagnostics.Append(req.State.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerProxySocks5 current value
	response, _, err := r.client.IndexerProxyAPI.GetIndexerProxyById(r.auth, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxySocks5ResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxySocks5ResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	proxy.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxySocks5Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var proxy *IndexerProxySocks5

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerProxySocks5
	request := proxy.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerProxyAPI.UpdateIndexerProxy(r.auth, strconv.Itoa(int(request.GetId()))).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerProxySocks5ResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerProxySocks5ResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxySocks5Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerProxySocks5 current value
	_, err := r.client.IndexerProxyAPI.DeleteIndexerProxy(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerProxySocks5ResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerProxySocks5ResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerProxySocks5Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerProxySocks5ResourceName+": "+req.ID)
}

func (i *IndexerProxySocks5) write(ctx context.Context, indexerProxy *prowlarr.IndexerProxyResource, diags *diag.Diagnostics) {
	genericIndexerProxy := i.toIndexerProxy()
	genericIndexerProxy.write(ctx, indexerProxy, diags)
	i.fromIndexerProxy(genericIndexerProxy)
}

func (i *IndexerProxySocks5) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.IndexerProxyResource {
	return i.toIndexerProxy().read(ctx, diags)
}
