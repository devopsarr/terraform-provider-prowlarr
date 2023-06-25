package provider

import (
	"context"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const applicationDataSourceName = "application"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ApplicationDataSource{}

func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

// ApplicationDataSource defines the application implementation.
type ApplicationDataSource struct {
	client *prowlarr.APIClient
}

func (d *ApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + applicationDataSourceName
}

func (d *ApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Applications -->Single [Application](../resources/application).",
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
				Required:            true,
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
	}
}

func (d *ApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Application

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get application current value
	response, _, err := d.client.ApplicationApi.ListApplications(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, applicationDataSourceName, err))

		return
	}

	data.find(ctx, data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+applicationDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (a *Application) find(ctx context.Context, name string, applications []*prowlarr.ApplicationResource, diags *diag.Diagnostics) {
	for _, app := range applications {
		if app.GetName() == name {
			a.write(ctx, app, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(applicationDataSourceName, "name", name))
}
