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
	err = conn.Migrator().AutoMigrate(&model.Company{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Companysetting{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Usercompany{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Property{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Calendar{})
	if err != nil {
		panic(err)
	}

	// view
	vUser := conn.Model(&model.User{}).
		Where("delete_dt is null")

	vCompany := conn.Model(&model.Company{}).
		Joins("join users on users.id = companies.user_id ").
		Where("companies.delete_dt is null")

	vCompanysetting := conn.Model(&model.Companysetting{}).
		Joins("join companies on companies.id = companysettings.id ").
		Where("companies.delete_dt is null")

	vUsercompany := conn.Model(&model.Usercompany{}).
		Joins("join users on users.id = usercompanies.user_id").
		Joins("join companies on companies.id = usercompanies.company_id ").
		Where("usercompanies.delete_dt is null")

	vProperty := conn.Model(&model.Property{}).
		Joins("join companies on companies.id = properties.company_id ").
		Where("companies.delete_dt is null")

	vCalendar := conn.Model(&model.Calendar{}).
		Select(
			"calendars.id",
			"calendars.company_id",
			"calendars.property_id",
			"calendars.name",
			"calendars.start_dt",
			"calendars.end_dt",
			"calendars.status",
			"calendars.create_by",
			"calendars.create_dt",
			"calendars.update_by",
			"calendars.update_dt",
		).
		Where("delete_dt is null")

	err = conn.Migrator().CreateView("users_view", gorm.ViewOption{
		Replace: true,
		Query:   vUser,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("companies_view", gorm.ViewOption{
		Replace: true,
		Query:   vCompany,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("companysettings_view", gorm.ViewOption{
		Replace: true,
		Query:   vCompanysetting,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("usercompanies_view", gorm.ViewOption{
		Replace: true,
		Query:   vUsercompany,
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

	err = conn.Migrator().CreateView("calendars_view", gorm.ViewOption{
		Replace: true,
		Query:   vCalendar,
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
	err = conn.Migrator().DropView("companies_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("companysettings_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("usercompanies_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("properties_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("calendars_view")
	if err != nil {
		panic(err)
	}

	// table
	err = conn.Migrator().DropTable(&model.User{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Company{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Companysetting{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Usercompany{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Property{})
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

	companies := []model.Company{
		{ID: utils.GetUniqueID(), UserID: users[0].ID, Name: "Company 1", CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}

	companysettings := []model.Companysetting{
		{ID: companies[0].ID, DefaultTimeStart: 12, DefaultTimeEnd: 12},
	}

	properties := []model.Property{
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, Name: "Lapangan 1", Description: "Description"},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, Name: "Lapangan 2", Description: "Description"},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, Name: "Lapangan 3", Description: "Description"},
	}

	calendars := []model.Calendar{
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[0].ID, Name: "Tes", StartDt: now.Add(-5 * time.Hour), EndDt: now.Add(-4 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[1].ID, Name: "Tes", StartDt: now.Add(-4 * time.Hour), EndDt: now.Add(-3 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[2].ID, Name: "Tes", StartDt: now.Add(-3 * time.Hour), EndDt: now.Add(-2 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[0].ID, Name: "Tes", StartDt: now.Add(-2 * time.Hour), EndDt: now.Add(-1 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[1].ID, Name: "Tes", StartDt: now.Add(-1 * time.Hour), EndDt: now.Add(0 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[2].ID, Name: "Tes", StartDt: now.Add(0 * time.Hour), EndDt: now.Add(1 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[0].ID, Name: "Tes", StartDt: now.Add(1 * time.Hour), EndDt: now.Add(2 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[1].ID, Name: "Tes", StartDt: now.Add(2 * time.Hour), EndDt: now.Add(3 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[2].ID, Name: "Tes", StartDt: now.Add(3 * time.Hour), EndDt: now.Add(4 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[0].ID, Name: "Tes", StartDt: now.Add(4 * time.Hour), EndDt: now.Add(5 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[1].ID, Name: "Tes", StartDt: now.Add(6 * time.Hour), EndDt: now.Add(8 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[2].ID, Name: "Tes", StartDt: now.Add(8 * time.Hour), EndDt: now.Add(10 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[0].ID, Name: "Tes", StartDt: now.Add(10 * time.Hour), EndDt: now.Add(12 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
		{ID: utils.GetUniqueID(), CompanyID: companies[0].ID, PropertyID: properties[1].ID, Name: "Tes", StartDt: now.Add(12 * time.Hour), EndDt: now.Add(14 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}

	usercompanies := []model.Usercompany{
		{ID: utils.GetUniqueID(), UserID: users[0].ID, CompanyID: companies[0].ID, DefaultCompany: true, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	tx.Create(&users)
	tx.Create(&companies)
	tx.Create(&companysettings)
	tx.Create(&usercompanies)
	tx.Create(&properties)
	tx.Create(&calendars)

	err = tx.Commit().Error
	if err != nil {
		panic(err)
	}
}
