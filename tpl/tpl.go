package tpl

func RepoTemplate() []byte {
	return []byte(`
package {{ .FileName }}

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"log"
    imp "{{ .EntityImport }}"
)

var ErrInvalidParam = errors.New("invalid parameter")

type {{ .EntityName }}Filter func(col string) string

type {{ .EntityName }}Repo interface {
	{{ .EntityName }}(ctx context.Context, id string) (*imp.{{ .EntityName }}, error)
	Update(ctx context.Context, id string, modify interface{}) error
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, entity *imp.{{ .EntityName }}) error
	Query(ctx context.Context, filters ...{{ .EntityName }}Filter) ([]imp.{{ .EntityName }}, error)
	ByPage(ctx context.Context, page, limit int64, filters ...{{ .EntityName }}Filter) ([]imp.{{ .EntityName }}, error)
}

type repo struct {
	db             driver.Database
	collectionName string
}

func New{{ .EntityName }}Repo(db driver.Database) {{ .EntityName }}Repo {
	return &repo{db, "{{ .CollectionName }}"}
}

func (r repo) {{ .EntityName }}(ctx context.Context, id string) (*imp.{{ .EntityName }}, error) {
	col, err := r.db.Collection(ctx, r.collectionName)
	if err != nil {
		return nil, err
	}

	var entity imp.{{ .EntityName }}
	if _, err = col.ReadDocument(ctx, id, &entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r repo) Update(ctx context.Context, id string, modify interface{}) error {
	col, err := r.db.Collection(ctx, r.collectionName)
	if err != nil {
		return err
	}

	if _, err := col.UpdateDocument(ctx, id, removeKey(modify)); err != nil {
		return err
	}

	return nil
}

func (r repo) Delete(ctx context.Context, id string) error {
	col, err := r.db.Collection(ctx, r.collectionName)
	if err != nil {
		return err
	}
	if _, err := col.RemoveDocument(ctx, id); err != nil {
		return err
	}

	return nil
}

func (r repo) Create(ctx context.Context, entity *imp.{{ .EntityName }}) error {
	col, err := r.db.Collection(ctx, r.collectionName)
	if err != nil {
		return err
	}
	if _, err := col.CreateDocument(ctx, removeKey(entity)); err != nil {
		return err
	}

	return nil
}

func (r repo) Query(ctx context.Context, filters ...{{ .EntityName }}Filter) ([]imp.{{ .EntityName }}, error) {
	var entities []imp.{{ .EntityName }}

	var f string
	for _, filterFunc := range filters {
		f = fmt.Sprintf("%s\n%s", f, filterFunc("e"))
	}

	query := fmt.Sprintf(` + "`" + `
	For e In %s
	%s
	Return e
	` + "`" + `, r.collectionName, f)

	cursor, err := r.db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("{{ .EntityName }}Repo query failed: %v", err)
	}
	defer func(cursor driver.Cursor) {
		if err := cursor.Close(); err != nil {
			log.Fatal(err)
		}
	}(cursor)

	for {
		var e imp.{{ .EntityName }}

		_, err := cursor.ReadDocument(ctx, &e)
		if driver.IsNotFound(err) {
			return nil, nil
		} else if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("{{ .EntityName }}Repo query failed: %v", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (r repo) ByPage(ctx context.Context, page, limit int64, filters ...{{ .EntityName }}Filter) ([]imp.{{ .EntityName }}, error) {
	var entities []imp.{{ .EntityName }}

	var f string
	for _, filterFunc := range filters {
		f = fmt.Sprintf("%s\n%s", f, filterFunc("e"))
	}

	if page < 1 {
		return nil, fmt.Errorf("page %w", ErrInvalidParam)
	}

	start := limit * (page - 1)

	query := fmt.Sprintf(` + "`" + `
	For e In %s
		%s
		SORT e._key
		LIMIT %d, %d
	Return e
	` + "`" + `, r.collectionName, f, start, limit)

	cursor, err := r.db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("{{ .EntityName }}Repo query failed: %v", err)
	}
	defer func(cursor driver.Cursor) {
		if err := cursor.Close(); err != nil {
			log.Fatal(err)
		}
	}(cursor)

	for {
		var e imp.{{ .EntityName }}

		_, err := cursor.ReadDocument(ctx, &e)
		if driver.IsNotFound(err) {
			return nil, nil
		} else if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("{{ .EntityName }}Repo query failed: %v", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func removeKey(d interface{}) map[string]interface{} {
	j, _ := json.Marshal(d)
	var res map[string]interface{}
	_ = json.Unmarshal(j, &res)

	delete(res, "_key")

	return res
}

func Filter(property, compare, value string) {{ .EntityName }}Filter {
	return func(col string) string {
		return fmt.Sprintf("FILTER %s.%s %s %s", col, property, compare, value)
	}
}
`)
}
