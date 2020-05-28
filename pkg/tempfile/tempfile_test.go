package tempfile

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestWriteFile(t *testing.T) {
	filename, err := WriteFile([]byte("hello\n"), "test", ".txt", true)
	require.NoError(t, err)
	t.Log("filename:", filename)
	var buf []byte
	buf, err = ioutil.ReadFile(filename)
	require.NoError(t, err)
	require.Equal(t, []byte("hello\n"), buf)
	DeleteAll()
	_, err = ioutil.ReadFile(filename)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}

func TestWriteDirFile(t *testing.T) {
	dirname, filename, err := WriteDirFile([]byte("hello\nworld\n"), "test", "text.txt", false)
	require.NoError(t, err)
	var buf []byte
	buf, err = ioutil.ReadFile(filename)
	require.NoError(t, err)
	require.Equal(t, []byte("hello\nworld\n"), buf)
	DeleteAll()
	_, err = ioutil.ReadFile(filename)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
	_, err = ioutil.ReadFile(dirname)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}
