package cmd

import (
	"fmt"
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
	err = conn.Migrator().AutoMigrate(&model.Item{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Calendar{})
	if err != nil {
		panic(err)
	}

	// view
	vUser := conn.Model(&model.User{}).
		Select("users.*, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("left join users u1 on u1.id = users.create_by").
		Joins("left join users u2 on u2.id = users.update_by").
		Joins("left join users u3 on u3.id = users.delete_by")

	vCompany := conn.Model(&model.Company{}).
		Select("companies.*, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("join users on users.id = companies.user_id ").
		Joins("left join users u1 on u1.id = companies.create_by").
		Joins("left join users u2 on u2.id = companies.update_by").
		Joins("left join users u3 on u3.id = companies.delete_by")

	vCompanysetting := conn.Model(&model.Companysetting{}).
		Select("companysettings.*").
		Joins("join companies on companies.id = companysettings.id ")

	vUsercompany := conn.Model(&model.Usercompany{}).
		Select("usercompanies.*, users.fullname, companies.name, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("join users on users.id = usercompanies.user_id").
		Joins("join companies on companies.id = usercompanies.company_id ").
		Joins("left join users u1 on u1.id = usercompanies.create_by").
		Joins("left join users u2 on u2.id = usercompanies.update_by").
		Joins("left join users u3 on u3.id = usercompanies.delete_by")

	vProperty := conn.Model(&model.Property{}).
		Select("properties.*, companies.name as company_name, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("join companies on companies.id = properties.company_id ").
		Joins("left join users u1 on u1.id = properties.create_by").
		Joins("left join users u2 on u2.id = properties.update_by").
		Joins("left join users u3 on u3.id = properties.delete_by")

	vItem := conn.Model(&model.Item{}).
		Select("items.*, companies.name as company_name, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("join companies on companies.id = items.company_id ").
		Joins("left join users u1 on u1.id = items.create_by").
		Joins("left join users u2 on u2.id = items.update_by").
		Joins("left join users u3 on u3.id = items.delete_by")

	vCalendar := conn.Model(&model.Calendar{}).
		Select("calendars.*, properties.name as property_name, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("join properties on properties.id = calendars.property_id").
		Joins("left join users u1 on u1.id = calendars.create_by").
		Joins("left join users u2 on u2.id = calendars.update_by").
		Joins("left join users u3 on u3.id = calendars.delete_by")

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

	err = conn.Migrator().CreateView("items_view", gorm.ViewOption{
		Replace: true,
		Query:   vItem,
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
	err = conn.Migrator().DropView("items_view")
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
	err = conn.Migrator().DropTable(&model.Item{})
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

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	users := []model.User{
		{RoleID: utils.GetUniqueID(), Email: "jihanlugas2@gmail.com", Username: "jihanlugas", NoHp: "6287770333043", Fullname: "Jihan Lugas", Passwd: password, PassVersion: 1, Active: true, LastLoginDt: nil, PhotoID: "", CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}

	tx.Create(&users)

	companies := []model.Company{
		{UserID: users[0].ID, Name: "Company 1", CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}
	tx.Create(&companies)

	companysettings := []model.Companysetting{
		{ID: companies[0].ID, DefaultTimeStart: 12, DefaultTimeEnd: 12},
	}
	tx.Create(&companysettings)

	properties := []model.Property{
		{CompanyID: companies[0].ID, Name: "Lapangan 1", Description: "Description"},
		{CompanyID: companies[0].ID, Name: "Lapangan 2", Description: "Description"},
		{CompanyID: companies[0].ID, Name: "Lapangan 3", Description: "Description"},
	}
	tx.Create(&properties)

	calendars := []model.Calendar{}
	startDedault := now.Add(-96 * time.Hour).Truncate(60 * time.Minute)
	for i := 0; i < 20; i++ {
		new1 := model.Calendar{CompanyID: companies[0].ID, PropertyID: properties[0].ID, Name: fmt.Sprintf("Tes data %d", (i*3)+1), StartDt: startDedault.Add(2 * time.Hour), EndDt: startDedault.Add(4 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now}
		new2 := model.Calendar{CompanyID: companies[0].ID, PropertyID: properties[1].ID, Name: fmt.Sprintf("Tes data %d", (i*3)+2), StartDt: startDedault.Add(6 * time.Hour), EndDt: startDedault.Add(8 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now}
		new3 := model.Calendar{CompanyID: companies[0].ID, PropertyID: properties[2].ID, Name: fmt.Sprintf("Tes data %d", (i*3)+3), StartDt: startDedault.Add(10 * time.Hour), EndDt: startDedault.Add(12 * time.Hour), Status: 1, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now}
		calendars = append(calendars, new1, new2, new3)
		startDedault = startDedault.Add(12 * time.Hour)
	}
	tx.Create(&calendars)

	usercompanies := []model.Usercompany{
		{UserID: users[0].ID, CompanyID: companies[0].ID, DefaultCompany: true, CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}
	tx.Create(&usercompanies)

	err = tx.Commit().Error
	if err != nil {
		panic(err)
	}
}
