package archivist

import (
    "os"
    "path/filepath"
)

var (
    perm os.FileMode = os.ModeDir|0700
    bufSize int = 1<<20 // ~1MB default
)

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

func Tar(src, dst string) error {
    t := &tarmonster{
        src: filepath.Clean(filepath.FromSlash(src)),
        dst: filepath.Clean(filepath.FromSlash(dst)),
    }
    return t.do()
}

func Untar(src, dst string) error {
    t := &untarmonster{
        src: filepath.Clean(filepath.FromSlash(src)),
        dst: filepath.Clean(filepath.FromSlash(dst)),
    }
    return t.do()
}

func SetFileMode(mode os.FileMode) {
    perm = mode
}

func SetBufferSize(size int) {
    bufSize = size
}