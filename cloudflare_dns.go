package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type CloudflareDNSProviderConfig struct {
	// Swtich to enable or disable this provider
	Enable bool `yaml:"enable"`

	// Cloudflare API Token with "All zones - DNS:Read, DNS:Edit" permissions
	APIToken string `yaml:"apiToken"`

	// Cloudfare Zone ID
	ZoneID string `yaml:"zoneID"`

	// List of A Records
	ARecords []string `yaml:"aRecords"`
}

type cloudfareListRecordsResponse struct {
	Result []cloudfareListRecordsResponseResult `json:"result"`
}

type cloudfareListRecordsResponseResult struct {
	ID         string                                 `json:"id"`
	ZoneID     string                                 `json:"zone_id"`
	ZoneName   string                                 `json:"zone_name"`
	Name       string                                 `json:"name"`
	Type       string                                 `json:"type"`
	Content    string                                 `json:"content"`
	Proxiable  bool                                   `json:"proxiable"`
	Proxied    bool                                   `json:"proxied"`
	TTL        int64                                  `json:"ttl"`
	Locked     bool                                   `json:"locked"`
	Meta       cloudfareListRecordsResponseResultMeta `json:"meta"`
	Comment    string                                 `json:"comment"`
	Tags       []string                               `json:"tags"`
	CreatedOn  time.Time                              `json:"created_on"`
	ModifiedOn time.Time                              `json:"modified_on"`
}

type cloudfareListRecordsResponseResultMeta struct {
	AutoAdded           bool   `json:"auto_added"`
	ManagedByApps       bool   `json:"managed_by_apps"`
	ManagedByArgoTunnel bool   `json:"managed_by_argo_tunnel"`
	Source              string `json:"source"`
}

type cloudfareUpdateRecordPayload struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

var defaultCloudflareDNSProviderConfig = &CloudflareDNSProviderConfig{
	Enable:   false,
	APIToken: "",
	ZoneID:   "",
	ARecords: nil,
}

type CloudflareDNSProvider struct {
	apiToken string
	zoneID   string
	aRecords []string
}

// Returns an instance of CloudflareDNSProvider based on the passed configuration
func NewCloudflareDNSProvider(config *CloudflareDNSProviderConfig) *CloudflareDNSProvider {
	return &CloudflareDNSProvider{
		apiToken: config.APIToken,
		zoneID:   config.ZoneID,
		aRecords: config.ARecords,
	}
}

// Return the RecordAddressMappings with the current ip addresses for names specified in the configuration
func (c *CloudflareDNSProvider) GetARecordAddresses() ([]RecordAddressMapping, error) {
	// Make request
	requestURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", c.zoneID)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse body
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code was %s, not 200", res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(bodyBytes)

	var r cloudfareListRecordsResponse
	if err := json.Unmarshal([]byte(bodyString), &r); err != nil {
		return nil, err
	}

	return GetContentsByNames(r, c.aRecords)
}

// Return the RecordAddressMappings with the current ip addresses for the provided names
func GetContentsByNames(response cloudfareListRecordsResponse, names []string) ([]RecordAddressMapping, error) {
	var ms []RecordAddressMapping

	for _, name := range names {

		m, err := GetContentByName(response, name)
		if err != nil {
			return nil, err
		}

		ms = append(ms, *m)
	}

	if len(ms) != len(names) {
		return nil, errors.New("not all a records could not found")
	}

	return ms, nil
}

// Return the RecordAddressMapping with the current ip address for the provided name
func GetContentByName(response cloudfareListRecordsResponse, name string) (*RecordAddressMapping, error) {
	for _, item := range response.Result {
		if item.Name == name {
			return &RecordAddressMapping{id: item.ID, aRecord: name, ipAddress: item.Content}, nil
		}
	}

	return nil, fmt.Errorf("no record with name %s was found", name)
}

// Set the provided A record to the provided ip address
func (c *CloudflareDNSProvider) SetARecordAddress(ipAddress string, m RecordAddressMapping) error {
	log.Info().Msgf("Setting A record %s to %s", m.aRecord, ipAddress)

	// Make request
	requestURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", c.zoneID, m.id)
	payload := &cloudfareUpdateRecordPayload{
		Name:    m.aRecord,
		Type:    "A",
		Content: ipAddress,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Executing put request against %s with payload %s", requestURL, json)
	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Parse body
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code from %s was %s, not 200", requestURL, res.Status)
	}

	return nil
}
