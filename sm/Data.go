package sm

type Users struct {
	Name             *string `form:"name"`
	Password         *string `form:"password"`
	Measurement_Unit *int    `form:"measurement_unit"`
	Temp_Unit        *int    `form:"temp_unit"`
}

type User_Config struct {
	Name                 *string `form:"name"`
	Password             *string `form:"password"`
	Admin_Password       *string `form:"admin_password"`
	Register_To_Lobby    *int    `form:"register_to_lobby"`
	Locked_Entry_List    *int    `form:"locked_entry_list"`
	Result_Screen_Time   *int    `form:"result_screen_time"`
	Udp_Port             *int    `form:"udp_port"`
	Tcp_Port             *int    `form:"tcp_port"`
	Http_Port            *int    `form:"http_port"`
	Client_Send_Interval *int    `form:"client_send_interval"`
	Num_Threads          *int    `form:"num_threads"`
	Max_Clients          *int    `form:"max_clients"`
	Welcome_Message      *string `form:"welcome_message"`
	Install_Path         *string `form:"install_path"`
	Csp_Required         *int    `form:"csp_required"`
	Csp_Version          *int    `form:"csp_version"`
	Csp_Phycars          *int    `form:"csp_phycars"`
	Csp_Phytracks        *int    `form:"csp_phytracks"`
	Csp_Hidepit          *int    `form:"csp_hidepit"`
	Cfg_Filled           *int
	Mod_Filled           *int
}

type User_Event struct {
	Id                 *int `form:"id"`
	Race_Laps          *int `form:"race_laps"`
	Strategy           *int `form:"strategy"`
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
	TruncWeather       *int
	Csp_Weather        *int
	Cache_Track_Key    *string `form:"track_key"`
	Cache_Track_Config *string `form:"track_config"`
	Difficulty_Id      *int    `form:"difficulty"`
	Session_Id         *int    `form:"session"`
	Class_Id           *int    `form:"class"`
	Time_Id            *int    `form:"time"`
	ServerCfg          *string
	EntryList          *string
	Started_At         *int64
	Finished           *int
}

type User_Difficulty struct {
	Id                         *int    `from:"id" json:"id"`
	Name                       *string `form:"name" json:"name"`
	Abs_Allowed                *int    `form:"abs_allowed" json:"abs_allowed"`
	Tc_Allowed                 *int    `form:"tc_allowed" json:"tc_allowed"`
	Stability_Allowed          *int    `form:"stability_allowed" json:"stability_allowed"`
	Autoclutch_Allowed         *int    `form:"autoclutch_allowed" json:"autoclutch_allowed"`
	Tyre_Blankets_Allowed      *int    `form:"tyre_blankets_allowed" json:"tyre_blankets_allowed"`
	Force_Virtual_Mirror       *int    `form:"force_virtual_mirror" json:"force_virtual_mirror"`
	Fuel_Rate                  *int    `form:"fuel_rate" json:"fuel_rate"`
	Damage_Multiplier          *int    `form:"damage_multiplier" json:"damage_multiplier"`
	Tyre_Wear_Rate             *int    `form:"tyre_wear_rate" json:"tyre_wear_rate"`
	Allowed_Tyres_Out          *int    `form:"allowed_tyres_out" json:"allowed_tyres_out"`
	Max_Ballast_Kg             *int    `form:"max_ballast_kg" json:"max_ballast_kg"`
	Start_Rule                 *int    `form:"start_rule" json:"start_rule"`
	Race_Gas_Penality_Disabled *int    `form:"race_gas_penality_disabled" json:"race_gas_penality_disabled"`
	Dynamic_Track              *int    `form:"dynamic_track" json:"dynamic_track"`
	Dynamic_Track_Preset       *int    `form:"dynamic_track_preset" json:"dynamic_track_preset"`
	Session_Start              *int    `form:"session_start" json:"session_start"`
	Randomness                 *int    `form:"randomness" json:"randomness"`
	Session_Transfer           *int    `form:"session_transfer" json:"session_transfer"`
	Lap_Gain                   *int    `form:"lap_gain" json:"lap_gain"`
	Kick_Quorum                *int    `form:"kick_quorum" json:"kick_quorum"`
	Vote_Duration              *int    `form:"vote_duration" json:"vote_duration"`
	Voting_Quorum              *int    `form:"voting_quorum" json:"voting_quorum"`
	Blacklist_Mode             *int    `form:"blacklist_mode" json:"blacklist_mode"`
	Max_Contacts_Per_Km        *int    `form:"max_contacts_per_km" json:"max_contacts_per_km"`
}

type User_Session struct {
	Id                      *int
	Name                    *string `form:"name" json:"name"`
	Booking_Enabled         *int    `form:"booking_enabled" json:"booking_enabled"`
	Booking_Time            *int    `form:"booking_time" json:"booking_time"`
	Practice_Enabled        *int    `form:"practice_enabled" json:"practice_enabled"`
	Practice_Time           *int    `form:"practice_time" json:"practice_time"`
	Practice_Is_Open        *int    `form:"practice_is_open" json:"practice_is_open"`
	Qualify_Enabled         *int    `form:"qualify_enabled" json:"qualify_enabled"`
	Qualify_Time            *int    `form:"qualify_time" json:"qualify_time"`
	Qualify_Is_Open         *int    `form:"qualify_is_open" json:"qualify_is_open"`
	Qualify_Max_Wait_Perc   *int    `form:"qualify_max_wait_perc" json:"qualify_max_wait_perc"`
	Race_Enabled            *int    `form:"race_enabled" json:"race_enabled"`
	Race_Time               *int    `form:"race_time" json:"race_time"`
	Race_Extra_Lap          *int    `form:"race_extra_lap" json:"race_extra_lap"`
	Race_Over_Time          *int    `form:"race_over_time" json:"race_over_time"`
	Race_Wait_Time          *int    `form:"race_wait_time" json:"race_wait_time"`
	Race_Is_Open            *int    `form:"race_is_open" json:"race_is_open"`
	Reversed_Grid_Positions *int    `form:"reversed_grid_positions" json:"reversed_grid_positions"`
	Race_Pit_Window_Start   *int    `form:"race_pit_window_start" json:"race_pit_window_start"`
	Race_Pit_Window_End     *int    `form:"race_pit_window_end" json:"race_pit_window_end"`
}

type User_Time struct {
	Id                *int
	Name              *string             `form:"name" json:"name"`
	Time              *string             `form:"time" json:"time"`
	Time_Of_Day_Multi *int                `form:"time_of_day_multi" json:"time_of_day_multi"`
	Csp_Enabled       *int                `form:"csp_enabled" json:"csp_enabled"`
	Weathers          []User_Time_Weather `json:"weathers"`
}

type User_Time_Weather struct {
	Id                       *int
	Name                     *string `json:"name"`
	User_Time_Id             *int
	Graphics                 *string `json:"graphics"`
	Base_Temperature_Ambient *int    `json:"base_temperature_ambient,string"`
	Base_Temperature_Road    *int    `json:"base_temperature_road,string"`
	Variation_Ambient        *int    `json:"variation_ambient,string"`
	Variation_Road           *int    `json:"variation_road,string"`
	Wind_Base_Speed_Min      *int    `json:"wind_base_speed_min,string"`
	Wind_Base_Speed_Max      *int    `json:"wind_base_speed_max,string"`
	Wind_Base_Direction      *int    `json:"wind_base_direction,string"`
	Wind_Variation_Direction *int    `json:"wind_variation_direction,string"`
	Csp_Time                 *string `json:"csp_time"`
	Csp_Time_Of_Day_Multi    *int    `json:"csp_time_of_day_multi,string"`
	Csp_Date                 *string `json:"csp_date"`
}

type User_Class struct {
	Id      *int
	Name    *string            `form:"name" json:"name"`
	Entries []User_Class_Entry `json:"entries"`
}

type User_Class_Entry struct {
	Id            *int
	User_Class_Id *int    `json:"user_class_id"`
	Cache_Car_Key *string `json:"cache_car_key"`
	Skin_Key      *string `json:"skin_key"`
	Ballast       *int    `json:"ballast"`
}

type DropDown_List struct {
	Id   *int    `json:"id"`
	Name *string `json:"name"`
}

type Cache_Car struct {
	Id    *int
	Key   *string   `json:"key"`
	Name  *string   `json:"name"`
	Brand *string   `json:"brand"`
	Desc  *string   `json:"description"`
	Tags  *[]string `json:"tags"`
	Class *string   `json:"class"`
	Specs struct {
		Bhp          string `json:"bhp"`
		Torque       string `json:"torque"`
		Weight       string `json:"weight"`
		Topspeed     string `json:"topspeed"`
		Acceleration string `json:"acceleration"`
		Pwratio      string `json:"pwratio"`
		Range        int    `json:"range"`
	} `json:"specs"`
	Torque [][]any `json:"torqueCurve"`
	Power  [][]any `json:"powerCurve"`
	Skins  []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"skins"`
}

type Cache_Track struct {
	Id       *int
	Key      *string   `json:"key"`
	Config   *string   `json:"config"`
	Name     *string   `json:"name"`
	Desc     *string   `json:"desc"`
	Tags     *[]string `json:"tags"`
	Country  *string   `json:"country"`
	City     *string   `json:"city"`
	Length   *int      `json:"length"`
	Width    *string   `json:"width"`
	Pitboxes *int      `json:"pitboxes,string"`
	Run      *string   `json:"run"`
}

type Cache_Weather struct {
	Id   *int
	Key  *string `json:"key"`
	Name *string `json:"name"`
}

var Dba Dbaccess
var Cr ConfigRenderer
var Status Server_Status
var Udp UDPPlugin
var SecretKey = []byte("XBLn0dUoXPVk742lkRVILa82hbRXz6Tx")
