package cfpsql_test

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"errors"
	. "github.com/jaecktec/cf-psql-plugin/cfpsql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SshRunner", func() {
	var cliConnection *pluginfakes.FakeCliConnection
	var sshRunner SshRunner
	service := PsqlService{
		Name:     "database-a",
		Hostname: "database-a.host",
		Port:     "5432",
		DbName:   "dbname-a",
		Username: "username-a",
		Password: "password-a",
	}

	BeforeEach(func() {
		cliConnection = new(pluginfakes.FakeCliConnection)
		sshRunner = NewSshRunner()
	})

	Context("When opening the tunnel", func() {
		It("Runs 'cf ssh'", func() {
			sshRunner.OpenSshTunnel(cliConnection, service, "app-name", 4242)

			Expect(cliConnection.CliCommandCallCount()).To(Equal(1))
			calledArgs := cliConnection.CliCommandArgsForCall(0)
			Expect(calledArgs).To(Equal([]string{"ssh", "app-name", "-N", "-L", "4242:database-a.host:5432"}))
		})
	})

	Context("When 'cf ssh' returns an error", func() {
		It("panics", func() {
			defer func() {
				if thePanic := recover(); thePanic != nil {
					Expect(thePanic).To(Equal(errors.New("SSH tunnel failed: PC LOAD LETTER")))
				} else {
					Fail("Expected function to panic")
				}
			}()
			cliConnection.CliCommandWithoutTerminalOutputStub = nil
			cliConnection.CliCommandReturns(nil, errors.New("PC LOAD LETTER"))

			sshRunner.OpenSshTunnel(cliConnection, service, "app-name", 4242)
		})
	})
})
