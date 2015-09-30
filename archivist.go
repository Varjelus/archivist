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
}

// zipper.do initialises output file, archive and directory tree walk function
func (z *zipper) do() error {
    out, err := os.Create(z.dst)
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

    // Not efficient FIXME
    _, err = io.Copy(w, file)

    return err
}

type unzipper struct {
    src    string
    dst    string
    reader *zip.ReadCloser
}

// unzipper.do initialises output file and unzips source there
func (z *unzipper) do() error {
    err := os.MkdirAll(z.dst, os.ModeDir)
    if err != nil { return err }

    z.reader, err = zip.OpenReader(z.src)
    if err != nil { return err }

    for _, f := range z.reader.File {
        r, err := f.Open()
        if err != nil { return err }

        w, err := os.Create(filepath.Join(z.dst, f.Name))
        if err != nil {
            z.reader.Close()
            r.Close()
            return err
        }

        if _, err := io.Copy(w, r); err != nil {
            z.reader.Close()
            w.Close()
            r.Close()
            return err
        }

        if err := r.Close(); err != nil {
            w.Close()
            z.reader.Close()
            return err
        }

        if err := w.Close(); err != nil {
            return err
        }
    }

    if err := z.reader.Close(); err != nil {
        return err
    }

    return nil
}

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
