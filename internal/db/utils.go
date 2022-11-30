package db

import (
	"fmt"
	"hash/crc32"
	"strings"
	"sync/atomic"

	"github.com/ToggyO/dbshaker/internal"
)

func GenerateLockID(dbName string, additional ...string) string {
	if len(additional) > 0 {
		dbName = strings.Join(append(additional, dbName), "\x00")
	}

	sum := crc32.ChecksumIEEE([]byte(dbName))
	sum *= uint32(internal.DBLockIDSalt)

	return fmt.Sprint(sum)
}

func CasRestoreOnError(lock *atomic.Bool, oldValue, newValue bool, casErr error, f func() error) error {
	if !lock.CompareAndSwap(oldValue, newValue) {
		return casErr
	}

	if err := f(); err != nil {
		lock.Store(oldValue)
		return err
	}

	return nil
}
