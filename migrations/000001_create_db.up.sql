CREATE TABLE users (
    id bigserial NOT NULL PRIMARY KEY,
    email varchar NOT NULL UNIQUE,
    encrypted_password varchar NOT NULL
);

CREATE TABLE platforms (
    id bigserial NOT NULL PRIMARY KEY,
    name varchar NOT NULL,
    full_name varchar NOT NULL
);

CREATE TABLE profiles (
    id bigserial NOT NULL PRIMARY KEY,
    user_id int NOT NULL,
    platform_id int NOT NULL,
    name varchar NOT NULL,
    email varchar NOT NULL,
    access_token varchar NOT NULL,
    refresh_token varchar NOT NULL,
    resumes varchar[]
);

CREATE TABLE groups (
    id bigserial NOT NULL PRIMARY KEY,
    user_id int NOT NULL,
    name varchar NOT NULL,
    resume varchar NOT NULL,
    letter varchar NOT NULL
);

CREATE TABLE letters (
    id bigserial NOT NULL PRIMARY KEY,
    user_id int NOT NULL,
    name varchar NOT NULL,
    body varchar NOT NULL
);

CREATE TABLE tasks (
    id bigserial NOT NULL PRIMARY KEY,
    group_id int NOT NULL,
    name varchar NOT NULL
);

CREATE TABLE vacancies (
    id bigserial NOT NULL PRIMARY KEY,
    task_id int NOT NULL,
    number varchar NOT NULL UNIQUE,
    link varchar NOT NULL,
    name varchar NOT NULL,
    salary_from numeric,
    salary_to numeric,
    area varchar NOT NULL,
    company varchar NOT NULL,
    description varchar NOT NULL,
    at_published timestamp NOT NULL,
    responsed boolean NOT NULL,
    selected boolean NOT NULL
);
