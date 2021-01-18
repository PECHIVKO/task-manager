package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/PECHIVKO/task-manager/models"
)

const (
	insertCommentQuery = `insert into comments (task_id, text)
								  values ($1, $2);`
	checkForTaskExists = `select exists
						(select 1 from tasks where task_id = $1 and project_id =
							(select project_id from tasks where task_id =
								(select task_id from comments where comment_id = $1)));`
	checkForCommentExists = "select exists (select 1 from comments where comment_id = $1);"
	deleteCommentQuery    = "delete from comments where comment_id = $1;"
	getCommentQuery       = "select * from comments where comment_id = $1;"
	fetchCommentsQuery    = "select * from comments where task_id = $1 order by creation_date;"
	updateCommentQuery    = "update comments set comment = $1 where comment_id = $3;"
)

type Comment struct {
	ID      int       `json:"comment_id"`
	Task    int       `json:"task_id"`
	Date    time.Time `json:"creation_date"`
	Comment string    `json:"text"`
}

type CommentRepository struct {
	DB *sql.DB
}

func NewCommentRepository(dbConn *sql.DB) *CommentRepository {
	var repo CommentRepository
	repo.DB = dbConn
	return &repo
}

func isCommentExists(tx *sql.Tx, ctx context.Context, commentID int) (isExists bool, err error) {
	err = tx.QueryRow(checkForCommentExists, commentID).Scan(&isExists)
	if err != nil {
		err = fmt.Errorf("comment repository: isCommentExists() func error : %w", err)
		return false, err
	}
	return isExists, err
}

func isTaskExists(tx *sql.Tx, ctx context.Context, taskID int) (isExists bool, err error) {
	err = tx.QueryRow(checkForTaskExists, taskID).Scan(&isExists)
	if err != nil {
		err = fmt.Errorf("comment repository: isTaskExists() func error : %w", err)
		return false, err
	}
	return isExists, err
}

// func getTaskID(tx *sql.Tx, ctx context.Context, id int) (taskID int, err error) {
// 	err = tx.QueryRowContext(ctx, getCurrentTaskIDQuery, id).Scan(&taskID)
// 	if err != nil {
// 		err = fmt.Errorf("comment repository: getTaskID() func error : %w", err)
// 		return -1, err
// 	}
// 	return taskID, err
// }

// func getMaxPriorityForTask(tx *sql.Tx, ctx context.Context, taskID int) (priority int, err error) {
// 	err = tx.QueryRowContext(ctx, getMaxPriority, taskID).Scan(&priority)
// 	if err != nil {
// 		err = fmt.Errorf("comment repository: getMaxPriorityForTask() func error : %w", err)
// 		return -1, err
// 	}
// 	return priority, err
// }

// func getTaskCommentsCount(tx *sql.Tx, ctx context.Context, id int) (count int, err error) {

// 	taskID, err := getTaskID(tx, ctx, id)
// 	if err != nil {
// 		return -1, err
// 	}

// 	err = tx.QueryRowContext(ctx, getTaskCommentsCountQuery, taskID).Scan(&count)
// 	if err != nil {
// 		err = fmt.Errorf("comment repository: getTaskCommentsCount() func error : %w", err)
// 		return -1, err
// 	}
// 	return count, err
// }

func (r CommentRepository) CreateComment(ctx context.Context, c *models.Comment) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	model := toComment(c)

	result, err := tx.ExecContext(ctx, insertCommentQuery, model.Task, model.Comment)
	if err != nil {
		err = fmt.Errorf("comment repository: CreateComment: exec insertCommentQuery error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("comment repository: CreateComment: get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("created comments: ", rows)
	return err
}

func (r CommentRepository) GetComment(ctx context.Context, commentID int) (*models.Comment, error) {

	comment := new(Comment)

	err := r.DB.QueryRowContext(ctx, getCommentQuery, commentID).Scan(&comment.ID, &comment.Task, &comment.Date, &comment.Comment)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return toModel(comment), err
}

func (r CommentRepository) FetchComments(ctx context.Context, taskID int) ([]*models.Comment, error) {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	isExists, err := isTaskExists(tx, ctx, taskID)
	if err != nil {
		return nil, err
	}

	if !isExists {
		err = fmt.Errorf("Task ID %d does not exists", taskID)
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, fetchCommentsQuery, taskID)
	if err != nil {
		err = fmt.Errorf("comment repository: FetchComments: get RowsAffected error : %w", err)
		return nil, err
	}
	defer rows.Close()

	comments := make([]*Comment, 0)

	for rows.Next() {
		comment := new(Comment)
		err := rows.Scan(&comment.ID, &comment.Task, &comment.Date, &comment.Comment)
		if err != nil {
			err = fmt.Errorf("comment repository: FetchComments: rows.Scan() error : %w", err)
			log.Println(err)
			continue
		}
		comments = append(comments, comment)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return toModels(comments), err
}

func (r CommentRepository) DeleteComment(ctx context.Context, commentID int) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isCommentExists(tx, ctx, commentID)
	if err != nil {
		return err
	}

	if !isExists {
		err = fmt.Errorf("Comment ID %d does not exists", commentID)
		return err
	}

	result, err := tx.ExecContext(ctx, deleteCommentQuery, commentID)
	if err != nil {
		err = fmt.Errorf("comment repository: DeleteComment: exec query error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("comment repository: DeleteComment: get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("deleted comments: ", rows)
	return err
}

func (r CommentRepository) UpdateComment(ctx context.Context, c *models.Comment) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isExists, err := isCommentExists(tx, ctx, c.ID)
	if err != nil {
		return err
	}
	if !isExists {
		err = fmt.Errorf("Comment ID %d does not exists", c.ID)
		return err
	}

	result, err := tx.ExecContext(ctx, updateCommentQuery, c.Comment, c.ID)
	if err != nil {
		err = fmt.Errorf("comment repository: UpdateCommentName: exec query error : %w", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("comment repository: UpdateCommentName: get RowsAffected error : %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("updated comments: ", rows)
	return err
}

func toComment(c *models.Comment) *Comment {
	return &Comment{
		ID:      c.ID,
		Task:    c.Task,
		Date:    c.Date,
		Comment: c.Comment,
	}
}

func toModel(c *Comment) *models.Comment {
	return &models.Comment{
		ID:      c.ID,
		Task:    c.Task,
		Date:    c.Date,
		Comment: c.Comment,
	}
}

func toModels(cs []*Comment) []*models.Comment {
	out := make([]*models.Comment, len(cs))

	for i, c := range cs {
		out[i] = toModel(c)
	}
	return out
}
