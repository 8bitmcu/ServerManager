import sqlite3, os, json
from flask import Flask, render_template, request, redirect, url_for, send_from_directory, jsonify
import json_repair

app = Flask(__name__, template_folder="htm")



#connection = sqlite3.connect('sqlite.db')
#with open('schema.sql') as f:
#    connection.executescript(f.read())
#connection.commit()
#connection.close()


def update_carlist():
    cars_path = os.path.join(get_uservalue("install_path"), "content/cars")

    conn = get_db_connection()
    conn.execute("DELETE FROM cache_vehicle")

    for key in os.listdir(cars_path):
        # if data.acd is missing, assume it's a missing dlc and avoid listing/saving it
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
            skins.append({'skin_id': skin, 'skin_name': skin_name})

        if data is not None:
            name = data.get('name')
            brand = data.get('brand')
            desc = data.get('description')
            tags = json.dumps(data.get('tags'))
            cls = data.get('class')
            specs = json.dumps(data.get('specs'))
            torque = json.dumps(data.get('torqueCurve'))
            power = json.dumps(data.get('powerCurve'))
            allskins = json.dumps(skins)

            conn.execute("INSERT INTO cache_vehicle (key, name, brand, desc, tags, class, specs, torque, power, skins) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", (key, name, brand, desc, tags, cls, specs, torque, power, allskins))

    conn.commit()
    conn.close()

def get_carlist():
    conn = get_db_connection()
    data = conn.execute("SELECT * FROM cache_vehicle ORDER BY Name ASC").fetchall()

    conn.close()
    return data

def get_car(car):
    conn = get_db_connection()
    data = conn.execute("SELECT * FROM cache_vehicle WHERE key = ?", (car, )).fetchone()

    conn.close()
    return data


def get_db_connection():
    conn = sqlite3.connect('sqlite.db')
    conn.row_factory = sqlite3.Row
    return conn

def insert_update_uservalues(category, values, id=None):
    conn = get_db_connection()

    if id is None:
        conn.execute("DELETE FROM user_value WHERE config_category_id = (SELECT id FROM config_category WHERE key = ?)", (category, ))
        conn.executemany("INSERT INTO user_value (config_category_id, key, value) VALUES ((SELECT id FROM config_category WHERE key = '" + category + "'), ?, ?)", values)
    else:
        conn.execute("DELETE FROM user_value WHERE config_category_id = (SELECT id FROM config_category WHERE key = ? and user_preset_id = ?)", (category, id))
        conn.executemany("INSERT INTO user_value (config_category_id, key, value, user_preset_id) VALUES ((SELECT id FROM config_category WHERE key = '" + category + "'), ?, ?, ?)", values)
    conn.commit()
    conn.close()

def get_uservalue(key, id=None):
    data = {}
    conn = get_db_connection()

    if id is None:
        data = conn.execute("SELECT value FROM user_value WHERE key = ?", (key, )).fetchone()[0]
    else:
        data = conn.execute("SELECT value FROM user_value WHERE key = ? and user_preset_id = ?", (key, id)).fetchone()[0]

    conn.close()
    return data

def get_uservalues(category, id=None):
    data = {}
    conn = get_db_connection()

    if id is None:
        for row in conn.execute("SELECT key, value FROM user_value WHERE config_category_id = (SELECT id FROM config_category WHERE key = ?)", (category,)).fetchall():
            data[row['key']] = row['value'];
    else:
        for row in conn.execute("SELECT key, value FROM user_value WHERE config_category_id = (SELECT id FROM config_category WHERE key = ? and user_preset_id = ?)", (category, id)).fetchall():
            data[row['key']] = row['value'];

    conn.close()
    return data

def delete_uservalues(category, id):
    conn = get_db_connection()
    conn.execute("DELETE FROM user_value WHERE config_category_id = (SELECT id FROM config_category WHERE key = ? and user_preset_id = ?)", (category, id))
    conn.execute("DELETE FROM user_preset WHERE id = ?", (id,))
    conn.commit()
    conn.close()

def get_presetlist(category):
    conn = get_db_connection()
    data = conn.execute("SELECT id, name FROM user_preset WHERE config_category_id = (SELECT id FROM config_category WHERE key = ?)", (category,)).fetchall()
    conn.close()
    return data

def add_preset(category, name):
    conn = get_db_connection()
    cur = conn.cursor();
    cur.execute("INSERT INTO user_preset (config_category_id, name) VALUES ((SELECT id FROM config_category WHERE key = ?), ?)", (category, name))
    conn.commit()
    conn.close()
    return cur.lastrowid;














@app.route("/")
def index():
    update_carlist()
    return render_template("index.htm")


@app.route("/config", methods=('GET', 'POST'))
def config():
    if request.method == 'POST':
        values = []
        for key in request.form: 
            values.append((key, request.form[key]))
        insert_update_uservalues('SERVER', values)
        data = request.form
    else:
        data = get_uservalues('SERVER')

    return render_template("config.htm", data=data)



@app.route("/difficulty", methods=('GET', 'POST'))
@app.route("/difficulty/<int:id>", methods=('GET', 'POST'))
def difficulty(id=None):
    data = {}
    conn = get_db_connection()

    if request.method == 'POST':
        # Add a new difficulty and navigate to it
        if id is None or "difficulty_name" in request.form:
            newid = add_preset('DIFF', request.form['difficulty_name'])
            return redirect(url_for("difficulty", id=newid));
        # Insert or Update the form
        else:
            values = []
            for key in request.form: 
                values.append((key, request.form[key], id))

            insert_update_uservalues('DIFF', values, id)

    if id is not None:
        data['id'] = id
        data['form'] = get_uservalues('DIFF', id)


    data['list'] = get_presetlist('DIFF') 

    conn.close()
    return render_template("difficulties.htm", data=data)

@app.route("/difficulty/delete/<int:id>", methods=('POST',))
def difficulty_delete(id):
    delete_uservalues('DIFF', id)
    return redirect(url_for("difficulty"));



@app.route("/event", methods=('GET', 'POST'))
@app.route("/event/<int:id>", methods=('GET', 'POST'))
def event(id=None):
    data = {}
    conn = get_db_connection()

    # POST can be from either form: 'Add new Event' or edit
    if request.method == 'POST':
        if id is None or "event_name" in request.form:
            newid = add_preset('EV_TYPE', request.form['event_name'])
            return redirect(url_for("event", id=newid));
        else:
            values = []
            for key in request.form: 
                values.append((key, request.form[key], id))

            insert_update_uservalues('EV_TYPE', values, id)

    if id is not None:
        data['id'] = id
        data['form'] = get_uservalues('EV_TYPE', id)


    data['list'] = get_presetlist('EV_TYPE') 

    conn.close()
    return render_template("events.htm", data=data)

@app.route("/event/delete/<int:id>", methods=('POST',))
def event_delete(id):
    delete_uservalues('EV_TYPE', id)
    return redirect(url_for("event"));



@app.route("/time", methods=('GET', 'POST'))
@app.route("/time/<int:id>", methods=('GET', 'POST'))
def time(id=None):
    data = {}
    conn = get_db_connection()

    # POST can be from either form: 'Add new Event' or edit
    if request.method == 'POST':
        if id is None or "time_name" in request.form:
            newid = add_preset('TIME', request.form['time_name'])
            return redirect(url_for("time", id=newid));
        else:
            values = []
            for key in request.form: 
                values.append((key, request.form[key], id))

            insert_update_uservalues('TIME', values, id)

    if id is not None:
        data['id'] = id
        data['form'] = get_uservalues('TIME', id)

        # add/have at least one panel
        if data['form'].get('weather') is None: 
            data['form']['weather'] = '[{}]'


    data['list'] = get_presetlist('TIME')

    conn.close()
    return render_template("times.htm", data=data)

@app.route("/time/delete/<int:id>", methods=('POST',))
def time_delete(id):
    delete_uservalues('TIME', id)
    return redirect(url_for("time"));


@app.route("/class", methods=('GET', 'POST'))
@app.route("/class/<int:id>", methods=('GET', 'POST'))
def veh_class(id=None):
    data = {}
    conn = get_db_connection()

    # POST can be from: 'Add new Class' or edit
    if request.method == 'POST':
        if id is None or "class_name" in request.form:
            newid = add_preset('VEH_CLASS', request.form['class_name'])
            return redirect(url_for("veh_class", id=newid));
        else:
            values = []
            for key in request.form: 
                values.append((key, request.form[key], id))

            insert_update_uservalues('VEH_CLASS', values, id)

    data['list'] = get_presetlist('VEH_CLASS') 
    if id is not None:
        data['id'] = id
        data['form'] = get_uservalues('VEH_CLASS', id)

        cars = get_carlist()
        data['vehicles'] = cars


        data['car_data'] = []
        for car in cars:
            data['car_data'].append({
                'key': car['key'],
                'skins': json.loads(car['skins']),
                'name': car['name'],
                'tags': json.loads(car['tags']),
                'specs': json.loads(car['specs'])
            })

    conn.close()
    return render_template("classes.htm", data=data)

@app.route("/class/delete/<int:id>", methods=('POST',))
def class_delete(id):
    delete_uservalues('VEH_CLASS', id)
    return redirect(url_for("veh_class"));


@app.route("/api/skin_img/<string:car>/<string:skin>")
def skin_img(car, skin): 
    cars_path = os.path.join(get_uservalue("install_path"), "content", "cars", car, "skins", skin)
    return send_from_directory(cars_path, "preview.jpg")


@app.route("/api/car_details/<string:car>")
def car_details(car): 
    car_data = get_car(car)

    raw_power = json.loads(car_data['power'])

    torque = []
    for t in json.loads(car_data['torque']):
        torque.append(int(t[1]))

    power = []
    for p in raw_power:
        power.append(int(p[1]))

    labels = []
    for p in raw_power: 
        labels.append(int(p[0]))

    json_data = {
        'id': car,
        'desc': car_data['desc'],
        'power': power,
        'torque': torque,
        'labels': labels
    }

    return jsonify(json_data)




if __name__ == "__main__":
    app.run(debug=True)
