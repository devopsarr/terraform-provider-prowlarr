package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const applicationsDataSourceName = "applications"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ApplicationsDataSource{}

func NewApplicationsDataSource() datasource.DataSource {
	return &ApplicationsDataSource{}
}

// ApplicationsDataSource defines the applications implementation.
type ApplicationsDataSource struct {
	client *prowlarr.APIClient
}

// Applications describes the applications data model.
type Applications struct {
	Applications types.Set    `tfsdk:"applications"`
	ID           types.String `tfsdk:"id"`
}

func (d *ApplicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + applicationsDataSourceName
}

func (d *ApplicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Applications -->\nList all available [Applications](../resources/application).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"applications": schema.SetNestedAttribute{
				MarkdownDescription: "Application list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"config_contract": schema.StringAttribute{
							MarkdownDescription: "Application configuration template.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "Application implementation name.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Application name.",
							Computed:            true,
						},
						"sync_level": schema.StringAttribute{
							MarkdownDescription: "Sync level.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Application ID.",
							Computed:            true,
						},
						// Field values
						"base_url": schema.StringAttribute{
							MarkdownDescription: "Base URL.",
							Computed:            true,
						},
						"prowlarr_url": schema.StringAttribute{
							MarkdownDescription: "Prowlarr URL.",
							Computed:            true,
						},
						"api_key": schema.StringAttribute{
							MarkdownDescription: "API key.",
							Computed:            true,
							Sensitive:           true,
						},
						"sync_categories": schema.SetAttribute{
							MarkdownDescription: "Sync categories.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"anime_sync_categories": schema.SetAttribute{
							MarkdownDescription: "Anime sync categories.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *ApplicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ApplicationsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get applications current value
	response, _, err := d.client.ApplicationAPI.ListApplications(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, applicationsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+applicationsDataSourceName)
	// Map response body to resource schema attribute
	applications := make([]Application, len(response))
	for i, a := range response {
		applications[i].write(ctx, &a, &resp.Diagnostics)
	}

	appList, diags := types.SetValueFrom(ctx, Application{}.getType(), applications)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, Applications{Applications: appList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
