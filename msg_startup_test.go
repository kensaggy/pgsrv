package pgsrv

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartupVersion(t *testing.T) {
	t.Run("typed message", func(t *testing.T) {
		m := &msg{'p', 0, 0, 0, 5}
		_, err := m.StartupVersion()

		require.Error(t, err)
	})

	t.Run("untyped message", func(t *testing.T) {
		m := &msg{
			0, 0, 0, 8,
			4, 210, 22, 47,
		}
		version, err := m.StartupVersion()

		require.NoError(t, err)
		require.Equal(t, "1234.5679", version)
	})
}

func TestStartupArgs(t *testing.T) {
	t.Run("typed message", func(t *testing.T) {
		m := &msg{'p', 0, 0, 0, 5}
		_, err := m.StartupArgs()

		require.Error(t, err)
	})

	t.Run("untyped message", func(t *testing.T) {
		// an original message sent by psql client
		m := &msg{
			0, 0, 0, 84, 0, 3, 0, 0, 117, 115,
			101, 114, 0, 112, 111, 115, 116, 103, 114, 101,
			115, 0, 100, 97, 116, 97, 98, 97, 115, 101,
			0, 112, 111, 115, 116, 103, 114, 101, 115, 0,
			97, 112, 112, 108, 105, 99, 97, 116, 105, 111,
			110, 95, 110, 97, 109, 101, 0, 112, 115, 113,
			108, 0, 99, 108, 105, 101, 110, 116, 95, 101,
			110, 99, 111, 100, 105, 110, 103, 0, 85, 84,
			70, 56, 0, 0,
		}

		args, err := m.StartupArgs()
		expectedArgs := map[string]interface{}{
			"user":             "postgres",
			"database":         "postgres",
			"application_name": "psql",
			"client_encoding":  "UTF8",
		}

		require.NoError(t, err)
		require.Equal(t, expectedArgs, args)
	})

	t.Run("untyped message, no params", func(t *testing.T) {
		m := &msg{
			0, 0, 0, 9,
			1, 2, 3, 4,
			5, // this is an extra byte without a null terminator
		}

		args, err := m.StartupArgs()
		expectedArgs := make(map[string]interface{})

		require.NoError(t, err)
		require.Equal(t, expectedArgs, args)
	})
}

func TestIsTLSRequest(t *testing.T) {
	t.Run("tls", func(t *testing.T) {
		// an actual message with version 1234.5679
		m := &msg{0, 0, 0, 8, 4, 210, 22, 47}

		isTLS := m.IsTLSRequest()
		require.True(t, isTLS)
	})

	t.Run("not tls", func(t *testing.T) {
		// an actual message with version 1234.5679 with modified last byte
		m := &msg{0, 0, 0, 8, 4, 210, 22, 42}

		isTLS := m.IsTLSRequest()
		require.False(t, isTLS)
	})
}

func TestIsTerminate(t *testing.T) {
	t.Run("terminate", func(t *testing.T) {
		// an actual message with version 1234.5679
		m := &msg{'X', 0, 0, 0, 5}

		isTLS := m.IsTerminate()
		require.True(t, isTLS)
	})

	t.Run("not tls", func(t *testing.T) {
		// an actual message with version 1234.5679 with modified last byte
		m := &msg{'x', 0, 0, 0, 5}

		isTLS := m.IsTerminate()
		require.False(t, isTLS)
	})
}

func TestTlsResponseMsg(t *testing.T) {
	t.Run("supported", func(t *testing.T) {
		m := tlsResponseMsg(true)

		require.Equal(t, msg{'S'}, m)
	})

	t.Run("not supported", func(t *testing.T) {
		m := tlsResponseMsg(false)

		require.Equal(t, msg{'N'}, m)
	})
}

func TestAuthOKMsg(t *testing.T) {
	m := authOKMsg()

	require.Equal(t, msg{'R', 0, 0, 0, 8, 0, 0, 0, 0}, m)
}
