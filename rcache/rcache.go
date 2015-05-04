// @author: Maksadbek
// @email: a.maksadbek@gmail.com:
/*
   пакет для кеширования данных 
*/

package rcache

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"github.com/garyburd/redigo/redis"
)

var (
	config conf.App   // config
	rc     redis.Conn // redis client
)

// структура для позиции трекера
type Pos struct {
	Id            int     `json:"id"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
	Time          string  `json:"time"`
	Owner         string  `json:"owner"`
	Number        string  `json:"number"`
	Name          string  `json:"name"`
	Direction     int     `json:"direction"`
	Speed         int     `json:"speed"`
	Sat           int     `json:"sat"`
	Ignition      int     `json:"ignition"`
	GsmSignal     int     `json:"gsmsignal"`
	Battery       int     `json:"battery66"`
	Seat          int     `json:"seat"`
	BatteryLvl    int     `json:"batterylvl"`
	Fuel          int     `json:"fuel"`
	FuelVal       int     `json:"fuel_val"`
	MuAdditional  string  `json:"mu_additional"`
	Customization string  `json:"customization"`
	Additional    string  `json:"additional"`
	Action        int     `json:"action"`
}

// структура для флита
type Fleet struct {
	Id     string
	Update map[string]Pos
}

// функция для инициализации пакета
// оно должна вызыватся первые перед исползованием пакета
func Initialize(c conf.App) (err error) {
	config = c
	rc, err = redis.Dial("tcp", c.DS.Redis.Host)
	if err != nil {
		return err
	}
	return
}

// функция используется для получения трекеров флита
func GetTrackers(fleet string, start, stop int) (trackers []string, err error) {
	v, err := redis.Strings(rc.Do("LRANGE", fleet, start, stop))
	if err != nil {
		return
	}

	for _, val := range v {
		trackers = append(trackers, val)
	}
	return
}

// функция исползуется для вставления данный в редис
func PushRedis(fleet Fleet) (err error) {
	for k, x := range fleet.Update {
		jpos, err := json.Marshal(x)
		if err != nil {
			return err
		}
		rc.Do("RPUSH", config.DS.Redis.TPrefix+":"+k, jpos)
	}
	return
}

// исползуется для получения позиции трекеров флита
func GetPositions(fleetNum string, start, stop int) (Fleet, error) {
	fleet := Fleet{}
	fleet.Id = fleetNum
	fleet.Update = make(map[string]Pos)
	trackers, err := GetTrackers(fleetNum, start, stop)
	if err != nil {
		return fleet, err
	}

	for _, v := range trackers {
		pBytes, err := rc.Do(
			"LINDEX",
			config.DS.Redis.TPrefix+":"+v,
			-1)
		if err != nil {
			return fleet, err
		}
		p := fmt.Sprintf("%s", pBytes)
		var pos Pos
		err = json.Unmarshal([]byte(p), &pos)
		if err != nil {
			return fleet, err
		}

		fleet.Update[v] = pos
	}
	return fleet, err
}
