
import os, json
import json_repair

class Fsaccess: 
    def __init__(self, dba):
        # base root folder
        self.dba = dba

    def get_basepath(self):
        return self.dba.select_config()["install_path"]


    def get_skin(self, car_id, skin_id):
        return os.path.join(self.get_basepath(), "content", "cars", car_id, "skins", skin_id)

    def parse_cars_folder(self, dba):
        db_data = []
        cars_path = os.path.join(self.get_basepath(), "content", "cars")

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

    
