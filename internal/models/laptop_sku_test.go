package models

import (
	"testing"
)

func TestGenerateSKU(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		cpu      string
		ram      string
		ssd      string
		expected string
	}{
		// Non-MacBook (Notebook) tests
		{
			name:     "Notebook i5 16GB 512GB",
			model:    "Dell Latitude 5520",
			cpu:      "i5",
			ram:      "16GB",
			ssd:      "512GB",
			expected: "C.NOT.0I5.016.2G",
		},
		{
			name:     "Notebook i7 16GB 512GB",
			model:    "HP EliteBook",
			cpu:      "i7",
			ram:      "16GB",
			ssd:      "512GB",
			expected: "C.NOT.0I7.016.2G",
		},
		{
			name:     "Notebook i7 32GB 1TB",
			model:    "Lenovo ThinkPad",
			cpu:      "i7",
			ram:      "32GB",
			ssd:      "1TB",
			expected: "C.NOT.0I7.032.1T",
		},
		{
			name:     "Notebook i9 32GB 512GB",
			model:    "Dell XPS 15",
			cpu:      "i9",
			ram:      "32GB",
			ssd:      "512GB",
			expected: "C.NOT.0I9.032.2G",
		},
		{
			name:     "Notebook i9 32GB 1TB",
			model:    "HP ZBook",
			cpu:      "i9",
			ram:      "32GB",
			ssd:      "1TB",
			expected: "C.NOT.0I9.032.1T",
		},
		// MacBook M1 tests
		{
			name:     "MacBook M1 16GB 512GB",
			model:    "MacBook Pro M1",
			cpu:      "M1",
			ram:      "16GB",
			ssd:      "512GB",
			expected: "C.MAC.M01.016.2G",
		},
		{
			name:     "MacBook M1 32GB 1TB",
			model:    "MacBook Air M1",
			cpu:      "M1",
			ram:      "32GB",
			ssd:      "1TB",
			expected: "C.MAC.M01.032.1T",
		},
		{
			name:     "MacBook M1 64GB 1TB",
			model:    "MacBook M1",
			cpu:      "M1",
			ram:      "64GB",
			ssd:      "1TB",
			expected: "C.MAC.M01.064.1T",
		},
		// MacBook M1 Pro tests
		{
			name:     "MacBook M1 Pro 16GB 512GB",
			model:    "MacBook Pro M1 Pro",
			cpu:      "M1 Pro",
			ram:      "16GB",
			ssd:      "512GB",
			expected: "C.MAC.MP1.016.2G",
		},
		{
			name:     "MacBook M1 Pro 32GB 512GB",
			model:    "MacBook M1 Pro",
			cpu:      "M1 Pro",
			ram:      "32GB",
			ssd:      "512GB",
			expected: "C.MAC.MP1.032.2G",
		},
		{
			name:     "MacBook M1 Pro 64GB 1TB",
			model:    "MacBook Pro M1 Pro",
			cpu:      "M1 Pro",
			ram:      "64GB",
			ssd:      "1TB",
			expected: "C.MAC.MP1.064.1T",
		},
		// MacBook M1 Max tests
		{
			name:     "MacBook M1 Max 16GB 512GB",
			model:    "MacBook Pro M1 Max",
			cpu:      "M1 Max",
			ram:      "16GB",
			ssd:      "512GB",
			expected: "C.MAC.MM1.016.2G",
		},
		{
			name:     "MacBook M1 Max 32GB 512GB",
			model:    "MacBook M1 Max",
			cpu:      "M1 Max",
			ram:      "32GB",
			ssd:      "512GB",
			expected: "C.MAC.MM1.032.2G",
		},
		{
			name:     "MacBook M1 Max 32GB 1TB",
			model:    "MacBook M1 Max",
			cpu:      "M1 Max",
			ram:      "32GB",
			ssd:      "1TB",
			expected: "C.MAC.MM1.032.1T",
		},
		{
			name:     "MacBook M1 Max 64GB 1TB",
			model:    "MacBook M1 Max",
			cpu:      "M1 Max",
			ram:      "64GB",
			ssd:      "1TB",
			expected: "C.MAC.MM1.064.1T",
		},
		// MacBook M1 Ultra tests (using MU1 as chip code)
		{
			name:     "MacBook M1 Ultra 16GB 512GB",
			model:    "MacBook Pro M1 Ultra",
			cpu:      "M1 Ultra",
			ram:      "16GB",
			ssd:      "512GB",
			expected: "C.MAC.MU1.016.2G",
		},
		{
			name:     "MacBook M1 Ultra 32GB 512GB",
			model:    "MacBook M1 Ultra",
			cpu:      "M1 Ultra",
			ram:      "32GB",
			ssd:      "512GB",
			expected: "C.MAC.MU1.032.2G",
		},
		{
			name:     "MacBook M1 Ultra 64GB 1TB",
			model:    "MacBook M1 Ultra",
			cpu:      "M1 Ultra",
			ram:      "64GB",
			ssd:      "1TB",
			expected: "C.MAC.MU1.064.1T",
		},
		// MacBook M2 tests
		{
			name:     "MacBook M2 16GB 512GB",
			model:    "MacBook Pro M2",
			cpu:      "M2",
			ram:      "16GB",
			ssd:      "512GB",
			expected: "C.MAC.M02.016.2G",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSKU(tt.model, tt.cpu, tt.ram, tt.ssd)
			if result != tt.expected {
				t.Errorf("GenerateSKU() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateSKU_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		cpu      string
		ram      string
		ssd      string
		expected string
	}{
		{
			name:     "Empty fields returns empty",
			model:    "",
			cpu:      "",
			ram:      "",
			ssd:      "",
			expected: "",
		},
		{
			name:     "Unknown CPU format",
			model:    "Dell Latitude",
			cpu:      "Unknown CPU",
			ram:      "16GB",
			ssd:      "512GB",
			expected: "C.NOT.UNK.016.2G",
		},
		{
			name:     "Case insensitive MacBook detection",
			model:    "macbook pro m1",
			cpu:      "m1",
			ram:      "16GB",
			ssd:      "512GB",
			expected: "C.MAC.M01.016.2G",
		},
		{
			name:     "Case insensitive CPU codes",
			model:    "Dell XPS",
			cpu:      "I7",
			ram:      "32GB",
			ssd:      "1TB",
			expected: "C.NOT.0I7.032.1T",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSKU(tt.model, tt.cpu, tt.ram, tt.ssd)
			if result != tt.expected {
				t.Errorf("GenerateSKU() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLaptop_GenerateAndSetSKU(t *testing.T) {
	tests := []struct {
		name        string
		laptop      *Laptop
		expectedSKU string
	}{
		{
			name: "Generate SKU for Dell laptop with i7",
			laptop: &Laptop{
				Model: "Dell Latitude 5520",
				CPU:   "i7",
				RAMGB: "16GB",
				SSDGB: "512GB",
			},
			expectedSKU: "C.NOT.0I7.016.2G",
		},
		{
			name: "Generate SKU for MacBook Pro M1",
			laptop: &Laptop{
				Model: "MacBook Pro",
				CPU:   "M1",
				RAMGB: "32GB",
				SSDGB: "1TB",
			},
			expectedSKU: "C.MAC.M01.032.1T",
		},
		{
			name: "Generate SKU for MacBook M1 Max",
			laptop: &Laptop{
				Model: "MacBook Pro M1 Max",
				CPU:   "M1 Max",
				RAMGB: "64GB",
				SSDGB: "1TB",
			},
			expectedSKU: "C.MAC.MM1.064.1T",
		},
		{
			name: "Do not overwrite existing SKU",
			laptop: &Laptop{
				SKU:   "CUSTOM-SKU-123",
				Model: "Dell XPS",
				CPU:   "i9",
				RAMGB: "32GB",
				SSDGB: "1TB",
			},
			expectedSKU: "CUSTOM-SKU-123",
		},
		{
			name: "Empty SKU should be generated",
			laptop: &Laptop{
				SKU:   "",
				Model: "HP EliteBook",
				CPU:   "i5",
				RAMGB: "16GB",
				SSDGB: "512GB",
			},
			expectedSKU: "C.NOT.0I5.016.2G",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.laptop.GenerateAndSetSKU()
			if tt.laptop.SKU != tt.expectedSKU {
				t.Errorf("Laptop.GenerateAndSetSKU() SKU = %v, want %v", tt.laptop.SKU, tt.expectedSKU)
			}
		})
	}
}
