package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cazier/wc/db"
	"github.com/cazier/wc/db/load"
)

var sqlite bool
var databasePath string

var importTeamPath string
var importMatchPath string
var importPlayerPath string

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
		databaseInit(true)
		importCmd.Run(cmd, args)
	},
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import details from a yaml file into the database",
	Run: func(cmd *cobra.Command, args []string) {
		databaseInit(true)

		if importTeamPath != "" {
			load.Teams(importTeamPath)
		}

		if importMatchPath != "" {
			load.Matches(importMatchPath)
		}

		if importPlayerPath != "" {
			load.Players(importPlayerPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(initializeCmd)
	databaseCmd.AddCommand(importCmd)

	databaseCommand(databaseCmd)

	for _, cmd := range []*cobra.Command{initializeCmd, importCmd} {
		cmd.Flags().StringVar(&importTeamPath, "teams", "", "team yaml file for importing")
		cmd.Flags().StringVar(&importMatchPath, "matches", "", "match yaml file for importing")
		cmd.Flags().StringVar(&importPlayerPath, "players", "", "player yaml file for importing")

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

func databaseCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&databasePath, "db", ".", "path to a database")
	cmd.Flags().BoolVar(&sqlite, "sqlite", true, "use a sqlite database instead of mariadb")
}

func databaseInit(purge bool) {
	if sqlite && databasePath == "memory" {
		db.InitSqlite(&db.SqliteDBOptions{Memory: true})
	} else if sqlite {
		info, _ := os.Stat(databasePath)
		if info.IsDir() {
			databasePath = fmt.Sprintf("%s/%s", databasePath, "wc.db")
		}
		db.InitSqlite(&db.SqliteDBOptions{Path: databasePath, LogLevel: 4, LogPath: "stdout"})

	} else {
		db.InitMariaDB(&db.MariaDBOptions{})
	}
	db.LinkTables(purge)
}
