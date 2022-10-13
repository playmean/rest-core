package restcore

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type Model struct {
	Database  *sqlx.DB `json:"-"`
	TableName string   `json:"-"`

	original interface{}
	set      *ModelSet
}

type ModelSortOptions struct {
	Field     string
	Direction string
}

type ModelLoadOptions struct {
	Where   map[string]interface{}
	Sort    *ModelSortOptions
	GroupBy string
}

func (m *Model) ModelInitialize(db *sqlx.DB, tableName string, input interface{}) {
	m.Database = db
	m.TableName = tableName

	m.original = input
	m.set = NewModelSet(input)
}

func (m *Model) ModelLoadByID(id string) error {
	err := m.Database.Get(m.original, fmt.Sprintf(`
		SELECT * FROM %s WHERE id = ?
	`, m.TableName), id)
	if err != nil {
		return NewApiError(&ApiErrorOptions{
			Code:     "SELECT",
			Subcode:  m.TableName,
			Message:  fmt.Sprintf("%s with id '%s' not found", m.TableName, id),
			Original: err,
		})
	}

	return nil
}

func (m *Model) ModelLoad(opts ModelLoadOptions) error {
	var sql = fmt.Sprintf("SELECT * FROM %s WHERE ", m.TableName)
	var values []interface{}

	for k, v := range opts.Where {
		sql += fmt.Sprintf("%s = ? AND ", k)
		values = append(values, v)
	}

	sql += "1"

	if opts.Sort != nil {
		sql += fmt.Sprintf(" ORDER BY %s %s", opts.Sort.Field, opts.Sort.Direction)
	}

	if len(opts.GroupBy) > 0 {
		sql += fmt.Sprintf(" GROUP BY %s", opts.GroupBy)
	}

	err := m.Database.Get(m.original, sql, values...)
	if err != nil {
		return NewApiError(&ApiErrorOptions{
			Code:     "SELECT",
			Subcode:  m.TableName,
			Message:  fmt.Sprintf("%s not found", m.TableName),
			Original: err,
		})
	}

	return nil
}

func (m *Model) ModelCreate() error {
	sqlQuery := fmt.Sprintf("INSERT INTO %s", m.TableName)

	dbNames, values := m.getValues()

	hasID := m.set.HasDbName("id")
	id := uuid.Must(uuid.NewV4()).String()

	if hasID {
		dbNames = append(dbNames, "id")
		values = append(values, id)
	}

	placeholders := strings.Repeat("?, ", len(dbNames)-1) + "?"

	sqlQuery += fmt.Sprintf("(%s) VALUES (%s)", strings.Join(dbNames, ", "), placeholders)

	_, err := m.Database.Exec(sqlQuery, values...)
	if err != nil {
		return NewApiError(&ApiErrorOptions{
			Code:     "INSERT",
			Subcode:  m.TableName,
			Message:  "cannot insert row",
			Original: err,
		})
	}

	if hasID {
		m.set.SetStringValue("id", id)
	}

	return nil
}

func (m *Model) ModelSave() error {
	sqlQuery := fmt.Sprintf("UPDATE %s SET ", m.TableName)

	assigns, values, idValue := m.getAssigns()

	values = append(values, idValue)

	sqlQuery += fmt.Sprintf("%s WHERE id = ?", strings.Join(assigns, ", "))

	_, err := m.Database.Exec(sqlQuery, values...)
	if err != nil {
		return NewApiError(&ApiErrorOptions{
			Code:     "UPDATE",
			Subcode:  m.TableName,
			Message:  fmt.Sprintf("cannot update row with id '%s'", idValue),
			Original: err,
		})
	}

	return nil
}

func (m *Model) getValues() ([]string, []interface{}) {
	dbNames := make([]string, 0)
	values := make([]interface{}, 0)

	for dbName, value := range m.set.GetValues() {
		dbNames = append(dbNames, dbName)
		values = append(values, value)
	}

	return dbNames, values
}

func (m *Model) getAssigns() ([]string, []interface{}, string) {
	assigns := make([]string, 0)
	values := make([]interface{}, 0)

	var idValue = ""

	for dbName, value := range m.set.GetValues() {
		if dbName == "id" {
			idValue = value.(string)

			continue
		}

		if dbName == "updated" {
			value = uint(time.Now().Unix())
		}

		assigns = append(assigns, fmt.Sprintf("%s = ?", dbName))
		values = append(values, value)
	}

	return assigns, values, idValue
}
