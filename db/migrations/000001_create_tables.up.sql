CREATE TABLE IF NOT EXISTS projects (
  project_id serial PRIMARY KEY,
  project_name VARCHAR(500) NOT NULL UNIQUE,
  project_description VARCHAR(1000)
);

CREATE TABLE IF NOT EXISTS columns (
  column_id serial PRIMARY KEY,
  project_id integer NOT NULL REFERENCES projects,
  position integer NOT NULL CHECK (position >= 0),
  column_name VARCHAR(255) NOT NULL CHECK (column_name <> '')
);

CREATE TABLE IF NOT EXISTS tasks (
  task_id serial PRIMARY KEY,
  column_id integer NOT NULL REFERENCES columns,
  priority integer NOT NULL CHECK (priority >= 0),
  task_name VARCHAR(500) NOT NULL CHECK (task_name <> ''),
  task_description VARCHAR(5000) NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
  comment_id serial PRIMARY KEY,
  task_id integer NOT NULL REFERENCES tasks,
  creation_date timestamp NOT NULL DEFAULT NOW(),
  comment VARCHAR(5000) NOT NULL CHECK (comment <> '')
);