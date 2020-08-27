# SSH

Start a new ssh for chosen EC2 instance based on its instance ID, name tag, or private DNS name.

```console
ssh [--type TYPE] ... [ { --gen-key-pair [--gen-key-dir DIR] | --pub-key PUB_KEY_PATH } ] TARGET
```

[Man](../man/sigil_ssh.md) page

## Sample config

Config file settings that affect the command

```toml
[default]
  type = "name"
  target = "Worker"
  region = "eu-west-1"
  profile = "dev"
  os-user = "ec2-user"
  gen-key-pair = false
  pub-key = "~/.ssh/dev.pub"
```

## Examples

`ssh_config` config file example:

```ssh_config
Host i-* mi-*
    IdentityFile /tmp/sigil/%h/temp_key
    IdentitiesOnly yes
    ProxyCommand sigil ssh --port %p --pub-key /tmp/sigil/%h/temp_key.pub --gen-key-pair --os-user %r --gen-key-dir /tmp/sigil/%h/ %h
Host *.compute.internal
    IdentityFile /tmp/sigil/%h/temp_key
    IdentitiesOnly yes
    ProxyCommand sigil ssh --type private-dns --port %p --pub-key /tmp/sigil/%h/temp_key.pub --gen-key-pair --os-user %r --gen-key-dir /tmp/sigil/%h/ %h
```

```console
$ ssh ec2-user@ip-10-0-0-5.eu-west-1.compute.internal
Last login: Tue Jun 18 20:50:59 2019 from 10.0.0.5
...
[ec2-user@example ~]$
```
