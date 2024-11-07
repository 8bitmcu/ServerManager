package sm

import (
	"encoding/binary"
	"log"
	"math"
	"net"
	"strconv"

	"golang.org/x/text/encoding/unicode/utf32"
)

const (
	ACSP_PROTOCOL_VERSION      = 4
	ACSP_NEW_SESSION           = 50
	ACSP_NEW_CONNECTION        = 51
	ACSP_CONNECTION_CLOSED     = 52
	ACSP_CAR_UPDATE            = 53
	ACSP_CAR_INFO              = 54
	ACSP_END_SESSION           = 55
	ACSP_LAP_COMPLETED         = 73
	ACSP_VERSION               = 56
	ACSP_CHAT                  = 57
	ACSP_CLIENT_LOADED         = 58
	ACSP_SESSION_INFO          = 59
	ACSP_ERROR                 = 60
	ACSP_CLIENT_EVENT          = 130
	ACSP_CE_COLLISION_WITH_CAR = 10
	ACSP_CE_COLLISION_WITH_ENV = 11
	ACSP_REALTIMEPOS_INTERVAL  = 200
	ACSP_GET_CAR_INFO          = 201
	ACSP_SEND_CHAT             = 202
	ACSP_BROADCAST_CHAT        = 203
	ACSP_GET_SESSION_INFO      = 204
	ACSP_SET_SESSION_INFO      = 205
	ACSP_KICK_USER             = 206
	ACSP_NEXT_SESSION          = 207
	ACSP_RESTART_SESSION       = 208
	ACSP_ADMIN_COMMAND         = 209
)

type SessionInfo struct {
	version               int
	session_index         int
	current_session_index int
	session_count         int
	server_name           string
	track                 string
	track_config          string
	name                  string
	typ                   int
	time                  int
	laps                  int
	wait_time             int
	ambient_temp          int
	road_temp             int
	weather_graphics      string
	elapsed_ms            int32
}

type ClientEvent struct {
	car_id       int
	other_car_id int
	event_type   int
	impact_speed float32
	world_pos    Vector
	rel_pos      Vector
}

type CarInfo struct {
	car_id       int
	is_connected bool
	car_model    string
	car_skin     string
	driver_name  string
	driver_team  string
	driver_guid  string
}

type CarUpdate struct {
	car_id                int
	position              Vector
	velocity              Vector
	gear                  int
	engine_rpm            int
	normalized_spline_pos float32
}

type NewConnection struct {
	driver_name string
	driver_guid string
	car_id      int
	car_model   string
	car_skin    string
}

type ConnectionClosed struct {
	driver_name string
	driver_guid string
	car_id      int
	car_model   string
	car_skin    string
}

type LapCompleted struct {
	car_id  int
	laptime uint32
	cuts    int
}

type Vector struct {
	x float32
	y float32
	z float32
}

type UdpReader struct {
	data      []byte
	readIndex int64
}

func (r *UdpReader) New(p []byte) {
	r.data = p
}

func (r *UdpReader) Read_Byte() int {
	val := r.data[r.readIndex]
	r.readIndex = r.readIndex + 1
	return int(val)
}

func (r *UdpReader) Read_Bytes(bytes int) []byte {
	val := r.data[r.readIndex : r.readIndex+int64(bytes)]
	r.readIndex = r.readIndex + int64(bytes)
	return val
}

func (r *UdpReader) Read_String() string {
	length := int(r.Read_Byte())
	bytes := r.Read_Bytes(length)
	return string(bytes)
}

func (r *UdpReader) Read_UTF32_String() string {
	length := int(r.Read_Byte() * 4)
	bytes := r.Read_Bytes(length)

	val, err := utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewDecoder().Bytes(bytes)

	if err != nil {
		log.Print(err)
	}

	return string(val)
}

func (r *UdpReader) Read_Uint16() int {
	bytes := r.Read_Bytes(2)
	data := binary.LittleEndian.Uint16(bytes)
	return int(data)
}

func (r *UdpReader) Read_Int32() int32 {
	bytes := r.Read_Bytes(4)
	data := binary.LittleEndian.Uint32(bytes)
	return int32(data)
}

func (r *UdpReader) Read_Uint32() uint32 {
	bytes := r.Read_Bytes(4)
	data := binary.LittleEndian.Uint32(bytes)
	return data
}

func (r *UdpReader) Read_Float() float32 {
	bytes := r.Read_Bytes(4)
	a := binary.LittleEndian.Uint32(bytes)
	a2 := math.Float32frombits(a)
	return a2
}

type UdpWriter struct {
	data []byte
}

func (w *UdpWriter) Write_Byte(d byte) {
	w.data = append(w.data, d)
}

func (w *UdpWriter) Write_UTF32_String(str string) {
	val, err := utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewEncoder().Bytes([]byte(str))

	w.Write_Byte(byte(len(str)))

	if err != nil {
		log.Print(err)
	}
	w.data = append(w.data, val[:]...)
}

func Read_SessionInfo(r UdpReader) SessionInfo {
	var s SessionInfo
	s.version = r.Read_Byte()
	s.session_index = r.Read_Byte()
	s.current_session_index = r.Read_Byte()
	s.session_count = r.Read_Byte()
	s.server_name = r.Read_UTF32_String()
	s.track = r.Read_String()
	s.track_config = r.Read_String()
	s.name = r.Read_String()
	s.typ = r.Read_Byte()
	s.time = r.Read_Uint16()
	s.laps = r.Read_Uint16()
	s.wait_time = r.Read_Uint16()
	s.ambient_temp = r.Read_Byte()
	s.road_temp = r.Read_Byte()
	s.weather_graphics = r.Read_String()
	s.elapsed_ms = r.Read_Int32()
	return s
}

type UDPPlugin struct {
	conn *net.UDPConn
}

func UdpListen() UDPPlugin {
	var udp UDPPlugin
	udpClient, err := net.ResolveUDPAddr("udp", ":5001")
	udpServer, err := net.ResolveUDPAddr("udp", ":5000")

	if err != nil {
		log.Print(err)
	}

	udp.conn, err = net.DialUDP("udp", udpClient, udpServer)
	if err != nil {
		log.Print(err)
	}

	return udp
}

func (udp UDPPlugin) Receive() {
	data := make([]byte, 1024)
	udp.conn.Read(data)

	r := UdpReader{}
	r.New(data)

	acsp := r.Read_Byte()

	switch acsp {
	case ACSP_ERROR:
		err := r.Read_UTF32_String()
		log.Print("ACSP_ERROR: " + err)

	case ACSP_CHAT:
		car := r.Read_Byte()
		msg := r.Read_UTF32_String()
		log.Print("ACSP_CHAT: " + strconv.Itoa(car) + "; " + msg)

	case ACSP_CLIENT_LOADED:
		car := r.Read_Byte()
		log.Print("ACSP_CLIENT_LOADED: " + strconv.Itoa(car))

	case ACSP_VERSION:
		v := r.Read_Byte()
		log.Print("ACSP_VERSION: " + strconv.Itoa(v))

	case ACSP_NEW_SESSION:
		sess := Read_SessionInfo(r)
		log.Print("ACSP_NEW_SESSION: ", sess)

	case ACSP_SESSION_INFO:
		sess := Read_SessionInfo(r)
		log.Print("ACSP_SESSION_INFO: ", sess)

	case ACSP_END_SESSION:
		file := r.Read_UTF32_String()
		log.Print("ACSP_END_SESSION: " + file)

	case ACSP_CLIENT_EVENT:
		var ce ClientEvent
		ce.event_type = r.Read_Byte()
		ce.car_id = r.Read_Byte()
		if ce.event_type == ACSP_CE_COLLISION_WITH_CAR {
			ce.other_car_id = r.Read_Byte()
		}
		ce.impact_speed = r.Read_Float()
		ce.world_pos = Vector{r.Read_Float(), r.Read_Float(), r.Read_Float()}
		ce.rel_pos = Vector{r.Read_Float(), r.Read_Float(), r.Read_Float()}
		log.Print("ACSP_CLIENT_EVENT: ", ce)

	case ACSP_CAR_INFO:
		var ci CarInfo
		ci.car_id = r.Read_Byte()
		ci.is_connected = r.Read_Byte() != 0
		ci.car_model = r.Read_UTF32_String()
		ci.car_skin = r.Read_UTF32_String()
		ci.driver_name = r.Read_UTF32_String()
		ci.driver_team = r.Read_UTF32_String()
		ci.driver_guid = r.Read_UTF32_String()
		log.Print("ACSP_CAR_INFO: ", ci)

	case ACSP_CAR_UPDATE:
		var cu CarUpdate
		cu.car_id = r.Read_Byte()
		cu.position = Vector{r.Read_Float(), r.Read_Float(), r.Read_Float()}
		cu.velocity = Vector{r.Read_Float(), r.Read_Float(), r.Read_Float()}
		cu.gear = r.Read_Byte()
		cu.engine_rpm = r.Read_Uint16()
		cu.normalized_spline_pos = r.Read_Float()
		log.Print("ACSP_CAR_UPDATE: ", cu)

	case ACSP_NEW_CONNECTION:
		var nc NewConnection
		nc.driver_name = r.Read_UTF32_String()
		nc.driver_guid = r.Read_UTF32_String()
		nc.car_id = r.Read_Byte()
		nc.car_model = r.Read_String()
		nc.car_skin = r.Read_String()
		log.Print("ACSP_NEW_CONNECTION: ", nc)

	case ACSP_CONNECTION_CLOSED:
		var cc ConnectionClosed
		cc.driver_name = r.Read_UTF32_String()
		cc.driver_guid = r.Read_UTF32_String()
		cc.car_id = r.Read_Byte()
		cc.car_model = r.Read_String()
		cc.car_skin = r.Read_String()
		log.Print("ACSP_CONNECTION_CLOSED: ", cc)

	case ACSP_LAP_COMPLETED:
		var lc LapCompleted
		lc.car_id = r.Read_Byte()
		lc.laptime = r.Read_Uint32()
		lc.cuts = r.Read_Byte()
		log.Print("ACSP_LAP_COMPLETED: ", lc)

	default:
		log.Print("ACSP Unknown code: "+strconv.Itoa(acsp), data)
	}
}

func (udp UDPPlugin) Write_AdminCommand(command string) {
	var w UdpWriter
	w.Write_Byte(ACSP_ADMIN_COMMAND)
	w.Write_UTF32_String(command)
	udp.conn.Write(w.data)
}

func (udp UDPPlugin) Write_BroadcastChat(message string) {
	var w UdpWriter
	w.Write_Byte(ACSP_BROADCAST_CHAT)
	w.Write_UTF32_String(message)
	udp.conn.Write(w.data)
}

func (udp UDPPlugin) Write_GetCarInfo(car_id int) {
	var w UdpWriter
	w.Write_Byte(ACSP_GET_CAR_INFO)
	w.Write_Byte(byte(car_id))
	udp.conn.Write(w.data)
}

func (udp UDPPlugin) Write_KickUser(car_id int) {
	var w UdpWriter
	w.Write_Byte(ACSP_KICK_USER)
	w.Write_Byte(byte(car_id))
	udp.conn.Write(w.data)
}

func (udp UDPPlugin) Write_NextSession() {
	var w UdpWriter
	w.Write_Byte(ACSP_NEXT_SESSION)
	udp.conn.Write(w.data)
}

func (udp UDPPlugin) Write_RestartSession() {
	var w UdpWriter
	w.Write_Byte(ACSP_RESTART_SESSION)
	udp.conn.Write(w.data)
}

func (udp UDPPlugin) Write_SendChat(car_id int, message string) {
	var w UdpWriter
	w.Write_Byte(ACSP_SEND_CHAT)
	w.Write_Byte(byte(car_id))
	w.Write_UTF32_String(message)
	udp.conn.Write(w.data)
}

