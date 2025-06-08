module github.com/detectviz/detectviz

go 1.24.3

require (
	github.com/robfig/cron/v3 v3.0.1 // @grafana/grafana-backend-group
	github.com/stretchr/testify v1.10.0 // @grafana/grafana-backend-group
	go.uber.org/zap v1.27.0 // @grafana/identity-access-team
	gopkg.in/yaml.v3 v3.0.1 // indirect; @grafana/alerting-backend
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/redis/go-redis/v9 v9.10.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// Use fork of crewjam/saml with fixes for some issues until changes get merged into upstream
replace github.com/crewjam/saml => github.com/grafana/saml v0.4.15-0.20240917091248-ae3bbdad8a56

// Use our fork of the upstream alertmanagers.
// This is required in order to get notification delivery errors from the receivers API.
replace github.com/prometheus/alertmanager => github.com/grafana/prometheus-alertmanager v0.25.1-0.20250417181314-6d0f5436a1fb

exclude github.com/mattn/go-sqlite3 v2.0.3+incompatible

// lock for mysql tsdb compat
replace github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.7.1

// v1.* versions were retracted, we need to stick with v0.*. This should work
// without the exclude, but this otherwise gets pulled in as a transitive
// dependency.
exclude github.com/prometheus/prometheus v1.8.2-0.20221021121301-51a44e6657c3

// This was retracted, but seems to be known by the Go module proxy, and is
// otherwise pulled in as a transitive dependency.
exclude k8s.io/client-go v12.0.0+incompatible
