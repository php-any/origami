package spl

import "testing"

func TestPhpPathJoinPreservesDotDot(t *testing.T) {
	dir := "/vendor/laravel/telescope/src/../database/migrations"
	got := phpPathJoin(dir, "2018_08_08_100000_create_telescope_entries_table.php")
	want := "/vendor/laravel/telescope/src/../database/migrations/2018_08_08_100000_create_telescope_entries_table.php"
	if got != want {
		t.Fatalf("phpPathJoin = %q, want %q", got, want)
	}
}

func TestPhpPathDirPreservesDotDot(t *testing.T) {
	path := "/vendor/laravel/telescope/src/../database/migrations/file.php"
	got := phpPathDir(path)
	want := "/vendor/laravel/telescope/src/../database/migrations"
	if got != want {
		t.Fatalf("phpPathDir = %q, want %q", got, want)
	}
}

func TestPhpPathRelRootFile(t *testing.T) {
	root := "/vendor/pkg/src/../database/migrations"
	full := root + "/2018_08_08_100000_create_telescope_entries_table.php"
	got := phpPathRel(root, full)
	want := "2018_08_08_100000_create_telescope_entries_table.php"
	if got != want {
		t.Fatalf("phpPathRel = %q, want %q", got, want)
	}
	if sub := phpPathDir(got); sub != "." {
		t.Fatalf("phpPathDir(rel) = %q, want \".\"", sub)
	}
}

func TestPhpPathRelNested(t *testing.T) {
	root := "/a/src/../migrations"
	full := root + "/sub/file.php"
	got := phpPathRel(root, full)
	if got != "sub/file.php" {
		t.Fatalf("phpPathRel = %q", got)
	}
	if phpPathDir(got) != "sub" {
		t.Fatalf("subPath = %q", phpPathDir(got))
	}
}
