package main

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/johndsheehan/n26balance/pkg/n26"
	"github.com/spf13/cobra"

	yaml "gopkg.in/yaml.v2"
)

func n26CmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use: "n26balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgFile, err := cmd.Flags().GetString("config")
			if err != nil {
				return errors.New("failed to determine configuration file")
			}

			if cfgFile == "" {
				return errors.New("no configuration file provided")
			}

			return n26CmdExecute(cfgFile)
		},
	}

	cmd.Flags().StringP("config", "c", "", "configuration file")

	return cmd
}

func n26CmdExecute(cfgFile string) error {
	yml, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	var cfg n26.Config
	err = yaml.Unmarshal(yml, &cfg)
	if err != nil {
		return err
	}

	n26Client, err := n26.NewClient(cfg)
	if err != nil {
		return err
	}

	balance, err := n26Client.Balance()
	if err != nil {
		return err
	}

	log.Printf("available: %0.2f\n", balance.AvailableBalance)
	log.Printf("usable: %0.2f\n", balance.UsableBalance)

	return nil
}

func main() {
	cmd := n26CmdCreate()

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
