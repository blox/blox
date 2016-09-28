package store

type Stores struct {
	TaskStore              TaskStore
	ContainerInstanceStore ContainerInstanceStore
}

func NewStores(datastore DataStore) (Stores, error) {
	taskStore, err := NewTaskStore(datastore)
	if err != nil {
		return Stores{}, err
	}

	containerInstanceStore, err := NewContainerInstanceStore(datastore)
	if err != nil {
		return Stores{}, err
	}

	return Stores{
		TaskStore:              taskStore,
		ContainerInstanceStore: containerInstanceStore,
	}, nil
}
