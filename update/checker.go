package update

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/Alia5/SISR/config"
	"github.com/Alia5/SISR/meta"
)

const releasesAPIURL = "https://api.github.com/repos/Alia5/SISR/releases"

var versionRegex = regexp.MustCompile(`v(\d+)\.(\d+)\.(\d+)(?:-(\d+)-g[0-9a-f]+)?`)

type Checker interface {
	GetVersionInfo() *VersionInfo
	CheckForUpdate(ctx context.Context) (*VersionInfo, error)
}

func NewChecker(updateChannel config.UpdateNotify) Checker {
	return &checker{
		updateChannel: updateChannel,
	}
}

type checker struct {
	updateAvailable bool
	newVersion      string

	updateChannel config.UpdateNotify

	mtx sync.Mutex
}

type VersionInfo struct {
	Version         string
	Commit          string
	Date            string
	UpdateAvailable bool
	NewVersion      string
}

type release struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Prerelease bool   `json:"prerelease"`
	HTMLURL    string `json:"html_url"`
}

func (c *checker) GetVersionInfo() *VersionInfo {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return &VersionInfo{
		Version:         meta.Version,
		Commit:          meta.Commit,
		Date:            meta.Date,
		UpdateAvailable: c.updateAvailable,
		NewVersion:      c.newVersion,
	}
}

func (c *checker) CheckForUpdate(ctx context.Context) (*VersionInfo, error) {
	cur, ok := parseVersion(meta.Version)
	if !ok && meta.Version != "dev" {
		slog.Error("failed to parse current version", "version", meta.Version)
		return &VersionInfo{
			Version:         meta.Version,
			Commit:          meta.Commit,
			Date:            meta.Date,
			UpdateAvailable: c.updateAvailable,
			NewVersion:      c.newVersion,
		}, nil
	}
	var r release
	client := &http.Client{Timeout: 10 * time.Second}

	c.mtx.Lock()
	notify := c.updateChannel
	c.mtx.Unlock()

	if notify == config.UpdateNotifyPrerelease {
		resp, err := client.Get(releasesAPIURL + "?per_page=1")
		if err != nil {
			slog.Error("failed to fetch releases", "error", err)
			return nil, err
		}
		defer resp.Body.Close() //nolint:errcheck
		var releases []release
		if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
			slog.Error("failed to decode releases", "error", err)
			return nil, err
		}
		if len(releases) == 0 {
			return &VersionInfo{
				Version:         meta.Version,
				Commit:          meta.Commit,
				Date:            meta.Date,
				UpdateAvailable: c.updateAvailable,
				NewVersion:      c.newVersion,
			}, nil
		}
		r = releases[0]
	} else {
		resp, err := client.Get(releasesAPIURL + "/latest")
		if err != nil {
			slog.Error("failed to fetch latest release", "error", err)
			return nil, err
		}
		defer resp.Body.Close() //nolint:errcheck
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			slog.Error("failed to decode latest release", "error", err)
			return nil, err
		}
	}

	versionSource := r.TagName
	if r.Prerelease {
		versionSource = r.Name
	}

	remote, ok := parseVersion(versionSource)
	if !ok {
		slog.Error("failed to parse remote version", "version", versionSource)
		return nil, errors.New("failed to parse remote version")
	}

	newer := remote.Major > cur.Major ||
		(remote.Major == cur.Major && remote.Minor > cur.Minor) ||
		(remote.Major == cur.Major && remote.Minor == cur.Minor && remote.Patch > cur.Patch) ||
		(remote.Major == cur.Major && remote.Minor == cur.Minor && remote.Patch == cur.Patch && remote.Commits > cur.Commits)

	if !newer {
		return &VersionInfo{
			Version:         meta.Version,
			Commit:          meta.Commit,
			Date:            meta.Date,
			UpdateAvailable: false,
			NewVersion:      "",
		}, nil
	}

	matched := versionRegex.FindString(versionSource)

	c.mtx.Lock()
	c.updateAvailable = true
	c.newVersion = matched
	c.mtx.Unlock()

	return &VersionInfo{
		Version:         meta.Version,
		Commit:          meta.Commit,
		Date:            meta.Date,
		UpdateAvailable: true,
		NewVersion:      matched,
	}, nil
}

type version struct {
	Major, Minor, Patch, Commits int
}

func parseVersion(s string) (version, bool) {
	m := versionRegex.FindStringSubmatch(s)
	if m == nil {
		return version{}, false
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])
	commits := 0
	if m[4] != "" {
		commits, _ = strconv.Atoi(m[4])
	}
	return version{major, minor, patch, commits}, true
}
