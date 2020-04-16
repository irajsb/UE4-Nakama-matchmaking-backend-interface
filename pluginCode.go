package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

var maximumlastping int64 = 9000 // maximum last ping in miliSeconds

type serverport struct {
	port int `json:"port"`
}
type dedicatedserver struct {
	ip              string
	port            string
	LastPing        time.Time
	IsPendingDelete bool
}

var ServerList []dedicatedserver

func RemoveIndex(s []dedicatedserver, index int) []dedicatedserver {
	return append(s[:index], s[index+1:]...)

}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {

	// Register as an  function, this call should be in InitModule.
	if err := initializer.RegisterMatchmakerMatched(MakeMatch); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	// Register as an RPC function, this call should be in InitModule.
	if err := initializer.RegisterRpc("http_register_server", RegisterServer); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	return nil
}
func RegisterServer(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	var message map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &message); err != nil {
		return "", err
	}
	clientIP, ok := ctx.Value(runtime.RUNTIME_CTX_CLIENT_IP).(string)
	if !ok {
		// handle error
	}

	var temp dedicatedserver
	temp.ip = clientIP
	temp.port = message["port"].(string)
	temp.LastPing = time.Now()
	var Serverresponse string
	logger.Info("Server  ip:", clientIP, temp.port)

	for i, It := range ServerList {
		if It.ip == temp.ip && It.port == temp.port {
			//if server already exists just set last ping to now
			Serverresponse = "Pinged"
			ServerList[i].LastPing = time.Now()
			if It.IsPendingDelete == true {
				ServerList = RemoveIndex(ServerList, i)
				Serverresponse = "Matched"
				break
			}
		} else {
			//create new server
			ServerList = append(ServerList, temp)
			Serverresponse = "Created"
		}
	}

	if len(ServerList) == 0 {
		//slice is empty so for each loop doenst work create server manually
		ServerList = append(ServerList, temp)
		Serverresponse = "Created"

	}
	response, err := json.Marshal(map[string]interface{}{"Result": Serverresponse})
	if err != nil {
		return "", err
	}
	//logger.Info("numOfSevers", len(ServerList))
	return string(response), nil

}

func MakeMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, entries []runtime.MatchmakerEntry) (string, error) {
	/*	for _, e := range entries {
		logger.Info("Matched user '%s' named '%s'", e.GetPresence().GetUserId(), e.GetPresence().GetUsername())
		for k, v := range e.GetProperties() {
			logger.Info("Matched on '%s' value '%v'", k, v)
		}
	}*/

	matchId := "Error"
	for i, ele := range ServerList {

		if (time.Now().Sub(ele.LastPing).Milliseconds()) > maximumlastping {
			ServerList = RemoveIndex(ServerList, i)

		} else if ele.IsPendingDelete == false {

			matchId = ele.ip + ":" + ele.port
			ServerList = RemoveIndex(ServerList, i)
			ele.IsPendingDelete = true
			break
		}

	}

	return matchId, nil
}
