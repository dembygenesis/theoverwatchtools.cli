package mysqlstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (m *Repository) GetClickTrackers(ctx context.Context, tx persistence.TransactionHandler, filters *model.ClickTrackerFilters) (*model.PaginatedClickTrackers, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	fmt.Println("the filters ---- ", strutil.GetAsJson(filters))

	res, err := m.getClickTrackers(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read clickTrackers: %v", err)
	}

	fmt.Println("the res -- ", strutil.GetAsJson(res))

	return res, nil
}

// DropClickTrackersTable drops the category table (for testing purposes).
func (m *Repository) DropClickTrackersTable(
	ctx context.Context,
	tx persistence.TransactionHandler,
) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("extract context executor: %v", err)
	}

	dropStmts := []string{
		"SET FOREIGN_KEY_CHECKS = 0;",
		"DROP TABLE click_trackers;",
		"SET FOREIGN_KEY_CHECKS = 1;",
	}

	for _, stmt := range dropStmts {
		if _, err := queries.Raw(stmt).Exec(ctxExec); err != nil {
			return fmt.Errorf("dropping click trackers table sql stmt: %v", err)
		}
	}

	return nil
}

func (m *Repository) getClickTrackers(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.ClickTrackerFilters,
) (*model.PaginatedClickTrackers, error) {
	var (
		paginated  model.PaginatedClickTrackers
		pagination = model.NewPagination()
		res        = make([]model.ClickTracker, 0)
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	fmt.Println("the id ----- ", strutil.GetAsJson(filters.ClickTrackerSetIdIn))

	queryMods := []qm.QueryMod{
		qm.InnerJoin(
			fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.ID,
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.ClickTrackerSetID,
			),
		),
		qm.Select(
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.ID,
				mysqlmodel.ClickTrackerColumns.ID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.Name,
				mysqlmodel.ClickTrackerColumns.Name,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.URLName,
				mysqlmodel.ClickTrackerColumns.URLName,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.RedirectURL,
				mysqlmodel.ClickTrackerColumns.RedirectURL,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.Clicks,
				mysqlmodel.ClickTrackerColumns.Clicks,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.UniqueClicks,
				mysqlmodel.ClickTrackerColumns.UniqueClicks,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.CreatedBy,
				mysqlmodel.ClickTrackerColumns.CreatedBy,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.UpdatedBy,
				mysqlmodel.ClickTrackerColumns.UpdatedBy,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.ClickTrackerSetID,
				mysqlmodel.ClickTrackerColumns.ClickTrackerSetID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.CountryID,
				mysqlmodel.ClickTrackerColumns.CountryID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.CreatedAt,
				mysqlmodel.ClickTrackerColumns.CreatedAt,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.UpdatedAt,
				mysqlmodel.ClickTrackerColumns.UpdatedAt,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackers,
				mysqlmodel.ClickTrackerColumns.DeletedAt,
				mysqlmodel.ClickTrackerColumns.DeletedAt,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.ID,
				"click_tracker_set_id",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.Name,
				"click_tracker_set_name",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.URLName,
				"click_tracker_set_url_name",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.CreatedBy,
				"click_tracker_set_created_by",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.UpdatedBy,
				"click_tracker_set_updated_by",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.OrganizationID,
				"click_tracker_set_organization_id",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.AnalyticsNumberOfLinks,
				"click_tracker_set_analytics_number_of_links",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.AnalyticsLastUpdatedAt,
				"click_tracker_set_analytics_last_updated_at",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.CreatedAt,
				"click_tracker_set_created_at",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.UpdatedAt,
				"click_tracker_set_updated_at",
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTrackerSets,
				mysqlmodel.ClickTrackerSetColumns.DeletedAt,
				"click_tracker_set_deleted_at",
			),
		),
	}

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.ID.IN(filters.IdsIn))
		}

		if len(filters.NameIn) > 0 {
			fmt.Println("did I make it here come on???")
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.Name.IN(filters.NameIn))
		}

		if len(filters.UrlNameIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.URLName.IN(filters.UrlNameIn))
		}

		if len(filters.RedirectUrlIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.RedirectURL.IN(filters.RedirectUrlIn))
		}

		if len(filters.ClicksIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.Clicks.IN(filters.ClicksIn))
		}

		if len(filters.UniqueClicksIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.UniqueClicks.IN(filters.UniqueClicksIn))
		}

		if len(filters.CreatedByIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.CreatedBy.IN(filters.CreatedByIn))
		}

		if len(filters.UpdatedByIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.UpdatedBy.IN(filters.UpdatedByIn))
		}

		if len(filters.ClickTrackerSetIdIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.ClickTrackerSetID.IN(filters.ClickTrackerSetIdIn))
		}

		if len(filters.CountryIdIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.CountryID.IN(filters.CountryIdIn))
		}
	}

	fmt.Println("the queryMods for CT ----- ", queryMods)

	q := mysqlmodel.ClickTrackers(queryMods...)
	totalCount, err := q.Count(ctx, ctxExec)
	if err != nil {
		return nil, fmt.Errorf("get click trackers count: %v", err)
	}

	page := pagination.Page
	maxRows := pagination.MaxRows
	if filters != nil {
		if filters.Page.Valid {
			page = filters.Page.Int
		}
		if filters.MaxRows.Valid {
			maxRows = filters.MaxRows.Int
		}
	}

	pagination.SetQueryBoundaries(page, maxRows, int(totalCount))

	queryMods = append(queryMods, qm.Limit(pagination.MaxRows), qm.Offset(pagination.Offset))
	q = mysqlmodel.ClickTrackers(queryMods...)

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get click trackers: %v", err)
	}

	pagination.RowCount = len(res)
	paginated.ClickTrackers = res
	paginated.Pagination = pagination

	return &paginated, nil
}

func (m *Repository) GetClickTrackerSetById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.ClickTrackerSet, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	var clickTrackerSet mysqlmodel.ClickTrackerSet
	err = queries.Raw(
		"SELECT * FROM click_tracker_sets WHERE id = ?",
		id,
	).Bind(ctx, ctxExec, &clickTrackerSet)

	if err != nil {
		return nil, fmt.Errorf("get click tracker set by id: %v", err)
	}

	if clickTrackerSet.ID == 0 {
		return nil, fmt.Errorf("click tracker set not found")
	}

	return &model.ClickTrackerSet{
		ID:                     clickTrackerSet.ID,
		Name:                   clickTrackerSet.Name,
		UrlName:                clickTrackerSet.URLName,
		CreatedBy:              clickTrackerSet.CreatedBy,
		UpdatedBy:              clickTrackerSet.UpdatedBy,
		OrganizationID:         clickTrackerSet.OrganizationID,
		AnalyticsNumberOfLinks: clickTrackerSet.AnalyticsNumberOfLinks,
		AnalyticsLastUpdatedAt: clickTrackerSet.AnalyticsLastUpdatedAt,
		CreatedAt:              clickTrackerSet.CreatedAt,
		UpdatedAt:              clickTrackerSet.UpdatedAt,
		DeletedAt:              clickTrackerSet.DeletedAt,
	}, nil
}

// GetClickTrackerById attempts to fetch the click tracker by its ID.
func (m *Repository) GetClickTrackerById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.ClickTracker, error) {
	paginated, err := m.GetClickTrackers(ctx, tx, &model.ClickTrackerFilters{
		IdsIn: []int{id},
	})
	if err != nil {
		return nil, fmt.Errorf("click tracker filtered by id: %v", err)
	}

	if paginated.Pagination.RowCount != 1 {
		return nil, fmt.Errorf(sysconsts.ErrExpectedExactlyOneEntry, id)
	}

	return &paginated.ClickTrackers[0], nil
}

func (m *Repository) GetClickTrackerByName(ctx context.Context, tx persistence.TransactionHandler, name string) (*model.ClickTracker, error) {
	paginated, err := m.GetClickTrackers(ctx, tx, &model.ClickTrackerFilters{
		NameIn: []string{name},
	})

	fmt.Println("the paginated ---- ", strutil.GetAsJson(paginated))

	if err != nil {
		return nil, fmt.Errorf("click tracker filtered by name: %v", err)
	}

	//if paginated.Pagination.RowCount == 0 {
	//	return nil, nil
	//}

	if paginated.Pagination.RowCount != 1 {
		return nil, errors.New("expected exactly one click tracker entry")
	}

	return &paginated.ClickTrackers[0], nil
}

func (m *Repository) CreateClickTracker(ctx context.Context, tx persistence.TransactionHandler, clickTracker *model.ClickTracker) (*model.ClickTracker, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	entry := mysqlmodel.ClickTracker{
		Name:              clickTracker.Name,
		URLName:           clickTracker.UrlName,
		RedirectURL:       clickTracker.RedirectUrl,
		Clicks:            clickTracker.Clicks,
		UniqueClicks:      clickTracker.UniqueClicks,
		CreatedBy:         clickTracker.CreatedBy,
		UpdatedBy:         clickTracker.UpdatedBy,
		ClickTrackerSetID: clickTracker.ClickTrackerSetId,
		CountryID:         clickTracker.CountryId,
	}
	if err = entry.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert click tracker: %v", err)
	}

	clickTracker, err = m.GetClickTrackerById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get click tracker by id: %v", err)
	}

	return clickTracker, nil
}
