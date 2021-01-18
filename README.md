# task-manager

Projects:

1. Create
curl localhost:8181/projects -d '{"project_name":"new project name", "project_description":"new project description"}'

2. Update
curl -X PUT  localhost:8181/projects/{project_id} -d '{"project_name":"updated name", "project_description":"updated description"}'

3. Get project info
curl localhost:8181/projects/{project_id}

4. Get all projects info
curl localhost:8181/projects

5. Delete project
curl -X DELETE localhost:8181/projects/{project_id}

Columns:

1. Create
curl localhost:8181/columns -d '{"project_id":{existing_project_id}, "column_name":"new column name"}'

2. Update column name
curl -X PUT  localhost:8181/columns/{column_id} -d '{"column_name":"updated column name"}'

3. Get column info
curl localhost:8181/columns/{column_id}

4. Get all columns for project
curl localhost:8181/columns/project/{project_id}

5. Move column to position
curl -X PUT  localhost:8181/columns/{column_id}/move/{new_position}

6. Delete column
curl -X DELETE localhost:8181/columns/{column_id}
