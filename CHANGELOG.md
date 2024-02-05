# Change Log

All notable changes to the "goploy" extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),

## [1.16.2] - 2024-02-05

### Added
- support periodic update password
- support project hyperlink

### Changed

- change load config publish event

### Bug fixed

- fix file compare missing server id
- fix get cron logs invalid request method

## [1.16.1] - 2024-01-10

### Added
- add process exec time
- add manage repository
- add docker env

### Changed

- change detail dialog css

### Bug fixed

- fix element-plus auto import

## [1.16.0] - 2023-12-29

### Added
- [project script yaml mode](https://github.com/zhenorzz/goploy/pull/67)
- add deploy time
- create user when ldap is enabled
- add cron logs

### Bug fixed

- fix element-plus locale
- fix toggle server state
- fix project set auto deploy


## [1.15.3] - 2023-10-23

### Bug fixed

- ldap validate password
- deploy page duplicate item

## [1.15.2] - 2023-08-17

### Added
- [add captcha](https://github.com/zhenorzz/goploy/pull/64)
- add predefine var ${PROJECT_LABEL}
- add predefine var ${PROJECT_ENV}
- support predefine var in transmitter
- support unbind projects in server page

### Changed

- changed file structure
- optimize projectUser struct
- optimize projectServer struct

### Bug fixed

- fix install agent


## [1.15.1] - 2023-07-18

### Added

- add user api key
- add publish progress
- [add vscode extension](https://github.com/goploy-devops/goploy-vscode)
- [add jetbrains extension](https://github.com/goploy-devops/goploy-jetbrains)

### Changed

- support media login log

### Bug fixed

- [fix dingtalk token cache #58](https://github.com/zhenorzz/goploy/commit/d823222a6e19bb9691caa419eee164404606cede)
- fix commit info


## [1.15.0] - 2023-06-28

### Added

- [lark & dingtalk login #55](https://github.com/zhenorzz/goploy/pull/55)
- support deploy finish script
- deploy project animation

### Changed

- [sftp open / after connected](https://github.com/zhenorzz/goploy/commit/d041f9e13ba756453216f74867c7af3b27ae3534)
- [load namespace list when the page opened](https://github.com/zhenorzz/goploy/commit/d171343857356ab2f55db07e96a7004f1d2fe30a)

### Bug fixed

- [fix missing project name in lark notify](https://github.com/zhenorzz/goploy/commit/0372c2e9b8542fb7e434a07b9a909fd52d60e58f)

## [1.14.0] - 2023-05-25

### Added

- Horizontal sidebar
- [Support monitor execute script #51](https://github.com/zhenorzz/goploy/pull/51)
- [Support project tag #51](https://github.com/zhenorzz/goploy/pull/51)
- [Add nginx manage #52](https://github.com/zhenorzz/goploy/pull/52)
- Support sftp copy rename edit file

### Changed

- Add app version in login page
- Optimize code

### Bug fixed

- [Jump server #54](https://github.com/zhenorzz/goploy/pull/54)
- Ping db after open connection
- Fix page css

## [1.13.1] - 2023-04-29

### Bug fixed

- Fix windows mkdir
- Refresh file sync project list
- [Fix check monitor config](https://github.com/zhenorzz/goploy/commit/dae371d96cbf241c4dbf9776c48dc1611e4dcff1)
- [Support sftp --delete option](https://github.com/zhenorzz/goploy/commit/5a0092412151212cab4adc8ca3d589eddb279645)
- [Fix monitor can not redo the task](https://github.com/zhenorzz/goploy/commit/4db5bd6f7bcfa4e647d923d10a59504548408e37)

## [1.13.0] - 2023-03-03

### Added

- Batch execute script

### Changed

- Show server name in monitor table

### Bug fixed

- Navbar dark mode
- Windows ssh connect

## [1.12.3] - 2023-02-03

### Changed

- Adjust commit list table column
- Refactor code
- Update Server process list UI

### Bug fixed

- Server process output newline

## [1.12.2] - 2022-11-02

### Changed

- Support delete file via sftp

### Bug fixed

- Key column 'server_id' doesn't exist in table


## [1.12.1] - 2022-10-22

### Changed

- Support double click for sftp directory
- Update package.json

### Bug fixed

- Translation
- Dark mode
## [1.12.0] - 2022-10-11

### Added

- Third login
- Transfer file across server

### Changed

- SFTP support mutiple server
- Server process manage

### Bug fixed

- Xterm disconnect before unmount

## [1.11.0] - 2022-09-19

### Added

- Notify tag after deploy
- Operation log
- serial/parallel publish to server

### Bug fixed

- Fix after deploy script

## [1.10.0] - 2022-08-31

### Added

- Server process manager
- Customize project transfer protocol

### Changed

- Tag view UI
- Support more predefined vars

## [1.9.1] - 2022-08-16

### Bug fixed

- Fix dockerfile
- Fix bat newline
- Fix sftp dropdown

## [1.9.0] - 2022-07-11

### Added

- Dark mode

### Changed

- Sidebar

## [1.8.1] - 2022-07-02

### Added

- Support after deploy script replace commit info

### Bug fixed

- Fix svn commit list

## [1.8.0] - 2022-05-17

### Added

- Support sftp transfer files
- Show deploy detail in realtime

### Bug fixed

- fix deploy list remove item failed
- fix ftp login anonymous

## [1.7.1] - 2022-05-09

### Changed

- show server on edit monitor item

### Bug fixed

- fix ace editor not found
- fix monitor nil pointer

## [1.7.0] - 2022-04-29

### Added

- Monitor

## [1.6.1] - 2022-04-20

### Added

- Support password login ssh (only work in linux)

### Bug fixed

- fix web script editor

## [1.6.0] - 2022-04-11

### Added

- RBAC

## [1.5.0] - 2022-03-24

### Added

- detect project name is link
- import csv in server page
- install agent in server page

### Changed

- update element-plus to 2.0
- goploy-agent check sign
- project dialog modify server and user
- file sync move to deploy page

### Bug fixed

- fix vue3 SFCs ref undefined
- fix cron task date popover
- fix publish detail filter popover
- fix namespace add user

## [1.4.7] - 2022-03-14

### Added

- support deploy ftp & sftp

### Changed

- script setup SFCs
- migrate docs to goploy-devops/goploy-doc

### Bug fixed

- sftp file upload

## [1.4.6] - 2022-02-24

### Added

- web log
- sftp file preview

### Bug fixed

- web cookies undefined

## [1.4.5] - 2022-01-26

### Added

- new web shell
- new sftp
- support copy server config

### Bug fixed

- git current branch

## [1.4.4] - 2022-01-17

### Added

- support jump server
- process manager
- split log

### Changed

- decode query

## [1.4.3] - 2021-12-24

### Changed

- code
- select db
### Fixed

- fix exit deploy script
- fix tag refresh
- fix deploy filter
- fix file upload

## [1.4.2] - 2021-12-15

### Added

- .env -> goploy.toml
- support ldap

### Fixed

- fix 飞书构建通知

## [1.4.1] - 2021-12-09

### Added

- file compare

## [1.4.0] - 2021-12-04

### Added

- second's cron
- support svn hook

### Changed

- route

### Fixed

- svn commit id length

## [1.3.8] - 2021-11-20

### Added

- server notify
- send command to all xterm

### Changed

- ts type

### Fixed

- fix web re-login

## [1.3.7] - 2021-11-09

### Added

- monitor server performance

### Fixed

- fix web redirect
- fix web date select i18n

## [1.3.6] - 2021-10-15

### Added

- monitor support http

### Changed

- delete cache

### Fixed

- fix task block
- fix symlink rollback

## [1.3.5] - 2021-09-18

### Added

- allow sort server ip
- add server configuration
- support multiple browser tabs

### Changed

- fix vite hot reload

## [1.3.4] - 2021-08-20

### Added

- support update app version
- support svn
- customize symlink backup number
- add cmd mode in pull script

### Changed

- repository factory (for support other protocol in the future)

## [1.3.3] - 2021-07-16

### Added

- web sftp
- support deploy table sorting

## [1.3.2] - 2021-06-25

### Changed

- web ssh

### Fixed

- fix illegal namespace
- fix web keep alive
- fix copy public key

## [1.3.1] - 2021-05-30

### Changed

- vite + vue3 + ts
- mobile compatible

### Fixed

- symlink in docker
- placeholder

## [1.2.2] - 2021-03-26

### Changed

- go embed static file
- more notify content

### Fixed

- fix symlink rollback

## [1.2.1] - 2021-03-03

### Added

- server terminal
- server can stay in any namespace

### Changed

- http.put for edit

### Fixed

- fix wss protocol
- delete trim rsync option

## [1.1.6] - 2021-02-12

### Added

- add ssh key path
- server host supports domain
- support graceful stop
- support symlink rebuild

### Changed

- add git url tips
- delete server install module

## [1.1.5] - 2021-01-20

### Changed

- deploy detail filters

### Fixed

- fix detail loading
- fix missing sql

## [1.1.4] - 2021-01-07

### Added

- add flag --asset-dir=

### Changed

- delete rsync option --delete-after
- unique project file

### Fixed

- fix copy project
- fix refresh tag view
- fix ssh fingerprint

## [1.1.3] - 2021-01-06

### Added

- upload project file

## [1.1.2] - 2020-12-26

### Added

- tags view

### Fixed

- fix get detail timeout

## [1.1.1] - 2020-12-03

### Added

- branch deploy

## [1.0.7] - 2020-11-27

### Added

- reset deploy state

## [1.0.6] - 2020-11-07

### Added

- grey publish

## [1.0.5] - 2020-11-06


### Added

- deploy tag list
- predefined vars

## [1.0.4] - 2020-10-25


### Added

- Project review

### Fixed

- Fix monitor bug

## [1.0.3] - 2020-10-11

### Added

- Monitor
- notify times
- error content

### Fixed

- Fix monitor bug

### Changed

- table loading

## [1.0.2] - 2020-09-04

### Added

- I18n

### Fixed

- Fix SQL error

### Changed

- project path

## [1.0.1] - 2020-08-21

### Added

- namespace

### Changed

- Auto deploy
- change project_name to project_id
