package ctl

import (
	"fmt"

	"github.com/leapp-to/leapp-go/pkg/msg"
	"github.com/spf13/cobra"
)

var (
	params           msg.MigrateMachine
	excludedPaths    string
	tcpPorts         string
	excludedTcpPorts string

	// migrateMachineCmd represents the migrate-machine command
	migrateMachineCmd = &cobra.Command{
		Use:   "migrate-machine",
		Short: "Executes a migration of a VM into a macrocontainer",
		Long: `This command migrates one or more application into containers by creating a macrocontainer.

This means that the entire system will be converted into a container, possibly bringing all the dirty with it.`,

		//Args: func(params *cobra.Command, args []string) error {
		////if sourceIP == "" {
		////return errors.New("source-ip is mandatory")
		////}
		//if targetIP == "" {
		//return errors.New("target-ip is mandatory")
		//}
		//return nil
		//},

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("source host: %s\n", params.SourceHost)
		},
	}
)

func init() {
	RootCmd.AddCommand(migrateMachineCmd)

	migrateMachineCmd.PersistentFlags().StringVar(&params.SourceHost, "source-host", "", "Source machine IP address")
	migrateMachineCmd.PersistentFlags().StringVar(&params.TargetHost, "target-host", "", "Target machine IP address")
	migrateMachineCmd.PersistentFlags().StringVar(&params.ContainerName, "container-name", "", "Container name")
	migrateMachineCmd.PersistentFlags().StringVar(&params.SourceUser, "source-user", "", "User in the source host")
	migrateMachineCmd.PersistentFlags().StringVar(&params.TargetUser, "target-user", "", "User in the target host")
	migrateMachineCmd.PersistentFlags().StringVar(&excludedPaths, "excluded-paths", "", "...")
	migrateMachineCmd.PersistentFlags().StringVar(&tcpPorts, "tcp-ports", "", "Mapping of TCP ports, ex: 80:8080")
	migrateMachineCmd.PersistentFlags().StringVar(&excludedTcpPorts, "excluded-tcp-ports", "", "TCP ports that shouldn't be used in the target")
	migrateMachineCmd.PersistentFlags().BoolVar(&params.ForceCreate, "force-create", false, "Force container creation")
	migrateMachineCmd.PersistentFlags().BoolVar(&params.DisableStart, "disable-start", true, "Don't start the container after creating it")
	migrateMachineCmd.PersistentFlags().BoolVar(&params.Debug, "debug", false, "Enable debug")
}
