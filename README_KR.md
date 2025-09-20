[<img src="https://developer.apple.com/assets/elements/badges/download-on-the-app-store.svg"
alt="Pushback App"
height="40">](https://apps.apple.com/us/app/id6615073345)
# NoLetServer

[中文](./README.md) | [English](./README_EN.md) | [日本語](./README_JP.md)

## 설치 및 실행

### GitHub Releases에서 다운로드

GitHub Releases 페이지에서 미리 컴파일된 바이너리를 다운로드할 수 있습니다:

1. [GitHub Releases](https://github.com/sunvc/NoLeterver/releases) 페이지 방문
2. 운영 체제 및 아키텍처에 맞는 버전 선택:
   - Windows (amd64)
   - macOS (amd64, arm64)
   - Linux (amd64, arm64, mips64, mips64le)
   - FreeBSD (amd64, arm64)
3. 다운로드한 파일 압축 해제
4. 구성 파일 생성(아래 구성 지침 참조)
5. 프로그램 실행:
   ```bash
   # Linux/macOS
   ./NoLetServer --config your_config.yaml
   
   # Windows
   NoLetServer.exe --config your_config.yaml
   ```

### Docker 사용

#### Docker 이미지

이 프로젝트는 다음 Docker 이미지 주소를 제공합니다:

- Docker Hub: `sunvc/nolet:latest`
- GitHub Container Registry: `ghcr.io/sunvc/nolet:latest`

다음 명령으로 이미지를 가져올 수 있습니다:

```bash
# Docker Hub에서 가져오기
docker pull sunvc/nolet:latest

# 또는 GitHub Container Registry에서 가져오기
docker pull ghcr.io/sunvc/nolet:latest

docker run -d --name NoLet-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  ghcr.io/sunvc/nolet:latest
```

#### Docker Compose 사용

프로젝트 루트 디렉토리의 `compose.yaml` 파일은 이미 Docker 이미지를 사용하도록 구성되어 있습니다:

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

다음 명령으로 서비스를 시작합니다:

```bash
docker-compose up -d
```

## 구성 파일

프로젝트의 `config.yaml`은 구성 파일의 예시일 뿐입니다. **사용자는 자신의 구성 파일을 생성하고 지정해야 합니다**. `--config` 또는 `-c` 매개변수를 사용하여 구성 파일 경로를 지정할 수 있습니다.

### 구성 파일 구조

```yaml
system:
  user: ""                  # 기본 인증 사용자 이름
  password: ""              # 기본 인증 비밀번호
  addr: "0.0.0.0:8080"      # 서버 리스닝 주소
  url_prefix: "/"           # 서비스 URL 접두사
  data: "./data"            # 데이터 저장 디렉토리
  name: "NoLetServer"            # 서비스 이름
  dsn: ""                   # MySQL DSN 연결 문자열
  cert: ""                  # TLS 인증서 경로
  key: ""                   # TLS 인증서 개인 키 경로
  reduce_memory_usage: false # 메모리 사용량 감소(CPU 사용량 증가)
  proxy_header: ""          # HTTP 헤더의 원격 IP 주소 소스
  max_batch_push_count: -1  # 배치 푸시 최대 수, -1은 무제한
  max_apns_client_count: 1  # APNs 클라이언트 연결 최대 수
  concurrency: 262144       # 최대 동시 연결 수(256 * 1024)
  read_timeout: 3s          # 읽기 타임아웃
  write_timeout: 3s         # 쓰기 타임아웃
  idle_timeout: 10s         # 유휴 타임아웃
  admins: []                # 관리자 ID 목록
  debug: true               # 디버그 모드 활성화
  expired: 0                # 음성 만료 시간(초)
  icp_info: ""              # ICP 등록 정보
  time_zone: "UTC"          # 시간대 설정

apple:
  apnsPrivateKey: ""        # APNs 개인 키 내용 또는 경로
  topic: ""                 # APNs Topic
  keyID: ""                 # APNs Key ID
  teamID: ""                # APNs Team ID
  develop: false            # APNs 개발 환경 활성화
```

## 서비스 구성 방법

서비스는 다음 3가지 방법으로 구성할 수 있으며, 우선순위는 높은 것부터 낮은 것 순입니다:

1. **명령줄 매개변수**: 시작 시 지정된 매개변수, 최우선
2. **환경 변수**: 시스템 환경 변수, 다음 우선순위
3. **구성 파일**: `config.yaml` 파일 또는 `--config`/`-c` 매개변수로 지정된 구성 파일

### 명령줄 매개변수 및 환경 변수

| 매개변수 | 환경 변수 | 설명 | 기본값 |
|------|----------|------|--------|
| `--addr` | `NOLET_SERVER_ADDRESS` | 서버 리스닝 주소 | `0.0.0.0:8080` |
| `--url-prefix` | `NOLET_SERVER_URL_PREFIX` | 서비스 URL 접두사 | `/` |
| `--dir` | `NOLET_SERVER_DATA_DIR` | 데이터 저장 디렉토리 | `./data` |
| `--dsn` | `NOLET_SERVER_DSN` | MySQL DSN, 형식: `user:pass@tcp(host)/dbname` | 비어 있음 |
| `--cert` | `NOLET_SERVER_CERT` | TLS 인증서 경로 | 비어 있음 |
| `--key` | `NOLET_SERVER_KEY` | TLS 인증서 개인 키 경로 | 비어 있음 |
| `--reduce-memory-usage` | `NOLET_SERVER_REDUCE_MEMORY_USAGE` | 메모리 사용량 감소(CPU 사용량 증가) | `false` |
| `--user, -u` | `NOLET_SERVER_BASIC_AUTH_USER` | 기본 인증 사용자 이름 | 비어 있음 |
| `--password, -p` | `NOLET_SERVER_BASIC_AUTH_PASSWORD` | 기본 인증 비밀번호 | 비어 있음 |
| `--proxy-header` | `NOLET_SERVER_PROXY_HEADER` | HTTP 헤더의 원격 IP 주소 소스 | 비어 있음 |
| `--max-batch-push-count` | `NOLET_SERVER_MAX_BATCH_PUSH_COUNT` | 배치 푸시 최대 수, `-1`은 무제한 | `-1` |
| `--max-apns-client-count` | `NOLET_SERVER_MAX_APNS_CLIENT_COUNT` | APNs 클라이언트 연결 최대 수 | `1` |
| `--admins` | `NOLET_SERVER_ADMINS` | 관리자 ID 목록 | 비어 있음 |
| `--debug` | `NOLET_DEBUG` | 디버그 모드 활성화 | `false` |
| `--apns-private-key` | `NOLET_APPLE_APNS_PRIVATE_KEY` | APNs 개인 키 경로 | 비어 있음 |
| `--topic` | `NOLET_APPLE_TOPIC` | APNs Topic | 비어 있음 |
| `--key-id` | `NOLET_APPLE_KEY_ID` | APNs Key ID | 비어 있음 |
| `--team-id` | `NOLET_APPLE_TEAM_ID` | APNs Team ID | 비어 있음 |
| `--develop, --dev` | `NOLET_APPLE_DEVELOP` | APNs 개발 환경 활성화 | `false` |
| `--Expired, --ex` | `NOLET_EXPIRED_TIME` | 음성 만료 시간(초) | `120` |
| `--help, -h` | - | 도움말 정보 표시 | - |
| `--config, -c` | - | 구성 파일 경로 지정 | - |

### 구성 파일 사용

1. 자신의 구성 파일 생성:
   - 프로젝트의 `config.yaml` 예시를 참조하여 자신의 구성 파일 생성
   - 구성 파일에 필요한 구성 항목이 포함되어 있는지 확인

2. 구성 파일 경로 지정:
   ```bash
   ./NoLetServer --config /path/to/your/config.yaml
   # 또는 약식 사용
   ./NoLetServer -c /path/to/your/config.yaml
   ```

3. 구성 파일과 명령줄 매개변수 혼합 사용:
   ```bash
   # 구성 파일의 설정은 명령줄 매개변수에 의해 재정의됩니다
   ./NoLetServer -c /path/to/your/config.yaml --debug --addr 127.0.0.1:8080
   ```