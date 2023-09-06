package fsplit

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWriter(t *testing.T) {
	entries, _ := os.ReadDir("testdata")
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "test.bin") {
			require.NoError(t, os.RemoveAll(entry.Name()))
		}
	}

	name := filepath.Join("testdata", "test.bin")

	w, err := NewWriter(name, WriterOptions{SplitSize: 10, Perm: 0640})
	require.NoError(t, err)
	defer w.Close()

	// write
	n, err := w.Write([]byte("012345678"))
	require.NoError(t, err)
	require.Equal(t, 9, n)
	require.NoError(t, w.Sync())

	// check
	buf, err := os.ReadFile(name)
	require.NoError(t, err)
	require.Equal(t, "012345678", string(buf))

	// files
	files, err := w.Files()
	require.NoError(t, err)
	require.Equal(t, []string{name}, files)

	// write
	n, err = w.Write([]byte("012345678"))
	require.NoError(t, err)
	require.Equal(t, 9, n)
	require.NoError(t, w.Sync())

	// check
	_, err = os.ReadFile(name)
	require.Error(t, err)
	require.ErrorIs(t, err, os.ErrNotExist)

	buf, err = os.ReadFile(name + ".1")
	require.NoError(t, err)
	require.Equal(t, "0123456780", string(buf))

	buf, err = os.ReadFile(name + ".2")
	require.NoError(t, err)
	require.Equal(t, "12345678", string(buf))

	// files
	files, err = w.Files()
	require.NoError(t, err)
	require.Equal(t, []string{name + ".1", name + ".2"}, files)

	// write
	n, err = w.Write([]byte("012345678"))
	require.NoError(t, err)
	require.Equal(t, 9, n)
	require.NoError(t, w.Sync())

	// check
	_, err = os.ReadFile(name)
	require.Error(t, err)
	require.ErrorIs(t, err, os.ErrNotExist)

	buf, err = os.ReadFile(name + ".1")
	require.NoError(t, err)
	require.Equal(t, "0123456780", string(buf))

	buf, err = os.ReadFile(name + ".2")
	require.NoError(t, err)
	require.Equal(t, "1234567801", string(buf))

	buf, err = os.ReadFile(name + ".3")
	require.NoError(t, err)
	require.Equal(t, "2345678", string(buf))

	// files
	files, err = w.Files()
	require.NoError(t, err)
	require.Equal(t, []string{name + ".1", name + ".2", name + ".3"}, files)
}
