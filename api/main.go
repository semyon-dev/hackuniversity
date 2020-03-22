package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/semyon-dev/hackuniversity/api/db"
	"github.com/semyon-dev/hackuniversity/api/model"
	"math"
	"strconv"
	"strings"
	"time"
)

func main() {

	r := gin.Default()
	r.Use(cors.Default())

	db.ConnectPostgres()
	db.ConnectClickhouse()

	r.GET("/criticals", func(context *gin.Context) {

		criticals := make(map[string]map[string]float64)
		for _, i := range db.GetCriticals() {
			criticals[i.Name] = map[string]float64{"min": i.Min, "max": i.Max}
		}

		context.JSON(200, criticals)
	})

	r.POST("/critical", func(context *gin.Context) {

		var critical model.Criticals
		err := context.ShouldBindJSON(&critical)
		if err != nil {
			fmt.Println(err)
			context.JSON(400, gin.H{
				"status":  "ERROR",
				"message": err,
			})
		}
		err = db.UpdateCritical(critical.Name, critical.Min, critical.Max)
		var status int
		var message string
		if err != nil {
			message = "ERROR"
			status = 500
			fmt.Println(err)
		} else {
			message = "OK"
			status = 200
		}
		// TODO: author name and min, max values
		db.NewEvent(critical.Name, critical.Author)
		context.JSON(status, gin.H{
			"message": message,
		})
	})

	// test url: /average?paramName=HUMIDITY&dateStart=2020-03-20&dateEnd=2020-03-30&timeStart=00:00:00&timeEnd=00:00:00
	// return average value between start date and time and end date and time
	r.GET("/average", func(context *gin.Context) {
		name, dateTimeStart, dateTimeEnd := nameDateTimes(context)
		fmt.Println(name, dateTimeStart, dateTimeEnd, " values from query")

		params := db.AverageValue(name, dateTimeStart, dateTimeEnd)
		context.JSON(200,
			gin.H{
				"parameters": params,
			})
	})

	r.GET("/max", func(context *gin.Context) {
		name, dateTimeStart, dateTimeEnd := nameDateTimes(context)
		fmt.Println(name, dateTimeStart, dateTimeEnd, " values from query")

		params := db.MaxValue(name, dateTimeStart, dateTimeEnd)
		context.JSON(200,
			gin.H{
				"parameters": params,
			})
	})

	r.GET("/min", func(context *gin.Context) {
		name, dateTimeStart, dateTimeEnd := nameDateTimes(context)
		fmt.Println(name, dateTimeStart, dateTimeEnd, " values from query")

		params := db.MinValue(name, dateTimeStart, dateTimeEnd)
		context.JSON(200,
			gin.H{
				"parameters": params,
			})
	})

	r.GET("/maindata", func(context *gin.Context) {
		name, dateTimeStart, dateTimeEnd := nameDateTimes(context)

		min := db.MinValue(name, dateTimeStart, dateTimeEnd)
		max := db.MaxValue(name, dateTimeStart, dateTimeEnd)
		avg := db.AverageValue(name, dateTimeStart, dateTimeEnd)

		if math.IsNaN(avg) {
			avg = 0
		}

		context.JSON(200,
			gin.H{
				"min": min,
				"avg": avg,
				"max": max,
			})
	})

	// old version
	r.GET("/hourly", func(context *gin.Context) {
		param := context.Query("param")
		_, dateTimeStart, dateTimeEnd := nameDateTimes(context)
		execStr := "SELECT " + param + " FROM journal WHERE action_time BETWEEN toDateTime('" + dateTimeStart + "', 'Europe/Moscow')  AND toDateTime('" + dateTimeEnd + "', 'Europe/Moscow')"
		rows, err := db.Clicconn.Query(execStr)
		if err != nil {
			fmt.Println(err)
		}

		var val float64
		var res = make([]float64, 0)
		var hours = make([]float64, 86400)
		temp := 3600
		for i := 0; rows.Next(); i++ {
			err = rows.Scan(&val)
			if err != nil {
				fmt.Println(err)
			}
			hours[i] = val
			if i == temp {
				summa := 0.0
				for t := temp - 3600; t <= temp; t++ {
					summa += hours[t]
				}
				res = append(res, summa/3600)
				temp += 3600
			}
		}

		max := 24 - len(res)
		for i := len(res); i <= max; i++ {
			res = append(res, float64(i))
		}





		context.JSON(200, gin.H{"data": res})
	})

	r.GET("/hourly2", func(context *gin.Context) {

		today:= strings.Split(time.Now().String()," ")[0]
		fmt.Println(today +" - is today!!!!!!!!!!!!!!!")
		param := context.Query("param")
		date:= context.DefaultQuery("date",today)


		values:=db.GetHourlyErrors(param,date)
		fmt.Println(values)

		context.JSON(200,gin.H{
			"data":values,
		})
	})

	r.GET("/errors", func(context *gin.Context) {

		dateStart := context.DefaultQuery("dateStart","today")
		dateEnd := context.DefaultQuery("dateEnd","now")
		limit := context.DefaultQuery("limit","10")
		limitInt,_:=strconv.Atoi(limit)

		periodErrors := db.GetErrors(dateStart, dateEnd,limitInt)

		context.JSON(200, gin.H{
			"errors": periodErrors,
		})
	})

	fmt.Println("запуск API на :5000...")
	err := r.Run(":5000")
	if err != nil {
		fmt.Println("ошибка при запуске API:", err)
	}
}

// получение границ даты и времени из URL
func nameDateTimes(context *gin.Context) (string, string, string) {
	name := context.Query("paramName")
	dateStart := context.Query("dateStart")
	var dateTimeStart, dateTimeEnd string
	if dateStart == "today" || len(dateStart) == 0 {
		currentTime := time.Now().String()
		strCurrTime := strings.Split(currentTime, ".")[0]
		dateTimeStart = strings.Split(strCurrTime, " ")[0] + " 00:00:00"
		dateTimeEnd = strCurrTime
	} else {
		dateEnd := context.Query("dateEnd")
		timeStart := context.Query("timeStart")
		timeEnd := context.Query("timeEnd")
		if len(timeEnd) == 0 {
			timeEnd = "00:00:00"
		}
		if len(timeStart) == 0 {
			timeStart = "00:00:00"
		}

		dateTimeStart = dateStart + " " + timeStart
		dateTimeEnd = dateEnd + " " + timeEnd
	}

	return name, dateTimeStart, dateTimeEnd
}
