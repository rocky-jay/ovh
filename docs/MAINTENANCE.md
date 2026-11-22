# 维护接手 / Maintenance handoff

- 来源 / Origin: [gokele/ovh](https://github.com/gokele/ovh), archived. The fork baseline is `025ac9629db0f50dae3133175fa8b6107564ecb3` (`VERSION=0.4.37`); this is a source baseline, not a release of this fork.
- 许可证 / License: the original history and [AGPL-3.0 text](../LICENSE) remain in place. This community fork is not affiliated with OVHcloud.
- 本仓 / Current fork: [rocky-jay/ovh](https://github.com/rocky-jay/ovh) is based on [luodaoyi/ovh](https://github.com/luodaoyi/ovh). New standalone builds check [rocky-jay/ovh Releases](https://github.com/rocky-jay/ovh/releases) by default. `OVH_UPDATE_API` remains an explicit override.

## 首个自有版本的验收 / Before the first fork release

1. Review the inherited update, release, and container paths together. Do not redirect installers or self-updaters to a replacement channel before its assets, integrity checks, and version handling have been verified. Keep existing installations unchanged during this review.
2. Establish CI for backend and frontend checks without real OVH credentials or external write calls. Use mocks for behavior that submits orders, initiates payment, or changes a server or account. Never exercise these paths against a real account in CI.
3. Verify builds and installation in an isolated environment. Test migration from a disposable copy of the data and its matching encryption key; do not use production data or credentials. Document supported platforms and rollback before publishing.

Do not paste credentials, payment information, database files, encryption keys, or unredacted logs into public PRs.
