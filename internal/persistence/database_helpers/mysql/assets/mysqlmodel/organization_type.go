// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package mysqlmodel

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// OrganizationType is an object representing the database table.
type OrganizationType struct {
	ID          int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	CreatedBy   null.Int  `boil:"created_by" json:"created_by,omitempty" toml:"created_by" yaml:"created_by,omitempty"`
	CreatedDate null.Time `boil:"created_date" json:"created_date,omitempty" toml:"created_date" yaml:"created_date,omitempty"`
	LastUpdated null.Time `boil:"last_updated" json:"last_updated,omitempty" toml:"last_updated" yaml:"last_updated,omitempty"`
	UpdatedBy   null.Int  `boil:"updated_by" json:"updated_by,omitempty" toml:"updated_by" yaml:"updated_by,omitempty"`
	IsActive    null.Int  `boil:"is_active" json:"is_active,omitempty" toml:"is_active" yaml:"is_active,omitempty"`
	Name        string    `boil:"name" json:"name" toml:"name" yaml:"name"`

	R *organizationTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L organizationTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var OrganizationTypeColumns = struct {
	ID          string
	CreatedBy   string
	CreatedDate string
	LastUpdated string
	UpdatedBy   string
	IsActive    string
	Name        string
}{
	ID:          "id",
	CreatedBy:   "created_by",
	CreatedDate: "created_date",
	LastUpdated: "last_updated",
	UpdatedBy:   "updated_by",
	IsActive:    "is_active",
	Name:        "name",
}

var OrganizationTypeTableColumns = struct {
	ID          string
	CreatedBy   string
	CreatedDate string
	LastUpdated string
	UpdatedBy   string
	IsActive    string
	Name        string
}{
	ID:          "organization_type.id",
	CreatedBy:   "organization_type.created_by",
	CreatedDate: "organization_type.created_date",
	LastUpdated: "organization_type.last_updated",
	UpdatedBy:   "organization_type.updated_by",
	IsActive:    "organization_type.is_active",
	Name:        "organization_type.name",
}

// Generated where

var OrganizationTypeWhere = struct {
	ID          whereHelperint
	CreatedBy   whereHelpernull_Int
	CreatedDate whereHelpernull_Time
	LastUpdated whereHelpernull_Time
	UpdatedBy   whereHelpernull_Int
	IsActive    whereHelpernull_Int
	Name        whereHelperstring
}{
	ID:          whereHelperint{field: "`organization_type`.`id`"},
	CreatedBy:   whereHelpernull_Int{field: "`organization_type`.`created_by`"},
	CreatedDate: whereHelpernull_Time{field: "`organization_type`.`created_date`"},
	LastUpdated: whereHelpernull_Time{field: "`organization_type`.`last_updated`"},
	UpdatedBy:   whereHelpernull_Int{field: "`organization_type`.`updated_by`"},
	IsActive:    whereHelpernull_Int{field: "`organization_type`.`is_active`"},
	Name:        whereHelperstring{field: "`organization_type`.`name`"},
}

// OrganizationTypeRels is where relationship names are stored.
var OrganizationTypeRels = struct {
	OrganizationTypeRefOrganizations string
}{
	OrganizationTypeRefOrganizations: "OrganizationTypeRefOrganizations",
}

// organizationTypeR is where relationships are stored.
type organizationTypeR struct {
	OrganizationTypeRefOrganizations OrganizationSlice `boil:"OrganizationTypeRefOrganizations" json:"OrganizationTypeRefOrganizations" toml:"OrganizationTypeRefOrganizations" yaml:"OrganizationTypeRefOrganizations"`
}

// NewStruct creates a new relationship struct
func (*organizationTypeR) NewStruct() *organizationTypeR {
	return &organizationTypeR{}
}

func (r *organizationTypeR) GetOrganizationTypeRefOrganizations() OrganizationSlice {
	if r == nil {
		return nil
	}
	return r.OrganizationTypeRefOrganizations
}

// organizationTypeL is where Load methods for each relationship are stored.
type organizationTypeL struct{}

var (
	organizationTypeAllColumns            = []string{"id", "created_by", "created_date", "last_updated", "updated_by", "is_active", "name"}
	organizationTypeColumnsWithoutDefault = []string{"created_by", "updated_by", "is_active", "name"}
	organizationTypeColumnsWithDefault    = []string{"id", "created_date", "last_updated"}
	organizationTypePrimaryKeyColumns     = []string{"id"}
	organizationTypeGeneratedColumns      = []string{}
)

type (
	// OrganizationTypeSlice is an alias for a slice of pointers to OrganizationType.
	// This should almost always be used instead of []OrganizationType.
	OrganizationTypeSlice []*OrganizationType
	// OrganizationTypeHook is the signature for custom OrganizationType hook methods
	OrganizationTypeHook func(context.Context, boil.ContextExecutor, *OrganizationType) error

	organizationTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	organizationTypeType                 = reflect.TypeOf(&OrganizationType{})
	organizationTypeMapping              = queries.MakeStructMapping(organizationTypeType)
	organizationTypePrimaryKeyMapping, _ = queries.BindMapping(organizationTypeType, organizationTypeMapping, organizationTypePrimaryKeyColumns)
	organizationTypeInsertCacheMut       sync.RWMutex
	organizationTypeInsertCache          = make(map[string]insertCache)
	organizationTypeUpdateCacheMut       sync.RWMutex
	organizationTypeUpdateCache          = make(map[string]updateCache)
	organizationTypeUpsertCacheMut       sync.RWMutex
	organizationTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var organizationTypeAfterSelectMu sync.Mutex
var organizationTypeAfterSelectHooks []OrganizationTypeHook

var organizationTypeBeforeInsertMu sync.Mutex
var organizationTypeBeforeInsertHooks []OrganizationTypeHook
var organizationTypeAfterInsertMu sync.Mutex
var organizationTypeAfterInsertHooks []OrganizationTypeHook

var organizationTypeBeforeUpdateMu sync.Mutex
var organizationTypeBeforeUpdateHooks []OrganizationTypeHook
var organizationTypeAfterUpdateMu sync.Mutex
var organizationTypeAfterUpdateHooks []OrganizationTypeHook

var organizationTypeBeforeDeleteMu sync.Mutex
var organizationTypeBeforeDeleteHooks []OrganizationTypeHook
var organizationTypeAfterDeleteMu sync.Mutex
var organizationTypeAfterDeleteHooks []OrganizationTypeHook

var organizationTypeBeforeUpsertMu sync.Mutex
var organizationTypeBeforeUpsertHooks []OrganizationTypeHook
var organizationTypeAfterUpsertMu sync.Mutex
var organizationTypeAfterUpsertHooks []OrganizationTypeHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *OrganizationType) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organizationTypeAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *OrganizationType) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organizationTypeBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *OrganizationType) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organizationTypeAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *OrganizationType) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organizationTypeBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *OrganizationType) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organizationTypeAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *OrganizationType) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organizationTypeBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *OrganizationType) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organizationTypeAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *OrganizationType) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organizationTypeBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *OrganizationType) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organizationTypeAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddOrganizationTypeHook registers your hook function for all future operations.
func AddOrganizationTypeHook(hookPoint boil.HookPoint, organizationTypeHook OrganizationTypeHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		organizationTypeAfterSelectMu.Lock()
		organizationTypeAfterSelectHooks = append(organizationTypeAfterSelectHooks, organizationTypeHook)
		organizationTypeAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		organizationTypeBeforeInsertMu.Lock()
		organizationTypeBeforeInsertHooks = append(organizationTypeBeforeInsertHooks, organizationTypeHook)
		organizationTypeBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		organizationTypeAfterInsertMu.Lock()
		organizationTypeAfterInsertHooks = append(organizationTypeAfterInsertHooks, organizationTypeHook)
		organizationTypeAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		organizationTypeBeforeUpdateMu.Lock()
		organizationTypeBeforeUpdateHooks = append(organizationTypeBeforeUpdateHooks, organizationTypeHook)
		organizationTypeBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		organizationTypeAfterUpdateMu.Lock()
		organizationTypeAfterUpdateHooks = append(organizationTypeAfterUpdateHooks, organizationTypeHook)
		organizationTypeAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		organizationTypeBeforeDeleteMu.Lock()
		organizationTypeBeforeDeleteHooks = append(organizationTypeBeforeDeleteHooks, organizationTypeHook)
		organizationTypeBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		organizationTypeAfterDeleteMu.Lock()
		organizationTypeAfterDeleteHooks = append(organizationTypeAfterDeleteHooks, organizationTypeHook)
		organizationTypeAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		organizationTypeBeforeUpsertMu.Lock()
		organizationTypeBeforeUpsertHooks = append(organizationTypeBeforeUpsertHooks, organizationTypeHook)
		organizationTypeBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		organizationTypeAfterUpsertMu.Lock()
		organizationTypeAfterUpsertHooks = append(organizationTypeAfterUpsertHooks, organizationTypeHook)
		organizationTypeAfterUpsertMu.Unlock()
	}
}

// One returns a single organizationType record from the query.
func (q organizationTypeQuery) One(ctx context.Context, exec boil.ContextExecutor) (*OrganizationType, error) {
	o := &OrganizationType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysqlmodel: failed to execute a one query for organization_type")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all OrganizationType records from the query.
func (q organizationTypeQuery) All(ctx context.Context, exec boil.ContextExecutor) (OrganizationTypeSlice, error) {
	var o []*OrganizationType

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "mysqlmodel: failed to assign all query results to OrganizationType slice")
	}

	if len(organizationTypeAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all OrganizationType records in the query.
func (q organizationTypeQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to count organization_type rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q organizationTypeQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mysqlmodel: failed to check if organization_type exists")
	}

	return count > 0, nil
}

// OrganizationTypeRefOrganizations retrieves all the organization's Organizations with an executor via organization_type_ref_id column.
func (o *OrganizationType) OrganizationTypeRefOrganizations(mods ...qm.QueryMod) organizationQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("`organization`.`organization_type_ref_id`=?", o.ID),
	)

	return Organizations(queryMods...)
}

// LoadOrganizationTypeRefOrganizations allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (organizationTypeL) LoadOrganizationTypeRefOrganizations(ctx context.Context, e boil.ContextExecutor, singular bool, maybeOrganizationType interface{}, mods queries.Applicator) error {
	var slice []*OrganizationType
	var object *OrganizationType

	if singular {
		var ok bool
		object, ok = maybeOrganizationType.(*OrganizationType)
		if !ok {
			object = new(OrganizationType)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeOrganizationType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeOrganizationType))
			}
		}
	} else {
		s, ok := maybeOrganizationType.(*[]*OrganizationType)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeOrganizationType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeOrganizationType))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &organizationTypeR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &organizationTypeR{}
			}
			args[obj.ID] = struct{}{}
		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`organization`),
		qm.WhereIn(`organization.organization_type_ref_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load organization")
	}

	var resultSlice []*Organization
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice organization")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on organization")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for organization")
	}

	if len(organizationAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.OrganizationTypeRefOrganizations = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &organizationR{}
			}
			foreign.R.OrganizationTypeRef = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.OrganizationTypeRefID {
				local.R.OrganizationTypeRefOrganizations = append(local.R.OrganizationTypeRefOrganizations, foreign)
				if foreign.R == nil {
					foreign.R = &organizationR{}
				}
				foreign.R.OrganizationTypeRef = local
				break
			}
		}
	}

	return nil
}

// AddOrganizationTypeRefOrganizations adds the given related objects to the existing relationships
// of the organization_type, optionally inserting them as new records.
// Appends related to o.R.OrganizationTypeRefOrganizations.
// Sets related.R.OrganizationTypeRef appropriately.
func (o *OrganizationType) AddOrganizationTypeRefOrganizations(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Organization) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.OrganizationTypeRefID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE `organization` SET %s WHERE %s",
				strmangle.SetParamNames("`", "`", 0, []string{"organization_type_ref_id"}),
				strmangle.WhereClause("`", "`", 0, organizationPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.OrganizationTypeRefID = o.ID
		}
	}

	if o.R == nil {
		o.R = &organizationTypeR{
			OrganizationTypeRefOrganizations: related,
		}
	} else {
		o.R.OrganizationTypeRefOrganizations = append(o.R.OrganizationTypeRefOrganizations, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &organizationR{
				OrganizationTypeRef: o,
			}
		} else {
			rel.R.OrganizationTypeRef = o
		}
	}
	return nil
}

// OrganizationTypes retrieves all the records using an executor.
func OrganizationTypes(mods ...qm.QueryMod) organizationTypeQuery {
	mods = append(mods, qm.From("`organization_type`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`organization_type`.*"})
	}

	return organizationTypeQuery{q}
}

// FindOrganizationType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindOrganizationType(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*OrganizationType, error) {
	organizationTypeObj := &OrganizationType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `organization_type` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, organizationTypeObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysqlmodel: unable to select from organization_type")
	}

	if err = organizationTypeObj.doAfterSelectHooks(ctx, exec); err != nil {
		return organizationTypeObj, err
	}

	return organizationTypeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *OrganizationType) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("mysqlmodel: no organization_type provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(organizationTypeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	organizationTypeInsertCacheMut.RLock()
	cache, cached := organizationTypeInsertCache[key]
	organizationTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			organizationTypeAllColumns,
			organizationTypeColumnsWithDefault,
			organizationTypeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(organizationTypeType, organizationTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(organizationTypeType, organizationTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `organization_type` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `organization_type` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `organization_type` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, organizationTypePrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "mysqlmodel: unable to insert into organization_type")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == organizationTypeMapping["id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "mysqlmodel: unable to populate default values for organization_type")
	}

CacheNoHooks:
	if !cached {
		organizationTypeInsertCacheMut.Lock()
		organizationTypeInsertCache[key] = cache
		organizationTypeInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the OrganizationType.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *OrganizationType) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	organizationTypeUpdateCacheMut.RLock()
	cache, cached := organizationTypeUpdateCache[key]
	organizationTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			organizationTypeAllColumns,
			organizationTypePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("mysqlmodel: unable to update organization_type, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `organization_type` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, organizationTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(organizationTypeType, organizationTypeMapping, append(wl, organizationTypePrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to update organization_type row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to get rows affected by update for organization_type")
	}

	if !cached {
		organizationTypeUpdateCacheMut.Lock()
		organizationTypeUpdateCache[key] = cache
		organizationTypeUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q organizationTypeQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to update all for organization_type")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to retrieve rows affected for organization_type")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o OrganizationTypeSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("mysqlmodel: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), organizationTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `organization_type` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, organizationTypePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to update all in organizationType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to retrieve rows affected all in update all organizationType")
	}
	return rowsAff, nil
}

var mySQLOrganizationTypeUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *OrganizationType) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("mysqlmodel: no organization_type provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(organizationTypeColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLOrganizationTypeUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	organizationTypeUpsertCacheMut.RLock()
	cache, cached := organizationTypeUpsertCache[key]
	organizationTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			organizationTypeAllColumns,
			organizationTypeColumnsWithDefault,
			organizationTypeColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			organizationTypeAllColumns,
			organizationTypePrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("mysqlmodel: unable to upsert organization_type, could not build update column list")
		}

		ret := strmangle.SetComplement(organizationTypeAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`organization_type`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `organization_type` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(organizationTypeType, organizationTypeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(organizationTypeType, organizationTypeMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "mysqlmodel: unable to upsert for organization_type")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == organizationTypeMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(organizationTypeType, organizationTypeMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "mysqlmodel: unable to retrieve unique values for organization_type")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "mysqlmodel: unable to populate default values for organization_type")
	}

CacheNoHooks:
	if !cached {
		organizationTypeUpsertCacheMut.Lock()
		organizationTypeUpsertCache[key] = cache
		organizationTypeUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single OrganizationType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *OrganizationType) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("mysqlmodel: no OrganizationType provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), organizationTypePrimaryKeyMapping)
	sql := "DELETE FROM `organization_type` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to delete from organization_type")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to get rows affected by delete for organization_type")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q organizationTypeQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("mysqlmodel: no organizationTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to delete all from organization_type")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to get rows affected by deleteall for organization_type")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o OrganizationTypeSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(organizationTypeBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), organizationTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `organization_type` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, organizationTypePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to delete all from organizationType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to get rows affected by deleteall for organization_type")
	}

	if len(organizationTypeAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *OrganizationType) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindOrganizationType(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *OrganizationTypeSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := OrganizationTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), organizationTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `organization_type`.* FROM `organization_type` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, organizationTypePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "mysqlmodel: unable to reload all in OrganizationTypeSlice")
	}

	*o = slice

	return nil
}

// OrganizationTypeExists checks if the OrganizationType row exists.
func OrganizationTypeExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `organization_type` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mysqlmodel: unable to check if organization_type exists")
	}

	return exists, nil
}

// Exists checks if the OrganizationType row exists.
func (o *OrganizationType) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return OrganizationTypeExists(ctx, exec, o.ID)
}
