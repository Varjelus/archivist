package archivist

import (
    "os"
    "path/filepath"
)

var perm os.FileMode = os.ModeDir|0700

// Zip is an exported method which sanitizes io paths and starts archiving
func Zip(src, dst string) error {
    z := &zipper{
        src: filepath.Clean(filepath.FromSlash(src)),
        dst: filepath.Clean(filepath.FromSlash(dst)),
    }
    return z.do()
}

// Unzip is an exported method which sanitizes io paths and starts unzipping
func Unzip(src, dst string) error {
    z := &unzipper{
        src: filepath.Clean(filepath.FromSlash(src)),
        dst: filepath.Clean(filepath.FromSlash(dst)),
    }
    return z.do()
}

func SetFileMode(mode int) {
    perm = os.FileMode(mode)
}
