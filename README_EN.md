[<img src="https://developer.apple.com/assets/elements/badges/download-on-the-app-store.svg"
alt="Pushback App"
height="40">](https://apps.apple.com/us/app/id6615073345)
# NoLetServer

[中文](./README.md) | [日本語](./README_JP.md) | [한국어](./README_KR.md)

## Installation and Running

### Download from GitHub Releases

You can download pre-compiled binaries from the GitHub Releases page:

1. Visit the [GitHub Releases](https://github.com/sunvc/NoLeterver/releases) page
2. Choose the appropriate version for your operating system and architecture:
   - Windows (amd64)
   - macOS (amd64, arm64)
   - Linux (amd64, arm64, mips64, mips64le)
   - FreeBSD (amd64, arm64)
3. Extract the downloaded file
4. Create a configuration file (refer to the configuration instructions below)
5. Run the program:
   ```bash
   # Linux/macOS
   ./NoLetServer --config your_config.yaml
   
   # Windows
   NoLetServer.exe --config your_config.yaml
   ```

   Common parameters:
   - `--addr`: Server listening address, default is 0.0.0.0:8080
   - `--url-prefix`: Service URL prefix, default is /
   - `--dir`: Data storage directory, default is ./data
   - `--dsn`: MySQL database connection string
   - `--debug`: Enable debug mode
   - `--config, -c`: Specify configuration file path

### Using Docker

#### Docker Image

This project provides the following Docker image addresses:

- Docker Hub: `sunvc/nolet:latest`
- GitHub Container Registry: `ghcr.io/sunvc/nolet:latest`

You can pull the image using the following command:

```bash
# Pull from Docker Hub
docker pull sunvc/nolet:latest

# Or pull from GitHub Container Registry
docker pull ghcr.io/sunvc/nolet:latest

docker run -d --name NoLet-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  ghcr.io/sunvc/nolet:latest
```

#### Using Docker Compose

The `compose.yaml` file in the project root directory is already configured to use the Docker image:

```yaml
services:
  NoLetServer:
    image: ghcr.io/sunvc/nolet:latest
    container_name: NoLetServer
    restart: always
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
```

Run the following command to start the service:

```bash
docker-compose up -d
```

## Configuration File

The `config.yaml` in the project is only a configuration file example. **Users need to create and specify their own configuration file** for service configuration. You can use the `--config` or `-c` parameter to specify the configuration file path.

### Configuration File Structure

```yaml
system:
  user: ""                  # Basic authentication username
  password: ""              # Basic authentication password
  addr: "0.0.0.0:8080"      # Server listening address
  url_prefix: "/"           # Service URL prefix
  data: "./data"            # Data storage directory
  name: "NoLetServer"            # Service name
  dsn: ""                   # MySQL DSN connection string
  cert: ""                  # TLS certificate path
  key: ""                   # TLS certificate private key path
  reduce_memory_usage: false # Reduce memory usage (increases CPU consumption)
  proxy_header: ""          # Remote IP address source in HTTP header
  max_batch_push_count: -1  # Maximum number of batch pushes, -1 means no limit
  max_apns_client_count: 1  # Maximum number of APNs client connections
  concurrency: 262144       # Maximum number of concurrent connections (256 * 1024)
  read_timeout: 3s          # Read timeout
  write_timeout: 3s         # Write timeout
  idle_timeout: 10s         # Idle timeout
  admins: []                # Administrator ID list
  debug: true               # Enable debug mode
  expired: 0                # Voice expiration time (seconds)
  icp_info: ""              # ICP filing information
  time_zone: "UTC"          # Time zone setting

apple:
  apnsPrivateKey: ""        # APNs private key content or path
  topic: ""                 # APNs Topic
  keyID: ""                 # APNs Key ID
  teamID: ""                # APNs Team ID
  develop: false            # Enable APNs development environment
```

## Service Configuration Methods

The service can be configured in the following three ways, with priority from high to low:

1. **Command-line parameters**: Parameters specified at startup, highest priority
2. **Environment variables**: System environment variables, second priority
3. **Configuration file**: `config.yaml` file or configuration file specified via `--config`/`-c` parameter

### Command-line Parameters and Environment Variables

| Parameter | Environment Variable | Description | Default Value |
|------|----------|------|--------|
| `--addr` | `NOLET_SERVER_ADDRESS` | Server listening address | `0.0.0.0:8080` |
| `--url-prefix` | `NOLET_SERVER_URL_PREFIX` | Service URL prefix | `/` |
| `--dir` | `NOLET_SERVER_DATA_DIR` | Data storage directory | `./data` |
| `--dsn` | `NOLET_SERVER_DSN` | MySQL DSN, format: `user:pass@tcp(host)/dbname` | Empty |
| `--cert` | `NOLET_SERVER_CERT` | TLS certificate path | Empty |
| `--key` | `NOLET_SERVER_KEY` | TLS certificate private key path | Empty |
| `--reduce-memory-usage` | `NOLET_SERVER_REDUCE_MEMORY_USAGE` | Reduce memory usage (increases CPU consumption) | `false` |
| `--user, -u` | `NOLET_SERVER_BASIC_AUTH_USER` | Basic authentication username | Empty |
| `--password, -p` | `NOLET_SERVER_BASIC_AUTH_PASSWORD` | Basic authentication password | Empty |
| `--proxy-header` | `NOLET_SERVER_PROXY_HEADER` | Remote IP address source in HTTP header | Empty |
| `--max-batch-push-count` | `NOLET_SERVER_MAX_BATCH_PUSH_COUNT` | Maximum number of batch pushes, `-1` means no limit | `-1` |
| `--max-apns-client-count` | `NOLET_SERVER_MAX_APNS_CLIENT_COUNT` | Maximum number of APNs client connections | `1` |
| `--admins` | `NOLET_SERVER_ADMINS` | Administrator ID list | Empty |
| `--debug` | `NOLET_DEBUG` | Enable debug mode | `false` |
| `--apns-private-key` | `NOLET_APPLE_APNS_PRIVATE_KEY` | APNs private key path | Empty |
| `--topic` | `NOLET_APPLE_TOPIC` | APNs Topic | Empty |
| `--key-id` | `NOLET_APPLE_KEY_ID` | APNs Key ID | Empty |
| `--team-id` | `NOLET_APPLE_TEAM_ID` | APNs Team ID | Empty |
| `--develop, --dev` | `NOLET_APPLE_DEVELOP` | Enable APNs development environment | `false` |
| `--Expired, --ex` | `NOLET_EXPIRED_TIME` | Voice expiration time (seconds) | `120` |
| `--help, -h` | - | Display help information | - |
| `--config, -c` | - | Specify configuration file path | - |

### Using Configuration File

1. Create your own configuration file:
   - Create your own configuration file referring to the `config.yaml` example in the project
   - Ensure the configuration file contains the required configuration items

2. Specify the configuration file path:
   ```bash
   ./NoLetServer --config /path/to/your/config.yaml
   # Or use the shorthand
   ./NoLetServer -c /path/to/your/config.yaml
   ```

3. Mixed use of configuration file and command-line parameters:
   ```bash
   # Settings in the configuration file will be overridden by command-line parameters
   ./NoLetServer -c /path/to/your/config.yaml --debug --addr 127.0.0.1:8080
   ```