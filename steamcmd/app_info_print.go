package steamcmd

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/frantjc/zoroark/encoding/vdf"
)

type AppInfo struct {
	Common   *AppInfoCommon   `json:"common,omitempty"`
	Extended *AppInfoExtended `json:"extended,omitempty"`
	Config   *AppInfoConfig   `json:"config,omitempty"`
}

type AppInfoCommon struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	OSList string `json:"oslist,omitempty"`
	GameID string `json:"gameid,omitempty"`
}

type AppInfoExtended struct {
	Developer                 string `json:"developer,omitempty"`
	GameDir                   string `json:"gamedir,omitempty"`
	Homepage                  string `json:"homepage,omitempty"`
	Icon                      string `json:"icon,omitempty"`
	NoServers                 string `json:"noservers,omitempty"`
	PrimaryCache              string `json:"primarycache,omitempty"`
	SourceGame                string `json:"sourcegame,omitempty"`
	State                     string `json:"state,omitempty"`
	VisibleOnlyWhenInstalled  string `json:"visibleonlywheninstalled,omitempty"`
	VisibleOnlyWhenSubscribed string `json:"visibleonlywhensubscribed,omitempty"`
}

type AppInfoConfig struct {
	Launch      map[string]AppInfoConfigLaunch `json:"launch,omitempty"`
	ContentType string                         `json:"contenttype,omitempty"`
	InstallDir  string                         `json:"installdir,omitempty"`
}

type AppInfoConfigLaunch struct {
	Executable string                     `json:"executable,omitempty"`
	Arguments  string                     `json:"arguments,omitempty"`
	Config     *AppInfoConfigLaunchConfig `json:"config,omitempty"`
}

type AppInfoConfigLaunchConfig struct {
	OSList string `json:"oslist"`
	OSArch string `json:"osarch"`
}

type AppInfoPrintCommand string

func (c AppInfoPrintCommand) String() string {
	return string(c)
}

var (
	appInfos map[string]AppInfo
)

func (c AppInfoPrintCommand) Check(flags *PromptFlags) error {
	return new(baseCommand).Check(flags)
}

func (c AppInfoPrintCommand) Args() ([]string, error) {
	if c == "" {
		return nil, fmt.Errorf("app_info_print requires app ID")
	} else if _, ok := appInfos[c.String()]; ok {
		return new(baseCommand).Args()
	}

	return []string{"app_info_print", c.String()}, nil
}

func (c AppInfoPrintCommand) ReadOutput(ctx context.Context, r io.Reader) error {
	if _, ok := appInfos[c.String()]; ok {
		return new(baseCommand).ReadOutput(ctx, r)
	}

	var (
		pr, pw    = io.Pipe()
		errC      = make(chan error, 1)
		buf       = new(bytes.Buffer)
		mw        = io.MultiWriter(pw, buf)
		errB      = []byte("ERROR! ")
		notFoundB = []byte("No app info for AppID")
	)

	go func() {
		defer close(errC)

		for {
			var b [4096]byte

			if _, err := r.Read(b[:]); err != nil {
				errC <- err
				return
			}

			if _, err := mw.Write(b[:]); err != nil {
				errC <- err
				return
			}

			p := buf.Bytes()
			if _, msgB, found := bytes.Cut(p, errB); found {
				msgB, _, _ = bytes.Cut(msgB, []byte("\n"))
				errC <- fmt.Errorf(string(msgB))
				return
			} else if bytes.Contains(p, notFoundB) {
				errC <- fmt.Errorf("app info for app ID %s not found", c.String())
				return
			}
		}
	}()

	go func() {
		defer close(errC)
		errC <- vdf.NewDecoder(pr).Decode(appInfos[c.String()])
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errC:
		return err
	}
}

func (c AppInfoPrintCommand) Modify(flags *PromptFlags) error {
	return new(baseCommand).Modify(flags)
}

func (c AppInfoPrintCommand) AppInfo() *AppInfo {
	if appInfo, ok := appInfos[c.String()]; ok {
		return &appInfo
	}

	return nil
}
