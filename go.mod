module github.com/haung921209/nhn-cloud-cli

go 1.24.0

require (
	github.com/AlecAivazis/survey/v2 v2.3.7
	github.com/haung921209/nhn-cloud-sdk-go v0.1.26
	github.com/jmespath/go-jmespath v0.4.0
	github.com/spf13/cobra v1.8.0
	golang.org/x/crypto v0.46.0
	golang.org/x/term v0.38.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-isatty v0.0.8 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)

replace github.com/haung921209/nhn-cloud-sdk-go => ../nhn-cloud-sdk-go
