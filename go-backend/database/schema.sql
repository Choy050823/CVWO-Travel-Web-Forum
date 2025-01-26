-- Schema for User Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    bio TEXT,
    profile_picture TEXT,
    role VARCHAR(50) NOT NULL,
    participation_score INTEGER DEFAULT 0
);

-- Schema for Category Table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Schema for Thread Table
CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    attached_images TEXT[],
    posted_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    likes INTEGER DEFAULT 0,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Schema for Comment Table
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    attached_images TEXT[],
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    thread_id INTEGER NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    upvotes INTEGER DEFAULT 0,
    downvotes INTEGER DEFAULT 0
);

-- Schema for Notification Table
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
