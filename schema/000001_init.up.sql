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

CREATE TABLE subscriptions
(
    id                serial                                      not null unique,
    parent_comment_id int references subscriptions (id) on delete cascade,
    post_id           int references posts (id) on delete cascade not null,
    content           varchar(2000)                               not null
);
