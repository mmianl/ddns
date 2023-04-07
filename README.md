# ddns
ddns is a golang daemon used to update dynamic DNS entries on supported dynamic DNS services.

## Command line flags
* `-help`: shows help
* `-config`: Relative or absolute path to the config file (default "./config.yml")
* `-logLevel`: Log level, possible values: trace, debug, info, warn, error, fatal, panic (default "info")

## Example Config File
```yaml
waitInterval: "1m"
retryInterval: "10s"

metricsServer:
  enable: true
  host: 127.0.0.1
  port: 8080

staticIPAddressProvider:
  enable: false
  address: "10.0.0.1"

urlIPAddressProvider:
  enable: true
  url: "www.example.com/ipaddress"
  https: true
  insecureSkipVerify: false
  regex: ""
  username: "username"
  password: "password"

cloudflareDNSProvider:
  enable: true
  apiToken: "12345"
  zoneID: "12345"
  aRecords:
    - "example.com"
    - "www.example.com"
```

If multiple DNS or IP Address providers are specified in the config file, only one will take effect, The order of precedence is the order in which the providers are listed below, with the first provider having the highest priority.

## Global Configuration Parameters

| Key             | Type            | Default Value | Required | Description |
|-----------------|-----------------|---------------|----------|-------------|
| `waitInterval`  | `time.Duration` | `1m`          | `false`  |             |
| `retryInterval` | `time.Duration` | `5s`          | `false`  |             |

## Metrics Server Configuration Parameters
Configuration Key: `metricsServer`

| Key      | Type     | Default Value | Required | Description                                          |
|----------|----------|---------------|----------|------------------------------------------------------|
| `enable` | `bool`   | `true`        | `true`   | Enable metrics endpoint listening on `/metrics` path |
| `host`   | `string` | `0.0.0.0`     | `false`  | Host to be bound by the metrics handler              |
| `port`   | `string` | `9097`        | `false`  | Port to be bound by the metrics handler              |

### Available Metrics
| Name                                    | Type    | Help                                                                                             |
|-----------------------------------------|---------|--------------------------------------------------------------------------------------------------|
| `ddns_build_info`                       | `Gauge` | Metric with a constant '1' value labeled by version and goversion from which ddns was built.     |
| `ddns_start_time_seconds`               | `Gauge` | Start time of the process since unix epoch in seconds.                                           |
| `ddns_dns_a_record_update_time_seconds` | `Gauge` | Time of last DNS A record update since unix epoch in seconds labeled by IP address and A Record. |

## Available Providers for Retrieving the IP Address

### StaticIPAddressProvider
Ip address provider that returns a the static ip address that is provided in the config file.

Configuration Key: `staticIPAddressProvider`

| Key                  | Type     | Default Value | Required | Description                 |
|----------------------|----------|---------------|----------|-----------------------------|
| `enable`             | `bool`   | `false`       | `true`   | Enable this provider        |
| `address`            | `string` | `127.0.0.1`   | `false`  | Static ip address to return |

### URLIPAddressProvider
Ip address provider that makes a get request against the url that is provided in the config file and parses the repsonse body using the regex if defined.

Configuration Key: `urlIPAddressProvider`

| Key                  | Type     | Default Value | Required | Description                                                                                                  |
|----------------------|----------|---------------|----------|--------------------------------------------------------------------------------------------------------------|
| `enable`             | `bool`   | `false`       | `true`   | Enable this provider                                                                                         |
| `url`                | `string` | `127.0.0.1`   | `false`  | URL to get the ip address from with a GET request                                                            |
| `https`              | `bool`   | `true`        | `false`  | Use https when accessing the url if true, http otherwise                                                     |
| `insecureSkipVerify` | `bool`   | `false`       | `false`  | Ignore bad certificates when accessing the url                                                               |
| `regex`              | `string` |               | `false`  | Regex to match the ip address containing a single numbered match group, see https://pkg.go.dev/regexp/syntax |
| `username`           | `string` |               | `false`  | Basic auth username to use when accessing the url, only set if required                                      |
| `password`           | `string` |               | `false`  | Basic auth password to use when accessing the url, only set if required                                      |

For example, if the `website https://www.example.com/ipaddress` retruned this json:
```json
{"address":"192.168.0.100"}
```

Then the following configuration could be used to get the ip address from there:
```yaml
urlIPAddressProvider:
  enable: true
  url: "www.example.com/ipaddress"
  regex: '"address":\s?"(.*)"'
```

## Availbe DNS Providers

### CloudflareDNSProvider
Configuration Key: `cloudflareDNSProvider`

| Key        | Type       | Default Value | Required | Description                                                            |
|------------|------------|---------------|----------|------------------------------------------------------------------------|
| `enable`   | `bool`     | `false`       | `true`   | Enable this provider                                                   |
| `apiToken` | `string`   |               | `true`   | Cloudflare API token with `All zones - DNS:Read, DNS:Edit` permissions |
| `zoneID`   | `string`   |               | `true`   | Cloudflare zone id                                                     |
| `aRecords` | `[]string` |               | `true`   | List of A records to update                                            |

## Build Docker Image
Docker image is available at [Docker Hub](https://hub.docker.com/repository/docker/mmianl/ddns/general).

```sh
export VERSION=`cat VERSION`
docker build . -t mmianl/ddns:v${VERSION}
```
