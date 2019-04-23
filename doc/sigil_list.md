## sigil list

List available EC2 instances

### Synopsis

Show list of all EC2 instances with AWS SSM Agent running.

```
sigil list [flags]
```

### Examples

```
sigil list --output-format wide -t Name=webapp1
```

### Options

```
  -h, --help                   help for list
  -i, --interactive            pick an instance from a list and start the session
      --output-format string   specify output format: text/json/yaml/wide (default "text")
  -t, --tags stringToString    specify tags to filter out results, e.g.: key1=value1,key2=value2 (default [])
```

### Options inherited from parent commands

```
  -c, --config string           full config file path
  -p, --config-profile string   pick the config profile (default "default")
      --config-type string      specify the type of a config file: json, yaml, toml, hcl, props (default "toml")
      --log-level string        specify the log level: trace/debug/info/warn/error/fatal/panic (default "debug")
  -m, --mfa string              specify MFA token
      --profile string          specify AWS profile
  -r, --region string           specify AWS region
```

### SEE ALSO

* [sigil](sigil.md)	 - AWS SSM Session manager client

