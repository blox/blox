package features

import (
	. "github.com/gucumber/gucumber"
)

//TODO: implement all features
//use either a library to call REST APIs or generate a go client for ESH
func init() {
	Given(`^I put (\d+) tasks in the queue$`, func(numTasks int) {
		T.Skip() // pending
	})

	When(`^I get task with the same arn$`, func() {
		T.Skip() // pending
	})

	Then(`^I get a task that matches the task I pushed to the queue$`, func() {
		T.Skip() // pending
	})

	When(`^I list tasks$`, func() {
		T.Skip() // pending
	})

	Then(`^I get a list of tasks that includes the tasks I pushed to the queue$`, func() {
		T.Skip() // pending
	})

	Given(`^I put the following tasks in the queue: (\d+) pending, (\d+) running, and (\d+) stopped$`, func(numPending int, numRunning int, numStopped int) {
		T.Skip() // pending
	})

	When(`^I filter tasks by (.+?)$`, func(status string) {
		T.Skip() // pending
	})

	Then(`^I get (\d+) number of tasks$`, func(numTasks int) {
		T.Skip() // pending
	})

	And(`^each task matches the corresponding task pushed to the queue$`, func() {
		T.Skip() // pending
	})
}
