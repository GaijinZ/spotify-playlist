package sql

import (
	"context"
	"fmt"
	"spf-playlist/users/handler/models"
	"time"

	"spf-playlist/pkg/config"
	"spf-playlist/pkg/logger"
	"spf-playlist/utils"

	"github.com/gocql/gocql"
)

type DBer interface {
	Insert(item interface{}) error
	Update(id int) error
	Get(email string) *gocql.Query
	Delete(id int) error
	Close()
}

type DB struct {
	ctx    context.Context
	Client *gocql.Session
}

func InitDB(log logger.Logger, cfg config.GlobalEnv, ctx context.Context) (DBer, error) {
	cluster := gocql.NewCluster(cfg.ClusterIP)
	cluster.Keyspace = cfg.KeySpace
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10

	session, err := cluster.CreateSession()
	if err != nil {
		log.Errorf("Failed to create session: %v", err)
		return nil, err
	}

	log.Infof("Connected to cluster: %s with keyspace: %s", cfg.ClusterIP, cfg.KeySpace)
	db := &DB{Client: session, ctx: ctx}

	return db, nil
}

func (d *DB) Close() {
	d.Client.Close()
}

func (d *DB) Insert(item interface{}) error {
	log := utils.GetLogger(d.ctx)

	switch v := item.(type) {
	case *models.User:
		id := gocql.TimeUUID()
		err := d.Client.Query(InsertUser, id, v.Name, v.Email, v.Password, v.Role).Exec()
		if err != nil {
			log.Errorf("Failed to insert user: %v", err)
			return err
		}
	default:
		log.Errorf("Unexpected type for Insert: %T", item)
		return fmt.Errorf("unexpected type for Insert")
	}

	log.Infof("Inserting item: %v", item)
	return nil
}

func (d *DB) Update(id int) error {
	return nil
}

func (d *DB) Get(email string) *gocql.Query {
	query := d.Client.Query(GetUser, email)

	return query
}

func (d *DB) Delete(id int) error {
	return nil
}
