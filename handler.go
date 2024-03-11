package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/microsoft/go-mssqldb"
)

type crewMemberInfo struct {
	Crew_id string         `json:"crew_id"`
	Data    crewMemberData `json:"data"`
}

type crewMemberData struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	Rank string `json:"rank"`
}

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type entireCompliment []crewMemberInfo

var db *sql.DB

func handleCrewMember(c *gin.Context) {
	var crewMember crewMemberInfo

	crewMember.Crew_id = c.Param("id")

	_, err := databaseConnection()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	defer closeDatabaseConnection()

	tsql := `
	SELECT *
	FROM [dbo].[EnterpriseComplement]
	WHERE crew_id = '` + crewMember.Crew_id + `';
	`

	// Execute query
	rows, err := db.Query(tsql)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	err = crewMember.checkCrewMember(rows)
	defer rows.Close()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"crew_member": crewMember,
	})
}

func handleEntireCrew(c *gin.Context) {
	var allCrew entireCompliment

	_, err := databaseConnection()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	defer closeDatabaseConnection()

	tsql := `
		SELECT *
		FROM [dbo].[EnterpriseComplement];
		`

	rows, err := db.Query(tsql)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	err = allCrew.checkAllCrew(rows)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, allCrew)
}

func handleCount(c *gin.Context) {

	var count int

	_, err := databaseConnection()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	defer closeDatabaseConnection()

	tsql := `
		SELECT COUNT(*)	
		FROM [dbo].[EnterpriseComplement];
		`

	// Execute query
	rows, err := db.Query(tsql)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	count, err = checkCount(rows)
	defer rows.Close()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"crew_count": count,
	})
}

func databaseConnection() (context.Context, error) {

	config, err := readConfigFile("config.json")
	if err != nil {
		panic(err)
	}

	var server = "sql-database-api.database.windows.net"
	var port = 1433
	var user = config.Username
	var password = config.Password
	var database = "SQLDatabaseAPI"

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %v", err.Error())
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify connect: %v", err.Error())
	}

	return ctx, nil
}

func closeDatabaseConnection() {
	if db != nil {
		db.Close()
	}
}

func readConfigFile(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func checkCount(rows *sql.Rows) (int, error) {
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return -1, err
		}
	}
	return count, nil
}

func (allCrew *entireCompliment) checkAllCrew(rows *sql.Rows) error {

	for rows.Next() {
		var crewID int
		var crewJSON string
		var dateInserted string

		// Scan the row into variables
		if err := rows.Scan(&crewID, &crewJSON, &dateInserted); err != nil {
			return err
		}

		// Unmarshal the crew info JSON string into a crewMemberData struct
		var crewData crewMemberData
		if err := json.Unmarshal([]byte(crewJSON), &crewData); err != nil {
			return err
		}

		// Create a crewMemberInfo object
		crewMemberInfo := crewMemberInfo{
			Crew_id: strconv.Itoa(crewID),
			Data:    crewData,
		}

		// Append crewMemberInfo to allCrew slice
		*allCrew = append(*allCrew, crewMemberInfo)
	}
	return nil
}

func (crewMember *crewMemberInfo) checkCrewMember(rows *sql.Rows) error {

	// Loop through rows
	for rows.Next() {
		var crewID int
		var crewJSON string
		var dateInserted string

		// Scan the row into variables
		if err := rows.Scan(&crewID, &crewJSON, &dateInserted); err != nil {
			return err
		}

		// Unmarshal the crew info JSON string into a crewMemberData struct
		var crewData crewMemberData
		if err := json.Unmarshal([]byte(crewJSON), &crewData); err != nil {
			return err
		}

		// Create a crewMemberInfo object
		crewMember.Data = crewData
	}
	return nil
}
