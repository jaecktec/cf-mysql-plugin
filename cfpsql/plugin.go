package cfpsql

import (
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/plugin/models"
	"fmt"
	"io"
)

type PsqlPlugin struct {
	In         io.Reader
	Out        io.Writer
	Err        io.Writer
	CfService  CfService
	PsqlRunner PsqlRunner
	PortFinder PortFinder
	exitCode   int
}

func NewPsqlPlugin(conf PluginConf) *PsqlPlugin {
	return &PsqlPlugin{
		In:         conf.In,
		Out:        conf.Out,
		Err:        conf.Err,
		CfService:  conf.CfService,
		PortFinder: conf.PortFinder,
		PsqlRunner: conf.PsqlRunner,
	}
}

func (self *PsqlPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "psql-plugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "psql",
				HelpText: "Connect to a PSQL database service",
				UsageDetails: plugin.Usage{
					Usage: "Open a psql client to a database:\n   " +
						"cf psql <service-name> [psql args...]",
				},
			},
		},
	}
}

func (self *PsqlPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	command := args[0]

	switch command {
	case "psql":
		if len(args) > 1 {
			dbName := args[1]

			var psqlArgs []string
			if len(args) > 2 {
				psqlArgs = args[2:]
			}

			self.connectTo(cliConnection, command, dbName, psqlArgs)
		} else {
			fmt.Fprint(self.Err, self.FormatUsage())
			self.setErrorExit()
		}

	default:
		// we don't handle "uninstall"
	}
}

func (self *PsqlPlugin) FormatUsage() string {
	var usage string
	for i, command := range self.GetMetadata().Commands {
		if i > 0 {
			usage += "\n\n"
		}

		usage += fmt.Sprintf(
			"cf %s - %s\n\nUSAGE:\n   %s\n",
			command.Name,
			command.HelpText,
			command.UsageDetails.Usage,
		)
	}

	return usage
}

func (self *PsqlPlugin) GetExitCode() int {
	return self.exitCode
}

func (self *PsqlPlugin) setErrorExit() {
	self.exitCode = 1
}

type StartedAppsResult struct {
	Apps []plugin_models.GetAppsModel
	Err  error
}

func (self *PsqlPlugin) connectTo(cliConnection plugin.CliConnection, command string, dbName string, psqlArgs []string) {
	appsChan := make(chan StartedAppsResult, 0)
	go func() {
		startedApps, err := self.CfService.GetStartedApps(cliConnection)
		appsChan <- StartedAppsResult{Apps: startedApps, Err: err}
	}()

	service, err := self.CfService.GetService(cliConnection, dbName)
	if err != nil {
		fmt.Fprintf(self.Err, "FAILED\nUnable to retrieve service credentials: %s\n", err)
		self.setErrorExit()
		return
	}

	appsResult := <-appsChan
	if appsResult.Err != nil {
		fmt.Fprintf(self.Err, "FAILED\nUnable to retrieve started apps: %s\n", appsResult.Err)
		self.setErrorExit()
		return
	}

	if len(appsResult.Apps) == 0 {
		fmt.Fprintf(self.Err, "FAILED\nUnable to connect to '%s': no started apps in current space\n", dbName)
		self.setErrorExit()
		return
	}

	tunnelPort := self.PortFinder.GetPort()
	self.CfService.OpenSshTunnel(cliConnection, service, appsResult.Apps, tunnelPort)

	err = self.runClient(command, "127.0.0.1", tunnelPort, service.DbName, service.Username, service.Password, psqlArgs...)
	if err != nil {
		fmt.Fprintf(self.Err, "FAILED\n%s", err)
		self.setErrorExit()
	}
}

func (self *PsqlPlugin) runClient(command string, hostname string, port int, dbName string, username string, password string, args ...string) error {
	switch command {
	case "psql":
		return self.PsqlRunner.RunPsql(hostname, port, dbName, username, password, args...)
	}

	panic(fmt.Errorf("command not implemented: %s", command))
}

type PluginConf struct {
	In         io.Reader
	Out        io.Writer
	Err        io.Writer
	CfService  CfService
	PsqlRunner PsqlRunner
	PortFinder PortFinder
}
