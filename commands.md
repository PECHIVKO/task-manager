*CAN BE USED FOR TESTS*

`curl localhost:8181/projects -d '{"project_name":"Task-manager", "project_description":"Backend api for task-manager app"}'`
`curl localhost:8181/projects -d '{"project_name":"Unroller", "project_description":"Unroller for MxN matrix"}'`
`curl localhost:8181/projects -d '{"project_name":"Encoder", "project_description":"Encode phrase to pig latin"}'`
`curl localhost:8181/projects -d '{"project_name":"go-koans", "project_description":"Meeting with golang"}'`

`curl -X PUT  localhost:8181/projects/3 -d '{"project_name":"Encoders", "project_description":"Pig-latin and alpha-numeric encoder"}'`

`curl localhost:8181/projects/3`
`curl localhost:8181/projects`

`curl localhost:8181/projects -d '{"project_name":"Anagram-finder", "project_description":"REST api for anagram-finder"}'`

`curl localhost:8181/columns -d '{"project_id":2, "column_name":"Done"}'`
`curl localhost:8181/columns -d '{"project_id":3, "column_name":"Done"}'`
`curl localhost:8181/columns -d '{"project_id":4, "column_name":"Done"}'`

`curl localhost:8181/columns -d '{"project_id":1, "column_name":"TODO 2021"}'`
`curl localhost:8181/columns -d '{"project_id":1, "column_name":"DONE 2020"}'`
`curl localhost:8181/columns -d '{"project_id":1, "column_name":"DONE 2021"}'`
`curl localhost:8181/columns -d '{"project_id":1, "column_name":"DONE 2020"}'` <!-- ERROR -->
`curl -X PUT  localhost:8181/columns/1 -d '{"column_name":"TODO 2020"}'`

`curl localhost:8181/columns -d '{"project_id":5, "column_name":"Done"}'`
`curl localhost:8181/columns -d '{"project_id":5, "column_name":"done sprint_1"}'`

`curl -X PUT  localhost:8181/columns/12 -d '{"column_name":"done sprint_2"}'`

`curl localhost:8181/columns/12`
`curl localhost:8181/columns/project/1`

`curl -X PUT  localhost:8181/columns/12/move/2`

`curl -X DELETE localhost:8181/columns/2`
`curl -X DELETE localhost:8181/columns/6` <!-- ERROR -->

`curl localhost:8181/tasks -d '{"column_id":5, "task_name":"make storage", "task_description":"create map storage"}'`
`curl localhost:8181/tasks -d '{"column_id":5, "task_name":"create router", "task_description":"create route usung http package"}'`
`curl localhost:8181/tasks -d '{"column_id":5, "task_name":"create handlers", "task_description":"create handlers"}'`
`curl localhost:8181/tasks -d '{"column_id":5, "task_name":"junk task", "task_description":"task to delete with comments"}'`

`curl -X PUT  localhost:8181/tasks/3 -d '{"task_name":"create handlers", "task_description":"create handlers and tests"}'`

`curl localhost:8181/tasks -d '{"column_id":1, "task_name":"create entities packages", "task_description":"create package for each entity"}'`
`curl localhost:8181/tasks -d '{"column_id":1, "task_name":"create db connection", "task_description":"create connection config and init db"}'`
`curl localhost:8181/tasks -d '{"column_id":1, "task_name":"create usecases", "task_description":"create usecases for each entity"}'`

`curl localhost:8181/tasks -d '{"column_id":9, "task_name":"create router", "task_description":"create subrouters for each entity"}'`
`curl localhost:8181/tasks -d '{"column_id":9, "task_name":"create repositories", "task_description":"create repositories for each entity"}'`
`curl localhost:8181/tasks -d '{"column_id":9, "task_name":"write tests", "task_description":"write unit tests"}'`

`curl -X PUT  localhost:8181/tasks/8 -d '{"task_name":"create routing", "task_description":"create subrouters and handlers for each entity"}'`

`curl localhost:8181/tasks/10`
`curl localhost:8181/tasks/column/5`

`curl -X PUT  localhost:8181/tasks/9/priority/0`
`curl -X PUT  localhost:8181/tasks/8/priority/1`

`curl -X PUT  localhost:8181/tasks/5/move/10`
`curl -X PUT  localhost:8181/tasks/6/move/10`
`curl -X PUT  localhost:8181/tasks/7/move/10`
`curl -X PUT  localhost:8181/tasks/8/move/11`
`curl -X PUT  localhost:8181/tasks/9/move/11`
`curl -X PUT  localhost:8181/tasks/10/move/5` <!-- ERROR -->

`curl -X PUT  localhost:8181/tasks/1/move/13`
`curl -X PUT  localhost:8181/tasks/2/move/13`
`curl -X PUT  localhost:8181/tasks/3/move/12`

`curl localhost:8181/comments -d '{"comment":"comment_1", "task_id":1}'`
`curl localhost:8181/comments -d '{"comment":"comment_2", "task_id":1}'`

`curl localhost:8181/comments -d '{"comment":"comment_1", "task_id":2}'`
`curl localhost:8181/comments -d '{"comment":"comment_2", "task_id":2}'`

`curl localhost:8181/comments -d '{"comment":"comment_1", "task_id":4}'`
`curl localhost:8181/comments -d '{"comment":"comment_2", "task_id":4}'`

`curl -X PUT  localhost:8181/comments/5 -d '{"comment":"updated comment_1"}'`

`curl localhost:8181/comments/5`
`curl localhost:8181/comments/task/4`

`curl -X DELETE localhost:8181/comments/2`

`curl -X DELETE localhost:8181/tasks/2`

`curl localhost:8181/tasks -d '{"column_id":1, "task_name":"task no one did", "task_description":"task left in 2020"}'`

`curl -X DELETE localhost:8181/columns/1`
`curl -X DELETE localhost:8181/columns/10`

`curl -X DELETE localhost:8181/columns/11`

`curl -X DELETE localhost:8181/projects/5`