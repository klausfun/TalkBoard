CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE posts
(
    id                 serial                                      not null unique,
    user_id            int references users (id) on delete cascade not null,
    title              varchar(255)                                not null,
    content            varchar(8191)                               not null,
    access_to_comments bool                                        not null
);

CREATE TABLE comments
(
    id                serial                                      not null unique,
    parent_comment_id int references comments (id) on delete cascade,
    post_id           int references posts (id) on delete cascade not null,
    user_id           int references users (id) on delete cascade not null,
    content           varchar(2000)                               not null
);

CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id);

CREATE INDEX idx_comments_post_id ON comments(post_id);