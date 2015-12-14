package api_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/api"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	fakeSysbenchClient "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client/fakes"
	"github.com/tedsuo/rata"
)

var _ = Describe("Api", func() {
	var (
		sysbenchClient *fakeSysbenchClient.FakeSysbenchClient
		ts             *httptest.Server
	)

	BeforeEach(func() {
		sysbenchClient = &fakeSysbenchClient.FakeSysbenchClient{}
		testConfig := &config.Config{}

		handler, err := api.NewRouter(api.Api{
			RootConfig:     testConfig,
			SysbenchClient: sysbenchClient,
		})

		Expect(err).ToNot(HaveOccurred())
		ts = httptest.NewServer(handler)
	})

	AfterEach(func() {
		ts.Close()
	})

	var createReq = func(routeName, targetNode string) *http.Request {
		routes := api.DefaultRoutes()
		requestGenerator := rata.NewRequestGenerator(ts.URL, routes)

		req, err := requestGenerator.CreateRequest(
			routeName,
			rata.Params{"node": targetNode},
			nil,
		)
		Expect(err).NotTo(HaveOccurred())

		return req
	}

	Context("POST /start", func() {
		It("sends a start command to the sysbench client", func() {
			req := createReq("start_test", "some-node")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(sysbenchClient.StartCallCount()).To(Equal(1))
			Expect(sysbenchClient.StartArgsForCall(0)).To(Equal("some-node"))
		})
	})
})
