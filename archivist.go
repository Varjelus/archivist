package archivist

import (
    "archive/zip"
    "errors"
    "io"
    "os"
    "path/filepath"
    "strings"
)

type zipper struct {
    src    string
    dst    string
    writer *zip.Writer
}

func (z *zipper) do() error {
    // Output file
    out, err := os.Create(z.dst)//os.OpenFile(z.dst, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
    if err != nil {
        return err
    }

    z.writer = zip.NewWriter(out)
    if err := filepath.Walk(z.src, z.walk); err != nil {
        return err
    }

    if err := z.writer.Close(); err != nil {
        return err
    }

    return out.Close()
}

func (z *zipper) walk(path string, info os.FileInfo, err error) error {
    if err != nil { return err }

    if info.IsDir() || info.Size() == 0 {
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

    _, err = io.Copy(w, file)
    return err
}

func Store(src, dst string) error {
    z := &zipper{
        src: filepath.Clean(filepath.FromSlash(src)),
        dst: filepath.Clean(filepath.FromSlash(dst)),
    }
    return z.do()
}
