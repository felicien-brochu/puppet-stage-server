package db

import (
	"encoding/json"
	"felicien/puppet-server/model"

	"github.com/garyburd/redigo/redis"
)

// GetPuppet returns the puppet with the given id if present in db
func GetPuppet(id string) (*model.Puppet, error) {
	exists, err := redis.Bool(conn.Do("EXISTS", "puppet:"+id))
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	puppetJSON, err := redis.Bytes(conn.Do("GET", "puppet:"+id))
	if err != nil {
		return nil, err
	}

	var puppet model.Puppet
	err = json.Unmarshal(puppetJSON, &puppet)
	if err != nil {
		return nil, err
	}

	return &puppet, nil
}

// ListPuppets retrieves puppets from redis
func ListPuppets() ([]model.Puppet, error) {
	puppetKeys, err := redis.Values(conn.Do("KEYS", "puppet:*"))
	if err != nil {
		return nil, err
	}
	if len(puppetKeys) == 0 {
		return make([]model.Puppet, 0), nil
	}
	puppetsJSON, err := redis.ByteSlices(conn.Do("MGET", puppetKeys...))
	if err != nil {
		return nil, err
	}

	var puppets []model.Puppet

	for _, puppetJSON := range puppetsJSON {
		var puppet model.Puppet
		err = json.Unmarshal(puppetJSON, &puppet)
		if err != nil {
			return nil, err
		}

		puppets = append(puppets, puppet)
	}

	return puppets, nil
}

// SavePuppet saves a puppet (CREATE or UPDATE)
func SavePuppet(puppet model.Puppet) error {
	puppetJSON, err := json.Marshal(puppet)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", "puppet:"+puppet.ID, puppetJSON)
	if err != nil {
		return err
	}
	return nil
}
