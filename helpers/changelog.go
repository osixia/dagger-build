package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ChangelogGetVersionSection(changelog, version string) (string, error) {
	f, err := os.Open(changelog)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Print(err.Error())
		}
	}()

	sectionHeaderPrefix := "## "
	sectionHeader := sectionHeaderPrefix + version

	var b strings.Builder
	inSection := false

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()

		if !inSection && strings.TrimSpace(line) == sectionHeader {
			inSection = true
			continue
		}

		if inSection {
			if strings.HasPrefix(line, sectionHeaderPrefix) {
				break
			}
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}
	if err := sc.Err(); err != nil {
		return "", err
	}

	if !inSection {
		return "", nil
	}

	return strings.TrimRight(b.String(), "\n"), nil
}
