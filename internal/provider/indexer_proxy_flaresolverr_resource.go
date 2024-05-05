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
	indexerProxyFlaresolverrResourceName   = "indexer_proxy_flaresolverr"
	indexerProxyFlaresolverrImplementation = "Flaresolverr"
	indexerProxyFlaresolverrConfigContract = "FlaresolverrSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerProxyFlaresolverrResource{}
	_ resource.ResourceWithImportState = &IndexerProxyFlaresolverrResource{}
)

func NewIndexerProxyFlaresolverrResource() resource.Resource {
	return &IndexerProxyFlaresolverrResource{}
}

// IndexerProxyFlaresolverrResource defines the indexer proxy implementation.
type IndexerProxyFlaresolverrResource struct {
	client *prowlarr.APIClient
	auth   context.Context
}

// IndexerProxyFlaresolverr describes the indexer proxy data model.
type IndexerProxyFlaresolverr struct {
	Tags           types.Set    `tfsdk:"tags"`
	Name           types.String `tfsdk:"name"`
	Host           types.String `tfsdk:"host"`
	RequestTimeout types.Int64  `tfsdk:"request_timeout"`
	ID             types.Int64  `tfsdk:"id"`
}

func (i IndexerProxyFlaresolverr) toIndexerProxy() *IndexerProxy {
	return &IndexerProxy{
		Tags:           i.Tags,
		Name:           i.Name,
		Host:           i.Host,
		RequestTimeout: i.RequestTimeout,
		ID:             i.ID,
		ConfigContract: types.StringValue(indexerProxyFlaresolverrConfigContract),
		Implementation: types.StringValue(indexerProxyFlaresolverrImplementation),
	}
}

func (i *IndexerProxyFlaresolverr) fromIndexerProxy(proxy *IndexerProxy) {
	i.Tags = proxy.Tags
	i.Name = proxy.Name
	i.Host = proxy.Host
	i.RequestTimeout = proxy.RequestTimeout
	i.ID = proxy.ID
}

func (r *IndexerProxyFlaresolverrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxyFlaresolverrResourceName
}

func (r *IndexerProxyFlaresolverrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->\nIndexer Proxy Flaresolverr resource.\nFor more information refer to [Indexer Proxy](https://wiki.servarr.com/prowlarr/settings#indexer-proxies) and [Flaresolverr](https://wiki.servarr.com/prowlarr/supported#flaresolverr).",
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
			"request_timeout": schema.Int64Attribute{
				MarkdownDescription: "Request timeout.",
				Required:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Required:            true,
			},
		},
	}
}

func (r *IndexerProxyFlaresolverrResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *IndexerProxyFlaresolverrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var proxy *IndexerProxyFlaresolverr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerProxyFlaresolverr
	request := proxy.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerProxyAPI.CreateIndexerProxy(r.auth).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerProxyFlaresolverrResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerProxyFlaresolverrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxyFlaresolverrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var proxy IndexerProxyFlaresolverr

	resp.Diagnostics.Append(req.State.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerProxyFlaresolverr current value
	response, _, err := r.client.IndexerProxyAPI.GetIndexerProxyById(r.auth, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxyFlaresolverrResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxyFlaresolverrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	proxy.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxyFlaresolverrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var proxy *IndexerProxyFlaresolverr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerProxyFlaresolverr
	request := proxy.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerProxyAPI.UpdateIndexerProxy(r.auth, strconv.Itoa(int(request.GetId()))).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerProxyFlaresolverrResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerProxyFlaresolverrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxyFlaresolverrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerProxyFlaresolverr current value
	_, err := r.client.IndexerProxyAPI.DeleteIndexerProxy(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerProxyFlaresolverrResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerProxyFlaresolverrResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerProxyFlaresolverrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerProxyFlaresolverrResourceName+": "+req.ID)
}

func (i *IndexerProxyFlaresolverr) write(ctx context.Context, indexerProxy *prowlarr.IndexerProxyResource, diags *diag.Diagnostics) {
	genericIndexerProxy := i.toIndexerProxy()
	genericIndexerProxy.write(ctx, indexerProxy, diags)
	i.fromIndexerProxy(genericIndexerProxy)
}

func (i *IndexerProxyFlaresolverr) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.IndexerProxyResource {
	return i.toIndexerProxy().read(ctx, diags)
}
