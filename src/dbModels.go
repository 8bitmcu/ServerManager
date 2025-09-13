package main

type Users struct {
	Name            *string `form:"name"`
	Password        *string `form:"password"`
	MeasurementUnit *int    `form:"measurement_unit"`
	TempUnit        *int    `form:"temp_unit"`
}

type UserConfig struct {
	Name               *string `form:"name"`
	Password           *string `form:"password"`
	AdminPassword      *string `form:"admin_password"`
	RegisterToLobby    *int    `form:"register_to_lobby"`
	LockedEntryList    *int    `form:"locked_entry_list"`
	ResultScreenTime   *int    `form:"result_screen_time"`
	UdpPort            *int    `form:"udp_port"`
	TcpPort            *int    `form:"tcp_port"`
	HttpPort           *int    `form:"http_port"`
	ClientSendInterval *int    `form:"client_send_interval"`
	NumThreads         *int    `form:"num_threads"`
	MaxClients         *int    `form:"max_clients"`
	WelcomeMessage     *string `form:"welcome_message"`
	AppendEventname    *int    `form:"append_eventname"`
	AppendModlinks     *int    `form:"append_modlinks"`
	InstallPath        *string `form:"install_path"`
	CspRequired        *int    `form:"csp_required"`
	CspVersion         *int    `form:"csp_version"`
	CspPhycars         *int    `form:"csp_phycars"`
	CspPhytracks       *int    `form:"csp_phytracks"`
	CspHidepit         *int    `form:"csp_hidepit"`
	CfgFilled          *int
	ModFilled          *int
	SecretKey          *string
}

type ServerEvent struct {
	Id        *int
	UserEvent UserEvent
	ServerCfg *string
	EntryList *string
	StartedAt *int64
	Finished  *int
}

type UserEvent struct {
	Id                *int `form:"id"`
	EventCategoryId   *int `form:"event_category_id"`
	CategoryName      *string
	RaceLaps          *int `form:"race_laps"`
	Strategy          *int `form:"strategy"`
	TrackName         *string
	TrackLength       *int
	Pitboxes          *int
	DifficultyName    *string
	AbsAllowed        *int
	TcAllowed         *int
	StabilityAllowed  *int
	AutoclutchAllowed *int
	SessionName       *string
	BookingEnabled    *int
	BookingTime       *int
	PracticeEnabled   *int
	PracticeTime      *int
	QualifyEnabled    *int
	QualifyTime       *int
	RaceEnabled       *int
	RaceTime          *int
	ClassName         *string
	Entries           *int
	TimeName          *string
	Time              *string
	Graphics          *string
	TruncWeather      *int
	CspWeather        *int
	CacheTrackKey     *string `form:"track_key"`
	CacheTrackConfig  *string `form:"track_config"`
	DifficultyId      *int    `form:"difficulty"`
	SessionId         *int    `form:"session"`
	ClassId           *int    `form:"class"`
	TimeId            *int    `form:"time"`
}

type UserEventList struct {
	Id              *int    `json:"id"`
	EventCategoryId *int    `json:"event_category_id"`
	TrackName       *string `json:"track_name"`
}

type UserEventCategory struct {
	Id     *int
	Name   *string `form:"name"`
	Events []UserEvent
}

type UserDifficulty struct {
	Id                      *int    `from:"id" json:"id"`
	Name                    *string `form:"name" json:"name"`
	AbsAllowed              *int    `form:"abs_allowed" json:"abs_allowed"`
	TcAllowed               *int    `form:"tc_allowed" json:"tc_allowed"`
	StabilityAllowed        *int    `form:"stability_allowed" json:"stability_allowed"`
	AutoclutchAllowed       *int    `form:"autoclutch_allowed" json:"autoclutch_allowed"`
	TyreBlanketsAllowed     *int    `form:"tyre_blankets_allowed" json:"tyre_blankets_allowed"`
	ForceVirtualMirror      *int    `form:"force_virtual_mirror" json:"force_virtual_mirror"`
	FuelRate                *int    `form:"fuel_rate" json:"fuel_rate"`
	DamageMultiplier        *int    `form:"damage_multiplier" json:"damage_multiplier"`
	TyreWearRate            *int    `form:"tyre_wear_rate" json:"tyre_wear_rate"`
	AllowedTyresOut         *int    `form:"allowed_tyres_out" json:"allowed_tyres_out"`
	MaxBallastKg            *int    `form:"max_ballast_kg" json:"max_ballast_kg"`
	StartRule               *int    `form:"start_rule" json:"start_rule"`
	RaceGasPenalityDisabled *int    `form:"race_gas_penality_disabled" json:"race_gas_penality_disabled"`
	DynamicTrack            *int    `form:"dynamic_track" json:"dynamic_track"`
	DynamicTrackPreset      *int    `form:"dynamic_track_preset" json:"dynamic_track_preset"`
	SessionStart            *int    `form:"session_start" json:"session_start"`
	Randomness              *int    `form:"randomness" json:"randomness"`
	SessionTransfer         *int    `form:"session_transfer" json:"session_transfer"`
	LapGain                 *int    `form:"lap_gain" json:"lap_gain"`
	KickQuorum              *int    `form:"kick_quorum" json:"kick_quorum"`
	VoteDuration            *int    `form:"vote_duration" json:"vote_duration"`
	VotingQuorum            *int    `form:"voting_quorum" json:"voting_quorum"`
	BlacklistMode           *int    `form:"blacklist_mode" json:"blacklist_mode"`
	MaxContactsPerKm        *int    `form:"max_contacts_per_km" json:"max_contacts_per_km"`
}

type UserSession struct {
	Id                    *int
	Name                  *string `form:"name" json:"name"`
	BookingEnabled        *int    `form:"booking_enabled" json:"booking_enabled"`
	BookingTime           *int    `form:"booking_time" json:"booking_time"`
	PracticeEnabled       *int    `form:"practice_enabled" json:"practice_enabled"`
	PracticeTime          *int    `form:"practice_time" json:"practice_time"`
	PracticeIsOpen        *int    `form:"practice_is_open" json:"practice_is_open"`
	QualifyEnabled        *int    `form:"qualify_enabled" json:"qualify_enabled"`
	QualifyTime           *int    `form:"qualify_time" json:"qualify_time"`
	QualifyIsOpen         *int    `form:"qualify_is_open" json:"qualify_is_open"`
	QualifyMaxWaitPerc    *int    `form:"qualify_max_wait_perc" json:"qualify_max_wait_perc"`
	RaceEnabled           *int    `form:"race_enabled" json:"race_enabled"`
	RaceTime              *int    `form:"race_time" json:"race_time"`
	RaceExtraLap          *int    `form:"race_extra_lap" json:"race_extra_lap"`
	RaceOverTime          *int    `form:"race_over_time" json:"race_over_time"`
	RaceWaitTime          *int    `form:"race_wait_time" json:"race_wait_time"`
	RaceIsOpen            *int    `form:"race_is_open" json:"race_is_open"`
	ReversedGridPositions *int    `form:"reversed_grid_positions" json:"reversed_grid_positions"`
	RacePitWindowStart    *int    `form:"race_pit_window_start" json:"race_pit_window_start"`
	RacePitWindowEnd      *int    `form:"race_pit_window_end" json:"race_pit_window_end"`
}

type UserTime struct {
	Id             *int
	Name           *string           `form:"name" json:"name"`
	Time           *string           `form:"time" json:"time"`
	TimeOfDayMulti *int              `form:"time_of_day_multi" json:"time_of_day_multi"`
	CspEnabled     *int              `form:"csp_enabled" json:"csp_enabled"`
	Weathers       []UserTimeWeather `json:"weathers"`
}

type UserTimeWeather struct {
	Id                     *int
	Name                   *string `json:"name"`
	UserTimeId             *int
	Graphics               *string `json:"graphics"`
	BaseTemperatureAmbient *int    `json:"base_temperature_ambient,string"`
	BaseTemperatureRoad    *int    `json:"base_temperature_road,string"`
	VariationAmbient       *int    `json:"variation_ambient,string"`
	VariationRoad          *int    `json:"variation_road,string"`
	WindBaseSpeedMin       *int    `json:"wind_base_speed_min,string"`
	WindBaseSpeedMax       *int    `json:"wind_base_speed_max,string"`
	WindBaseDirection      *int    `json:"wind_base_direction,string"`
	WindVariationDirection *int    `json:"wind_variation_direction,string"`
	CspTime                *string `json:"csp_time"`
	CspTimeOfDayMulti      *int    `json:"csp_time_of_day_multi,string"`
	CspDate                *string `json:"csp_date"`
}

type UserClass struct {
	Id      *int
	Name    *string          `form:"name" json:"name"`
	Entries []UserClassEntry `json:"entries"`
}

type UserClassEntry struct {
	Id          *int
	UserClassId *int    `json:"user_class_id"`
	CacheCarKey *string `json:"cache_car_key"`
	SkinKey     *string `json:"skin_key"`
	Ballast     *int    `json:"ballast"`
}

type DropDownList struct {
	Id   *int    `json:"id"`
	Name *string `json:"name"`
}

type CacheCar struct {
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

type CacheTrack struct {
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

type CacheWeather struct {
	Id   *int
	Key  *string `json:"key"`
	Name *string `json:"name"`
}
