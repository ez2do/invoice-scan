package debug

import (
	"fmt"
	"gorm.io/gorm"
	"invoice-scan/backend/pkg"
	"invoice-scan/backend/pkg/log"
	"time"
)

func NewGormPlugin() *GormPlugin {
	return &GormPlugin{}
}

type GormPlugin struct{}

func (p GormPlugin) Name() string {
	return "gorm_debugger"
}

func (p GormPlugin) Initialize(db *gorm.DB) (err error) {
	// create
	err = db.Callback().Create().Before("gorm:create").Register("debugGormPluginBeforeCreate", p.before)
	if err != nil {
		return err
	}

	err = db.Callback().Create().After("gorm:create").Register("debugGormPluginAfterCreate", p.after)
	if err != nil {
		return err
	}

	// update
	err = db.Callback().Update().Before("gorm:update").Register("debugGormPluginBeforeUpdate", p.before)
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("gorm:update").Register("debugGormPluginAfterUpdate", p.after)
	if err != nil {
		return err
	}

	// query
	err = db.Callback().Query().Before("gorm:query").Register("debugGormPluginBeforeQuery", p.before)
	if err != nil {
		return err
	}

	err = db.Callback().Query().After("gorm:query").Register("debugGormPluginAfterQuery", p.after)
	if err != nil {
		return err
	}

	// delete
	err = db.Callback().Delete().Before("gorm:delete").Register("debugGormPluginBeforeDelete", p.before)
	if err != nil {
		return err
	}

	err = db.Callback().Delete().After("gorm:delete").Register("debugGormPluginAfterDelete", p.after)
	if err != nil {
		return err
	}

	// row
	err = db.Callback().Row().Before("gorm:row").Register("debugGormPluginBeforeRow", p.before)
	if err != nil {
		return err
	}

	err = db.Callback().Row().After("gorm:row").Register("debugGormPluginAfterRow", p.after)
	if err != nil {
		return err
	}

	// raw
	err = db.Callback().Raw().Before("gorm:raw").Register("debugGormPluginBeforeRaw", p.before)
	if err != nil {
		return err
	}

	err = db.Callback().Raw().After("gorm:raw").Register("debugGormPluginAfterRaw", p.after)
	if err != nil {
		return err
	}

	return nil
}

func (p GormPlugin) before(db *gorm.DB) {
	// make sure context could be used
	if db == nil {
		return
	}

	if db.Statement == nil || db.Statement.Context == nil {
		return
	}

	//requestID := GetTracingIDFromContext(db.Statement.Context)
	db.InstanceSet("debugGormPlugin:startTime", time.Now())
	//db.InstanceSet("debugGormPlugin:tracingID", requestID)
}

func (p GormPlugin) after(db *gorm.DB) {
	startTime, _ := db.InstanceGet("debugGormPlugin:startTime")
	tracingID, _ := db.InstanceGet("debugGormPlugin:tracingID")
	if t, ok := startTime.(time.Time); ok && tracingID != nil {
		log.Infow(
			fmt.Sprintf("Gorm (tracingID %s)", tracingID),
			"error", db.Error,
			"ex_qstr", db.Statement.SQL.String(),
			"ex_qpts", time.Since(t),
			"ex_qtable", db.Statement.Table,
			"context",
			map[string]interface{}{
				"payload": pkg.ToJSONString(map[string]interface{}{
					"vars": db.Statement.Vars,
				}),
			},
		)
	}
}
