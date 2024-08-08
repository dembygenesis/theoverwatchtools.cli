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
)

type testCaseGetOrganizations struct {
	name       string
	filter     *model.OrganizationFilters
	mutations  func(t *testing.T, db *sqlx.DB, organization *model.Organization, userId int)
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error)
}

func getTestCasesGetOrganizations() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name: "success-filter-ids-in",
			filter: &model.OrganizationFilters{
				IdsIn: []int{4},
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
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization, userId int) {
				entry := mysqlmodel.Organization{
					ID:            organization.Id,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(userId),
					LastUpdatedBy: null.IntFrom(userId),
					CreatedAt:     organization.CreatedAt,
					LastUpdatedAt: organization.LastUpdatedAt,
					IsActive:      organization.IsActive,
				}

				err := entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
		},
		//{
		//	name: "success-filter-names-in",
		//	filter: &model.OrganizationFilters{
		//		OrganizationNameIn: []string{"Lawrence", "Younes"},
		//	},
		//	mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization, user *model.User, categoryType *model.CategoryType, category *model.Category) {
		//		entryCatType := mysqlmodel.CategoryType{
		//			ID:   categoryType.Id,
		//			Name: categoryType.Name,
		//		}
		//
		//		err := entryCatType.Insert(context.Background(), db, boil.Infer())
		//		require.NoError(t, err, "error inserting into category type")
		//
		//		entryCategory := mysqlmodel.Category{
		//			ID:                category.Id,
		//			CategoryTypeRefID: category.CategoryTypeRefId,
		//			Name:              category.Name,
		//		}
		//
		//		err = entryCategory.Insert(context.Background(), db, boil.Infer())
		//		require.NoError(t, err, "error inserting into category table")
		//
		//		entryUser := mysqlmodel.User{
		//			Firstname:         user.Firstname,
		//			Lastname:          user.Lastname,
		//			Email:             user.Email,
		//			Password:          user.Password,
		//			CategoryTypeRefID: user.CategoryTypeRefId,
		//		}
		//
		//		err = entryUser.Insert(context.Background(), db, boil.Infer())
		//		require.NoError(t, err, "error inserting in the user db")
		//
		//		var createdByID, lastUpdatedByID int
		//
		//		err = db.QueryRow("SELECT id FROM user WHERE firstname = ?", organization.CreatedBy).Scan(&createdByID)
		//		require.NoError(t, err, "error converting CreatedBy to ID")
		//
		//		err = db.QueryRow("SELECT id FROM user WHERE firstname = ?", organization.LastUpdatedBy).Scan(&lastUpdatedByID)
		//		require.NoError(t, err, "error converting LastUpdatedBy to ID")
		//
		//		entryOrg1 := mysqlmodel.Organization{
		//			ID:            1,
		//			Name:          "Lawrence",
		//			CreatedBy:     null.IntFrom(createdByID),
		//			LastUpdatedBy: null.IntFrom(lastUpdatedByID),
		//			CreatedAt:     organization.CreatedAt,
		//			LastUpdatedAt: organization.LastUpdatedAt,
		//			IsActive:      organization.IsActive,
		//		}
		//
		//		entryOrg2 := mysqlmodel.Organization{
		//			ID:            2,
		//			Name:          "Younes",
		//			CreatedBy:     null.IntFrom(createdByID),
		//			LastUpdatedBy: null.IntFrom(lastUpdatedByID),
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

			userId := 1

			organization := &model.Organization{
				Id:            4,
				CreatedBy:     "Demby",
				LastUpdatedBy: "Demby",
				IsActive:      true,
			}

			testCase.mutations(t, db, organization, userId)
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
//				CreatedBy: null.IntFrom(1),
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
//
//type testCaseUpdateOrganization struct {
//	name                string
//	initialOrganization mysqlmodel.Organization
//	updateOrganization  model.UpdateOrganization
//	assertions          func(t *testing.T, db *sqlx.DB, id int, err error)
//	mutations           func(t *testing.T, db *sqlx.DB, organization *model.UpdateOrganization) (updateData model.UpdateOrganization)
//}
//
//func getUpdateOrganizationTestCases() []testCaseUpdateOrganization {
//	return []testCaseUpdateOrganization{
//		{
//			name: "success",
//			initialOrganization: mysqlmodel.Organization{
//				ID:            1,
//				Name:          "Organization A",
//				CreatedBy:     null.IntFrom(1),
//				LastUpdatedBy: null.IntFrom(1),
//				CreatedAt:     time.Now(),
//				LastUpdatedAt: null.TimeFrom(time.Now()),
//				IsActive:      true,
//			},
//			updateOrganization: model.UpdateOrganization{
//				Id:   1,
//				Name: null.StringFrom("Organization A new"),
//			},
//			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
//				require.Nil(t, err, "unexpected non-nil error")
//				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
//				require.NoError(t, err, "unexpected error fetching the organization")
//
//				assert.Equal(t, true, entry.IsActive)
//			},
//			mutations: func(t *testing.T, db *sqlx.DB, organization *model.UpdateOrganization) (updateData model.UpdateOrganization) {
//				updateData = model.UpdateOrganization{
//					Id:   organization.Id,
//					Name: organization.Name,
//				}
//				return updateData
//			},
//		},
//		{
//			name: "fail",
//			initialOrganization: mysqlmodel.Organization{
//				ID:            1,
//				Name:          "Organization A",
//				CreatedBy:     null.IntFrom(1),
//				LastUpdatedBy: null.IntFrom(1),
//				CreatedAt:     time.Now(),
//				LastUpdatedAt: null.TimeFrom(time.Now()),
//				IsActive:      true,
//			},
//			updateOrganization: model.UpdateOrganization{
//				Id:   32123,
//				Name: null.StringFrom("Wrong Name"),
//			},
//			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
//				require.Error(t, err, "expected error when updating a non-existent organization")
//				assert.Contains(t, err.Error(), "expected exactly one entry for entity: 32123")
//			},
//			mutations: func(t *testing.T, db *sqlx.DB, organization *model.UpdateOrganization) (updateData model.UpdateOrganization) {
//				updateData = model.UpdateOrganization{
//					Id:   organization.Id,
//					Name: organization.Name,
//				}
//				return updateData
//			},
//		},
//	}
//}
//
//func Test_UpdateOrganization(t *testing.T) {
//	for _, testCase := range getUpdateOrganizationTestCases() {
//		db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
//		defer cleanup()
//
//		cfg := &Config{
//			Logger:        testLogger,
//			QueryTimeouts: testQueryTimeouts,
//		}
//
//		m, err := New(cfg)
//		require.NoError(t, err, "unexpected error")
//		require.NotNil(t, m, "unexpected nil")
//
//		txHandler, err := mysqltx.New(&mysqltx.Config{
//			Logger:       testLogger,
//			Db:           db,
//			DatabaseName: cp.Database,
//		})
//		require.NoError(t, err, "unexpected error creating the tx handler")
//
//		txHandlerDb, err := txHandler.Db(testCtx)
//		require.NoError(t, err, "unexpected error fetching the db from the tx handler")
//		require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")
//
//		err = testCase.initialOrganization.Insert(context.Background(), db, boil.Infer())
//		require.NoError(t, err, "unexpected error inserting initial organization")
//
//		updateData := testCase.mutations(t, db, &testCase.updateOrganization)
//
//		_, err = m.UpdateOrganization(testCtx, txHandlerDb, &updateData)
//		testCase.assertions(t, db, updateData.Id, err)
//	}
//}

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
