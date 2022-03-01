package main

import (
	_deliver "BaseWebSocket/service/server"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	_ "net/http/pprof"
)

func main() {
	cmd.AddCommand(Server)

	cmd.Execute()
}

var cmd = &cobra.Command{
	Use:   "WebSocketServer",
	Short: "WebSocketServer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var Server = &cobra.Command{
	Use:   "server [file_location option]",
	Short: "Execute router server.",
	Run: func(cmd *cobra.Command, args []string) {
		server := _deliver.NewDeliverServer()
		server.Start()
	},
}
