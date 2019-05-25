package cmd

import (
	"database/sql"
	"fmt"
	"go-api-base/app/user"
	"go-api-base/mysql"
	"go-api-base/router"

	// Load MySQL driver
	_ "github.com/go-sql-driver/mysql"
	drvMySql "github.com/go-sql-driver/mysql"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("JWT_SECRET_KEY") == "" {
			panic("JWT_SECRET_KEY in .env is empty")
		}
		jwtSecretKey := []byte(viper.GetString("JWT_SECRET_KEY"))

		if debugMode := viper.GetBool("DEBUG_MODE"); !debugMode {
			gin.SetMode(gin.ReleaseMode)
		}

		dataSourceName := drvMySql.NewConfig()
		dataSourceName.Net = viper.GetString("DB_NET")
		dataSourceName.Addr = viper.GetString("DB_ADDRESS")
		dataSourceName.User = viper.GetString("DB_USERNAME")
		dataSourceName.Passwd = viper.GetString("DB_PASSWORD")
		dataSourceName.DBName = viper.GetString("DB_DATABASE")

		db, err := sql.Open("mysql", dataSourceName.FormatDSN())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			panic(err)
		}

		userRepository := mysql.NewUserRepository(db)

		routerConfig := &router.Config{
			UserService: user.NewService(
				jwtSecretKey,
				userRepository,
			),
		}

		fmt.Println("App is runnning")
		fmt.Println("Ctrl + C to terminate")

		app := router.New(routerConfig)

		port := viper.GetString("PORT")

		app.Run(":" + port)
	},
}
