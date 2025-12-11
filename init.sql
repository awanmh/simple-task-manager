-- 1. Tabel Users
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    password varchar(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NOT NULL DEFAULT (now())
);

-- 2. Tabel Tasks (Lengkap dengan fitur baru + recurring)
CREATE TABLE IF NOT EXISTS tasks (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    title varchar(255) NOT NULL,
    description text,
    status varchar(50) NOT NULL DEFAULT 'pending',

    -- Kolom tambahan Phase 1
    priority varchar(20) DEFAULT 'medium', -- 'low', 'medium', 'high'
    labels text[],                         -- kategori task
    reminder_time timestamptz,             -- waktu reminder

    -- FITUR BARU: Recurring
    recurrence_pattern varchar(20),        -- 'daily', 'weekly', 'monthly'
    next_run timestamptz,                  -- kapan task ini dijadwalkan ulang

    -- Metadata
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NOT NULL DEFAULT (now()),

    CONSTRAINT fk_user 
        FOREIGN KEY(user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

-- 3. Tabel Subtasks
CREATE TABLE IF NOT EXISTS subtasks (
    id bigserial PRIMARY KEY,
    task_id bigint NOT NULL,
    title varchar(255) NOT NULL,
    is_done boolean DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT (now()),

    CONSTRAINT fk_task 
        FOREIGN KEY(task_id) 
        REFERENCES tasks(id) 
        ON DELETE CASCADE
);
