package dbmysql

import (
	"EventFlow/common/tool/arraytool"
	"EventFlow/common/tool/cachetool"
	"EventFlow/common/tool/parametertool"
	"database/sql"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type MySQLPlugin struct {
	Setting SettingConfig
}

var replacer = strings.NewReplacer("[", "", "]", "")

func (filter *MySQLPlugin) DoFilter(messageFromTrigger *string, parameters *map[string]interface{}) (doNextPipeline bool) {

	cacheKey := parametertool.ReplaceWithParameter(&filter.Setting.Cache.CacheKey, parameters)
	recordMaps, existed := cachetool.GetCache(cacheKey)

	if existed {

		if recordMaps, ok := recordMaps.([]map[string]interface{}); ok {
			for metadataKey, metadataParm := range filter.Setting.AddMetadata {

				metadataParm = replacer.Replace(metadataParm)
				results := []interface{}{}

				for _, record := range recordMaps {
					results = append(results, record[metadataParm])
				}

				(*parameters)[metadataKey] = results
			}

			return true
		}
	}

	dbConfig := mysql.Config{
		Net:    "tcp",
		Addr:   filter.Setting.IP,
		User:   filter.Setting.User,
		Passwd: filter.Setting.Password,
		DBName: filter.Setting.Database,
	}

	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	defer close(db)

	if err != nil {
		log.Printf("[filter][mysql] open db connection failed: %v", err)
	} else {
		command := parametertool.ReplaceWithParameter(&filter.Setting.Command, parameters)

		rows, err := db.Query(command)
		defer rows.Close()

		if err != nil {
			log.Printf("[filter][mysql] query db failed: %v", err)
		} else {

			recordMaps := []map[string]interface{}{}
			columnList, _ := rows.Columns()
			columnCount := len(columnList)

			for rows.Next() {
				columns := make([]interface{}, columnCount)
				columnPointers := make([]interface{}, columnCount)

				for index := range columnList {
					columnPointers[index] = &columns[index]
				}

				if err := rows.Scan(columnPointers...); err != nil {
					continue
				}

				recordMap := make(map[string]interface{})

				for index, columnName := range columnList {

					value, ok := columns[index].([]byte)

					if ok {
						recordMap[columnName] = string(value)
					} else {
						recordMap[columnName] = value
					}

				}

				recordMaps = append(recordMaps, recordMap)
			}

			if cacheKey != "" {
				cachetool.CreateCache(cacheKey, filter.Setting.Cache.TimeoutSecond, recordMaps)
			}

			for metadataKey, metadataParm := range filter.Setting.AddMetadata {

				metadataParm = replacer.Replace(metadataParm)
				results := []interface{}{}

				if existed, _ := arraytool.InArray(metadataParm, columnList); existed {
					for _, record := range recordMaps {
						results = append(results, record[metadataParm])
					}
				}

				(*parameters)[metadataKey] = results
			}
		}
	}

	return true
}

func close(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
