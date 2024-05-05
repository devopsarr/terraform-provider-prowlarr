package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// needed for tf debug mode
// var stderr = os.Stderr

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &ProwlarrProvider{}

// ProwlarrProvider defines the provider implementation.
type ProwlarrProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Prowlarr describes the provider data model.
type Prowlarr struct {
	ExtraHeaders types.Set    `tfsdk:"extra_headers"`
	APIKey       types.String `tfsdk:"api_key"`
	URL          types.String `tfsdk:"url"`
}

// ExtraHeader is part of Prowlarr.
type ExtraHeader struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// ProwlarrData defines auth and client to be used when connecting to Prowlarr.
type ProwlarrData struct {
	Auth   context.Context
	Client *prowlarr.APIClient
}

func (p *ProwlarrProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "prowlarr"
	resp.Version = p.version
}

func (p *ProwlarrProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Prowlarr provider is used to interact with any [Prowlarr](https://prowlarr.com/) installation. You must configure the provider with the proper credentials before you can use it. Use the left navigation to read about the available resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for Prowlarr authentication. Can be specified via the `PROWLARR_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Full Prowlarr URL with protocol and port (e.g. `https://test.prowlarr.audio:8686`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `PROWLARR_URL` environment variable.",
				Optional:            true,
			},
			"extra_headers": schema.SetNestedAttribute{
				MarkdownDescription: "Extra headers to be sent along with all Prowlarr requests. If this attribute is unset, it can be specified via environment variables following this pattern `PROWLARR_EXTRA_HEADER_${Header-Name}=${Header-Value}`.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Header name.",
							Required:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Header value.",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (p *ProwlarrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Prowlarr

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Extract URL
	APIURL := data.URL.ValueString()
	if APIURL == "" {
		APIURL = os.Getenv("PROWLARR_URL")
	}

	parsedAPIURL, err := url.Parse(APIURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find valid URL",
			"URL cannot parsed",
		)

		return
	}

	// Extract key
	key := data.APIKey.ValueString()
	if key == "" {
		key = os.Getenv("PROWLARR_API_KEY")
	}

	if key == "" {
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}

	// Init config
	config := prowlarr.NewConfiguration()
	// Check extra headers
	if len(data.ExtraHeaders.Elements()) > 0 {
		headers := make([]ExtraHeader, len(data.ExtraHeaders.Elements()))
		resp.Diagnostics.Append(data.ExtraHeaders.ElementsAs(ctx, &headers, false)...)

		for _, header := range headers {
			config.AddDefaultHeader(header.Name.ValueString(), header.Value.ValueString())
		}
	} else {
		env := os.Environ()
		for _, v := range env {
			if strings.HasPrefix(v, "PROWLARR_EXTRA_HEADER_") {
				header := strings.Split(v, "=")
				config.AddDefaultHeader(strings.TrimPrefix(header[0], "PROWLARR_EXTRA_HEADER_"), header[1])
			}
		}
	}

	// Set context for API calls
	auth := context.WithValue(
		context.Background(),
		prowlarr.ContextAPIKeys,
		map[string]prowlarr.APIKey{
			"X-Api-Key": {Key: key},
		},
	)
	auth = context.WithValue(auth, prowlarr.ContextServerVariables, map[string]string{
		"protocol": parsedAPIURL.Scheme,
		"hostpath": parsedAPIURL.Host,
	})

	prowlarrData := ProwlarrData{
		Auth:   auth,
		Client: prowlarr.NewAPIClient(config),
	}
	resp.DataSourceData = &prowlarrData
	resp.ResourceData = &prowlarrData
}

func (p *ProwlarrProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Applications
		NewSyncProfileResource,
		NewApplicationResource,
		NewApplicationLazyLibrarianResource,
		NewApplicationLidarrResource,
		NewApplicationMylarResource,
		NewApplicationRadarrResource,
		NewApplicationReadarrResource,
		NewApplicationSonarrResource,
		NewApplicationWhisparrResource,

		// Download Clients
		NewDownloadClientResource,
		NewDownloadClientAria2Resource,
		NewDownloadClientDelugeResource,
		NewDownloadClientFloodResource,
		NewDownloadClientFreeboxResource,
		NewDownloadClientHadoukenResource,
		NewDownloadClientNzbgetResource,
		NewDownloadClientNzbvortexResource,
		NewDownloadClientPneumaticResource,
		NewDownloadClientQbittorrentResource,
		NewDownloadClientRtorrentResource,
		NewDownloadClientSabnzbdResource,
		NewDownloadClientTorrentBlackholeResource,
		NewDownloadClientTorrentDownloadStationResource,
		NewDownloadClientTransmissionResource,
		NewDownloadClientUsenetBlackholeResource,
		NewDownloadClientUsenetDownloadStationResource,
		NewDownloadClientUtorrentResource,
		NewDownloadClientVuzeResource,

		// Indexer Proxies
		NewIndexerProxyResource,
		NewIndexerProxyFlaresolverrResource,
		NewIndexerProxyHTTPResource,
		NewIndexerProxySocks4Resource,
		NewIndexerProxySocks5Resource,

		// Indexer
		NewIndexerResource,

		// Notifications
		NewNotificationResource,
		NewNotificationAppriseResource,
		NewNotificationCustomScriptResource,
		NewNotificationDiscordResource,
		NewNotificationEmailResource,
		NewNotificationGotifyResource,
		NewNotificationJoinResource,
		NewNotificationMailgunResource,
		NewNotificationNotifiarrResource,
		NewNotificationNtfyResource,
		NewNotificationProwlResource,
		NewNotificationPushbulletResource,
		NewNotificationPushoverResource,
		NewNotificationSendgridResource,
		NewNotificationSignalResource,
		NewNotificationSimplepushResource,
		NewNotificationSlackResource,
		NewNotificationTelegramResource,
		NewNotificationTwitterResource,
		NewNotificationWebhookResource,

		// System
		NewHostResource,

		// Tags
		NewTagResource,
	}
}

func (p *ProwlarrProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Applications
		NewSyncProfileDataSource,
		NewSyncProfilesDataSource,
		NewApplicationDataSource,
		NewApplicationsDataSource,

		// Download Clients
		NewDownloadClientDataSource,
		NewDownloadClientsDataSource,

		// Indexer Proxies
		NewIndexerProxyDataSource,
		NewIndexerProxiesDataSource,

		// Indexer
		NewIndexerDataSource,
		NewIndexersDataSource,
		NewIndexerSchemaDataSource,
		NewIndexerSchemasDataSource,

		// Notifications
		NewNotificationDataSource,
		NewNotificationsDataSource,

		// System
		NewHostDataSource,
		NewSystemStatusDataSource,

		// Tags
		NewTagDataSource,
		NewTagsDataSource,
		NewTagDetailsDataSource,
		NewTagsDetailsDataSource,
	}
}

// New returns the provider with a specific version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ProwlarrProvider{
			version: version,
		}
	}
}

// ResourceConfigure is a helper function to set the client for a specific resource.
func resourceConfigure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) (context.Context, *prowlarr.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(*ProwlarrData)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *ProwlarrData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData.Auth, providerData.Client
}

func dataSourceConfigure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) (context.Context, *prowlarr.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(*ProwlarrData)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *ProwlarrData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData.Auth, providerData.Client
}
