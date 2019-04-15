package amqp

import (
	"encoding/binary"
	"testing"
)

func TestLinkOptions(t *testing.T) {
	tests := []struct {
		label string
		opts  []LinkOption

		wantSource     *source
		wantProperties map[symbol]interface{}
	}{
		{
			label: "no options",
		},
		{
			label: "link-filters",
			opts: []LinkOption{
				LinkSelectorFilter("amqp.annotation.x-opt-offset > '100'"),
				LinkProperty("x-opt-test1", "test1"),
				LinkProperty("x-opt-test2", "test2"),
				LinkProperty("x-opt-test1", "test3"),
				LinkPropertyInt64("x-opt-test4", 1),
				LinkSourceFilter("com.microsoft:session-filter", 0x00000137000000C, "123"),
			},

			wantSource: &source{
				Filter: map[symbol]*describedType{
					"apache.org:selector-filter:string": {
						descriptor: binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x46, 0x8C, 0x00, 0x00, 0x00, 0x04}),
						value:      "amqp.annotation.x-opt-offset > '100'",
					},
					"com.microsoft:session-filter": {
						descriptor: binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x00, 0x13, 0x70, 0x00, 0x00, 0x0C}),
						value:      "123",
					},
				},
			},
			wantProperties: map[symbol]interface{}{
				"x-opt-test1": "test3",
				"x-opt-test2": "test2",
				"x-opt-test4": int64(1),
			},
		},
		{
			label: "more-link-filters",
			opts: []LinkOption{
				LinkSourceFilter("com.microsoft:session-filter", 0x00000137000000C, nil),
			},

			wantSource: &source{
				Filter: map[symbol]*describedType{
					"com.microsoft:session-filter": {
						descriptor: binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x00, 0x13, 0x70, 0x00, 0x00, 0x0C}),
						value:      nil,
					},
				},
			},
		},
		{
			label: "link-source-capabilities",
			opts: []LinkOption{
				LinkSourceCapabilities("cap1", "cap2", "cap3"),
			},
			wantSource: &source{
				Capabilities: []symbol{"cap1", "cap2", "cap3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			got, err := newLink(nil, nil, tt.opts)
			if err != nil {
				t.Fatal(err)
			}

			if !testEqual(got.source, tt.wantSource) {
				t.Errorf("Source properties don't match expected:\n %s", testDiff(got.source, tt.wantSource))
			}

			if !testEqual(got.properties, tt.wantProperties) {
				t.Errorf("Link properties don't match expected:\n %s", testDiff(got.properties, tt.wantProperties))
			}
		})
	}
}

func TestSourceName(t *testing.T) {
	expectedSourceName := "source-name"
	opts := []LinkOption{
		LinkName(expectedSourceName),
	}

	got, err := newLink(nil, nil, opts)
	if err != nil {
		t.Fatal(err)
	}

	if got.name != expectedSourceName {
		t.Errorf("Link Source Name does not match expected: %v got: %v", expectedSourceName, got.name)
	}
}
