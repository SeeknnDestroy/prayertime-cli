package version

var (
	commit = "dev"
	date   = "unknown"
	tag    = "v0.1.0-dev"
)

func String() string {
	return tag + " (" + commit + ", " + date + ")"
}
