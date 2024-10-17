DROP TABLE IF EXISTS config_category;
DROP TABLE IF EXISTS config_key;
DROP TABLE IF EXISTS user_config;
DROP TABLE IF EXISTS user_config_value;
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS user_value;
DROP TABLE IF EXISTS user_preset;
DROP TABLE IF EXISTS cache_vehicle;


CREATE TABLE config_category (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key TEXT NOT NULL,
  name TEXT NOT NULL
);
INSERT INTO config_category (key, name) VALUES ("VEH_CLASS", "Vehicle Class");
INSERT INTO config_category (key, name) VALUES ("EV_TYPE", "Event Types");
INSERT INTO config_category (key, name) VALUES ("DIFF", "Difficulties");
INSERT INTO config_category (key, name) VALUES ("TIME", "Time and Weather");
INSERT INTO config_category (key, name) VALUES ("SERVER", "Server Config");


CREATE TABLE user_value (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  config_category_id INTEGER NOT NULL,
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  user_preset_id INTEGER
);


CREATE TABLE user_preset (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  config_category_id INTEGER NOT NULL,
  name TEXT NOT NULL
);


CREATE TABLE cache_vehicle (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key TEXT NOT NULL,
  name TEXT NOT NULL,
  brand TEXT,
  desc TEXT,
  tags TEXT,
  class TEXT,
  specs TEXT,
  torque TEXT,
  power TEXT,
  skins TEXT
)

