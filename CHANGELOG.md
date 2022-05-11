
<a name="EdgeX Command Line Interface (found in edgex-cli) Changelog"></a>
## EdgeX Command Line Interface (CLI)
[Github repository](https://github.com/edgexfoundry/edgex-cli)

### Change Logs for EdgeX Dependencies
- [go-mod-core-contracts](https://github.com/edgexfoundry/go-mod-core-contracts/blob/main/CHANGELOG.md)

## [v2.2.0] Kamakura - 2022-05-11  (Not Compatible with 1.x releases)

### Bug Fixes 🐛
- Replace slack.edgexfoundry.org with edgexfoundry.slack.com ([#b653c32](https://github.com/edgexfoundry/edgex-cli/commits/b653c32))
- Wrong list param ([#56017f9](https://github.com/edgexfoundry/edgex-cli/commits/56017f9))
- Wrong device param of command list ([#c7441b1](https://github.com/edgexfoundry/edgex-cli/commits/c7441b1))
### Code Refactoring ♻
- **snap:** remove redundant content identifier ([#8c7365f](https://github.com/edgexfoundry/device-modbus-go/commits/8c7365f))

### Documentation 📖
- **snap:** Move usage instructions to docs ([#440](https://github.com/edgexfoundry/edgex-cli/issues/440)) ([#081a76a](https://github.com/edgexfoundry/edgex-cli/commits/081a76a))
- **snap:** add snap README.md ([#9fc8be8](https://github.com/edgexfoundry/edgex-cli/commits/9fc8be8))

### Build 👷
- **snap:** Source metadata from central repo ([#bca6647](https://github.com/edgexfoundry/edgex-cli/commits/bca6647))

### Continuous Integration 🔄
- gomod changes related for Go 1.17 ([#433](https://github.com/edgexfoundry/edgex-cli/issues/433)) ([#cb94420](https://github.com/edgexfoundry/edgex-cli/commits/cb94420))
- Go 1.17 related changes ([#432](https://github.com/edgexfoundry/edgex-cli/issues/432)) ([#10c2baf](https://github.com/edgexfoundry/edgex-cli/commits/10c2baf))

## [v2.1.0] Jakarta - 2021-11-18  (Not Compatible with 1.x releases)
### Features ✨
- Add V2 support-scheduler ([#6d2939f](https://github.com/edgexfoundry/edgex-cli/commits/6d2939f))
- Add V2 support-notifications ([#c21b0cc](https://github.com/edgexfoundry/edgex-cli/commits/c21b0cc))
- Add support for V2 core-metadata endpoints ([#419](https://github.com/edgexfoundry/edgex-cli/issues/419)) ([#a77b052](https://github.com/edgexfoundry/edgex-cli/commits/a77b052))
- Add support for core-data endpoints ([#082273a](https://github.com/edgexfoundry/edgex-cli/commits/082273a))
- Add support for core-command list ([#337d45e](https://github.com/edgexfoundry/edgex-cli/commits/337d45e))
- Add support for core-command write ([#87a5b2f](https://github.com/edgexfoundry/edgex-cli/commits/87a5b2f))
- Add support for core-command read ([#83c667c](https://github.com/edgexfoundry/edgex-cli/commits/83c667c))
- Add support for common v2 endpoints (ping/metrics/version/config) ([#1b3752f](https://github.com/edgexfoundry/edgex-cli/commits/1b3752f))
- Remove go-mod-core-contracts v1 from attribution.txt ([#1615725](https://github.com/edgexfoundry/edgex-cli/commits/1615725))
- Upgrade to v2 apis with base files and status and version commands working ([#e5e45d9](https://github.com/edgexfoundry/edgex-cli/commits/e5e45d9))
- Add context with correlation-ID Value ([#f6afed9](https://github.com/edgexfoundry/edgex-cli/commits/f6afed9))

### Bug Fixes 🐛
- `make install` creates $home/.edgex-cli/configuration.toml ([#372](https://github.com/edgexfoundry/edgex-cli/issues/372)) ([#b6bf22c](https://github.com/edgexfoundry/edgex-cli/commits/b6bf22c))
- Instead of passing the entire cmd, send just the context ([#3ce8009](https://github.com/edgexfoundry/edgex-cli/commits/3ce8009))
- Snap-fix-build-after-removal-of-snap-local ([#4e86fd6](https://github.com/edgexfoundry/edgex-cli/commits/4e86fd6))
- Add snap related changes ([#a4c6c2e](https://github.com/edgexfoundry/edgex-cli/commits/a4c6c2e))
- Propagate configuration file to distribution archives ([#271f6fd](https://github.com/edgexfoundry/edgex-cli/commits/271f6fd))
- Missing configuration.toml in release packages ([#357](https://github.com/edgexfoundry/edgex-cli/issues/357)) ([#7729e0d](https://github.com/edgexfoundry/edgex-cli/commits/7729e0d))
- Rewrite sample files to use json format only ([#1c434c9](https://github.com/edgexfoundry/edgex-cli/commits/1c434c9))
- Make Device Profile message more accurate ([#f97e0f4](https://github.com/edgexfoundry/edgex-cli/commits/f97e0f4))
- **snap:** fix initial configuration file install ([#368](https://github.com/edgexfoundry/edgex-cli/issues/368)) ([#1f3ac9c](https://github.com/edgexfoundry/edgex-cli/commits/1f3ac9c))

### Continuous Integration 🔄
- Standardize dockerfiles ([#352](https://github.com/edgexfoundry/edgex-cli/issues/352)) ([#2a116c9](https://github.com/edgexfoundry/edgex-cli/commits/2a116c9))

## [v1.0.1] - 2021-02-08
### Bug Fixes 🐛
- Make install creates $home/.edgex-cli/configuration.toml ([#ba2c4a5](https://github.com/edgexfoundry/edgex-cli/commits/ba2c4a5))
- Fix missing configuration file in distribution archives ([#4b89f1c](https://github.com/edgexfoundry/edgex-cli/commits/4b89f1c))

## [v1.0.0] - 2020-11-30
### Features ✨
- Add snap packaging ([#333](https://github.com/edgexfoundry/edgex-cli/issues/333)) ([#796f864](https://github.com/edgexfoundry/edgex-cli/commits/796f864))
- Add support of Provisioning watchers ([#c1be68d](https://github.com/edgexfoundry/edgex-cli/commits/c1be68d))
- Add attribution and LICENSE files in archive files ([#348d03a](https://github.com/edgexfoundry/edgex-cli/commits/348d03a))
- Issue pet/put commands ([#266433f](https://github.com/edgexfoundry/edgex-cli/commits/266433f))
- Add command for update Device Profile ([#92857c8](https://github.com/edgexfoundry/edgex-cli/commits/92857c8))
- Provides code organization and sample template information for developers ([#316](https://github.com/edgexfoundry/edgex-cli/issues/316)) ([#2d0fa2c](https://github.com/edgexfoundry/edgex-cli/commits/2d0fa2c))
- Add device operation status update command ([#351a6c3](https://github.com/edgexfoundry/edgex-cli/commits/351a6c3))
- Add support of multiple delete by id ([#b0b4e8f](https://github.com/edgexfoundry/edgex-cli/commits/b0b4e8f))
- Update device admin state ([#3dab669](https://github.com/edgexfoundry/edgex-cli/commits/3dab669))
- Define -u and -v flags ([#bb862d5](https://github.com/edgexfoundry/edgex-cli/commits/bb862d5))

### Bug Fixes 🐛
- Update edgex-cli documentation ([#a99930f](https://github.com/edgexfoundry/edgex-cli/commits/a99930f))
- Limits Makefile clean action to removing edgex-cli binary ([#322](https://github.com/edgexfoundry/edgex-cli/issues/322)) ([#2b9df61](https://github.com/edgexfoundry/edgex-cli/commits/2b9df61))
- Added y/n checkboxes for new imports ([#323](https://github.com/edgexfoundry/edgex-cli/issues/323)) ([#a097169](https://github.com/edgexfoundry/edgex-cli/commits/a097169))
- Updates command documentation, including that CLI works only with localhost ([#312](https://github.com/edgexfoundry/edgex-cli/issues/312)) ([#d5c6e11](https://github.com/edgexfoundry/edgex-cli/commits/d5c6e11))
- Wrong err msg ([#245635c](https://github.com/edgexfoundry/edgex-cli/commits/245635c))

### Code Refactoring ♻
- Make update device adminstate cmd shorter ([#53c3bc0](https://github.com/edgexfoundry/edgex-cli/commits/53c3bc0))
- create Device profile to use Json ([#7fa708b](https://github.com/edgexfoundry/edgex-cli/commits/7fa708b))
