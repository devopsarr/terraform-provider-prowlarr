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

const downloadClientDataSourceName = "download_client"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DownloadClientDataSource{}

func NewDownloadClientDataSource() datasource.DataSource {
	return &DownloadClientDataSource{}
}

// DownloadClientDataSource defines the download_client implementation.
type DownloadClientDataSource struct {
	client *prowlarr.APIClient
}

func (d *DownloadClientDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientDataSourceName
}

func (d *DownloadClientDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Download Clients -->Single [Download Client](../resources/download_client).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "DownloadClient configuration template.",
				Computed:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "DownloadClient implementation name.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"categories": schema.SetNestedAttribute{
				MarkdownDescription: "List of mapped categories.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of client category.",
							Computed:            true,
						},
						"categories": schema.SetAttribute{
							MarkdownDescription: "List of categories.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
			},
			// Field values
			"add_paused": schema.BoolAttribute{
				MarkdownDescription: "Add paused flag.",
				Computed:            true,
			},
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Computed:            true,
			},
			"start_on_add": schema.BoolAttribute{
				MarkdownDescription: "Start on add flag.",
				Computed:            true,
			},
			"add_stopped": schema.BoolAttribute{
				MarkdownDescription: "Add stopped flag.",
				Computed:            true,
			},
			"save_magnet_files": schema.BoolAttribute{
				MarkdownDescription: "Save magnet files flag.",
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Read only flag.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Computed:            true,
			},
			"item_priority": schema.Int64Attribute{
				MarkdownDescription: "Priority. `0` Last, `1` First.",
				Computed:            true,
			},
			"initial_state": schema.Int64Attribute{
				MarkdownDescription: "Initial state. `0` Start, `1` ForceStart, `2` Pause.",
				Computed:            true,
			},
			"intial_state": schema.Int64Attribute{
				MarkdownDescription: "Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Computed:            true,
				Sensitive:           true,
			},
			"rpc_path": schema.StringAttribute{
				MarkdownDescription: "RPC path.",
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "API URL.",
				Computed:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "App ID.",
				Computed:            true,
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "App Token.",
				Computed:            true,
				Sensitive:           true,
			},
			"secret_token": schema.StringAttribute{
				MarkdownDescription: "Secret token.",
				Computed:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Computed:            true,
				Sensitive:           true,
			},
			"tv_imported_category": schema.StringAttribute{
				MarkdownDescription: "TV imported category.",
				Computed:            true,
			},
			"destination_directory": schema.StringAttribute{
				MarkdownDescription: "Movie directory.",
				Computed:            true,
			},
			"directory": schema.StringAttribute{
				MarkdownDescription: "Directory.",
				Computed:            true,
			},
			"station_directory": schema.StringAttribute{
				MarkdownDescription: "Directory.",
				Computed:            true,
			},
			"destination": schema.StringAttribute{
				MarkdownDescription: "Destination.",
				Computed:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Category.",
				Computed:            true,
			},
			"nzb_folder": schema.StringAttribute{
				MarkdownDescription: "NZB folder.",
				Computed:            true,
			},
			"strm_folder": schema.StringAttribute{
				MarkdownDescription: "STRM folder.",
				Computed:            true,
			},
			"torrent_folder": schema.StringAttribute{
				MarkdownDescription: "Torrent folder.",
				Computed:            true,
			},
			"magnet_file_extension": schema.StringAttribute{
				MarkdownDescription: "Magnet file extension.",
				Computed:            true,
			},
			"additional_tags": schema.SetAttribute{
				MarkdownDescription: "Additional tags, `0` TitleSlug, `1` Quality, `2` Language, `3` ReleaseGroup, `4` Year, `5` Indexer, `6` Network.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"field_tags": schema.SetAttribute{
				MarkdownDescription: "Field tags.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"post_im_tags": schema.SetAttribute{
				MarkdownDescription: "Post import tags.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *DownloadClientDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *DownloadClientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *DownloadClient

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get downloadClient current value
	response, _, err := d.client.DownloadClientApi.ListDownloadClient(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientDataSourceName, err))

		return
	}

	data.find(ctx, data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+downloadClientDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DownloadClient) find(ctx context.Context, name string, downloadClients []*prowlarr.DownloadClientResource, diags *diag.Diagnostics) {
	for _, client := range downloadClients {
		if client.GetName() == name {
			d.write(ctx, client, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(downloadClientDataSourceName, "name", name))
}
