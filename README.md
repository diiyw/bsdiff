# BSDiff

BSDiff is a binary diff tool that generates a patch between two binary files. The patch can be used to transform the first binary file into the second binary file. The patch is usually much smaller than the second binary file, making it easier to distribute updates to the binary file.

## Usage

```go
// create a patch
patch := bsdiff.Diff(old, new)
// returns the patch as a byte slice
// patch.ToBytes()
// bsdiff.FromBytes(patch.ToBytes())
patch.Apply(oldFile)
```