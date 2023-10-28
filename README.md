# DDNS
DDNS is a golang tool used to update dynamic DNS entries on supported dynamic DNS services.

## Command line flags
```sh
The DDNS CLI lets you interact with the DDNS service

Usage:
  ddns [flags]
  ddns [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Run A record synchronization once
  serve       Serve daemon that periodically performs A record synchronization

Flags:
      --config string     relative or absolute path to the config file (default "./config.yml")
  -h, --help              help for ddns
      --loglevel string   log level, possible values: trace, debug, info, warn, error, fatal, panic (default "info")
  -v, --version           version for ddns

Use "ddns [command] --help" for more information about a command.
```

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

Configuration parameters specified via environment variables take precedence over those specified in the config file.

## Global Configuration Parameters

| Key             | Env Var               | Type            | Default Value | Required | Description                                                    |
|-----------------|-----------------------|-----------------|---------------|----------|----------------------------------------------------------------|
| `waitInterval`  | `DDNS_WAIT_INTERVAL`  | `time.Duration` | `1m`          | `false`  | time.Duration to wait after successfully updating records      |
| `retryInterval` | `DDNS_RETRY_INTERVAL` | `time.Duration` | `5s`          | `false`  | time.Duration to wait after a failed attempt to update records |

## Metrics Server Configuration Parameters
Configuration Key: `metricsServer`

| Key      | Env Var               | Type     | Default Value | Required | Description                                          |
|----------|-----------------------|----------|---------------|----------|------------------------------------------------------|
| `enable` | `DDNS_METRICS_ENABLE` | `bool`   | `true`        | `true`   | Enable metrics endpoint listening on `/metrics` path |
| `host`   | `DDNS_METRICS_HOST`   | `string` | `0.0.0.0`     | `false`  | Host to be bound by the metrics handler              |
| `port`   | `DDNS_METRICS_PORT`   | `string` | `9097`        | `false`  | Port to be bound by the metrics handler              |

### Available Metrics
| Name                                    | Type    | Help                                                                                         |
|-----------------------------------------|---------|----------------------------------------------------------------------------------------------|
| `ddns_build_info`                       | `Gauge` | Metric with a constant '1' value labeled by version and goversion from which ddns was built. |
| `ddns_start_time_seconds`               | `Gauge` | Start time of the process since unix epoch in seconds.                                       |
| `ddns_dns_a_record_info`                | `Gauge` | Metric with a constant '1' value showing the current a records and their ip addresses.       |

## Available Providers for Retrieving the IP Address

### StaticIPAddressProvider
Ip address provider that returns the static ip address that is provided in the config file.

Configuration Key: `staticIPAddressProvider`

| Key       | Env Var                        | Type     | Default Value | Required | Description                 |
|-----------|--------------------------------|----------|---------------|----------|-----------------------------|
| `enable`  | `DDNS_STATIC_PROVIDER_ENABLE`  | `bool`   | `false`       | `true`   | Enable this provider        |
| `address` | `DDNS_STATIC_PROVIDER_ADDRESS` | `string` | `127.0.0.1`   | `false`  | Static ip address to return |

### URLIPAddressProvider
Ip address provider that makes a get request against the url that is provided in the config file and parses the response body using the regex if defined.

Configuration Key: `urlIPAddressProvider`

| Key                  | Env Var                      | Type     | Default Value | Required | Description                                                                                                  |
|----------------------|------------------------------|----------|---------------|----------|--------------------------------------------------------------------------------------------------------------|
| `enable`             | `DDNS_URL_PROVIDER_ENABLE`   | `bool`   | `false`       | `true`   | Enable this provider                                                                                         |
| `url`                | `DDNS_URL_PROVIDER_URL`      | `string` | `127.0.0.1`   | `false`  | URL to get the ip address from with a GET request                                                            |
| `https`              | `DDNS_URL_PROVIDER_HTTPS`    | `bool`   | `true`        | `false`  | Use https when accessing the url if true, http otherwise                                                     |
| `insecureSkipVerify` | `DDNS_URL_PROVIDER_INSECURE` | `bool`   | `false`       | `false`  | Ignore bad certificates when accessing the url                                                               |
| `regex`              | `DDNS_URL_PROVIDER_REGEX`    | `string` |               | `false`  | Regex to match the ip address containing a single numbered match group, see https://pkg.go.dev/regexp/syntax |
| `username`           | `DDNS_URL_PROVIDER_USERNAME` | `string` |               | `false`  | Basic auth username to use when accessing the url, only set if required                                      |
| `password`           | `DDNS_URL_PROVIDER_PASSWORD` | `string` |               | `false`  | Basic auth password to use when accessing the url, only set if required                                      |

For example, if the `website https://www.example.com/ipaddress` returned this json:
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

## Available DNS Providers

### CloudflareDNSProvider
Configuration Key: `cloudflareDNSProvider`

| Key        | Env Var                            | Type       | Default Value | Required | Description                                                            |
|------------|------------------------------------|------------|---------------|----------|------------------------------------------------------------------------|
| `enable`   | `DDNS_CLOUDFLARE_PROVIDER_ENABLE`  | `bool`     | `false`       | `true`   | Enable this provider                                                   |
| `apiToken` | `DDNS_CLOUDFLARE_API_TOKEN`        | `string`   |               | `true`   | Cloudflare API token with `All zones - DNS:Read, DNS:Edit` permissions |
| `zoneID`   | `DDNS_CLOUDFLARE_PROVIDER_ZONE_ID` | `string`   |               | `true`   | Cloudflare zone id                                                     |
| `aRecords` | `DDNS_CLOUDFLARE_PROVIDER_RECORDS` | `[]string` |               | `true`   | List of A records to update                                            |

## Build Docker Image
Docker image is available at [Docker Hub](https://hub.docker.com/r/mmianl/ddns).

```sh
export VERSION=`cat VERSION`
docker build . -t ddns:v${VERSION}
```

## Example Systemd Service File
This service file assumes that a user called `ddns` exists, and that the config file is located at `/etc/ddns/ddns.yaml`.
```sh
sudo groupadd ddns
sudo useradd -r -g ddns ddns
sudo mkdir /etc/ddns
# Write config file to /etc/ddns/ddns.yaml
sudo chmod 600 /etc/ddns/ddns.yaml
sudo chown ddns:ddns -R /etc/ddns/
```

```sh
[Unit]
Description=Dynamic DNS Client
After=network.target

[Service]
Type=simple
User=ddns
ExecStart=/usr/local/bin/ddns -config /etc/ddns/ddns.yaml

[Install]
WantedBy=multi-user.target
```
