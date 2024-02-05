package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/samber/lo"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v3"
)

func main() {

	app := cli.NewApp()
	app.Usage = "Simple git repo clone tool with workspace support"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "workspace",
			Value: filepath.Join(lo.Must(os.UserHomeDir()), "workspace"),
			Usage: "workspace path",
		},
		cli.StringFlag{
			Name:  "config",
			Value: filepath.Join(lo.Must(os.UserHomeDir()), ".gtconfig.yaml"),
			Usage: "gt config for multiple gitconfig support",
		},
	}
	app.Action = func(c *cli.Context) error {
		gitpath := c.Args().First()
		if gitpath == "" {
			return errors.New("empty git path is invalid")
		}

		parsedPath := lo.Must(url.Parse(gitpath))
		clonePath := filepath.Join(
			lo.Must(filepath.Abs(c.String("workspace"))), parsedPath.Host, parsedPath.Path)

		out := new(config)
		lo.Must0(yaml.Unmarshal(createOfGetConfig(c.String("config")), out))

		_, err := os.Stat(clonePath)
		if os.IsNotExist(err) {
			fmt.Println("git", "clone", gitpath, clonePath)
			lo.Must0(exec.Command("git", "clone", gitpath, clonePath).Run())
		} else {
			fmt.Println(clonePath, "already exists")
		}

		if cfg, ok := lo.Find(out.GitConfig, func(c gitConfig) bool {
			return c.Host == parsedPath.Host
		}); ok {
			fmt.Println("git", "config", "--local", "user.email", cfg.Email)
			cmd := exec.Command("git", "config", "--local", "user.email", cfg.Email)
			cmd.Dir = clonePath
			lo.Must0(cmd.Run())
		}

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type gitConfig struct {
	Host  string
	Email string
}

type config struct {
	GitConfig []gitConfig
}

func createOfGetConfig(config string) []byte {
	// if config file exists, read it
	// if not, create it

	if _, err := os.Stat(config); os.IsNotExist(err) {
		// create config file
		lo.Must0(os.WriteFile(config, []byte(""), 0644))
	}

	// read config file
	data := lo.Must(os.ReadFile(config))
	if len(data) == 0 {
		return []byte("")
	}

	return data
}
