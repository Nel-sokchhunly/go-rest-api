package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Nel-sokchhunly/go-rest-api/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	dialect     = "pgx"
	fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory with migration files")
)

func main() {
	flags.Usage = usage      // func to print usage info
	flags.Parse(os.Args[1:]) // parse command line args

	// print usage info if no command is provided or help flag is set
	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	c := config.NewDB()
	dbString := fmt.Sprintf(fmtDBString, c.Host, c.User, c.Password, c.DBName, c.Port)
	db, err := goose.OpenDBWithDriver(dialect, dbString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err.Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err.Error())
		}
	}()

	if err := goose.Run(command, db, *dir, args[1:]...); err != nil {
		log.Fatalf("migrate %v: %v", command, err)
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
