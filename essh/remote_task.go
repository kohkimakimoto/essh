package essh

// WIP...

type RemoteTask struct {
	URL string
	Namespace string
}

var RemoteTasks []*RemoteTask = []*RemoteTask{}

func NewRemoteTask(url string) *RemoteTask {
	return &RemoteTask{URL:url}
}
