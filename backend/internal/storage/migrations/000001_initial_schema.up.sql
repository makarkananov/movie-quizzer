-- ============================================
-- USERS
-- ============================================
CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    nickname      TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================
-- SESSIONS
-- ============================================
CREATE TABLE sessions
(
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    mode            TEXT        NOT NULL,                       -- frame | video | quote
    status          TEXT        NOT NULL DEFAULT 'in_progress', -- finished
    total_questions INT         NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at     TIMESTAMPTZ
);

CREATE INDEX idx_sessions_user ON sessions (user_id);

-- ============================================
-- QUESTIONS BANK
-- ============================================
CREATE TABLE questions
(
    id             BIGSERIAL PRIMARY KEY,
    type           TEXT NOT NULL, -- frame | video | quote
    text           TEXT,
    image_url      TEXT,
    video_url      TEXT,
    options        TEXT[],        -- массив вариантов
    correct_answer TEXT NOT NULL
);

CREATE INDEX idx_questions_type ON questions (type);

-- ============================================
-- SESSION QUESTIONS (порядок вопросов в раунде)
-- ============================================
CREATE TABLE session_questions
(
    id          BIGSERIAL PRIMARY KEY,
    session_id  BIGINT NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL REFERENCES questions (id),
    position    INT    NOT NULL
);

CREATE UNIQUE INDEX session_question_unique
    ON session_questions (session_id, position);

-- ============================================
-- ANSWERS (ответы пользователя на вопросы)
-- ============================================
CREATE TABLE answers
(
    id          BIGSERIAL PRIMARY KEY,
    session_id  BIGINT      NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
    question_id BIGINT      NOT NULL,
    user_answer TEXT        NOT NULL,
    correct     BOOLEAN     NOT NULL,
    score_delta INT         NOT NULL DEFAULT 0,
    elapsed_ms  BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_answers_session ON answers (session_id);

-- ============================================
-- ACHIEVEMENTS
-- ============================================
CREATE TABLE achievements
(
    id          BIGSERIAL PRIMARY KEY,
    code        TEXT UNIQUE NOT NULL,
    title       TEXT        NOT NULL,
    description TEXT        NOT NULL
);

CREATE TABLE user_achievements
(
    user_id        BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    achievement_id BIGINT      NOT NULL REFERENCES achievements (id),
    earned_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, achievement_id)
);
