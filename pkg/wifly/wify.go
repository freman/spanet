package wifly

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

type Command struct {
	c net.Conn
	s *bufio.Scanner
}

func NewCommand(c net.Conn) (*Command, error) {
	// Drain the sawmp
	buf := make([]byte, 1024)
	if _, err := c.Read(buf); err != nil {
		return nil, err
	}

	cmd := Command{c: c, s: bufio.NewScanner(c)}
	if os.Getenv("WIFLY_DEBUG") == "true" {
		cmd = Command{c: c, s: bufio.NewScanner(io.TeeReader(c, os.Stdout))}
	}

	cmd.s.Split(scanWifly)

	if err := cmd.command("$$$", "CMD"); err != nil {
		return nil, err
	}

	return &cmd, nil
}

func (c *Command) readln() (string, error) {
	c.s.Scan()
	if err := c.s.Err(); err != nil {
		return "", err
	}

	return strings.TrimSpace(c.s.Text()), nil
}

func (c *Command) writeln(s string) error {
	_, err := c.c.Write([]byte(s + "\r\n"))

	return err
}

func (c *Command) command(cmd string, expect ...string) error {
	if err := c.writeln(cmd); err != nil {
		return err
	}

	s, err := c.readln()
	if err != nil {
		return err
	}

	if len(expect) > 0 {
		if !strings.HasPrefix(s, expect[0]) {
			return fmt.Errorf("unexpected response: %s", s)
		}

		return nil
	}

	if !strings.HasPrefix(s, cmd) {
		return fmt.Errorf("unexpected response: %s", s)
	}

	return nil
}

func (c *Command) promptCommand(cmd string) (res []string, err error) {
	if err := c.command(cmd); err != nil {
		return nil, err
	}

	empty, err := c.readln()
	if err != nil {
		return nil, err
	}

	if empty != "" {
		return nil, fmt.Errorf("unexpected response: %s", empty)
	}

	for c.s.Scan() {
		line := strings.TrimSpace(c.s.Text())
		if line[0] == '<' && line[len(line)-1] == '>' {
			break
		}

		// capture first error but don't leave yet because we want to keep reading
		if strings.HasPrefix(line, "ERR:") && err != nil {
			err = errors.New(line)
		}
		res = append(res, line)
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Command) Get(what string) (map[string]string, error) {
	res, err := c.promptCommand("get " + what)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string, len(res))
	for _, r := range res {
		s := strings.Split(r, "=")
		if len(s) > 1 {
			m[s[0]] = strings.Join(s[1:], "=")
		}
	}

	return m, nil
}

func (c *Command) Set(what, value string) error {
	res, err := c.promptCommand("set " + what + " " + value)
	if err != nil {
		return err
	}

	if len(res) != 1 || res[0] != "AOK" {
		return errors.New("unexpected response")
	}

	return err
}

func (c *Command) Scan() (res []WIFIRecord, err error) {
	_, err = c.promptCommand("scan")
	if err != nil {
		return nil, err
	}

	for c.s.Scan() {
		line := strings.TrimSpace(c.s.Text())
		if strings.HasPrefix(line, "END:") {
			break
		}
		if strings.HasPrefix(line, "SCAN:Found") {
			continue
		}

		// capture first error but don't leave yet because we want to keep reading
		if strings.HasPrefix(line, "ERR:") && err != nil {
			err = errors.New(line)
			continue
		}

		rec, err := parseWifiRecord(line)
		if err != nil {
			continue
		}
		res = append(res, rec)
	}

	return res, err
}

func (c *Command) Save() error {
	_, err := c.promptCommand("save")
	if err != nil {
		return err
	}
	return nil
}

func (c *Command) Reboot() error {
	_, err := c.promptCommand("reboot")
	if err != nil {
		return err
	}

	return c.c.Close()
}

var reVersion = regexp.MustCompile(`<\d+\.\d+> `)

func scanWifly(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}

	if i := reVersion.FindIndex(data); i != nil {
		return i[1], dropCR(data[0:i[1]]), nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}

	// Request more data.
	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

type WIFIRecord struct {
	Index, Channel, RSSI, Auth int
	BSSID, SSID                string
}

// 01,06,-71,04,1114,28,c0,d4:35:1d:26:2d:31,Telstra262D31
// Index, Channel, RSSI, ?, ?, ?, ?, BSSID, SSID
func parseWifiRecord(str string) (rec WIFIRecord, err error) {
	var unknown int
	_, err = fmt.Sscanf(str,
		"%d,%d,%d,%d,%d,%d,%x,%17s",
		&rec.Index,
		&rec.Channel,
		&rec.RSSI,
		&rec.Auth,
		&unknown,
		&unknown,
		&unknown,
		&rec.BSSID,
	)

	rec.SSID = str[1+strings.LastIndex(str, ","):]

	return
}
