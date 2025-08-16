package SherryWeather

import(
   "context"
   "github.com/mark3labs/mcp-go/mcp"
)

// MCP Interface
func(app *Weather) ReadWeather(ctx context.Context, request mcp.CallToolRequest)(*mcp.CallToolResult, error) {
   region, ok := request.Params.Arguments["region"].(string)
   if !ok {
      return mcp.NewToolResultError("name must be a string"), nil
   }
   if len(app.WeatherData) == 0 {
      if err := app.ReadWeatherJson(); err != nil {
         return mcp.NewToolResultError(err.Error()), nil
      }
   }
   data := app.SearchStatus(region)
   if data == "" {
      return mcp.NewToolResultError(region + " not found"), nil
   }
   return mcp.NewToolResultText(data), nil
}

type Weather struct {
   Name		string					// 工具名稱
   Description	string					// 工具描述
   Params	string					// 參數名稱
   ParamsDesc	string					// 參數描述
   WeatherData	[]WeatherData				// 天氣資料
   MCPHandler	func(ctx context.Context, request mcp.CallToolRequest)(*mcp.CallToolResult, error)
}

func(app *Weather) AddTools()(mcp.Tool) {
   tool := mcp.NewTool(app.Name,
           mcp.WithDescription(app.Description),
           mcp.WithString(app.Params, mcp.Required(), mcp.Description(app.ParamsDesc)),
       )
   return tool
}

func NewWeather()(*Weather) {
   w := Weather {
      Name: "weather_status",
      Description: "Get the weekly weather forecast by region",
      Params: "regign",
      ParamsDesc: "Region name to get weather conditions for",
   }
   w.MCPHandler = w.ReadWeather
   return &w
}
