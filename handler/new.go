package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

func New(c *cli.Context) error {
	err := os.Mkdir(migrationFolderName, 0777)
	if err != nil && !os.IsExist(err) {
		return err
	}
	fileName := makeFileName(c)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	fmt.Println("New migration file has been created in " + fileName)
	return file.Close()
}

func makeFileName(c *cli.Context) string {
	datePrefix := time.Now().Format("2006-01-02_15-04-05")
	idPrefix := uuid.New().String()[:5] // Full uuid is not necessary, 5 letters are enough to avoid the file name conflict
	dmlFlag := ""
	if c.IsSet(FlagNameDml) {
		dmlFlag = Dml + "_"
	} else if c.IsSet(FlagNamePartitionedDml) {
		dmlFlag = PartitionedDml + "_"
	}
	return migrationFolderName + "/" + datePrefix + "_" + dmlFlag + idPrefix + ".sql"
}
