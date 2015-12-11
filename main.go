package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"

	"github.com/cloudfoundry-incubator/switchboard/api"
	"github.com/cloudfoundry-incubator/switchboard/config"
	"github.com/cloudfoundry-incubator/switchboard/domain"
	"github.com/cloudfoundry-incubator/switchboard/health"
	"github.com/cloudfoundry-incubator/switchboard/proxy"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"

	"time"

	"github.com/pivotal-golang/lager"
)

func main() {

}
