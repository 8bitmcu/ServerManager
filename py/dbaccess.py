import sqlite3, json

class Dbaccess: 
    def __init__(self, name):
        # db filename
        self.name = name

    def get_db_connection(self):
        conn = sqlite3.connect(self.name)
        conn.row_factory = sqlite3.Row
        return conn


    def select_config(self, ):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_config LIMIT 1").fetchone()
        conn.close()
        return data

    def insertupdate_config(self, form): 
        conn = self.get_db_connection()

        values = [
            form["name"],
            form["password"],
            form["admin_password"],
            form["register_to_lobby"],
            form["pickup_mode_enabled"],
            form["locked_entry_list"],
            form["result_screen_time"],
            form["udp_port"],
            form["tcp_port"],
            form["http_port"],
            form["client_send_interval"],
            form["num_threads"],
            form["measurement_unit"],
            form["temp_unit"],
            form["install_path"]
        ]
        # TODO: change to `update if exists else insert` pattern
        conn.execute("DELETE FROM user_config");
        conn.execute("INSERT INTO user_config (name, password, admin_password, register_to_lobby, pickup_mode_enabled, locked_entry_list, result_screen_time, udp_port, tcp_port, http_port, client_send_interval, num_threads, measurement_unit, temp_unit, install_path) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", values)

        conn.commit()
        conn.close()


    def get_session(self):
        conn = self.get_db_connection()
        data = conn.execute("""
            SELECT
                s.id as id,
                s.laps as laps,
                s.strategy as strategy,
                t.name as track_name,
                t.length as track_length,
                t.pitboxes as pitboxes,
                d.name as difficulty_name,
                d.abs_allowed as abs_allowed,
                d.tc_allowed as tc_allowed,
                d.stability_allowed as stability_allowed,
                d.autoclutch_allowed as autoclutch_allowed,
                e.name as event_name,
                e.booking_enabled as booking_enabled,
                e.practice_enabled as practice_enabled,
                e.qualify_enabled as qualify_enabled,
                e.race_enabled as race_enabled,
                c.name as class_name,
                COUNT(ce.id) as entries,
                tw.name as time_name,
                tw.time as time
            FROM user_session s
            JOIN cache_track t
                on s.cache_track_id = t.id
            JOIN user_difficulty d
                on s.difficulty_id = d.id
            JOIN user_event e
                on s.event_id = e.id
            JOIN user_class c
                on s.class_id = c.id
            JOIN user_class_entry ce
                on s.class_id = ce.user_class_id
            JOIN user_time tw
                on s.time_id = tw.id
            GROUP BY (s.id)
""").fetchall()
        conn.close()
        return data

    def insertupdate_session(self, form):
        conn = self.get_db_connection()

        values = [
            form["track"],
            form["difficulty"],
            form["event"],
            form["class"],
            form["time"],
            form["laps"],
            form["strategy"]
        ]

        conn.execute("INSERT INTO user_session (cache_track_id, difficulty_id, event_id, class_id, time_id, laps, strategy) VALUES (?, ?, ?, ?, ?, ?, ?)", values)
        conn.commit()
        conn.close()

    def delete_session(self, id):
        conn = self.get_db_connection()
        conn.execute("DELETE FROM user_session WHERE id = ?", (id, ))
        conn.commit()
        conn.close()

    def get_difficulty(self, id):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_difficulty WHERE id = ? LIMIT 1", (id, )).fetchone()
        conn.close()
        return data

    def get_difficulty_list(self, filled=False):
        conn = self.get_db_connection()
        where = ""
        if filled:
            where = "WHERE filled = 1"
        data = conn.execute("SELECT id, name from user_difficulty " + where).fetchall()
        conn.close()
        return data

    def insert_difficulty(self, name):
        conn = self.get_db_connection()
        cur = conn.cursor();
        cur.execute("INSERT INTO user_difficulty (name) VALUES (?)", (name, ))
        conn.commit()
        conn.close()
        return cur.lastrowid;

    def delete_difficulty(self, id):
        conn = self.get_db_connection()
        conn.execute("DELETE FROM user_difficulty WHERE id = ?", (id, ))
        conn.commit()
        conn.close()

    def update_difficulty(self, id, form):
        print(form)
        conn = self.get_db_connection()
        values = [
            form["abs_allowed"],
            form["tc_allowed"],
            form["stability_allowed"],
            form["autoclutch_allowed"],
            form["tyre_blankets_allowed"],
            form["force_virtual_mirror"],
            form["fuel_rate"],
            form["damage_multiplier"],
            form["tyre_wear_rate"],
            form["allowed_tyres_out"],
            form["max_ballast_kg"],
            form["start_rule"],
            form["race_gas_penality_disabled"],
            form["dynamic_track"],
            form["dynamic_track_preset"],
            form["session_start"],
            form["randomness"],
            form["session_transfer"],
            form["lap_gain"],
            form["kick_quorum"],
            form["voting_quorum"],
            form["vote_duration"],
            form["blacklist_mode"],
            form["max_contacts_per_km"],
            id
        ]

        conn.execute("UPDATE user_difficulty SET abs_allowed = ?, tc_allowed = ?, stability_allowed = ?, autoclutch_allowed = ?, tyre_blankets_allowed = ?, force_virtual_mirror = ?, fuel_rate = ?, damage_multiplier = ?, tyre_wear_rate = ?, allowed_tyres_out = ?, max_ballast_kg = ?, start_rule = ?, race_gas_penality_disabled = ?, dynamic_track = ?, dynamic_track_preset = ?, session_start = ?, randomness = ?, session_transfer = ?, lap_gain = ?, kick_quorum = ?, voting_quorum = ?, vote_duration = ?, blacklist_mode = ?, max_contacts_per_km = ?, filled = 1 WHERE id = ?", values)

        conn.commit()
        conn.close()



    def get_event(self, id):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_event WHERE id = ? LIMIT 1", (id, )).fetchone()
        conn.close()
        return data

    def get_event_list(self, filled=False):
        conn = self.get_db_connection()
        where = ""
        if filled:
            where = "WHERE filled = 1"
        data = conn.execute("SELECT id, name from user_event " + where).fetchall()
        conn.close()
        return data

    def insert_event(self, name):
        conn = self.get_db_connection()
        cur = conn.cursor();
        cur.execute("INSERT INTO user_event (name) VALUES (?)", (name, ))
        conn.commit()
        conn.close()
        return cur.lastrowid;

    def delete_event(self, id):
        conn = self.get_db_connection()
        conn.execute("DELETE FROM user_event WHERE id = ?", (id, ))
        conn.commit()
        conn.close()

    def update_event(self, id, form):
        conn = self.get_db_connection()
        values = [
            form["booking_enabled"],
            form["booking_time"],
            form["practice_enabled"],
            form["practice_time"],
            form["practice_is_open"],
            form["qualify_enabled"],
            form["qualify_time"],
            form["qualify_is_open"],
            form["qualify_max_wait_perc"],
            form["race_enabled"],
            form["race_laps"],
            form["race_time"],
            form["race_over_time"],
            form["race_wait_time"],
            form["race_is_open"],
            form["reversed_grid_positions"],
            form["race_pit_window_start"],
            form["race_pit_window_end"],
            id
        ]

        conn.execute("UPDATE user_event SET booking_enabled = ?, booking_time = ?, practice_enabled = ?, practice_time = ?, practice_is_open = ?, qualify_enabled = ?, qualify_time = ?, qualify_is_open = ?, qualify_max_wait_perc = ?, race_enabled = ?, race_laps = ?, race_time = ?, race_over_time = ?, race_wait_time = ?, race_is_open = ?, reversed_grid_positions = ?, race_pit_window_start = ?, race_pit_window_end = ?, filled = 1 WHERE id = ?", values)

        conn.commit()
        conn.close()






    def get_time(self, id):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_time WHERE id = ? LIMIT 1", (id, )).fetchone()
        conn.close()
        return data

    def get_weather(self, id):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_time_weather WHERE user_time_id = ?", (id, )).fetchall()
        conn.close()
        return data

    def get_time_list(self, filled=False):
        conn = self.get_db_connection()
        where = ""
        if filled:
            where = "WHERE filled = 1"
        data = conn.execute("SELECT id, name from user_time " + where).fetchall()
        conn.close()
        return data

    def insert_time(self, name):
        conn = self.get_db_connection()
        cur = conn.cursor();
        cur.execute("INSERT INTO user_time (name) VALUES (?)", (name, ))
        conn.commit()
        conn.close()
        return cur.lastrowid;

    def delete_time(self, id):
        conn = self.get_db_connection()
        conn.execute("DELETE FROM user_time WHERE id = ?", (id, ))
        conn.commit()
        conn.close()

    def update_time(self, id, form):
        conn = self.get_db_connection()
        values = [
            form["time"],
            form["time_of_day_multi"],
            id
        ]

        conn.execute("UPDATE user_time SET time = ?, time_of_day_multi = ?, filled = 1 WHERE id = ?", values)
        conn.execute("DELETE FROM user_time_weather WHERE user_time_id = ?", (id, ))

        for weather in json.loads(form["weather"]):
            values = [
                id,
                weather["graphics"],
                weather["base_temperature_ambient"],
                weather["base_temperature_road"],
                weather["variation_ambient"],
                weather["variation_road"],
                weather["wind_base_speed_min"],
                weather["wind_base_speed_max"],
                weather["wind_base_direction"],
                weather["wind_variation_direction"]
            ]
            conn.execute("INSERT INTO user_time_weather (user_time_id, graphics, base_temperature_ambient, base_temperature_road, variation_ambient, variation_road, wind_base_speed_min, wind_base_speed_max, wind_base_direction, wind_variation_direction) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", values)

        conn.commit()
        conn.close()


    def get_class(self, id):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_class WHERE id = ? LIMIT 1", (id, )).fetchone()
        conn.close()
        return data

    def get_class_entries(self, id):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_class_entry WHERE user_class_id = ?", (id, )).fetchall()
        conn.close()
        return data

    def get_class_list(self, filled=False):
        conn = self.get_db_connection()
        where = ""
        if filled:
            where = "WHERE filled = 1"
        data = conn.execute("SELECT id, name from user_class " + where).fetchall()
        conn.close()
        return data

    def insert_class(self, name):
        conn = self.get_db_connection()
        cur = conn.cursor();
        cur.execute("INSERT INTO user_class (name) VALUES (?)", (name, ))
        conn.commit()
        conn.close()
        return cur.lastrowid;

    def delete_class(self, id):
        conn = self.get_db_connection()
        conn.execute("DELETE FROM user_class WHERE id = ?", (id, ))
        conn.commit()
        conn.close()

    def update_class(self, id, form):
        conn = self.get_db_connection()

        conn.execute("DELETE FROM user_class_entry WHERE user_class_id = ?", (id, ))
        conn.execute("UPDATE user_class SET filled = 1 WHERE id = ?", (id, ))

        for car in json.loads(form['cars']):
            values = [
                id,
                car['key'],
                car['skin'],
                # TODO
                #car['ballast']
            ]
            conn.execute("INSERT INTO user_class_entry (user_class_id, cache_vehicle_id, skin_id) VALUES (?, ?, ?)", values)

        conn.commit()
        conn.close()




    def update_vehicles(self, data):
        conn = self.get_db_connection()
        conn.execute("DELETE FROM cache_vehicle")
        conn.executemany("INSERT INTO cache_vehicle (key, name, brand, desc, tags, class, specs, torque, power, skins) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", data)

        conn.commit()
        conn.close()

    def get_carlist(self):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * FROM cache_vehicle ORDER BY Name ASC").fetchall()

        conn.close()
        return data

    def get_car(self, car):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * FROM cache_vehicle WHERE key = ?", (car, )).fetchone()

        conn.close()
        return data



    def update_tracks(self, data):
        conn = self.get_db_connection()
        conn.execute("DELETE FROM cache_track")
        conn.executemany("INSERT INTO cache_track (key, config, name, desc, tags, country, city, length, width, pitboxes, run) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", data)

        conn.commit()
        conn.close()

    def get_tracklist(self):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * FROM cache_track ORDER BY Name ASC").fetchall()

        conn.close()
        return data
