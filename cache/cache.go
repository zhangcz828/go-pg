package cache

import (
	"context"
	"encoding/json"
	"github.com/golang/groupcache"
	"go-pg/modules"
	"go-pg/pkg/connection"
	"log"
)

var HerosGroup *groupcache.Group

func SetCacheHttpPool(addr string, peers_addr []string) {
	peers := groupcache.NewHTTPPool(addr)
	peers.Set(peers_addr...)
}

// 单例模式， groupcache.Group 必须唯一
func CreateHerosCacheGroup() *groupcache.Group {
	if HerosGroup == nil {
		HerosGroup = groupcache.NewGroup("getHeros", 8<<30,
			groupcache.GetterFunc(
				func(ctx context.Context, key string, dest groupcache.Sink) error {
					heros, err := getAllHero()
					if err != nil {
						log.Fatalf("Unable to get all heros. %v", err)
					}

					log.Printf("Get %s of value from database.\n",key)
					herosB, _ := json.Marshal(heros)
					// dest.SetBytes([]byte(fmt.Sprintf("%v", heros)))
					dest.SetBytes(herosB)
					return nil
				}))
	}
	// 获取group对象
	return HerosGroup
}

func getAllHero() ([]modules.Hero, error){
	db := connection.CreateConnection()
	defer db.Close()

	var heros []modules.Hero

	// create the select sql query
	sqlStatement := `SELECT * FROM Hero`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	//iterate over the rows
	for rows.Next() {
		var hero modules.Hero

		// unmarshal the row object to hero
		err = rows.Scan(&hero.Name, &hero.Detail, &hero.AttackPower, &hero.DefensePower, &hero.Blood)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the hero to the heros slice
		heros = append(heros, hero)
	}

	return heros, err
}