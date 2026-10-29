package pdf

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadInlineImage(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
		want dict
	}{
		{
			name: `normal case`,
			data: []byte("BI /W 1 /H 1 ID a EI"),
			want: dict{name("W"): int64(1), name("H"): int64(1), name("ID"): "a"},
		},
		{
			name: `malformed data: BI only`,
			data: []byte("BI"),
			want: dict{name("ID"): ""},
		},
		{
			name: `malformed data: no ID`,
			data: []byte("BI /W 1 /H 1"),
			want: dict{name("W"): int64(1), name("H"): int64(1), name("ID"): ""},
		},
		{
			name: `malformed data: end with ID`,
			data: []byte("BI /W 1 /H 1 ID"),
			want: dict{name("W"): int64(1), name("H"): int64(1), name("ID"): ""},
		},
		{
			name: `malformed data: no image data`,
			data: []byte("BI /W 1 /H 1 ID "),
			want: dict{name("W"): int64(1), name("H"): int64(1), name("ID"): ""},
		},
		{
			name: `malformed data: no EI`,
			data: []byte("BI /W 2 /H 1 ID a"),
			want: dict{name("W"): int64(2), name("H"): int64(1), name("ID"): "a"},
		},
		{
			name: `malformed data: no space before EI`,
			data: []byte("BI /W 1 /H 1 ID aEI Q"),
			want: dict{name("W"): int64(1), name("H"): int64(1), name("ID"): "aEI Q"},
		},
		{
			name: `malformed data: the space between ID and EI is \r\n`,
			data: []byte("BI /W 1 /H 1 ID\r\na\r\nEI"),
			want: dict{name("W"): int64(1), name("H"): int64(1), name("ID"): "\na\r"},
		},
		{
			name: `malformed data: EI in image data`,
			data: []byte("BI /W 3 /H 1 ID  EI EI"),
			want: dict{name("W"): int64(3), name("H"): int64(1), name("ID"): ""},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.data)
			b := newBuffer(r, 0)
			b.allowEOF = true
			b.allowObjptr = false
			b.allowStream = false

			if tok := b.readToken(); tok != keyword("BI") {
				t.Fatalf("expected BI, got %v", tok)
			}
			got := b.readInlineImage()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}
