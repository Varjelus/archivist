package archivist

import (
    "fmt"
    "path/filepath"
    "math/rand"
    "testing"
)

func TestZipUnzip(t *testing.T) {
    var (
        err      error
        original string
        zipped   string
        unzipped string
    )

    zipped, err = filepath.Abs("./archivist_" + randString(rand.Intn(16)+5) + ".zip")
    if err != nil {
        t.Errorf("Initializing error: %s", err.Error())
    }

    unzipped, err = filepath.Abs("./archivist_" + randString(rand.Intn(16)+5) + "_unz")
    if err != nil {
        t.Errorf("Initializing error: %s", err.Error())
    }

    fmt.Println("ZipUnzip \t[Step 1/5] Spin up rand fs")
    original, err = randTempDirStruct()
    if err != nil {
        t.Errorf("Fabricating failed: %s", err.Error())
        goto cleanup
    }

    fmt.Println("ZipUnzip \t[Step 2/5] Zip")
    if err := Zip(original, zipped); err != nil {
        t.Errorf("Zipping failed: %s", err.Error())
        goto cleanup
    }

    fmt.Println("ZipUnzip \t[Step 3/5] Unzip")
    if err := Unzip(zipped, unzipped); err != nil {
        t.Errorf("Unzipping failed: %s", err.Error())
        goto cleanup
    }

    fmt.Println("ZipUnzip \t[Step 4/5] Compare")
    if err := (&comparer{original, unzipped}).compare(); err != nil {
        t.Errorf("Comparison failed: %s", err.Error())
        goto cleanup
    }

    fmt.Println("ZipUnzip \t[Step 5/5] Cleanup")
 cleanup:
    if err := cleanup(zipped, unzipped, original); err != nil {
        t.Errorf("cleanup failed: %s", err.Error())
    }
}