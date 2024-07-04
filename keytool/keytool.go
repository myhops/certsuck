package keytool

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	Command = "keytool"
)

var (
	ErrNotEnoughAliases = errors.New("not enough aliases")
)

// IsInstalled return true if keytool is present on the system.
type Runner interface {
	RunContext(context.Context)
}

type ImportCerts struct {
	Aliases   []string
	Keystore  string
	Files     []string
	Keypass   string
	Storepass string
}

type ImportCert struct {
	Alias     string
	Keystore  string
	File      string
	Keypass   string
	Storepass string
}

func (c ImportCert) RunContext(ctx context.Context) error {
	// Build and execute the command.
	if len(c.Alias) == 0 {
		c.Alias = aliasFromFile(c.File)
	}
	cmd := exec.CommandContext(ctx, Command)
	cmd.Args = []string{
		Command,
		"-importcert",
		"-keystore", c.Keystore,
		"-alias", c.Alias,
		"-file", c.File,
		"-storepass", c.Storepass,
		"-keypass", c.Keypass,
	}
	// Connect the stdin
	cmd.Stdin = NewYesReader()
	// Collect std out
	out := bytes.Buffer{}
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Start()
	if err != nil {
		return err
	}
	// Wait for cmd to end
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func aliasFromFile(file string) string {
	return strings.ReplaceAll(filepath.Base(file), " ", "_")
}

func (c ImportCerts) RunContext(ctx context.Context) error {
	cmd := &ImportCert{
		Keystore:  c.Keystore,
		Keypass:   c.Keypass,
		Storepass: c.Storepass,
	}
	aliases := c.Aliases
	if len(aliases) > 0 && len(aliases) < len(c.Files) {
		return ErrNotEnoughAliases
	}
	if len(aliases) == 0 {
		for _, file := range c.Files {
			aliases = append(aliases, aliasFromFile(file))
		}
	}

	for i := range c.Files {
		cmd.Alias = aliases[i]
		cmd.File = c.Files[i]
		err := cmd.RunContext(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
