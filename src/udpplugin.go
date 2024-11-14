package main

import (
	"encoding/binary"
	"log"
	"math"
	"net"
	"strconv"
	"time"

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
	Version               int
	Session_index         int
	Current_session_index int
	Session_count         int
	Server_name           string
	Track                 string
	Track_config          string
	Name                  string
	Typ                   int
	Time                  int
	Laps                  int
	Wait_time             int
	Ambient_temp          int
	Road_temp             int
	Weather_graphics      string
	Elapsed_ms            int32
}

type ClientEvent struct {
	Car_id       int
	Other_car_id int
	Event_type   int
	Impact_speed float32
	World_pos    Vector
	Rel_pos      Vector
}

type CarInfo struct {
	Car_id       int
	Is_connected bool
	Car_model    string
	Car_skin     string
	Driver_name  string
	Driver_team  string
	Driver_guid  string
}

type CarUpdate struct {
	Car_id                int
	Position              Vector
	Velocity              Vector
	Gear                  int
	Engine_rpm            int
	Normalized_spline_pos float32
}

type NewConnection struct {
	Driver_name string
	Driver_guid string
	Car_id      int
	Car_model   string
	Car_skin    string
}

type ConnectionClosed struct {
	Driver_name string
	Driver_guid string
	Car_id      int
	Car_model   string
	Car_skin    string
}

type LapCompleted struct {
	Car_id  int
	Laptime uint32
	Cuts    int
}

type Vector struct {
	X float32
	Y float32
	Z float32
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
	s.Version = r.Read_Byte()
	s.Session_index = r.Read_Byte()
	s.Current_session_index = r.Read_Byte()
	s.Session_count = r.Read_Byte()
	s.Server_name = r.Read_UTF32_String()
	s.Track = r.Read_String()
	s.Track_config = r.Read_String()
	s.Name = r.Read_String()
	s.Typ = r.Read_Byte()
	s.Time = r.Read_Uint16()
	s.Laps = r.Read_Uint16()
	s.Wait_time = r.Read_Uint16()
	s.Ambient_temp = r.Read_Byte()
	s.Road_temp = r.Read_Byte()
	s.Weather_graphics = r.Read_String()
	s.Elapsed_ms = r.Read_Int32()
	return s
}

type UDPPlugin struct {
	conn   *net.UDPConn
	online bool
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
		log.Print("ACSP_ERROR: ", err)

	case ACSP_CHAT:
		car := r.Read_Byte()
		msg := r.Read_UTF32_String()
		log.Print("ACSP_CHAT: " + strconv.Itoa(car) + "; " + msg)

	case ACSP_CLIENT_LOADED:
		car := r.Read_Byte()
		log.Print("ACSP_CLIENT_LOADED: ", car)

	case ACSP_VERSION:
		v := r.Read_Byte()
		log.Print("ACSP_VERSION: ", v)
		Udp.online = true

	case ACSP_NEW_SESSION:
		sess := Read_SessionInfo(r)
		log.Print("ACSP_NEW_SESSION: ")
		Print_Interface(sess)
		Status.Session = sess

	case ACSP_SESSION_INFO:
		sess := Read_SessionInfo(r)
		log.Print("ACSP_SESSION_INFO: ")
		Print_Interface(sess)
		Status.Session = sess

	case ACSP_END_SESSION:
		file := r.Read_UTF32_String()
		log.Print("ACSP_END_SESSION: " + file)

		if Status.Session.Current_session_index == Status.Session.Session_count-1 {
			// This was the last session, let's kill the server and change track
			log.Print("Kicking players")
			Udp.Write_KickUser(0)
			time.Sleep(2 * time.Second)

			// TODO: kick all players


			log.Print("Track change")
			Stop()
			Status.Server_ApplyTrack()
			Start()

		}

	case ACSP_CLIENT_EVENT:
		var ce ClientEvent
		ce.Event_type = r.Read_Byte()
		ce.Car_id = r.Read_Byte()
		if ce.Event_type == ACSP_CE_COLLISION_WITH_CAR {
			ce.Other_car_id = r.Read_Byte()
		}
		ce.Impact_speed = r.Read_Float()
		ce.World_pos = Vector{r.Read_Float(), r.Read_Float(), r.Read_Float()}
		ce.Rel_pos = Vector{r.Read_Float(), r.Read_Float(), r.Read_Float()}
		log.Print("ACSP_CLIENT_EVENT: ")
		Print_Interface(ce)

	case ACSP_CAR_INFO:
		var ci CarInfo
		ci.Car_id = r.Read_Byte()
		ci.Is_connected = r.Read_Byte() != 0
		ci.Car_model = r.Read_UTF32_String()
		ci.Car_skin = r.Read_UTF32_String()
		ci.Driver_name = r.Read_UTF32_String()
		ci.Driver_team = r.Read_UTF32_String()
		ci.Driver_guid = r.Read_UTF32_String()
		log.Print("ACSP_CAR_INFO: ")
		Print_Interface(ci)

	case ACSP_CAR_UPDATE:
		var cu CarUpdate
		cu.Car_id = r.Read_Byte()
		cu.Position = Vector{r.Read_Float(), r.Read_Float(), r.Read_Float()}
		cu.Velocity = Vector{r.Read_Float(), r.Read_Float(), r.Read_Float()}
		cu.Gear = r.Read_Byte()
		cu.Engine_rpm = r.Read_Uint16()
		cu.Normalized_spline_pos = r.Read_Float()
		log.Print("ACSP_CAR_UPDATE: ")
		Print_Interface(cu)

	case ACSP_NEW_CONNECTION:
		var nc NewConnection
		nc.Driver_name = r.Read_UTF32_String()
		nc.Driver_guid = r.Read_UTF32_String()
		nc.Car_id = r.Read_Byte()
		nc.Car_model = r.Read_String()
		nc.Car_skin = r.Read_String()
		Status.Players = Status.Players + 1
		log.Print("ACSP_NEW_CONNECTION: ")
		Print_Interface(nc)

	case ACSP_CONNECTION_CLOSED:
		var cc ConnectionClosed
		cc.Driver_name = r.Read_UTF32_String()
		cc.Driver_guid = r.Read_UTF32_String()
		cc.Car_id = r.Read_Byte()
		cc.Car_model = r.Read_String()
		cc.Car_skin = r.Read_String()
		Status.Players = Status.Players - 1
		log.Print("ACSP_CONNECTION_CLOSED: ")
		Print_Interface(cc)

	case ACSP_LAP_COMPLETED:
		var lc LapCompleted
		lc.Car_id = r.Read_Byte()
		lc.Laptime = r.Read_Uint32()
		lc.Cuts = r.Read_Byte()
		log.Print("ACSP_LAP_COMPLETED: ")
		Print_Interface(lc)

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
