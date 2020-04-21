package seafile

import (
	"os"
	"strings"
	"time"
)

type Library struct {
	OwnerContactEmail     string `json:"owner_contact_email,omitempty"`
	OwnerName             string `json:"owner_name,omitempty"`
	ModifierEmail         string `json:"modifier_email"`
	FName                 string `json:"name"`
	Permission            string `json:"permission"`
	SizeFormatted         string `json:"size_formatted,omitempty"`
	Virtual               bool   `json:"virtual,omitempty"`
	MtimeRelative         string `json:"mtime_relative"`
	HeadCommitID          string `json:"head_commit_id"`
	Encrypted             bool   `json:"encrypted"`
	Version               int    `json:"version"`
	Mtime                 int64  `json:"mtime"`
	FOwner                string `json:"owner"`
	ModifierContactEmail  string `json:"modifier_contact_email"`
	Root                  string `json:"root"`
	Type                  string `json:"type"`
	ID                    string `json:"id"`
	ModifierName          string `json:"modifier_name"`
	FSize                 int64  `json:"size"`
	ShareFromName         string `json:"share_from_name,omitempty"`
	ShareFrom             string `json:"share_from,omitempty"`
	ShareFromContactEmail string `json:"share_from_contact_email,omitempty"`
	Groupid               int    `json:"groupid,omitempty"`
	GroupName             string `json:"group_name,omitempty"`
}

func (f Library) Name() string {
	return f.FName
}
func (f Library) Size() int64 {
	return int64(f.FSize)
}
func (f Library) Mode() (m os.FileMode) {
	m = 0111
	if strings.Contains(f.Permission, "r") {
		m = m | 0444
	}
	if strings.Contains(f.Permission, "w") {
		m = m | 0222
	}
	return m
}
func (f Library) ModTime() time.Time {
	return time.Unix(f.Mtime, 0)
}
func (f Library) IsDir() bool {
	return true
}
func (f Library) Sys() interface{} {
	return nil
}

func (f Library) Owner() string {
	if f.OwnerName == "" {
		return f.ShareFromName
	}
	return f.OwnerName
}

func (f Library) Group() string {
	return ""
}

type directoryEntry struct {
	Permission           string `json:"permission"`
	Mtime                int64  `json:"mtime"`
	Type                 string `json:"type"`
	FName                string `json:"name"`
	ID                   string `json:"id"`
	ModifierEmail        string `json:"modifier_email,omitempty"`
	ModifierContactEmail string `json:"modifier_contact_email,omitempty"`
	Starred              bool   `json:"starred,omitempty"`
	ModifierName         string `json:"modifier_name,omitempty"`
	FSize                int64  `json:"size,omitempty"`
}

func (f directoryEntry) Name() string {
	return f.FName
}
func (f directoryEntry) Size() int64 {
	return int64(f.FSize)
}
func (f directoryEntry) Mode() (m os.FileMode) {
	m = 0111
	if strings.Contains(f.Permission, "r") {
		m = m | 0444
	}
	if strings.Contains(f.Permission, "w") {
		m = m | 0222
	}
	return m
}
func (f directoryEntry) ModTime() time.Time {
	return time.Unix(f.Mtime, 0)
}
func (f directoryEntry) IsDir() bool {
	return f.Type == "dir"
}
func (f directoryEntry) Sys() interface{} {
	return nil
}
func (f directoryEntry) Owner() string {
	return f.ModifierName
}

func (f directoryEntry) Group() string {
	return ""
}

type directoryDetail struct {
	Path   string    `json:"path"`
	RepoID string    `json:"repo_id"`
	FName  string    `json:"name"`
	Mtime  time.Time `json:"mtime"`
}

func (f directoryDetail) Name() string {
	return f.FName
}
func (f directoryDetail) Size() int64 {
	return 0
}
func (f directoryDetail) Mode() (m os.FileMode) {
	return 0777
}
func (f directoryDetail) ModTime() time.Time {
	return f.Mtime
}
func (f directoryDetail) IsDir() bool {
	return true
}
func (f directoryDetail) Sys() interface{} {
	return nil
}
func (f directoryDetail) Owner() string {
	return ""
}
func (f directoryDetail) Group() string {
	return ""
}
