# List

This command list all active instances and sessions that match the filter.
When interactive mode is enabled you can start a session in a given instance or terminate active session.

```console
sigil list [flags]
```

[Man](../man/sigil_list.md) page

## Sample config

Config file settings that affect the command

```toml
[default]
  output-format = "text"
  region = "eu-west-1"
  profile = "dev"
  interactive = false
  [default.filters.session]
    after="2018-08-29T00:00:00Z"
    before="2019-08-29T00:00:00Z"
    target="i-xxxxxxxxxxxxxxxx1"
    owner="user@example.com"
  [default.filters.instance]
    ids=["i-xxxxxxxxxxxxxxxx1","i-xxxxxxxxxxxxxxxx2"]
    tags = [
        {key="Name", values=["Web","DB"] }
    ]
```

## Examples

List instances

```console
$ sigil list --instance-tags '[{"key":"Name","values":["Web","DB"]}]'
Index  Name       Instance ID          IP Address  Private DNS Name
1      Web        i-xxxxxxxxxxxxxxxx1  10.10.10.1  test1.local
2      DB         i-xxxxxxxxxxxxxxxx2  10.10.10.2  test2.local
```

List sessions

```console
$ sigil list -t sessions'
Index  Session ID       Target               Start Date
1      test-1234567890  i-xxxxxxxxxxxxxxxx1  2019-05-03T10:08:44Z
```
