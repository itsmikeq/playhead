package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"playhead/api"
	"playhead/app"
)

func init() {
	rootCmd.AddCommand(routesCmd)
}

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Print the routes",
	Run: func(cmd *cobra.Command, args []string) {
		a, err := app.New()
		if err != nil {
			panic(err)
		}
		defer a.Close()

		api, err := api.New(a)
		if err != nil {
			panic(err)
		}
		router := mux.NewRouter()
		api.Init(router.PathPrefix("/api").Subrouter())
		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			t, err := route.GetPathTemplate()
			if err != nil {
				return err
			}
			fmt.Println(t)
			return nil
		})

	},
}
