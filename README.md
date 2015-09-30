# archivist
Straightforward zipping for golang

`go get github.com/Varjelus/archivist`


## Notes
Using inefficient `io.Copy` at the moment.

## Usage
`import github.com/Varjelus/archivist`

### Methods
`archivist.Zip(sourcePath, destinationPath)`, handle error

`archivist.Unzip(sourcePath, destinationPath)`, handle error
