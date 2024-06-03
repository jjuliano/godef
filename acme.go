package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"strconv"
	"strings"

	"9fans.net/go/acme"
)

// acmeFile represents an Acme file with its name, body, and cursor offsets.
type acmeFile struct {
	name       string
	body       []byte
	offset     int
	runeOffset int
}

// acmeCurrentFile retrieves the current file in Acme, including its name, body, and cursor offsets.
func acmeCurrentFile() (*acmeFile, error) {
	win, err := acmeCurrentWin()
	if err != nil {
		return nil, err
	}
	defer win.CloseFiles()

	// Ensure the address file is already open by reading it
	if _, _, err = win.ReadAddr(); err != nil {
		return nil, fmt.Errorf("cannot read address: %w", err)
	}

	// Set the address to dot (current selection)
	if err := win.Ctl("addr=dot"); err != nil {
		return nil, fmt.Errorf("cannot set addr=dot: %w", err)
	}

	q0, _, err := win.ReadAddr()
	if err != nil {
		return nil, fmt.Errorf("cannot read address: %w", err)
	}

	body, err := readBody(win)
	if err != nil {
		return nil, fmt.Errorf("cannot read body: %w", err)
	}

	tag, err := readTag(win)
	if err != nil {
		return nil, fmt.Errorf("cannot read tag: %w", err)
	}

	name, err := extractNameFromTag(tag)
	if err != nil {
		return nil, err
	}

	return &acmeFile{
		name:       name,
		body:       body,
		offset:     runeOffsetToByteOffset(body, q0),
		runeOffset: q0,
	}, nil
}

// readBody reads the entire content of the body of an Acme window.
func readBody(win *acme.Win) ([]byte, error) {
	var body []byte
	buf := make([]byte, 8000)
	for {
		n, err := win.Read("body", buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		body = append(body, buf[:n]...)
	}
	return body, nil
}

// readTag reads the tag content of an Acme window.
func readTag(win *acme.Win) (string, error) {
	tagb, err := win.ReadAll("tag")
	if err != nil {
		return "", err
	}
	return string(tagb), nil
}

// extractNameFromTag extracts the name from the Acme window tag.
func extractNameFromTag(tag string) (string, error) {
	i := strings.Index(tag, " ")
	if i == -1 {
		return "", fmt.Errorf("invalid tag format: no spaces found")
	}
	return tag[:i], nil
}

// acmeCurrentWin opens the current Acme window using the winid environment variable.
func acmeCurrentWin() (*acme.Win, error) {
	winid := os.Getenv("winid")
	if winid == "" {
		return nil, fmt.Errorf("$winid not set - not running inside acme?")
	}

	id, err := strconv.Atoi(winid)
	if err != nil {
		return nil, fmt.Errorf("invalid $winid %q: %w", winid, err)
	}

	if err := setNameSpace(); err != nil {
		return nil, err
	}

	win, err := acme.Open(id, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot open acme window: %w", err)
	}
	return win, nil
}

// runeOffsetToByteOffset converts a rune offset to a byte offset within a byte slice.
func runeOffsetToByteOffset(b []byte, off int) int {
	r := 0
	for i := range string(b) {
		if r == off {
			return i
		}
		r++
	}
	return len(b)
}

// setNameSpace sets the NAMESPACE environment variable if not already set.
func setNameSpace() error {
	if ns := os.Getenv("NAMESPACE"); ns != "" {
		return nil
	}

	ns, err := nsFromDisplay()
	if err != nil {
		return fmt.Errorf("cannot get namespace: %w", err)
	}

	os.Setenv("NAMESPACE", ns)
	return nil
}

// nsFromDisplay retrieves the namespace based on the DISPLAY environment variable.
func nsFromDisplay() (string, error) {
	disp := os.Getenv("DISPLAY")
	if disp == "" {
		disp = ":0.0"
	}

	if i := strings.LastIndex(disp, ":"); i >= 0 && strings.HasSuffix(disp, ".0") {
		disp = disp[:len(disp)-2]
	}

	disp = strings.ReplaceAll(disp, "/", "_")

	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("cannot get current user name: %w", err)
	}

	ns := fmt.Sprintf("/tmp/ns.%s.%s", u.Username, disp)
	if _, err := os.Stat(ns); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("namespace directory does not exist")
		}
		return "", fmt.Errorf("cannot stat namespace directory: %w", err)
	}
	return ns, nil
}
