package archivist

import (
    "path/filepath"
    "fmt"
    "archive/tar"
    "io"
    "os"
)

type untarmonster struct {
    src    string
    dst    string
    reader *tar.Reader
    buffer []byte
}

func (t *untarmonster) do() error {
    err := os.MkdirAll(t.dst, perm)
    if err != nil { return err }

    t.buffer = make([]byte, bufSize)

    tarball, err := os.Open(t.src)
    if err != nil {
        return err
    }
    defer tarball.Close()

    t.reader = tar.NewReader(tarball)
    for {
        h, err := t.reader.Next()
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }

        if err := t.untar(h); err != nil {
            return err
        }
    }

    return tarball.Close()
}

func (t *untarmonster) untar(h *tar.Header) error {
    path := filepath.Join(t.dst, h.Name)

    fi := h.FileInfo()

    if fi.IsDir() {
        return os.MkdirAll(path, perm)
    }

    if !fi.Mode().IsRegular() {
        return nil
    }

    dir, _ := filepath.Split(path)
    err := os.MkdirAll(dir, perm)
    if err != nil {
        return fmt.Errorf("creating %s", dir)
    }

    file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, perm)
    if err != nil {
        return err
    }
    defer file.Close()

    if _, err := io.CopyBuffer(file, t.reader, t.buffer); err != nil {
        return err
    }

    if err := file.Sync(); err != nil {
        return err
    }

    return file.Close()
}