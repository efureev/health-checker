module github.com/efureev/health-checker

go 1.18

replace (
	github.com/efureev/go-multierror => ../../packages/go-multierror
)

require (
	github.com/efureev/go-multierror v0.0.0-20220428202111-3a28ecba9fb0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/jackc/pgconn v1.12.1
	github.com/jackc/pgx/v4 v4.16.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
	golang.org/x/text v0.3.7 // indirect
)
