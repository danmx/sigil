package aws

//go:generate go run github.com/golang/mock/mockgen -self_package=github.com/danmx/sigil -destination mocks/aws_mocks.go -package=mocks github.com/danmx/sigil/pkg/aws Cloud,CloudInstances,CloudSessions,CloudSSH

// Cloud wraps init methods used from the aws package
type Cloud interface {
	NewWithConfig(c *Config) error
}

// CloudInstances wraps instances methods used from the aws package
type CloudInstances interface {
	Cloud
	ListInstances() ([]*Instance, error)
	StartSession(targetType, target string) error
}

// CloudSessions wraps sessions methods used from the aws package
type CloudSessions interface {
	Cloud
	ListSessions() ([]*Session, error)
	TerminateSession(sessionID string) error
}

// CloudSSH wraps ssh methods used from the aws package
type CloudSSH interface {
	Cloud
	StartSSH(targetType, target, osUser string, portNumber uint64, publicKey []byte) error
}
