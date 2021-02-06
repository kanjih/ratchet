package handler

import (
	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

const migrationSelectQuery = "SELECT id FROM `Migrations`"

func getTargetDb(projectId, instanceName, dbName string) string {
	return fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectId, instanceName, dbName)
}

func createClients(ctx context.Context, databaseId string) (*database.DatabaseAdminClient, *spanner.Client, error) {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	dataClient, err := spanner.NewClient(ctx, databaseId)
	if err != nil {
		return nil, nil, err
	}
	return adminClient, dataClient, nil
}

func execSql(ctx context.Context, client *spanner.Client, sql string) ([]*spanner.Row, error) {
	stmt := spanner.Statement{SQL: sql}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var rows []*spanner.Row
	for {
		row, err := iter.Next()
		switch err {
		case nil:
			rows = append(rows, row)
		case iterator.Done:
			return rows, nil
		default:
			return nil, err
		}
	}
}

func execDdl(ctx context.Context, adminClient *database.DatabaseAdminClient, targetDb string, query string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database:   targetDb,
		Statements: []string{query},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func execDml(ctx context.Context, client *spanner.Client, dml string) error {
	return execDmls(ctx, client, []string{dml})
}

func execDmls(ctx context.Context, client *spanner.Client, dmls []string) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		for _, dml := range dmls {
			stmt := spanner.Statement{SQL: dml}
			_, err := txn.Update(ctx, stmt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func execPartitionedDml(ctx context.Context, client *spanner.Client, dml string) error {
	_, err := client.PartitionedUpdate(ctx, spanner.Statement{SQL: dml})
	return err
}
