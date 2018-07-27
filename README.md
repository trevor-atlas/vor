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

`export JIRA_API_KEY={your jira api key}`

id.atlassian.net get API key
