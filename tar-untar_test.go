package archivist

import (
    "fmt"
    "path/filepath"
    "math/rand"
    "testing"
)

func init() {
    
}

func TestTarUntar(t *testing.T) {
    var (
        tarred   string
        untarred string
        err      error
        original string
    )

    tarred, err = filepath.Abs("./archivist_" + randString(rand.Intn(16)+5) + ".tar.gz")
    if err != nil {
        t.Errorf("Initializing error: %s", err.Error())
    }

    untarred, err = filepath.Abs("./archivist_" + randString(rand.Intn(16)+5) + "_unt")
    if err != nil {
        t.Errorf("Initializing error: %s", err.Error())
    }

    fmt.Println("TarUntar \t[Step 1/5] Spin up rand fs")
    original, err = randTempDirStruct()
    if err != nil {
        t.Errorf("Fabricating failed: %s", err.Error())
        goto cleanup
    }

    fmt.Println("TarUntar \t[Step 2/5] Tar")
    if err := Tar(original, tarred); err != nil {
        t.Errorf("Tarring failed: %s", err.Error())
        goto cleanup
    }

    fmt.Println("TarUntar \t[Step 3/5] Untar")
    if err := Untar(tarred, untarred); err != nil {
        t.Errorf("Untarring failed: %s", err.Error())
        goto cleanup
    }

    fmt.Println("TarUntar \t[Step 4/5] Compare")
    if err := (&comparer{original, untarred}).compare(); err != nil {
        t.Errorf("Comparison failed: %s", err.Error())
        goto cleanup
    }

    fmt.Println("TarUntar \t[Step 5/5] Cleanup")
 cleanup:
    if err := cleanup(tarred, untarred, original); err != nil {
        t.Errorf("cleanup failed: %s", err.Error())
    }
}
