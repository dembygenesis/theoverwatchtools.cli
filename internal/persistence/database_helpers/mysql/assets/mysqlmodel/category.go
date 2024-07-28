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

// Category is an object representing the database table.
type Category struct {
	ID                int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	CreatedBy         null.Int  `boil:"created_by" json:"created_by,omitempty" toml:"created_by" yaml:"created_by,omitempty"`
	CreatedDate       null.Time `boil:"created_date" json:"created_date,omitempty" toml:"created_date" yaml:"created_date,omitempty"`
	LastUpdated       null.Time `boil:"last_updated" json:"last_updated,omitempty" toml:"last_updated" yaml:"last_updated,omitempty"`
	UpdatedBy         null.Int  `boil:"updated_by" json:"updated_by,omitempty" toml:"updated_by" yaml:"updated_by,omitempty"`
	IsActive          int       `boil:"is_active" json:"is_active" toml:"is_active" yaml:"is_active"`
	CategoryTypeRefID int       `boil:"category_type_ref_id" json:"category_type_ref_id" toml:"category_type_ref_id" yaml:"category_type_ref_id"`
	Name              string    `boil:"name" json:"name" toml:"name" yaml:"name"`

	R *categoryR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L categoryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CategoryColumns = struct {
	ID                string
	CreatedBy         string
	CreatedDate       string
	LastUpdated       string
	UpdatedBy         string
	IsActive          string
	CategoryTypeRefID string
	Name              string
}{
	ID:                "id",
	CreatedBy:         "created_by",
	CreatedDate:       "created_date",
	LastUpdated:       "last_updated",
	UpdatedBy:         "updated_by",
	IsActive:          "is_active",
	CategoryTypeRefID: "category_type_ref_id",
	Name:              "name",
}

var CategoryTableColumns = struct {
	ID                string
	CreatedBy         string
	CreatedDate       string
	LastUpdated       string
	UpdatedBy         string
	IsActive          string
	CategoryTypeRefID string
	Name              string
}{
	ID:                "category.id",
	CreatedBy:         "category.created_by",
	CreatedDate:       "category.created_date",
	LastUpdated:       "category.last_updated",
	UpdatedBy:         "category.updated_by",
	IsActive:          "category.is_active",
	CategoryTypeRefID: "category.category_type_ref_id",
	Name:              "category.name",
}

// Generated where

type whereHelperint struct{ field string }

func (w whereHelperint) EQ(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint) NEQ(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint) LT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint) LTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint) GT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint) GTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpernull_Int struct{ field string }

func (w whereHelpernull_Int) EQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Int) NEQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Int) LT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Int) LTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Int) GT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Int) GTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelpernull_Int) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelpernull_Int) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

func (w whereHelpernull_Int) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Int) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

type whereHelpernull_Time struct{ field string }

func (w whereHelpernull_Time) EQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Time) NEQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Time) LT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Time) LTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Time) GT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Time) GTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Time) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Time) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod   { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod   { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod   { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) LIKE(x string) qm.QueryMod  { return qm.Where(w.field+" LIKE ?", x) }
func (w whereHelperstring) NLIKE(x string) qm.QueryMod { return qm.Where(w.field+" NOT LIKE ?", x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var CategoryWhere = struct {
	ID                whereHelperint
	CreatedBy         whereHelpernull_Int
	CreatedDate       whereHelpernull_Time
	LastUpdated       whereHelpernull_Time
	UpdatedBy         whereHelpernull_Int
	IsActive          whereHelperint
	CategoryTypeRefID whereHelperint
	Name              whereHelperstring
}{
	ID:                whereHelperint{field: "`category`.`id`"},
	CreatedBy:         whereHelpernull_Int{field: "`category`.`created_by`"},
	CreatedDate:       whereHelpernull_Time{field: "`category`.`created_date`"},
	LastUpdated:       whereHelpernull_Time{field: "`category`.`last_updated`"},
	UpdatedBy:         whereHelpernull_Int{field: "`category`.`updated_by`"},
	IsActive:          whereHelperint{field: "`category`.`is_active`"},
	CategoryTypeRefID: whereHelperint{field: "`category`.`category_type_ref_id`"},
	Name:              whereHelperstring{field: "`category`.`name`"},
}

// CategoryRels is where relationship names are stored.
var CategoryRels = struct {
	CategoryTypeRef      string
	CategoryTypeRefUsers string
}{
	CategoryTypeRef:      "CategoryTypeRef",
	CategoryTypeRefUsers: "CategoryTypeRefUsers",
}

// categoryR is where relationships are stored.
type categoryR struct {
	CategoryTypeRef      *CategoryType `boil:"CategoryTypeRef" json:"CategoryTypeRef" toml:"CategoryTypeRef" yaml:"CategoryTypeRef"`
	CategoryTypeRefUsers UserSlice     `boil:"CategoryTypeRefUsers" json:"CategoryTypeRefUsers" toml:"CategoryTypeRefUsers" yaml:"CategoryTypeRefUsers"`
}

// NewStruct creates a new relationship struct
func (*categoryR) NewStruct() *categoryR {
	return &categoryR{}
}

func (r *categoryR) GetCategoryTypeRef() *CategoryType {
	if r == nil {
		return nil
	}
	return r.CategoryTypeRef
}

func (r *categoryR) GetCategoryTypeRefUsers() UserSlice {
	if r == nil {
		return nil
	}
	return r.CategoryTypeRefUsers
}

// categoryL is where Load methods for each relationship are stored.
type categoryL struct{}

var (
	categoryAllColumns            = []string{"id", "created_by", "created_date", "last_updated", "updated_by", "is_active", "category_type_ref_id", "name"}
	categoryColumnsWithoutDefault = []string{"created_by", "last_updated", "updated_by", "category_type_ref_id", "name"}
	categoryColumnsWithDefault    = []string{"id", "created_date", "is_active"}
	categoryPrimaryKeyColumns     = []string{"id"}
	categoryGeneratedColumns      = []string{}
)

type (
	// CategorySlice is an alias for a slice of pointers to Category.
	// This should almost always be used instead of []Category.
	CategorySlice []*Category

	categoryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	categoryType                 = reflect.TypeOf(&Category{})
	categoryMapping              = queries.MakeStructMapping(categoryType)
	categoryPrimaryKeyMapping, _ = queries.BindMapping(categoryType, categoryMapping, categoryPrimaryKeyColumns)
	categoryInsertCacheMut       sync.RWMutex
	categoryInsertCache          = make(map[string]insertCache)
	categoryUpdateCacheMut       sync.RWMutex
	categoryUpdateCache          = make(map[string]updateCache)
	categoryUpsertCacheMut       sync.RWMutex
	categoryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single category record from the query.
func (q categoryQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Category, error) {
	o := &Category{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysqlmodel: failed to execute a one query for category")
	}

	return o, nil
}

// All returns all Category records from the query.
func (q categoryQuery) All(ctx context.Context, exec boil.ContextExecutor) (CategorySlice, error) {
	var o []*Category

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "mysqlmodel: failed to assign all query results to Category slice")
	}

	return o, nil
}

// Count returns the count of all Category records in the query.
func (q categoryQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to count category rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q categoryQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mysqlmodel: failed to check if category exists")
	}

	return count > 0, nil
}

// CategoryTypeRef pointed to by the foreign key.
func (o *Category) CategoryTypeRef(mods ...qm.QueryMod) categoryTypeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`id` = ?", o.CategoryTypeRefID),
	}

	queryMods = append(queryMods, mods...)

	return CategoryTypes(queryMods...)
}

// CategoryTypeRefUsers retrieves all the user's Users with an executor via category_type_ref_id column.
func (o *Category) CategoryTypeRefUsers(mods ...qm.QueryMod) userQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("`user`.`category_type_ref_id`=?", o.ID),
	)

	return Users(queryMods...)
}

// LoadCategoryTypeRef allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (categoryL) LoadCategoryTypeRef(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCategory interface{}, mods queries.Applicator) error {
	var slice []*Category
	var object *Category

	if singular {
		var ok bool
		object, ok = maybeCategory.(*Category)
		if !ok {
			object = new(Category)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCategory))
			}
		}
	} else {
		s, ok := maybeCategory.(*[]*Category)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCategory))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &categoryR{}
		}
		args[object.CategoryTypeRefID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &categoryR{}
			}

			args[obj.CategoryTypeRefID] = struct{}{}

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
		qm.From(`category_type`),
		qm.WhereIn(`category_type.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load CategoryType")
	}

	var resultSlice []*CategoryType
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice CategoryType")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for category_type")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for category_type")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.CategoryTypeRef = foreign
		if foreign.R == nil {
			foreign.R = &categoryTypeR{}
		}
		foreign.R.CategoryTypeRefCategories = append(foreign.R.CategoryTypeRefCategories, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CategoryTypeRefID == foreign.ID {
				local.R.CategoryTypeRef = foreign
				if foreign.R == nil {
					foreign.R = &categoryTypeR{}
				}
				foreign.R.CategoryTypeRefCategories = append(foreign.R.CategoryTypeRefCategories, local)
				break
			}
		}
	}

	return nil
}

// LoadCategoryTypeRefUsers allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (categoryL) LoadCategoryTypeRefUsers(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCategory interface{}, mods queries.Applicator) error {
	var slice []*Category
	var object *Category

	if singular {
		var ok bool
		object, ok = maybeCategory.(*Category)
		if !ok {
			object = new(Category)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCategory))
			}
		}
	} else {
		s, ok := maybeCategory.(*[]*Category)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCategory))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &categoryR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &categoryR{}
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
		qm.From(`user`),
		qm.WhereIn(`user.category_type_ref_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load user")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice user")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on user")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for user")
	}

	if singular {
		object.R.CategoryTypeRefUsers = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &userR{}
			}
			foreign.R.CategoryTypeRef = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.CategoryTypeRefID {
				local.R.CategoryTypeRefUsers = append(local.R.CategoryTypeRefUsers, foreign)
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.CategoryTypeRef = local
				break
			}
		}
	}

	return nil
}

// SetCategoryTypeRef of the category to the related item.
// Sets o.R.CategoryTypeRef to related.
// Adds o to related.R.CategoryTypeRefCategories.
func (o *Category) SetCategoryTypeRef(ctx context.Context, exec boil.ContextExecutor, insert bool, related *CategoryType) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `category` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"category_type_ref_id"}),
		strmangle.WhereClause("`", "`", 0, categoryPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CategoryTypeRefID = related.ID
	if o.R == nil {
		o.R = &categoryR{
			CategoryTypeRef: related,
		}
	} else {
		o.R.CategoryTypeRef = related
	}

	if related.R == nil {
		related.R = &categoryTypeR{
			CategoryTypeRefCategories: CategorySlice{o},
		}
	} else {
		related.R.CategoryTypeRefCategories = append(related.R.CategoryTypeRefCategories, o)
	}

	return nil
}

// AddCategoryTypeRefUsers adds the given related objects to the existing relationships
// of the category, optionally inserting them as new records.
// Appends related to o.R.CategoryTypeRefUsers.
// Sets related.R.CategoryTypeRef appropriately.
func (o *Category) AddCategoryTypeRefUsers(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*User) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.CategoryTypeRefID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE `user` SET %s WHERE %s",
				strmangle.SetParamNames("`", "`", 0, []string{"category_type_ref_id"}),
				strmangle.WhereClause("`", "`", 0, userPrimaryKeyColumns),
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

			rel.CategoryTypeRefID = o.ID
		}
	}

	if o.R == nil {
		o.R = &categoryR{
			CategoryTypeRefUsers: related,
		}
	} else {
		o.R.CategoryTypeRefUsers = append(o.R.CategoryTypeRefUsers, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &userR{
				CategoryTypeRef: o,
			}
		} else {
			rel.R.CategoryTypeRef = o
		}
	}
	return nil
}

// Categories retrieves all the records using an executor.
func Categories(mods ...qm.QueryMod) categoryQuery {
	mods = append(mods, qm.From("`category`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`category`.*"})
	}

	return categoryQuery{q}
}

// FindCategory retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCategory(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Category, error) {
	categoryObj := &Category{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `category` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, categoryObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysqlmodel: unable to select from category")
	}

	return categoryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Category) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("mysqlmodel: no category provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(categoryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	categoryInsertCacheMut.RLock()
	cache, cached := categoryInsertCache[key]
	categoryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			categoryAllColumns,
			categoryColumnsWithDefault,
			categoryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(categoryType, categoryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(categoryType, categoryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `category` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `category` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `category` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, categoryPrimaryKeyColumns))
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
		return errors.Wrap(err, "mysqlmodel: unable to insert into category")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == categoryMapping["id"] {
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
		return errors.Wrap(err, "mysqlmodel: unable to populate default values for category")
	}

CacheNoHooks:
	if !cached {
		categoryInsertCacheMut.Lock()
		categoryInsertCache[key] = cache
		categoryInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Category.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Category) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	categoryUpdateCacheMut.RLock()
	cache, cached := categoryUpdateCache[key]
	categoryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			categoryAllColumns,
			categoryPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("mysqlmodel: unable to update category, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `category` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, categoryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(categoryType, categoryMapping, append(wl, categoryPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "mysqlmodel: unable to update category row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to get rows affected by update for category")
	}

	if !cached {
		categoryUpdateCacheMut.Lock()
		categoryUpdateCache[key] = cache
		categoryUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q categoryQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to update all for category")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to retrieve rows affected for category")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CategorySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), categoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `category` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, categoryPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to update all in category slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to retrieve rows affected all in update all category")
	}
	return rowsAff, nil
}

var mySQLCategoryUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Category) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("mysqlmodel: no category provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(categoryColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLCategoryUniqueColumns, o)

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

	categoryUpsertCacheMut.RLock()
	cache, cached := categoryUpsertCache[key]
	categoryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			categoryAllColumns,
			categoryColumnsWithDefault,
			categoryColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			categoryAllColumns,
			categoryPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("mysqlmodel: unable to upsert category, could not build update column list")
		}

		ret := strmangle.SetComplement(categoryAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`category`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `category` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(categoryType, categoryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(categoryType, categoryMapping, ret)
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
		return errors.Wrap(err, "mysqlmodel: unable to upsert for category")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == categoryMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(categoryType, categoryMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "mysqlmodel: unable to retrieve unique values for category")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "mysqlmodel: unable to populate default values for category")
	}

CacheNoHooks:
	if !cached {
		categoryUpsertCacheMut.Lock()
		categoryUpsertCache[key] = cache
		categoryUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Category record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Category) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("mysqlmodel: no Category provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), categoryPrimaryKeyMapping)
	sql := "DELETE FROM `category` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to delete from category")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to get rows affected by delete for category")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q categoryQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("mysqlmodel: no categoryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to delete all from category")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to get rows affected by deleteall for category")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CategorySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), categoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `category` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, categoryPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: unable to delete all from category slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysqlmodel: failed to get rows affected by deleteall for category")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Category) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindCategory(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CategorySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CategorySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), categoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `category`.* FROM `category` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, categoryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "mysqlmodel: unable to reload all in CategorySlice")
	}

	*o = slice

	return nil
}

// CategoryExists checks if the Category row exists.
func CategoryExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `category` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mysqlmodel: unable to check if category exists")
	}

	return exists, nil
}

// Exists checks if the Category row exists.
func (o *Category) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return CategoryExists(ctx, exec, o.ID)
}
