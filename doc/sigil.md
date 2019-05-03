## sigil

AWS SSM Session manager client

### Synopsis

A tool for establishing a session in EC2 instances with AWS SSM Agent installed

### Options

```
  -c, --config string           full config file path
  -p, --config-profile string   pick the config profile (default "default")
      --config-type string      specify the type of a config file: json, yaml, toml, hcl, props (default "toml")
  -h, --help                    help for sigil
      --log-level string        specify the log level: trace/debug/info/warn/error/fatal/panic (default "panic")
  -m, --mfa string              specify MFA token
      --profile string          specify AWS profile
  -r, --region string           specify AWS region
```

### SEE ALSO

* [sigil gendoc](sigil_gendoc.md)	 - Generate the documentation in Markdown
* [sigil list](sigil_list.md)	 - List available EC2 instances or SSM sessions
* [sigil session](sigil_session.md)	 - Start or resume a session
* [sigil verify](sigil_verify.md)	 - Verify if all external dependencies are available

