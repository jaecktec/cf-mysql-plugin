package cfpsql_test

import (
	. "github.com/jaecktec/cf-psql-plugin/cfpsql"

	"errors"
	"github.com/jaecktec/cf-psql-plugin/cfpsql/cfpsqlfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("PsqlRunner", func() {
	Context("RunPsql", func() {
		var exec *cfpsqlfakes.FakeExecWrapper
		var runner PsqlRunner

		BeforeEach(func() {
			exec = new(cfpsqlfakes.FakeExecWrapper)
			runner = NewPsqlRunner(exec)
		})

		Context("When psql is not in PATH", func() {
			It("Returns an error", func() {
				exec.LookPathReturns("", errors.New("PC LOAD LETTER"))

				err := runner.RunPsql("hostname", 42, "dbname", "username", "password")

				Expect(err).To(Equal(errors.New("'psql' client not found in PATH")))
				Expect(exec.LookPathArgsForCall(0)).To(Equal("psql"))
			})
		})

		Context("When Run returns an error", func() {
			It("Forwards the error", func() {
				exec.LookPathReturns("/path/to/psql", nil)
				exec.RunReturns(errors.New("PC LOAD LETTER"))

				err := runner.RunPsql("hostname", 42, "dbname", "username", "password")

				Expect(err).To(Equal(errors.New("error running psql client: PC LOAD LETTER")))
			})
		})

		Context("When psql is in PATH", func() {
			It("Calls psql with the right arguments", func() {
				exec.LookPathReturns("/path/to/psql", nil)

				err := runner.RunPsql("hostname", 42, "dbname", "username", "password")

				Expect(err).To(BeNil())
				Expect(exec.LookPathCallCount()).To(Equal(1))
				Expect(exec.RunCallCount()).To(Equal(1))

				cmd := exec.RunArgsForCall(0)
				Expect(cmd.Path).To(Equal("/path/to/psql"))
				Expect(cmd.Args).To(Equal([]string{"/path/to/psql", "-U", "username", "-h", "hostname", "-p", "42", "-d", "dbname"}))
				Expect(cmd.Env).To(ContainElement("PGPASSWORD=password"))
				Expect(cmd.Stdin).To(Equal(os.Stdin))
				Expect(cmd.Stdout).To(Equal(os.Stdout))
				Expect(cmd.Stderr).To(Equal(os.Stderr))
			})
		})

		Context("When psql is in PATH and additional arguments are passed", func() {
			It("Calls psql with the right arguments", func() {
				exec.LookPathReturns("/path/to/psql", nil)

				err := runner.RunPsql("hostname", 42, "dbname", "username", "password", "--foo", "bar", "--baz")

				Expect(err).To(BeNil())
				Expect(exec.LookPathCallCount()).To(Equal(1))
				Expect(exec.RunCallCount()).To(Equal(1))

				cmd := exec.RunArgsForCall(0)
				Expect(cmd.Path).To(Equal("/path/to/psql"))
				Expect(cmd.Args).To(Equal([]string{"/path/to/psql", "-U", "username", "-h", "hostname", "-p", "42", "-d", "dbname", "--foo", "bar", "--baz"}))
				Expect(cmd.Env).To(ContainElement("PGPASSWORD=password"))
				Expect(cmd.Stdin).To(Equal(os.Stdin))
				Expect(cmd.Stdout).To(Equal(os.Stdout))
				Expect(cmd.Stderr).To(Equal(os.Stderr))
			})
		})
	})
})
