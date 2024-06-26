package cmd

import (
	"os"
	"testing"
	"time"
)

func TestSchedulerCommand(t *testing.T) {
	t.Run("Start the scheduler", func(t *testing.T) {
		tmpDir, _, _, _ := setupTest(t)
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		go func() {
			testRunCommand(t, schedulerCmd(), cmdTest{
				args:        []string{"scheduler"},
				expectedOut: []string{"starting dagu scheduler"},
			})
		}()

		time.Sleep(time.Millisecond * 500)
	})
}
