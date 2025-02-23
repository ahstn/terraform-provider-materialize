package materialize

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DatabaseBuilder struct {
	ddl          Builder
	databaseName string
}

func NewDatabaseBuilder(conn *sqlx.DB, obj MaterializeObject) *DatabaseBuilder {
	return &DatabaseBuilder{
		ddl:          Builder{conn, Database},
		databaseName: obj.Name,
	}
}

func (b *DatabaseBuilder) QualifiedName() string {
	return QualifiedName(b.databaseName)
}

func (b *DatabaseBuilder) Create() error {
	q := fmt.Sprintf(`CREATE DATABASE %s;`, b.QualifiedName())
	return b.ddl.exec(q)
}

func (b *DatabaseBuilder) Drop() error {
	qn := b.QualifiedName()
	return b.ddl.drop(qn)
}

type DatabaseParams struct {
	DatabaseId   sql.NullString `db:"id"`
	DatabaseName sql.NullString `db:"database_name"`
	OwnerName    sql.NullString `db:"owner_name"`
	Privileges   sql.NullString `db:"privileges"`
}

var databaseQuery = NewBaseQuery(`
	SELECT
		mz_databases.id,
		mz_databases.name AS database_name,
		mz_roles.name AS owner_name,
		mz_databases.privileges
	FROM mz_databases
	JOIN mz_roles
		ON mz_databases.owner_id = mz_roles.id`)

func DatabaseId(conn *sqlx.DB, obj MaterializeObject) (string, error) {
	q := databaseQuery.QueryPredicate(map[string]string{"mz_databases.name": obj.Name})

	var c DatabaseParams
	if err := conn.Get(&c, q); err != nil {
		return "", err
	}

	return c.DatabaseId.String, nil
}

func ScanDatabase(conn *sqlx.DB, id string) (DatabaseParams, error) {
	q := databaseQuery.QueryPredicate(map[string]string{"mz_databases.id": id})

	var c DatabaseParams
	if err := conn.Get(&c, q); err != nil {
		return c, err
	}

	return c, nil
}

func ListDatabases(conn *sqlx.DB) ([]DatabaseParams, error) {
	q := databaseQuery.QueryPredicate(map[string]string{})

	var c []DatabaseParams
	if err := conn.Select(&c, q); err != nil {
		return c, err
	}

	return c, nil
}
