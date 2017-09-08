package db

type MySqlDb interface {
	// init db connection
	InitDB()

	// get db version for migrate
	CheckMigrate()

	// create tables  or migrate old table
	migrate(oldVersion int32)
}

type DBVersion struct {
	Version int32
}
