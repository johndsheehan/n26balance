package main

import (
	"errors"
	"io/ioutil"
	"log"
	"time"

	"github.com/johndsheehan/n26balance/pkg/n26"
	"github.com/spf13/cobra"

	yaml "gopkg.in/yaml.v2"
)

const queryInterval = 600 * time.Second

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

			watch, err := cmd.Flags().GetBool("watch")
			if err != nil {
				return errors.New("failed to determine watch value")
			}

			interval := queryInterval
			if watch {
				i, err := cmd.Flags().GetInt("interval")
				if err != nil {
					return errors.New("failed to determine interval value")
				}

				if i < 1 {
					log.Printf("integer should be positive interval, using default: %d seconds",
						queryInterval/(1000*time.Millisecond))
				} else {
					interval = time.Duration(i) * time.Second
				}
			}

			return n26CmdExecute(cfgFile, watch, interval)
		},
	}

	cmd.Flags().StringP("config", "c", "", "configuration file")
	cmd.Flags().IntP("interval", "i", 0, "interval between queries")
	cmd.Flags().BoolP("watch", "w", false, "query balance every \"interval\" minutes, report if changed")

	return cmd
}

func n26CmdExecute(cfgFile string, watch bool, interval time.Duration) error {
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

	if watch {
		previous := balance.UsableBalance

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case _ = <-ticker.C:
				balance, err = n26Client.Balance()
				if err != nil {
					return err
				}

				if balance.UsableBalance != previous {
					log.Printf("available: %0.2f\n", balance.AvailableBalance)
					log.Printf("usable: %0.2f\n", balance.UsableBalance)
					previous = balance.UsableBalance
				}
			}
		}
	}

	return nil
}

func main() {
	cmd := n26CmdCreate()

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
