package resources

import (
	"context"
	"testing"

	"github.com/MaterializeInc/terraform-provider-materialize/pkg/testhelpers"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var inPostgres = map[string]interface{}{
	"name":                      "conn",
	"schema_name":               "schema",
	"database_name":             "database",
	"database":                  "default",
	"host":                      "postgres_host",
	"port":                      5432,
	"user":                      []interface{}{map[string]interface{}{"secret": []interface{}{map[string]interface{}{"name": "user"}}}},
	"password":                  []interface{}{map[string]interface{}{"name": "password"}},
	"ssh_tunnel":                []interface{}{map[string]interface{}{"name": "ssh_conn"}},
	"ssl_certificate_authority": []interface{}{map[string]interface{}{"secret": []interface{}{map[string]interface{}{"name": "root"}}}},
	"ssl_certificate":           []interface{}{map[string]interface{}{"secret": []interface{}{map[string]interface{}{"name": "cert"}}}},
	"ssl_key":                   []interface{}{map[string]interface{}{"name": "key"}},
	"ssl_mode":                  "verify-full",
	"aws_privatelink":           []interface{}{map[string]interface{}{"name": "link"}},
}

func TestResourceConnectionPostgresCreate(t *testing.T) {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, ConnectionPostgres().Schema, inPostgres)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		// Create
		mock.ExpectExec(
			`CREATE CONNECTION "database"."schema"."conn" TO POSTGRES \(HOST 'postgres_host', PORT 5432, USER SECRET "database"."schema"."user", PASSWORD SECRET "database"."schema"."password", SSL MODE 'verify-full', SSH TUNNEL "database"."schema"."ssh_conn", SSL CERTIFICATE AUTHORITY SECRET "database"."schema"."root", SSL CERTIFICATE SECRET "database"."schema"."cert", SSL KEY SECRET "database"."schema"."key", AWS PRIVATELINK "database"."schema"."link", DATABASE 'default'\);`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Query Id
		ip := `WHERE mz_connections.name = 'conn' AND mz_databases.name = 'database' AND mz_schemas.name = 'schema'`
		testhelpers.MockConnectionScan(mock, ip)

		// Query Params
		pp := `WHERE mz_connections.id = 'u1'`
		testhelpers.MockConnectionScan(mock, pp)

		if err := connectionPostgresCreate(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}
	})
}
