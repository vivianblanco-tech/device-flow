package models

import (
	"fmt"
	"regexp"
	"strings"
)

// GenerateSKU generates a SKU code based on laptop specifications
// Format for non-MacBook: C.NOT.{CPU}.{RAM}.{SSD}
// Format for MacBook: C.MAC.{CHIP}.{RAM}.{SSD}
func GenerateSKU(model, cpu, ram, ssd string) string {
	// Return empty if any field is empty
	if model == "" || cpu == "" || ram == "" || ssd == "" {
		return ""
	}

	// Normalize inputs to lowercase for comparison
	modelLower := strings.ToLower(model)
	cpuLower := strings.ToLower(cpu)

	// Determine if it's a MacBook
	isMacBook := strings.Contains(modelLower, "macbook")

	// Parse CPU code
	cpuCode := parseCPUCode(cpuLower, isMacBook)

	// Parse RAM code
	ramCode := parseRAMCode(ram)

	// Parse SSD code
	ssdCode := parseSSDCode(ssd)

	// Build SKU
	if isMacBook {
		return fmt.Sprintf("C.MAC.%s.%s.%s", cpuCode, ramCode, ssdCode)
	}
	return fmt.Sprintf("C.NOT.%s.%s.%s", cpuCode, ramCode, ssdCode)
}

// parseCPUCode converts CPU string to SKU CPU code
func parseCPUCode(cpu string, isMacBook bool) string {
	cpu = strings.ToLower(strings.TrimSpace(cpu))

	if isMacBook {
		// MacBook chip codes
		if strings.Contains(cpu, "m1 ultra") {
			return "MU1"
		}
		if strings.Contains(cpu, "m1 max") {
			return "MM1"
		}
		if strings.Contains(cpu, "m1 pro") {
			return "MP1"
		}
		if strings.Contains(cpu, "m1") {
			return "M01"
		}
		if strings.Contains(cpu, "m2 ultra") {
			return "MU2"
		}
		if strings.Contains(cpu, "m2 max") {
			return "MM2"
		}
		if strings.Contains(cpu, "m2 pro") {
			return "MP2"
		}
		if strings.Contains(cpu, "m2") {
			return "M02"
		}
		if strings.Contains(cpu, "m3 ultra") {
			return "MU3"
		}
		if strings.Contains(cpu, "m3 max") {
			return "MM3"
		}
		if strings.Contains(cpu, "m3 pro") {
			return "MP3"
		}
		if strings.Contains(cpu, "m3") {
			return "M03"
		}
		if strings.Contains(cpu, "m4 ultra") {
			return "MU4"
		}
		if strings.Contains(cpu, "m4 max") {
			return "MM4"
		}
		if strings.Contains(cpu, "m4 pro") {
			return "MP4"
		}
		if strings.Contains(cpu, "m4") {
			return "M04"
		}
		// Unknown MacBook chip
		return "UNK"
	}

	// Non-MacBook (Intel/AMD) CPU codes
	if strings.Contains(cpu, "i9") {
		return "0I9"
	}
	if strings.Contains(cpu, "i7") {
		return "0I7"
	}
	if strings.Contains(cpu, "i5") {
		return "0I5"
	}
	if strings.Contains(cpu, "i3") {
		return "0I3"
	}
	if strings.Contains(cpu, "ryzen 9") || strings.Contains(cpu, "r9") {
		return "0R9"
	}
	if strings.Contains(cpu, "ryzen 7") || strings.Contains(cpu, "r7") {
		return "0R7"
	}
	if strings.Contains(cpu, "ryzen 5") || strings.Contains(cpu, "r5") {
		return "0R5"
	}
	if strings.Contains(cpu, "ryzen 3") || strings.Contains(cpu, "r3") {
		return "0R3"
	}

	// Unknown CPU
	return "UNK"
}

// parseRAMCode converts RAM string to SKU RAM code (e.g., "16GB" -> "016")
func parseRAMCode(ram string) string {
	// Extract numeric value from RAM string
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(ram)
	if len(matches) > 1 {
		ramValue := matches[1]
		// Pad to 3 digits
		return fmt.Sprintf("%03s", ramValue)
	}
	return "000"
}

// parseSSDCode converts SSD string to SKU SSD code (e.g., "512GB" -> "2G", "1TB" -> "1T")
func parseSSDCode(ssd string) string {
	ssd = strings.ToUpper(strings.TrimSpace(ssd))

	// Check for TB first
	if strings.Contains(ssd, "TB") {
		re := regexp.MustCompile(`(\d+)`)
		matches := re.FindStringSubmatch(ssd)
		if len(matches) > 1 {
			return matches[1] + "T"
		}
	}

	// Check for GB
	if strings.Contains(ssd, "GB") {
		re := regexp.MustCompile(`(\d+)`)
		matches := re.FindStringSubmatch(ssd)
		if len(matches) > 1 {
			ssdValue := matches[1]
			// Convert GB to format: 512GB -> 2G, 256GB -> G, 1024GB -> 4G
			switch ssdValue {
			case "128":
				return "G" // 128GB -> G (half of 256)
			case "256":
				return "1G"
			case "512":
				return "2G"
			case "1024":
				return "4G"
			case "2048":
				return "8G"
			default:
				// For other values, calculate: value / 256
				// This is a simple heuristic
				return "XG"
			}
		}
	}

	return "0G"
}

