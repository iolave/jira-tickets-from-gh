# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [unreleased]
### Added
- Config file support including all previous cli options with the following new features:
  - New `sync[].jira.estimateField` option that allow users to specify the jira issue field name that stores story points (estimate). 
  - New `sync[].jira.issuePrefix` option that pre-appends the given prefix to jira issue titles. 
  - New `sync[].jira.issues[].transitionsToWip` and `sync[].jira.issues[].transitionsToWip` option that allows to set transitions per issue type.
- Verbose environment variable for docker.
- Docker compose volume to persist executions data.
- New `--user <string>` option in the `github list-projects` command.

### Changed
- `github-projects list` command was renamed to `github list-projects`
- Retrieving items from the local storage with(out) an url now checks if the url field is a valid url. This will cause that item's urls that are invalid are going to be replaced with a valid url that will point to a new jira issue.  
- GitHub issues that lack of status no longer creates Jira issues to prevent malformed Jira issues.
- Dockerfile now uses go implementation instead of bun's one. 

### Fixed
- GitHub task's estimate field now is reflected in Jira issues.


## [v0.3.0]
### Added
- Github project state is now stored locally (experimental).

## [v0.2.0]
### Added
- Docker config files.
- Added `--transitions-to-wip` option to specify Jira cloud transitions required in order to transition a task to a "Dev in progress" state.
- Added `--transitions-to-done` option to specify Jira cloud transitions required in order to transition a task to a "Done" state.
- Added `--sleep-time` that enables new executions tiggered when the sleep time passes by.

### Changed
- Log entries now have date and time.

## [v0.1.0]

### Added
- Sync command that creates Jira tickets from GitHub project cards.
- GitHub utility to list organization projects in order to extract a project id.

[unreleased]: https://github.com/iolave/bun-jira-tickets-from-gh/compare/v0.3.0...staging
[v0.3.0]: https://github.com/iolave/bun-jira-tickets-from-gh/releases/tag/v0.3.0
[v0.2.0]: https://github.com/iolave/bun-jira-tickets-from-gh/releases/tag/v0.2.0
[v0.1.0]: https://github.com/iolave/bun-jira-tickets-from-gh/releases/tag/v0.1.0
