from datetime import datetime
from flask import render_template

class Utilities: 
    def __init__(self, ):
        pass

    def clean_dbdata(self, dbdata):
        newdata = {}
        if dbdata is None:
            return None
        for i in dict(dbdata):
            if dbdata[i] is not None:
                newdata[i] = dbdata[i]
        return newdata


# 8:00 AM = -80
# 18:00 PM = 80
# increment of 8 every 30 minutes
    def convert_time_to_sunangle(self, time_str):
        dt = datetime.strptime(time_str,'%H:%M')
        angle = -80 + (16 * (dt.hour - 8))
        angle = angle + round(dt.minute/15)*4

        return angle


    def render_servercfg(self, dba, session_id):
        session = dba.get_session(session_id)
        data = {}
        # TODO refactor into single query ?
        data['session'] = session
        data['config'] = dba.select_config()
        data['diff'] = dba.get_difficulty(session['difficulty_id'])
        data['event'] = dba.get_event(session['event_id'])
        data['time'] = dba.get_time(session['time_id'])
        data['sunangle'] = self.convert_time_to_sunangle(data['time']['time'])
        data['weather'] = dba.get_weather(session['time_id'])
        data['class'] = dba.get_class_entries_cache(session['class_id'])
        data['track'] = dba.get_track(session['cache_track_id'])

        return render_template("server_cfg.ini", data=data)

    def render_entrylist(self, dba, session_id):
        session = dba.get_session(session_id)
        data = {}
        data['class'] = dba.get_class_entries_cache(session['class_id'])

        return render_template("entry_list.ini", data=data)
