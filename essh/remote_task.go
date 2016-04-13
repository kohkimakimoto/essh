package essh

type RemoteTask struct {
	URL string
}

var RemoteTasks []*RemoteTask = []*RemoteTask{}

func NewRemoteTask(url string) *RemoteTask {
	return &RemoteTask{URL:url}
}
