package archivist

import (
	"path/filepath"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Essentially from stackoverflow.com/questions/22892120
const charBytes = "01234567890_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    charIdxBits = 6                    // 6 bits to represent a char index
    charIdxMask = 1<<charIdxBits - 1   // All 1-bits, as many as charIdxBits
    charIdxMax  = 63 / charIdxBits    // # of char indices fitting in 63 bits
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func randString(n int) string {
    b := make([]byte, n)
    // A src.Int63() generates 63 random bits, enough for charIdxMax characters!
    for i, cache, remain := n-1, rand.Int63(), charIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = rand.Int63(), charIdxMax
        }
        if idx := int(cache & charIdxMask); idx < len(charBytes) {
            b[i] = charBytes[idx]
            i--
        }
        cache >>= charIdxBits
        remain--
    }

    return string(b)
}

type comparer struct {
    first  string
    second string
}

func (c *comparer) compare() error {
    return filepath.Walk(c.second, func(path string, info os.FileInfo, err error) error {
        if err != nil { return err }
        if info.IsDir() { return nil }

        originalPath := filepath.Join(c.first, strings.TrimPrefix(path, c.second + string(filepath.Separator)))

        shouldBeSame, err := os.Stat(originalPath)
        if err != nil {
            return err
        }

        if info.Size() != shouldBeSame.Size() {
            return fmt.Errorf("Unmatching FileInfos")
        }

        return nil
    })
}

func randTempDirStruct() (string, error) {
    root, err := filepath.Abs("./archivist_" + randString(rand.Intn(16)+5))
    if err != nil {
        return root, err
    }

    if err := os.MkdirAll(root, perm); err != nil {
        return root, err
    }

    if err := populateWithFiles(root); err != nil {
        return root, err
    }

    for i := 0; i < (rand.Intn(10)+1); i++ {
        child := randString(rand.Intn(16)+5)
        cPath := filepath.Join(root, child)

        if err := os.MkdirAll(cPath, perm); err != nil {
            return root, err
        }
        if err := populateWithFiles(cPath); err != nil {
            return root, err
        }

        for j := 0; j < (rand.Intn(10)+1); j++ {
            grandchild := randString(rand.Intn(16)+5)
            gcPath := filepath.Join(cPath, grandchild)

            if err := os.MkdirAll(gcPath, perm); err != nil {
                return root, err
            }
            if err := populateWithFiles(gcPath); err != nil {
                return root, err
            }
        }
    }

    return root, nil
}

func populateWithFiles(path string) error {
    for i := 0; i < (rand.Intn(10)+1); i++ {
        f, err := os.Create(filepath.Join(path, randString(rand.Intn(16)+5)))
        if err != nil {
            return err
        }

        if _, err = f.Write([]byte(randString(rand.Intn(5000000)+4096))); err != nil {
            return err
        }

        return f.Close()
    }

    return nil
}

func cleanup(paths ...string) error {
    for _, path := range paths {
        if err := os.RemoveAll(path); err != nil {
            return err
        }
    }
    return nil
}