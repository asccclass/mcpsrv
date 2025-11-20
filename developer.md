# MCPsrv 開發手冊

本手冊旨在協助開發者快速上手、開發與維護 MCPsrv 專案。

## 1. 專案簡介

MCPsrv 是一個基於 Go 語言開發的輕量級 MCP (Model Context Protocol) Server。它提供了一個靈活的框架來處理 MCP 指令並管理連線。

**核心技術：**
*   **語言**：Go (Golang)
*   **Web 框架**：`sherryserver` (基於標準庫 `net/http` 的封裝)
*   **MCP 實作**：`mcp-go`
*   **依賴管理**：Go Modules

## 2. 環境需求

在開始開發之前，請確保您的開發環境已安裝以下工具：

*   **Go**: 版本 1.24 或更高。
*   **Docker** (選用): 用於容器化部署與測試。
*   **Make**: 用於執行自動化指令 (Windows 使用者可安裝 MinGW 或使用 Git Bash)。

## 3. 快速開始 (Quick Start)

### 3.1 取得專案

```bash
git clone https://github.com/asccclass/mcpsrv.git
cd mcpsrv
```

### 3.2 設定環境變數

複製範例設定檔並依需求修改：

```bash
cp envfile.example envfile
```

`envfile` 重要設定說明：
*   `PORT`: 伺服器監聽埠號 (預設 8080 或 11042)
*   `DocumentRoot`: 靜態網頁根目錄 (預設 `www/html`)
*   `TemplateRoot`: 模板檔案目錄 (預設 `www/template`)
*   `ToDoAPIEndpoint`: TODO 工具的 API 端點
*   `weatherAPIEndpoint`: 天氣工具的 API 端點
*   `SystemName`: 系統名稱 (用於 Health Check)

### 3.3 初始化與編譯

使用 `make` 指令進行初始化與編譯：

```bash
# 下載 Go 模組依賴
make init

# 編譯專案 (產出執行檔 app)
make build
```

### 3.4 啟動服務

```bash
# 直接執行編譯後的程式
./app

# 或使用 Docker 啟動 (會自動編譯並打包)
make run
```

服務啟動後，預設可透過 `http://localhost:11042` (或您設定的 PORT) 存取。

## 4. 專案結構說明

```
mcpsrv/
├── libs/               # 共用函式庫 (Library)
│   └── mcpserver/      # MCP Server 核心封裝
├── tools/              # MCP 工具實作 (Tools)
│   ├── todo/           # 待辦事項工具
│   └── weather/        # 天氣查詢工具
├── www/                # 靜態資源目錄
│   └── html/           # 網頁檔案 (index.html 等)
├── auth.go             # 認證相關邏輯
├── router.go           # 路由定義與工具註冊
├── server.go           # 程式進入點 (Main)
├── makefile            # 自動化指令腳本
├── go.mod              # Go 模組定義
└── README.md           # 專案說明文件
```

## 5. 開發指南

### 5.1 新增 MCP 工具 (Add New Tool)

若要新增一個 MCP 工具，請依照以下步驟：

1.  **建立工具目錄**：在 `tools/` 下建立新的目錄，例如 `tools/mytool/`。
2.  **實作工具邏輯**：在該目錄下撰寫 Go 程式碼，實作您的工具功能。通常需要實作一個結構體並提供 `AddTools(s *mcpserver.MCPServer)` 方法。
3.  **註冊工具**：開啟 `router.go`，在 `NewRouter` 函式中註冊您的新工具。

**範例 (router.go):**

```go
// ... import your tool package

// 註冊 MyTool
myToolSrv, err := mytool.NewMyToolServer(os.Getenv("MyToolEndpoint"))
if err == nil {
    myToolSrv.AddTools(mcpserver)
} else {
    fmt.Println("service mytool is not initial:", err.Error())
}
```

### 5.1.1 AddTools 架構詳解

`AddTools(s *SryMCPServer.MCPServer)` 是每個 MCP 工具的核心註冊入口。開發者需在此函式中完成以下五大要素的定義與註冊：

#### 1. 定義意圖判斷 Prompt (`IsRelatedPrompt`)
這是一段給 LLM 的指令，用於判斷使用者的輸入是否與此工具有關。
*   **判斷標準**：必須明確定義何時回答「是」或「否」。
*   **回傳格式**：需定義回傳的 JSON 格式，通常包含 `is_related` (bool), `action` (string), `parameters` (object)。

#### 2. 定義 Host Capabilities (`HostCapabilities`)
描述此工具集提供的功能列表。
*   `Version`: 版本號。
*   `ServerID`: 伺服器識別碼。
*   `Tools`: `HostTool` 的切片 (Slice)，列出所有可用工具的描述。

#### 3. 定義 MCP Host (`MCPHost`)
代表此工具服務的主體物件。
*   `ID`: 服務 ID (如 "Weather001")。
*   `Name`: 服務名稱 (如 "WEATHER")。
*   `IsRelatedPrompt`: 填入上述定義的 Prompt。

#### 4. 註冊個別工具 (`HostTool` & `s.AddTool`)
對於每一個功能 (Function)，需要做兩件事：
1.  **描述工具 (HostTool)**：加入到 `capz.Tools` 中。這是給 LLM 看的，包含名稱、描述與參數定義。
2.  **註冊實作 (s.AddTool)**：這是給 Server 執行的。包含：
    *   `Name`: 工具名稱 (需與 HostTool 一致)。
    *   `Handler`: 實際執行的 Go 函式。
    *   `InputSchema`: 參數的驗證規則 (JSON Schema)。

#### 5. 掛載至 Server (`s.Tools`)
最後將設定好的 `mcphost` 物件指定給 `s.Tools` Map，Key 為工具模組名稱 (如 "weather")。

**程式碼結構範例：**

```go
func(app *MyToolServer) AddTools(s *SryMCPServer.MCPServer) {
    // 1. 定義 Prompt
    prompt := `...你的 Prompt...`

    // 2. 初始化 Capabilities
    capz := &SryMCPServer.HostCapabilities{
        Version: "1.0",
        ServerID: "0",
    }

    // 3. 初始化 Host
    mcphost := &SryMCPServer.MCPHost{
        ID: "MyTool001",
        Name: "MYTOOL",
        IsRelatedPrompt: prompt,
    }

    // 4. 定義並註冊工具 (例如: get_info)
    // 4.1 描述工具 (給 LLM)
    capz.Tools = append(capz.Tools, SryMCPServer.HostTool{
        Name: "get_info",
        Description: "取得資訊",
        Parameters: map[string]string{"id": "string"},
    })

    // 4.2 註冊實作 (給 Server)
    s.AddTool(SryMCPServer.Tool{
        Name: "get_info",
        Handler: app.getInfoHandler, // 實際執行的 Go func
        InputSchema: SryMCPServer.InputSchema{
             Type: "object",
             Properties: map[string]SryMCPServer.PropertySchema{
                 "id": {Type: "string", Description: "ID"},
             },
        },
    })

    // 5. 完成掛載
    mcphost.Capabilities = *capz
    s.Tools["mytool"] = mcphost
}
```

### 5.2 修改路由 (Routing)

路由定義位於 `router.go` 的 `NewRouter` 函式中。

*   **靜態檔案**：透過 `SherryServer.StaticFileServer` 處理。
*   **MCP 服務**：透過 `mcpserver.AddRouter(router)` 掛載。
*   **API 路由**：使用 `router.HandleFunc` 新增自定義 API，例如 `/healthz`。

### 5.3 靜態網頁開發

靜態網頁檔案放置於 `www/html` 目錄下。預設首頁為 `index.html`。
若您修改了 `envfile` 中的 `DocumentRoot`，請將檔案放置於對應目錄。

## 6. 部署與維運

### 6.1 Docker 部署

本專案提供完整的 Docker 支援。

*   **建置映像檔**：`make docker`
*   **啟動容器**：`make run` (會掛載當前目錄的 `www`, `data`, `envfile` 到容器中)
*   **停止容器**：`make stop`
*   **查看 Log**：`make log`

### 6.2 Health Check

服務提供 `/healthz` 端點用於健康狀態檢查。
回應範例：`{"status": "ok", "system": "YourSystemName"}`

## 7. 常見問題 (FAQ)

*   **Q: 修改了程式碼但沒有生效？**
    *   A: 請確認是否重新執行了 `make build` 或重啟了 Docker 容器 (`make re`)。
*   **Q: 缺少依賴套件？**
    *   A: 請執行 `make init` 或 `go mod tidy`。

---
*文件最後更新日期：2025-11-21*
