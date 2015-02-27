package scheduler

import (
	"10gen.com/mci"
	"10gen.com/mci/model"
	"github.com/10gen-labs/slogger/v1"
)

// Interface responsible for finding all tasks that are ready to be run.
type TaskFinder interface {
	// Returns a slice of tasks that are ready to be run, and an error if
	// appropriate.
	FindRunnableTasks() ([]model.Task, error)
}

// Implementation that fetches tasks from the database.
type DBTaskFinder struct{}

// Find all tasks that are ready to be run.  Works by fetching all undispatched
// tasks from the database, and filtering out any whose dependencies are not
// met.
func (self *DBTaskFinder) FindRunnableTasks() ([]model.Task, error) {

	// find all of the undispatched tasks
	undispatchedTasks, err := model.FindUndispatchedTasks()
	if err != nil {
		return nil, err
	}

	// filter out any tasks whose dependencies are not met
	runnableTasks := make([]model.Task, 0, len(undispatchedTasks))
	dependencyCaches := make(map[string]model.Task)
	for _, task := range undispatchedTasks {
		depsMet, err := task.DependenciesMet(dependencyCaches)
		if err != nil {
			mci.Logger.Logf(slogger.ERROR, "Error checking dependencies for"+
				" task %v: %v", task.Id, err)
			continue
		}
		if depsMet {
			runnableTasks = append(runnableTasks, task)
		}
	}

	return runnableTasks, nil
}