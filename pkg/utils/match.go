package util

import "strings"

// TODO: move this to constants file
const (
	OK_TO_TEST string = "/ok-to-test"
	TEST       string = "/test"
	RETEST     string = "/retest"
)

func DefaultTests() []string {
	return []string{"install-vc7", "install-azure", "install-aws"}
}

func FindTestsFromCommentBody(body string) []string {
	var tests []string

	if len(strings.TrimSpace(body)) == 0 {
		return tests
	}

	for _, line := range strings.Split(strings.TrimSuffix(body, "\n"), "\n") {
		// TODO: ideally RETEST should only run previous set of tests commented
		if strings.HasPrefix(line, OK_TO_TEST) || strings.HasPrefix(line, RETEST) {
			tests = append(tests, DefaultTests()...)
			// do not proceed forward
			return tests
		}
		if strings.HasPrefix(line, TEST) {
			rest := strings.TrimLeft(line, TEST)
			rest = strings.TrimSpace(rest)
			if len(rest) == 0 {
				tests = append(tests, DefaultTests()...)
			} else {
				tests = append(tests, strings.Split(rest, ",")...)
			}
		}
	}
	// TODO: set (remove duplication) before returning
	return tests
}
