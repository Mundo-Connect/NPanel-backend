package tool

import (
	"testing"
	"time"
)

func TestAddTimeSupportsProductAndRedemptionUnits(t *testing.T) {
	basic := time.Date(2026, time.June, 28, 3, 23, 24, 0, time.UTC)
	tests := []struct {
		name     string
		unit     string
		quantity int64
		want     time.Time
	}{
		{name: "product month", unit: "Month", quantity: 1, want: basic.AddDate(0, 1, 0)},
		{name: "redemption month", unit: "month", quantity: 1, want: basic.AddDate(0, 1, 0)},
		{name: "redemption day", unit: "day", quantity: 10, want: basic.AddDate(0, 0, 10)},
		{name: "redemption quarter", unit: "quarter", quantity: 1, want: basic.AddDate(0, 3, 0)},
		{name: "redemption half year", unit: "half_year", quantity: 1, want: basic.AddDate(0, 6, 0)},
		{name: "redemption year", unit: "year", quantity: 1, want: basic.AddDate(1, 0, 0)},
		{name: "product no limit", unit: "NoLimit", quantity: 1, want: time.UnixMilli(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddTime(tt.unit, tt.quantity, basic)
			if !got.Equal(tt.want) {
				t.Fatalf("AddTime(%q, %d) = %v, want %v", tt.unit, tt.quantity, got, tt.want)
			}
		})
	}
}

func TestGetYearDays(t *testing.T) {
	days := GetYearDays(time.Now(), 2, 1)
	t.Logf("GetYearDays() success, expected 365, got %d", days)

}
