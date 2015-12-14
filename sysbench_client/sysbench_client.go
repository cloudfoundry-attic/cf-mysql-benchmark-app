package sysbench_client

type SysbenchClient interface {
	Start(string) (string, error)
}

type sysbenchClient struct {}

func New() *sysbenchClient {
	return &sysbenchClient{}
}
