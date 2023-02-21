package provider

import (
	"context"
	"os"

	"github.com/devopsarr/prowlarr-go/prowlarr"

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

// ScaffoldingProvider defines the provider implementation.
type ProwlarrProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Prowlarr describes the provider data model.
type Prowlarr struct {
	APIKey types.String `tfsdk:"api_key"`
	URL    types.String `tfsdk:"url"`
}

func (p *ProwlarrProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "prowlarr"
	resp.Version = p.version
}

func (p *ProwlarrProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Prowlarr provider is used to interact with any [Prowlarr](https://prowlarr.com/) installation. You must configure the provider with the proper credentials before you can use it. Use the left navigation to read about the available resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for Prowlarr authentication. Can be specified via the `PROWLARR_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Full Prowlarr URL with protocol and port (e.g. `https://test.prowlarr.com:9696`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `PROWLARR_URL` environment variable.",
				Optional:            true,
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

	// User must provide URL to the provider
	if data.URL.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as url",
		)

		return
	}

	var url string
	if data.URL.IsNull() {
		url = os.Getenv("PROWLARR_URL")
	} else {
		url = data.URL.ValueString()
	}

	if url == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find URL",
			"URL cannot be an empty string",
		)

		return
	}

	// User must provide API key to the provider
	if data.APIKey.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_key",
		)

		return
	}

	var key string
	if data.APIKey.IsNull() {
		key = os.Getenv("PROWLARR_API_KEY")
	} else {
		key = data.APIKey.ValueString()
	}

	if key == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}

	// Configuring client. API Key management could be changed once new options avail in sdk.
	config := prowlarr.NewConfiguration()
	config.AddDefaultHeader("X-Api-Key", key)
	config.Servers[0].URL = url
	client := prowlarr.NewAPIClient(config)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ProwlarrProvider) Resources(ctx context.Context) []func() resource.Resource {
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
		NewDownloadClientTransmissionResource,

		// Indexer Proxies
		NewIndexerProxyResource,
		NewIndexerProxyFlaresolverrResource,
		NewIndexerProxyHTTPResource,
		NewIndexerProxySocks4Resource,
		NewIndexerProxySocks5Resource,

		// Notifications
		NewNotificationResource,
		NewNotificationCustomScriptResource,
		NewNotificationWebhookResource,

		// Tags
		NewTagResource,
	}
}

func (p *ProwlarrProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
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

		// Notifications
		NewNotificationDataSource,
		NewNotificationsDataSource,

		// System Status
		NewSystemStatusDataSource,

		// Tags
		NewTagDataSource,
		NewTagsDataSource,
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
