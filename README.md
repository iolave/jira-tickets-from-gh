# Sync Jira and GitHub projects
I was forced by my job to use Jira to keep track of tasks, but I kept forgetting to update status within Jira. Most of time im in GitHub and it's just easier to have issues and close them via PR's.

...so i built this tool for myself and hopefully, it will help you too.

> [!WARNING]
> All versions released prior to `v1.0.0` are to be considered [breaking changes](https://semver.org/#how-do-i-know-when-to-release-100) (I'll try my best to not push breaking changes btw).

## Installation
```bash
bun install -g jira-tickets-from-gh
```

## Pre-requisites for running the CLI
### A GitHub project with required fields
Make sure you have a GitHub project with the following fields:

- `Title`: Title for the task.
- `Jira issue type`: choice field with available jira issue types.
- `Jira URL`: text field to store jira url.
- `Status`: `Todo`, `In Progress`, `Done` choice field.
- `Estimate`: Number field.
- `Repository`: Default field for repository info.

### Get a Jira cloud token
To get a jira api token from your jira cloud account please refer to the [docs](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/).

### Environment variables
- `GITHUB_TOKEN`: Your GitHub token. If the project you're trying to sync is in an organization, make sure the token have access to it.
- `JIRA_TOKEN`: Jira api token used for auth.
- `JIRA_EMAIL`: Jira email used for auth.

### Get the id of your github project
The cli is shipped with a utility that's going to help you to search a GitHub project id.

#### List organization projects
```bash
jira-tickets-from-gh --gh-token=GH_TOKEN github list-projects --org=<ORG>
# or with envs
# export GITHUB_TOKEN=token
# jira-tickets-from-gh github list-projects --org=<ORG>
```

#### List user projects
```bash
jira-tickets-from-gh --gh-token=GH_TOKEN github list-projects --user=<GH_USER>
# or with envs
# export GITHUB_TOKEN=token
# jira-tickets-from-gh github list-projects --user=<GH_USER>
```

## Using the CLI to sync projects
Use `jira-tickets-from-gh sync` command to sync a GitHub project with a Jira cloud project.

### Config YAML Schema definition

| Property                                    | Required | Description |
|---------------------------------------------|:--------:|-------------|
| `sleepTime`                                 |`false`	 | sleep time between executions (if not specified the program will run once) |
| `enableApi`                                 |`false`	 | serves an api to interact with the projects storage and manage tasks manually (like moving a task to done) |
| `sync[].name`                               |`true`	 | tag to identify a sync project |
| `sync[].assignees[]`                        |`false`	 | map of GitHub users to Jira ones (email)  |
| `sync[].assignees[].jiraEmail`	      |`true`	 | Jira email |
| `sync[].assignees[].ghUer`    	      |`true`	 | GitHub user |
| `sync[].github.projectId`		      |`true`	 | Github project ID |
| `sync[].jira.subdomain`		      |`true`	 | Jira subdomain |
| `sync[].jira.projectKey`		      |`true`	 | Jira project key (usually a the short name) |
| `sync[].jira.estimateField`		      |`false`	 | Jira field name within the api response that stores story points (estimate) |
| `sync[].jira.issuePrefix`		      |`false`	 | Prefix to be added to Jira issues |
| `sync[].jira.issues`    		      |`true`	 | Jira issues type definition |
| `sync[].jira.issues.type`    		      |`true`	 | Jira issue name (ie. Task) |
| `sync[].jira.issues.transitionsToWip[]`     |`false`	 | Jira issue transitions to get to a WIP status (in the future a jira utility will be added to retrieve this transitions) |
| `sync[].jira.issues.transitionsToDone[]`    |`false`	 | Jira issue transitions to get to a DONE status (in the future a jira utility will be added to retrieve this transitions) |

### Example
*Using environment variables*
```bash
export GITHUB_TOKEN=token
export JIRA_TOKEN=token
export JIRA_EMAIL=email@org.com
jira-tickets-from-gh sync --config ./config.yml
```

*Passing tokens via the cli*
```bash
export GITHUB_TOKEN=token
export JIRA_TOKEN=token
jira-tickets-from-gh --gh-token=TOKEN --jira-token=TOKEN sync --config ./config.yml
```

*Execution example*
```
[2024-08-01 10:45:46][INFO]	syncCmd.action                          	creating jira issue                                         	{"title":"TEST: bun-jira-tickets-from-gh"}
[2024-08-01 10:45:47][INFO]	syncCmd.action                          	created issue                                               	{"url":"https://mfhnet.atlassian.net/browse/TEST3-112"}
```

## Running using Docker
### Environment variables
<!-- TODO: Update this part of the docs -->
| Env                   | Maps to option          |
|-----------------------|-------------------------|
| GITHUB_TOKEN          | `--gh-token` |
| GH_PROJECT_ID         | `--gh-project-id` |
| GH_USERS_MAP          | `--gh-assignees-map` |
| JIRA_TOKEN            | `--jira-token` |
| JIRA_SUBDOMAIN        | `--jira-subdomain` |
| JIRA_PROJECT_KEY      | `--jira-project-key` |
| JIRA_WIP_TRANSITIONS  | `--transitions-to-wip` |
| JIRA_DONE_TRANSITIONS | `--transitions-to-done` |
| JIRA_ISSUE_PREFIX	| `--jira-issue-prefix` |
| JIRA_ESTIMATE_FIELD	| `--jira-estimate-field` |
| SLEEP_TIME            | `--sleep-time` |
| VERBOSE               | if value is set to `true` then `-v` option is mapped |

### Example env file
<!-- TODO: Update this part of the docs -->
```
GITHUB_TOKEN=TOKEN
GH_PROJECT_ID=PROJECT_ID
GH_USERS_MAP=GH_USER:JIRA_EMAIL
JIRA_TOKEN=TOKEN
JIRA_SUBDOMAIN=SUBDOMAIN
JIRA_PROJECT_KEY=TEST
JIRA_WIP_TRANSITIONS=2
JIRA_DONE_TRANSITIONS=3
JIRA_ISSUE_PREFIX=[BACKEND]
JIRA_ESTIMATE_FIELD=customfield_10016
SLEEP_TIME=600000
VERBOSE=false
```

### Build
```bash
docker compose build
```

### Run
```bash
docker compose --env-file=path/to/env up -d
```

