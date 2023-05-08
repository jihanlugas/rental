package cmd

import (
	"github.com/jihanlugas/rental/cryption"
	"github.com/jihanlugas/rental/db"
	"github.com/jihanlugas/rental/model"
	"github.com/jihanlugas/rental/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"time"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Run server",
	Long: `With this command you can
	up : create database table
	down :  drop database table
	seed :  insert data table
	`,
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Up table",
	Long:  "Up table",
	Run: func(cmd *cobra.Command, args []string) {
		up()
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Down table",
	Long:  "Down table",
	Run: func(cmd *cobra.Command, args []string) {
		down()
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed data table",
	Long:  "Seed data table",
	Run: func(cmd *cobra.Command, args []string) {
		seed()
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Down, up, seed table",
	Long:  "Down, up, seed table",
	Run: func(cmd *cobra.Command, args []string) {
		down()
		up()
		seed()
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(upCmd)
	dbCmd.AddCommand(downCmd)
	dbCmd.AddCommand(resetCmd)
	dbCmd.AddCommand(seedCmd)
}

func up() {
	var err error
	conn, closeConn := db.GetConnection()
	defer closeConn()

	// table
	err = conn.Migrator().AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Property{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Propertysetting{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Calendar{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Userproperties{})
	if err != nil {
		panic(err)
	}

	// view
	vUser := conn.Model(&model.User{}).
		Where("delete_dt is null")

	vProperty := conn.Model(&model.Property{}).
		Joins("join users on users.id = properties.user_id ").
		Where("properties.delete_dt is null")

	vPropertysetting := conn.Model(&model.Propertysetting{}).
		Joins("join properties on properties.id = propertysettings.id ").
		Where("properties.delete_dt is null")

	vUserproperty := conn.Model(&model.Userproperties{}).
		Joins("join users on users.id = userproperties.user_id").
		Joins("join properties on properties.id = userproperties.property_id ").
		Where("userproperties.delete_dt is null")

	err = conn.Migrator().CreateView("users_view", gorm.ViewOption{
		Replace: true,
		Query:   vUser,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("properties_view", gorm.ViewOption{
		Replace: true,
		Query:   vProperty,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("propertysettings_view", gorm.ViewOption{
		Replace: true,
		Query:   vPropertysetting,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("userproperties_view", gorm.ViewOption{
		Replace: true,
		Query:   vUserproperty,
	})
	if err != nil {
		panic(err)
	}
}

func down() {
	var err error
	conn, closeConn := db.GetConnection()
	defer closeConn()

	// view
	err = conn.Migrator().DropView("users_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("properties_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("propertysettings_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("userproperties_view")
	if err != nil {
		panic(err)
	}

	// table
	err = conn.Migrator().DropTable(&model.User{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Property{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Propertysetting{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Userproperties{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Calendar{})
	if err != nil {
		panic(err)
	}

}

func seed() {
	now := time.Now()
	password, err := cryption.EncryptAES64("123456")
	if err != nil {
		panic(err)
	}
	users := []model.User{
		{ID: utils.GetUniqueID(), RoleID: utils.GetUniqueID(), Email: "jihanlugas2@gmail.com", Username: "jihanlugas", NoHp: "6287770333043", Fullname: "Jihan Lugas", Passwd: password, PassVersion: 1, Active: true, LastLoginDt: nil, PhotoID: "", CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}

	properties := []model.Property{
		{ID: utils.GetUniqueID(), UserID: users[0].ID, Name: "Property 1", CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}

	propertysettings := []model.Propertysetting{
		{ID: properties[0].ID, DefaultTimeStart: 12, DefaultTimeEnd: 12},
	}

	events := []model.Calendar{
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "1", Name: "Tes", StartDt: now.Add(-5 * time.Hour), EndDt: now.Add(-4 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "2", Name: "Tes", StartDt: now.Add(-4 * time.Hour), EndDt: now.Add(-3 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "1", Name: "Tes", StartDt: now.Add(-3 * time.Hour), EndDt: now.Add(-2 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "2", Name: "Tes", StartDt: now.Add(-2 * time.Hour), EndDt: now.Add(-1 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "1", Name: "Tes", StartDt: now.Add(-1 * time.Hour), EndDt: now.Add(0 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "2", Name: "Tes", StartDt: now.Add(0 * time.Hour), EndDt: now.Add(1 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "1", Name: "Tes", StartDt: now.Add(1 * time.Hour), EndDt: now.Add(2 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "2", Name: "Tes", StartDt: now.Add(2 * time.Hour), EndDt: now.Add(3 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "1", Name: "Tes", StartDt: now.Add(3 * time.Hour), EndDt: now.Add(4 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), PropertyID: properties[0].ID, GroupID: "2", Name: "Tes", StartDt: now.Add(4 * time.Hour), EndDt: now.Add(5 * time.Hour), CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}

	userproperties := []model.Userproperties{
		{ID: utils.GetUniqueID(), UserID: users[0].ID, PropertyID: properties[0].ID, DefaultProperty: true, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	tx.Create(&users)
	tx.Create(&events)
	tx.Create(&properties)
	tx.Create(&propertysettings)
	tx.Create(&userproperties)

	err = tx.Commit().Error
	if err != nil {
		panic(err)
	}
}
