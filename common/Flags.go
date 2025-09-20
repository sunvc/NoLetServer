package common

import (
	"context"
	"time"

	"github.com/urfave/cli/v3"
)

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "addr",
			Usage:   "Server listen address",
			Sources: cli.EnvVars("NOLET_SERVER_ADDRESS"),
			Value:   "0.0.0.0:8080",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.Addr = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "url-prefix",
			Usage:   "Serve URL Prefix",
			Sources: cli.EnvVars("NOLET_SERVER_URL_PREFIX"),
			Value:   "/",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.URLPrefix = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "dir",
			Usage:   "Server data storage dir",
			Sources: cli.EnvVars("NOLET_SERVER_DATA_DIR"),
			Value:   "./data",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.DataDir = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "dsn",
			Usage:   "MySQL DSN user:pass@tcp(host)/dbname",
			Sources: cli.EnvVars("NOLET_SERVER_DSN"),
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.DSN = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "cert",
			Usage:   "Server TLS certificate",
			Sources: cli.EnvVars("NOLET_SERVER_CERT"),
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.Cert = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "key",
			Usage:   "Server TLS certificate key",
			Sources: cli.EnvVars("NOLET_SERVER_KEY"),
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.Key = s
				return nil
			},
		},
		&cli.BoolFlag{
			Name:    "reduce-memory-usage",
			Usage:   "Aggressively reduces memory usage at the cost of higher CPU usage if set to true",
			Sources: cli.EnvVars("NOLET_SERVER_REDUCE_MEMORY_USAGE"),
			Value:   false,
			Action: func(ctx context.Context, command *cli.Command, b bool) error {
				LocalConfig.System.ReduceMemoryUsage = b
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "user",
			Usage:   "Basic auth username",
			Sources: cli.EnvVars("NOLET_SERVER_BASIC_AUTH_USER"),
			Aliases: []string{"u"},
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.User = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "password",
			Usage:   "Basic auth password",
			Sources: cli.EnvVars("NOLET_SERVER_BASIC_AUTH_PASSWORD"),
			Aliases: []string{"p"},
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.Password = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "proxy-header",
			Usage:   "The remote IP address used by the NOLET server http header",
			Sources: cli.EnvVars("NOLET_SERVER_PROXY_HEADER"),
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.ProxyHeader = s
				return nil
			},
		},
		&cli.IntFlag{
			Name:    "max-batch-push-count",
			Usage:   "Maximum number of batch pushes allowed, -1 means no limit",
			Sources: cli.EnvVars("NOLET_SERVER_MAX_BATCH_PUSH_COUNT"),
			Value:   -1,
			Action: func(ctx context.Context, command *cli.Command, v int) error {
				LocalConfig.System.MaxBatchPushCount = v
				return nil
			},
		},
		&cli.IntFlag{
			Name:    "max-apns-client-count",
			Usage:   "Maximum number of APNs client connections",
			Sources: cli.EnvVars("NOLET_SERVER_MAX_APNS_CLIENT_COUNT"),
			Value:   1,
			Action: func(ctx context.Context, command *cli.Command, v int) error {
				LocalConfig.System.MaxAPNSClientCount = v
				return nil
			},
		},
		&cli.IntFlag{
			Name:    "concurrency",
			Usage:   "Maximum number of concurrent connections",
			Sources: cli.EnvVars("NOLET_SERVER_CONCURRENCY"),
			Value:   256 * 1024,
			Hidden:  true,
			Action: func(ctx context.Context, command *cli.Command, b int) error {
				LocalConfig.System.Concurrency = b
				return nil
			},
		},
		&cli.DurationFlag{
			Name:    "read-timeout",
			Usage:   "The amount of time allowed to read the full request, including the body",
			Sources: cli.EnvVars("NOLET_SERVER_READ_TIMEOUT"),
			Value:   3 * time.Second,
			Hidden:  true,
			Action: func(ctx context.Context, command *cli.Command, duration time.Duration) error {
				LocalConfig.System.ReadTimeout = duration
				return nil
			},
		},
		&cli.DurationFlag{
			Name:    "write-timeout",
			Usage:   "The maximum duration before timing out writes of the response",
			Sources: cli.EnvVars("NOLET_SERVER_WRITE_TIMEOUT"),
			Value:   3 * time.Second,
			Hidden:  true,
			Action: func(ctx context.Context, command *cli.Command, duration time.Duration) error {
				LocalConfig.System.WriteTimeout = duration
				return nil
			},
		},
		&cli.DurationFlag{
			Name:    "idle-timeout",
			Usage:   "The maximum amount of time to wait for the next request when keep-alive is enabled",
			Sources: cli.EnvVars("NOLET_SERVER_IDLE_TIMEOUT"),
			Value:   10 * time.Second,
			Hidden:  true,
			Action: func(ctx context.Context, command *cli.Command, duration time.Duration) error {
				LocalConfig.System.IdleTimeout = duration
				return nil
			},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Value:   false,
			Usage:   "enable debug mode",
			Sources: cli.EnvVars("NOLET_DEBUG"),
			Action: func(ctx context.Context, command *cli.Command, b bool) error {
				LocalConfig.System.Debug = b
				return nil
			},
		},
		&cli.BoolFlag{
			Name:    "voice",
			Value:   false,
			Usage:   "Support voice",
			Sources: cli.EnvVars("NOLET_VOICE"),
			Hidden:  true,
			Action: func(ctx context.Context, command *cli.Command, b bool) error {
				LocalConfig.System.Voice = b
				return nil
			},
		},
		&cli.StringSliceFlag{
			Name:    "auths",
			Value:   []string{},
			Usage:   "auth id list",
			Sources: cli.EnvVars("NOLET_AUTHS"),
			Action: func(ctx context.Context, command *cli.Command, strings []string) error {
				LocalConfig.System.Auths = strings
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "apns-private-key",
			Usage:   "APNs private key path",
			Sources: cli.EnvVars("NOLET_APPLE_APNS_PRIVATE_KEY"),
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.Apple.ApnsPrivateKey = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "topic",
			Usage:   "APNs topic",
			Sources: cli.EnvVars("NOLET_APPLE_TOPIC"),
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.Apple.Topic = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "key-id",
			Usage:   "APNs key ID",
			Sources: cli.EnvVars("NOLET_APPLE_KEY_ID"),
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.Apple.KeyID = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "team-id",
			Usage:   "APNs team ID",
			Sources: cli.EnvVars("NOLET_APPLE_TEAM_ID"),
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.Apple.TeamID = s
				return nil
			},
		},
		&cli.BoolFlag{
			Name:    "develop",
			Usage:   "Use APNs development environment",
			Sources: cli.EnvVars("NOLET_APPLE_DEVELOP"),
			Aliases: []string{"dev"},
			Value:   false,
			Action: func(ctx context.Context, command *cli.Command, b bool) error {
				LocalConfig.Apple.Develop = b
				return nil
			},
		},
		&cli.Float64Flag{
			Name:    "Expired",
			Usage:   "Voice Expired Time",
			Sources: cli.EnvVars("NOLET_EXPIRED_TIME"),
			Aliases: []string{"ex"},
			Value:   60 * 2,
			Action: func(ctx context.Context, command *cli.Command, f float64) error {
				LocalConfig.System.Expired = f
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "ICP",
			Usage:   "Icp Footer Info",
			Sources: cli.EnvVars("NOLET_ICP_Info"),
			Aliases: []string{"icp"},
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.ICPInfo = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "timeZone",
			Usage:   "Time Zone",
			Aliases: []string{"tz"},
			Value:   "",
			Action: func(ctx context.Context, command *cli.Command, s string) error {
				LocalConfig.System.TimeZone = s
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Config file Dir",
			Aliases: []string{"c"},
			Value:   "",
		},
	}
}
