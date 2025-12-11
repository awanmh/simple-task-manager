-- Create Users Table
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    password varchar(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NOT NULL DEFAULT (now())
);

-- Create Tasks Table
CREATE TABLE IF NOT EXISTS tasks (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    title varchar(255) NOT NULL,
    description text,
    status varchar(50) NOT NULL DEFAULT 'pending',
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NOT NULL DEFAULT (now()),
    
    -- Foreign Key Constraint (One-to-Many)
    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Create Index for performance on filtering tasks by user
CREATE INDEX idx_tasks_user_id ON tasks(user_id);