package installable

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"path"
	"runtime"
	"strings"

	"github.com/codeclysm/extract/v3"
)

const HTTPArchiveType = "http:archive"

func newHTTPArchive(entry Entry, opt *HTTPArchiveOption) *HTTPArchive {
	return &HTTPArchive{
		Name:    entry.Name,
		Version: entry.Version,
		Source:  entry.Source,

		option:    *opt,
		versioned: entry.Name + "@" + entry.Version,
	}
}

type HTTPArchiveOption struct {
	StripPrefix string `yaml:"stripPrefix"`

	Overrides struct {
		// TODO(dio): Make it typed.
		OS   map[string]string `yaml:"os"`
		Arch map[string]string `yaml:"arch"`
		Ext  map[string]string `yaml:"ext"`
	} `yaml:"overrides"`

	// TODO(dio): Have a way to set main binary and put it in a "bin" directory.
	// This is for the case when an archive doesn't have "bin" directory, or the
	// main binary is not in the "bin" directory.

	SHAs map[string]string `yaml:"shas"`

	CI string `yaml:"ci"`
}

type HTTPArchive struct {
	Name    string
	Version string
	Source  string

	option    HTTPArchiveOption
	versioned string
}

func (a *HTTPArchive) Install(dst string) (string, error) {
	versioned := path.Join(dst, a.versioned)
	installed := path.Join(versioned, "bin")

	if err := checkInstalled(dst, a.Name, a.versioned, a.option.CI); err != nil {
		if err == ErrAlreadyInstalled {
			return installed, nil
		}
		return installed, err
	}

	source, err := a.expand(a.Name+":url", a.Source)
	if err != nil {
		return installed, err
	}
	data, _, err := readRemoteFile(source)
	if err != nil {
		return installed, err
	}

	if err := a.checksum(data); err != nil {
		return installed, err
	}

	br := bufio.NewReader(bytes.NewBuffer(data))
	prefix, err := a.expand(a.Name+":stripPrefix", a.option.StripPrefix)
	if err != nil {
		return installed, err
	}

	return installed, extract.Archive(context.Background(), br, versioned, func(s string) string {
		return strings.TrimPrefix(s, prefix)
	})
}

func (a *HTTPArchive) checksum(data []byte) error {
	// TODO(dio): Add checksum.
	name := runtime.GOOS + "-" + runtime.GOARCH
	value := infer(a.option.SHAs, name, "")
	if value == "" {
		return ErrInvalid
	}

	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("failed to checksum: %w", ErrInvalid)
	}

	h := sha256.New()
	h.Write(data)
	sum := h.Sum(nil)
	if hex.EncodeToString(sum) != parts[1] {
		return ErrInvalid
	}
	return nil
}

func (a *HTTPArchive) expand(name, text string) (string, error) {
	u, err := template.New(name).Parse(text)
	if err != nil {
		return "", err
	}
	var rendered bytes.Buffer
	if err = u.Execute(&rendered, map[string]string{
		"Version": a.Version,
		"OS":      infer(a.option.Overrides.OS, runtime.GOOS, runtime.GOOS),
		"Arch":    infer(a.option.Overrides.Arch, runtime.GOARCH, runtime.GOARCH),
		"Ext":     infer(a.option.Overrides.Ext, runtime.GOOS, ".tar.gz"), // We default to .tar.gz
	}); err != nil {
		return "", err
	}
	return rendered.String(), nil
}

func (n *HTTPArchive) Runtime() Installable {
	return nil
}
