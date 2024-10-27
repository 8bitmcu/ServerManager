import os, json, re
import json_repair

class Fsaccess: 
    def __init__(self, dba):
        # base root folder
        self.dba = dba
        self.SM_DATA = os.environ.get('SM_DATA')

    def get_basepath(self):
        return self.dba.select_config()["install_path"]


    def get_datadir(self):
        if self.SM_DATA is None:
            sm_data = os.path.join(os.path.curdir, "smdata")
            if not os.path.exists(sm_data):
                os.mkdir(sm_data)
            return sm_data
        else:
            return self.SM_DATA

    def get_database(self, ):
        return os.path.join(self.get_datadir(), "sm.db")

    def get_serverexe(self, osa):
        if osa.is_unix():
            return os.path.join(self.get_serverpath(), "acServer")
        else:
            return os.path.join(self.get_serverpath(), "acServer.exe")

    def get_serverpath(self):
        return os.path.join(self.get_basepath(), "server")

    def get_skin(self, car_id, skin_id):
        return os.path.join(self.get_basepath(), "content", "cars", car_id, "skins", skin_id)

    def get_track(self, track, config):
        if config == None:
            return os.path.join(self.get_basepath(), "content", "tracks", track, "ui")
        else:
            return os.path.join(self.get_basepath(), "content", "tracks", track, "ui", config)

    def set_server_ini(self, server_cfg, entry_list):
        server_cfg_ini = os.path.join(self.get_serverpath(), "cfg", "server_cfg.ini")
        f = open(server_cfg_ini, "w")
        f.write(server_cfg)
        f.close()

        entry_list_ini = os.path.join(self.get_serverpath(), "cfg", "entry_list.ini")
        f = open(entry_list_ini, "w")
        f.write(entry_list)
        f.close()

    def validate_installpath(self, server_path):
        return os.path.exists(os.path.join(server_path, "acs.exe"))

    def parse_cars_folder(self, dba):
        db_data = []
        cars_path = os.path.join(self.get_basepath(), "content", "cars")

        for key in os.listdir(cars_path):
            # if data.acd is missing, assume it"s a missing dlc and avoid listing/saving it
            data_acd = os.path.join(cars_path, key, "data.acd")
            if not os.path.isfile(data_acd):
                continue

            try:
                json_path = os.path.join(cars_path, key, "ui", "ui_car.json")
                data = json_repair.from_file(json_path)
            except:
                try:
                    json_path = os.path.join(cars_path, key, "ui", "dlc_ui_car.json")
                    data = json_repair.from_file(json_path)
                except:
                    print("No data for " + key)
                    data = None

            skins_path = os.path.join(cars_path, key, "skins")

            skins = []
            for skin in os.listdir(skins_path):

                try:
                    skin_json = os.path.join(skins_path, skin, "ui_skin.json")
                    skin_data = json_repair.from_file(skin_json)

                    skin_name = skin_data.get("skinname")
                except:
                    skin_name = ""

                if skin_name == "":
                    skin_name = skin
                skins.append({"skin_id": skin, "skin_name": skin_name})

            if data is not None:
                db_data.append([
                    key,
                    data.get("name"),
                    data.get("brand"),
                    data.get("description"),
                    json.dumps(data.get("tags")),
                    data.get("class"),
                    json.dumps(data.get("specs")),
                    json.dumps(data.get("torqueCurve")),
                    json.dumps(data.get("powerCurve")),
                    json.dumps(skins)
                ])
        dba.update_vehicles(db_data)
        return len(db_data)

    def parse_tracks_folder(self, dba):
        db_data = []
        tracks_path = os.path.join(self.get_basepath(), "content", "tracks")

        for key in os.listdir(tracks_path):
            # if skins is missing, assume it"s a missing dlc and avoid listing/saving it
            data_acd = os.path.join(tracks_path, key, "skins")
            if not os.path.isdir(data_acd):
                continue

            try:
                json_path = os.path.join(tracks_path, key, "ui", "ui_track.json")
                data = json_repair.from_file(json_path)

                # EXCEPTION: for some reason laguna seca's length is a float instead of an int
                length = float(data.get("length"))
                if key == "ks_laguna_seca":
                    length = length * 1000

                db_data.append([
                    key,
                    '',
                    data.get("name"),
                    data.get("desc"),
                    json.dumps(data.get("tags")),
                    data.get("country"),
                    data.get("city"),
                    length,
                    data.get("width"),
                    data.get("pitboxes"),
                    data.get("run")
                ])
            except Exception as e:
                print(e)
                for subtrack in os.listdir(os.path.join(tracks_path, key, "ui")):
                    if os.path.isfile(os.path.join(tracks_path, key, "ui", subtrack)):
                        continue
                    try:
                        json_path = os.path.join(tracks_path, key, "ui", subtrack, "ui_track.json")

                        with open(json_path, 'rb') as f:
                            lines = f.read()

                        data = json_repair.loads(lines.decode('latin-1'))

                    except Exception as e:
                        json_path = os.path.join(tracks_path, key, "ui", subtrack, "dlc_ui_track.json")

                        with open(json_path, 'rb') as f:
                            lines = f.read()

                        data = json_repair.loads(lines.decode('latin-1'))

                    db_data.append([
                        key,
                        subtrack,
                        data.get("name"),
                        data.get("desc"),
                        json.dumps(data.get("tags")),
                        data.get("country"),
                        data.get("city"),
                        data.get("length"),
                        data.get("width"),
                        data.get("pitboxes"),
                        data.get("run")
                    ])

        dba.update_tracks(db_data)
        return len(db_data)


    def parse_weathers_folder(self, dba):
        db_data = []
        weather_path = os.path.join(self.get_basepath(), "content", "weather")

        for key in os.listdir(weather_path):
            if not os.path.isdir(os.path.join(weather_path, key)):
                continue

            try:
                ini_path = os.path.join(weather_path, key, "weather.ini")
                with open(ini_path, 'rb') as f:
                    name = None
                    while name is None:
                        line = f.readline().decode('utf-8')
                        if re.search("^NAME", line):
                            name = re.findall("^NAME=(.*)\r", line)
                            name = re.sub("; .*", "", name[0])

                db_data.append([
                    key,
                    name
                ])
            except Exception as e:
                print(e)

        dba.update_weathers(db_data)
        return len(db_data)
