package cfpsql_test

import (
	. "github.com/jaecktec/cf-psql-plugin/cfpsql"

	"errors"
	"github.com/jaecktec/cf-psql-plugin/cfpsql/cfpsqlfakes"
	"github.com/jaecktec/cf-psql-plugin/cfpsql/netfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
)

var _ = Describe("PortWaiter", func() {
	var netWrapper *cfpsqlfakes.FakeNetWrapper
	var portWaiter PortWaiter

	const SucceedAfterTries = 5
	dialCount := 0
	mockDial := func(network, address string) (net.Conn, error) {
		if dialCount < SucceedAfterTries-1 {
			dialCount++
			return nil, errors.New("GURU MEDITATION")
		}

		return new(netfakes.FakeConn), nil
	}

	BeforeEach(func() {
		netWrapper = new(cfpsqlfakes.FakeNetWrapper)
		portWaiter = NewPortWaiter(netWrapper)
	})

	It("Waits until the port is open", func() {
		netWrapper.DialStub = mockDial

		portWaiter.WaitUntilOpen(523)

		Expect(netWrapper.DialCallCount()).To(Equal(SucceedAfterTries))
	})

	It("Closes the connection", func() {
		mockConn := new(netfakes.FakeConn)
		netWrapper.DialReturns(mockConn, nil)

		portWaiter.WaitUntilOpen(523)

		Expect(netWrapper.CloseCallCount()).To(Equal(1))
		Expect(netWrapper.CloseArgsForCall(0)).To(Equal(mockConn))
		Expect(netWrapper.DialCallCount()).To(Equal(1))
	})
})
