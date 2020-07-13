package trello

import "testing"

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		token   string
		wantErr bool
	}{
		{name: "both empty", key: "", token: "", wantErr: true},
		{name: "empty key", key: "val", token: "", wantErr: true},
		{name: "empty token", key: "", token: "val", wantErr: true},
		{name: "valid", key: "val", token: "val", wantErr: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api, err := New(test.key, test.token)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if api != nil {
					t.Fatalf("unexpected api: %+v", api)
				}
			} else if !test.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if api == nil {
					t.Fatal("unexpected nil api")
				}
			}
		})
	}
}
