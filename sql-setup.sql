CREATE TABLE IF NOT EXISTS users (

    username VARCHAR(32) NOT NULL,
    password VARCHAR(32) NOT NULL,
    id INTEGER AUTO_INCREMENT,

    CONSTRAINT pk_id PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS classes (
    player_id INTEGER NOT NULL,
    class_id INTEGER NOT NULL,
    title VARCHAR(32) NOT NULL,
    icon VARCHAR(64),

    CONSTRAINT pk_class_player_id PRIMARY KEY (player_id , class_id),
    CONSTRAINT fk_player_id FOREIGN KEY (player_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS homework (
    player_id INTEGER NOT NULL,
    class_id INTEGER NOT NULL,
    homework_id INTEGER AUTO_INCREMENT,
    description TEXT,
    due_date TIMESTAMP,

    CONSTRAINT fk_player_id FOREIGN KEY (player_id) REFERENCES users(player_id),
    CONSTRAINT pk_homework_id PRIMARY KEY (homework_id)
);