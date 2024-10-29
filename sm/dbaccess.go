package sm

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type dbaccess struct {
	name string
	db   *sql.DB
}

type User_Config struct {
	Name                 *string
	Password             *string
	Admin_Password       *string
	Register_To_Lobby    *int
	Pickup_Mode_Enabled  *int
	Locked_Entry_List    *int
	Result_Screen_Time   *int
	Udp_Port             *int
	Tcp_Port             *int
	Http_Port            *int
	Client_Send_Interval *int
	Num_Threads          *int
	Max_Clients          *int
	Welcome_Message      *string
	Measurement_Unit     *int
	Temp_Unit            *int
	Install_Path         *string
	Csp_Required         *int
	Csp_Version          *int
	Csp_Phycars          *int
	Csp_Phytracks        *int
	Csp_Hidepit          *int
}

type User_Event struct {
	Id                 *int
	Race_Laps          *int
	Strategy           *int
	Track_Name         *string
	Track_Length       *int
	Pitboxes           *int
	Difficulty_Name    *string
	Abs_Allowed        *int
	Tc_Allowed         *int
	Stability_Allowed  *int
	Autoclutch_Allowed *int
	Session_Name       *string
	Booking_Enabled    *int
	Booking_Time       *int
	Practice_Enabled   *int
	Practice_Time      *int
	Qualify_Enabled    *int
	Qualify_Time       *int
	Race_Enabled       *int
	Race_Time          *int
	Class_Name         *string
	Entries            *int
	Time_Name          *string
	Time               *string
	Graphics           *string
	Cache_Track_Key    *string
	Cache_Track_Config *string
	Difficulty_Id      *int
	Session_Id         *int
	Class_Id           *int
	Time_Id            *int
	Started_At         time.Time
	Finished           *int
}

type User_Difficulty struct {
	Id                         *int
	Name                       *string
	Abs_Allowed                *int
	Tc_Allowed                 *int
	Stability_Allowed          *int
	Autoclutch_Allowed         *int
	Tyre_Blankets_Allowed      *int
	Force_Virtual_Mirror       *int
	Fuel_Rate                  *int
	Damage_Multiplier          *int
	Tyre_Wear_Rate             *int
	Allowed_Tyres_Out          *int
	Max_Ballast_Kg             *int
	Start_Rule                 *int
	Race_Gas_Penality_Disabled *int
	Dynamic_Track              *int
	Dynamic_Track_Preset       *int
	Session_Start              *int
	Randomness                 *int
	Session_Transfer           *int
	Lap_Gain                   *int
	Kick_Quorum                *int
	Vote_Duration              *int
	Voting_Quorum              *int
	Blacklist_Mode             *int
	Max_Contacts_Per_Km        *int
}

type User_Session struct {
	Id                      *int
	Name                    *string
	Booking_Enabled         *int
	Booking_Time            *int
	Practice_Enabled        *int
	Practice_Time           *int
	Practice_Is_Open        *int
	Qualify_Enabled         *int
	Qualify_Time            *int
	Qualify_Is_Open         *int
	Qualify_Max_Wait_Perc   *int
	Race_Enabled            *int
	Race_Time               *int
	Race_Extra_Lap          *int
	Race_Over_Time          *int
	Race_Wait_Time          *int
	Race_Is_Open            *int
	Reversed_Grid_Positions *int
	Race_Pit_Window_Start   *int
	Race_Pit_Window_End     *int
}

type User_Time struct {
	Id                *int
	Name              *string
	Time              *string
	Time_of_Day_Multi *int
	Weathers          []User_Time_Weather
}

type User_Time_Weather struct {
	Id                       *int
	User_Time_Id             *int
	Graphics                 *string
	Base_Temperature_Ambient *int
	Base_Temperature_Road    *int
	Variation_Ambient        *int
	Variation_Road           *int
	Wind_Base_Speed_Min      *int
	Wind_Base_Speed_Max      *int
	Wind_Base_Direction      *int
	Wind_Variation_Direction *int
}

type User_Class struct {
	Id      *int
	Name    *string
	Entries []User_Class_Entry
}

type User_Class_Entry struct {
	Id            *int
	User_Class_Id *int
	Cache_Car_Key *string
	Skin_Key      *string
	Ballast       *int
}

type DropDown_List struct {
	Id   *int
	Name *string
}

type Cache_Car struct {
	Id     *int
	Key    *string
	Name   *string
	Brand  *string
	Desc   *string
	Tags   *string
	Class  *string
	Specs  *string
	Torque *string
	Power  *string
	Skins  *string
}

type Cache_Track struct {
	Id       *int
	Key      *string
	Config   *string
	Name     *string
	Desc     *string
	Tags     *string
	Country  *string
	City     *string
	Length   *int
	Width    *int
	Pitboxes *int
	Run      *string
}

type Cache_Weather struct {
	Id   *int
	Key  *string
	Name *string
}

func Open(name string) dbaccess {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatal(err)
	}
	// TODO close db
	//defer db.Close()

	dba := dbaccess{name, db}

	return dba
}

func (dba dbaccess) Apply_Schema() {
	sqlBytes, err := os.ReadFile("schema.sql")
	if err != nil {
		fmt.Print(err)
	}
	sqlStmt := string(sqlBytes)
	_, err = dba.db.Exec(sqlStmt)
	if err != nil {
		log.Print(err)
		return
	}
}

func (dba dbaccess) Table_Exists(tablename string) int {
	stmt, err := dba.db.Prepare("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?")
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

func (dba dbaccess) Select_DropDownList(filled bool, tableName string) []DropDown_List {
	where := ""
	if filled {
		where = " WHERE filled = 1"
	}
	rows, err := dba.db.Query("SELECT id, name from " + tableName + where)

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
			log.Fatal(err)
		}

		ddl = append(ddl, item)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return ddl

}

func (dba dbaccess) Delete_From(id int, tableName string) int64 {
	stmt, err := dba.db.Prepare("DELETE FROM " + tableName + " WHERE id = ?")
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

func (dba dbaccess) Insert_Name_Into(name string, tableName string) int64 {
	stmt, err := dba.db.Prepare("INSERT INTO " + tableName + " (name) VALUES (?)")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(name)
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

func (dba dbaccess) Select_Config() User_Config {
	row := dba.db.QueryRow("SELECT name, password, admin_password, register_to_lobby, pickup_mode_enabled, locked_entry_list, result_screen_time, udp_port, tcp_port, http_port, client_send_interval, num_threads, max_clients, welcome_message, measurement_unit, temp_unit, install_path, csp_required, csp_version, csp_phycars, csp_phytracks, csp_hidepit FROM user_config")

	err := row.Err()
	if err != nil {
		log.Print(err)
	}

	cfg := User_Config{}
	err = row.Scan(&cfg.Name, &cfg.Password, &cfg.Admin_Password, &cfg.Register_To_Lobby, &cfg.Pickup_Mode_Enabled, &cfg.Locked_Entry_List, &cfg.Result_Screen_Time, &cfg.Udp_Port, &cfg.Tcp_Port, &cfg.Http_Port, &cfg.Client_Send_Interval, &cfg.Num_Threads, &cfg.Max_Clients, &cfg.Welcome_Message, &cfg.Measurement_Unit, &cfg.Temp_Unit, &cfg.Install_Path, &cfg.Csp_Required, &cfg.Csp_Version, &cfg.Csp_Phycars, &cfg.Csp_Phytracks, &cfg.Csp_Hidepit)
	if err != nil {
		log.Print(err)
	}

	log.Print(cfg)
	return cfg
}

func (dba dbaccess) Update_Config(cfg User_Config) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_config SET name = ?, password = ?, admin_password = ?, register_to_lobby = ?, pickup_mode_enabled = ?, locked_entry_list = ?, result_screen_time = ?, udp_port = ?, tcp_port = ?, http_port = ?, client_send_interval = ?, num_threads = ?, measurement_unit = ?, temp_unit = ?, install_path = ?, max_clients = ?, welcome_message = ?")

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(&cfg.Name, &cfg.Password, &cfg.Admin_Password, &cfg.Register_To_Lobby, &cfg.Pickup_Mode_Enabled, &cfg.Locked_Entry_List, &cfg.Result_Screen_Time, &cfg.Udp_Port, &cfg.Tcp_Port, &cfg.Http_Port, &cfg.Client_Send_Interval, &cfg.Num_Threads, &cfg.Max_Clients, &cfg.Welcome_Message, &cfg.Measurement_Unit, &cfg.Temp_Unit, &cfg.Install_Path, &cfg.Csp_Required, &cfg.Csp_Version, &cfg.Csp_Phycars, &cfg.Csp_Phytracks, &cfg.Csp_Hidepit)
	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba dbaccess) Update_Content(cfg User_Config) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_config SET csp_required = ?, csp_phycars = ?, csp_phytracks = ?, csp_hidepit = ?, csp_version = ?")

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(&cfg.Csp_Required, &cfg.Csp_Phycars, &cfg.Csp_Phytracks, &cfg.Csp_Hidepit, &cfg.Csp_Version)
	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba dbaccess) Select_Events() []User_Event {
	rows, err := dba.db.Query(`
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
	tw.time as time,
	(SELECT
		GROUP_CONCAT(b.name, ', ')
	FROM user_time_weather a
	JOIN cache_weather b
		on a.graphics = b.key
	WHERE user_time_id = tw.id) as graphics
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
GROUP BY (s.id)
`)

	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	events := make([]User_Event, 0)
	for rows.Next() {
		evt := User_Event{}
		err = rows.Scan(&evt.Id, &evt.Race_Laps, &evt.Strategy, &evt.Track_Name, &evt.Track_Length, &evt.Pitboxes, &evt.Difficulty_Name, &evt.Abs_Allowed, &evt.Tc_Allowed, &evt.Stability_Allowed, &evt.Autoclutch_Allowed, &evt.Session_Name, &evt.Booking_Enabled, &evt.Booking_Time, &evt.Practice_Enabled, &evt.Practice_Time, &evt.Qualify_Enabled, &evt.Qualify_Time, &evt.Race_Enabled, &evt.Race_Time, &evt.Class_Name, &evt.Entries, &evt.Time_Name, &evt.Time, &evt.Graphics)
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

func (dba dbaccess) Select_Event(id int) User_Event {
	evt := User_Event{}
	stmt, err := dba.db.Prepare("SELECT id, cache_track_key, cache_track_config, difficulty_id, session_id, class_id, time_id, race_laps, strategy, started_at, finished FROM user_event LIMIT 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&evt.Id, &evt.Cache_Track_Key, &evt.Cache_Track_Config, &evt.Difficulty_Id, &evt.Session_Id, &evt.Class_Id, &evt.Time_Id, &evt.Race_Laps, &evt.Strategy, &evt.Started_At, &evt.Finished)
	if err != nil {
		log.Print(err)
	}

	return evt
}

func (dba dbaccess) Insert_Event(evt User_Event) int64 {
	stmt, err := dba.db.Prepare("INSERT INTO user_event (cache_track_key, cache_track_config, difficulty_id, session_id, class_id, time_id, race_laps, strategy) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
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

func (dba dbaccess) Update_Event(evt User_Event) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_event SET cache_track_key = ?, cache_track_config = ?, difficulty_id = ?, session_id = ?, class_id = ?, time_id = ?, race_laps = ?, strategy = ? WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(&evt.Cache_Track_Key, &evt.Cache_Track_Config, &evt.Difficulty_Id, &evt.Session_Id, &evt.Class_Id, &evt.Time_Id, &evt.Race_Laps, &evt.Strategy, &evt.Id)
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

func (dba dbaccess) Delete_Event(id int) int64 {
	return dba.Delete_From(id, "user_event")
}

func (dba dbaccess) Select_Difficulty(id int) User_Difficulty {
	dif := User_Difficulty{}
	stmt, err := dba.db.Prepare("SELECT id, name, abs_allowed, tc_allowed, stability_allowed, autoclutch_allowed, tyre_blankets_allowed, force_virtual_mirror, fuel_rate, damage_multiplier, tyre_wear_rate, allowed_tyres_out, max_ballast_kg, start_rule, race_gas_penality_disabled, dynamic_track, dynamic_track_preset, session_start, randomness, session_transfer, lap_gain, kick_quorum, vote_duration, voting_quorum, blacklist_mode, max_contacts_per_km LIMIT 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&dif.Id, &dif.Name, &dif.Abs_Allowed, &dif.Tc_Allowed, &dif.Stability_Allowed, &dif.Autoclutch_Allowed, &dif.Tyre_Blankets_Allowed, &dif.Force_Virtual_Mirror, &dif.Fuel_Rate, &dif.Damage_Multiplier, &dif.Tyre_Wear_Rate, &dif.Allowed_Tyres_Out, &dif.Max_Ballast_Kg, &dif.Start_Rule, &dif.Race_Gas_Penality_Disabled, &dif.Dynamic_Track, &dif.Dynamic_Track_Preset, &dif.Session_Start, &dif.Randomness, &dif.Session_Transfer, &dif.Lap_Gain, &dif.Kick_Quorum, &dif.Vote_Duration, &dif.Voting_Quorum, &dif.Tyre_Blankets_Allowed, &dif.Max_Contacts_Per_Km)
	if err != nil {
		log.Print(err)
	}

	return dif
}

func (dba dbaccess) Select_DifficultyList(filled bool) []DropDown_List {
	return dba.Select_DropDownList(filled, "user_difficulty")
}

func (dba dbaccess) Insert_Difficulty(difficulty_name string) int64 {
	return dba.Insert_Name_Into(difficulty_name, "user_difficulty")
}

func (dba dbaccess) Delete_Difficulty(id int) int64 {
	return dba.Delete_From(id, "user_difficulty")
}

func (dba dbaccess) Update_Difficulty(dif User_Difficulty) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_difficulty SET abs_allowed = ?, tc_allowed = ?, stability_allowed = ?, autoclutch_allowed = ?, tyre_blankets_allowed = ?, force_virtual_mirror = ?, fuel_rate = ?, damage_multiplier = ?, tyre_wear_rate = ?, allowed_tyres_out = ?, max_ballast_kg = ?, start_rule = ?, race_gas_penality_disabled = ?, dynamic_track = ?, dynamic_track_preset = ?, session_start = ?, randomness = ?, session_transfer = ?, lap_gain = ?, kick_quorum = ?, voting_quorum = ?, vote_duration = ?, blacklist_mode = ?, max_contacts_per_km = ?, filled = 1 WHERE id = ?")
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

func (dba dbaccess) Select_Session(id int) User_Session {
	ses := User_Session{}
	stmt, err := dba.db.Prepare("SELECT id, name, booking_enabled, booking_time, practice_enabled, practice_time, practice_is_open, qualify_enabled, qualify_time, qualify_is_open, qualify_max_wait_perc, race_enabled, race_time, race_extra_lap, race_over_time, race_wait_time, race_is_open, reversed_grid_positions, race_pit_window_start, race_pit_window_end FROM user_session WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&ses.Id, &ses.Name, &ses.Booking_Enabled, &ses.Booking_Time, &ses.Practice_Enabled, &ses.Practice_Is_Open, &ses.Qualify_Enabled, &ses.Qualify_Time, &ses.Qualify_Is_Open, &ses.Qualify_Max_Wait_Perc, &ses.Race_Time, &ses.Race_Extra_Lap, &ses.Race_Over_Time, &ses.Race_Wait_Time, &ses.Race_Is_Open, &ses.Reversed_Grid_Positions, &ses.Race_Pit_Window_Start, &ses.Race_Pit_Window_End)
	if err != nil {
		log.Print(err)
	}

	return ses
}

func (dba dbaccess) Select_SessionList(filled bool) []DropDown_List {
	return dba.Select_DropDownList(filled, "user_session")
}

func (dba dbaccess) Insert_Session(difficulty_name string) int64 {
	return dba.Insert_Name_Into(difficulty_name, "user_session")
}

func (dba dbaccess) Delete_Session(id int) int64 {
	return dba.Delete_From(id, "user_session")
}

func (dba dbaccess) Update_Session(ses User_Session) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_session SET booking_enabled = ?, booking_time = ?, practice_enabled = ?, practice_time = ?, practice_is_open = ?, qualify_enabled = ?, qualify_time = ?, qualify_is_open = ?, qualify_max_wait_perc = ?, race_enabled = ?, race_time = ?, race_extra_lap = ?, race_over_time = ?, race_wait_time = ?, race_is_open = ?, reversed_grid_positions = ?, race_pit_window_start = ?, race_pit_window_end = ?, filled = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(&ses.Booking_Enabled, &ses.Booking_Time, &ses.Practice_Enabled, &ses.Practice_Time, &ses.Practice_Is_Open, &ses.Qualify_Enabled, &ses.Qualify_Time, &ses.Qualify_Is_Open, &ses.Qualify_Max_Wait_Perc, &ses.Race_Enabled, &ses.Race_Time, &ses.Race_Extra_Lap, &ses.Race_Is_Open, &ses.Race_Wait_Time, &ses.Race_Is_Open, &ses.Reversed_Grid_Positions, &ses.Race_Pit_Window_Start, &ses.Race_Pit_Window_End, &ses.Id)
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

func (dba dbaccess) Select_Time(id int) User_Time {
	time := User_Time{}
	stmt, err := dba.db.Prepare("SELECT id, name, time, time_of_day_multi FROM user_time WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&time.Id, &time.Name, &time.Time, &time.Time_of_Day_Multi)
	if err != nil {
		log.Print(err)
	}

	return time
}

func (dba dbaccess) Select_Time_Weather(id int) User_Time_Weather {
	wt := User_Time_Weather{}
	stmt, err := dba.db.Prepare("SELECT id, user_time_id, graphics, base_temperature_ambient, base_temperature_road, variation_ambient, variation_road, wind_base_speed_min, wind_base_speed_max, wind_base_direction, wind_variation_direction FROM user_time_weather WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&wt.Id, &wt.User_Time_Id, &wt.Graphics, &wt.Base_Temperature_Ambient, &wt.Base_Temperature_Road, &wt.Variation_Ambient, &wt.Variation_Road, &wt.Wind_Base_Speed_Min, &wt.Wind_Base_Speed_Max, &wt.Wind_Base_Direction, &wt.Wind_Variation_Direction)
	if err != nil {
		log.Print(err)
	}

	return wt
}

func (dba dbaccess) Select_TimeList(filled bool) []DropDown_List {
	return dba.Select_DropDownList(filled, "user_time")
}

func (dba dbaccess) Insert_Time(time_name string) int64 {
	return dba.Insert_Name_Into(time_name, "user_time")
}

func (dba dbaccess) Delete_Time(id int) int64 {
	return dba.Delete_From(id, "user_time")
}

func (dba dbaccess) Update_Time(time User_Time) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_time SET time = ?, time_of_day_multi = ?, filled = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(&time.Time, &time.Time_of_Day_Multi, &time.Id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.db.Prepare("DELETE FROM user_time_weather WHERE user_time_id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(time.Id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	for _, w := range time.Weathers {
		stmt, err = dba.db.Prepare("INSERT INTO user_time_weather (user_time_id, graphics, base_temperature_ambient, base_temperature_road, variation_ambient, variation_road, wind_base_speed_min, wind_base_speed_max, wind_base_direction, wind_variation_direction) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&time.Id, &w.Graphics, &w.Base_Temperature_Ambient, &w.Base_Temperature_Road, &w.Variation_Ambient, &w.Variation_Road, &w.Wind_Base_Speed_Min, &w.Wind_Base_Speed_Max, &w.Wind_Base_Direction, &w.Wind_Variation_Direction)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba dbaccess) Select_Class(id int) User_Class {
	cls := User_Class{}
	stmt, err := dba.db.Prepare("SELECT id, name FROM user_class WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&cls.Id, &cls.Name)
	if err != nil {
		log.Print(err)
	}

	return cls
}

func (dba dbaccess) Select_Class_Entries(id int) []User_Class_Entry {
	stmt, err := dba.db.Prepare("SELECT id, user_class_id, cache_car_key, skin_key, ballast FROM user_class_entry WHERE user_class_id = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)

	entries := make([]User_Class_Entry, 0)
	for rows.Next() {
		ent := User_Class_Entry{}
		err = rows.Scan(&ent.Id, &ent.User_Class_Id, &ent.Cache_Car_Key, &ent.Skin_Key, &ent.Ballast)
		if err != nil {
			log.Print(err)
		}
		entries = append(entries, ent)
	}

	if err != nil {
		log.Print(err)
	}

	return entries
}

func (dba dbaccess) Select_ClassList(filled bool) []DropDown_List {
	return dba.Select_DropDownList(filled, "user_class")
}

func (dba dbaccess) Insert_Class(time_name string) int64 {
	return dba.Insert_Name_Into(time_name, "user_class")
}

func (dba dbaccess) Delete_Class(id int) int64 {
	dba.Delete_From(id, "user_class")

	stmt, err := dba.db.Prepare("DELETE FROM user_class_entry WHERE user_class_id = ?")
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

func (dba dbaccess) Update_Class(cls User_Class) int64 {
	stmt, err := dba.db.Prepare("DELETE FROM user_class_entry WHERE user_class_id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(cls.Id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	for _, ent := range cls.Entries {
		stmt, err = dba.db.Prepare("INSERT INTO user_class_entry (user_class_id, cache_car_key, skin_key) VALUES (?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&ent.User_Class_Id, &ent.Cache_Car_Key, &ent.Skin_Key)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba dbaccess) Update_Cache_Cars(cars []Cache_Car) int64 {
	_, err := dba.db.Query("DELETE FROM cache_car")
	if err != nil {
		log.Print(err)
	}

	for _, car := range cars {
		stmt, err := dba.db.Prepare("INSERT INTO cache_vehicle (key, name, brand, desc, tags, class, specs, torque, power, skins) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&car.Key, &car.Name, &car.Brand, &car.Desc, &car.Tags, &car.Class, &car.Specs, &car.Torque, &car.Power, &car.Skins)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba dbaccess) Select_Cache_Cars() []Cache_Car {
	rows, err := dba.db.Query("SELECT key, name, brand, desc, tags, class, specs, torque, power, skins FROM cache_car")
	if err != nil {
		log.Print(err)
	}

	cars := make([]Cache_Car, 0)
	for rows.Next() {
		car := Cache_Car{}
		err = rows.Scan(&car.Key, &car.Name, &car.Brand, &car.Desc, &car.Tags, &car.Class, &car.Specs, &car.Torque, &car.Power, &car.Skins)
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

func (dba dbaccess) Select_Cache_Car(car_key string) Cache_Car {
	car := Cache_Car{}
	stmt, err := dba.db.Prepare("SELECT key, name, brand, desc, tags, class, specs, torque, power, skins FROM cache_car WHERE key = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(car_key).Scan(&car.Key, &car.Name, &car.Brand, &car.Desc, &car.Tags, &car.Class, &car.Specs, &car.Torque, &car.Power, &car.Skins)
	if err != nil {
		log.Print(err)
	}

	return car
}

func (dba dbaccess) Update_Cache_Tracks(tracks []Cache_Track) int64 {
	_, err := dba.db.Query("DELETE FROM cache_track")
	if err != nil {
		log.Print(err)
	}

	for _, track := range tracks {
		stmt, err := dba.db.Prepare("INSERT INTO cache_track (key, config, name, desc, tags, country, city, length, width, pitboxes, run) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&track.Key, &track.Config, &track.Name, &track.Desc, &track.Tags, &track.Country, &track.City, &track.Length, &track.Width, &track.Pitboxes, &track.Run)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba dbaccess) Select_Cache_Tracks() []Cache_Track {
	rows, err := dba.db.Query("SELECT key, config, name, desc, tags, country, city, length, width, pitboxes, run FROM cache_track")
	if err != nil {
		log.Print(err)
	}

	tracks := make([]Cache_Track, 0)
	for rows.Next() {
		t := Cache_Track{}
		err = rows.Scan(&t.Key, &t.Config, &t.Name, &t.Desc, &t.Tags, &t.Country, &t.City, &t.Length, &t.Width, &t.Pitboxes, &t.Run)
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

func (dba dbaccess) Select_Cache_Track(track_key string, track_config string) Cache_Track {
	t := Cache_Track{}

	if track_config == "" {
		stmt, err := dba.db.Prepare("SELECT key, config, name, desc, tags, country, city, length, width, pitboxes, run FROM cache_track WHERE key = ?")
		if err != nil {
			log.Print(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(track_key).Scan(&t.Key, &t.Config, &t.Name, &t.Desc, &t.Tags, &t.Country, &t.City, &t.Length, &t.Width, &t.Pitboxes, &t.Run)
		if err != nil {
			log.Print(err)
		}
	} else {
		stmt, err := dba.db.Prepare("SELECT key, config, name, desc, tags, country, city, length, width, pitboxes, run FROM cache_track WHERE key = ? AND config = ?")
		if err != nil {
			log.Print(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(track_key, track_config).Scan(&t.Key, &t.Config, &t.Name, &t.Desc, &t.Tags, &t.Country, &t.City, &t.Length, &t.Width, &t.Pitboxes, &t.Run)
		if err != nil {
			log.Print(err)
		}

	}

	return t
}

func (dba dbaccess) Update_Cache_Weathers(weathers []Cache_Weather) int64 {
	_, err := dba.db.Query("DELETE FROM cache_weather")
	if err != nil {
		log.Print(err)
	}

	for _, w := range weathers {
		stmt, err := dba.db.Prepare("INSERT INTO cache_weather (key, name) VALUES (?, ?)")
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

func (dba dbaccess) Select_Cache_Weathers() []Cache_Weather {
	rows, err := dba.db.Query("SELECT key, name FROM cache_track")
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

func (dba dbaccess) Select_Cache_Weather(weather_key string) Cache_Weather {
	w := Cache_Weather{}
	stmt, err := dba.db.Prepare("SELECT key, name FROM cache_weather WHERE key = ?")
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
