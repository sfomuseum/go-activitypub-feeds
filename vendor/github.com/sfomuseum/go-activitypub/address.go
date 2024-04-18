package activitypub

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
)

// copied from: https://github.com/mcnijman/go-emailaddress
// which in turn was copied from: https://stackoverflow.com/a/201378/5405453
// updated to add leading '@'

const pat_rfc5322 string = "(?i)@(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"

// To do: update to use pat_rfc5322

const pat_addr string = `(?:acct\:)?@?([^@]+)(?:@(.*))?`

var re_addr = regexp.MustCompile(fmt.Sprintf(`^%s$`, pat_addr))

var re_addresses = regexp.MustCompile(pat_rfc5322)

func ParseAddress(addr string) (string, string, error) {

	if !re_addr.MatchString(addr) {
		return "", "", fmt.Errorf("Failed to parse address")
	}

	m := re_addr.FindStringSubmatch(addr)
	return m[1], m[2], nil
}

func ParseAddressFromRequest(req *http.Request) (string, string, error) {

	resource := req.PathValue("resource")

	if resource == "" {
		return "", "", fmt.Errorf("request is missing {resource} path value")
	}

	slog.Debug("Parse address from request", "path", req.URL.Path, "resource", resource)
	return ParseAddress(resource)
}

func ParseAddressesFromString(body string) ([]string, error) {

	matches := re_addresses.FindAllStringSubmatch(body, -1)
	addresses := make([]string, len(matches))

	for idx, m := range matches {
		addresses[idx] = m[0]
	}

	return addresses, nil
}
