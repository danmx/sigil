# Configuration

The configuration varies depending on the command. For more defails (like default values) check [usage](usage/README.md) section or [man](man/sigil.md) pages.

## Config file

Description of different values of the configuration file.

| Parameter                |   Command(s)    | Description                                                                  |
| ------------------------ | :-------------: | :--------------------------------------------------------------------------- |
| `type`                   | `session`/`ssh` | specify target type                                                          |
| `target`                 | `session`/`ssh` | specify the target depending on the type                                     |
| `os-user`                |      `ssh`      | specify an instance OS user which will be using sent public key              |
| `port`                   |      `ssh`      | specify ssh port                                                             |
| `gen-key-pair`           |      `ssh`      | generate a temporary key pair that will be send and used                     |
| `gen-key-dir`            |      `ssh`      | the directory where temporary keys will be generated                         |
| `pub-key`                |      `ssh`      | local public key that will be send to the instance                           |
| `output-format`          |     `list`      | specify output format                                                        |
| `interactive`            |     `list`      | pick an instance or a session from a list and start or terminate the session |
| `filters.session.after`  |     `list`      | show only sessions that started after given datetime                         |
| `filters.session.before` |     `list`      | show only sessions that started before given datetime                        |
| `filters.session.target` |     `list`      | show only sessions for given target                                          |
| `filters.session.owner`  |     `list`      | show only sessions owned by given owner                                      |
| `filters.instance.ids`   |     `list`      | show only instances with matching IDs                                        |
| `filters.instance.tags`  |     `list`      | show only instances with matching tags                                       |

## Example

An example of a fully configured `default` profile.

```toml
[default]
  type = "name"
  target = "Worker"
  type = "instance-id"
  output-format = "text"
  region = "eu-west-1"
  profile = "dev"
  interactive = false
  os-user = "ec2-user"
  gen-key-pair = false
  gen-key-dir = "/tmp/sigil"
  pub-key = "~/.ssh/dev.pub"
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
