package api_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/api"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	fakeSysbenchClient "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client/fakes"
	"github.com/tedsuo/rata"
	"io/ioutil"
	"strconv"
)

var _ = Describe("Api", func() {

	const nodeIndex = 1

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

	var createReq = func(routeName string) *http.Request {
		routes := api.DefaultRoutes()
		requestGenerator := rata.NewRequestGenerator(ts.URL, routes)

		req, err := requestGenerator.CreateRequest(
			routeName,
			rata.Params{"node": strconv.Itoa(nodeIndex)},
			nil,
		)
		Expect(err).NotTo(HaveOccurred())

		return req
	}

	Describe("POST /start", func() {
		Context("when sysbench runs successfully", func() {

			BeforeEach(func() {
				sysbenchClient.StartReturns("ran successfully", nil)
			})

			It("sends a start command to the sysbench client", func() {
				req := createReq("start_test")
				resp, err := http.DefaultClient.Do(req)
				Expect(err).ToNot(HaveOccurred())

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(body).To(ContainSubstring("ran successfully"))
				Expect(sysbenchClient.StartCallCount()).To(Equal(1))
				Expect(sysbenchClient.StartArgsForCall(0)).To(Equal(nodeIndex))
			})
		})

		Context("when sysbench returns an error", func() {

			BeforeEach(func() {
				sysbenchClient.StartReturns("fake-stderr", errors.New("fake-error"))
			})

			It("sends a start command to the sysbench client", func() {
				req := createReq("start_test")
				resp, err := http.DefaultClient.Do(req)
				Expect(err).ToNot(HaveOccurred())

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))

				Expect(body).To(ContainSubstring("fake-stderr"))
				Expect(body).To(ContainSubstring("fake-error"))
				Expect(sysbenchClient.StartCallCount()).To(Equal(1))
				Expect(sysbenchClient.StartArgsForCall(0)).To(Equal(nodeIndex))
			})
		})
	})

	Describe("POST /prepare", func() {
		Context("when sysbench runs successfully", func() {

			BeforeEach(func() {
				sysbenchClient.PrepareReturns("", nil)
			})

			It("sends a prepare command to the sysbench client", func() {
				req := createReq("prepare_test")
				resp, err := http.DefaultClient.Do(req)
				Expect(err).ToNot(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(sysbenchClient.PrepareCallCount()).To(Equal(1))
				Expect(sysbenchClient.PrepareArgsForCall(0)).To(Equal(nodeIndex))
			})
		})

		Context("when sysbench returns an error", func() {

			BeforeEach(func() {
				sysbenchClient.PrepareReturns("", errors.New("fake-error"))
			})

			It("sends a prepare command to the sysbench client", func() {
				req := createReq("prepare_test")
				resp, err := http.DefaultClient.Do(req)
				Expect(err).ToNot(HaveOccurred())

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))

				Expect(body).To(ContainSubstring("fake-error"))
				Expect(sysbenchClient.PrepareCallCount()).To(Equal(1))
				Expect(sysbenchClient.PrepareArgsForCall(0)).To(Equal(nodeIndex))
			})
		})
	})
})
