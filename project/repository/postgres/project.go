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
	getProjectQuery    = "select * from projects where project_id = $1"
	fetchProjectsQuery = "select * from projects order by project_name"
	deleteProjectQuery = "delete from projects where project_id = $1"
	deleteColumnsQuery = "delete from columns where project_id = $1"
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

func (r ProjectRepository) CreateProject(ctx context.Context, p *models.Project) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	model := toProject(p)

	var projectID int

	err = tx.QueryRowContext(ctx, insertProjectQuery, model.Name, model.Description).Scan(&projectID)

	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, initColumnQuery, projectID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Println("created columns: ", rows)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (r ProjectRepository) GetProject(ctx context.Context, id string) (*models.Project, error) {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	project := new(Project)

	err = tx.QueryRowContext(ctx, getProjectQuery, id).Scan(&project.ID, &project.Name, &project.Description)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return toModel(project), nil
}

func (r ProjectRepository) FetchProjects(ctx context.Context) ([]*models.Project, error) {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, fetchProjectsQuery)
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

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return toModels(projects), nil
}

// TODO delete coments tasks and columns
func (r ProjectRepository) DeleteProject(ctx context.Context, id string) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// deleteTasksQuery := "delete from tasks where column_id in ( select column_id from columns where project_id = $1"
	result, err := tx.ExecContext(ctx, deleteColumnsQuery, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Println("deleted columns: ", rows)

	result, err = tx.ExecContext(ctx, deleteProjectQuery, id)
	if err != nil {
		return err
	}
	rows, err = result.RowsAffected()
	if err != nil {
		return err
	}
	log.Println("deleted projects: ", rows)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (r ProjectRepository) UpdateProject(ctx context.Context, p *models.Project) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateProjectQuery, p.Name, p.Description, p.ID)

	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("result err: ", err)
	}
	log.Println("updated projects: ", rows)

	err = tx.Commit()
	if err != nil {
		return err
	}

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
