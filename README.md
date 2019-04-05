# Control Tower

[![asciicast](https://asciinema.org/a/xVKD0dQuXdEmOcExt4A9WfbEN.svg)](https://asciinema.org/a/xVKD0dQuXdEmOcExt4A9WfbEN)

A tool for easily deploying [Concourse](https://concourse-ci.org) in a single command.

![CI build badge](https://ci.engineerbetter.com/api/v1/teams/main/pipelines/control-tower/jobs/system-test/badge)

## TL;DR

### AWS

```sh
$ AWS_ACCESS_KEY_ID=<access-key-id> \
  AWS_SECRET_ACCESS_KEY=<secret-access-key> \
  control-tower deploy --iaas aws <your-project-name>
```

### GCP

```sh
$ GOOGLE_APPLICATION_CREDENTIALS=<path/to/googlecreds.json> \
  control-tower deploy --iaas gcp <your-project-name>
```

## Why Control Tower?

The goal of Control Tower is to be the world's easiest way to deploy and operate Concourse CI in production.

In just one command you can deploy a new Concourse environment for your team, on either AWS or GCP. Your Control Tower deployment will *upgrade itself* and self-heal, restoring the underlying VMs if needed. Using the same command-line tool you can do things like manage DNS, scale your environment, or manage firewall policy. CredHub is provided for secrets management and Grafana for viewing your Concourse metrics.

You can keep up to date on Control Tower announcements by reading the [EngineerBetter Blog](http://www.engineerbetter.com/blog/) and by joining the discussion on our [Community Slack](https://join.slack.com/t/concourse-up/shared_invite/enQtNDMzNjY1MjczNDU3LTA1NzIxYTZkYjFkMjA2ODBmY2E2OTM3OGE3YTc2OGViNTMxYTY4MjYwNGNjOTAxNDNiOGE5NzhmMTQ2NWVhNzQ).

## Features

| **Feature** | **AWS** | **GCP** |
|:------------|:-------:|:-------:|
| Concourse IP whitelisting | **+** | **+** |
| Credhub | **+** | **+** |
| Custom domains | **+** | **+** |
| Custom tagging | **BOSH only** | **BOSH only** |
| Custom TLS certificates | **+** | **+** |
| Database vertical scaling | **+** | **+** |
| GitHub authentication | **+** | **+** |
| Grafana (on port 3000) | **+** | **+** |
| Interruptable worker support | **+** | **+** |
| Letsencrypt integration | **+** | **+** |
| Namespace support | **+** | **+** |
| Region selection | **+** | **+** |
| Retrieving deployment information | **+** | **+** |
| Retrieving deployment information as shell exports | **+** | **+** |
| Retrieving deployment information in JSON | **+** | **+** |
| Retrieving director NATS cert expiration | **+** | **+** |
| Rotating director NATS cert | **+** | **+** |
| Self-Update support | **+** | **+** |
| Teardown deployment | **+** | **+** |
| Web server vertical scaling | **+** | **+** |
| Worker horizontal scaling | **+** | **+** |
| Worker type selection | **+** | **N/A** |
| Worker vertical scaling | **+** | **+** |
| Zone selection | **+** | **+** |
| Customised networking | **+** | **+** |

## Detailed Documentation

| | |
|:-|:-|
|Before you start|[Prerequisites](docs/prerequisites.md)|
|Installing Control Tower|[Installation](docs/installation.md)|
|Flags on all commands|[Global flags](docs/global.md)|
|Deploying a Concourse|[Deploy](docs/deploy.md)|
|Retrieving info from a deployment|[Info](docs/info.md)|
|Destroying a Concourse|[Destroy](docs/destroy.md)|
|Maintaining your Concourse|[Maintain](docs/maintain.md)|
|Updating|[Updating](docs/updating.md)|
|Metrics|[Metrics](docs/metrics.md)|
|Credential Management|[Credhub](docs/credhub.md)|
|How much will this cost?|[Cost Estimation](cost.md)|
|What is it doing? - deep dive|[Walkthrough](docs/walkthrough.md)|
|Want to Contribute?|[Development](docs/development.md)|
