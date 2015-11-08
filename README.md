# Slack File Cleaner

Use this to cleanup the files of your team to keep storage low.

## Description

Copy `config.json.example` to `config.json`, add token of your team and run:

```bash
$ go install
$ $GOPATH/bin/slack-file-cleaner -days 30 -force
```

## Help

```bash
$ $GOPATH/bin/slack-file-cleaner -help
```
