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

func (d IndexerProxySocks5) toIndexerProxy() *IndexerProxy {
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

func (d *IndexerProxySocks5) fromIndexerProxy(proxy *IndexerProxy) {
	d.Tags = proxy.Tags
	d.Name = proxy.Name
	d.Host = proxy.Host
	d.Username = proxy.Username
	d.Password = proxy.Password
	d.Port = proxy.Port
	d.ID = proxy.ID
}

func (r *IndexerProxySocks5Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxySocks5ResourceName
}

func (r *IndexerProxySocks5Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->Indexer Proxy Socks5 resource.\nFor more information refer to [Indexer Proxy](https://wiki.servarr.com/prowlarr/settings#indexer-proxys) and [Socks5](https://wiki.servarr.com/prowlarr/supported#socks5).",
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
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
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
	request := proxy.read(ctx)

	response, _, err := r.client.IndexerProxyApi.CreateIndexerProxy(ctx).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerProxySocks5ResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerProxySocks5ResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response)
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
	response, _, err := r.client.IndexerProxyApi.GetIndexerProxyById(ctx, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxySocks5ResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxySocks5ResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	proxy.write(ctx, response)
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
	request := proxy.read(ctx)

	response, _, err := r.client.IndexerProxyApi.UpdateIndexerProxy(ctx, strconv.Itoa(int(request.GetId()))).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerProxySocks5ResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerProxySocks5ResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	proxy.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &proxy)...)
}

func (r *IndexerProxySocks5Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var proxy *IndexerProxySocks5

	resp.Diagnostics.Append(req.State.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerProxySocks5 current value
	_, err := r.client.IndexerProxyApi.DeleteIndexerProxy(ctx, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxySocks5ResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerProxySocks5ResourceName+": "+strconv.Itoa(int(proxy.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerProxySocks5Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerProxySocks5ResourceName+": "+req.ID)
}

func (d *IndexerProxySocks5) write(ctx context.Context, indexerProxy *prowlarr.IndexerProxyResource) {
	genericIndexerProxy := IndexerProxy{
		ID:   types.Int64Value(int64(indexerProxy.GetId())),
		Name: types.StringValue(indexerProxy.GetName()),
		Tags: types.SetValueMust(types.Int64Type, nil),
	}

	tfsdk.ValueFrom(ctx, indexerProxy.Tags, genericIndexerProxy.Tags.Type(ctx), &genericIndexerProxy.Tags)
	genericIndexerProxy.writeFields(indexerProxy.GetFields())
	d.fromIndexerProxy(&genericIndexerProxy)
}

func (d *IndexerProxySocks5) read(ctx context.Context) *prowlarr.IndexerProxyResource {
	tags := make([]*int32, len(d.Tags.Elements()))

	tfsdk.ValueAs(ctx, d.Tags, &tags)

	proxy := prowlarr.NewIndexerProxyResource()
	proxy.SetId(int32(d.ID.ValueInt64()))
	proxy.SetConfigContract(indexerProxySocks5ConfigContract)
	proxy.SetImplementation(indexerProxySocks5Implementation)
	proxy.SetName(d.Name.ValueString())
	proxy.SetTags(tags)
	proxy.SetFields(d.toIndexerProxy().readFields())

	return proxy
}
