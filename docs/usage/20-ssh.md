# SSH

Start a new ssh for chosen EC2 instance based on its instance ID, name tag, or private DNS name.

```console
sigil ssh [flags]
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
    IdentityFile ~/.sigil/temp_key
    IdentitiesOnly yes
    ProxyCommand sigil ssh --target %h --port %p --pub-key "${HOME}"/.sigil/temp_key.pub --gen-key-pair --os-user %r
Host *.compute.internal
    IdentityFile ~/.sigil/temp_key
    IdentitiesOnly yes
    ProxyCommand sigil ssh --type private-dns --target %h --port %p --pub-key "${HOME}"/.sigil/temp_key.pub --gen-key-pair --os-user %r
```

```console
$ ssh ec2-user@ip-10-0-0-5.eu-west-1.compute.internal
Last login: Tue Jun 18 20:50:59 2019 from 10.0.0.5
...
[ec2-user@example ~]$
```
