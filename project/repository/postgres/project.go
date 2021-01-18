package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/PECHIVKO/task-manager/models"
)

const (
	insertProjectQuery = `insert into projects (project_name, project_description)
						  values ($1, $2) returning project_id;`
	initColumnQuery = `insert into columns (column_name, project_id, position)
					   values ('TODO', $1, 0);`
	updateProjectQuery = `update projects set project_name = $1, project_description = $2
							 where project_id = $3;`
	deleteTasksQuery = `delete from tasks where column_id in (
							select column_id from columns where project_id = $1)`
	deleteCommentsQuery = `delete from comments where task_id (
								select task_id from tasks where column_id in (
									select column_id from columns where project_id = $1))`
	deleteColumnsQuery    = "delete from columns where project_id = $1"
	checkForProjectExists = "select exists (select 1 from projects where project_id = $1);"
	getProjectQuery       = "select * from projects where project_id = $1"
	fetchProjectsQuery    = "select * from projects order by project_name"
	deleteProjectQuery    = "delete from projects where project_id = $1"
)

type Project struct {
	ID          int    `json:"project_id"`
	Name        string `json:"project_name"`
	Description string `json:"project_description"`
}

type ProjectRepository struct {
	DB *sql.DB
}

func NewProjectRepository(dbConn *sql.DB) *ProjectRepository {
	var repo ProjectRepository
	repo.DB = dbConn
	return &repo
}

func isProjectExists(tx *sql.Tx, ctx context.Context, projectID int) (isExists bool, err error) {
	err = tx.QueryRow(checkForProjectExists, projectID).Scan(&isExists)
	if err != nil {
		err = fmt.Errorf("project repository: isProjectExists() func error : %w", err)
		return false, err
	}
	return isExists, err
}

func (r ProjectRepository) CreateProject(ctx context.Context, p *models.Project) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	model := toProject(p)

	var projectID int

	err = tx.QueryRowContext(ctx, insertProjectQuery, model.Name, model.Description).Scan(&projectID)
	if err != nil {
		err = fmt.Errorf("project repository: CreateProject: exec insertProjectQuery error : %w", err)
		return err
	}

	result, err := tx.ExecContext(ctx, initColumnQuery, projectID)
	if err != nil {
		err = fmt.Errorf("project repository: CreateProject: exec initColumnQuery error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("project repository: CreateProject: initColumnQuery get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("created projects: 1")
	log.Println("created columns: ", rows)
	return err
}

func (r ProjectRepository) GetProject(ctx context.Context, projectID int) (*models.Project, error) {

	project := new(Project)

	err := r.DB.QueryRowContext(ctx, getProjectQuery, projectID).Scan(&project.ID, &project.Name, &project.Description)
	if err != nil {
		err = fmt.Errorf("project repository: GetProject: exec query error : %w", err)
		return nil, err
	}

	return toModel(project), nil
}

func (r ProjectRepository) FetchProjects(ctx context.Context) ([]*models.Project, error) {

	rows, err := r.DB.QueryContext(ctx, fetchProjectsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]*Project, 0)

	for rows.Next() {
		project := new(Project)
		err := rows.Scan(&project.ID, &project.Name, &project.Description)
		if err != nil {
			fmt.Println(err)
			continue
		}
		projects = append(projects, project)
	}

	return toModels(projects), nil
}

// TODO delete coments tasks and columns
func (r ProjectRepository) DeleteProject(ctx context.Context, projectID int) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isProjectExists(tx, ctx, projectID)
	if err != nil {
		return err
	}

	if !isExists {
		err = fmt.Errorf("Project ID %d is not exists", projectID)
		return err
	}

	result, err := tx.ExecContext(ctx, deleteCommentsQuery, projectID)
	if err != nil {
		err = fmt.Errorf("project repository: DeleteProject: exec deleteCommentsQuery error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("project repository: DeleteProject: deleteCommentsQuery get RowsAffected error : %w", err)
		return err
	}
	log.Println("deleted comments: ", rows)

	result, err = tx.ExecContext(ctx, deleteTasksQuery, projectID)
	if err != nil {
		err = fmt.Errorf("project repository: DeleteProject: exec deleteTasksQuery error : %w", err)
		return err
	}
	rows, err = result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("project repository: DeleteProject: deleteTasksQuery get RowsAffected error : %w", err)
		return err
	}
	log.Println("deleted tasks: ", rows)

	result, err = tx.ExecContext(ctx, deleteColumnsQuery, projectID)
	if err != nil {
		err = fmt.Errorf("project repository: DeleteProject: exec deleteColumnsQuery error : %w", err)
		return err
	}
	rows, err = result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("project repository: DeleteProject: deleteColumnsQuery get RowsAffected error : %w", err)
		return err
	}
	log.Println("deleted columns: ", rows)

	result, err = tx.ExecContext(ctx, deleteProjectQuery, projectID)
	if err != nil {
		err = fmt.Errorf("project repository: DeleteProject: exec deleteProjectQuery error : %w", err)
		return err
	}
	rows, err = result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("project repository: DeleteProject: deleteProjectQuery get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("deleted projects: ", rows)
	return err
}

func (r ProjectRepository) UpdateProject(ctx context.Context, p *models.Project) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isProjectExists(tx, ctx, p.ID)
	if err != nil {
		return err
	}

	if !isExists {
		err = fmt.Errorf("Project ID %d does not exists", p.ID)
		return err
	}

	result, err := tx.ExecContext(ctx, updateProjectQuery, p.Name, p.Description, p.ID)
	if err != nil {
		err = fmt.Errorf("project repository: UpdateProject: exec updateProjectQuery error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("project repository: UpdateProject:  get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("updated projects: ", rows)
	return err
}

func toProject(p *models.Project) *Project {
	return &Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}
}

func toModel(p *Project) *models.Project {
	return &models.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}
}

func toModels(ps []*Project) []*models.Project {
	out := make([]*models.Project, len(ps))

	for i, p := range ps {
		out[i] = toModel(p)
	}
	return out
}
