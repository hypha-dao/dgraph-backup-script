module github.com/sebastianmontero/dgraph-backup-script

go 1.16

require (
	github.com/machinebox/graphql v0.2.2
	github.com/matryer/is v1.4.0 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/robfig/cron/v3 v3.0.0
	github.com/rs/zerolog v1.20.0
	github.com/sebastianmontero/slog-go v0.0.0-20210801140624-25c2da708d0b
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
)

// replace github.com/sebastianmontero/slog-go => ../slog-go
