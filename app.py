import sqlite3, os, json
from flask import Flask, render_template, request, redirect, url_for, send_from_directory, jsonify
from py.dbaccess import Dbaccess
from py.fsaccess import Fsaccess

app = Flask(__name__, template_folder="htm")

dba = Dbaccess("sqlite.db")
fsa = Fsaccess(dba)


def clean_dbdata(dbdata):
    newdata = {}
    if dbdata is None:
        return None
    for i in dict(dbdata):
        if dbdata[i] is not None:
            newdata[i] = dbdata[i]
    return newdata

# Dashboard / Index
@app.route("/")
def index():
    return render_template("index.htm")

# Configuration
@app.route("/config", methods=('GET', 'POST'))
def config():
    if request.method == 'POST':
        dba.insertupdate_config(request.form)

    data = clean_dbdata(dba.select_config())
    return render_template("config.htm", data=data)

# Difficulties
@app.route("/difficulty", methods=('GET', 'POST'))
@app.route("/difficulty/<int:id>", methods=('GET', 'POST'))
def difficulty(id=None):
    data = {}

    if request.method == 'POST':
        if id is None or "difficulty_name" in request.form:
            newid = dba.insert_difficulty(request.form['difficulty_name'])
            return redirect(url_for("difficulty", id=newid))
        else:
            dba.update_difficulty(id, request.form)

    if id is not None:
        data['id'] = id
        data['form'] = clean_dbdata(dba.get_difficulty(id))

    data['list'] = dba.get_difficulty_list()

    return render_template("difficulties.htm", data=data)

@app.route("/difficulty/delete/<int:id>", methods=('POST',))
def difficulty_delete(id):
    dba.delete_difficulty(id)
    return redirect(url_for("difficulty"));

# Events
@app.route("/event", methods=('GET', 'POST'))
@app.route("/event/<int:id>", methods=('GET', 'POST'))
def event(id=None):
    data = {}

    if request.method == 'POST':
        if id is None or "event_name" in request.form:
            newid = dba.insert_event(request.form['event_name'])
            return redirect(url_for("event", id=newid));
        else:
            dba.update_event(id, request.form)

    if id is not None:
        data['id'] = id
        data['form'] = clean_dbdata(dba.get_event(id))


    data['list'] = dba.get_event_list()

    return render_template("events.htm", data=data)

@app.route("/event/delete/<int:id>", methods=('POST',))
def event_delete(id):
    dba.delete_event(id)
    return redirect(url_for("event"));

# Time & Weather
@app.route("/time", methods=('GET', 'POST'))
@app.route("/time/<int:id>", methods=('GET', 'POST'))
def time(id=None):
    data = {}

    if request.method == 'POST':
        if id is None or "time_name" in request.form:
            newid = dba.insert_time(request.form['time_name'])
            return redirect(url_for("time", id=newid));
        else:
            dba.update_time(id, request.form)

    if id is not None:
        data['id'] = id
        data['form'] = clean_dbdata(dba.get_time(id))
        wdata = [dict(ix) for ix in dba.get_weather(id)]

        # add at least one weather panel
        if len(wdata) == 0: 
            wdata = [{}]

        data['weather'] = json.dumps(wdata)

    data['list'] = dba.get_time_list()

    return render_template("times.htm", data=data)

@app.route("/time/delete/<int:id>", methods=('POST',))
def time_delete(id):
    dba.delete_time(id)
    return redirect(url_for("time"));

# Vehicle Classes
@app.route("/class", methods=('GET', 'POST'))
@app.route("/class/<int:id>", methods=('GET', 'POST'))
def veh_class(id=None):
    data = {}

    if request.method == 'POST':
        if id is None or "class_name" in request.form:
            newid = dba.insert_class(request.form['class_name'])
            return redirect(url_for("veh_class", id=newid));
        else:
            dba.update_class(id, request.form)

    data['list'] = dba.get_class_list()
    if id is not None:
        data['id'] = id
        data['form'] = dba.get_class(id)

        # TODO: refactor vehicles
        data['vehicles'] = dba.get_carlist() 


        # add at least one car panel

        data['entries'] =[dict(ix) for ix in dba.get_class_entries(id)] 

        data['car_data'] = []
        for car in dba.get_carlist():
            data['car_data'].append({
                'key': car['key'],
                'skins': json.loads(car['skins']),
                'name': car['name'],
                'tags': json.loads(car['tags']),
                'specs': json.loads(car['specs'])
            })

    return render_template("classes.htm", data=data)

@app.route("/class/delete/<int:id>", methods=('POST',))
def class_delete(id):
    dba.delete_class(id)
    return redirect(url_for("veh_class"));


@app.route("/api/skin_img/<string:car>/<string:skin>.jpg")
def skin_img(car, skin):
    return send_from_directory(fsa.get_skin(car, skin), "preview.jpg")


@app.route("/api/car_details/<string:car>")
def car_details(car): 
    car_data = dba.get_car(car)

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

@app.route("/api/update_carlist")
def api_update_carlist():
    fsa.parse_cars_folder(dba)
    return jsonify({'ok': 'OK'})


@app.route("/api/applyschema")
def api_applyschema(): 
    connection = sqlite3.connect('sqlite.db')
    with open('schema.sql') as f:
        connection.executescript(f.read())
    connection.commit()
    connection.close()
    return jsonify({'ok': 'OK'})

if __name__ == "__main__":
    app.run(debug=True)
