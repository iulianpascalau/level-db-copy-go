package main

import (
	"os"

	"iulianpascalau/level-db-copy-go/process"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/urfave/cli"
)

var (
	logLevel = cli.StringFlag{
		Name: "log-level",
		Usage: "This flag specifies the logger `level(s)`. It can contain multiple comma-separated value. For example" +
			", if set to *:INFO the logs for all packages will have the INFO level. However, if set to *:INFO,api:DEBUG" +
			" the logs for all packages will have the INFO level, excepting the api package which will receive a DEBUG" +
			" log level.",
		Value: "*:" + logger.LogInfo.String(),
	}
	logSaveFile = cli.BoolFlag{
		Name:  "log-save",
		Usage: "Boolean option for enabling log saving. If set, it will automatically save all the logs into a file.",
	}
	sourceDir = cli.StringFlag{
		Name:  "source",
		Usage: "The source directory to read data from",
		Value: "source",
	}
	destinationDir = cli.StringFlag{
		Name:  "destination",
		Usage: "The destination directory to write the missing data to",
		Value: "destination",
	}

	log          = logger.GetOrCreate("tool")
	helpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`
)

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = helpTemplate
	app.Name = "Level DB copy missing data tool"
	app.Usage = ""
	app.Flags = []cli.Flag{
		logLevel,
		logSaveFile,
		sourceDir,
		destinationDir,
	}

	app.Authors = []cli.Author{
		{
			Name:  "Iulian Pascalau",
			Email: "",
		},
	}

	app.Action = copyProcess

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func copyProcess(ctx *cli.Context) error {
	log.Info("Level DB copy missing data tool. Copying data",
		"from", ctx.GlobalString(sourceDir.Name),
		"to", ctx.GlobalString(destinationDir.Name))

	dirHandler, err := process.NewDirectoriesHandler(
		ctx.GlobalString(sourceDir.Name),
		ctx.GlobalString(destinationDir.Name),
	)
	if err != nil {
		return err
	}

	dbCopyHandler, err := process.NewDataCopyHandler(
		dirHandler,
		process.NewDBWrapper(),
		process.NewDBWrapper(),
	)
	if err != nil {
		return err
	}

	return dbCopyHandler.Process()
}
