package record

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/results/pkg/api/server/db"
	"github.com/tektoncd/results/pkg/api/server/internal/protoutil"
	ppb "github.com/tektoncd/results/proto/pipeline/v1beta1/pipeline_go_proto"
	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestParseName(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		// if want is nil, assume error
		want []string
	}{
		{
			name: "simple",
			in:   "a/results/b/records/c",
			want: []string{"a", "b", "c"},
		},
		{
			name: "resource name reuse",
			in:   "results/results/records/records/records",
			want: []string{"results", "records", "records"},
		},
		{
			name: "missing name",
			in:   "a/results/b/records/",
		},
		{
			name: "missing name, no slash",
			in:   "a/results/b/records/",
		},
		{
			name: "missing parent",
			in:   "/records/b",
		},
		{
			name: "missing parent, no slash",
			in:   "records/b",
		},
		{
			name: "wrong resource",
			in:   "a/tacocat/b/records/c",
		},
		{
			name: "result resource",
			in:   "a/results/b",
		},
		{
			name: "invalid parent",
			in:   "a/b/results/c",
		},
		{
			name: "invalid name",
			in:   "a/results/b/records/c/d",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parent, result, name, err := ParseName(tc.in)
			if err != nil {
				if tc.want == nil {
					// error was expected, continue
					return
				}
				t.Fatal(err)
			}
			if tc.want == nil {
				t.Fatalf("expected error, got: [%s, %s]", parent, name)
			}

			if parent != tc.want[0] || result != tc.want[1] || name != tc.want[2] {
				t.Errorf("want: %v, got: [%s, %s, %s]", tc.want, parent, result, name)
			}
		})
	}
}

func TestToStorage(t *testing.T) {
	data := &ppb.TaskRun{Metadata: &ppb.ObjectMeta{Name: "tacocat"}}

	for _, tc := range []struct {
		name string
		in   *pb.Record
		want *db.Record
	}{
		{
			name: "full",
			in: &pb.Record{
				Name: "foo/results/bar",
				Id:   "a",
				Data: protoutil.Any(t, data),
				// These fields are ignored for now.
				Etag: "tacocat",
			},
			want: &db.Record{
				Parent:     "foo",
				ResultID:   "1",
				ResultName: "bar",
				Name:       "baz",
				ID:         "a",
				Data:       protoutil.AnyBytes(t, data),
			},
		},
		{
			name: "missing data",
			in: &pb.Record{
				Name: "foo/results/bar",
				Id:   "a",
			},
			want: &db.Record{
				Parent:     "foo",
				ResultID:   "1",
				ResultName: "bar",
				Name:       "baz",
				ID:         "a",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToStorage("foo", "bar", "1", "baz", tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("-want,+got: %s", diff)
			}
		})
	}

}

func TestToAPI(t *testing.T) {
	data := &ppb.TaskRun{Metadata: &ppb.ObjectMeta{Name: "tacocat"}}
	for _, tc := range []struct {
		name string
		in   *db.Record
		want *pb.Record
	}{
		{
			name: "full",
			in: &db.Record{
				Parent:     "foo",
				ResultID:   "1",
				ResultName: "bar",
				Name:       "baz",
				ID:         "a",
				Data:       protoutil.AnyBytes(t, data),
			},
			want: &pb.Record{
				Name: "foo/results/bar/records/baz",
				Id:   "a",
				Data: protoutil.Any(t, data),
			},
		},
		{
			in: &db.Record{
				Parent:     "foo",
				ResultID:   "1",
				ResultName: "bar",
				Name:       "baz",
				ID:         "a",
			},
			want: &pb.Record{
				Name: "foo/results/bar/records/baz",
				Id:   "a",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToAPI(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("-want,+got: %s", diff)
			}
		})
	}
}

func TestFormatName(t *testing.T) {
	got := FormatName("a", "b")
	want := "a/records/b"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}
