package db

import (
	"fmt"
	"hash/crc32"
	"strings"
	"sync/atomic"

	"github.com/ToggyO/dbshaker/internal"
)

func GenerateLockId(dbName string, additional ...string) string {
	if len(additional) > 0 {
		dbName = strings.Join(append(additional, dbName), "\x00")
	}

	sum := crc32.ChecksumIEEE([]byte(dbName))
	sum = sum * uint32(internal.DbLockIDSalt)

	return fmt.Sprint(sum)
}

func CasRestoreOnError(lock *atomic.Bool, old, new bool, casErr error, f func() error) error {
	if !lock.CompareAndSwap(old, new) {
		return casErr
	}

	if err := f(); err != nil {
		lock.Store(old)
		return err
	}

	return nil
}
