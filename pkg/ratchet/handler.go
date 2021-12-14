package ratchet

import (
	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"context"
	"github.com/kanjih/ratchet/handler"
)

type Migration struct {
	Id      string
	Content []byte
}

func Init(ctx context.Context, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client, targetDb string) error {
	return handler.ExecInit(ctx, adminClient, dataClient, targetDb)
}

func Run(ctx context.Context, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client, targetDb string, migrations []Migration) error {
	var parsedMigrations []handler.Migrations
	for _, m := range migrations {
		parsedMigrations = append(parsedMigrations, handler.Migrations{Id: m.Id, Content: m.Content})
	}
	return handler.ExecRun(ctx, adminClient, dataClient, targetDb, parsedMigrations)
}

func RunWithFiles(ctx context.Context, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client, targetDb string, migrationFilePaths []string) error {
	return handler.ExecRunWithFiles(ctx, adminClient, dataClient, targetDb, migrationFilePaths)
}
