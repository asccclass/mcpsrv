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

// WeatherEntry holds a single weather record
// 天氣資料結構
type WeatherEntry struct {
    StartTime                    string `json:"StartTime,omitempty"`
    EndTime                      string `json:"EndTime,omitempty"`
    AvgTemp                      string `json:"AvgTemp,omitempty"` // 注意：即使是空字串，也用 string
    MaxTemp                      string `json:"MaxTemp,omitempty"`
    MinTemp                      string `json:"MinTemp,omitempty"`
    RelativeHumidity             string `json:"RelativeHumidity,omitempty"`
    MaxApparentTemperature       string `json:"MaxApparentTemperature,omitempty"`
    MinApparentTemperature       string `json:"MinApparentTemperature,omitempty"`
    MaxComfortIndexDescription   string `json:"MaxComfortIndexDescription,omitempty"`
    ProbabilityOfPrecipitation   string `json:"ProbabilityOfPrecipitation,omitempty"`
    WeatherDesc                  string `json:"WeatherDesc,omitempty"`
    UVExposureLevel              string `json:"UVExposureLevel,omitempty"`
    // 這裡要注意 Key 之間有空格，在 Struct 欄位上要用駝峰式命名，並在 Tag 中完整保留
    WeatherDescription           string `json:"Weather Description,omitempty"`
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
   Weathers []WeatherEntry `json:"weathers"`
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
func(app *WeatherMCPServer) judgeTime(data []WeatherEntry)(string) {
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
   city, ok := args["city"].(string)
   if !ok {
      return nil, fmt.Errorf("Missing required parameter: city")
   }
   url := fmt.Sprintf("%s/status/%s", app.API, city)
   resp, err := app.makeAPIRequest("GET", url, nil)
   if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to get todos: %v", err)), nil
   }
   var statuz WeatherEntry
   if err := json.Unmarshal([]byte(resp), &statuz); err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to parse response: %v", err)), nil
   }
   desc := ""
   if statuz.StartTime == "" {
      desc = "無" + city + "的氣候資料"
      return mcp.NewToolResultText(desc), nil
   }
   desc = fmt.Sprintf("%s目前為%s", city, statuz.WeatherDescription)
   if statuz.UVExposureLevel != "" {
      desc = desc + fmt.Sprintf("紫外線等級為%s", statuz.UVExposureLevel)
   }
   return mcp.NewToolResultText(desc), nil
}
