package provider

import (
	"context"
	"math/big"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerResourceName = "indexer"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerResource{}
	_ resource.ResourceWithImportState = &IndexerResource{}
)

func NewIndexerResource() resource.Resource {
	return &IndexerResource{}
}

// IndexerResource defines the indexer implementation.
type IndexerResource struct {
	client *prowlarr.APIClient
}

// Indexer describes the indexer data model.
type Indexer struct {
	Tags types.Set `tfsdk:"tags"`
	// IndexerURLs    types.Set    `tfsdk:"indexer_urls"`
	Fields         types.Set    `tfsdk:"fields"`
	ConfigContract types.String `tfsdk:"config_contract"`
	Implementation types.String `tfsdk:"implementation"`
	Name           types.String `tfsdk:"name"`
	Protocol       types.String `tfsdk:"protocol"`
	Language       types.String `tfsdk:"language"`
	Privacy        types.String `tfsdk:"privacy"`
	AppProfileID   types.Int64  `tfsdk:"app_profile_id"`
	Priority       types.Int64  `tfsdk:"priority"`
	ID             types.Int64  `tfsdk:"id"`
	Enable         types.Bool   `tfsdk:"enable"`
}

// Field is part of Indexer.
type Field struct {
	SetValue       types.Set    `tfsdk:"set_value"`
	NumberValue    types.Number `tfsdk:"number_value"`
	Name           types.String `tfsdk:"name"`
	TextValue      types.String `tfsdk:"text_value"`
	SensitiveValue types.String `tfsdk:"sensitive_value"`
	BoolValue      types.Bool   `tfsdk:"bool_value"`
}

func (r *IndexerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerResourceName
}

func (r *IndexerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Generic Indexer resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/prowlarr/indexers) documentation.",
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
			"app_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Application profile ID.",
				Optional:            true,
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Indexer configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Indexer implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Indexer name.",
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
			"language": schema.StringAttribute{
				MarkdownDescription: "Language.",
				Computed:            true,
			},
			"privacy": schema.StringAttribute{
				MarkdownDescription: "Privacy.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Indexer ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"fields": schema.SetNestedAttribute{
				Required:            true,
				MarkdownDescription: "Set of configuration fields.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getFieldSchema().Attributes,
				},
			},
		},
	}
}

func (r IndexerResource) getFieldSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Field name.",
				Required:            true,
			},
			"text_value": schema.StringAttribute{
				MarkdownDescription: "Text value. Only one value must be filled out.",
				Optional:            true,
				Computed:            true,
			},
			"sensitive_value": schema.StringAttribute{
				MarkdownDescription: "Sensitive string value. Only one value must be filled out.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"number_value": schema.NumberAttribute{
				MarkdownDescription: "Number value. Only one value must be filled out.",
				Optional:            true,
				Computed:            true,
			},
			"bool_value": schema.BoolAttribute{
				MarkdownDescription: "Bool value. Only one value must be filled out.",
				Optional:            true,
				Computed:            true,
			},
			"set_value": schema.SetAttribute{
				MarkdownDescription: "Set value. Only one value must be filled out.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *IndexerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Indexer
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, indexer)...)
}

func (r *IndexerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *Indexer

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Indexer current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, indexer)...)
}

func (r *IndexerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Indexer
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, indexer)...)
}

func (r *IndexerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Indexer current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerResourceName+": "+req.ID)
}

func (i *Indexer) write(ctx context.Context, indexer *prowlarr.IndexerResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	i.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, indexer.Tags)
	diags.Append(localDiag...)

	i.Enable = types.BoolValue(indexer.GetEnable())
	i.Priority = types.Int64Value(int64(indexer.GetPriority()))
	i.AppProfileID = types.Int64Value(int64(indexer.GetAppProfileId()))
	i.ID = types.Int64Value(int64(indexer.GetId()))
	i.ConfigContract = types.StringValue(indexer.GetConfigContract())
	i.Implementation = types.StringValue(indexer.GetImplementation())
	i.Name = types.StringValue(indexer.GetName())
	i.Protocol = types.StringValue(string(indexer.GetProtocol()))
	i.Language = types.StringValue(indexer.GetLanguage())
	i.Privacy = types.StringValue(string(indexer.GetPrivacy()))

	var fields []Field

	for _, f := range indexer.GetFields() {
		if _, ok := f.GetValueOk(); ok {
			var field Field

			field.write(ctx, f, diags)
			fields = append(fields, field)
		}
	}

	i.Fields, localDiag = types.SetValueFrom(ctx, IndexerResource{}.getFieldSchema().Type(), fields)
	diags.Append(localDiag...)
}

func (f *Field) write(ctx context.Context, field *prowlarr.Field, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	f.Name = types.StringValue(field.GetName())
	// init all values to null
	f.BoolValue = types.BoolNull()
	f.NumberValue = types.NumberNull()
	f.SensitiveValue = types.StringNull()
	f.TextValue = types.StringNull()
	f.SetValue = types.SetNull(types.Int64Type)

	if _, ok := field.GetValueOk(); ok {
		switch v := field.GetValue().(type) {
		case bool:
			f.BoolValue = types.BoolValue(v)
		case float64:
			f.NumberValue = types.NumberValue(big.NewFloat(v))
		case string:
			if field.GetType() == "password" {
				f.SensitiveValue = types.StringValue(v)
			} else {
				f.TextValue = types.StringValue(v)
			}
		case []interface{}:
			setValue := make([]*int64, len(v))

			for i, value := range v {
				if element, ok := value.(float64); ok {
					point := int64(element)
					setValue[i] = &point
				}
			}

			f.SetValue, tempDiag = types.SetValueFrom(ctx, types.Int64Type, setValue)
			diags.Append(tempDiag...)
		}
	}
}

func (i *Indexer) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.IndexerResource {
	fieldList := make([]Field, len(i.Fields.Elements()))
	diags.Append(i.Fields.ElementsAs(ctx, &fieldList, true)...)

	fields := make([]*prowlarr.Field, len(fieldList))

	for n, f := range fieldList {
		fields[n] = f.read(ctx, diags)
	}

	indexer := prowlarr.NewIndexerResource()
	indexer.SetEnable(i.Enable.ValueBool())
	indexer.SetPriority(int32(i.Priority.ValueInt64()))
	indexer.SetAppProfileId(int32(i.AppProfileID.ValueInt64()))
	indexer.SetId(int32(i.ID.ValueInt64()))
	indexer.SetConfigContract(i.ConfigContract.ValueString())
	indexer.SetImplementation(i.Implementation.ValueString())
	indexer.SetName(i.Name.ValueString())
	indexer.SetProtocol(prowlarr.DownloadProtocol(i.Protocol.ValueString()))
	diags.Append(i.Tags.ElementsAs(ctx, &indexer.Tags, true)...)
	indexer.SetFields(fields)

	return indexer
}

func (f *Field) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.Field {
	field := prowlarr.NewField()
	field.SetName(f.Name.ValueString())

	if !f.BoolValue.IsNull() && !f.BoolValue.IsUnknown() {
		field.SetValue(f.BoolValue.ValueBool())
	}

	if !f.NumberValue.IsNull() && !f.NumberValue.IsUnknown() {
		field.SetValue(f.NumberValue.ValueBigFloat())
	}

	if !f.TextValue.IsNull() && !f.TextValue.IsUnknown() {
		field.SetValue(f.TextValue.ValueString())
	}

	if !f.SensitiveValue.IsNull() && !f.SensitiveValue.IsUnknown() {
		field.SetValue(f.SensitiveValue.ValueString())
	}

	if !f.SetValue.IsNull() && !f.SetValue.IsUnknown() {
		set := make([]*int64, len(f.SetValue.Elements()))
		diags.Append(f.SetValue.ElementsAs(ctx, &set, true)...)
		field.SetValue(set)
	}

	return field
}
