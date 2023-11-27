CREATE TABLE follow_relationships (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    followee_id uuid NOT NULL UNIQUE,
    follower_id uuid NOT NULL
);

CREATE TABLE post_likes (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    post_id uuid NOT NULL UNIQUE,
    user_id uuid NOT NULL,
    CONSTRAINT FK_post_likes_post_id FOREIGN KEY (post_id) REFERENCES posts (id),
    CONSTRAINT FK_post_likes_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);