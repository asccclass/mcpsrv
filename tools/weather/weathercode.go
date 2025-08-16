package weatherMCPServer

import(
   "os"
   "fmt"
   "time"
   "io/ioutil"
   "encoding/json"
   "github.com/mark3labs/mcp-go/mcp"
)

type ElementValue struct {
   Temperature  string `json:"Temperature,omitempty"`
   MaxTemperature string `json:"MaxTemperature,omitempty"`
}

type Time struct {
   StartTime  string        `json:"StartTime"`
   EndTime    string        `json:"EndTime"`
   ElementValue []ElementValue `json:"ElementValue"`
}

type WeatherElement struct {
   ElementName string `json:"ElementName"`  //WeatherElement
   Time        []Time   `json:"Time"`
}

// WeatherRecord holds a single weather record
// 天氣資料結構
type WeatherRecord struct {
   StartTime                 string  `json:"StartTime"`
   EndTime                   string  `json:"EndTime"`
   AvgTemp                   string `json:"AvgTemp (°C)"`
   MaxTemp                   string `json:"MaxTemp (°C)"`
   MinTemp                   string `json:"MinTemp (°C)"`
   RelativeHumidity          string  `json:"RelativeHumidity"`
   MaxApparentTemperature    string `json:"MaxApparentTemperature"`
   MinApparentTemperature    string `json:"MinApparentTemperature"`
   MaxComfortIndexDescription string  `json:"MaxComfortIndexDescription"`
   ProbabilityOfPrecipitation string     `json:"ProbabilityOfPrecipitation"`
   WeatherDesc               string  `json:"WeatherDesc"`
   UVExposureLevel           string  `json:"UVExposureLevel"`
   WeatherDescription        string  `json:"Weather Description"`
}

// DailyForecast represents the weather forecast for a day
type DailyForecast struct {
   Date      string
   DayOfWeek string
   MaxTemp   float64
   MinTemp   float64
   Weather   string
   Humidity  int
   RainProb  int
}

// WeatherData holds the weather forecast data structure
// 天氣資料結構
type WeatherData struct {
   City     string          `json:"city"`
   Weathers []WeatherRecord `json:"weathers"`
}

// 判斷目前時間是否落在兩個時間字串中
func(app *WeatherMCPServer) IsNowBetween(startStr, endStr string) (bool, error) {
   layout := "2006-01-02 15:04:05"  // 定義時間格式
   loc, err := time.LoadLocation("Local")
   if err != nil {
      return false, fmt.Errorf("parse local Time failed: %s", err.Error())
   }
   startTime, err := time.ParseInLocation(layout, startStr, loc)
   if err != nil {
      return false, fmt.Errorf("parse StartTime failed: %s", err.Error())
   }
   endTime, err := time.ParseInLocation(layout, endStr, loc)
   if err != nil {
      return false, fmt.Errorf("parse EndTime failed: %s", err.Error())
   }
   now := time.Now()  // time.Now().UTC() // 獲取現在時間
   return (now.After(startTime) || now.Equal(startTime)) && (now.Before(endTime) || now.Equal(endTime)), nil
}

// 找到最接近時間的天候狀態
func(app *WeatherMCPServer) judgeTime(data []WeatherRecord)(string) {
   result := ""
   for _, w := range data {
      if in, _ := app.IsNowBetween(w.StartTime, w.EndTime); in {
fmt.Println("search ok.")
         result = w.WeatherDescription
         break
      }
   }
   return result
}

// 轉換日期格式
// 將時間轉為指定格式的函式
func(app *WeatherMCPServer) convertTimeFormat(input string) (string) {
   parsedTime, err := time.Parse(time.RFC3339, input) // 解析 ISO 8601 格式時間
   if err != nil {
      fmt.Println(err.Error())
      return input
   }
   return parsedTime.Format("2006-01-02 15:04:05") // 轉換為目標格式
}

// 搜尋區域名稱
func(app *WeatherMCPServer) SearchStatus(region string)(string) {
   for _, data := range app.WeatherData {
      if data.City == region {
         return app.judgeTime(data.Weathers)
	 /*
         jsonBytes, err := json.Marshal(data.Weathers)
	 if err != nil {
            fmt.Println(err.Error())
	    return "" 
	 }
         return string(jsonBytes)
	 */
      }
   }
   return ""
}

// 讀取氣候檔
func(app *WeatherMCPServer) ReadWeatherJson()(error) {
   filePath := os.Getenv("DataRoot") + "weather.json"
   file, err := os.Open(filePath)
   if err != nil {
      fmt.Println(err.Error())
      return err
   }
   defer file.Close()
   byteValue, err := ioutil.ReadAll(file)
   if err != nil {
      fmt.Println(err.Error())
      return err
   }
   var weatherData []WeatherData
   if err := json.Unmarshal(byteValue, &weatherData); err != nil {
      fmt.Println(err.Error())
      return err
   }
   app.WeatherData = weatherData
   return nil
}

// 取得目前天氣狀態
func(s *WeatherMCPServer) getWeather(args map[string]interface{}) (*mcp.CallToolResult, error) {
	/*
   url := fmt.Sprintf("%s/todos", s.API)
   resp, err := s.makeAPIRequest("GET", url, nil)
   if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to get todos: %v", err)), nil
   }

   var todos []Todo
   if err := json.Unmarshal(resp, &todos); err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to parse response: %v", err)), nil
   }
   // 格式化輸出
   result := "# 現有待辦事項如下:\n"
   for _, todo := range todos {
      status := "❌ 未完成"
      if todo.IsFinish == "1" {
         status = "處理中"
      } else if todo.IsFinish == "2" {
         status = "✅ 完成"
      } else if todo.IsFinish == "3" {
         status = "擱置"
      }
      result += fmt.Sprintf("\n## ID: %d\n* 負責者: %s\n* 事件狀態: %s\n* 到期時間: %s", todo.ID, todo.User, status, todo.DueTime)
      result += fmt.Sprintf("\n* 內容: %s\n* 建立時間:%s）", todo.Context, todo.CreateDate.Format("2006-01-02 15:04:05"))
   }
   */
   result := "請直接查詢區域"
   return mcp.NewToolResultText(result), nil
}

// 查詢特定城市的天氣狀態
func(app *WeatherMCPServer) getWeatherByCity(args map[string]interface{}) (*mcp.CallToolResult, error) {
   city, ok := args["city"]
   if !ok {
      return nil, fmt.Errorf("Missing required parameter: city")
   }
   desc := app.SearchStatus(city.(string))
   if desc == "" {
      desc = "無" + city.(string) + "的氣候資料"
   }
   fmt.Println("get weather by city:", desc)
   return mcp.NewToolResultText(desc), nil
}
