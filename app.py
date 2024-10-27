import sqlite3, json
from flask import Flask, render_template, request, redirect, url_for, send_from_directory, jsonify
from functools import wraps
from py.dbaccess import Dbaccess
from py.fsaccess import Fsaccess
from py.psaccess import Psaccess
from py.osaccess import Osaccess
from py.utilities import Utilities

app = Flask(__name__, template_folder="htm")


dba = Dbaccess()
osa = Osaccess()
fsa = Fsaccess(dba)
psa = Psaccess(dba, fsa)
util = Utilities()

@app.before_request
def before_request_func():
    dba.open(fsa.get_database())

def require_config_set(f):
    @wraps(f)

    def check_config_exists(*args, **kwargs):
        if not dba.table_exists("user_config"):
            dba.apply_schema()
        if dba.select_config() is None:
            return redirect(url_for("config", welcome=True))
        return f(*args, **kwargs)
    return check_config_exists





# Dashboard / Index
@app.route("/", methods=('GET', 'POST'))
@require_config_set
def index():
    data = {}
    data['is_running'] = psa.is_running()
    return render_template("index.htm", data=data)

# Session
@app.route("/event", methods=('GET', 'POST'))
@app.route("/event/<int:id>", methods=('GET', 'POST'))
@require_config_set
def event(id=None):
    if request.method == 'POST' and request.form['id']:
        dba.update_event(request.form)
    elif request.method == 'POST':
        dba.insert_event(request.form)

    data = {}

    # TODO: refactor?
    data['difficulties'] = [dict(ix) for ix in dba.get_difficulty_list(filled=True)]
    data['sessions'] = [dict(ix) for ix in dba.get_session_list(filled=True)]
    data['times'] = [dict(ix) for ix in dba.get_time_list(filled=True)]
    data['classes'] = [dict(ix) for ix in dba.get_class_list(filled=True)]

    data['events'] = dba.get_events()

    data['track_data'] = [dict(ix) for ix in dba.get_tracklist()]
    for x in data['track_data']:
        x['tags'] = json.loads(x['tags'])

    data['form'] = {}
    if id is not None:
        data['form'] = dba.get_event(id)

    return render_template("events.htm", data=data)

@app.route("/event_delete/<int:id>", methods=('GET', 'POST'))
@require_config_set
def event_delete(id):
    dba.delete_event(id)
    return redirect(url_for("index"));

# Configuration
@app.route("/config", methods=('GET', 'POST'))
def config():
    if request.method == 'POST':
        dba.insertupdate_config(request.form)

    # TODO: when install_path changes, refresh content

    data = {}
    data['form'] = util.clean_dbdata(dba.select_config())
    data['welcome'] = request.args.get("welcome")
    return render_template("config.htm", data=data)

# Content and Mods
@app.route("/content", methods=('GET', 'POST'))
def content():
    data = {}
    data['form'] = {}
    return render_template("content.htm", data=data)

# Difficulties
@app.route("/difficulty", methods=('GET', 'POST'))
@app.route("/difficulty/<int:id>", methods=('GET', 'POST'))
@require_config_set
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
        data['form'] = util.clean_dbdata(dba.get_difficulty(id))

    data['list'] = dba.get_difficulty_list()

    return render_template("difficulties.htm", data=data)

@app.route("/difficulty/delete/<int:id>", methods=('POST',))
@require_config_set
def difficulty_delete(id):
    dba.delete_difficulty(id)
    return redirect(url_for("difficulty"));

# Events
@app.route("/session", methods=('GET', 'POST'))
@app.route("/session/<int:id>", methods=('GET', 'POST'))
@require_config_set
def session(id=None):
    data = {}

    if request.method == 'POST':
        if id is None or "session_name" in request.form:
            newid = dba.insert_session(request.form['session_name'])
            return redirect(url_for("session", id=newid));
        else:
            dba.update_session(id, request.form)

    if id is not None:
        data['id'] = id
        data['form'] = util.clean_dbdata(dba.get_session(id))


    data['list'] = dba.get_session_list()

    return render_template("sessions.htm", data=data)

@app.route("/session/delete/<int:id>", methods=('POST',))
@require_config_set
def session_delete(id):
    dba.delete_session(id)
    return redirect(url_for("session"));

# Time & Weather
@app.route("/time", methods=('GET', 'POST'))
@app.route("/time/<int:id>", methods=('GET', 'POST'))
@require_config_set
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
        data['form'] = util.clean_dbdata(dba.get_time(id))
        wdata = [dict(ix) for ix in dba.get_time_weather(id)]

        # add at least one weather panel
        if len(wdata) == 0: 
            wdata = [{}]

        data['weather'] = json.dumps(wdata)

    data['list'] = dba.get_time_list()
    data['weatherlist'] = dba.get_weatherlist()

    return render_template("times.htm", data=data)

@app.route("/time/delete/<int:id>", methods=('POST',))
@require_config_set
def time_delete(id):
    dba.delete_time(id)
    return redirect(url_for("time"));

# Vehicle Classes
@app.route("/class", methods=('GET', 'POST'))
@app.route("/class/<int:id>", methods=('GET', 'POST'))
@require_config_set
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
            'id': car['id'],
            'key': car['key'],
            'skins': json.loads(car['skins']),
            'name': car['name'],
            'tags': json.loads(car['tags']),
            'specs': json.loads(car['specs'])
        })

    return render_template("classes.htm", data=data)

@app.route("/class/delete/<int:id>", methods=('POST',))
@require_config_set
def class_delete(id):
    dba.delete_class(id)
    return redirect(url_for("veh_class"));

###########
### API ###
###########

@app.route("/api/console", methods=('GET', 'POST'))
@require_config_set
def console():
    data = {
        'is_running': psa.is_running(),
        'text': psa.process_content()
    }
    return jsonify(data);

@app.route("/api/console/start", methods=('GET', 'POST'))
@require_config_set
def console_start():
    id = 2
    server_cfg = util.render_servercfg(dba, id)
    entry_list = util.render_entrylist(dba, id)
    fsa.set_server_ini(server_cfg, entry_list)
    psa.start_server(osa)
    return console()

@app.route("/api/console/stop", methods=('GET', 'POST'))
@require_config_set
def console_stop():
    psa.stop_server()
    return console()

@app.route("/api/server_cfg.ini", methods=('GET',))
@require_config_set
def server_cfg():
    id = request.args.get('id')
    if id is None:
        return "invalid id"
    return util.render_servercfg(dba, id)

@app.route("/api/entry_list.ini", methods=('GET',))
@require_config_set
def entry_list():
    id = request.args.get('id')
    if id is None:
        return "invalid id"
    return util.render_entrylist(dba, id)

@app.route("/api/skin_img/<string:car>/<string:skin>.jpg")
@require_config_set
def skin_img(car, skin):
    return send_from_directory(fsa.get_skin(car, skin), "preview.jpg")

@app.route("/api/track_preview/<string:key>.png")
@app.route("/api/track_preview/<string:key>/<string:config>.png")
@require_config_set
def track_preview(key, config=None):
    return send_from_directory(fsa.get_track(key, config), "preview.png")

@app.route("/api/track_outline/<string:key>.png")
@app.route("/api/track_outline/<string:key>/<string:config>.png")
@require_config_set
def track_layout(key, config=None):
    return send_from_directory(fsa.get_track(key, config), "outline.png")

@app.route("/api/get_vehicle/<string:id>")
@require_config_set
def get_vehicle(id): 
    car_data = dba.get_car(id)

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
        'id': id,
        'desc': car_data['desc'],
        'power': power,
        'torque': torque,
        'labels': labels
    }

    return jsonify(json_data)

@app.route("/api/get_difficulty/<int:id>")
@require_config_set
def get_difficulty(id):
    diff_data = dict(dba.get_difficulty(id))
    return jsonify(diff_data)

@app.route("/api/get_session/<int:id>")
@require_config_set
def get_session(id):
    session_data = dict(dba.get_session(id))
    return jsonify(session_data)

@app.route("/api/get_class/<int:id>")
@require_config_set
def get_class(id):
    class_data = dict(dba.get_class(id))
    class_data['entries'] = [dict(ix) for ix in dba.get_class_entries(id)] 
    return jsonify(class_data)

@app.route("/api/get_time/<int:id>")
@require_config_set
def get_time(id):
    time_data = dict(dba.get_time(id))
    time_data['weathers'] = [dict(ix) for ix in dba.get_time_weather_names(id)] 
    return jsonify(time_data)


@app.route("/api/validate_installpath", methods=('POST', ))
def validate_installpath():
    if fsa.validate_installpath(request.form['path']):
        return jsonify({"result": "ok"})
    else:
        return jsonify({"result": "no"})

@app.route("/api/recache_vehicles")
def recache_vehicles():
    result = fsa.parse_cars_folder(dba)
    return jsonify({'result': 'ok', 'value': result})

@app.route("/api/recache_tracks")
def recache_tracks():
    result = fsa.parse_tracks_folder(dba)
    return jsonify({'result': 'ok', 'value': result})

@app.route("/api/recache_weathers")
def recache_weathers():
    result = fsa.parse_weathers_folder(dba)
    return jsonify({'result': 'ok', 'value': result})


if __name__ == "__main__":
    app.run(debug=True)
