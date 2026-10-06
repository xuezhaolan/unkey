package validation

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

var numericBoundReasonPattern = regexp.MustCompile(`^(maximum|minimum|exclusiveMaximum|exclusiveMinimum): got ([^,]+), want ([^,]+)$`)

func formatSchemaValidationReason(reason string) string {
	matches := numericBoundReasonPattern.FindStringSubmatch(reason)
	if matches == nil {
		return reason
	}

	return matches[1] + ": got " + formatDecimal(matches[2]) + ", want " + formatDecimal(matches[3])
}

func formatDecimal(value string) string {
	if strings.Contains(value, "×10") {
		return formatSuperscriptScientific(value)
	}
	if !strings.ContainsAny(value, "eE") {
		return value
	}

	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return value
	}

	return strconv.FormatFloat(number, 'f', -1, 64)
}

func formatSuperscriptScientific(value string) string {
	mantissaText, exponentText, ok := strings.Cut(value, "×10")
	if !ok {
		return value
	}

	mantissa, err := strconv.ParseFloat(mantissaText, 64)
	if err != nil {
		return value
	}

	exponent, ok := parseSuperscriptInteger(exponentText)
	if !ok {
		return value
	}

	return strconv.FormatFloat(mantissa*math.Pow10(exponent), 'f', -1, 64)
}

func parseSuperscriptInteger(value string) (int, bool) {
	replacer := strings.NewReplacer(
		"⁰", "0",
		"¹", "1",
		"²", "2",
		"³", "3",
		"⁴", "4",
		"⁵", "5",
		"⁶", "6",
		"⁷", "7",
		"⁸", "8",
		"⁹", "9",
		"⁺", "+",
		"⁻", "-",
	)
	exponent, err := strconv.Atoi(replacer.Replace(value))
	if err != nil {
		return 0, false
	}
	return exponent, true
}
