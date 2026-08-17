CREATE SCHEMA IF NOT EXISTS users;

create table users.users(
	id uuid primary key,
	username text unique not null,
	password_hash text,
	name text not null,
	surname text not null,
	score integer,
	created_at timestamp not null default current_timestamp
);

CREATE SCHEMA IF NOT EXISTS posts;

create table posts.posts(
	id uuid primary key,
	author_id uuid not null,
	text text not null,
	created_at timestamp not null default current_timestamp
);

CREATE SCHEMA IF NOT EXISTS comments;

create table comments.comments(
	id uuid primary key,
	post_id uuid not null,
	author_id uuid not null,
	text text not null,
	created_at timestamp not null default current_timestamp
);

CREATE SCHEMA IF NOT EXISTS likes;

CREATE TABLE likes.likes (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    liker_id UUID NOT NULL,
    CONSTRAINT unique_post_liker UNIQUE (post_id, liker_id)
);

GRANT ALL PRIVILEGES ON SCHEMA users TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA posts TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA comments TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA likes TO postgres;