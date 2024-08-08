package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strings"
	"testing"
	"time"
)

type testCaseGetOrganizations struct {
	name       string
	filter     *model.OrganizationFilters
	mutations  func(t *testing.T, db *sqlx.DB)
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error)
}

func getTestCasesGetOrganizations() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name: "success-filter-ids-in",
			filter: &model.OrganizationFilters{
				IdsIn: []int{1},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
				fmt.Printf("Organizations returned: %v\n", paginated.Organizations)

				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 1, "unexpected number of organizations")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected row count")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Demby",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entry := mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}

				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
		},
		//{
		//	name: "success-filter-names-in",
		//	filter: &model.OrganizationFilters{
		//		OrganizationNameIn: []string{"Lawrence", "Younes"},
		//	},
		//	mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) ([]mysqlmodel.Organization, []mysqlmodel.User) {
		//		entryUser := mysqlmodel.User{
		//			Firstname: "Demby",
		//			Lastname:  "Abella",
		//			Email:     "demby@test.com",
		//			Password:  "password",
		//		}
		//
		//		err := entryUser.Insert(context.Background(), db, boil.Infer())
		//		require.NoError(t, err, "error inserting in the user db")
		//
		//		foundUser, err := mysqlmodel.Users(
		//			mysqlmodel.UserWhere.Firstname.EQ(organization.CreatedBy),
		//		).One(context.TODO(), db)
		//		require.NoError(t, err, "error converting CreatedBy to ID")
		//
		//		entryOrg1 := mysqlmodel.Organization{
		//			ID:            1,
		//			Name:          "Lawrence",
		//			CreatedBy:     null.IntFrom(foundUser.ID),
		//			LastUpdatedBy: null.IntFrom(foundUser.ID),
		//			CreatedAt:     organization.CreatedAt,
		//			LastUpdatedAt: organization.LastUpdatedAt,
		//			IsActive:      organization.IsActive,
		//		}
		//
		//		entryOrg2 := mysqlmodel.Organization{
		//			ID:            2,
		//			Name:          "Younes",
		//			CreatedBy:     null.IntFrom(foundUser.ID),
		//			LastUpdatedBy: null.IntFrom(foundUser.ID),
		//			CreatedAt:     organization.CreatedAt,
		//			LastUpdatedAt: organization.LastUpdatedAt,
		//			IsActive:      organization.IsActive,
		//		}
		//
		//		err = entryOrg1.Insert(context.Background(), db, boil.Infer())
		//		require.NoError(t, err, "error inserting first organization")
		//
		//		err = entryOrg2.Insert(context.Background(), db, boil.Infer())
		//		require.NoError(t, err, "error inserting second organization")
		//
		//		return []mysqlmodel.Organization{entryOrg1, entryOrg2}, []mysqlmodel.User{entryUser}
		//	},
		//	assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
		//		require.NoError(t, err, "unexpected error")
		//		require.NotNil(t, paginated, "unexpected nil paginated")
		//		require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
		//		require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
		//		assert.True(t, len(paginated.Organizations) == 2, "unexpected number of organizations")
		//		assert.True(t, paginated.Pagination.RowCount == 2, "unexpected row count")
		//
		//		modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
		//	},
		//},
	}
}

func Test_GetOrganizations(t *testing.T) {
	for _, testCase := range getTestCasesGetOrganizations() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			testCase.mutations(t, db)
			paginated, err := m.GetOrganizations(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}

//type testCaseDeleteOrganization struct {
//	name       string
//	id         int
//	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
//	mutations  func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int)
//}
//
//func getDeleteOrganizationTestCases() []testCaseDeleteOrganization {
//	return []testCaseDeleteOrganization{
//		{
//			name: "success",
//			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
//				require.Nil(t, err, "unexpected non-nil error")
//				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
//				require.NoError(t, err, "unexpected error fetching the organization")
//
//				assert.Equal(t, false, entry.IsActive)
//			},
//			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int) {
//				entry := mysqlmodel.Organization{
//					ID:        organization.Id,
//					Name:      organization.Name,
//					CreatedBy: organization.CreatedBy,
//					IsActive:  organization.IsActive,
//				}
//				err := entry.Insert(context.TODO(), db, boil.Infer())
//				assert.NoError(t, err, "unexpected insert error")
//				return organization.Id
//			},
//		},
//	}
//}
//
//func Test_DeleteOrganization(t *testing.T) {
//	for _, testCase := range getDeleteOrganizationTestCases() {
//		t.Run(testCase.name, func(t *testing.T) {
//			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
//			defer cleanup()
//			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
//			require.NotNil(t, testCase.assertions, "unexpected nil assertions")
//
//			cfg := &Config{
//				Logger:        testLogger,
//				QueryTimeouts: testQueryTimeouts,
//			}
//
//			m, err := New(cfg)
//			require.NoError(t, err, "unexpected error")
//			require.NotNil(t, m, "unexpected nil")
//
//			txHandler, err := mysqltx.New(&mysqltx.Config{
//				Logger:       testLogger,
//				Db:           db,
//				DatabaseName: cp.Database,
//			})
//			require.NoError(t, err, "unexpected error creating the tx handler")
//
//			txHandlerDb, err := txHandler.Db(testCtx)
//			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
//			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")
//
//			organization := &model.Organization{
//				Id:        1,
//				Name:      "Organization A",
//				CreatedBy: "Super Admin",
//				IsActive:  true,
//			}
//
//			id := testCase.mutations(t, db, organization)
//
//			err = m.DeleteOrganization(testCtx, txHandlerDb, id)
//			testCase.assertions(t, db, id, err)
//		})
//	}
//}

type testCaseUpdateOrganization struct {
	name                string
	initialOrganization mysqlmodel.Organization
	updateOrganization  model.UpdateOrganization
	assertions          func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations           func(t *testing.T, db *sqlx.DB, organization *model.UpdateOrganization) (updateData model.UpdateOrganization)
}

func getUpdateOrganizationTestCases() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
		{
			name: "success",
			initialOrganization: mysqlmodel.Organization{
				ID:            1,
				Name:          "Organization A",
				CreatedBy:     null.IntFrom(1),
				LastUpdatedBy: null.IntFrom(1),
				CreatedAt:     time.Now(),
				LastUpdatedAt: null.TimeFrom(time.Now()),
				IsActive:      true,
			},
			updateOrganization: model.UpdateOrganization{
				Id:   1,
				Name: null.StringFrom("Organization A new"),
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.UpdateOrganization) (updateData model.UpdateOrganization) {
				updateData = model.UpdateOrganization{
					Id:   organization.Id,
					Name: organization.Name,
				}
				return updateData
			},
		},
		{
			name: "fail",
			initialOrganization: mysqlmodel.Organization{
				ID:            1,
				Name:          "Organization A",
				CreatedBy:     null.IntFrom(1),
				LastUpdatedBy: null.IntFrom(1),
				CreatedAt:     time.Now(),
				LastUpdatedAt: null.TimeFrom(time.Now()),
				IsActive:      true,
			},
			updateOrganization: model.UpdateOrganization{
				Id:   32123,
				Name: null.StringFrom("Wrong Name"),
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "expected error when updating a non-existent organization")
				assert.Contains(t, err.Error(), "expected exactly one entry for entity: 32123")
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.UpdateOrganization) (updateData model.UpdateOrganization) {
				updateData = model.UpdateOrganization{
					Id:   organization.Id,
					Name: organization.Name,
				}
				return updateData
			},
		},
	}
}

func Test_UpdateOrganization(t *testing.T) {
	for _, testCase := range getUpdateOrganizationTestCases() {
		db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
		defer cleanup()

		cfg := &Config{
			Logger:        testLogger,
			QueryTimeouts: testQueryTimeouts,
		}

		m, err := New(cfg)
		require.NoError(t, err, "unexpected error")
		require.NotNil(t, m, "unexpected nil")

		txHandler, err := mysqltx.New(&mysqltx.Config{
			Logger:       testLogger,
			Db:           db,
			DatabaseName: cp.Database,
		})
		require.NoError(t, err, "unexpected error creating the tx handler")

		txHandlerDb, err := txHandler.Db(testCtx)
		require.NoError(t, err, "unexpected error fetching the db from the tx handler")
		require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

		err = testCase.initialOrganization.Insert(context.Background(), db, boil.Infer())
		require.NoError(t, err, "unexpected error inserting initial organization")

		updateData := testCase.mutations(t, db, &testCase.updateOrganization)

		_, err = m.UpdateOrganization(testCtx, txHandlerDb, &updateData)
		testCase.assertions(t, db, updateData.Id, err)
	}
}

type testCaseCreateOrganization struct {
	name             string
	organizationName string
	createdBy        null.Int
	assertions       func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error)
	mutations        func(t *testing.T, db *sqlx.DB, organization *model.CreateOrganization)
}

func getAddOrganizationTestCases() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
		{
			name:             "success",
			organizationName: "Example Organization",
			createdBy:        null.IntFrom(3),
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.NotNil(t, organization, "unexpected nil organization")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
		{
			name:             "fail-name-too-long",
			organizationName: strings.Repeat("Demby", 256),
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.Nil(t, organization, "unexpected non-nil organization")
				assert.Error(t, err, "expected an error due to name exceeding length limit")
			},
		},
	}
}

func Test_AddOrganization(t *testing.T) {
	for _, testCase := range getAddOrganizationTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			organization := &model.CreateOrganization{
				Name:   testCase.organizationName,
				UserId: 1,
			}

			createOrganization, err := m.AddOrganization(testCtx, txHandlerDb, organization)
			testCase.assertions(t, db, createOrganization, err)
		})
	}
}

type testCaseRestoreOrganization struct {
	name                string
	id                  int
	restoreOrganization model.Organization
	assertions          func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations           func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int)
}

func getRestoreOrganizationTestCases() []testCaseRestoreOrganization {
	return []testCaseRestoreOrganization{
		{
			name: "success",
			restoreOrganization: model.Organization{
				Id:       1,
				Name:     "Organization A",
				IsActive: false,
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")
				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int) {
				entry := mysqlmodel.Organization{
					ID:       organization.Id,
					Name:     organization.Name,
					IsActive: true,
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")
				assert.Equal(t, true, entry.IsActive)
				return organization.Id
			},
		},
		{
			name: "fail-missing-entry-to-update",
			restoreOrganization: model.Organization{
				Id:       99999,
				Name:     "Non-existent Organization",
				IsActive: false,
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				_, fetchErr := mysqlmodel.FindOrganization(context.TODO(), db, id)
				assert.Error(t, fetchErr, "organization should not exist")
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int) {
				return organization.Id
			},
		},
	}
}

func Test_RestoreOrganization(t *testing.T) {
	for _, testCase := range getRestoreOrganizationTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			id := testCase.mutations(t, db, &testCase.restoreOrganization)
			err = m.RestoreOrganization(testCtx, txHandlerDb, id)
			testCase.assertions(t, db, id, err)
		})
	}
}
