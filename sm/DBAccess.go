package sm

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"

	"github.com/jessevdk/go-assets"
	_ "github.com/mattn/go-sqlite3"
)

type Dbaccess struct {
	Name string
	Db   *sql.DB
}

func Open(name string) Dbaccess {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatal(err)
	}

	dba := Dbaccess{name, db}
	return dba
}

func (Dba Dbaccess) Basepath() string {
	return *Dba.Select_Config().Install_Path
}

func (dba Dbaccess) Apply_Schema(f *assets.File) {
	sqlBytes, err := io.ReadAll(f)
	if err != nil {
		log.Print(err)
	}
	if err != nil {
		log.Print(err)
	}
	sqlStmt := string(sqlBytes)
	_, err = dba.Db.Exec(sqlStmt)
	if err != nil {
		log.Print(err)
		return
	}
}

func (dba Dbaccess) Table_Exists(tablename string) int {
	stmt, err := dba.Db.Prepare("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	var count int
	err = stmt.QueryRow(tablename).Scan(&count)
	if err != nil {
		log.Print(err)
	}

	return count
}

func (dba Dbaccess) Select_DropDownList(filled bool, tableName string) []DropDown_List {
	where := ""
	if filled {
		where = " WHERE filled = 1 AND deleted = 0"
	} else {
		where = " WHERE deleted = 0"
	}
	rows, err := dba.Db.Query("SELECT id, name from " + tableName + where + " ORDER BY name")

	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	ddl := make([]DropDown_List, 0)
	for rows.Next() {
		item := DropDown_List{}
		err = rows.Scan(&item.Id, &item.Name)
		if err != nil {
			log.Print(err)
		}

		ddl = append(ddl, item)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return ddl
}

func (dba Dbaccess) Delete_From(id int, tableName string) int64 {
	stmt, err := dba.Db.Prepare("UPDATE " + tableName + " SET deleted = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) Insert_Name_Into(name string, tableName string) int64 {
	stmt, err := dba.Db.Prepare("INSERT INTO " + tableName + " (name) VALUES (?)")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(name)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
	}

	return id
}

func (dba Dbaccess) Select_User(username string) Users {
	user := Users{}
	stmt, err := dba.Db.Prepare("SELECT name, password, measurement_unit, temp_unit FROM users WHERE name = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.Name, &user.Password, &user.Measurement_Unit, &user.Temp_Unit)
	if err != nil {
		log.Print(err)
	}

	return user
}

func (dba Dbaccess) Update_User(usr Users) int64 {
	stmt, err := dba.Db.Prepare("UPDATE users SET password = ?, measurement_unit = ?, temp_unit = ? WHERE name = ?")

	if err != nil {
		log.Print(err)
	}

	res, err := stmt.Exec(&usr.Password, &usr.Measurement_Unit, &usr.Temp_Unit, &usr.Name)
	defer stmt.Close()
	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) Select_Config_Filled() bool {
	row := dba.Db.QueryRow("SELECT COUNT(*) FROM user_config WHERE cfg_filled = 1 AND mod_filled = 1")

	err := row.Err()
	if err != nil {
		log.Print(err)
	}

	var count int
	err = row.Scan(&count)

	if err != nil {
		log.Print(err)
	}

	return count > 0
}

func (dba Dbaccess) Select_Config() User_Config {
	row := dba.Db.QueryRow("SELECT name, password, admin_password, register_to_lobby, pickup_mode_enabled, locked_entry_list, result_screen_time, udp_port, tcp_port, http_port, client_send_interval, num_threads, max_clients, welcome_message, install_path, csp_required, csp_version, csp_phycars, csp_phytracks, csp_hidepit, cfg_filled, mod_filled FROM user_config")

	err := row.Err()
	if err != nil {
		log.Print(err)
	}

	cfg := User_Config{}
	err = row.Scan(&cfg.Name, &cfg.Password, &cfg.Admin_Password, &cfg.Register_To_Lobby, &cfg.Pickup_Mode_Enabled, &cfg.Locked_Entry_List, &cfg.Result_Screen_Time, &cfg.Udp_Port, &cfg.Tcp_Port, &cfg.Http_Port, &cfg.Client_Send_Interval, &cfg.Num_Threads, &cfg.Max_Clients, &cfg.Welcome_Message, &cfg.Install_Path, &cfg.Csp_Required, &cfg.Csp_Version, &cfg.Csp_Phycars, &cfg.Csp_Phytracks, &cfg.Csp_Hidepit, &cfg.Cfg_Filled, &cfg.Mod_Filled)
	if err != nil {
		log.Print(err)
	}

	return cfg
}

func (dba Dbaccess) Update_Config(cfg User_Config) int64 {
	stmt, err := dba.Db.Prepare("UPDATE user_config SET name = ?, password = ?, admin_password = ?, register_to_lobby = ?, pickup_mode_enabled = ?, locked_entry_list = ?, result_screen_time = ?, udp_port = ?, tcp_port = ?, http_port = ?, client_send_interval = ?, num_threads = ?, max_clients = ?, welcome_message = ?, cfg_filled = 1")

	if err != nil {
		log.Print(err)
	}

	res, err := stmt.Exec(&cfg.Name, &cfg.Password, &cfg.Admin_Password, &cfg.Register_To_Lobby, &cfg.Pickup_Mode_Enabled, &cfg.Locked_Entry_List, &cfg.Result_Screen_Time, &cfg.Udp_Port, &cfg.Tcp_Port, &cfg.Http_Port, &cfg.Client_Send_Interval, &cfg.Num_Threads, &cfg.Max_Clients, &cfg.Welcome_Message)
	defer stmt.Close()
	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) Update_Content(cfg User_Config) int64 {
	stmt, err := dba.Db.Prepare("UPDATE user_config SET csp_required = ?, csp_phycars = ?, csp_phytracks = ?, csp_hidepit = ?, csp_version = ?, install_path = ?, mod_filled = 1")

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(&cfg.Csp_Required, &cfg.Csp_Phycars, &cfg.Csp_Phytracks, &cfg.Csp_Hidepit, &cfg.Csp_Version, &cfg.Install_Path)
	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) Select_Events(started bool, finished bool) []User_Event {

	and := ""
	groupby := " GROUP BY (s.id)"

	if finished {
		and = " AND s.finished = 1"
	} else if started {
		and = " AND s.started_at is not null AND s.finished = 0"
	} else {
		and = " AND s.started_at is null"
	}

	query := `
SELECT
	s.id as id,
	s.race_laps as race_laps,
	s.strategy as strategy,
	t.name as track_name,
	t.length as track_length,
	t.pitboxes as pitboxes,
	d.name as difficulty_name,
	d.abs_allowed as abs_allowed,
	d.tc_allowed as tc_allowed,
	d.stability_allowed as stability_allowed,
	d.autoclutch_allowed as autoclutch_allowed,
	e.name as session_name,
	e.booking_enabled as booking_enabled,
	e.booking_time as booking_time,
	e.practice_enabled as practice_enabled,
	e.practice_time as practice_time,
	e.qualify_enabled as qualify_enabled,
	e.qualify_time as qualify_time,
	e.race_enabled as race_enabled,
	e.race_time as race_time,
	c.name as class_name,
	COUNT(ce.id) as entries,
	tw.name as time_name,
	concat((SELECT
		GROUP_CONCAT(a.csp_time, ', ')
	FROM user_time_weather a
	WHERE user_time_id = tw.id), tw.time) as time,
	(SELECT
		GROUP_CONCAT(b.name, ', ')
	FROM user_time_weather a
	JOIN cache_weather b
		on a.graphics = b.key
	WHERE user_time_id = tw.id) as graphics,
	(SELECT COUNT (*) FROM user_time_weather WHERE user_time_id = tw.id) as trunc_weather,
	tw.csp_enabled as csp_weather,
	s.started_at as started_at
FROM user_event s
JOIN cache_track t
	on s.cache_track_key = t.key
	AND s.cache_track_config = t.config
JOIN user_difficulty d
	on s.difficulty_id = d.id
JOIN user_session e
	on s.session_id = e.id
JOIN user_class c
	on s.class_id = c.id
JOIN user_class_entry ce
	on s.class_id = ce.user_class_id
JOIN user_time tw
	on s.time_id = tw.id
WHERE s.deleted = 0`

	rows, err := dba.Db.Query(query + and + groupby)
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	events := make([]User_Event, 0)
	for rows.Next() {
		evt := User_Event{}
		err = rows.Scan(&evt.Id, &evt.Race_Laps, &evt.Strategy, &evt.Track_Name, &evt.Track_Length, &evt.Pitboxes, &evt.Difficulty_Name, &evt.Abs_Allowed, &evt.Tc_Allowed, &evt.Stability_Allowed, &evt.Autoclutch_Allowed, &evt.Session_Name, &evt.Booking_Enabled, &evt.Booking_Time, &evt.Practice_Enabled, &evt.Practice_Time, &evt.Qualify_Enabled, &evt.Qualify_Time, &evt.Race_Enabled, &evt.Race_Time, &evt.Class_Name, &evt.Entries, &evt.Time_Name, &evt.Time, &evt.Graphics, &evt.TruncWeather, &evt.Csp_Weather, &evt.Started_At)
		if err != nil {
			log.Fatal(err)
		}

		events = append(events, evt)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return events
}

func (dba Dbaccess) Select_Event(id int) User_Event {
	evt := User_Event{}
	stmt, err := dba.Db.Prepare("SELECT id, cache_track_key, cache_track_config, difficulty_id, session_id, class_id, time_id, race_laps, strategy, started_at, finished, servercfg, entrylist FROM user_event WHERE id = ? AND deleted = 0 LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&evt.Id, &evt.Cache_Track_Key, &evt.Cache_Track_Config, &evt.Difficulty_Id, &evt.Session_Id, &evt.Class_Id, &evt.Time_Id, &evt.Race_Laps, &evt.Strategy, &evt.Started_At, &evt.Finished, &evt.ServerCfg, &evt.EntryList)
	if err != nil {
		log.Print(err)
	}

	return evt
}

func (dba Dbaccess) Select_Event_Next() User_Event {
	evt := User_Event{}
	row := dba.Db.QueryRow("SELECT id, cache_track_key, cache_track_config, difficulty_id, session_id, class_id, time_id, race_laps, strategy, started_at, finished, servercfg, entrylist FROM user_event WHERE started_at IS NULL AND finished != 1 AND deleted = 0 LIMIT 1")
	err := row.Err()
	if err != nil {
		log.Print(err)
	}
	err = row.Scan(&evt.Id, &evt.Cache_Track_Key, &evt.Cache_Track_Config, &evt.Difficulty_Id, &evt.Session_Id, &evt.Class_Id, &evt.Time_Id, &evt.Race_Laps, &evt.Strategy, &evt.Started_At, &evt.Finished, &evt.ServerCfg, &evt.EntryList)
	if err != nil {
		log.Print(err)
	}

	return evt
}

func (dba Dbaccess) Insert_Event(evt User_Event) int64 {
	stmt, err := dba.Db.Prepare("INSERT INTO user_event (cache_track_key, cache_track_config, difficulty_id, session_id, class_id, time_id, race_laps, strategy) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(&evt.Cache_Track_Key, &evt.Cache_Track_Config, &evt.Difficulty_Id, &evt.Session_Id, &evt.Class_Id, &evt.Time_Id, &evt.Race_Laps, &evt.Strategy)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) Update_Event(evt User_Event) int64 {
	stmt, err := dba.Db.Prepare("UPDATE user_event SET cache_track_key = ?, cache_track_config = ?, difficulty_id = ?, session_id = ?, class_id = ?, time_id = ?, race_laps = ?, strategy = ?, started_at = ?, servercfg = ?, entrylist = ? WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(evt.Cache_Track_Key, evt.Cache_Track_Config, evt.Difficulty_Id, evt.Session_Id, evt.Class_Id, evt.Time_Id, evt.Race_Laps, evt.Strategy, evt.Started_At, evt.ServerCfg, evt.EntryList, evt.Id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) Update_Event_SetComplete() int64 {
	res, err := dba.Db.Exec("UPDATE user_event SET finished = 1 WHERE started_at IS NOT NULL")
	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) Delete_Event(id int) int64 {
	return dba.Delete_From(id, "user_event")
}

func (dba Dbaccess) Select_Difficulty(id int) User_Difficulty {
	dif := User_Difficulty{}
	stmt, err := dba.Db.Prepare("SELECT id, name, abs_allowed, tc_allowed, stability_allowed, autoclutch_allowed, tyre_blankets_allowed, force_virtual_mirror, fuel_rate, damage_multiplier, tyre_wear_rate, allowed_tyres_out, max_ballast_kg, start_rule, race_gas_penality_disabled, dynamic_track, dynamic_track_preset, session_start, randomness, session_transfer, lap_gain, kick_quorum, vote_duration, voting_quorum, blacklist_mode, max_contacts_per_km FROM user_difficulty WHERE id = ? AND deleted = 0 LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&dif.Id, &dif.Name, &dif.Abs_Allowed, &dif.Tc_Allowed, &dif.Stability_Allowed, &dif.Autoclutch_Allowed, &dif.Tyre_Blankets_Allowed, &dif.Force_Virtual_Mirror, &dif.Fuel_Rate, &dif.Damage_Multiplier, &dif.Tyre_Wear_Rate, &dif.Allowed_Tyres_Out, &dif.Max_Ballast_Kg, &dif.Start_Rule, &dif.Race_Gas_Penality_Disabled, &dif.Dynamic_Track, &dif.Dynamic_Track_Preset, &dif.Session_Start, &dif.Randomness, &dif.Session_Transfer, &dif.Lap_Gain, &dif.Kick_Quorum, &dif.Vote_Duration, &dif.Voting_Quorum, &dif.Blacklist_Mode, &dif.Max_Contacts_Per_Km)
	if err != nil {
		log.Print(err)
	}

	return dif
}

func (dba Dbaccess) Select_DifficultyList(filled bool) []DropDown_List {
	return dba.Select_DropDownList(filled, "user_difficulty")
}

func (dba Dbaccess) Insert_Difficulty(difficulty_name string) int64 {
	return dba.Insert_Name_Into(difficulty_name, "user_difficulty")
}

func (dba Dbaccess) Delete_Difficulty(id int) int64 {
	return dba.Delete_From(id, "user_difficulty")
}

func (dba Dbaccess) Update_Difficulty(dif User_Difficulty) int64 {
	stmt, err := dba.Db.Prepare("UPDATE user_difficulty SET abs_allowed = ?, tc_allowed = ?, stability_allowed = ?, autoclutch_allowed = ?, tyre_blankets_allowed = ?, force_virtual_mirror = ?, fuel_rate = ?, damage_multiplier = ?, tyre_wear_rate = ?, allowed_tyres_out = ?, max_ballast_kg = ?, start_rule = ?, race_gas_penality_disabled = ?, dynamic_track = ?, dynamic_track_preset = ?, session_start = ?, randomness = ?, session_transfer = ?, lap_gain = ?, kick_quorum = ?, voting_quorum = ?, vote_duration = ?, blacklist_mode = ?, max_contacts_per_km = ?, filled = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(&dif.Abs_Allowed, &dif.Tc_Allowed, &dif.Stability_Allowed, &dif.Autoclutch_Allowed, &dif.Tyre_Blankets_Allowed, &dif.Force_Virtual_Mirror, &dif.Fuel_Rate, &dif.Damage_Multiplier, &dif.Tyre_Wear_Rate, &dif.Allowed_Tyres_Out, &dif.Max_Ballast_Kg, &dif.Start_Rule, &dif.Race_Gas_Penality_Disabled, &dif.Dynamic_Track, &dif.Dynamic_Track_Preset, &dif.Session_Start, &dif.Randomness, &dif.Session_Transfer, &dif.Lap_Gain, &dif.Kick_Quorum, &dif.Vote_Duration, &dif.Vote_Duration, &dif.Blacklist_Mode, &dif.Max_Contacts_Per_Km, &dif.Id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) Select_Session(id int) User_Session {
	ses := User_Session{}
	stmt, err := dba.Db.Prepare("SELECT id, name, booking_enabled, booking_time, practice_enabled, practice_time, practice_is_open, qualify_enabled, qualify_time, qualify_is_open, qualify_max_wait_perc, race_enabled, race_time, race_extra_lap, race_over_time, race_wait_time, race_is_open, reversed_grid_positions, race_pit_window_start, race_pit_window_end FROM user_session WHERE id = ? AND deleted = 0 LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	err = stmt.QueryRow(id).Scan(&ses.Id, &ses.Name, &ses.Booking_Enabled, &ses.Booking_Time, &ses.Practice_Enabled, &ses.Practice_Time, &ses.Practice_Is_Open, &ses.Qualify_Enabled, &ses.Qualify_Time, &ses.Qualify_Is_Open, &ses.Qualify_Max_Wait_Perc, &ses.Race_Enabled, &ses.Race_Time, &ses.Race_Extra_Lap, &ses.Race_Over_Time, &ses.Race_Wait_Time, &ses.Race_Is_Open, &ses.Reversed_Grid_Positions, &ses.Race_Pit_Window_Start, &ses.Race_Pit_Window_End)
	defer stmt.Close()
	if err != nil {
		log.Print(err)
	}

	return ses
}

func (dba Dbaccess) Select_SessionList(filled bool) []DropDown_List {
	return dba.Select_DropDownList(filled, "user_session")
}

func (dba Dbaccess) Insert_Session(difficulty_name string) int64 {
	return dba.Insert_Name_Into(difficulty_name, "user_session")
}

func (dba Dbaccess) Delete_Session(id int) int64 {
	return dba.Delete_From(id, "user_session")
}

func (dba Dbaccess) Update_Session(ses User_Session) int64 {
	stmt, err := dba.Db.Prepare("UPDATE user_session SET booking_enabled = ?, booking_time = ?, practice_enabled = ?, practice_time = ?, practice_is_open = ?, qualify_enabled = ?, qualify_time = ?, qualify_is_open = ?, qualify_max_wait_perc = ?, race_enabled = ?, race_time = ?, race_extra_lap = ?, race_over_time = ?, race_wait_time = ?, race_is_open = ?, reversed_grid_positions = ?, race_pit_window_start = ?, race_pit_window_end = ?, filled = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(&ses.Booking_Enabled, &ses.Booking_Time, &ses.Practice_Enabled, &ses.Practice_Time, &ses.Practice_Is_Open, &ses.Qualify_Enabled, &ses.Qualify_Time, &ses.Qualify_Is_Open, &ses.Qualify_Max_Wait_Perc, &ses.Race_Enabled, &ses.Race_Time, &ses.Race_Extra_Lap, &ses.Race_Over_Time, &ses.Race_Wait_Time, &ses.Race_Is_Open, &ses.Reversed_Grid_Positions, &ses.Race_Pit_Window_Start, &ses.Race_Pit_Window_End, &ses.Id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) Select_Time_Weather(id int) User_Time {
	time := User_Time{}
	stmt, err := dba.Db.Prepare("SELECT id, name, time, time_of_day_multi, csp_enabled FROM user_time WHERE id = ? AND deleted = 0 LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&time.Id, &time.Name, &time.Time, &time.Time_Of_Day_Multi, &time.Csp_Enabled)
	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.Db.Prepare("SELECT a.id, name, user_time_id, graphics, base_temperature_ambient, base_temperature_road, variation_ambient, variation_road, wind_base_speed_min, wind_base_speed_max, wind_base_direction, wind_variation_direction, csp_time, csp_time_of_day_multi, csp_date FROM user_time_weather a JOIN cache_weather b on a.graphics = b.key WHERE user_time_id = ? AND deleted = 0")
	if err != nil {
		log.Print(err)
	}
	rows, err := stmt.Query(id)

	time.Weathers = make([]User_Time_Weather, 0)
	wt := User_Time_Weather{}
	for rows.Next() {
		err = rows.Scan(&wt.Id, &wt.Name, &wt.User_Time_Id, &wt.Graphics, &wt.Base_Temperature_Ambient, &wt.Base_Temperature_Road, &wt.Variation_Ambient, &wt.Variation_Road, &wt.Wind_Base_Speed_Min, &wt.Wind_Base_Speed_Max, &wt.Wind_Base_Direction, &wt.Wind_Variation_Direction, &wt.Csp_Time, &wt.Csp_Time_Of_Day_Multi, &wt.Csp_Date)

		time.Weathers = append(time.Weathers, wt)
	}

	return time
}

func (dba Dbaccess) Select_TimeList(filled bool) []DropDown_List {
	return dba.Select_DropDownList(filled, "user_time")
}

func (dba Dbaccess) Insert_Time(time_name string) int64 {
	return dba.Insert_Name_Into(time_name, "user_time")
}

func (dba Dbaccess) Delete_Time(id int) int64 {
	rows := dba.Delete_From(id, "user_time")

	stmt, err := dba.Db.Prepare("UPDATE user_time_weather SET deleted = 1 WHERE user_time_id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	return rows
}

func (dba Dbaccess) Update_Time(time User_Time) int64 {
	stmt, err := dba.Db.Prepare("UPDATE user_time SET time = ?, time_of_day_multi = ?, csp_enabled = ?, filled = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(&time.Time, &time.Time_Of_Day_Multi, &time.Csp_Enabled, &time.Id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.Db.Prepare("DELETE FROM user_time_weather WHERE user_time_id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(time.Id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	for _, w := range time.Weathers {
		stmt, err = dba.Db.Prepare("INSERT INTO user_time_weather (user_time_id, graphics, base_temperature_ambient, base_temperature_road, variation_ambient, variation_road, wind_base_speed_min, wind_base_speed_max, wind_base_direction, wind_variation_direction, csp_time, csp_time_of_day_multi, csp_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&time.Id, &w.Graphics, &w.Base_Temperature_Ambient, &w.Base_Temperature_Road, &w.Variation_Ambient, &w.Variation_Road, &w.Wind_Base_Speed_Min, &w.Wind_Base_Speed_Max, &w.Wind_Base_Direction, &w.Wind_Variation_Direction, &w.Csp_Time, &w.Csp_Time_Of_Day_Multi, &w.Csp_Date)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba Dbaccess) Select_Class_Entries(id int) User_Class {
	cls := User_Class{}
	stmt, err := dba.Db.Prepare("SELECT id, name FROM user_class WHERE id = ? AND deleted = 0 LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&cls.Id, &cls.Name)
	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.Db.Prepare("SELECT id, user_class_id, cache_car_key, skin_key, ballast FROM user_class_entry WHERE user_class_id = ? AND deleted = 0 LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)

	cls.Entries = make([]User_Class_Entry, 0)
	for rows.Next() {
		ent := User_Class_Entry{}
		err = rows.Scan(&ent.Id, &ent.User_Class_Id, &ent.Cache_Car_Key, &ent.Skin_Key, &ent.Ballast)
		if err != nil {
			log.Print(err)
		}
		cls.Entries = append(cls.Entries, ent)
	}

	if err != nil {
		log.Print(err)
	}

	return cls
}

func (dba Dbaccess) Select_ClassList(filled bool) []DropDown_List {
	return dba.Select_DropDownList(filled, "user_class")
}

func (dba Dbaccess) Insert_Class(time_name string) int64 {
	return dba.Insert_Name_Into(time_name, "user_class")
}

func (dba Dbaccess) Delete_Class(id int) int64 {
	dba.Delete_From(id, "user_class")

	stmt, err := dba.Db.Prepare("UPDATE user_class_entry SET deleted = 1 WHERE user_class_id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	return 1
}

func (dba Dbaccess) Update_Class(cls User_Class) int64 {
	stmt, err := dba.Db.Prepare("UPDATE user_class SET filled = 1 WHERE id = ?")

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(cls.Id)
	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.Db.Prepare("DELETE FROM user_class_entry WHERE user_class_id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(cls.Id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	for _, ent := range cls.Entries {
		stmt, err = dba.Db.Prepare("INSERT INTO user_class_entry (user_class_id, cache_car_key, skin_key) VALUES (?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(cls.Id, ent.Cache_Car_Key, ent.Skin_Key)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba Dbaccess) Update_Cache_Cars(cars []Cache_Car) int64 {
	_, err := dba.Db.Exec("DELETE FROM cache_car")
	if err != nil {
		log.Print(err)
	}

	for _, car := range cars {
		stmt, err := dba.Db.Prepare("INSERT INTO cache_car (key, name, brand, desc, tags, class, specs, torque, power, skins) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}

		tagsRes, err := json.Marshal(&car.Tags)
		if err != nil {
			log.Print(err)
		}
		tags := string(tagsRes)

		specsRes, err := json.Marshal(&car.Specs)
		if err != nil {
			log.Print(err)
		}
		specs := string(specsRes)

		powerRes, err := json.Marshal(&car.Power)
		if err != nil {
			log.Print(err)
		}
		power := string(powerRes)

		torqueRes, err := json.Marshal(&car.Torque)
		if err != nil {
			log.Print(err)
		}
		torque := string(torqueRes)

		skinsRes, err := json.Marshal(&car.Skins)
		if err != nil {
			log.Print(err)
		}
		skins := string(skinsRes)
		_, err = stmt.Exec(&car.Key, &car.Name, &car.Brand, &car.Desc, &tags, &car.Class, &specs, &torque, &power, &skins)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba Dbaccess) Select_Cache_Cars() []Cache_Car {
	rows, err := dba.Db.Query("SELECT key, name, brand, desc, tags, class, specs, torque, power, skins FROM cache_car ORDER BY name")
	if err != nil {
		log.Print(err)
	}

	var tags string
	var specs string
	var power string
	var torque string
	var skins string

	cars := make([]Cache_Car, 0)
	for rows.Next() {
		car := Cache_Car{}
		err = rows.Scan(&car.Key, &car.Name, &car.Brand, &car.Desc, &tags, &car.Class, &specs, &torque, &power, &skins)
		if err != nil {
			log.Print(err)
		}
		err = json.Unmarshal([]byte(tags), &car.Tags)
		if err != nil {
			log.Print(err)
		}
		err = json.Unmarshal([]byte(specs), &car.Specs)
		if err != nil {
			log.Print(err)
		}
		err = json.Unmarshal([]byte(power), &car.Power)
		if err != nil {
			log.Print(err)
		}
		err = json.Unmarshal([]byte(torque), &car.Torque)
		if err != nil {
			log.Print(err)
		}
		err = json.Unmarshal([]byte(skins), &car.Skins)
		if err != nil {
			log.Print(err)
		}
		cars = append(cars, car)
	}

	if err != nil {
		log.Print(err)
	}

	return cars
}

func (dba Dbaccess) Select_Cache_Car(car_key string) Cache_Car {
	var tags string
	var specs string
	var power string
	var torque string
	var skins string

	car := Cache_Car{}
	stmt, err := dba.Db.Prepare("SELECT key, name, brand, desc, tags, class, specs, torque, power, skins FROM cache_car WHERE key = ? ORDER BY name")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(car_key).Scan(&car.Key, &car.Name, &car.Brand, &car.Desc, &tags, &car.Class, &specs, &torque, &power, &skins)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal([]byte(tags), &car.Tags)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal([]byte(specs), &car.Specs)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal([]byte(power), &car.Power)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal([]byte(torque), &car.Torque)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal([]byte(skins), &car.Skins)
	if err != nil {
		log.Print(err)
	}

	return car
}

func (dba Dbaccess) Update_Cache_Tracks(tracks []Cache_Track) int64 {
	_, err := dba.Db.Exec("DELETE FROM cache_track")
	if err != nil {
		log.Print(err)
	}

	for _, track := range tracks {

		tagsRes, err := json.Marshal(&track.Tags)
		if err != nil {
			log.Print(err)
		}
		tags := string(tagsRes)

		stmt, err := dba.Db.Prepare("INSERT INTO cache_track (key, config, name, desc, tags, country, city, length, width, pitboxes, run) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&track.Key, &track.Config, &track.Name, &track.Desc, tags, &track.Country, &track.City, &track.Length, &track.Width, &track.Pitboxes, &track.Run)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba Dbaccess) Select_Cache_Tracks() []Cache_Track {
	rows, err := dba.Db.Query("SELECT key, config, name, desc, tags, country, city, length, width, pitboxes, run FROM cache_track ORDER BY name")
	if err != nil {
		log.Print(err)
	}

	var tags string
	tracks := make([]Cache_Track, 0)
	for rows.Next() {
		t := Cache_Track{}
		err = rows.Scan(&t.Key, &t.Config, &t.Name, &t.Desc, &tags, &t.Country, &t.City, &t.Length, &t.Width, &t.Pitboxes, &t.Run)
		if err != nil {
			log.Print(err)
		}

		err = json.Unmarshal([]byte(tags), &t.Tags)
		if err != nil {
			log.Print(err)
		}

		tracks = append(tracks, t)
	}

	if err != nil {
		log.Print(err)
	}

	return tracks
}

func (dba Dbaccess) Select_Cache_Track(track_key string, track_config string) Cache_Track {
	t := Cache_Track{}
	var tags string

	if track_config == "" {
		stmt, err := dba.Db.Prepare("SELECT key, config, name, desc, tags, country, city, length, width, pitboxes, run FROM cache_track WHERE key = ? LIMIT 1")
		if err != nil {
			log.Print(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(track_key).Scan(&t.Key, &t.Config, &t.Name, &t.Desc, &tags, &t.Country, &t.City, &t.Length, &t.Width, &t.Pitboxes, &t.Run)
		if err != nil {
			log.Print(err)
		}
	} else {
		stmt, err := dba.Db.Prepare("SELECT key, config, name, desc, tags, country, city, length, width, pitboxes, run FROM cache_track WHERE key = ? AND config = ? LIMIT 1")
		if err != nil {
			log.Print(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(track_key, track_config).Scan(&t.Key, &t.Config, &t.Name, &t.Desc, &tags, &t.Country, &t.City, &t.Length, &t.Width, &t.Pitboxes, &t.Run)
		if err != nil {
			log.Print(err)
		}
	}

	err := json.Unmarshal([]byte(tags), &t.Tags)
	if err != nil {
		log.Print(err)
	}
	return t
}

func (dba Dbaccess) Update_Cache_Weathers(weathers []Cache_Weather) int64 {
	_, err := dba.Db.Exec("DELETE FROM cache_weather")
	if err != nil {
		log.Print(err)
	}

	for _, w := range weathers {
		stmt, err := dba.Db.Prepare("INSERT INTO cache_weather (key, name) VALUES (?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&w.Key, &w.Name)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba Dbaccess) Select_Cache_Weathers() []Cache_Weather {
	rows, err := dba.Db.Query("SELECT key, name FROM cache_weather ORDER BY name")
	if err != nil {
		log.Print(err)
	}

	weathers := make([]Cache_Weather, 0)
	for rows.Next() {
		w := Cache_Weather{}
		err = rows.Scan(&w.Key, &w.Name)
		if err != nil {
			log.Print(err)
		}
		weathers = append(weathers, w)
	}

	if err != nil {
		log.Print(err)
	}

	return weathers
}

func (dba Dbaccess) Select_Cache_Weather(weather_key string) Cache_Weather {
	w := Cache_Weather{}
	stmt, err := dba.Db.Prepare("SELECT key, name FROM cache_weather WHERE key = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(weather_key).Scan(&w.Key, &w.Name)
	if err != nil {
		log.Print(err)
	}

	return w
}
