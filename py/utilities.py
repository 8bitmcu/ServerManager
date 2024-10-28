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


    def render_servercfg(self, dba, event_id):
        event = dba.get_event(event_id)

        data = {}
        data['event'] = event
        data['config'] = dba.select_config()
        data['diff'] = dba.get_difficulty(event['difficulty_id'])
        data['session'] = dba.get_session(event['session_id'])
        data['time'] = dba.get_time(event['time_id'])
        data['sunangle'] = self.convert_time_to_sunangle(data['time']['time'])
        data['weather'] = dba.get_time_weather(event['time_id'])
        data['class'] = dba.get_class_entries_cache(event['class_id'])
        data['track'] = dba.get_track(event['cache_track_key'], event['cache_track_config'])

        if data['config']['csp_required']:
            cars = data['config']['csp_phycars']
            tracks = data['config']['csp_phytracks']
            pit = data['config']['csp_hidepit']
            version = data['config']['csp_version']

            csp_letter = ""
            if cars and tracks and pit:
                csp_letter = "/../H"
            elif cars and tracks:
                csp_letter = "/../D"
            elif cars and pit:
                csp_letter = "/../F"
            elif tracks and pit:
                csp_letter = "/../G"
            elif cars:
                csp_letter = "/../B"
            elif tracks:
                csp_letter = "/../C"
            elif pit:
                csp_letter = "/../E"

            data['cspstr'] = 'csp/' + str(version) + csp_letter + '/../'


        return render_template("server_cfg.ini", data=data)

    def render_entrylist(self, dba, session_id):
        session = dba.get_session(session_id)
        data = {}
        data['class'] = dba.get_class_entries_cache(session['class_id'])

        return render_template("entry_list.ini", data=data)




