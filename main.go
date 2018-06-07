package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/jaecktec/cf-psql-plugin/cfpsql"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "This executable is a cf plugin. "+
			"Run `cf install-plugin %s` to install it\nand `cf psql service-name` "+
			"to use it.\n",
			os.Args[0])
		os.Exit(1)
	}

	psqlPlugin := newPlugin()
	plugin.Start(psqlPlugin)

	os.Exit(psqlPlugin.GetExitCode())
}

func newPlugin() *cfpsql.PsqlPlugin {
	httpClientFactory := cfpsql.NewHttpClientFactory()
	osWrapper := cfpsql.NewOsWrapper()
	requestDumper := cfpsql.NewRequestDumper(osWrapper, os.Stderr)
	http := cfpsql.NewHttpWrapper(httpClientFactory, requestDumper)
	apiClient := cfpsql.NewApiClient(http)

	sshRunner := cfpsql.NewSshRunner()
	netWrapper := cfpsql.NewNetWrapper()
	waiter := cfpsql.NewPortWaiter(netWrapper)
	randWrapper := cfpsql.NewRandWrapper()
	cfService := cfpsql.NewCfService(apiClient, sshRunner, waiter, http, randWrapper, os.Stderr)

	execWrapper := cfpsql.NewExecWrapper()
	runner := cfpsql.NewPsqlRunner(execWrapper)

	portFinder := cfpsql.NewPortFinder()

	return cfpsql.NewPsqlPlugin(cfpsql.PluginConf{
		In:         os.Stdin,
		Out:        os.Stdout,
		Err:        os.Stderr,
		CfService:  cfService,
		PortFinder: portFinder,
		PsqlRunner: runner,
	})
}
