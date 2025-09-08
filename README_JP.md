# NoLetServer

[中文](./README.md) | [English](./README_EN.md) | [한국어](./README_KR.md)

## インストールと実行

### GitHub Releasesからダウンロード

GitHub Releasesページからプリコンパイルされたバイナリをダウンロードできます：

1. [GitHub Releases](https://github.com/uuneo/NoLetServer/releases) ページにアクセス
2. お使いのオペレーティングシステムとアーキテクチャに適したバージョンを選択：
   - Windows (amd64)
   - macOS (amd64, arm64)
   - Linux (amd64, arm64, mips64, mips64le)
   - FreeBSD (amd64, arm64)
3. ダウンロードしたファイルを解凍
4. 設定ファイルを作成（以下の設定手順を参照）
5. プログラムを実行：
   ```bash
   # Linux/macOS
   ./NoLetServer --config your_config.yaml
   
   # Windows
   NoLetServer.exe --config your_config.yaml
   ```

### Dockerの使用

#### Dockerイメージ

このプロジェクトでは、以下のDockerイメージアドレスを提供しています：

- Docker Hub: `sunvc/nolet:latest`
- GitHub Container Registry: `ghcr.io/sunvc/nolet:latest`

以下のコマンドでイメージをプルできます：

```bash
# Docker Hubからプル
docker pull sunvc/nolet:latest

# または、GitHub Container Registryからプル
docker pull ghcr.io/sunvc/nolet:latest

docker run -d --name NoLet-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  ghcr.io/sunvc/nolet:latest
```

#### Docker Composeの使用

プロジェクトのルートディレクトリにある`compose.yaml`ファイルは、Dockerイメージを使用するように既に設定されています：

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

以下のコマンドでサービスを起動します：

```bash
docker-compose up -d
```

## 設定ファイル

プロジェクト内の`config.yaml`は設定ファイルの例にすぎません。**ユーザーは自分で設定ファイルを作成して指定する必要があります**。`--config`または`-c`パラメータを使用して設定ファイルのパスを指定できます。

### 設定ファイルの構造

```yaml
system:
  user: ""                  # 基本認証ユーザー名
  password: ""              # 基本認証パスワード
  addr: "0.0.0.0:8080"      # サーバーリスニングアドレス
  url_prefix: "/"           # サービスURLプレフィックス
  data: "./data"            # データストレージディレクトリ
  name: "NoLetServer"            # サービス名
  dsn: ""                   # MySQL DSN接続文字列
  cert: ""                  # TLS証明書パス
  key: ""                   # TLS証明書秘密鍵パス
  reduce_memory_usage: false # メモリ使用量を削減（CPU消費量が増加）
  proxy_header: ""          # HTTPヘッダーのリモートIPアドレスソース
  max_batch_push_count: -1  # バッチプッシュの最大数、-1は無制限
  max_apns_client_count: 1  # APNsクライアント接続の最大数
  concurrency: 262144       # 最大同時接続数（256 * 1024）
  read_timeout: 3s          # 読み取りタイムアウト
  write_timeout: 3s         # 書き込みタイムアウト
  idle_timeout: 10s         # アイドルタイムアウト
  admins: []                # 管理者IDリスト
  debug: true               # デバッグモードを有効にする
  expired: 0                # 音声の有効期限（秒）
  icp_info: ""              # ICP登録情報
  time_zone: "UTC"          # タイムゾーン設定

apple:
  apnsPrivateKey: ""        # APNs秘密鍵の内容またはパス
  topic: ""                 # APNs Topic
  keyID: ""                 # APNs Key ID
  teamID: ""                # APNs Team ID
  develop: false            # APNs開発環境を有効にする
```

## サービス設定方法

サービスは以下の3つの方法で設定でき、優先順位は高いものから低いものへ：

1. **コマンドラインパラメータ**：起動時に指定されるパラメータ、最優先
2. **環境変数**：システム環境変数、次の優先順位
3. **設定ファイル**：`config.yaml`ファイルまたは`--config`/`-c`パラメータで指定された設定ファイル

### コマンドラインパラメータと環境変数

| パラメータ | 環境変数 | 説明 | デフォルト値 |
|------|----------|------|--------|
| `--addr` | `NOLET_SERVER_ADDRESS` | サーバーリスニングアドレス | `0.0.0.0:8080` |
| `--url-prefix` | `NOLET_SERVER_URL_PREFIX` | サービスURLプレフィックス | `/` |
| `--dir` | `NOLET_SERVER_DATA_DIR` | データストレージディレクトリ | `./data` |
| `--dsn` | `NOLET_SERVER_DSN` | MySQL DSN、形式：`user:pass@tcp(host)/dbname` | 空 |
| `--cert` | `NOLET_SERVER_CERT` | TLS証明書パス | 空 |
| `--key` | `NOLET_SERVER_KEY` | TLS証明書秘密鍵パス | 空 |
| `--reduce-memory-usage` | `NOLET_SERVER_REDUCE_MEMORY_USAGE` | メモリ使用量を削減（CPU消費量が増加） | `false` |
| `--user, -u` | `NOLET_SERVER_BASIC_AUTH_USER` | 基本認証ユーザー名 | 空 |
| `--password, -p` | `NOLET_SERVER_BASIC_AUTH_PASSWORD` | 基本認証パスワード | 空 |
| `--proxy-header` | `NOLET_SERVER_PROXY_HEADER` | HTTPヘッダーのリモートIPアドレスソース | 空 |
| `--max-batch-push-count` | `NOLET_SERVER_MAX_BATCH_PUSH_COUNT` | バッチプッシュの最大数、`-1`は無制限 | `-1` |
| `--max-apns-client-count` | `NOLET_SERVER_MAX_APNS_CLIENT_COUNT` | APNsクライアント接続の最大数 | `1` |
| `--admins` | `NOLET_SERVER_ADMINS` | 管理者IDリスト | 空 |
| `--debug` | `NOLET_DEBUG` | デバッグモードを有効にする | `false` |
| `--apns-private-key` | `NOLET_APPLE_APNS_PRIVATE_KEY` | APNs秘密鍵パス | 空 |
| `--topic` | `NOLET_APPLE_TOPIC` | APNs Topic | 空 |
| `--key-id` | `NOLET_APPLE_KEY_ID` | APNs Key ID | 空 |
| `--team-id` | `NOLET_APPLE_TEAM_ID` | APNs Team ID | 空 |
| `--develop, --dev` | `NOLET_APPLE_DEVELOP` | APNs開発環境を有効にする | `false` |
| `--Expired, --ex` | `NOLET_EXPIRED_TIME` | 音声の有効期限（秒） | `120` |
| `--help, -h` | - | ヘルプ情報を表示 | - |
| `--config, -c` | - | 設定ファイルパスを指定 | - |

### 設定ファイルの使用

1. 自分の設定ファイルを作成：
   - プロジェクト内の`config.yaml`例を参考に自分の設定ファイルを作成
   - 設定ファイルに必要な設定項目が含まれていることを確認

2. 設定ファイルパスを指定：
   ```bash
   ./NoLetServer --config /path/to/your/config.yaml
   # または省略形を使用
   ./NoLetServer -c /path/to/your/config.yaml
   ```

3. 設定ファイルとコマンドラインパラメータの混合使用：
   ```bash
   # 設定ファイル内の設定はコマンドラインパラメータによって上書きされます
   ./NoLetServer -c /path/to/your/config.yaml --debug --addr 127.0.0.1:8080
   ```