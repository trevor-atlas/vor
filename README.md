# Vor – Jira & Git made simple
// Vor, Norse Goddess who knows all.

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

**Open reviews for a ticket number**

```
vor review AQ-1234
```

## Setup

`export JIRA_API_KEY={your jira api key}`

id.atlassian.net get API key
