package storage

type FakeRobotStorage struct {
}

func NewFakeRobotStorage() *FakeRobotStorage {
	return &FakeRobotStorage{}
}

func (r *FakeRobotStorage) Allowed(url string) bool {
	return true
}
