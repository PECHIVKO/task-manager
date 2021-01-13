package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/PECHIVKO/task-manager/models"
)

const (
	insertColumnQuery = `insert into columns (column_name, project_id, position)
								  values ($1, $2, (select count(*)
								  from columns where project_id = $2));`
	moveLeftQuery = `update columns
					 set position = position + 1
					 where project_id = $3
					 	and position >= $2
						and position < $1;`
	moveRightQuery = `update columns
					  set position = position - 1
					  where project_id = $3
						and position <= $2
						and position > $1;`
	checkForUniqueColumnName    = "select not exists (select 1 from columns where project_id = $1 and column_name = $2);"
	deleteColumnQuery           = "delete from columns where column_id = $1;"
	getProjectColumnsCountQuery = "select count(*) from columns where project_id = $1;"
	getColumnQuery              = "select * from columns where column_id = $1;"
	fetchColumnsQuery           = "select * from columns where project_id = $1 order by position;"
	updateColumnNameQuery       = "update columns set column_name = $1 where column_id = $2;"
	moveToPositionQuery         = "update columns set position = $2 where column_id = $1;"
	getCurrentPositionQuery     = "select position from columns where column_id = $1;"
	getCurrentProjectIDQuery    = "select project_id from columns where column_id = $1;"
)

type Column struct {
	ID       int    `json:"column_id"`
	Project  int    `json:"project_id"`
	Position int    `json:"position"`
	Name     string `json:"column_name"`
}

type ColumnRepository struct {
	DB *sql.DB
}

func NewColumnRepository(dbConn *sql.DB) *ColumnRepository {
	var repo ColumnRepository
	repo.DB = dbConn
	return &repo
}

func isUniqueColumnName(tx *sql.Tx, ctx context.Context, projectID int, name string) (isUnique bool, err error) {
	err = tx.QueryRow(checkForUniqueColumnName, projectID, name).Scan(&isUnique)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return isUnique, err
}

func (r ColumnRepository) CreateColumn(ctx context.Context, c *models.Column) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	isUnique, err := isUniqueColumnName(tx, ctx, c.Project, c.Name)
	if err != nil {
		return err
	}

	if !isUnique {
		err = fmt.Errorf("Column named %s is already exists in this project", c.Name)
		return err
	}

	model := toColumn(c)
	_, err = tx.ExecContext(ctx, insertColumnQuery, model.Name, model.Project)
	if err != nil {
		//wrap  everything with errors package!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (r ColumnRepository) GetColumn(ctx context.Context, id string) (*models.Column, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	column := new(Column)

	err = tx.QueryRowContext(ctx, getColumnQuery, id).Scan(&column.ID, &column.Project, &column.Position, &column.Name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return toModel(column), err
}

func (r ColumnRepository) FetchColumns(ctx context.Context, projectID string) ([]*models.Column, error) {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, fetchColumnsQuery, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]*Column, 0)

	for rows.Next() {
		column := new(Column)
		err := rows.Scan(&column.ID, &column.Project, &column.Position, &column.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		columns = append(columns, column)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return toModels(columns), err
}

func getProjectColumnsCount(tx *sql.Tx, ctx context.Context, id int) (count int, err error) {

	projectID, err := getProjectID(tx, ctx, id)
	if err != nil {
		return -1, err
	}

	err = tx.QueryRowContext(ctx, getProjectColumnsCountQuery, projectID).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, err
}

func (r ColumnRepository) DeleteColumn(ctx context.Context, id int) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	projectColumnsCount, err := getProjectColumnsCount(tx, ctx, id)
	if err != nil {
		return err
	}
	if projectColumnsCount > 1 {
		// ADD tasks transporation
		err = r.MoveColumnToPosition(ctx, id, projectColumnsCount-1)
		if err != nil {
			return err
		}
		result, err := tx.ExecContext(ctx, deleteColumnQuery, id)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		log.Println("deleted rows: ", rows)
	} else {
		err = fmt.Errorf("Can not delete last column")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

func (r ColumnRepository) UpdateColumnName(ctx context.Context, c *models.Column) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateColumnNameQuery, c.Name, c.ID)

	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Println("updated rows: ", rows)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func getProjectID(tx *sql.Tx, ctx context.Context, id int) (projectID int, err error) {
	err = tx.QueryRowContext(ctx, getCurrentProjectIDQuery, id).Scan(&projectID)
	if err != nil {
		return -1, err
	}
	return projectID, err
}

func (r ColumnRepository) MoveColumnToPosition(ctx context.Context, id, position int) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	projectColumnsCount, err := getProjectColumnsCount(tx, ctx, id)
	if err != nil {
		return err
	}
	if position >= projectColumnsCount {
		err = fmt.Errorf("columnrepo: Can not move further then last position %d", projectColumnsCount)
		return err
	}

	var currentPosition int
	err = tx.QueryRowContext(ctx, getCurrentPositionQuery, id).Scan(&currentPosition)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(tx, ctx, id)
	if err != nil {
		return err
	}

	var displaceColumnsQuery string
	switch {
	case currentPosition > position:
		displaceColumnsQuery = moveLeftQuery
	case currentPosition < position:
		displaceColumnsQuery = moveRightQuery
	default:
		err = fmt.Errorf("Column is already on wanted position")
		return err
	}

	result, err := tx.ExecContext(ctx, displaceColumnsQuery, currentPosition, position, projectID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Println("displaced rows: ", rows)

	result, err = tx.ExecContext(ctx, moveToPositionQuery, id, position)
	if err != nil {
		return err
	}
	rows, err = result.RowsAffected()
	if err != nil {
		return err
	}
	log.Println("moved rows: ", rows)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func toColumn(c *models.Column) *Column {
	return &Column{
		ID:       c.ID,
		Name:     c.Name,
		Project:  c.Project,
		Position: c.Position,
	}
}

func toModel(c *Column) *models.Column {
	return &models.Column{
		ID:       c.ID,
		Name:     c.Name,
		Project:  c.Project,
		Position: c.Position,
	}
}

func toModels(cs []*Column) []*models.Column {
	out := make([]*models.Column, len(cs))

	for i, c := range cs {
		out[i] = toModel(c)
	}
	return out
}
