package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "Birthday was 10 days ago (already occurred this year)",
			dob:      time.Now().AddDate(-25, 0, -10),
			expected: 25,
		},
		{
			name:     "Birthday is in 10 days (not occurred yet this year)",
			dob:      time.Now().AddDate(-25, 0, 10),
			expected: 24,
		},
		{
			name:     "Birthday is today",
			dob:      time.Now().AddDate(-25, 0, 0),
			expected: 25,
		},
		{
			name:     "Born exactly 1 year ago today",
			dob:      time.Now().AddDate(-1, 0, 0),
			expected: 1,
		},
		{
			name:     "Born 1 year ago, but birthday is tomorrow",
			dob:      time.Now().AddDate(-1, 0, 1),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CalculateAge(tt.dob)
			if actual != tt.expected {
				t.Errorf("CalculateAge(%v) = %d; expected %d", tt.dob.Format("2006-01-02"), actual, tt.expected)
			}
		})
	}
}
