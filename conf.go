package config

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type conf struct {
	filePath string
	groups   []*Group
}

func newConf(path string) *conf {
	return &conf{
		groups:   make([]*Group, 0),
		filePath: path,
	}
}

func (c *conf) loadConfiguration() error {
	filebytes, err := os.ReadFile(c.filePath)
	if err != nil {
		return fmt.Errorf("read file data failed: %v", err)
	}
	rendered := bytes.NewReader(filebytes)

	var group *Group
	scanner := bufio.NewScanner(rendered)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "[") && strings.HasSuffix(text, "]") && len(text) >= 3 {
			name := text[1 : len(text)-2]
			group = c.getGroup(name)
			continue
		}

		parts := strings.Split(text, "=")
		if len(parts) != 2 {
			return fmt.Errorf("the config message is error, does not support the format: %s", text)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if group != nil {
			v := strings.Trim(value, "\"")
			group.set(key, v)
		}
	}

	if scanner.Err() != nil {
		return err
	}

	return nil
}

func (c *conf) getGroup(name string) *Group {
	for _, g := range c.groups {
		if name == g.name {
			return g
		}
	}
	g := NewGroup(name)
	c.groups = append(c.groups, g)
	return g
}

func (c *conf) RegisterGroup(g *Group) {
	c.registerGroup(g)
}

func (c *conf) registerGroup(g *Group) {
	for _, gro := range c.groups {
		if gro.name == g.name {
			gro.copy(g)
			return
		}
	}
	c.groups = append(c.groups, g)
}

// if read failed, will return the empty string.
func (c *conf) ReadString(group, key string) string {
	for _, g := range c.groups {
		if g.name == group {
			return g.readString(key)
		}
	}
	return defaultString
}

// if read failed, ReadBool will return false.
func (c *conf) ReadBool(group, key string) bool {
	for _, g := range c.groups {
		if g.name == group {
			return g.readBool(key)
		}
	}
	return defaultBool
}

// if read failed, ReadInt will return -1.
func (c *conf) ReadInt(group, key string) int64 {
	for _, g := range c.groups {
		if g.name == group {
			return g.readInt(key)
		}
	}
	return defaultInt
}
