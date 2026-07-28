package match

import (
	"fmt"

	"github.com/bmatcuk/doublestar/v4"
)

// MatchWatchPaths returns true if any of the changedFiles match at least one of
// the given watch path patterns. Patterns use doublestar syntax (e.g. "src/**",
// "*.go", "services/api/**/*.ts").
//
// If patterns is empty the function returns true (no filtering).
// If changedFiles is empty the function returns false (nothing to match).
func MatchWatchPaths(patterns []string, changedFiles []string) bool {
	matched, err := MatchWatchPathsErr(patterns, changedFiles)
	if err != nil {
		// Preserve the legacy behavior for callers that do not need to surface
		// invalid watch path configuration yet.
		return false
	}

	return matched
}

// MatchWatchPathsErr returns true if any changed file matches a watch path
// pattern and returns an error when any configured pattern is invalid.
func MatchWatchPathsErr(patterns []string, changedFiles []string) (bool, error) {
	if err := ValidateWatchPaths(patterns); err != nil {
		return false, err
	}

	if len(patterns) == 0 {
		return true, nil
	}
	if len(changedFiles) == 0 {
		return false, nil
	}

	for _, file := range changedFiles {
		for _, pattern := range patterns {
			matched, err := doublestar.Match(pattern, file)
			// Pattern validity is checked above, so any error here means this
			// pattern/file pair is malformed in a way validation missed.
			if err != nil {
				return false, fmt.Errorf("match watch path %q against %q: %w", pattern, file, err)
			}
			if matched {
				return true, nil
			}
		}
	}
	return false, nil
}

// ValidateWatchPaths returns an error if any watch path pattern is invalid.
func ValidateWatchPaths(patterns []string) error {
	for _, pattern := range patterns {
		if _, err := doublestar.Match(pattern, ""); err != nil {
			return fmt.Errorf("invalid watch path pattern %q: %w", pattern, err)
		}
	}

	return nil
}
