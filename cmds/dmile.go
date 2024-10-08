package cmds

type DmileCommand struct {
	CreateData    CreateDataCommand    `cmd:"" name:"create-data" help:"create new data"`
	RegisterModel RegisterModelCommand `cmd:"" name:"register-model" help:"register did model"`
	MigrateData   MigrateDataCommand   `cmd:"" name:"migrate-data" help:"migrate data"`
}
