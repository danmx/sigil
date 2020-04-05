## sigil list

List available EC2 instances or SSM sessions

### Synopsis

Show list of all EC2 instances with AWS SSM Agent running or active SSM sessions.

Supported groups of filters:
- instances:
	- tags - list of tag keys with a list of values for given keys
	- ids - list of instastance ids
- sessions:
	- after - the timestamp, in ISO-8601 Extended format, to see sessions that started after given date
	- before - the timestamp, in ISO-8601 Extended format, to see sessions that started before given date
	- target - an instance to which session connections have been made
	- owner - an AWS user account to see a list of sessions started by that user

Filter format examples:
[default.filters.session]
  after="2018-08-29T00:00:00Z"
  before="2019-08-29T00:00:00Z"
  target="i-xxxxxxxxxxxxxxxx1"
  owner="user@example.com"
[default.filters.instance]
  ids=["i-xxxxxxxxxxxxxxxx1","i-xxxxxxxxxxxxxxxx2"]
  tags=[{key="Name",values=["WebApp1","WebApp2"]}]


```
sigil list [flags]
```

### Examples

```
sigil list --output-format wide --instance-tags '[{"key":"Name","values":["Web","DB"]}]'
```

### Options

```
  -h, --help                             help for list
      --instance-ids strings             specify instance ids to limit results
      --instance-tags string             specify instance tags, in JSON format, to limit results
  -i, --interactive                      pick an instance or a session from a list and start or terminate the session
      --output-format string             specify output format: text/wide/json/yaml (default "text")
      --session-filters stringToString   specify session filters to limit results (default [after=,before=,target=,owner=])
  -t, --type string                      specify list type: instances/sessions (default "instances")
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

