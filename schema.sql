CREATE TABLE IF NOT EXISTS user_config (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  password TEXT,
  admin_password TEXT,
  register_to_lobby INTEGER,
  pickup_mode_enabled INTEGER,
  locked_entry_list INTEGER,
  result_screen_time INTEGER,
  udp_port INTEGER,
  tcp_port INTEGER,
  http_port INTEGER,
  client_send_interval INTEGER,
  num_threads INTEGER,
  max_clients INTEGER,
  welcome_message TEXT,

  install_path TEXT,
  csp_required INTEGER,
  csp_phycars INTEGER,
  csp_phytracks INTEGER,
  csp_hidepit INTEGER,
  csp_version INTEGER,

  cfg_filled INTEGER,
  mod_filled INTEGER
);

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  password TEXT,

  measurement_unit INTEGER,
  temp_unit INTEGER
);

CREATE TABLE IF NOT EXISTS user_difficulty (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  abs_allowed INTEGER,
  tc_allowed INTEGER,
  stability_allowed INTEGER,
  autoclutch_allowed INTEGER,
  tyre_blankets_allowed INTEGER,
  force_virtual_mirror INTEGER,
  fuel_rate INTEGER,
  damage_multiplier INTEGER,
  tyre_wear_rate INTEGER,
  allowed_tyres_out INTEGER,
  max_ballast_kg INTEGER,
  start_rule INTEGER,
  race_gas_penality_disabled INTEGER,
  dynamic_track INTEGER,
  dynamic_track_preset INTEGER,
  session_start INTEGER,
  randomness INTEGER,
  session_transfer INTEGER,
  lap_gain INTEGER,
  kick_quorum INTEGER,
  vote_duration INTEGER,
  voting_quorum INTEGER,
  blacklist_mode INTEGER,
  max_contacts_per_km INTEGER,

  filled INTEGER,
  deleted INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_time (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  time TEXT,
  time_of_day_multi INTEGER,

  filled INTEGER,
  deleted INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_time_weather (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_time_id INTEGER NOT NULL,
  graphics TEXT,
  base_temperature_ambient INTEGER,
  base_temperature_road INTEGER,
  variation_ambient INTEGER,
  variation_road INTEGER,
  wind_base_speed_min INTEGER,
  wind_base_speed_max INTEGER,
  wind_base_direction INTEGER,
  wind_variation_direction INTEGER,

  deleted INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_session (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  booking_enabled INTEGER,
  booking_time INTEGER,
  practice_enabled INTEGER,
  practice_time INTEGER,
  practice_is_open INTEGER,
  qualify_enabled INTEGER,
  qualify_time INTEGER,
  qualify_is_open INTEGER,
  qualify_max_wait_perc INTEGER,
  race_enabled INTEGER,
  race_time INTEGER,
  race_extra_lap INTEGER,
  race_over_time INTEGER,
  race_wait_time INTEGER,
  race_is_open INTEGER,
  reversed_grid_positions INTEGER,
  race_pit_window_start INTEGER,
  race_pit_window_end INTEGER,

  filled INTEGER,
  deleted INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_class (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,

  filled INTEGER,
  deleted INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_class_entry (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_class_id INTEGER NOT NULL,
  cache_car_key TEXT NOT NULL,
  skin_key TEXT NOT NULL,
  ballast INTEGER,

  deleted INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_event (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  cache_track_key TEXT NOT NULL,
  cache_track_config TEXT NOT NULL,
  difficulty_id INTEGER NOT NULL,
  session_id INTEGER NOT NULL,
  class_id INTEGER NOT NULL,
  time_id INTEGER NOT NULL,

  race_laps INTEGER,
  strategy INTEGER,

  started_at DATETIME,
  finished INTEGER,

  deleted INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS cache_track (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key TEXT NOT NULL,
  config TEXT NOT NULL,
  name TEXT NOT NULL,
  desc TEXT,
  tags TEXT,
  country TEXT,
  city TEXT,
  length INTEGER,
  width INTEGER,
  pitboxes INTEGER,
  run TEXT
);

CREATE TABLE IF NOT EXISTS cache_car (
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
);

CREATE TABLE IF NOT EXISTS cache_weather (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key TEXT NOT NULL,
  name TEXT NOT NULL
);


-- INDEXES

CREATE UNIQUE INDEX IF NOT EXISTS idx_cache_track_key_config
ON cache_track (key, config);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cache_car_key
ON cache_car (key);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cache_weather_key
ON cache_weather (key);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_name
ON users (name);


-- DEFAULT VALUES
INSERT OR IGNORE INTO user_config (id, name, udp_port, tcp_port, http_port, client_send_interval, num_threads) VALUES (1, 'AC Server', 9600, 9600, 8081, 18, 2);

-- DEFAULT USERNAME admin PASSWORD admin
INSERT OR IGNORE INTO users (id, name, password) VALUES (1, 'admin', '$2a$08$BvgMQY6H60BhcK9wM79RBu9IlURIP26BWYcCiWJjs06L1yEdkUif2')
