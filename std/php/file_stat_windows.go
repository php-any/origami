//go:build windows

package php

import "os"

func fileInodeNumber(os.FileInfo) (int, bool) { return 0, false }
func fileOwnerID(os.FileInfo) (int, bool)     { return 0, false }
func fileGroupID(os.FileInfo) (int, bool)     { return 0, false }
