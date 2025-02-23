package resources

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/MaterializeInc/terraform-provider-materialize/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestResourceGrantConnectionDefaultPrivilegeCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"grantee_name":     "project_managers",
		"target_role_name": "developers",
		"privilege":        "USAGE",
	}
	d := schema.TestResourceDataRaw(t, GrantConnectionDefaultPrivilege().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		// Create
		mock.ExpectExec(
			`ALTER DEFAULT PRIVILEGES FOR ROLE "developers" GRANT USAGE ON CONNECTIONS TO "project_managers";`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// DefaultPrivilegeId - Query grantee role
		gp := `WHERE mz_roles.name = 'project_managers'`
		testhelpers.MockRoleScan(mock, gp)

		// DefaultPrivilegeId - Query target role
		tp := `WHERE mz_roles.name = 'developers'`
		testhelpers.MockRoleScan(mock, tp)

		// Query Params
		qp := `
			WHERE mz_default_privileges.grantee = 'u1'
			AND mz_default_privileges.object_type = 'connection'
			AND mz_default_privileges.role_id = 'u1'`
		testhelpers.MockDefaultPrivilegeScan(mock, qp, "connection")

		if err := grantConnectionDefaultPrivilegeCreate(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}

		if d.Id() != "GRANT DEFAULT|CONNECTION|u1|u1|||USAGE" {
			t.Fatalf("unexpected id of %s", d.Id())
		}
	})
}

func TestResourceGrantConnectionDefaultPrivilegeDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"grantee_name":     "project_managers",
		"target_role_name": "developers",
		"privilege":        "USAGE",
	}
	d := schema.TestResourceDataRaw(t, GrantConnectionDefaultPrivilege().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`ALTER DEFAULT PRIVILEGES FOR ROLE "developers" REVOKE USAGE ON CONNECTIONS FROM "project_managers";`).WillReturnResult(sqlmock.NewResult(1, 1))

		if err := grantConnectionDefaultPrivilegeDelete(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}
	})
}
