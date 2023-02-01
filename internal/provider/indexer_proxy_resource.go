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
	"golang.org/x/exp/slices"
)

const indexerProxyResourceName = "indexer_proxy"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerProxyResource{}
	_ resource.ResourceWithImportState = &IndexerProxyResource{}
)

var (
	indexerProxyIntFields    = []string{"port", "requestTimeout"}
	indexerProxyStringFields = []string{"host", "password", "username"}
)

func NewIndexerProxyResource() resource.Resource {
	return &IndexerProxyResource{}
}

// IndexerProxyResource defines the indexer proxy implementation.
type IndexerProxyResource struct {
	client *prowlarr.APIClient
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

// ProxyCategory is part of IndexerProxy.
type ProxyCategory struct {
	Categories types.Set    `tfsdk:"categories"`
	Name       types.String `tfsdk:"name"`
}

func (r *IndexerProxyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerProxyResourceName
}

func (r *IndexerProxyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexer Proxies -->Generic Indexer Proxy resource. When possible use a specific resource instead.\nFor more information refer to [Indexer Proxy](https://wiki.servarr.com/prowlarr/settings#indexer-proxies).",
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
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
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
	request := proxy.read(ctx)

	response, _, err := r.client.IndexerProxyApi.CreateIndexerProxy(ctx).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerProxyResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerProxyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state IndexerProxy

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerProxyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var proxy IndexerProxy

	resp.Diagnostics.Append(req.State.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerProxy current value
	response, _, err := r.client.IndexerProxyApi.GetIndexerProxyById(ctx, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxyResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerProxyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state IndexerProxy

	state.write(ctx, response)
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
	request := proxy.read(ctx)

	response, _, err := r.client.IndexerProxyApi.UpdateIndexerProxy(ctx, strconv.Itoa(int(request.GetId()))).IndexerProxyResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerProxyResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerProxyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state IndexerProxy

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerProxyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var proxy *IndexerProxy

	resp.Diagnostics.Append(req.State.Get(ctx, &proxy)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerProxy current value
	_, err := r.client.IndexerProxyApi.DeleteIndexerProxy(ctx, int32(proxy.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerProxyResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerProxyResourceName+": "+strconv.Itoa(int(proxy.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerProxyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerProxyResourceName+": "+req.ID)
}

func (i *IndexerProxy) write(ctx context.Context, indexerProxy *prowlarr.IndexerProxyResource) {
	i.ID = types.Int64Value(int64(indexerProxy.GetId()))
	i.ConfigContract = types.StringValue(indexerProxy.GetConfigContract())
	i.Implementation = types.StringValue(indexerProxy.GetImplementation())
	i.Name = types.StringValue(indexerProxy.GetName())
	i.Tags = types.SetValueMust(types.Int64Type, nil)

	tfsdk.ValueFrom(ctx, indexerProxy.Tags, i.Tags.Type(ctx), &i.Tags)
	i.writeFields(indexerProxy.GetFields())
}

func (i *IndexerProxy) writeFields(fields []*prowlarr.Field) {
	for _, f := range fields {
		if f.Value == nil {
			continue
		}

		if slices.Contains(indexerProxyStringFields, f.GetName()) {
			helpers.WriteStringField(f, i)

			continue
		}

		if slices.Contains(indexerProxyIntFields, f.GetName()) {
			helpers.WriteIntField(f, i)

			continue
		}
	}
}

func (i *IndexerProxy) read(ctx context.Context) *prowlarr.IndexerProxyResource {
	tags := make([]*int32, len(i.Tags.Elements()))

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	proxy := prowlarr.NewIndexerProxyResource()
	proxy.SetId(int32(i.ID.ValueInt64()))
	proxy.SetConfigContract(i.ConfigContract.ValueString())
	proxy.SetImplementation(i.Implementation.ValueString())
	proxy.SetName(i.Name.ValueString())
	proxy.SetTags(tags)
	proxy.SetFields(i.readFields())

	return proxy
}

func (i *IndexerProxy) readFields() []*prowlarr.Field {
	var output []*prowlarr.Field

	for _, j := range indexerProxyIntFields {
		if field := helpers.ReadIntField(j, i); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range indexerProxyStringFields {
		if field := helpers.ReadStringField(s, i); field != nil {
			output = append(output, field)
		}
	}

	return output
}
