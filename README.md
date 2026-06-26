# zipstream

Small helper package for streaming filtered zip archives.

It is able to rebuild zip files while keeping only selected entries. It copies entry data 1:1, so files are not decompressed or recompressed.
This is useful when you want to dynamically serve a modified zip archive, without first writing the complete output to disk.

## Example

```go
// Assuming you have some sort of source zip file, e.g. a file object
source, _ := os.Open("source.zip")
info, _ := file.Stat()
size := info.Size()

resultStream, resultSize, err := zipstream.StreamFiltered(
	source, size,
	// Determine what to keep, e.g. only txt files in this case
	func(file *zip.File) bool {
		return strings.HasSuffix(file.Name, ".txt")
	},
)

// do whatever ...
```
