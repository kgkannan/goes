// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package bang

import (
	"fmt"
	"io"
	neturl "net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"

	"github.com/platinasystems/goes/cmd"
	"github.com/platinasystems/goes/lang"

	"github.com/platinasystems/url"
)

var tmpDir string

type Command struct{}

func (Command) String() string { return "!" }

func (Command) Usage() string {
	return "! COMMAND [ARGS]... [&]"
}

func (Command) Apropos() lang.Alt {
	return lang.Alt{
		lang.EnUS: "run an external command",
	}
}

func (Command) Man() lang.Alt {
	return lang.Alt{
		lang.EnUS: `
DESCRIPTION
	Sh-bang!

	Command executes in background if last argument ends with '&'.
	The standard i/o redirections apply.`,
	}
}

func (Command) Kind() cmd.Kind { return cmd.DontFork }

func (Command) Main(args ...string) error {
	var background bool

	if len(args) == 0 {
		return fmt.Errorf("COMMAND: missing")
	}

	tmpDir = "/var/run/goes/bang-" + strconv.Itoa(os.Getppid())
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return fmt.Errorf("Error in os.MkdirAll(%s): %w", tmpDir, err)
	}

	if n := len(args); args[n-1] == "&" {
		background = true
		args = args[:n-1]
	}

	filepath, u, err := url.FilePathFromUrl(args[0])
	if err != nil {
		return fmt.Errorf("Error from url.FilePathFromUrl(%s): %w",
			args[0], err)
	}
	execpath := args[0]
	command := args[0]
	if filepath == "" {
		execpath, command, err = loadNetExec(u)
		if err != nil {
			return fmt.Errorf("Error from loadNetExec(%v): %w",
				u, err)
		}
	}
	cmd := exec.Command(execpath, args[1:]...)
	cmd.Args[0] = command
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if background {
		go func() {
			err := cmd.Run()
			if err != nil {
				fmt.Fprintln(os.Stderr, cmd.Args[0], ": ", err)
			}
		}()
		return nil
	} else {
		return cmd.Run()
	}
}

func loadNetExec(u *neturl.URL) (execpath, command string, err error) {
	command = path.Base(u.String())
	execpath = filepath.Join(tmpDir, command)
	r, err := url.OpenUrl(u)
	if err != nil {
		return "", "", fmt.Errorf("Error from url.OpenUrl(%v): %w",
			u, err)
	}
	f, err := os.OpenFile(execpath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	if err != nil {
		return "", "", err
	}

	return
}
