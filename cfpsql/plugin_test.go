package cfpsql_test

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"fmt"
	. "github.com/jaecktec/cf-psql-plugin/cfpsql"
	"github.com/jaecktec/cf-psql-plugin/cfpsql/cfpsqlfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Plugin", func() {
	var appList []plugin_models.GetAppsModel
	usage := "cf psql - Connect to a PSQL database service\n\nUSAGE:\n   Open a psql client to a database:\n   cf psql <service-name> [psql args...]\n"

	BeforeEach(func() {
		appList = []plugin_models.GetAppsModel{
			{
				Name: "app-name-1",
			},
			{
				Name: "app-name-2",
			},
		}
	})

	Context("When calling 'cf plugins'", func() {
		It("Shows the psql plugin with the current version", func() {
			psqlPlugin, _ := NewPluginAndMocks()

			Expect(psqlPlugin.GetMetadata().Name).To(Equal("psql"))
			Expect(psqlPlugin.GetMetadata().Version).To(Equal(plugin.VersionType{
				Major: 2,
				Minor: 0,
				Build: 0,
			}))
		})
	})

	Context("When calling 'cf psql -h'", func() {
		It("Shows instructions for 'cf psql'", func() {
			psqlPlugin, _ := NewPluginAndMocks()

			Expect(psqlPlugin.GetMetadata().Commands).To(HaveLen(1))
			Expect(psqlPlugin.GetMetadata().Commands[0].Name).To(Equal("psql"))
		})
	})

	Context("When calling 'cf psql' without arguments", func() {
		It("Prints usage instructions to STDERR and exits with 1", func() {
			psqlPlugin, mocks := NewPluginAndMocks()

			psqlPlugin.Run(mocks.CliConnection, []string{"psql"})

			Expect(mocks.Out).To(gbytes.Say(""))
			Expect(string(mocks.Err.Contents())).To(Equal(usage))
			Expect(psqlPlugin.GetExitCode()).To(Equal(1))
		})
	})

	Context("When calling 'cf psql db-name'", func() {
		var serviceA PsqlService

		BeforeEach(func() {
			serviceA = PsqlService{
				Name:     "database-a",
				Hostname: "database-a.host",
				Port:     "123",
				DbName:   "dbname-a",
				Username: "username",
				Password: "password",
			}
		})

		Context("When the database is available", func() {
			It("Opens an SSH tunnel through a started app", func() {
				psqlPlugin, mocks := NewPluginAndMocks()

				mocks.CfService.GetServiceReturns(serviceA, nil)
				mocks.CfService.GetStartedAppsReturns(appList, nil)
				mocks.PortFinder.GetPortReturns(2342)

				psqlPlugin.Run(mocks.CliConnection, []string{"psql", "database-a"})

				Expect(mocks.CfService.GetServiceCallCount()).To(Equal(1))
				calledCliConnection, calledName := mocks.CfService.GetServiceArgsForCall(0)
				Expect(calledName).To(Equal("database-a"))
				Expect(calledCliConnection).To(Equal(mocks.CliConnection))

				Expect(mocks.CfService.GetStartedAppsCallCount()).To(Equal(1))
				Expect(mocks.PortFinder.GetPortCallCount()).To(Equal(1))
				Expect(mocks.CfService.OpenSshTunnelCallCount()).To(Equal(1))

				calledCliConnection, calledService, calledAppList, localPort := mocks.CfService.OpenSshTunnelArgsForCall(0)
				Expect(calledCliConnection).To(Equal(mocks.CliConnection))
				Expect(calledService).To(Equal(serviceA))
				Expect(calledAppList).To(Equal(appList))
				Expect(localPort).To(Equal(2342))
			})

			It("Opens a PSQL client connecting through the tunnel", func() {
				psqlPlugin, mocks := NewPluginAndMocks()

				mocks.CfService.GetServiceReturns(serviceA, nil)
				mocks.CfService.GetStartedAppsReturns(appList, nil)
				mocks.PortFinder.GetPortReturns(2342)

				psqlPlugin.Run(mocks.CliConnection, []string{"psql", "database-a"})

				Expect(mocks.PortFinder.GetPortCallCount()).To(Equal(1))
				Expect(mocks.PsqlRunner.RunPsqlCallCount()).To(Equal(1))

				hostname, port, dbName, username, password, _ := mocks.PsqlRunner.RunPsqlArgsForCall(0)
				Expect(hostname).To(Equal("127.0.0.1"))
				Expect(port).To(Equal(2342))
				Expect(dbName).To(Equal(serviceA.DbName))
				Expect(username).To(Equal(serviceA.Username))
				Expect(password).To(Equal(serviceA.Password))
			})

			Context("When passing additional arguments", func() {
				It("Passes the arguments to psql", func() {
					psqlPlugin, mocks := NewPluginAndMocks()

					mocks.CfService.GetServiceReturns(serviceA, nil)
					mocks.CfService.GetStartedAppsReturns(appList, nil)
					mocks.PortFinder.GetPortReturns(2342)

					psqlPlugin.Run(mocks.CliConnection, []string{"psql", "database-a", "--foo", "bar", "--baz"})

					Expect(mocks.PortFinder.GetPortCallCount()).To(Equal(1))
					Expect(mocks.PsqlRunner.RunPsqlCallCount()).To(Equal(1))

					hostname, port, dbName, username, password, args := mocks.PsqlRunner.RunPsqlArgsForCall(0)
					Expect(hostname).To(Equal("127.0.0.1"))
					Expect(port).To(Equal(2342))
					Expect(dbName).To(Equal(serviceA.DbName))
					Expect(username).To(Equal(serviceA.Username))
					Expect(password).To(Equal(serviceA.Password))
					Expect(args).To(Equal([]string{"--foo", "bar", "--baz"}))
				})
			})
		})

		Context("When a service key cannot be retrieved", func() {
			It("Shows an error message and exits with 1", func() {
				psqlPlugin, mocks := NewPluginAndMocks()

				mocks.CfService.GetServiceReturns(PsqlService{}, errors.New("database not found"))

				psqlPlugin.Run(mocks.CliConnection, []string{"psql", "db-name"})

				Expect(mocks.CfService.GetServiceCallCount()).To(Equal(1))
				Expect(mocks.Out).To(gbytes.Say(""))
				Expect(mocks.Err).To(gbytes.Say("^FAILED\nUnable to retrieve service credentials: database not found\n$"))
				Expect(psqlPlugin.GetExitCode()).To(Equal(1))
			})
		})

		Context("When there are no started apps", func() {
			It("Shows an error message and exits with 1", func() {
				psqlPlugin, mocks := NewPluginAndMocks()

				mocks.CfService.GetServiceReturns(serviceA, nil)
				mocks.CfService.GetStartedAppsReturns([]plugin_models.GetAppsModel{}, nil)

				psqlPlugin.Run(mocks.CliConnection, []string{"psql", "database-a"})

				Expect(mocks.CfService.GetServiceCallCount()).To(Equal(1))
				Expect(mocks.CfService.GetStartedAppsCallCount()).To(Equal(1))
				Expect(mocks.Out).To(gbytes.Say("^$"))
				Expect(mocks.Err).To(gbytes.Say("^FAILED\nUnable to connect to 'database-a': no started apps in current space\n$"))
				Expect(psqlPlugin.GetExitCode()).To(Equal(1))
			})
		})

		Context("When GetStartedApps returns an error", func() {
			It("Shows an error message and exits with 1", func() {
				psqlPlugin, mocks := NewPluginAndMocks()

				mocks.CfService.GetServiceReturns(serviceA, nil)
				mocks.CfService.GetStartedAppsReturns(nil, fmt.Errorf("PC LOAD LETTER"))

				psqlPlugin.Run(mocks.CliConnection, []string{"psql", "database-a"})

				Expect(mocks.CfService.GetServiceCallCount()).To(Equal(1))
				Expect(mocks.CfService.GetStartedAppsCallCount()).To(Equal(1))
				Expect(mocks.Out).To(gbytes.Say(""))
				Expect(mocks.Err).To(gbytes.Say("^FAILED\nUnable to retrieve started apps: PC LOAD LETTER\n$"))
				Expect(psqlPlugin.GetExitCode()).To(Equal(1))
			})
		})

	})

	Context("When uninstalling the plugin", func() {
		It("Does not give any output or call the API", func() {
			psqlPlugin, mocks := NewPluginAndMocks()

			psqlPlugin.Run(mocks.CliConnection, []string{"CLI-MESSAGE-UNINSTALL"})

			Expect(mocks.Out).To(gbytes.Say("^$"))
			Expect(mocks.Err).To(gbytes.Say("^$"))
			Expect(psqlPlugin.GetExitCode()).To(Equal(0))
		})
	})
})

type Mocks struct {
	In            *gbytes.Buffer
	Out           *gbytes.Buffer
	Err           *gbytes.Buffer
	CfService     *cfpsqlfakes.FakeCfService
	PortFinder    *cfpsqlfakes.FakePortFinder
	CliConnection *pluginfakes.FakeCliConnection
	PsqlRunner   *cfpsqlfakes.FakePsqlRunner
}

func NewPluginAndMocks() (*PsqlPlugin, Mocks) {
	mocks := Mocks{
		In:            gbytes.NewBuffer(),
		Out:           gbytes.NewBuffer(),
		Err:           gbytes.NewBuffer(),
		CfService:     new(cfpsqlfakes.FakeCfService),
		CliConnection: new(pluginfakes.FakeCliConnection),
		PsqlRunner:   new(cfpsqlfakes.FakePsqlRunner),
		PortFinder:    new(cfpsqlfakes.FakePortFinder),
	}

	psqlPlugin := NewPsqlPlugin(PluginConf{
		In:         mocks.In,
		Out:        mocks.Out,
		Err:        mocks.Err,
		CfService:  mocks.CfService,
		PsqlRunner: mocks.PsqlRunner,
		PortFinder: mocks.PortFinder,
	})

	return psqlPlugin, mocks
}
