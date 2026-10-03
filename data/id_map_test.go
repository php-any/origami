package data

import "testing"

func TestIDMapSparseBranchSnapshots(t *testing.T) {
	var original IDMap[int]
	ids := []uint32{0, 63, 64, 4095, 4096, 8192, 60000}
	versions := []IDMap[int]{original}
	for i, id := range ids {
		versions = append(versions, versions[len(versions)-1].With(id, i+1))
	}
	latest := versions[len(versions)-1].With(4096, 99)
	for version, snapshot := range versions {
		for i, id := range ids {
			value, found := snapshot.Get(id)
			if found != (i < version) || found && value != i+1 {
				t.Fatalf("snapshot %d, id %d: %d/%v", version, id, value, found)
			}
		}
		if _, found := snapshot.Get(50000); found {
			t.Fatal("hole in sparse branch appeared present")
		}
	}
	if value, found := latest.Get(4096); !found || value != 99 {
		t.Fatal("branch overwrite lost new value")
	}
}

func BenchmarkIDMapSparsePublication(b *testing.B) {
	var snapshot IDMap[int]
	for id := uint32(0); id < 60000; id += 97 {
		snapshot = snapshot.With(id, 1)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = snapshot.With(59999, i)
	}
}
