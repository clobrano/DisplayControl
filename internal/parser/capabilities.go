package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/clobrano/DisplayControl/internal/config"
)

// ParseCapabilities parses ddcutil capabilities output and returns VCP features
func ParseCapabilities(output string) (map[string]config.VCPFeature, error) {
	features := make(map[string]config.VCPFeature)

	// Regex to match feature lines: "Feature: XX (Description)"
	featureRegex := regexp.MustCompile(`Feature:\s+([A-Fa-f0-9]+)\s+\(([^)]+)\)`)
	// Regex to match discrete value lines: "XX: Description"
	valueRegex := regexp.MustCompile(`^\s+([A-Fa-f0-9]+):\s+(.+)$`)
	// Regex to match range values: "00..64"
	rangeRegex := regexp.MustCompile(`Values:\s+([A-Fa-f0-9]+)\.\.([A-Fa-f0-9]+)`)

	lines := strings.Split(output, "\n")

	var currentFeature *config.VCPFeature
	var currentFeatureName string
	var inValues bool

	for _, line := range lines {
		// Check for feature definition
		if matches := featureRegex.FindStringSubmatch(line); matches != nil {
			// Save previous feature if it exists and has values
			if currentFeature != nil && len(currentFeature.Values) > 0 {
				features[currentFeatureName] = *currentFeature
			}

			code := matches[1]
			description := matches[2]

			currentFeature = &config.VCPFeature{
				Code:        decimalToHex(code),
				Description: description,
				Values:      make(map[string]string),
			}
			currentFeatureName = sanitizeFeatureName(description)
			inValues = false
			continue
		}

		// Check for range values (continuous features)
		if currentFeature != nil && rangeRegex.MatchString(line) {
			matches := rangeRegex.FindStringSubmatch(line)
			if len(matches) == 3 {
				startHex := matches[1]
				endHex := matches[2]

				// Convert hex to decimal
				start, _ := strconv.ParseInt(startHex, 16, 64)
				end, _ := strconv.ParseInt(endHex, 16, 64)

				// Generate values at intervals (0, 25, 50, 75, 100 for 0-100 range)
				intervals := generateRangeValues(int(start), int(end))
				for name, value := range intervals {
					currentFeature.Values[name] = value
				}
			}
			continue
		}

		// Check for "Values:" line
		if currentFeature != nil && strings.Contains(line, "Values:") {
			inValues = true
			continue
		}

		// Check for discrete value entries
		if currentFeature != nil && inValues {
			if matches := valueRegex.FindStringSubmatch(line); matches != nil {
				hexCode := strings.ToLower(matches[1])
				valueName := strings.TrimSpace(matches[2])
				currentFeature.Values[valueName] = "0x" + hexCode
			}
		}
	}

	// Don't forget the last feature
	if currentFeature != nil && len(currentFeature.Values) > 0 {
		features[currentFeatureName] = *currentFeature
	}

	return features, nil
}

// decimalToHex converts a hex or decimal string to lowercase hex with 0x prefix
// ddcutil output provides feature codes as hex strings (e.g., "10", "DC", "E0")
func decimalToHex(input string) string {
	input = strings.TrimSpace(input)

	// Parse as hex (ddcutil capabilities uses hex notation)
	if val, err := strconv.ParseInt(input, 16, 64); err == nil {
		return fmt.Sprintf("0x%02x", val)
	}

	// If all else fails, return with 0x prefix
	return "0x" + strings.ToLower(input)
}

// sanitizeFeatureName converts a feature description to a valid map key
func sanitizeFeatureName(description string) string {
	// Convert to lowercase
	name := strings.ToLower(description)

	// Remove common prefixes
	name = strings.TrimPrefix(name, "select ")

	// Replace spaces and special characters with underscores
	name = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(name, "_")

	// Remove leading/trailing underscores
	name = strings.Trim(name, "_")

	return name
}

// generateRangeValues generates a map of value names to hex codes for continuous ranges
func generateRangeValues(start, end int) map[string]string {
	values := make(map[string]string)

	// For typical ranges like 0-100, generate at 0, 25, 50, 75, 100
	if end == 100 {
		percentages := []int{0, 25, 50, 75, 100}
		for _, pct := range percentages {
			values[fmt.Sprintf("%d", pct)] = fmt.Sprintf("0x%02x", pct)
		}
		return values
	}

	// For other ranges, generate at start, 25%, 50%, 75%, end
	range_ := end - start
	steps := []float64{0, 0.25, 0.50, 0.75, 1.0}
	for _, step := range steps {
		val := start + int(float64(range_)*step)
		values[fmt.Sprintf("%d", val)] = fmt.Sprintf("0x%02x", val)
	}

	return values
}
