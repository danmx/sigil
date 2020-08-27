# Session

Start a new session in chosen EC2 instance based on its instance ID, name tag, or private DNS name.

```console
sigil session [--type TYPE] ... TARGET
```

[Man](../man/sigil_session.md) page

## Sample config

Config file settings that affect the command

```toml
[default]
  type = "name"
  target = "Worker"
  region = "eu-west-1"
  profile = "dev"
```

## Examples

```console
$ sigil -r eu-west-1 session --type instance-id i-xxxxxxxxxxxxxxxxx
Starting session with SessionId: example
sh-4.2$
```
