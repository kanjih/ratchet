package handler

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"time"
)

func Init(c *cli.Context) error {
	fmt.Println("Creating migration table...")

	targetDb := getTargetDb(
		c.String(FlagNameProjectId),
		c.String(FlagNameInstanceName),
		c.String(FlagNameDatabaseName),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adminClient, dataClient, err := createClients(ctx, targetDb)
	if err != nil {
		return err
	}

	exists, err := isMigrationTableExists(ctx, dataClient)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("Migration table already exists.")
		return nil
	}

	if err = execDdl(ctx, adminClient, targetDb, migrationTableDdl); err != nil {
		return err
	}

	fmt.Println("Migration table has been created!!")
	return nil
}

func isMigrationTableExists(ctx context.Context, client *spanner.Client) (bool, error) {
	rows, err := execSql(ctx, client, "SELECT count(*) as cnt FROM information_schema.Tables WHERE table_name = 'Migrations';")
	if err != nil {
		return false, err
	}
	var count int64
	if err = rows[0].ColumnByName("cnt", &count); err != nil {
		return false, err
	}
	return count == 1, nil
}
