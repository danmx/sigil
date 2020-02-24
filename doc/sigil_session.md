## sigil session

Start a session

### Synopsis

Start a new session in chosen EC2 instance.

```
sigil session [flags]
```

### Examples

```
sigil session --type instance-id --target i-xxxxxxxxxxxxxxxxx
```

### Options

```
  -h, --help            help for session
      --target string   specify the target depending on the type
      --type string     specify target type: instance-id/private-dns/name-tag (default "instance-id")
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

