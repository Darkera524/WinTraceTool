package g

import (
	"github.com/jackc/pgx"
)

func initConfig() (config *pgx.ConnConfig){
	config = &pgx.ConnConfig{
		Host: "10.201.50.82",
		Port: 5432,
		User: "postgres",
		Database: "tutorial",
		Password: "",
	}
	return config
}

func ExecSql(sql string) error {
	config := initConfig()
	conn, err := pgx.Connect(*config)
	if err != nil {
		Logger().Println("conn err")
		return err
	}
	tx, _ := conn.Begin()
	defer  tx.Rollback()

	_,err = tx.Exec(sql)
	if err != nil {
		Logger().Println(err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		Logger().Println(err.Error())
		return err
	}

	return nil
}
