package filesystem

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

type perm struct {
}

func (s *perm) GetOwner(file string) (string, error) {
	info, err := os.Lstat(file)
	if err != nil {
		return "", err
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uID := strconv.Itoa(int(stat.Uid))
		user, err := user.LookupId(uID)
		if err != nil {
			return uID, nil
		}
		return user.Username, nil
	}

	return "unknown", nil
}

func (s *perm) GetGroup(file string) (string, error) {
	info, err := os.Lstat(file)
	if err != nil {
		return "", err
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		gID := strconv.Itoa(int(stat.Gid))
		group, err := user.LookupGroupId(gID)
		if err != nil {
			return gID, nil
		}
		return group.Name, nil
	}

	return "unknown", nil
}

func (s *perm) GetMode(file string) (os.FileMode, error) {
	info, err := os.Lstat(file)
	if err != nil {
		return 0, err
	}
	return info.Mode(), nil
}

func (s *perm) ChOwner(string, string) error {
	return nil
}

func (s *perm) ChGroup(string, string) error {
	return nil
}

func (s *perm) ChMode(string, os.FileMode) error {
	return nil
}
