package gen

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
)

func TempFilename(prefix string) string {
	return fmt.Sprintf("%s/%s%s", os.TempDir(), prefix, uuid.NewString())
}

func CopyFile(from string, to string) error {
	in, err := os.Open(from)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

// JoinNonzeroStrings concatenates the given string slice using sep as the
// delimiter. It ignores empty strings in ss.
func JoinNonzeroStrings(ss []string, sep string) string {
	sb := strings.Builder{}

	for _, s := range ss {
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteString(sep)
			}
			sb.WriteString(s)
		}
	}

	return sb.String()
}
