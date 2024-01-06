package migrate

import (
	"fmt"
	"github.com/lanyulei/toolkit/db"
	"github.com/lanyulei/toolkit/logger"
	"os"
	"strings"

	"gorm.io/gorm"
)

func InitDB(db *gorm.DB, filePath string) (err error) {
	err = ExecSQL(db, filePath)
	return err
}

func ExecSQL(db *gorm.DB, filePath string) error {
	sql, err := Ioutil(filePath)
	if err != nil {
		logger.Errorf("database base data initialization script read failed, error: %s", err.Error())
		return err
	}
	if err = db.Exec(sql).Error; err != nil {
		logger.Errorf("error sql: %s", sql)
		if !strings.Contains(err.Error(), "Query was empty") {
			return err
		}
	}
	return nil
}

func Ioutil(filePath string) (string, error) {
	if contents, err := os.ReadFile(filePath); err == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := strings.Replace(string(contents), "\n", "", 1)
		return result, nil
	} else {
		logger.Errorf("database base data initialization script read failed, error: %s", err.Error())
		return "", err
	}
}

func ExecSQLValue(sql string) (err error) {
	fmt.Println(sql)
	tx := db.Orm().Begin()
	if err = tx.Exec(sql).Error; err != nil {
		tx.Rollback()
		logger.Errorf("error sql: %s", sql)
		if !strings.Contains(err.Error(), "Query was empty") {
			return
		}
	}
	tx.Commit()
	return
}
