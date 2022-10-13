package restcore

import "reflect"

type ModelSet struct {
	original   interface{}
	fieldNames map[string]string
}

func NewModelSet(input interface{}) *ModelSet {
	m := &ModelSet{
		original:   input,
		fieldNames: make(map[string]string),
	}

	m.resolveNames()

	return m
}

func (m *ModelSet) HasDbName(name string) bool {
	for _, dbName := range m.fieldNames {
		if dbName == name {
			return true
		}
	}

	return false
}

func (m *ModelSet) GetNames() ([]string, []string) {
	names := make([]string, 0, len(m.fieldNames))
	dbNames := make([]string, 0, len(m.fieldNames))

	for name, dbName := range m.fieldNames {
		names = append(names, name)
		dbNames = append(dbNames, dbName)
	}

	return names, dbNames
}

func (m *ModelSet) GetValues() map[string]interface{} {
	dbValues := make(map[string]interface{})

	pointer := reflect.Indirect(reflect.ValueOf(m.original))

	for structName, dbName := range m.fieldNames {
		field := pointer.FieldByName(structName)

		if field.IsZero() {
			continue
		}

		dbValues[dbName] = field.Interface()
	}

	return dbValues
}

func (m *ModelSet) SetStringValue(dbFieldName string, value string) {
	pointer := reflect.Indirect(reflect.ValueOf(m.original))

	for structName, dbName := range m.fieldNames {
		field := pointer.FieldByName(structName)

		if dbName == dbFieldName {
			field.SetString(value)

			return
		}
	}
}

func (m *ModelSet) resolveNames() {
	el := reflect.TypeOf(m.original).Elem()

	el.FieldByNameFunc(func(name string) bool {
		field, _ := el.FieldByName(name)

		dbName := field.Tag.Get("db")

		if dbName == "" || dbName == "-" {
			return false
		}

		m.fieldNames[name] = dbName

		return false
	})
}
