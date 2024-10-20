import sqlite3, os, json
import json_repair



class DbAccess: 
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
        # todo: change to `update if exists else insert` pattern
        conn.execute("DELETE FROM user_config");
        conn.execute("INSERT INTO user_config (name, password, admin_password, register_to_lobby, pickup_mode_enabled, locked_entry_list, result_screen_time, udp_port, tcp_port, http_port, client_send_interval, num_threads, measurement_unit, temp_unit, install_path) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", values)

        conn.commit()
        conn.close()




    def get_difficulty(self, id):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_difficulty WHERE id = ? LIMIT 1", (id, )).fetchone()
        conn.close()
        return data

    def get_difficulty_list(self, ):
        conn = self.get_db_connection()
        data = conn.execute("SELECT id, name from user_difficulty").fetchall()
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
            form["session_start"],
            form["randomness"],
            form["session_transfer"],
            form["lap_gain"],
            form["kick_quorum"],
            form["blacklist_mode"],
            form["max_contacts_per_km"],
            id
        ]

        conn.execute("UPDATE user_difficulty SET abs_allowed = ?, tc_allowed = ?, stability_allowed = ?, autoclutch_allowed = ?, tyre_blankets_allowed = ?, force_virtual_mirror = ?, fuel_rate = ?, damage_multiplier = ?, tyre_wear_rate = ?, allowed_tyres_out = ?, max_ballast_kg = ?, start_rule = ?, race_gas_penality_disabled = ?, dynamic_track = ?, session_start = ?, randomness = ?, session_transfer = ?, lap_gain = ?, kick_quorum = ?, blacklist_mode = ?, max_contacts_per_km = ? WHERE id = ?", values)

        conn.commit()
        conn.close()



    def get_event(self, id):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_event WHERE id = ? LIMIT 1", (id, )).fetchone()
        conn.close()
        return data

    def get_event_list(self, ):
        conn = self.get_db_connection()
        data = conn.execute("SELECT id, name from user_event").fetchall()
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

        conn.execute("UPDATE user_event SET booking_enabled = ?, booking_time = ?, practice_enabled = ?, practice_time = ?, practice_is_open = ?, qualify_enabled = ?, qualify_time = ?, qualify_is_open = ?, qualify_max_wait_perc = ?, race_enabled = ?, race_laps = ?, race_time = ?, race_over_time = ?, race_wait_time = ?, race_is_open = ?, reversed_grid_positions = ?, race_pit_window_start = ?, race_pit_window_end = ? WHERE id = ?", values)

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

    def get_time_list(self, ):
        conn = self.get_db_connection()
        data = conn.execute("SELECT id, name from user_time").fetchall()
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

        conn.execute("UPDATE user_time SET time = ?, time_of_day_multi = ? WHERE id = ?", values)
        conn.execute("DELETE FROM user_time_weather WHERE user_time_id = ?", (id, ))

        for x in json.loads(form["weather"]):
            values = [
                id,
                x["graphics"],
                x["base_temperature_ambient"],
                x["base_temperature_road"],
                x["variation_ambient"],
                x["variation_road"],
                x["wind_base_speed_min"],
                x["wind_base_speed_max"],
                x["wind_base_direction"],
                x["wind_variation_direction"]
            ]
            conn.execute("INSERT INTO user_time_weather (user_time_id, graphics, base_temperature_ambient, base_temperature_road, variation_ambient, variation_road, wind_base_speed_min, wind_base_speed_max, wind_base_direction, wind_variation_direction) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", values)

        conn.commit()
        conn.close()


    def get_class(self, id):
        conn = self.get_db_connection()
        data = conn.execute("SELECT * from user_class WHERE id = ? LIMIT 1", (id, )).fetchone()
        conn.close()
        return data

    def get_class_list(self, ):
        conn = self.get_db_connection()
        data = conn.execute("SELECT id, name from user_class").fetchall()
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
        values = [
            id
        ]

        conn.execute("UPDATE user_class WHERE id = ?", values)

        # TODO: linked objects

        conn.commit()
        conn.close()




    def update_carlist(self):
        cars_path = os.path.join(self.select_config()["install_path"], "content/cars")

        conn = self.get_db_connection()
        conn.execute("DELETE FROM cache_vehicle")

        for key in os.listdir(cars_path):
            # if data.acd is missing, assume it"s a missing dlc and avoid listing/saving it
            data_acd = os.path.join(cars_path, key, "data.acd")
            if not os.path.isfile(data_acd):
                continue

            try:
                json_path = os.path.join(cars_path, key, "ui/ui_car.json")
                data = json_repair.from_file(json_path)
            except:
                try:
                    json_path = os.path.join(cars_path, key, "ui/dlc_ui_car.json")
                    data = json_repair.from_file(json_path)
                except:
                    data = None

            skins_path = os.path.join(cars_path, key, "skins")

            skins = []
            for skin in os.listdir(skins_path):
                skin_json = os.path.join(skins_path, skin, "ui_skin.json")
                skin_data = json_repair.from_file(skin_json)

                skin_name = skin_data.get("skinname")
                if skin_name == "":
                    skin_name = skin
                skins.append({"skin_id": skin, "skin_name": skin_name})

            if data is not None:
                name = data.get("name")
                brand = data.get("brand")
                desc = data.get("description")
                tags = json.dumps(data.get("tags"))
                cls = data.get("class")
                specs = json.dumps(data.get("specs"))
                torque = json.dumps(data.get("torqueCurve"))
                power = json.dumps(data.get("powerCurve"))
                allskins = json.dumps(skins)

                conn.execute("INSERT INTO cache_vehicle (key, name, brand, desc, tags, class, specs, torque, power, skins) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", (key, name, brand, desc, tags, cls, specs, torque, power, allskins))

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
