package config

import (
	"fmt"
	"strings"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
)

// One spelling only, so a document cannot name the instance twice with no rule for which wins.
const clientNameKeyCanon = "client_name"

// UnnamedIndexPattern labels a credentials entry that carries no client name.
const UnnamedIndexPattern = "<index %d>"

// ReadClientName returns one credentials entry's client name, whatever shape the entry has.
func ReadClientName(block any) string {
	switch typed := block.(type) {
	case SdkConfig:
		return typed.Auth.ClientName
	case *SdkConfig:
		if typed == nil {
			return ""
		}
		return typed.Auth.ClientName
	case AuthConfig:
		return typed.ClientName
	case *AuthConfig:
		if typed == nil {
			return ""
		}
		return typed.ClientName
	case map[string]any:
		if value, ok := typed[clientNameKeyCanon]; ok && value != nil {
			return fmt.Sprint(value)
		}
		return ""
	case map[string]string:
		return typed[clientNameKeyCanon]
	default:
		return ""
	}
}

// SelectCredentials returns one entry: nil is the first, an int a position (negative from the end), a string a name.
func SelectCredentials(credentials []any, instance any) (any, error) {
	if len(credentials) == 0 {
		return nil, wcerrors.NewUsageError("no credentials were supplied")
	}
	if instance == nil {
		return credentials[0], nil
	}
	switch typed := instance.(type) {
	case int:
		size := len(credentials)
		if typed >= -size && typed < size {
			return credentials[((typed%size)+size)%size], nil
		}
		return nil, wcerrors.NewUsageError(fmt.Sprintf(
			"no credentials at index %d; %d supplied (0-%d)", typed, size, size-1))
	case string:
		return selectByName(credentials, typed)
	default:
		return nil, wcerrors.NewUsageError(
			"instance must be nil, an index or a clientName")
	}
}

func selectByName(credentials []any, name string) (any, error) {
	var matched []int
	for index, block := range credentials {
		if ReadClientName(block) == name {
			matched = append(matched, index)
		}
	}
	if len(matched) > 1 {
		positions := make([]string, 0, len(matched))
		for _, index := range matched {
			positions = append(positions, fmt.Sprintf("%d", index))
		}
		return nil, wcerrors.NewUsageError(fmt.Sprintf(
			"instance '%s' is defined more than once, at index %s; clientName has to identify "+
				"one entry", name, strings.Join(positions, " and ")))
	}
	if len(matched) == 1 {
		return credentials[matched[0]], nil
	}
	available := make([]string, 0, len(credentials))
	for index, block := range credentials {
		label := ReadClientName(block)
		if label == "" {
			label = fmt.Sprintf(UnnamedIndexPattern, index)
		}
		available = append(available, label)
	}
	return nil, wcerrors.NewUsageError(fmt.Sprintf(
		"unknown instance '%s'; available: %s", name, strings.Join(available, ", ")))
}
