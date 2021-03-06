# Changelog
All notable changes to this project will be documented in this file using [CHANGELOG](https://keepachangelog.com/en/0.3.0/) format.
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

For details about the features changes, you can look at their CHANGELOG.

## 0.10.14 - 2021-05-25
### Changed
- Base release to fix missing package

## 0.10.13 - 2021-05-25
### Changed
- Base release to fix python installation

## 0.10.12 - 2021-05-25
### Changed
- Terraform release 0.9.1

## 0.8.1 - 2020-06-24
### Changed
- Updated fluentd chart

## 0.8.0 - 2019-11-18
### Added
- Vault release 0.1.0

## 0.7.1 - 2019-10-30
### Changed
- Terraform release 0.2.0

## 0.7.0 - 2019-10-24
### Changed
- **WARNING** (breaking change): `bom.yaml` moved from `$GOPATH/src/github.com/scalair/robokops/bom.yaml` to `/etc/robokops/bom.yaml`
- **WARNING**: binary moved from `$GOPATH/bin/robokops` to `/usr/local/bin/robokops`. You must remove the previous binary in order to use that version!
- Installation process has evolved (moving from `go get` to a proper package installation). Installation via `go get` is no longer supported!

### Fixed
- robokops.go: Colors were not initialized soon enough

### Added
- robokops.go: Support for yellow color

## 0.6.1 - 2019-10-24
### Fixed
-  Fix go-releaser with Github actions

## 0.6.0 - 2019-10-23
### Added
-  Aws-efs-csi-driver release 0.1.0

## 0.5.3 - 2019-10-15
### Changed
- Cluster-init release 0.1.2

## 0.5.2 - 2019-10-11
### Changed
- Cluster-init release 0.1.1

## 0.5.1 - 2019-10-02
### Fixed
- [PR12](https://github.com/scalair/robokops/pull/12) Remove ended container from docker host

## 0.5.0 - 2019-10-01
### Added
- Jenkins release 0.1.0

## 0.4.0 - 2019-10-01
### Added
- Kubewatch release 0.1.0

## 0.3.0 - 2019-09-30
### Added
- Velero release 0.1.0

## 0.2.0 - 2019-09-25
### Changed
- Elastic-stack release 0.2.0

## 0.1.0 - 2019-08-26
### Added
- Initial commit
