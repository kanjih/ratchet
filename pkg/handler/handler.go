package handler

import (
	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"context"
	"ratchet/handler"
)

func Init(ctx context.Context, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client, targetDb string) error {
	return handler.ExecInit(ctx, adminClient, dataClient, targetDb)
}

func Run(ctx context.Context, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client, targetDb string, migrationFilePaths []string) error {
	return handler.ExecInit(ctx, adminClient, dataClient, targetDb)
}
