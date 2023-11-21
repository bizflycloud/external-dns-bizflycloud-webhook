package dnsprovider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v8"

	"github.com/ionos-cloud/external-dns-ionos-webhook/cmd/webhook/init/configuration"
	"github.com/ionos-cloud/external-dns-ionos-webhook/internal/bizflycloud"
	"github.com/ionos-cloud/external-dns-ionos-webhook/pkg/endpoint"
	"github.com/ionos-cloud/external-dns-ionos-webhook/pkg/provider"
	log "github.com/sirupsen/logrus"
)

const (
	webtokenIonosISSValue = "ionoscloud"
)

func Init(config configuration.Config) (provider.Provider, error) {
	var domainFilter endpoint.DomainFilter
	createMsg := "Creating BIZFLYCLOUD provider with "

	if config.RegexDomainFilter != "" {
		createMsg += fmt.Sprintf("Regexp domain filter: '%s', ", config.RegexDomainFilter)
		if config.RegexDomainExclusion != "" {
			createMsg += fmt.Sprintf("with exclusion: '%s', ", config.RegexDomainExclusion)
		}
		domainFilter = endpoint.NewRegexDomainFilter(
			regexp.MustCompile(config.RegexDomainFilter),
			regexp.MustCompile(config.RegexDomainExclusion),
		)
	} else {
		if config.DomainFilter != nil && len(config.DomainFilter) > 0 {
			createMsg += fmt.Sprintf("zoneNode filter: '%s', ", strings.Join(config.DomainFilter, ","))
		}
		if config.ExcludeDomains != nil && len(config.ExcludeDomains) > 0 {
			createMsg += fmt.Sprintf("Exclude domain filter: '%s', ", strings.Join(config.ExcludeDomains, ","))
		}
		domainFilter = endpoint.NewDomainFilterWithExclusions(config.DomainFilter, config.ExcludeDomains)
	}

	createMsg = strings.TrimSuffix(createMsg, ", ")
	if strings.HasSuffix(createMsg, "with ") {
		createMsg += "no kind of domain filters"
	}
	log.Info(createMsg)
	bizflycloudConfig := bizflycloud.Configuration{}
	if err := env.Parse(&bizflycloudConfig); err != nil {
		return nil, fmt.Errorf("reading ionos ionosConfig failed: %v", err)
	}
	return bizflycloud.NewBizflyCloudProvider(domainFilter, &bizflycloudConfig)
}
