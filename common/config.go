package common

import (
	_ "embed"
	"time"
)

var LocalConfig = &Config{
	System: System{
		User:               "",
		Password:           "",
		Addr:               "0.0.0.0:8080",
		URLPrefix:          "/",
		DataDir:            "./data",
		Name:               "NoLet",
		DSN:                "",
		Cert:               "",
		Key:                "",
		ReduceMemoryUsage:  false,
		ProxyHeader:        "",
		MaxBatchPushCount:  -1,
		MaxAPNSClientCount: 1,
		Concurrency:        256 * 1024,
		ReadTimeout:        3 * time.Second,
		WriteTimeout:       3 * time.Second,
		IdleTimeout:        10 * time.Second,
		Debug:              false,
		Version:            "",
		BuildDate:          "",
		CommitID:           "",
		ICPInfo:            "",
		TimeZone:           "UTC",
		Voice:              false,
		Auths:              []string{},
	},
	Apple: Apple{
		ApnsPrivateKey: `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgvjopbchDpzJNojnc
o7ErdZQFZM7Qxho6m61gqZuGVRigCgYIKoZIzj0DAQehRANCAAQ8ReU0fBNg+sA+
ZdDf3w+8FRQxFBKSD/Opt7n3tmtnmnl9Vrtw/nUXX4ldasxA2gErXR4YbEL9Z+uJ
REJP/5bp
-----END PRIVATE KEY-----`,
		Topic:   "me.uuneo.Meoworld",
		KeyID:   "BNY5GUGV38",
		TeamID:  "FUWV6U942Q",
		Develop: false,
	},
}
