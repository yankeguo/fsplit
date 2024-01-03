# fsplit

a file writer in Golang with auto split by size

## Usage

```go
package main

import "github.com/yankeguo/fsplit"

func main() {
	w, _ := fsplit.NewWriter("test.bin", fsplit.WriterOptions{
		Perm:      0640,
		SplitSize: 10,
	})
	_, _ = w.Write([]byte("012345678"))
	_ = w.Sync()
	// Files:
	//   test.bin (9 byte): 012345678
	_, _ = w.Write([]byte("012345678"))
	_ = w.Sync()
	// Files:
	//   test.bin.1 (10 byte): 0123456780
	//   test.bin.2 (8 byte): 12345678
	_ = w.Close()
}
```

## Credits

GUO YANKE, MIT License
