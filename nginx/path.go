package nginx

import (
	"fmt"
)

func ArchivePath(version string) string {
	return fmt.Sprintf("%s.tar.gz", SourcePath(version))
}

func SourcePath(version string) string {
	return fmt.Sprintf("nginx-%s", version)
}
