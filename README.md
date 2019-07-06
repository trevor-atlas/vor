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


## What is it?
Vor is a CLI tool that makes it really easy to connect jira and git/github
it provides commands to create branches from a given jira ticket, github pull requests from that branch and makes it easy to view your assigned tickets - all without leaving the command line! 

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
// output additional logging information at runtime
devmode: false

// the template to use for new branches
branchtemplate: {jira-issue-number}/{jira-issue-type}/{jira-issue-title}

// the path to your local git
git:
    path: /usr/local/bin/git

    // the base branch to make pull requests against
    pull-request-base: master

author: your name
email: you@yourdomain.xyz
jira:
    orgname: <your company name, usually contained in the url or your jira install>
    username: <your jira username (sometimes an email)>
    apikey: <your jira api key from id.atlassian.net>
github:
    owner: <the owner of the repository>
    apikey: <your github api key (get this from github.com/settings/tokens)>
    
```
you can also export the config options in your `bash_profile` or elsewhere like so:

```
export JIRA_API_KEY={your jira api key}
```

---

[Copyright © 2018 Trevor Atlas](https://github.com/trevor-atlas/vor/blob/master/LICENSE)
Vör is a command line tool to make working to Jira and Git/Github easier
