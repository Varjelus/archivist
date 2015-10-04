package archivist

import (
	"path/filepath"
	"io"
	"os"
	"strings"
	"archive/tar"
)

type tarmonster struct {
    src    string
    dst    string
    writer *tar.Writer
    buffer []byte
}

func (t *tarmonster) do() error {
    out, err := os.Create(t.dst)
    if err != nil {
        return err
    }

    t.buffer = make([]byte, bufSize)

    t.writer = tar.NewWriter(out)
    if err := filepath.Walk(t.src, t.walk); err != nil {
        return err
    }

    if err := t.writer.Close(); err != nil {
        return err
    }

    return out.Close()
}

func (t *tarmonster) walk(path string, info os.FileInfo, err error) error {
    if err != nil { return err }

    if !info.Mode().IsRegular() || info.Size() == 0 {
        return nil
    }

    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    // Get tar.Header
    fih, err := tar.FileInfoHeader(info, "")
    if err != nil {
        return err
    }
    fih.Name = strings.TrimPrefix(path, t.src + string(filepath.Separator))
    
    // Begin a new file
    if err := t.writer.WriteHeader(fih); err != nil {
    	return err
    }

    // Write the file
    if _, err := io.CopyBuffer(t.writer, file, t.buffer); err != nil {
    	return err
    }

    return err
}
