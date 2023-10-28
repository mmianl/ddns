package internal

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"

	"github.com/rs/zerolog/log"
)

type URLIPAddressProviderConfig struct {
	// Switch to enable or disable this provider
	Enable bool `yaml:"enable" envconfig:"DDNS_URL_PROVIDER_ENABLE" required:"false"`

	// URL to get the ip address from excluding the protocol
	URL string `yaml:"url" envconfig:"DDNS_URL_PROVIDER_URL" required:"false"`

	// Switch to turn on or off https, will be http if off
	HTTPS bool `yaml:"https" envconfig:"DDNS_URL_PROVIDER_HTTPS" required:"false"`

	// Switch to turn on or off https, will be http if off
	InsecureSkipVerify bool `yaml:"insecureSkipVerify" envconfig:"DDNS_URL_PROVIDER_INSECURE" required:"false"`

	// Regex containing a single numbered match group, see https://pkg.go.dev/regexp/syntax
	Regex string `yaml:"regex" envconfig:"DDNS_URL_PROVIDER_REGEX" required:"false"`

	// BasicAuth username
	Username string `yaml:"username" envconfig:"DDNS_URL_PROVIDER_USERNAME" required:"false"`

	// BasicAuth password
	Password string `yaml:"password" envconfig:"DDNS_URL_PROVIDER_PASSWORD" required:"false"`
}

var defaultURLIPAddressProviderConfig = &URLIPAddressProviderConfig{
	Enable:             false,
	URL:                "127.0.0.1",
	HTTPS:              true,
	InsecureSkipVerify: false,
	Regex:              "",
	Username:           "",
	Password:           "",
}

type URLIPAddressProvider struct {
	url                string
	https              bool
	insecureSkipVerify bool
	regex              string
	username           string
	password           string
}

// NewURLIPAddressProvider Returns an instance of URLIPAddressProvider based on the passed configuration
func NewURLIPAddressProvider(config *URLIPAddressProviderConfig) *URLIPAddressProvider {
	return &URLIPAddressProvider{
		url:                config.URL,
		https:              config.HTTPS,
		insecureSkipVerify: config.InsecureSkipVerify,
		regex:              config.Regex,
		username:           config.Username,
		password:           config.Password,
	}
}

// GetIPAddress Returns the ip address returned by the url, and parsed using the regex provided via the configuration
func (u *URLIPAddressProvider) GetIPAddress() (*string, error) {
	// Make request
	proto := "http"
	if u.https {
		proto = "https"
	}
	requestURL := fmt.Sprintf("%s://%s", proto, u.url)

	log.Debug().Msgf("Setting InsecureSkipVerify to %v", u.insecureSkipVerify)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: u.insecureSkipVerify},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	if u.username != "" && u.password != "" {
		log.Debug().Msg("Setting basic auth credentials")
		req.SetBasicAuth(u.username, u.password)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Msgf("error %s occurred while closing response body", err)
		}
	}(res.Body)

	// Parse body
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code from %s was %s, not 200", requestURL, res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(bodyBytes)

	// Apply regex if configured
	if u.regex == "" {
		addr := &bodyString
		return addr, nil
	}

	addr, err := GetRegexSubstring(u.regex, bodyString)
	if err != nil {
		return nil, err
	}

	validated := net.ParseIP(*addr)
	if validated == nil {
		return nil, fmt.Errorf("did not get a valid ip address, got %s", *addr)
	}
	return addr, nil
}

// GetRegexSubstring Returns the first numbered match group, or an error if there are no matches or if there are more than 1
func GetRegexSubstring(regex string, s string) (*string, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	r := re.FindStringSubmatch(s)
	if len(r) != 2 {
		return nil, fmt.Errorf("unexpected result when applying regex to %s", s)
	}

	return &r[1], nil
}
