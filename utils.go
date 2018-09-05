package gitreport

import (
	"errors"
	"os"
	"os/user"
)

func ResolvePath(path string) (string, error) {
	if path[0] == '.' {
		return "", errors.New("Must use absolute path for repo")
	}

	if path[0] == '~' {
		return resolveHomePath(path), nil
	}

	return path, nil
}

func resolveHomePath(path string) string {
	if path == "~" {
		return getUserHome()
	}
	return getUserHome() + path[1:]
}

func getUserHome() string {
	if os.Getenv("HOME") != "" {
		return os.Getenv("HOME")
	}

	usr, _ := user.Current()
	return usr.HomeDir
}
