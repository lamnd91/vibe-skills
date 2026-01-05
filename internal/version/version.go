package version

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func GetVersion() string {
	return Version
}

func GetFullVersion() string {
	return Version + " (" + Commit + ") built at " + Date
}
