package archivist

import (
    "archive/zip"
    "io"
    "os"
    "path/filepath"
    "strings"
)

type zipper struct {
    src    string
    dst    string
    writer *zip.Writer
    buffer []byte
}

// zipper.do initialises output file, archive and directory tree walk function
func (z *zipper) do() error {
    out, err := os.Create(z.dst)
    if err != nil {
        return err
    }

    z.buffer = make([]byte, 1<<20)

    z.writer = zip.NewWriter(out)
    if err := filepath.Walk(z.src, z.walk); err != nil {
        return err
    }

    if err := z.writer.Close(); err != nil {
        return err
    }

    return out.Close()
}

// zipper.walk gets called for each file in given directory tree
func (z *zipper) walk(path string, info os.FileInfo, err error) error {
    if err != nil { return err }

    if !info.Mode().IsRegular() || info.Size() == 0 {
        return nil
    }

    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    fileName := strings.TrimPrefix(path, z.src + string(filepath.Separator))
    w, err := z.writer.Create(fileName)
    if err != nil {
        return err
    }

    _, err = io.CopyBuffer(w, file, z.buffer)

    return err
}
