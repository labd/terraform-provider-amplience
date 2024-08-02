package hub

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/amplience-go-sdk/content"
)

type Hub struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
	Settings    *Settings    `tfsdk:"settings"`
}

type Settings struct {
	Publishing                       *Publishing                       `tfsdk:"publishing"`
	Devices                          Devices                           `tfsdk:"devices"`
	Localization                     *Localization                     `tfsdk:"localization"`
	Applications                     Applications                      `tfsdk:"applications"`
	PreviewVirtualStagingEnvironment *PreviewVirtualStagingEnvironment `tfsdk:"preview_virtual_staging_environment"`
	VirtualStagingEnvironment        *VirtualStagingEnvironment        `tfsdk:"virtual_staging_environment"`
	AssetManagement                  *AssetManagement                  `tfsdk:"asset_management"`
}

type Publishing struct {
	Platforms *Platforms `tfsdk:"platforms"`
}

type Platforms struct {
	AmplienceDAM *AmplienceDAM `tfsdk:"amplience_dam"`
}

type AmplienceDAM struct {
	APIKey    types.String `tfsdk:"api_key"`
	APISecret types.String `tfsdk:"api_secret"`
	Endpoint  types.String `tfsdk:"endpoint"`
}

type Devices []DeviceSettings

type DeviceSettings struct {
	Name      types.String `tfsdk:"name"`
	Width     types.Int64  `tfsdk:"width"`
	Height    types.Int64  `tfsdk:"height"`
	Orientate types.Bool   `tfsdk:"orientate"`
}

type Localization struct {
	Locales []types.String `tfsdk:"locales"`
}

type Applications []Application

type Application struct {
	Name        types.String `tfsdk:"name"`
	TemplateURI types.String `tfsdk:"template_uri"`
}

type PreviewVirtualStagingEnvironment struct {
	Hostname types.String `tfsdk:"hostname"`
}

type VirtualStagingEnvironment struct {
	Hostname types.String `tfsdk:"hostname"`
}

type AssetManagement struct {
	Enabled      types.Bool   `tfsdk:"enabled"`
	ClientConfig types.String `tfsdk:"client_config"`
}

func (h *Hub) ToUpdateInput() content.HubUpdateInput {
	return content.HubUpdateInput{
		Name:        h.Name.ValueString(),
		Label:       h.Label.ValueString(),
		Description: h.Description.ValueStringPointer(),
		Settings:    h.Settings.ToUpdateInput(),
	}
}

func (h *Hub) setSecretValuesFromState(s Hub) {
	if h.Settings == nil || h.Settings.Publishing == nil || h.Settings.Publishing.Platforms == nil || h.Settings.Publishing.Platforms.AmplienceDAM == nil {
		return
	}

	if s.Settings == nil || s.Settings.Publishing == nil || s.Settings.Publishing.Platforms == nil || s.Settings.Publishing.Platforms.AmplienceDAM == nil {
		return
	}

	h.Settings.Publishing.Platforms.AmplienceDAM.APISecret = s.Settings.Publishing.Platforms.AmplienceDAM.APISecret
}

func (s *Settings) ToUpdateInput() *content.Settings {
	if s == nil {
		return nil
	}

	return &content.Settings{
		Publishing:                       s.Publishing.ToUpdateInput(),
		Devices:                          s.Devices.ToUpdateInput(),
		Localization:                     s.Localization.ToUpdateInput(),
		Applications:                     s.Applications.ToUpdateInput(),
		PreviewVirtualStagingEnvironment: s.PreviewVirtualStagingEnvironment.ToUpdateInput(),
		VirtualStagingEnvironment:        s.VirtualStagingEnvironment.ToUpdateInput(),
		AssetManagement:                  s.AssetManagement.ToUpdateInput(),
	}
}

func (p *Publishing) ToUpdateInput() *content.PublishingSettings {
	if p == nil {
		return nil
	}

	return &content.PublishingSettings{
		Platforms: p.Platforms.ToUpdateInput(),
	}
}

func (p *Platforms) ToUpdateInput() *content.PlatformSettings {
	if p == nil {
		return nil
	}

	return &content.PlatformSettings{
		AmplienceDam: p.AmplienceDAM.ToUpdateInput(),
	}
}

func (a *AmplienceDAM) ToUpdateInput() *content.AmplienceDamSettings {
	if a == nil {
		return nil
	}

	return &content.AmplienceDamSettings{
		ApiKey:    a.APIKey.ValueString(),
		ApiSecret: a.APISecret.ValueString(),
		Endpoint:  a.Endpoint.ValueString(),
	}
}

func (d Devices) ToUpdateInput() []content.DeviceSettings {
	var devices = make([]content.DeviceSettings, 0, len(d))
	for _, device := range d {
		devices = append(devices, content.DeviceSettings{
			Name:      device.Name.ValueString(),
			Width:     int(device.Width.ValueInt64()),
			Height:    int(device.Height.ValueInt64()),
			Orientate: device.Orientate.ValueBool(),
		})
	}

	return devices
}

func (l *Localization) ToUpdateInput() *content.LocalizationSettings {
	if l == nil {
		return nil
	}

	var locales = make([]string, 0, len(l.Locales))
	for _, locale := range l.Locales {
		locales = append(locales, locale.ValueString())
	}

	return &content.LocalizationSettings{
		Locales: locales,
	}
}

func (a Applications) ToUpdateInput() []content.ApplicationSettings {
	var applications = make([]content.ApplicationSettings, 0, len(a))
	for _, app := range a {
		applications = append(applications, content.ApplicationSettings{
			Name:         app.Name.ValueString(),
			TemplatedUri: app.TemplateURI.ValueString(),
		})
	}

	return applications
}

func (p *PreviewVirtualStagingEnvironment) ToUpdateInput() *content.PreviewVirtualStagingEnvironmentSettings {
	if p == nil {
		return nil
	}

	return &content.PreviewVirtualStagingEnvironmentSettings{
		Hostname: p.Hostname.ValueString(),
	}
}

func (v *VirtualStagingEnvironment) ToUpdateInput() *content.VirtualStagingEnvironmentSettings {
	if v == nil {
		return nil
	}

	return &content.VirtualStagingEnvironmentSettings{
		Hostname: v.Hostname.ValueString(),
	}
}

func (a *AssetManagement) ToUpdateInput() *content.AssetManagementSettings {
	if a == nil {
		return nil
	}

	return &content.AssetManagementSettings{
		Enabled:      a.Enabled.ValueBoolPointer(),
		ClientConfig: a.ClientConfig.ValueStringPointer(),
	}
}

func NewHubFromNative(hub *content.Hub) *Hub {
	return &Hub{
		ID:          types.StringValue(hub.ID),
		Name:        types.StringValue(hub.Name),
		Label:       types.StringValue(hub.Label),
		Description: types.StringPointerValue(hub.Description),
		Settings:    NewSettingsFromNative(hub.Settings),
	}
}

func NewSettingsFromNative(settings *content.Settings) *Settings {
	if settings == nil {
		return nil
	}
	return &Settings{
		Publishing:                       NewPublishingFromNative(settings.Publishing),
		Devices:                          NewDevicesFromNative(settings.Devices),
		Localization:                     NewLocalizationFromNative(settings.Localization),
		Applications:                     NewApplicationsFromNative(settings.Applications),
		PreviewVirtualStagingEnvironment: NewPreviewVirtualStagingEnvironmentFromNative(settings.PreviewVirtualStagingEnvironment),
		VirtualStagingEnvironment:        NewVirtualStagingEnvironmentFromNative(settings.VirtualStagingEnvironment),
		AssetManagement:                  NewAssetManagementFromNative(settings.AssetManagement),
	}
}

func NewPublishingFromNative(publishing *content.PublishingSettings) *Publishing {
	if publishing == nil {
		return nil
	}
	return &Publishing{
		Platforms: NewPlatformsFromNative(publishing.Platforms),
	}
}

func NewPlatformsFromNative(platforms *content.PlatformSettings) *Platforms {
	if platforms == nil {
		return nil
	}
	return &Platforms{
		AmplienceDAM: NewAmplienceDAMFromNative(platforms.AmplienceDam),
	}
}

func NewAmplienceDAMFromNative(dam *content.AmplienceDamSettings) *AmplienceDAM {
	if dam == nil {
		return nil
	}
	return &AmplienceDAM{
		APIKey:    types.StringValue(dam.ApiKey),
		APISecret: types.StringValue(dam.ApiSecret),
		Endpoint:  types.StringValue(dam.Endpoint),
	}
}

func NewDevicesFromNative(devices []content.DeviceSettings) Devices {
	var result Devices
	for _, device := range devices {
		result = append(result, DeviceSettings{
			Name:      types.StringValue(device.Name),
			Width:     types.Int64Value(int64(device.Width)),
			Height:    types.Int64Value(int64(device.Height)),
			Orientate: types.BoolValue(device.Orientate),
		})
	}
	return result
}

func NewLocalizationFromNative(localization *content.LocalizationSettings) *Localization {
	if localization == nil {
		return nil
	}
	var locales []types.String
	for _, locale := range localization.Locales {
		locales = append(locales, types.StringValue(locale))
	}
	return &Localization{
		Locales: locales,
	}
}

func NewApplicationsFromNative(applications []content.ApplicationSettings) Applications {
	var result Applications
	for _, app := range applications {
		result = append(result, Application{
			Name:        types.StringValue(app.Name),
			TemplateURI: types.StringValue(app.TemplatedUri),
		})
	}
	return result
}

func NewPreviewVirtualStagingEnvironmentFromNative(env *content.PreviewVirtualStagingEnvironmentSettings) *PreviewVirtualStagingEnvironment {
	if env == nil {
		return nil
	}
	return &PreviewVirtualStagingEnvironment{
		Hostname: types.StringValue(env.Hostname),
	}
}

func NewVirtualStagingEnvironmentFromNative(env *content.VirtualStagingEnvironmentSettings) *VirtualStagingEnvironment {
	if env == nil {
		return nil
	}
	return &VirtualStagingEnvironment{
		Hostname: types.StringValue(env.Hostname),
	}
}

func NewAssetManagementFromNative(assetManagement *content.AssetManagementSettings) *AssetManagement {
	if assetManagement == nil {
		return nil
	}
	return &AssetManagement{
		Enabled:      types.BoolPointerValue(assetManagement.Enabled),
		ClientConfig: types.StringPointerValue(assetManagement.ClientConfig),
	}
}
