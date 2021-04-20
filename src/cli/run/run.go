package run

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vj396/milton/lib/file"
	"github.com/vj396/milton/lib/logger"
	"github.com/vj396/milton/lib/utils"
	"github.com/vj396/milton/pkg/milton"
	"github.com/vj396/milton/src/cli/root"
	"github.com/vj396/milton/src/types"
)

type runCmd struct {
	configFile      string
	modelsDirectory string
	debug           bool

	cmd *cobra.Command
}

func (c *runCmd) run(cmd *cobra.Command, args []string) {
	done := make(chan struct{})
	go utils.TrapOSInterrupt(done)
	log, err := logger.BuildLogger(c.debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not build logger, err: %+v", err)
		os.Exit(1)
	}
	conf := new(types.Config)
	err = file.ReadYAMLFile(c.configFile, conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading custom config: %s\n", err.Error())
		os.Exit(1)
	}
	milton.Start(done, log, conf, c.modelsDirectory)
}

func init() {
	run := &runCmd{}
	run.cmd = &cobra.Command{
		Use:   "run",
		Short: "Start generating files",
		Long:  ``,
		Run:   run.run,
	}
	run.cmd.Flags().StringVarP(&run.configFile, "config", "c", "", "Path to config yaml file")
	run.cmd.Flags().StringVarP(&run.modelsDirectory, "models-dir", "d", "", "Path to models directory")
	run.cmd.Flags().BoolVarP(&run.debug, "debug", "v", false, "Enable debug logging")
	run.cmd.MarkFlagRequired("config")
	run.cmd.MarkFlagRequired("models-dir")
	root.GetRoot().AddCommand(run.cmd)
}
