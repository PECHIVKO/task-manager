CREATE TABLE users (
   name VARCHAR(20) NOT NULL UNIQUE PRIMARY KEY
);

CREATE TABLE projects (
  project_id serial PRIMARY KEY,
  project_name VARCHAR(500) NOT NULL UNIQUE,
  project_description VARCHAR(1000)
);

CREATE TABLE columns (
  column_id serial PRIMARY KEY,
  project_id integer NOT NULL REFERENCES projects,
  position integer NOT NULL,
  column_name VARCHAR(255) NOT NULL
);

CREATE TABLE tasks (
  task_id serial PRIMARY KEY,
  column_id integer NOT NULL REFERENCES columns,
  priority serial,
  task_name VARCHAR(500) NOT NULL UNIQUE,
  task_description VARCHAR(5000) NOT NULL,
  username VARCHAR(20) NOT NULL
);

CREATE TABLE comments (
  comment_id serial PRIMARY KEY,
  task_id integer NOT NULL REFERENCES tasks,
  creation_date timestamp NOT NULL DEFAULT NOW(),
  comment VARCHAR(5000) NOT NULL,
  username VARCHAR(20) NOT NULL
);

