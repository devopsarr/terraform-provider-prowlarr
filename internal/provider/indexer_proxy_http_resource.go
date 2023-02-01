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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	indexerProxyHTTPResourceName   = "indexer_proxy_http"
	indexerProxyHTTPImplementation = "HTTP"
	indexerProxyHTTPConfigContract = "HTTPSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerProxyHTTPResource{}
	_ resource.ResourceWithImportState = &IndexerProxyHTTPResource{}
)

func NewIndexerProxyHTTPResource() resource.Resource {
	return &IndexerProxyHTTPResource{}
}

// IndexerProxyHTTPResource defines the indexer proxy implementation.
type IndexerProxyHTTPResource struct {
	client *prowlarr.APIClient
}

// IndexerProxyHTTP describes the indexer proxy data model.
type IndexerProxyHTTP struct {
	Tags     types.Set    `tfsdk:"tags"`
	Name     types.String `tfsdk:"name"`
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Port     types.Int64  `tfsdk:"port"`
	ID       types.Int64  `tfsdk:"id"`
}

func (d IndexerProxyHTTP) toIndexerProxy() *IndexerProxy {
	return &IndexerProxy{
		Tags:     d.Tags,
		Name:     d.Name,
		Host:     d.Host,
		Username: d.Username,
		Password: d.Password,
		Port:     d.Port,
		ID:       d.ID,
	}
}

func (d *IndexerProxyHTTP) fromIndexerProxy(proxy *IndexerProxy) {
	d.Tags = proxy.Tags
	d.Name = proxy.Name
	d.Host = proxy.Host
	d.Username = proxy.Username
	d.Password = proxy.Password
	d.Port = proxy.Port
	d.ID = proxy.ID
}

func (r *IndexerProxyHTTPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxyHTTPResourceName
}

func (r *IndexerProxyHTTPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->Indexer Proxy HTTP resource.\nFor more information refer to [Indexer Proxy](https://wiki.servarr.com/prowlarr/settings#indexer-proxies) and [HTTP](https://wiki.servarr.com/prowlarr/supported#http).",
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

func (r *IndexerProxyHTTPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerProxyHTTPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var proxy *IndexerProxyHTTP

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerProxyHTTP
	request := proxy.read(ctx)

	response, _, err := r.client.IndexerProxyApi.CreateIndexerProxy(ctx).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerProxyHTTPResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerProxyHTTPResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxyHTTPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var proxy IndexerProxyHTTP

	resp.Diagnostics.Append(req.State.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerProxyHTTP current value
	response, _, err := r.client.IndexerProxyApi.GetIndexerProxyById(ctx, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxyHTTPResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxyHTTPResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	proxy.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxyHTTPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var proxy *IndexerProxyHTTP

	resp.Diagnostics.Append(req.Plan.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerProxyHTTP
	request := proxy.read(ctx)

	response, _, err := r.client.IndexerProxyApi.UpdateIndexerProxy(ctx, strconv.Itoa(int(request.GetId()))).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerProxyHTTPResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerProxyHTTPResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxyHTTPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var proxy *IndexerProxyHTTP

	resp.Diagnostics.Append(req.State.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerProxyHTTP current value
	_, err := r.client.IndexerProxyApi.DeleteIndexerProxy(ctx, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxyHTTPResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerProxyHTTPResourceName+": "+strconv.Itoa(int(proxy.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerProxyHTTPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerProxyHTTPResourceName+": "+req.ID)
}

func (d *IndexerProxyHTTP) write(ctx context.Context, indexerProxy *prowlarr.IndexerProxyResource) {
	genericIndexerProxy := IndexerProxy{
		ID:   types.Int64Value(int64(indexerProxy.GetId())),
		Name: types.StringValue(indexerProxy.GetName()),
		Tags: types.SetValueMust(types.Int64Type, nil),
	}

	tfsdk.ValueFrom(ctx, indexerProxy.Tags, genericIndexerProxy.Tags.Type(ctx), &genericIndexerProxy.Tags)
	genericIndexerProxy.writeFields(indexerProxy.GetFields())
	d.fromIndexerProxy(&genericIndexerProxy)
}

func (d *IndexerProxyHTTP) read(ctx context.Context) *prowlarr.IndexerProxyResource {
	tags := make([]*int32, len(d.Tags.Elements()))

	tfsdk.ValueAs(ctx, d.Tags, &tags)

	proxy := prowlarr.NewIndexerProxyResource()
	proxy.SetId(int32(d.ID.ValueInt64()))
	proxy.SetConfigContract(indexerProxyHTTPConfigContract)
	proxy.SetImplementation(indexerProxyHTTPImplementation)
	proxy.SetName(d.Name.ValueString())
	proxy.SetTags(tags)
	proxy.SetFields(d.toIndexerProxy().readFields())

	return proxy
}
