package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"syscall"
	"vulcanlabs-assignment/pkg/client"
	"vulcanlabs-assignment/pkg/db"
	"vulcanlabs-assignment/pkg/server"
	"vulcanlabs-assignment/pkg/utils"
)

var rootCmd = &cobra.Command{
	Use:   "cine",
	Short: "cine command line interface",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "cine server command",
	Run: func(cmd *cobra.Command, args []string) {
		err := db.InitDB()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		go server.StartServer(port, rows, cols, distance, reset)
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
	},
}
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "cine client command",
	Run: func(cmd *cobra.Command, args []string) {
		logger := utils.Logger("cine-client")
		logger.Info("cine client command")
	},
}
var getSeatMapCmd = &cobra.Command{
	Use:   "get-seat-map",
	Short: "cine get seat map command",
	Run: func(cmd *cobra.Command, args []string) {
		cine, err := client.NewCineClient(port)
		if err != nil {
			os.Exit(1)
		}
		defer cine.Close()
		_, err = cine.GetSeatMap()
		if err != nil {
			os.Exit(1)
		}
	},
}
var reserveCmd = &cobra.Command{
	Use:   "reserve x y [flags]",
	Short: "cine client reserve x y [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args)%2 != 0 {
			fmt.Println("Usage: cine reserve X Y [flags] ")
			os.Exit(1)
		}
		cine, err := client.NewCineClient(port)
		if err != nil {
			os.Exit(1)
		}
		defer cine.Close()
		parsed := slices.Chunk(args, 2)
		seatList := make([][2]int32, 0)
		for p := range parsed {
			x, _ := strconv.Atoi(p[0])
			y, _ := strconv.Atoi(p[1])
			seatList = append(seatList, [2]int32{int32(x), int32(y)})
		}
		cine.Reserve(seatList)

	},
}
var (
	port, rows, cols, distance int
	reset                      bool
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "server port")
	rootCmd.PersistentFlags().BoolVarP(&reset, "reset", "s", false, "reset cine app")
	serverCmd.PersistentFlags().IntVarP(&rows, "rows", "r", 10, "number of rows")
	serverCmd.PersistentFlags().IntVarP(&cols, "cols", "c", 10, "number of columns")
	serverCmd.PersistentFlags().IntVarP(&distance, "distance", "d", 1, "distance to server")
	// client
	clientCmd.AddCommand(getSeatMapCmd)
	clientCmd.AddCommand(reserveCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(clientCmd)
}

func Execute() {
	rootCmd.Execute()
}
