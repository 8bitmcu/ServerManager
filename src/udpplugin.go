package main

import (
	"encoding/binary"
	"log"
	"math"
	"net"
	"strconv"

	"golang.org/x/text/encoding/unicode/utf32"
)

const (
	acspProtocolVersion     = 4
	acspNewSession          = 50
	acspNewConnection       = 51
	acspConnectionClosed    = 52
	acspCarUpdate           = 53
	acspCarInfo             = 54
	acspEndSession          = 55
	acspLapCompleted        = 73
	acspVersion             = 56
	acspChat                = 57
	acspClientLoaded        = 58
	acspSessionInfo         = 59
	acspError               = 60
	acspClientEvent         = 130
	acspCeCollisionWithCar  = 10
	acspCeCollisionWithEnv  = 11
	acspRealtimeposInterval = 200
	acspGetCarInfo          = 201
	acspSendChat            = 202
	acspBroadcastChat       = 203
	acspGetSessionInfo      = 204
	acspSetSessionInfo      = 205
	acspKickUser            = 206
	acspNextSession         = 207
	acspRestartSession      = 208
	acspAdminCommand        = 209
)

type SessionInfo struct {
	version             int
	sessionIndex        int
	currentSessionIndex int
	sessionCount        int
	serverName          string
	track               string
	trackConfig         string
	name                string
	typ                 int
	time                int
	laps                int
	waitTime            int
	ambientTemp         int
	roadTemp            int
	weatherGraphics     string
	elapsedMs           int32
}

type ClientEvent struct {
	carId       int
	otherCarId  int
	eventType   int
	impactSpeed float32
	worldPos    Vector
	relPos      Vector
}

type CarInfo struct {
	carId       int
	isConnected bool
	carModel    string
	carSkin     string
	driverName  string
	driverTeam  string
	driverGuid  string
}

type CarUpdate struct {
	carId               int
	position            Vector
	velocity            Vector
	gear                int
	engineRpm           int
	normalizedSplinePos float32
}

type NewConnection struct {
	driverName string
	driverGuid string
	carId      int
	carModel   string
	carSkin    string
}

type ConnectionClosed struct {
	driverName string
	driverGuid string
	carId      int
	carModel   string
	carSkin    string
}

type LapCompleted struct {
	carId   int
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

func (r *UdpReader) ReadByte() int {
	val := r.data[r.readIndex]
	r.readIndex = r.readIndex + 1
	return int(val)
}

func (r *UdpReader) ReadBytes(bytes int) []byte {
	val := r.data[r.readIndex : r.readIndex+int64(bytes)]
	r.readIndex = r.readIndex + int64(bytes)
	return val
}

func (r *UdpReader) ReadString() string {
	length := int(r.ReadByte())
	bytes := r.ReadBytes(length)
	return string(bytes)
}

func (r *UdpReader) ReadUTF32String() string {
	length := int(r.ReadByte() * 4)
	bytes := r.ReadBytes(length)

	val, err := utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewDecoder().Bytes(bytes)

	if err != nil {
		log.Print(err)
	}

	return string(val)
}

func (r *UdpReader) ReadUint16() int {
	bytes := r.ReadBytes(2)
	data := binary.LittleEndian.Uint16(bytes)
	return int(data)
}

func (r *UdpReader) ReadInt32() int32 {
	bytes := r.ReadBytes(4)
	data := binary.LittleEndian.Uint32(bytes)
	return int32(data)
}

func (r *UdpReader) ReadUint32() uint32 {
	bytes := r.ReadBytes(4)
	data := binary.LittleEndian.Uint32(bytes)
	return data
}

func (r *UdpReader) ReadFloat() float32 {
	bytes := r.ReadBytes(4)
	a := binary.LittleEndian.Uint32(bytes)
	a2 := math.Float32frombits(a)
	return a2
}

type UdpWriter struct {
	data []byte
}

func (w *UdpWriter) WriteByte(d byte) {
	w.data = append(w.data, d)
}

func (w *UdpWriter) WriteUTF32String(str string) {
	val, err := utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewEncoder().Bytes([]byte(str))

	w.WriteByte(byte(len(str)))

	if err != nil {
		log.Print(err)
	}
	w.data = append(w.data, val[:]...)
}

func readSessionInfo(r UdpReader) SessionInfo {
	var s SessionInfo
	s.version = r.ReadByte()
	s.sessionIndex = r.ReadByte()
	s.currentSessionIndex = r.ReadByte()
	s.sessionCount = r.ReadByte()
	s.serverName = r.ReadUTF32String()
	s.track = r.ReadString()
	s.trackConfig = r.ReadString()
	s.name = r.ReadString()
	s.typ = r.ReadByte()
	s.time = r.ReadUint16()
	s.laps = r.ReadUint16()
	s.waitTime = r.ReadUint16()
	s.ambientTemp = r.ReadByte()
	s.roadTemp = r.ReadByte()
	s.weatherGraphics = r.ReadString()
	s.elapsedMs = r.ReadInt32()
	return s
}

type UdpPlugin struct {
	conn   *net.UDPConn
	online bool
}

func udpListen() UdpPlugin {
	var udp UdpPlugin
	udpClient, err := net.ResolveUDPAddr("udp", ":5001")
	if err != nil {
		log.Print(err)
	}

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

func (udp UdpPlugin) Receive() {
	data := make([]byte, 1024)
	udp.conn.Read(data)

	r := UdpReader{}
	r.New(data)

	acsp := r.ReadByte()

	switch acsp {
	case acspError:
		err := r.ReadUTF32String()
		log.Print("ACSP_ERROR: ", err)

	case acspChat:
		car := r.ReadByte()
		msg := r.ReadUTF32String()
		log.Print("ACSP_CHAT: " + strconv.Itoa(car) + "; " + msg)

	case acspClientLoaded:
		car := r.ReadByte()
		log.Print("ACSP_CLIENT_LOADED: ", car)

	case acspVersion:
		v := r.ReadByte()
		log.Print("ACSP_VERSION: ", v)
		Udp.online = true

	case acspNewSession:
		sess := readSessionInfo(r)
		log.Print("ACSP_NEW_SESSION: ")
		PrintInterface(sess)
		Status.Session = sess

	case acspSessionInfo:
		sess := readSessionInfo(r)
		log.Print("ACSP_SESSION_INFO: ")
		PrintInterface(sess)
		Status.Session = sess

	case acspEndSession:
		file := r.ReadUTF32String()
		log.Print("ACSP_END_SESSION: " + file)

		if Status.Session.currentSessionIndex == Status.Session.sessionCount-1 {
			log.Print("UDP Plugin triggers Server Change Track")
			Status.serverChangeTrack()
		}

	case acspClientEvent:
		var ce ClientEvent
		ce.eventType = r.ReadByte()
		ce.carId = r.ReadByte()
		if ce.eventType == acspCeCollisionWithCar {
			ce.otherCarId = r.ReadByte()
		}
		ce.impactSpeed = r.ReadFloat()
		ce.worldPos = Vector{r.ReadFloat(), r.ReadFloat(), r.ReadFloat()}
		ce.relPos = Vector{r.ReadFloat(), r.ReadFloat(), r.ReadFloat()}
		log.Print("ACSP_CLIENT_EVENT: ")
		PrintInterface(ce)

	case acspCarInfo:
		var ci CarInfo
		ci.carId = r.ReadByte()
		ci.isConnected = r.ReadByte() != 0
		ci.carModel = r.ReadUTF32String()
		ci.carSkin = r.ReadUTF32String()
		ci.driverName = r.ReadUTF32String()
		ci.driverTeam = r.ReadUTF32String()
		ci.driverGuid = r.ReadUTF32String()
		log.Print("ACSP_CAR_INFO: ")
		PrintInterface(ci)

	case acspCarUpdate:
		var cu CarUpdate
		cu.carId = r.ReadByte()
		cu.position = Vector{r.ReadFloat(), r.ReadFloat(), r.ReadFloat()}
		cu.velocity = Vector{r.ReadFloat(), r.ReadFloat(), r.ReadFloat()}
		cu.gear = r.ReadByte()
		cu.engineRpm = r.ReadUint16()
		cu.normalizedSplinePos = r.ReadFloat()
		log.Print("ACSP_CAR_UPDATE: ")
		PrintInterface(cu)

	case acspNewConnection:
		var nc NewConnection
		nc.driverName = r.ReadUTF32String()
		nc.driverGuid = r.ReadUTF32String()
		nc.carId = r.ReadByte()
		nc.carModel = r.ReadString()
		nc.carSkin = r.ReadString()
		Status.Players = Status.Players + 1
		log.Print("ACSP_NEW_CONNECTION: ")
		PrintInterface(nc)

	case acspConnectionClosed:
		var cc ConnectionClosed
		cc.driverName = r.ReadUTF32String()
		cc.driverGuid = r.ReadUTF32String()
		cc.carId = r.ReadByte()
		cc.carModel = r.ReadString()
		cc.carSkin = r.ReadString()
		Status.Players = Status.Players - 1
		log.Print("ACSP_CONNECTION_CLOSED: ")
		PrintInterface(cc)

	case acspLapCompleted:
		var lc LapCompleted
		lc.carId = r.ReadByte()
		lc.laptime = r.ReadUint32()
		lc.cuts = r.ReadByte()
		log.Print("ACSP_LAP_COMPLETED: ")
		PrintInterface(lc)

	default:
		log.Print("ACSP Unknown code: "+strconv.Itoa(acsp), data)
	}
}

func (udp UdpPlugin) WriteAdminCommand(command string) {
	var w UdpWriter
	w.WriteByte(acspAdminCommand)
	w.WriteUTF32String(command)
	udp.conn.Write(w.data)
}

func (udp UdpPlugin) WriteBroadcastChat(message string) {
	var w UdpWriter
	w.WriteByte(acspBroadcastChat)
	w.WriteUTF32String(message)
	udp.conn.Write(w.data)
}

func (udp UdpPlugin) WriteGetCarInfo(carid int) {
	var w UdpWriter
	w.WriteByte(acspGetCarInfo)
	w.WriteByte(byte(carid))
	udp.conn.Write(w.data)
}

func (udp UdpPlugin) WriteKickUser(carid int) {
	var w UdpWriter
	w.WriteByte(acspKickUser)
	w.WriteByte(byte(carid))
	udp.conn.Write(w.data)
}

func (udp UdpPlugin) WriteNextSession() {
	var w UdpWriter
	w.WriteByte(acspNextSession)
	udp.conn.Write(w.data)
}

func (udp UdpPlugin) WriteRestartSession() {
	var w UdpWriter
	w.WriteByte(acspRestartSession)
	udp.conn.Write(w.data)
}

func (udp UdpPlugin) WriteSendChat(carid int, message string) {
	var w UdpWriter
	w.WriteByte(acspSendChat)
	w.WriteByte(byte(carid))
	w.WriteUTF32String(message)
	udp.conn.Write(w.data)
}
