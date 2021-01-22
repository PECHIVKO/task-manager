# task-manager

**HOW TO RUN**

use `docker-compose build` & `docker-compose up` commands  
*OR*  
use commands:  
`docker run --name postgresqldb -e POSTGRES_USER=manager -e POSTGRES_PASSWORD=task -p 5432:5432 -v /data:/var/lib/postgresql/data -d postgres`
`docker exec -i postgresqldb psql -U manager -c "CREATE DATABASE taskmanager WITH ENCODING='UTF8' OWNER=manager;"`  
to create database and  
`cd cmd/api`  
`go run main.go`  
to run lockaly  

**HOW TO USE**

*Projects:*  

1. Create:  
`curl localhost:8181/projects -d '{"project_name":"new project name", "project_description":"new project description"}'`
2. Update:  
`curl -X PUT  localhost:8181/projects/{project_id} -d '{"project_name":"updated name", "project_description":"updated description"}'`
3. Get:  
`curl localhost:8181/projects/{project_id}`
4. Get all projects:  
`curl localhost:8181/projects`
5. Delete:  
`curl -X DELETE localhost:8181/projects/{project_id}`

*Columns:*  

1. Create:  
`curl localhost:8181/columns -d '{"project_id":{existing_project_id}, "column_name":"new column name"}'`
2. Update column name:  
`curl -X PUT  localhost:8181/columns/{column_id} -d '{"column_name":"updated column name"}'`
3. Get:  
`curl localhost:8181/columns/{column_id}`
4. Get all columns for project:  
`curl localhost:8181/columns/project/{project_id}`
5. Move column to position:  
`curl -X PUT  localhost:8181/columns/{column_id}/move/{new_position}`
6. Delete:  
`curl -X DELETE localhost:8181/columns/{column_id}`

*Tasks:*  

1. Create:  
`curl localhost:8181/tasks -d '{"column_id":{existing_column_id}, "task_name":"new task name", "task_description":"new task description"}'`
2. Update:  
`curl -X PUT  localhost:8181/tasks/{task_id} -d '{"task_name":"updated task name", "task_description":"updated task description"}'`
3. Get:  
`curl localhost:8181/tasks/{task_id}`
4. Get all tasks for column:  
`curl localhost:8181/tasks/column/{column_id}`
5. Change task priority:  
`curl -X PUT  localhost:8181/tasks/{task_id}/priority/{new_priority}`
6. Move task to column:  
`curl -X PUT  localhost:8181/tasks/{task_id}/move/{column_id}`
7. Delete:  
`curl -X DELETE localhost:8181/tasks/{task_id}`

*Comments:*  

1. Create:  
`curl localhost:8181/comments -d '{"comment":"new comment", "task_id":{existing_task_id}}'`
2. Update:  
`curl -X PUT  localhost:8181/comments/{comment_id} -d '{"comment":"updated comment"}'`
3. Get:  
`curl localhost:8181/comments/{comment_id}`
4. Get all comments for task:  
`curl localhost:8181/comments/task/{task_id}`
5. Delete:  
`curl -X DELETE localhost:8181/comments/{comment_id}`