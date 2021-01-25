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
								  values ($1, $2, (select count(*) from columns where project_id = $2));`
	moveLeftQuery = `update columns
					 set position = position + 1
					 where project_id = (select project_id from columns where column_id = $1)
					 	and position >= $2
						and position < (select position from columns where column_id = $1);`
	moveRightQuery = `update columns
					  set position = position - 1
					  where project_id = (select project_id from columns where column_id = $1)
						and position <= $2
						and position > (select position from columns where column_id = $1);`
	deleteColumnQuery = "delete from columns where column_id = $1;"
	moveTasksQuery    = `update tasks set column_id = (
							select coalesce (
								(select column_id from columns where project_id =
										(select project_id from columns where column_id = $1)
									and position <
										(select position from columns where column_id = $1)
									order by position desc limit 1),
								(select column_id from columns where project_id =
										(select project_id from columns where column_id = $1)
									and position >
								(select position from columns where column_id = $1)
									order by position limit 1)))
						 where column_id =$1`
	checkForUniqueColumnName    = "select not exists (select 1 from columns where project_id = $1 and column_name = $2);"
	checkForColumnExists        = "select exists (select 1 from columns where column_id = $1);"
	checkForProjectExists       = "select exists (select 1 from projects where project_id = $1);"
	getProjectColumnsCountQuery = "select count(*) from columns where project_id = (select project_id from columns where column_id = $1);"
	getColumnQuery              = "select * from columns where column_id = $1;"
	fetchColumnsQuery           = "select * from columns where project_id = $1 order by position;"
	updateColumnNameQuery       = "update columns set column_name = $1 where column_id = $2;"
	moveToPositionQuery         = "update columns set position = $2 where column_id = $1;"
	getCurrentPositionQuery     = "select position from columns where column_id = $1;"
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
		err = fmt.Errorf("column repository: isUniqueColumnName() func error : %w", err)
		return false, err
	}
	return isUnique, err
}

func isColumnExists(tx *sql.Tx, ctx context.Context, columnID int) (isExists bool, err error) {
	err = tx.QueryRow(checkForColumnExists, columnID).Scan(&isExists)
	if err != nil {
		err = fmt.Errorf("column repository: isColumnExists() func error : %w", err)
		return false, err
	}
	return isExists, err
}

func isProjectExists(tx *sql.Tx, ctx context.Context, projectID int) (isExists bool, err error) {
	err = tx.QueryRow(checkForProjectExists, projectID).Scan(&isExists)
	if err != nil {
		err = fmt.Errorf("column repository: isProjectExists() func error : %w", err)
		return false, err
	}
	return isExists, err
}

func getProjectColumnsCount(tx *sql.Tx, ctx context.Context, columnID int) (count int, err error) {

	err = tx.QueryRowContext(ctx, getProjectColumnsCountQuery, columnID).Scan(&count)
	if err != nil {
		err = fmt.Errorf("column repository: getProjectColumnsCount() func error : %w", err)
		return -1, err
	}
	return count, err
}

func (r ColumnRepository) CreateColumn(ctx context.Context, c *models.Column) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
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

	result, err := tx.ExecContext(ctx, insertColumnQuery, model.Name, model.Project)
	if err != nil {
		err = fmt.Errorf("column repository: CreateColumn: exec insertColumnQuery error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("column repository: CreateColumn: get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("created columns: ", rows)
	return err
}

func (r ColumnRepository) GetColumn(ctx context.Context, columnID int) (*models.Column, error) {

	column := new(Column)

	err := r.DB.QueryRowContext(ctx, getColumnQuery, columnID).Scan(&column.ID, &column.Project, &column.Position, &column.Name)

	if err != nil {
		err = fmt.Errorf("column repository: GetColumn: exec getColumnQuery error : %w", err)
		return nil, err
	}

	return toModel(column), err
}

func (r ColumnRepository) FetchColumns(ctx context.Context, projectID int) ([]*models.Column, error) {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	isExists, err := isProjectExists(tx, ctx, projectID)
	if err != nil {
		return nil, err
	}

	if !isExists {
		err = fmt.Errorf("Project ID %d does not exists", projectID)
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, fetchColumnsQuery, projectID)
	if err != nil {
		err = fmt.Errorf("column repository: FetchColumns: get RowsAffected error : %w", err)
		return nil, err
	}
	defer rows.Close()

	columns := make([]*Column, 0)

	for rows.Next() {
		column := new(Column)
		err := rows.Scan(&column.ID, &column.Project, &column.Position, &column.Name)
		if err != nil {
			err = fmt.Errorf("column repository: FetchColumns: rows.Scan() error : %w", err)
			log.Println(err)
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

func (r ColumnRepository) DeleteColumn(ctx context.Context, columnID int) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isColumnExists(tx, ctx, columnID)
	if err != nil {
		return err
	}

	if !isExists {
		err = fmt.Errorf("Column ID %d does not exists", columnID)
		return err
	}

	projectColumnsCount, err := getProjectColumnsCount(tx, ctx, columnID)
	if err != nil {
		return err
	}

	var deletedColumns, movedTasks int64
	if projectColumnsCount > 1 {
		result, err := tx.ExecContext(ctx, moveTasksQuery, columnID)
		if err != nil {
			err = fmt.Errorf("column repository: DeleteColumn: exec moveTasksQuery error : %w", err)
			return err
		}
		movedTasks, err = result.RowsAffected()
		if err != nil {
			err = fmt.Errorf("column repository: DeleteColumn: get moveTasksQuery RowsAffected error : %w", err)
			return err
		}

		err = r.MoveColumnToPosition(ctx, columnID, projectColumnsCount-1)
		if err != nil {
			return err
		}

		result, err = tx.ExecContext(ctx, deleteColumnQuery, columnID)
		if err != nil {
			err = fmt.Errorf("column repository: DeleteColumn: exec deleteColumnQuery error : %w", err)
			return err
		}
		deletedColumns, err = result.RowsAffected()
		if err != nil {
			err = fmt.Errorf("column repository: DeleteColumn: get deleteColumnQuery RowsAffected error : %w", err)
			return err
		}
	} else {
		err = fmt.Errorf("Can not delete last column")
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("tasks moved: ", movedTasks)
	log.Println("columns deleted: ", deletedColumns)
	return err
}

func (r ColumnRepository) UpdateColumnName(ctx context.Context, c *models.Column) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isColumnExists(tx, ctx, c.ID)
	if err != nil {
		return err
	}
	if !isExists {
		err = fmt.Errorf("Column ID %d does not exists", c.ID)
		return err
	}

	result, err := tx.ExecContext(ctx, updateColumnNameQuery, c.Name, c.ID)
	if err != nil {
		err = fmt.Errorf("column repository: UpdateColumnName: exec query error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("column repository: UpdateColumnName: get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("updated columns: ", rows)
	return err
}

func (r ColumnRepository) MoveColumnToPosition(ctx context.Context, columnID, position int) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isColumnExists(tx, ctx, columnID)
	if err != nil {
		return err
	}

	if !isExists {
		err = fmt.Errorf("Column ID %d does not exists", columnID)
		return err
	}

	projectColumnsCount, err := getProjectColumnsCount(tx, ctx, columnID)
	if err != nil {
		return err
	}
	if position >= projectColumnsCount {
		err = fmt.Errorf("column repository: Can not move further then last position %d", projectColumnsCount)
		return err
	}

	var currentPosition int
	err = tx.QueryRowContext(ctx, getCurrentPositionQuery, columnID).Scan(&currentPosition)
	if err != nil {
		err = fmt.Errorf("column repository: MoveColumnToPosition: getCurrentPosition query error : %w", err)
		return err
	}

	var displacedRows, movedRows int64
	if currentPosition != position {
		var displaceColumnsQuery string
		switch {
		case currentPosition > position:
			displaceColumnsQuery = moveLeftQuery
		case currentPosition < position:
			displaceColumnsQuery = moveRightQuery
		}

		result, err := tx.ExecContext(ctx, displaceColumnsQuery, columnID, position)
		if err != nil {
			err = fmt.Errorf("column repository: MoveColumnToPosition: exec displaceColumnsQuery error : %w", err)
			return err
		}
		displacedRows, err = result.RowsAffected()
		if err != nil {
			err = fmt.Errorf("column repository: MoveColumnToPosition: get displacedRowsAffected error : %w", err)
			return err
		}

		result, err = tx.ExecContext(ctx, moveToPositionQuery, columnID, position)
		if err != nil {
			err = fmt.Errorf("column repository: MoveColumnToPosition: exec moveToPositionQuery error : %w", err)
			return err
		}
		movedRows, err = result.RowsAffected()
		if err != nil {
			err = fmt.Errorf("column repository: MoveColumnToPosition: get movedRowsAffected error : %w", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("displaced columns: ", displacedRows)
	log.Println("moved columns: ", movedRows)
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
