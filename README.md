# edgex-cli
[![Build Status](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/edgex-cli/job/main/badge/icon)](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/edgex-cli/job/main/) [![Code Coverage](https://codecov.io/gh/edgexfoundry/edgex-cli/branch/main/graph/badge.svg?token=wWeDPW5a81)](https://codecov.io/gh/edgexfoundry/edgex-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/edgexfoundry/edgex-cli)](https://goreportcard.com/report/github.com/edgexfoundry/edgex-cli) [![GitHub Latest Dev Tag)](https://img.shields.io/github/v/tag/edgexfoundry/edgex-cli?include_prereleases&sort=semver&label=latest-dev)](https://github.com/edgexfoundry/edgex-cli/tags) ![GitHub Latest Stable Tag)](https://img.shields.io/github/v/tag/edgexfoundry/edgex-cli?sort=semver&label=latest-stable) [![GitHub License](https://img.shields.io/github/license/edgexfoundry/edgex-cli)](https://choosealicense.com/licenses/apache-2.0/) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/edgexfoundry/edgex-cli) [![GitHub Pull Requests](https://img.shields.io/github/issues-pr-raw/edgexfoundry/edgex-cli)](https://github.com/edgexfoundry/edgex-cli/pulls) [![GitHub Contributors](https://img.shields.io/github/contributors/edgexfoundry/edgex-cli)](https://github.com/edgexfoundry/edgex-cli/contributors) [![GitHub Committers](https://img.shields.io/badge/team-committers-green)](https://github.com/orgs/edgexfoundry/teams/edgex-cli-committers/members) [![GitHub Commit Activity](https://img.shields.io/github/commit-activity/m/edgexfoundry/edgex-cli)](https://github.com/edgexfoundry/edgex-cli/commits)

## Introduction

EdgeX CLI is a command-line interface tool for developers, used for interacting with the EdgeX microservices.

See the [CLI Getting started documentation](https://docs.edgexfoundry.org/2.2/getting-started/tools/Ch-CommandLineInterface/) and the [EdgeX-CLI V2 Design ADR](https://github.com/edgexfoundry/edgex-docs/blob/main/docs_src/design/adr/core/0019-EdgeX-CLI-V2.md) for more information about the client.


## Installing EdgeX CLI

The client can be installed using a [snap](https://github.com/edgexfoundry/edgex-cli/tree/main/snap):

```
sudo snap install edgex-cli
```

To build, install and run `edgex-cli` natively, do the following:
```
git clone http://github.com/edgexfoundry/edgex-cli.git
cd edgex-cli
make tidy
make build
./bin/edgex-cli
```

## Limitations
- The client requires all services to run on the local host. It does not support a distributed configuration or using the API gateway ([#427](https://github.com/edgexfoundry/edgex-cli/issues/427))
- The `db` command from the v1 client is not supported ([#383](https://github.com/edgexfoundry/edgex-cli/issues/383))
- See this list of [all current enhancement issues](https://github.com/edgexfoundry/edgex-cli/issues?q=is%3Aissue+is%3Aopen+label%3Aenhancement) 

## Community
- [EdgeXFoundry Slack](https://slack.edgexfoundry.org/)
- [Mailing lists](https://lists.edgexfoundry.org/g/main)

## License
[Apache-2.0](LICENSE)
