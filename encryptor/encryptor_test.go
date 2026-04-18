package encryptor_test

import (
	"os"
	"path/filepath"
	"testing"

	E "github.com/follow1123/secret-book/encryptor"
	"github.com/stretchr/testify/require"
)

const password = "123456"

func TestNewSuccess(t *testing.T) {
	dataPath := filepath.Join(os.TempDir(), "secrets")
	encryptor, err := E.New(dataPath, password)
	require.NoError(t, err)

	data, err := encryptor.GetData()
	require.NoError(t, err)
	require.Nil(t, data)
}

func TestNewFailure(t *testing.T) {

}
