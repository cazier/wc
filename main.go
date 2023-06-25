package main

import "github.com/cazier/wc/cmd"

func main() {
	// viper.SetConfigName("config.toml")
	// viper.AddConfigPath("$HOME/.config/wc")
	// viper.AddConfigPath(".")
	// err := viper.ReadInConfig
	// if err != nil {
	// 	panic(fmt.Errorf("Could not open the configuration file: %w", err))
	// }
	cmd.Execute()
}
