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
	indexerProxySocks4ResourceName   = "indexer_proxy_socks4"
	indexerProxySocks4Implementation = "Socks4"
	indexerProxySocks4ConfigContract = "Socks4Settings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerProxySocks4Resource{}
	_ resource.ResourceWithImportState = &IndexerProxySocks4Resource{}
)

func NewIndexerProxySocks4Resource() resource.Resource {
	return &IndexerProxySocks4Resource{}
}

// IndexerProxySocks4Resource defines the indexer proxy implementation.
type IndexerProxySocks4Resource struct {
	client *prowlarr.APIClient
}

// IndexerProxySocks4 describes the indexer proxy data model.
type IndexerProxySocks4 struct {
	Tags     types.Set    `tfsdk:"tags"`
	Name     types.String `tfsdk:"name"`
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Port     types.Int64  `tfsdk:"port"`
	ID       types.Int64  `tfsdk:"id"`
}

func (i IndexerProxySocks4) toIndexerProxy() *IndexerProxy {
	return &IndexerProxy{
		Tags:           i.Tags,
		Name:           i.Name,
		Host:           i.Host,
		Username:       i.Username,
		Password:       i.Password,
		Port:           i.Port,
		ID:             i.ID,
		ConfigContract: types.StringValue(indexerProxySocks4ConfigContract),
		Implementation: types.StringValue(indexerProxySocks4Implementation),
	}
}

func (i *IndexerProxySocks4) fromIndexerProxy(proxy *IndexerProxy) {
	i.Tags = proxy.Tags
	i.Name = proxy.Name
	i.Host = proxy.Host
	i.Username = proxy.Username
	i.Password = proxy.Password
	i.Port = proxy.Port
	i.ID = proxy.ID
}

func (r *IndexerProxySocks4Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxySocks4ResourceName
}

func (r *IndexerProxySocks4Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->Indexer Proxy Socks4 resource.\nFor more information refer to [Indexer Proxy](https://wiki.servarr.com/prowlarr/settings#indexer-proxies) and [Socks4](https://wiki.servarr.com/prowlarr/supported#socks4).",
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

func (r *IndexerProxySocks4Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerProxySocks4Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var proxy *IndexerProxySocks4

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerProxySocks4
	request := proxy.read(ctx)

	response, _, err := r.client.IndexerProxyApi.CreateIndexerProxy(ctx).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerProxySocks4ResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerProxySocks4ResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxySocks4Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var proxy IndexerProxySocks4

	resp.Diagnostics.Append(req.State.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerProxySocks4 current value
	response, _, err := r.client.IndexerProxyApi.GetIndexerProxyById(ctx, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxySocks4ResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxySocks4ResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	proxy.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxySocks4Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var proxy *IndexerProxySocks4

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerProxySocks4
	request := proxy.read(ctx)

	response, _, err := r.client.IndexerProxyApi.UpdateIndexerProxy(ctx, strconv.Itoa(int(request.GetId()))).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerProxySocks4ResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerProxySocks4ResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxySocks4Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerProxySocks4 current value
	_, err := r.client.IndexerProxyApi.DeleteIndexerProxy(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerProxySocks4ResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerProxySocks4ResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerProxySocks4Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerProxySocks4ResourceName+": "+req.ID)
}

func (i *IndexerProxySocks4) write(ctx context.Context, indexerProxy *prowlarr.IndexerProxyResource) {
	genericIndexerProxy := i.toIndexerProxy()
	genericIndexerProxy.write(ctx, indexerProxy)
	i.fromIndexerProxy(genericIndexerProxy)
}

func (i *IndexerProxySocks4) read(ctx context.Context) *prowlarr.IndexerProxyResource {
	return i.toIndexerProxy().read(ctx)
}
