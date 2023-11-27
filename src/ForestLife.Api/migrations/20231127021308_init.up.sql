ALTER TABLE follow_relationships ADD CONSTRAINT fk_follow_relationships_followee_id
FOREIGN KEY (followee_id)
REFERENCES users(id);

ALTER TABLE follow_relationships ADD CONSTRAINT fk_follow_relationships_follower_id
FOREIGN KEY (follower_id)
REFERENCES users(id)