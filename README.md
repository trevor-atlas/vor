[📋 changelog](https://github.com/trevor-atlas/vor/blob/master/CHANGELOG.md)
```
                  ___          ___
      ___        /\  \        /\  \
     /\  \      /::\  \      /::\  \
     \:\  \    /:/\:\  \    /:/\:\__\
      \:\  \  /:/  \:\  \  /:/ /:/  /
  ___  \:\__\/:/__/ \:\__\/:/_/:/__/___
 /\  \ |:|  |\:\  \ /:/  /\:\/:::::/  /
 \:\  \|:|  | \:\  /:/  /  \::/~~/~~~~
  \:\__|:|__|  \:\/:/  /    \:\~~\
   \::::/__/    \::/  /      \:\__\
    ~~~~         \/__/        \/__/
```

# Vör – Jira & Git made simple
In Norse mythology, Vör is the seeress. She is wise and of searching spirit, so that none can conceal anything from her. Her name means "awareness" or "to become aware of something", and she can be prayed to for intuitive information that cannot be acquired by normal means.

## Commands
**Create a branch for a specific jira issue**

```
vor branch AQ-1234
```
creates a branch of the form `{repo-name}/{issue-type}/{issue-number}/{issue-title}`
so for the aquicore repo issue 4753 would result in
`aquicore/story/AQ-4653/do-some-stuff-with-the-thing`

**Create a pr with my current branch**

```
vor pull-request
```

**list my issues in Jira**

```
vor issues
```

**View details of a specific issue**

```
vor issue AQ-1234
```

## Setup

Vör is project specific and usually requires a config at the root of that project (though this is not always true, it is recommended)

see an example config file [here](https://github.com/trevor-atlas/vor/blob/master/example.vor)

The supported configuration options and their defaults are:

```
devmode: false // output additional logging information at runtime

global.branchtemplate "{jira-issue-number}/{jira-issue-type}/{jira-issue-title}
branchtemplate: ""

global.jira.orgname: ""
jira.orgname: ""

global.jira.apikey: ""
jira.apikey: ""

global.jira.username: ""
jira.username: ""

global.github.apikey: ""
github.apikey: ""

github.owner: ""

global.git.path: "/usr/local/bin/git"
git.path: ""

global.git.pull-request-base: "master"
git.pull-request-base: ""
```
you can also export the config options in your `bash_profile` or elsewhere like so:

```
export JIRA_API_KEY={your jira api key}
```

Vör also supports global configs for some options:
these are loaded from your `~` directory or environment variables prefixed with `GLOBAL_`
```
export GLOBAL_JIRA_API_KEY=<your jira api key>
```

---

[Copyright © 2018 Trevor Atlas](https://github.com/trevor-atlas/vor/blob/master/LICENSE)
Vör is a command line tool to make working to Jira and Git/Github easier
