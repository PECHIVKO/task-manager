package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/PECHIVKO/task-manager/models"
)

const (
	insertTaskQuery = `insert into tasks (task_name, task_description, column_id, priority)
								  values ($1, $2, $3, (
									select coalesce (
										(select max (priority) from tasks where column_id = $3)+1,
										0)));`
	checkColumnInThisProject = `select exists
								(select 1 from columns where column_id = $1 and project_id =
									(select project_id from columns where column_id =
										(select column_id from tasks where task_id = $2)));`
	checkForColumnExists    = "select exists (select 1 from columns where column_id = $1);"
	checkForTaskExists      = "select exists (select 1 from tasks where task_id = $1);"
	deleteTaskQuery         = "delete from tasks where task_id = $1;"
	deleteCommentsQuery     = "delete from comments where task_id =$1"
	getTaskQuery            = "select * from tasks where task_id = $1;"
	fetchTasksQuery         = "select * from tasks where column_id = $1 order by priority;"
	updateTaskPriorityQuery = "update tasks set priority = $1 where task_id = $2;"
	updateTaskQuery         = "update tasks set task_name = $1, task_description = $2 where task_id = $3;"
	moveToColumnQuery       = "update tasks set column_id = $2 where task_id = $1;"
)

type Task struct {
	ID          int    `json:"task_id"`
	Column      int    `json:"column_id"`
	Priority    int    `json:"priority"`
	Name        string `json:"task_name"`
	Description string `json:"description"`
}

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(dbConn *sql.DB) *TaskRepository {
	var repo TaskRepository
	repo.DB = dbConn
	return &repo
}

func isTaskExists(tx *sql.Tx, ctx context.Context, taskID int) (isExists bool, err error) {
	err = tx.QueryRow(checkForTaskExists, taskID).Scan(&isExists)
	if err != nil {
		err = fmt.Errorf("task repository: isTaskExists() func error : %w", err)
		return false, err
	}
	return isExists, err
}

func isColumnExists(tx *sql.Tx, ctx context.Context, columnID int) (isExists bool, err error) {
	err = tx.QueryRow(checkForColumnExists, columnID).Scan(&isExists)
	if err != nil {
		err = fmt.Errorf("task repository: isColumnExists() func error : %w", err)
		return false, err
	}
	return isExists, err
}

func isInThisProject(tx *sql.Tx, ctx context.Context, columnID, taskID int) (isExists bool, err error) {
	err = tx.QueryRow(checkColumnInThisProject, columnID, taskID).Scan(&isExists)
	if err != nil {
		err = fmt.Errorf("task repository: isInThisProject() func error : %w", err)
		return false, err
	}
	return isExists, err
}

func (r TaskRepository) CreateTask(ctx context.Context, t *models.Task) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	model := toTask(t)

	result, err := tx.ExecContext(ctx, insertTaskQuery, model.Name, model.Description, model.Column)
	if err != nil {
		err = fmt.Errorf("task repository: CreateTask: exec insertTaskQuery error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("task repository: CreateTask: get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("created tasks: ", rows)
	return err
}

func (r TaskRepository) GetTask(ctx context.Context, taskID int) (*models.Task, error) {

	task := new(Task)

	err := r.DB.QueryRowContext(ctx, getTaskQuery, taskID).Scan(&task.ID, &task.Column, &task.Priority, &task.Name, &task.Description)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return toModel(task), err
}

func (r TaskRepository) FetchTasks(ctx context.Context, columnID int) ([]*models.Task, error) {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	isExists, err := isColumnExists(tx, ctx, columnID)
	if err != nil {
		return nil, err
	}

	if !isExists {
		err = fmt.Errorf("Column ID %d does not exists", columnID)
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, fetchTasksQuery, columnID)
	if err != nil {
		err = fmt.Errorf("task repository: FetchTasks: get RowsAffected error : %w", err)
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		task := new(Task)
		err := rows.Scan(&task.ID, &task.Column, &task.Priority, &task.Name, &task.Description)
		if err != nil {
			err = fmt.Errorf("task repository: FetchTasks: rows.Scan() error : %w", err)
			log.Println(err)
			continue
		}
		tasks = append(tasks, task)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return toModels(tasks), err
}

func (r TaskRepository) DeleteTask(ctx context.Context, taskID int) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isTaskExists(tx, ctx, taskID)
	if err != nil {
		return err
	}

	if !isExists {
		err = fmt.Errorf("Task ID %d does not exists", taskID)
		return err
	}

	result, err := tx.ExecContext(ctx, deleteCommentsQuery, taskID)
	if err != nil {
		err = fmt.Errorf("task repository: DeleteTask: exec deleteCommentsQuery error : %w", err)
		return err
	}
	commentsDeleted, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("task repository: DeleteTask: get deleteCommentsQuery RowsAffected error : %w", err)
		return err
	}

	result, err = tx.ExecContext(ctx, deleteTaskQuery, taskID)
	if err != nil {
		err = fmt.Errorf("task repository: DeleteTask: exec deleteTaskQuery error : %w", err)
		return err
	}
	tasksDeleted, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("task repository: DeleteTask: get deleteTaskQuery RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Println("deleted comments: ", commentsDeleted)
	log.Println("deleted tasks: ", tasksDeleted)
	return err
}

func (r TaskRepository) UpdateTask(ctx context.Context, t *models.Task) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isTaskExists(tx, ctx, t.ID)
	if err != nil {
		return err
	}
	if !isExists {
		err = fmt.Errorf("Task ID %d does not exists", t.ID)
		return err
	}

	result, err := tx.ExecContext(ctx, updateTaskQuery, t.Name, t.Description, t.ID)
	if err != nil {
		err = fmt.Errorf("task repository: UpdateTaskName: exec query error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("task repository: UpdateTaskName: get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("updated tasks: ", rows)
	return err
}

func (r TaskRepository) ChangeTaskPriority(ctx context.Context, taskID, priority int) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isTaskExists(tx, ctx, taskID)
	if err != nil {
		return err
	}
	if !isExists {
		err = fmt.Errorf("Task ID %d does not exists", taskID)
		return err
	}

	result, err := tx.ExecContext(ctx, updateTaskPriorityQuery, priority, taskID)
	if err != nil {
		err = fmt.Errorf("task repository: ChangeTaskPriority: exec query error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("task repository: ChangeTaskPriority: get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("updated tasks: ", rows)
	return err
}

func (r TaskRepository) MoveToColumn(ctx context.Context, taskID, columnID int) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isTaskExists(tx, ctx, taskID)
	if err != nil {
		return err
	}
	if !isExists {
		err = fmt.Errorf("Task ID %d does not exists", taskID)
		return err
	}

	isExists, err = isInThisProject(tx, ctx, columnID, taskID)
	if err != nil {
		return err
	}
	if !isExists {
		err = fmt.Errorf("Column ID %d does not exists in this project", columnID)
		return err
	}

	result, err := tx.ExecContext(ctx, moveToColumnQuery, taskID, columnID)
	if err != nil {
		err = fmt.Errorf("task repository: MoveToColumn: exec moveToColumnQuery error : %w", err)
		return err
	}
	movedRows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("task repository: MoveToColumn: get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("moved tasks: ", movedRows)
	return err
}

func toTask(t *models.Task) *Task {
	return &Task{
		ID:          t.ID,
		Column:      t.Column,
		Priority:    t.Priority,
		Name:        t.Name,
		Description: t.Description,
	}
}

func toModel(t *Task) *models.Task {
	return &models.Task{
		ID:          t.ID,
		Column:      t.Column,
		Priority:    t.Priority,
		Name:        t.Name,
		Description: t.Description,
	}
}

func toModels(cs []*Task) []*models.Task {
	out := make([]*models.Task, len(cs))

	for i, t := range cs {
		out[i] = toModel(t)
	}
	return out
}
