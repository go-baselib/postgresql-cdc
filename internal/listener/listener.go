package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/1005281342/postgresql-cdc/internal/config"
	"github.com/1005281342/postgresql-cdc/internal/publisher"
	"github.com/1005281342/postgresql-cdc/model"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx"
)

type Listener struct {
	// 生产/发布者
	publisher message.Publisher

	buffer []byte
}

func build() *Listener {
	var pub = publisher.CreatePublisher(config.GetConfig().EventQueen)
	return &Listener{publisher: pub}
}

func (l *Listener) start() error {
	var connConfig = buildConnConfig()
	log.Printf("connConfig %+v \n", connConfig)
	var conn, err = pgx.ReplicationConnect(connConfig)
	if err != nil {
		return err
	}

	var lsnAsInt uint64
	lsnAsInt, err = pgx.ParseLSN(l.lsn())
	if err != nil {
		return err
	}

	var cfg = config.GetConfig()
	var outputParams []string
	if cfg.Source.TableNames != "" {
		outputParams = append(outputParams, fmt.Sprintf("\"add-tables\" '%s'", cfg.Source.TableNames))
	}
	if cfg.Source.Chunks != "" {
		outputParams = append(outputParams, fmt.Sprintf("\"write-in-chunks\" '%s'", cfg.Source.Chunks))
	}

	log.Printf("slot: %s, lsnAsInt: %d, outputParams: %+v \n", cfg.Slot, lsnAsInt, outputParams)
	if err = conn.StartReplication(cfg.Source.Slot, lsnAsInt, -1, outputParams...); err != nil {
		return err
	}

	l.listen(conn)
	return nil
}

//listen is main infinite process for handle this replication slot
func (l *Listener) listen(rc *pgx.ReplicationConn) {
	log.Print("listen start \n")
	for {
		var r, err = rc.WaitForReplicationMessage(context.Background())
		if err != nil {
			panic(err)
		}

		if r == nil {
			continue
		}

		if r.ServerHeartbeat != nil {
			log.Printf("server heartbeat received: %d", r.ServerHeartbeat.ServerWalEnd)
			l.sendStandBy(rc, r.ServerHeartbeat.ServerWalEnd)
		}

		if r.WalMessage != nil {
			log.Printf("wal message start: %d data: %s", r.WalMessage.WalStart, string(r.WalMessage.WalData))
			l.handleMessage(rc, r.WalMessage)
		}
	}
}

//handleMessage parses WAL data, send message to Kafka, and sends standby status
func (l *Listener) handleMessage(rc *pgx.ReplicationConn, wm *pgx.WalMessage) {
	if len(wm.WalData) == 0 {
		return
	}

	l.buffer = append(l.buffer, wm.WalData...)

	var (
		change model.Wal2JsonMessage
		err    error
	)
	//trying to deserialize to JSON
	if err = json.Unmarshal(l.buffer, &change); err == nil {
		var uid = watermill.NewUUID()
		var msg = message.NewMessage(
			uid, // internal uuid of the message, useful for debugging
			l.buffer,
		)
		log.Printf("buffer: %s\n", l.buffer)
		log.Printf("change: %+v\n", change)

		var (
			cfg   = config.GetConfig()
			retry int
		)
		for retry < cfg.EventQueen.MaxRetry {
			if err = l.publisher.Publish(cfg.EventQueen.Name, msg); err != nil {
				retry++
				time.Sleep(cfg.EventQueen.RetryInterval * time.Millisecond)
				continue
			}

			log.Printf("EventQueen.Name: %s, msg: %+v \n", cfg.EventQueen.Name, msg)
			l.buffer = []byte{}
			break
		}

		if err != nil {
			log.Printf("uuid: %+v 消息: %+v 发送失败: %+v\n", uid, msg, err)
		}
	}

	l.sendStandBy(rc, wm.WalStart)
}

func (l *Listener) sendStandBy(rc *pgx.ReplicationConn, lastWal uint64) {
	var (
		status *pgx.StandbyStatus
		err    error
	)
	if status, err = pgx.NewStandbyStatus(lastWal); err != nil {
		log.Printf("%+v\n", err)
		return
	}

	if err = rc.SendStandbyStatus(status); err != nil {
		log.Printf("%+v\n", err)
		return
	}
}

func (l *Listener) lsn() string {
	var s = config.GetConfig().Lsn
	if s != "" {
		return s
	}
	return "0/0"
}

func buildConnConfig() pgx.ConnConfig {
	var cfg = config.GetConfig()

	var port, _ = strconv.ParseUint(cfg.Source.DbPort, 10, 16)

	return pgx.ConnConfig{
		Host:     cfg.Source.DbHost,
		Port:     uint16(port),
		Database: cfg.Source.DbName,
		User:     cfg.Source.DbUser,
		Password: cfg.Source.DbPass,
	}
}

func Start() error {
	return build().start()
}
