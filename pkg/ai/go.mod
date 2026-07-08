module github.com/kyanite/ai

go 1.26.0

toolchain go1.26.5

require (
	github.com/kyanite/appnames v0.0.0-00010101000000-000000000000
	github.com/kyanite/config v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.12.3
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/knadh/koanf/parsers/yaml v1.1.0 // indirect
	github.com/knadh/koanf/providers/confmap v1.0.0 // indirect
	github.com/knadh/koanf/providers/env v1.1.0 // indirect
	github.com/knadh/koanf/providers/file v1.2.1 // indirect
	github.com/knadh/koanf/v2 v2.3.5 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.42.0 // indirect
)

replace github.com/kyanite/config => ../../pkg/config

replace github.com/kyanite/appnames => ../../pkg/appnames
