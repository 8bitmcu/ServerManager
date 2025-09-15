package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"

	"github.com/jessevdk/go-assets"
	_ "github.com/mattn/go-sqlite3"
)

type Dbaccess struct {
	name string
	db   *sql.DB
}

func open(name string) Dbaccess {
	db, err := sql.Open("sqlite3", "file:"+name+"?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}

	dba := Dbaccess{name, db}
	return dba
}

func (dba Dbaccess) basepath() string {
	return *dba.selectConfig().InstallPath
}

func (dba Dbaccess) applySchema(f *assets.File) {
	sqlBytes, err := io.ReadAll(f)
	if err != nil {
		log.Print(err)
	}
	if err != nil {
		log.Print(err)
	}
	sqlStmt := string(sqlBytes)
	_, err = dba.db.Exec(sqlStmt)
	if err != nil {
		log.Print(err)
		return
	}
}

func (dba Dbaccess) tableExists(tablename string) int {
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

func (dba Dbaccess) selectDropDownList(filled bool, tableName string) []DropDownList {
	where := ""
	if filled {
		where = " WHERE filled = 1"
	}
	rows, err := dba.db.Query("SELECT id, name from " + tableName + where + " ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	ddl := make([]DropDownList, 0)
	for rows.Next() {
		item := DropDownList{}
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

func (dba Dbaccess) deleteFrom(id int, tableName string) (int64, error) {
	stmt, err := dba.db.Prepare("DELETE FROM " + tableName + " WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected, nil
}

func (dba Dbaccess) insertNameInto(name string, tableName string) int64 {
	stmt, err := dba.db.Prepare("INSERT INTO " + tableName + " (name) VALUES (?)")
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

func (dba Dbaccess) selectUser(username string) Users {
	user := Users{}
	stmt, err := dba.db.Prepare("SELECT name, password, measurement_unit, temp_unit FROM users WHERE name = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.Name, &user.Password, &user.MeasurementUnit, &user.TempUnit)
	if err != nil {
		log.Print(err)
	}

	return user
}

func (dba Dbaccess) updateUser(usr Users) int64 {
	stmt, err := dba.db.Prepare("UPDATE users SET password = ?, measurement_unit = ?, temp_unit = ? WHERE name = ?")

	if err != nil {
		log.Print(err)
	}

	res, err := stmt.Exec(&usr.Password, &usr.MeasurementUnit, &usr.TempUnit, &usr.Name)
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

func (dba Dbaccess) selectConfigFilled() bool {
	row := dba.db.QueryRow("SELECT COUNT(*) FROM user_config WHERE cfg_filled = 1 AND mod_filled = 1")

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

func (dba Dbaccess) selectConfig() UserConfig {
	row := dba.db.QueryRow("SELECT name, password, admin_password, register_to_lobby, locked_entry_list, result_screen_time, udp_port, tcp_port, http_port, client_send_interval, num_threads, max_clients, welcome_message, append_eventname, append_modlinks, install_path, csp_required, csp_version, csp_phycars, csp_phytracks, csp_hidepit, cfg_filled, mod_filled, secret_key FROM user_config")

	err := row.Err()
	if err != nil {
		log.Print(err)
	}

	cfg := UserConfig{}
	err = row.Scan(&cfg.Name, &cfg.Password, &cfg.AdminPassword, &cfg.RegisterToLobby, &cfg.LockedEntryList, &cfg.ResultScreenTime, &cfg.UdpPort, &cfg.TcpPort, &cfg.HttpPort, &cfg.ClientSendInterval, &cfg.NumThreads, &cfg.MaxClients, &cfg.WelcomeMessage, &cfg.AppendEventname, &cfg.AppendModlinks, &cfg.InstallPath, &cfg.CspRequired, &cfg.CspVersion, &cfg.CspPhycars, &cfg.CspPhytracks, &cfg.CspHidepit, &cfg.CfgFilled, &cfg.ModFilled, &cfg.SecretKey)
	if err != nil {
		log.Print(err)
	}

	return cfg
}

func (dba Dbaccess) updateConfig(cfg UserConfig) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_config SET name = ?, append_eventname = ?, password = ?, admin_password = ?, register_to_lobby = ?, locked_entry_list = ?, result_screen_time = ?, udp_port = ?, tcp_port = ?, http_port = ?, client_send_interval = ?, num_threads = ?, max_clients = ?, welcome_message = ?, append_modlinks = ?, cfg_filled = 1")

	if err != nil {
		log.Print(err)
	}

	res, err := stmt.Exec(&cfg.Name, &cfg.AppendEventname, &cfg.Password, &cfg.AdminPassword, &cfg.RegisterToLobby, &cfg.LockedEntryList, &cfg.ResultScreenTime, &cfg.UdpPort, &cfg.TcpPort, &cfg.HttpPort, &cfg.ClientSendInterval, &cfg.NumThreads, &cfg.MaxClients, &cfg.WelcomeMessage, &cfg.AppendModlinks)
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

func (dba Dbaccess) updateContent(cfg UserConfig) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_config SET csp_required = ?, csp_phycars = ?, csp_phytracks = ?, csp_hidepit = ?, csp_version = ?, install_path = ?, mod_filled = 1")

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(&cfg.CspRequired, &cfg.CspPhycars, &cfg.CspPhytracks, &cfg.CspHidepit, &cfg.CspVersion, &cfg.InstallPath)
	if err != nil {
		log.Print(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}

	return affected
}

func (dba Dbaccess) selectServerEvents(notfinished bool) []ServerEvent {

	orderby := " ORDER BY orderby ASC"

	where := ""
	if notfinished {
		where = " WHERE finished = 0"
	}

	rows, err := dba.db.Query(`
SELECT
	s.id as id,
	u.id as event_id,
	t.name as track_name,
	d.name as difficulty_name,
	e.name as session_name,
	c.name as class_name,
	tw.name as time_name,
	ct.name as category_name,
	s.started_at as started_at,
	s.finished as finished
FROM server_event s
JOIN user_event u
	on s.user_event_id = u.id
JOIN user_event_category ct
	on u.event_category_id = ct.id
JOIN cache_track t
	on u.cache_track_key = t.key
	AND u.cache_track_config = t.config
JOIN user_difficulty d
	on u.difficulty_id = d.id
JOIN user_session e
	on u.session_id = e.id
JOIN user_class c
	on u.class_id = c.id
JOIN user_time tw
	on u.time_id = tw.id` + where + orderby)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	list := make([]ServerEvent, 0)
	for rows.Next() {
		se := ServerEvent{}
		err = rows.Scan(&se.Id, &se.UserEvent.Id, &se.UserEvent.TrackName, &se.UserEvent.DifficultyName, &se.UserEvent.SessionName, &se.UserEvent.ClassName, &se.UserEvent.TimeName, &se.UserEvent.CategoryName, &se.StartedAt, &se.Finished)
		if err != nil {
			log.Print(err)
		}

		list = append(list, se)
	}

	return list
}

func (dba Dbaccess) insertServerEvent(event int) int64 {
	sql := "INSERT INTO server_event (user_event_id, orderby) SELECT ?, (SELECT ifnull(MAX(orderby)+1, 1) FROM server_event)"

	stmt, err := dba.db.Prepare(sql)
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(event)
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

func (dba Dbaccess) insertServerEventCategory(category int) int64 {
	stmt, err := dba.db.Prepare("SELECT id FROM user_event WHERE event_category_id = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(category)
	if err != nil {
		log.Print(err)
	}

	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Print(err)
		}
		ids = append(ids, id)
	}

	for _, id := range ids {
		stmt, err := dba.db.Prepare("INSERT INTO server_event (user_event_id, orderby) VALUES (?, (SELECT ifnull(MAX(orderby)+1, 1) FROM server_event))")
		if err != nil {
			log.Print(err)
		}
		res, err := stmt.Exec(id)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}

		_, err = res.RowsAffected()
		if err != nil {
			log.Print(err)
		}
	}

	if err != nil {
		log.Print(err)
	}

	return 0
}

func (dba Dbaccess) updateServerEvent(se ServerEvent) int64 {
	stmt, err := dba.db.Prepare("UPDATE server_event SET started_at = ?, servercfg = ?, entrylist = ?, finished = ? WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(se.StartedAt, se.ServerCfg, se.EntryList, se.Finished, se.Id)
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

func (dba Dbaccess) updateServerEventMoveUp(id int) {
	stmt, err := dba.db.Prepare("SELECT id, MAX(orderby) as orderby FROM server_event WHERE orderby < (SELECT orderby FROM server_event WHERE id = ?)")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	var oldid int
	var orderby int
	err = stmt.QueryRow(id).Scan(&oldid, &orderby)

	if err != nil {
		log.Print(err)
	}

	log.Print("old id", oldid)
	log.Print("id", id)

	stmt, err = dba.db.Prepare("UPDATE server_event SET orderby = ? WHERE id = ?;")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(orderby, id)
	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.db.Prepare("UPDATE server_event SET orderby = ? WHERE id = ?;")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	res, err := stmt.Exec(orderby+1, oldid)
	if err != nil {
		log.Print(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Print(err)
	}
}

func (dba Dbaccess) updateServerEventMoveDown(id int) {
	stmt, err := dba.db.Prepare("SELECT id, MIN(orderby) as orderby FROM server_event WHERE orderby > (SELECT orderby FROM server_event WHERE id = ?)")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	var oldid int
	var orderby int
	err = stmt.QueryRow(id).Scan(&oldid, &orderby)

	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.db.Prepare("UPDATE server_event SET orderby = ? WHERE id = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(orderby, id)
	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.db.Prepare("UPDATE server_event SET orderby = ? WHERE id = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(orderby-1, oldid)
	if err != nil {
		log.Print(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Print(err)
	}
}

func (dba Dbaccess) deleteServerEventsCompleted() (int64, error) {
	stmt, err := dba.db.Prepare("DELETE FROM server_event WHERE finished = 1")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec()
	defer stmt.Close()

	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (dba Dbaccess) deleteServerEvent(id int) (int64, error) {
	stmt, err := dba.db.Prepare("DELETE FROM server_event WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (dba Dbaccess) selectEvent(id int) UserEvent {
	evt := UserEvent{}
	stmt, err := dba.db.Prepare("SELECT id, event_category_id, cache_track_key, cache_track_config, difficulty_id, session_id, class_id, time_id, race_laps, strategy FROM user_event WHERE id = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&evt.Id, &evt.EventCategoryId, &evt.CacheTrackKey, &evt.CacheTrackConfig, &evt.DifficultyId, &evt.SessionId, &evt.ClassId, &evt.TimeId, &evt.RaceLaps, &evt.Strategy)
	if err != nil {
		log.Print(err)
	}

	return evt
}

func (dba Dbaccess) selectEventList() []UserEventList {
	rows, err := dba.db.Query("SELECT s.id, t.name, s.event_category_id from user_event s JOIN cache_track t on s.cache_track_key = t.key AND s.cache_track_config = t.config")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	ddl := make([]UserEventList, 0)
	for rows.Next() {
		item := UserEventList{}
		err = rows.Scan(&item.Id, &item.TrackName, &item.EventCategoryId)
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

func (dba Dbaccess) insertEvent(evt UserEvent) int64 {
	stmt, err := dba.db.Prepare("INSERT INTO user_event (event_category_id, cache_track_key, cache_track_config, difficulty_id, session_id, class_id, time_id, race_laps, strategy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(&evt.EventCategoryId, &evt.CacheTrackKey, &evt.CacheTrackConfig, &evt.DifficultyId, &evt.SessionId, &evt.ClassId, &evt.TimeId, &evt.RaceLaps, &evt.Strategy)
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

func (dba Dbaccess) updateEvent(evt UserEvent) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_event SET cache_track_key = ?, cache_track_config = ?, difficulty_id = ?, session_id = ?, class_id = ?, time_id = ?, race_laps = ?, strategy = ? WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(evt.CacheTrackKey, evt.CacheTrackConfig, evt.DifficultyId, evt.SessionId, evt.ClassId, evt.TimeId, evt.RaceLaps, evt.Strategy, evt.Id)
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

func (dba Dbaccess) deleteEvent(id int) (int64, error) {
	return dba.deleteFrom(id, "user_event")
}

func (dba Dbaccess) selectEventCategory(id int) UserEventCategory {
	cat := UserEventCategory{}
	stmt, err := dba.db.Prepare("SELECT id, name FROM user_event_category WHERE id = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&cat.Id, &cat.Name)
	if err != nil {
		log.Print(err)
	}

	return cat
}

func (dba Dbaccess) selectCategoryEvents(id int) UserEventCategory {
	cat := dba.selectEventCategory(id)

	query := `
SELECT
	s.id as id,
	s.race_laps as race_laps,
	s.strategy as strategy,
	t.name as track_name,
	t.length as track_length,
	t.pitboxes as pitboxes,
	d.id as difficulty_id,
	d.name as difficulty_name,
	d.abs_allowed as abs_allowed,
	d.tc_allowed as tc_allowed,
	d.stability_allowed as stability_allowed,
	d.autoclutch_allowed as autoclutch_allowed,
	e.id as session_id,
	e.name as session_name,
	e.booking_enabled as booking_enabled,
	e.booking_time as booking_time,
	e.practice_enabled as practice_enabled,
	e.practice_time as practice_time,
	e.qualify_enabled as qualify_enabled,
	e.qualify_time as qualify_time,
	e.race_enabled as race_enabled,
	e.race_time as race_time,
	c.id as class_id,
	c.name as class_name,
	COUNT(ce.id) as entries,
	tw.id as time_id,
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
	tw.csp_enabled as csp_weather
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
WHERE s.event_category_id = ?
GROUP BY s.id`

	rows, err := dba.db.Query(query, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		evt := UserEvent{}
		err = rows.Scan(&evt.Id, &evt.RaceLaps, &evt.Strategy, &evt.TrackName, &evt.TrackLength, &evt.Pitboxes, &evt.DifficultyId, &evt.DifficultyName, &evt.AbsAllowed, &evt.TcAllowed, &evt.StabilityAllowed, &evt.AutoclutchAllowed, &evt.SessionId, &evt.SessionName, &evt.BookingEnabled, &evt.BookingTime, &evt.PracticeEnabled, &evt.PracticeTime, &evt.QualifyEnabled, &evt.QualifyTime, &evt.RaceEnabled, &evt.RaceTime, &evt.ClassId, &evt.ClassName, &evt.Entries, &evt.TimeId, &evt.TimeName, &evt.Time, &evt.Graphics, &evt.TruncWeather, &evt.CspWeather)
		if err != nil {
			log.Fatal(err)
		}

		cat.Events = append(cat.Events, evt)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return cat
}

func (dba Dbaccess) selectEventCategoryList(filled bool) []DropDownList {
	return dba.selectDropDownList(filled, "user_event_category")
}

func (dba Dbaccess) insertEventCategory(categoryname string) int64 {
	return dba.insertNameInto(categoryname, "user_event_category")
}

func (dba Dbaccess) updateEventCategory(cat UserEventCategory) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_event_category SET name = ?, filled = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(&cat.Name, &cat.Id)
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

func (dba Dbaccess) deleteEventCategory(id int) (int64, error) {
	return dba.deleteFrom(id, "user_event_category")
}

func (dba Dbaccess) selectDifficulty(id int) UserDifficulty {
	dif := UserDifficulty{}
	stmt, err := dba.db.Prepare("SELECT id, name, abs_allowed, tc_allowed, stability_allowed, autoclutch_allowed, tyre_blankets_allowed, force_virtual_mirror, fuel_rate, damage_multiplier, tyre_wear_rate, allowed_tyres_out, max_ballast_kg, start_rule, race_gas_penality_disabled, dynamic_track, dynamic_track_preset, session_start, randomness, session_transfer, lap_gain, kick_quorum, vote_duration, voting_quorum, blacklist_mode, max_contacts_per_km FROM user_difficulty WHERE id = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&dif.Id, &dif.Name, &dif.AbsAllowed, &dif.TcAllowed, &dif.StabilityAllowed, &dif.AutoclutchAllowed, &dif.TyreBlanketsAllowed, &dif.ForceVirtualMirror, &dif.FuelRate, &dif.DamageMultiplier, &dif.TyreWearRate, &dif.AllowedTyresOut, &dif.MaxBallastKg, &dif.StartRule, &dif.RaceGasPenalityDisabled, &dif.DynamicTrack, &dif.DynamicTrackPreset, &dif.SessionStart, &dif.Randomness, &dif.SessionTransfer, &dif.LapGain, &dif.KickQuorum, &dif.VoteDuration, &dif.VotingQuorum, &dif.BlacklistMode, &dif.MaxContactsPerKm)
	if err != nil {
		log.Print(err)
	}

	return dif
}

func (dba Dbaccess) selectDifficultyList(filled bool) []DropDownList {
	return dba.selectDropDownList(filled, "user_difficulty")
}

func (dba Dbaccess) insertDifficulty(difficultyname string) int64 {
	return dba.insertNameInto(difficultyname, "user_difficulty")
}

func (dba Dbaccess) deleteDifficulty(id int) (int64, error) {
	return dba.deleteFrom(id, "user_difficulty")
}

func (dba Dbaccess) updateDifficulty(dif UserDifficulty) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_difficulty SET name = ?, abs_allowed = ?, tc_allowed = ?, stability_allowed = ?, autoclutch_allowed = ?, tyre_blankets_allowed = ?, force_virtual_mirror = ?, fuel_rate = ?, damage_multiplier = ?, tyre_wear_rate = ?, allowed_tyres_out = ?, max_ballast_kg = ?, start_rule = ?, race_gas_penality_disabled = ?, dynamic_track = ?, dynamic_track_preset = ?, session_start = ?, randomness = ?, session_transfer = ?, lap_gain = ?, kick_quorum = ?, voting_quorum = ?, vote_duration = ?, blacklist_mode = ?, max_contacts_per_km = ?, filled = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(&dif.Name, &dif.AbsAllowed, &dif.TcAllowed, &dif.StabilityAllowed, &dif.AutoclutchAllowed, &dif.TyreBlanketsAllowed, &dif.ForceVirtualMirror, &dif.FuelRate, &dif.DamageMultiplier, &dif.TyreWearRate, &dif.AllowedTyresOut, &dif.MaxBallastKg, &dif.StartRule, &dif.RaceGasPenalityDisabled, &dif.DynamicTrack, &dif.DynamicTrackPreset, &dif.SessionStart, &dif.Randomness, &dif.SessionTransfer, &dif.LapGain, &dif.KickQuorum, &dif.VoteDuration, &dif.VoteDuration, &dif.BlacklistMode, &dif.MaxContactsPerKm, &dif.Id)
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

func (dba Dbaccess) selectSession(id int) UserSession {
	ses := UserSession{}
	stmt, err := dba.db.Prepare("SELECT id, name, booking_enabled, booking_time, practice_enabled, practice_time, practice_is_open, qualify_enabled, qualify_time, qualify_is_open, qualify_max_wait_perc, race_enabled, race_time, race_extra_lap, race_over_time, race_wait_time, race_is_open, reversed_grid_positions, race_pit_window_start, race_pit_window_end FROM user_session WHERE id = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	err = stmt.QueryRow(id).Scan(&ses.Id, &ses.Name, &ses.BookingEnabled, &ses.BookingTime, &ses.PracticeEnabled, &ses.PracticeTime, &ses.PracticeIsOpen, &ses.QualifyEnabled, &ses.QualifyTime, &ses.QualifyIsOpen, &ses.QualifyMaxWaitPerc, &ses.RaceEnabled, &ses.RaceTime, &ses.RaceExtraLap, &ses.RaceOverTime, &ses.RaceWaitTime, &ses.RaceIsOpen, &ses.ReversedGridPositions, &ses.RacePitWindowStart, &ses.RacePitWindowEnd)
	defer stmt.Close()
	if err != nil {
		log.Print(err)
	}

	return ses
}

func (dba Dbaccess) selectSessionList(filled bool) []DropDownList {
	return dba.selectDropDownList(filled, "user_session")
}

func (dba Dbaccess) insertSession(difficultyname string) int64 {
	return dba.insertNameInto(difficultyname, "user_session")
}

func (dba Dbaccess) deleteSession(id int) (int64, error) {
	return dba.deleteFrom(id, "user_session")
}

func (dba Dbaccess) updateSession(ses UserSession) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_session SET name = ?, booking_enabled = ?, booking_time = ?, practice_enabled = ?, practice_time = ?, practice_is_open = ?, qualify_enabled = ?, qualify_time = ?, qualify_is_open = ?, qualify_max_wait_perc = ?, race_enabled = ?, race_time = ?, race_extra_lap = ?, race_over_time = ?, race_wait_time = ?, race_is_open = ?, reversed_grid_positions = ?, race_pit_window_start = ?, race_pit_window_end = ?, filled = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(&ses.Name, &ses.BookingEnabled, &ses.BookingTime, &ses.PracticeEnabled, &ses.PracticeTime, &ses.PracticeIsOpen, &ses.QualifyEnabled, &ses.QualifyTime, &ses.QualifyIsOpen, &ses.QualifyMaxWaitPerc, &ses.RaceEnabled, &ses.RaceTime, &ses.RaceExtraLap, &ses.RaceOverTime, &ses.RaceWaitTime, &ses.RaceIsOpen, &ses.ReversedGridPositions, &ses.RacePitWindowStart, &ses.RacePitWindowEnd, &ses.Id)
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

func (dba Dbaccess) selectTimeWeather(id int) UserTime {
	time := UserTime{}
	stmt, err := dba.db.Prepare("SELECT id, name, time, time_of_day_multi, csp_enabled FROM user_time WHERE id = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&time.Id, &time.Name, &time.Time, &time.TimeOfDayMulti, &time.CspEnabled)
	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.db.Prepare("SELECT a.id, name, user_time_id, graphics, base_temperature_ambient, base_temperature_road, variation_ambient, variation_road, wind_base_speed_min, wind_base_speed_max, wind_base_direction, wind_variation_direction, csp_time, csp_time_of_day_multi, csp_date FROM user_time_weather a JOIN cache_weather b on a.graphics = b.key WHERE user_time_id = ?")
	if err != nil {
		log.Print(err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err)
	}

	time.Weathers = make([]UserTimeWeather, 0)
	wt := UserTimeWeather{}
	for rows.Next() {
		err = rows.Scan(&wt.Id, &wt.Name, &wt.UserTimeId, &wt.Graphics, &wt.BaseTemperatureAmbient, &wt.BaseTemperatureRoad, &wt.VariationAmbient, &wt.VariationRoad, &wt.WindBaseSpeedMin, &wt.WindBaseSpeedMax, &wt.WindBaseDirection, &wt.WindVariationDirection, &wt.CspTime, &wt.CspTimeOfDayMulti, &wt.CspDate)
		if err != nil {
			log.Fatal(err)
		}
		time.Weathers = append(time.Weathers, wt)
	}

	return time
}

func (dba Dbaccess) selectTimeList(filled bool) []DropDownList {
	return dba.selectDropDownList(filled, "user_time")
}

func (dba Dbaccess) insertTime(timename string) int64 {
	return dba.insertNameInto(timename, "user_time")
}

func (dba Dbaccess) deleteTime(id int) (int64, error) {
	rows, err := dba.deleteFrom(id, "user_time")

	if err != nil {
		return 0, err
	}

	stmt, err := dba.db.Prepare("DELETE FROM user_time_weather WHERE user_time_id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	return rows, err
}

func (dba Dbaccess) updateTime(time UserTime) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_time SET name = ?, time = ?, time_of_day_multi = ?, csp_enabled = ?, filled = 1 WHERE id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(&time.Name, &time.Time, &time.TimeOfDayMulti, &time.CspEnabled, &time.Id)
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
		stmt, err = dba.db.Prepare("INSERT INTO user_time_weather (user_time_id, graphics, base_temperature_ambient, base_temperature_road, variation_ambient, variation_road, wind_base_speed_min, wind_base_speed_max, wind_base_direction, wind_variation_direction, csp_time, csp_time_of_day_multi, csp_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&time.Id, &w.Graphics, &w.BaseTemperatureAmbient, &w.BaseTemperatureRoad, &w.VariationAmbient, &w.VariationRoad, &w.WindBaseSpeedMin, &w.WindBaseSpeedMax, &w.WindBaseDirection, &w.WindVariationDirection, &w.CspTime, &w.CspTimeOfDayMulti, &w.CspDate)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba Dbaccess) selectClassEntries(id int) UserClass {
	cls := UserClass{}
	stmt, err := dba.db.Prepare("SELECT id, name FROM user_class WHERE id = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&cls.Id, &cls.Name)
	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.db.Prepare("SELECT id, user_class_id, cache_car_key, skin_key, ballast FROM user_class_entry WHERE user_class_id = ?")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)

	cls.Entries = make([]UserClassEntry, 0)
	for rows.Next() {
		ent := UserClassEntry{}
		err = rows.Scan(&ent.Id, &ent.UserClassId, &ent.CacheCarKey, &ent.SkinKey, &ent.Ballast)
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

func (dba Dbaccess) selectClassList(filled bool) []DropDownList {
	return dba.selectDropDownList(filled, "user_class")
}

func (dba Dbaccess) insertClass(timename string) int64 {
	return dba.insertNameInto(timename, "user_class")
}

func (dba Dbaccess) deleteClass(id int) (int64, error) {
	_, err := dba.deleteFrom(id, "user_class")

	if err != nil {
		return 0, err
	}

	stmt, err := dba.db.Prepare("DELETE FROM user_class_entry WHERE user_class_id = ?")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		log.Print(err)
	}

	return 1, nil
}

func (dba Dbaccess) updateClass(cls UserClass) int64 {
	stmt, err := dba.db.Prepare("UPDATE user_class SET name = ?, filled = 1 WHERE id = ?")

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(cls.Name, cls.Id)
	if err != nil {
		log.Print(err)
	}

	stmt, err = dba.db.Prepare("DELETE FROM user_class_entry WHERE user_class_id = ?")
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
		_, err = stmt.Exec(cls.Id, ent.CacheCarKey, ent.SkinKey)
		defer stmt.Close()

		if err != nil {
			log.Print(err)
		}
	}

	return 1
}

func (dba Dbaccess) updateCacheCars(cars []CacheCar) int64 {
	for _, car := range cars {
		stmt, err := dba.db.Prepare("INSERT INTO cache_car (key, name, brand, desc, tags, class, specs, torque, power, skins) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
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

		// err is assumed to be FK error
		if err != nil {
			stmt, err = dba.db.Prepare("UPDATE cache_car SET name = ?, brand = ?, desc = ?, tags = ?, class = ?, specs = ?, torque = ?, power = ?, skins = ? WHERE key = ?")
			if err != nil {
				log.Print(err)
			}
			_, err = stmt.Exec(&car.Name, &car.Brand, &car.Desc, &tags, &car.Class, &specs, &torque, &power, &skins, &car.Key)
			defer stmt.Close()
			if err != nil {
				log.Print(err)
			}
		}
	}

	return 1
}

func (dba Dbaccess) selectCacheCars() []CacheCar {
	rows, err := dba.db.Query("SELECT key, name, brand, desc, tags, class, specs, torque, power, skins FROM cache_car ORDER BY name")
	if err != nil {
		log.Print(err)
	}

	var tags string
	var specs string
	var power string
	var torque string
	var skins string

	cars := make([]CacheCar, 0)
	for rows.Next() {
		car := CacheCar{}
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

func (dba Dbaccess) selectCacheCar(carkey string) CacheCar {
	var tags string
	var specs string
	var power string
	var torque string
	var skins string

	car := CacheCar{}
	stmt, err := dba.db.Prepare("SELECT key, name, brand, desc, tags, class, specs, torque, power, skins FROM cache_car WHERE key = ? ORDER BY name")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(carkey).Scan(&car.Key, &car.Name, &car.Brand, &car.Desc, &tags, &car.Class, &specs, &torque, &power, &skins)
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

func (dba Dbaccess) updateCacheTracks(tracks []CacheTrack) int64 {
	for _, track := range tracks {

		tagsRes, err := json.Marshal(&track.Tags)
		if err != nil {
			log.Print(err)
		}
		tags := string(tagsRes)

		stmt, err := dba.db.Prepare("INSERT INTO cache_track (key, config, name, desc, tags, country, city, length, width, pitboxes, run) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&track.Key, &track.Config, &track.Name, &track.Desc, tags, &track.Country, &track.City, &track.Length, &track.Width, &track.Pitboxes, &track.Run)
		defer stmt.Close()

		// err is assumed to be FK error
		if err != nil {
			stmt, err := dba.db.Prepare("UPDATE cache_track SET name = ?, desc = ?, tags = ?, country = ?, city = ?, length = ?, width = ?, pitboxes = ?, run = ? WHERE key = ? AND config = ?")
			if err != nil {
				log.Print(err)
			}
			_, err = stmt.Exec(&track.Name, &track.Desc, tags, &track.Country, &track.City, &track.Length, &track.Width, &track.Pitboxes, &track.Run, &track.Key, &track.Config)
			defer stmt.Close()

			if err != nil {
				log.Print(err)
			}
		}
	}

	return 1
}

func (dba Dbaccess) selectCacheTracks() []CacheTrack {
	rows, err := dba.db.Query("SELECT key, config, name, desc, tags, country, city, length, width, pitboxes, run FROM cache_track ORDER BY name")
	if err != nil {
		log.Print(err)
	}

	var tags string
	tracks := make([]CacheTrack, 0)
	for rows.Next() {
		t := CacheTrack{}
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

func (dba Dbaccess) selectCacheTrack(trackkey string, trackconfig string) CacheTrack {
	t := CacheTrack{}
	var tags string

	if trackconfig == "" {
		stmt, err := dba.db.Prepare("SELECT key, config, name, desc, tags, country, city, length, width, pitboxes, run FROM cache_track WHERE key = ? LIMIT 1")
		if err != nil {
			log.Print(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(trackkey).Scan(&t.Key, &t.Config, &t.Name, &t.Desc, &tags, &t.Country, &t.City, &t.Length, &t.Width, &t.Pitboxes, &t.Run)
		if err != nil {
			log.Print(err)
		}
	} else {
		stmt, err := dba.db.Prepare("SELECT key, config, name, desc, tags, country, city, length, width, pitboxes, run FROM cache_track WHERE key = ? AND config = ? LIMIT 1")
		if err != nil {
			log.Print(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(trackkey, trackconfig).Scan(&t.Key, &t.Config, &t.Name, &t.Desc, &tags, &t.Country, &t.City, &t.Length, &t.Width, &t.Pitboxes, &t.Run)
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

func (dba Dbaccess) updateCacheWeathers(weathers []CacheWeather) int64 {
	for _, w := range weathers {
		stmt, err := dba.db.Prepare("INSERT INTO cache_weather (key, name) VALUES (?, ?)")
		if err != nil {
			log.Print(err)
		}
		_, err = stmt.Exec(&w.Key, &w.Name)
		defer stmt.Close()

		// err is assumed to be FK error
		if err != nil {
			stmt, err := dba.db.Prepare("UPDATE cache_weather set name = ? WHERE key = ?")
			if err != nil {
				log.Print(err)
			}
			_, err = stmt.Exec(&w.Key, &w.Name)
			defer stmt.Close()

			if err != nil {
				log.Print(err)
			}
		}
	}

	return 1
}

func (dba Dbaccess) selectCacheWeathers() []CacheWeather {
	rows, err := dba.db.Query("SELECT key, name FROM cache_weather ORDER BY name")
	if err != nil {
		log.Print(err)
	}

	weathers := make([]CacheWeather, 0)
	for rows.Next() {
		w := CacheWeather{}
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

func (dba Dbaccess) selectCacheWeather(weatherkey string) CacheWeather {
	w := CacheWeather{}
	stmt, err := dba.db.Prepare("SELECT key, name FROM cache_weather WHERE key = ? LIMIT 1")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(weatherkey).Scan(&w.Key, &w.Name)
	if err != nil {
		log.Print(err)
	}

	return w
}
