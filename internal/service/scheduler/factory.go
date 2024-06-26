package scheduler

import (
	"time"

	"github.com/dagu-dev/dagu/internal/dag"
	"github.com/dagu-dev/dagu/internal/engine"
	"github.com/dagu-dev/dagu/internal/service/scheduler/job"
	"github.com/dagu-dev/dagu/internal/service/scheduler/scheduler"
)

type jobFactory struct {
	Executable string
	WorkDir    string
	Engine     engine.Engine
}

func (jf jobFactory) NewJob(dg *dag.DAG, next time.Time) scheduler.Job {
	return &job.Job{
		DAG:        dg,
		Executable: jf.Executable,
		WorkDir:    jf.WorkDir,
		Next:       next,
		Engine:     jf.Engine,
	}
}
