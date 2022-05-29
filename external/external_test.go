package external

import "testing"

func Test_compatibleLibraryVersion(t *testing.T) {
	tests := []struct {
		name      string
		versions  map[string]string
		cgVersion string
		want      string
	}{
		{name: "exact match", versions: map[string]string{
			"0.1": "1.2",
			"0.3": "1.4",
			"0.5": "1.6",
			"0.7": "1.8",
			"1.1": "2.2",
			"1.5": "2.6",
			"1.7": "2.8",
		}, cgVersion: "0.5", want: "1.6"},
		{name: "next highest minor", versions: map[string]string{
			"0.1": "1.2",
			"0.3": "1.4",
			"0.6": "1.6",
			"0.7": "1.8",
			"1.1": "2.2",
			"1.5": "2.6",
			"1.7": "2.8",
		}, cgVersion: "0.5", want: "1.6"},
		{name: "next lowest minor", versions: map[string]string{
			"0.1": "1.2",
			"0.3": "1.4",
			"1.1": "2.2",
			"1.5": "2.6",
			"1.7": "2.8",
		}, cgVersion: "0.5", want: "1.4"},
		{name: "none", versions: map[string]string{
			"0.1": "1.2",
			"0.3": "1.4",
			"1.1": "2.2",
			"1.5": "2.6",
			"1.7": "2.8",
		}, cgVersion: "2.5", want: "latest"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compatibleLibraryVersion(tt.versions, tt.cgVersion); got != tt.want {
				t.Errorf("compatibleLibraryVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
