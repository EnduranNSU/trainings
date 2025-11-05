-- Exercise table
CREATE TABLE exercise (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    description TEXT NOT NULL,
    href TEXT NOT NULL,
    tags BIGINT NOT NULL
);

CREATE INDEX idx_exercise_id ON exercise(id);
CREATE INDEX idx_exercise_tags ON exercise(tags);

-- Tag table
CREATE TABLE tag (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    type VARCHAR(255) CHECK (type IN('')) NOT NULL
);

CREATE INDEX idx_tag_id ON tag(id);
CREATE INDEX idx_tag_type ON tag(type);

-- Exercise to tag (many-to-many relationship)
CREATE TABLE exercise_to_tag (
    execrise_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    PRIMARY KEY (execrise_id, tag_id)
);

CREATE INDEX idx_exercise_to_tag_execrise_id ON exercise_to_tag(execrise_id);
CREATE INDEX idx_exercise_to_tag_tag_id ON exercise_to_tag(tag_id);
CREATE INDEX idx_exercise_to_tag_both ON exercise_to_tag(execrise_id, tag_id);

-- Training table
CREATE TABLE training (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    isDone BOOLEAN NOT NULL,
    planned DATE NOT NULL,
    done DATE NULL,
    total_time INTERVAL NULL,
    rating INTEGER NULL
);

CREATE INDEX idx_training_id ON training(id);
CREATE INDEX idx_training_user_id ON training(user_id);
CREATE INDEX idx_training_planned ON training(planned DESC);
CREATE INDEX idx_training_user_planned ON training(user_id, planned DESC);

-- Trained exercise table
CREATE TABLE trained_exercise (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    training_id BIGINT NOT NULL,
    exercise_id BIGINT NOT NULL,
    weight FLOAT(53) NULL,
    approaches BIGINT NULL,
    reps BIGINT NULL,
    time TIME(0) WITHOUT TIME ZONE NULL,
    notes TEXT NULL
);

CREATE INDEX idx_trained_exercise_id ON trained_exercise(id);
CREATE INDEX idx_trained_exercise_training_id ON trained_exercise(training_id);
CREATE INDEX idx_trained_exercise_exercise_id ON trained_exercise(exercise_id);

-- Foreign key constraints
ALTER TABLE exercise_to_tag 
    ADD CONSTRAINT exercise_to_tag_tag_id_foreign 
    FOREIGN KEY (tag_id) REFERENCES tag(id);

ALTER TABLE exercise_to_tag 
    ADD CONSTRAINT exercise_to_tag_execrise_id_foreign 
    FOREIGN KEY (execrise_id) REFERENCES exercise(id);

ALTER TABLE trained_exercise 
    ADD CONSTRAINT trained_exercise_exercise_id_foreign 
    FOREIGN KEY (exercise_id) REFERENCES exercise(id);

ALTER TABLE trained_exercise 
    ADD CONSTRAINT trained_exercise_training_id_foreign 
    FOREIGN KEY (training_id) REFERENCES training(id);