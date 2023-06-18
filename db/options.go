package db

type Purge uint

const (
	Nil Purge = iota
	Yes
	No
)

type MariaDBOptions struct {
	Username string
	Password string
	Database string
	Host     string
	Port     int
	Other    string

	LogLevel int
	LogPath  string

	Purge Purge
}

// This sucks
func (o *MariaDBOptions) validate() {
	if o.Username == "" {
		o.Username = "wc"
	}
	if o.Password == "" {
		o.Password = "password"
	}
	if o.Host == "" {
		o.Host = "127.0.0.1"
	}
	if o.Port == 0 {
		o.Port = 3306
	}
	if o.Database == "" {
		o.Database = "wc"
	}
	if o.Other == "" {
		o.Other = "charset=utf8&parseTime=True"
	}
	if o.LogLevel == 0 {
		o.LogLevel = 4
	}
	if o.LogPath == "" {
		o.LogPath = "<stdout>"
	}
}

type SqliteDBOptions struct {
	Path   string
	Memory bool
	Other  string

	LogLevel int
	LogPath  string

	Purge Purge
}

func (o *SqliteDBOptions) validate() {
	if o.Memory {
		o.Path = ":memory:"
	}
	if o.Path == "" {
		o.Path = "storage/test.db"
	}
	if o.LogLevel == 0 {
		o.LogLevel = 4
	}
	if o.LogPath == "" {
		o.LogPath = "<stdout>"
	}
}
