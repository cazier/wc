package cmd

import (
	"github.com/cazier/wc/db"
	"github.com/spf13/cobra"
)

var databasePath string
var importTeamPath string
var importMatchPath string

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage the backend database",
}

var initializeCmd = &cobra.Command{
	Use:     "initialize",
	Aliases: []string{"init"},
	Short:   "Initialize an empty database for use",
	Long: `Create a database with all of its tables. Optionally, you can supply
a set of import flags to fill the database with values`,
	Run: func(cmd *cobra.Command, args []string) {
		db.Init(databasePath)
		db.LinkTables()

		if importTeamPath != "" || importMatchPath != "" {
			db.LoadYaml(importTeamPath, importMatchPath)
		}
	},
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import details from a yaml file into the database",
	Run: func(cmd *cobra.Command, args []string) {
		db.Init(databasePath)
		db.LoadYaml(importTeamPath, importMatchPath)
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(initializeCmd)
	databaseCmd.AddCommand(importCmd)

	databaseCmd.PersistentFlags().StringVar(&databasePath, "db", ".", "path to a database")

	for _, cmd := range []*cobra.Command{initializeCmd, importCmd} {
		cmd.Flags().StringVar(&importTeamPath, "teams", "", "team yaml file for importing")
		cmd.Flags().StringVar(&importMatchPath, "matches", "", "match yaml file for importing")

		// if cmd == importCmd {
		// 	// cmd.MarkFlagRequired("teams")
		// 	// cmd.MarkFlagRequired("matches")
		// }
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// databaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// databaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
