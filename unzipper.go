package archivist

import (
    "archive/zip"
    "io"
    "os"
    "path/filepath"
)

type unzipper struct {
    src    string
    dst    string
    reader *zip.ReadCloser
    buffer []byte
}

// unzipper.do initialises output file and unzips source there
func (z *unzipper) do() error {
    err := os.MkdirAll(z.dst, perm)
    if err != nil { return err }

    z.buffer = make([]byte, bufSize)

    z.reader, err = zip.OpenReader(z.src)
    if err != nil { return err }

    for _, f := range z.reader.File {
        if err := z.unzip(f); err != nil {
            z.reader.Close()
            return err
        }
    }

    if err := z.reader.Close(); err != nil {
        return err
    }

    return nil
}

func (z *unzipper) unzip(f *zip.File) error {
    if f.FileInfo().IsDir() { return nil }

    fName := filepath.Join(z.dst, f.Name)
    dir, _ := filepath.Split(fName)

    if err := os.MkdirAll(dir, perm); err != nil && os.IsNotExist(err) {
        return err
    }

    r, err := f.Open()
    if err != nil {
        return err
    }
    defer r.Close()

    w, err := os.Create(filepath.Join(z.dst, f.Name))
    if err != nil {
        return err
    }
    defer w.Close()

    if _, err := io.CopyBuffer(w, r, z.buffer); err != nil {
        w.Close()
        return err
    }

    if err := r.Close(); err != nil {
        return err
    }

    return w.Close()
}
