package handler

import (
	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

const (
	migrationFolderName = "migrations"
	migrationTableDdl   = `
CREATE TABLE Migrations (
	id STRING(MAX) NOT NULL,
	executed_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (id)`
	migrationInsertBaseQuery = "INSERT INTO `Migrations`(`id`, `executed_at`) VALUES ('%s', PENDING_COMMIT_TIMESTAMP())"
)

func Run(c *cli.Context) error {
	fmt.Println("Migration started.")

	targetDb := getTargetDb(c.String(FlagNameProjectId), c.String(FlagNameInstanceName), c.String(FlagNameDatabaseName))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	adminClient, dataClient, err := createClients(ctx, targetDb)
	if err != nil {
		return err
	}

	executedMigrationIds, err := fetchExecutedMigrationIds(ctx, dataClient)
	if err != nil {
		return err
	}

	migrationFilePaths, err := getMigrationFilePaths(migrationFolderName)
	if err != nil {
		return err
	}
	for _, migrationFilePath := range migrationFilePaths {
		if err := runEachMigration(ctx, adminClient, dataClient, migrationFilePath, targetDb, executedMigrationIds); err != nil {
			return err
		}
	}

	fmt.Println("Migration completed!")
	return nil
}

func runEachMigration(ctx context.Context, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client, migrationFilePath, targetDb string, executedMigrationIds map[string]struct{}) error {
	migrationId := makeMigrationIdFromFilePath(migrationFilePath)
	if _, exists := executedMigrationIds[migrationId]; exists {
		return nil
	}
	fmt.Print("running " + migrationId + " ... ")

	fileContent, err := readFile(migrationFilePath)
	if err != nil {
		return err
	}

	for _, query := range parseToQueries(fileContent) {
		switch {
		case strings.Contains(migrationFilePath, PartitionedDml):
			err = execPartitionedDml(ctx, dataClient, query)
		case strings.Contains(migrationFilePath, Dml):
			err = execDml(ctx, dataClient, query)
		default:
			err = execDdl(ctx, adminClient, targetDb, query)
		}

		if err != nil {
			fmt.Println("following query failed.\n" + query + ";")
			return err
		}
	}

	// insert into Migration table to mark the migration is executed
	if err = execDml(ctx, dataClient, fmt.Sprintf(migrationInsertBaseQuery, migrationId)); err != nil {
		return err
	}
	fmt.Println("done.")

	return nil
}

func parseToQueries(fileContent string) []string {
	var queries []string

	fileContent = trimComments(fileContent)
	for _, query := range strings.Split(fileContent, ";") {
		query = strings.TrimLeft(query, "\n")
		if query != "" {
			queries = append(queries, query)
		}
	}

	return queries
}

func getMigrationFilePaths(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, nil
}

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func trimComments(fileContent string) string {
	var res []string
	for _, line := range strings.Split(fileContent, "\n") {
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "--") {
			res = append(res, line)
		}
	}
	return strings.Join(res, "\n")
}

func makeMigrationIdFromFilePath(filePath string) string {
	fileName := strings.Split(filePath, "/")[1] // remove folder name
	return strings.Split(fileName, ".")[0]      // remove ext
}

func fetchExecutedMigrationIds(ctx context.Context, client *spanner.Client) (map[string]struct{}, error) {
	rows, err := execSql(ctx, client, migrationSelectQuery)
	if err != nil {
		return nil, err
	}

	idSet := map[string]struct{}{}
	for _, row := range rows {
		var migrationId string
		if err = row.ColumnByName("id", &migrationId); err != nil {
			return nil, err
		} else {
			idSet[migrationId] = struct{}{}
		}
	}
	return idSet, nil
}
