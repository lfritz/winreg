package main

// A command-line tool for Windows that lists shared folders.

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func main() {
	// To get the share folders exported by Windows we need to use the Registry, open
	// Computer\HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\LanmanServer\Shares
	// and look for entries with Type=0.

	// open Shares key
	path := `SYSTEM\CurrentControlSet\Services\LanmanServer\Shares`
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		fmt.Printf("error opening Windows Registry key: %v\n", err)
		os.Exit(1)
	}
	defer key.Close()

	// get value names
	names, err := key.ReadValueNames(-1)
	if err != nil {
		fmt.Printf("error reading Windows Registry value names: %v\n", err)
		os.Exit(1)
	}

	// print share name and path for each
	for _, name := range names {
		values, _, err := key.GetStringsValue(name)
		if err != nil {
			fmt.Printf("error reading Windows Registry value: %v\n", err)
			os.Exit(1)
		}
		m := mapFromStrings(values)
		if m["Type"] != "0" {
			// Type=0 means network drives; if it's not 0 it could be, for example, a shared printer
			continue
		}
		fmt.Printf("%s -> %s\n", name, m["Path"])
	}
}

// mapFromStrings takes a slice of mappings of the form "key=value" and turns it into a map.
func mapFromStrings(values []string) map[string]string {
	result := make(map[string]string)
	for _, v := range values {
		key, value, ok := strings.Cut(v, "=")
		if !ok {
			continue
		}
		result[key] = value
	}
	return result
}
