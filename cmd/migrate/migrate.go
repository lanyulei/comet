package migrate

import (
	"comet/app/system/models"
	migrate "comet/cmd/migrate/models"
	"comet/cmd/migrate/sql"
	"fmt"
	"github.com/lanyulei/toolkit/db"
	"github.com/lanyulei/toolkit/logger"
	"github.com/liushuochen/gotable"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

func autoMigrate(data bool) {
	var (
		migrateList   []*models.Migrate
		migrateModels []interface{}
		migrateMaps   = make(map[string]struct{})
	)

	// 同步数据结构
	logger.Info("start synchronizing data structures...")

	migrateModels = append(migrateModels, migrate.SystemModels...)

	err := db.Orm().AutoMigrate(
		migrateModels...,
	)

	if err != nil {
		logger.Fatalf("failed to synchronize data structure, error: %v", err)
	}
	logger.Info("data structure synchronization completed")
	if data {
		// 同步初始数据
		logger.Info("start synchronizing base data...")

		// 查询迁移记录
		err = db.Orm().Select("id, name").Find(&migrateList).Error
		if err != nil {
			logger.Fatalf("failed to query migration records, error: %v", err)
		}

		// 将迁移记录转换为map
		for _, v := range migrateList {
			migrateMaps[v.Name] = struct{}{}
		}

		for _, m := range sql.ListSQL {
			for k, v := range m {
				if _, ok := migrateMaps[k]; !ok {
					logger.Infof("start migration of data structures with file [%s]", k)
					err = ExecSQLValue(v)
					if err != nil {
						logger.Fatalf("failed to synchronize initial data, error: %v", err)
						err = createMigrateRecord(k, "failed", err.Error())
						if err != nil {
							logger.Errorf("failed to create migration record, error: %v", err)
						}
					}

					// 创建迁移记录
					err = createMigrateRecord(k, "success", "")
					if err != nil {
						logger.Fatalf("failed to create migration record, error: %v", err)
					}
				}
			}
		}

		logger.Info("basic data synchronization completed")
	}
	return
}

func createMigrateRecord(name, status, result string) (err error) {
	err = db.Orm().Create(&models.Migrate{
		Name:   name,
		Status: status,
		Result: result,
	}).Error
	return
}

func clearMigrate(value string) {
	var (
		err         error
		recordCount int64
	)

	if value == "all" {
		// emptying migration records
		err = db.Orm().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Migrate{}).Error
		if err != nil {
			logger.Fatalf("failed to empty migration records, error: %v", err)
		}
		logger.Info("migration records emptied")
	} else {
		err = db.Orm().Model(&models.Migrate{}).Where("name = ?", value).Count(&recordCount).Error
		if err != nil {
			logger.Fatalf("failed to query migration records, error: %v", err)
		}
		if recordCount == 0 {
			logger.Fatalf("no migration records found")
		}

		// clear the specified migration record
		err = db.Orm().Delete(&models.Migrate{}, "name = ?", value).Error
		if err != nil {
			logger.Fatalf("failed to clear migration records, error: %v", err)
		}
		logger.Infof("[%s] migration record cleared successfully", value)
	}
}

func listMigrate() {
	var (
		migrateList []map[string]interface{}
		values      []string
	)
	err := db.Orm().Table("system_migrate").
		Where("delete_time IS NULL").
		Find(&migrateList).Error
	if err != nil {
		logger.Fatalf("failed to query migration records, error: %v", err)
	}
	if len(migrateList) > 0 {
		table, err := gotable.Create("id", "name", "status", "create_time", "update_time")
		if err != nil {
			logger.Fatalf("failed to create table, error: %v", err)
		}

		for _, value := range migrateList {
			values = make([]string, 0)
			idString := strconv.Itoa(int(value["id"].(int64)))
			values = append(values,
				idString,
				value["name"].(string),
				value["status"].(string),
				value["create_time"].(time.Time).Format("2006-01-02 15:04:05"),
				value["update_time"].(time.Time).Format("2006-01-02 15:04:05"),
			)
			err := table.AddRow(values)
			if err != nil {
				logger.Fatalf("failed to add row, error: %v", err)
			}
		}
		fmt.Println("\nmigration records: ")
		fmt.Println(table)
	} else {
		fmt.Println("\nno migration records")
		fmt.Println()
	}
}

func generateSQL() {
	var (
		err error
	)

	for _, m := range sql.ListSQL {
		for k, v := range m {
			_, err = os.Stat(migrateDir)
			if err != nil {
				if os.IsNotExist(err) {
					err = os.MkdirAll(migrateDir, 0755)
					if err != nil {
						logger.Fatalf("failed to create migration directory, error: %v", err.Error())
					}
				} else {
					logger.Fatalf("failed to create migration directory, error: %v", err.Error())
				}
			}

			// 写入文件
			err = os.WriteFile(fmt.Sprintf("%s/%s.sql", migrateDir, k), []byte(v), 0644)
			if err != nil {
				logger.Fatalf("failed to write file, error: %v", err.Error())
			}
		}
	}
}
