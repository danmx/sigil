## sigil ssh

Start ssh session

### Synopsis

Start a new ssh for chosen EC2 instance.

```
sigil ssh [flags]
```

### Options

```
      --gen-key-pair     generate a temporary key pair that will be send and used. Use ${HOME}/.sigil/temp_key as an identity file
  -h, --help             help for ssh
      --os-user string   specify an instance OS user which will be using sent public key (default "ec2-user")
      --port int         specify ssh port (default 22)
      --pub-key string   local public key that will be send to the instance, ignored when gen-key-pair is true
      --target string    specify the target depedning on the type
      --type string      specify target type: instance-id/private-dns/name-tag (default "instance-id")
```

### Options inherited from parent commands

```
  -c, --config string           full config file path
  -p, --config-profile string   pick the config profile (default "default")
      --config-type string      specify the type of a config file: json, yaml, toml, hcl, props (default "toml")
      --log-level string        specify the log level: trace/debug/info/warn/error/fatal/panic (default "panic")
  -m, --mfa string              specify MFA token
      --profile string          specify AWS profile
  -r, --region string           specify AWS region
```

### SEE ALSO

* [sigil](sigil.md)	 - AWS SSM Session manager client

