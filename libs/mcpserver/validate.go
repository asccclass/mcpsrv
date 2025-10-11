package SryMCPServer

import (
   "fmt"
   "strconv"
)

// 驗證參數類型
func(s *MCPServer) validateArgumentType(name string, value interface{}, schema PropertySchema)(error) {
   switch schema.Type {
      case "string":
         if _, ok := value.(string); !ok {
            return fmt.Errorf("parameter %s must be a string", name)
         }
         // 檢查枚舉值
         if len(schema.Enum) > 0 {
            strVal := value.(string)
            valid := false
            for _, enumVal := range schema.Enum {
               if strVal == enumVal {
                  valid = true
                  break
               }
            }
            if !valid {
               return fmt.Errorf("parameter %s must be one of: %v", name, schema.Enum)
            }
         }
      case "number":
         switch value.(type) {
            case string:
               _, err := strconv.Atoi(value.(string))
	       return err
	       break
            case float64, int, int64, int32: // 有效的數字類型
	       break
            default:
               return fmt.Errorf("parameter %s must be a number", name)
         }
      case "boolean":
         if _, ok := value.(bool); !ok {
            return fmt.Errorf("parameter %s must be a boolean", name)
         }
      case "array":
         if _, ok := value.([]interface{}); !ok {
            return fmt.Errorf("parameter %s must be an array", name)
         }
      case "object":
         if _, ok := value.(map[string]interface{}); !ok {
            return fmt.Errorf("parameter %s must be an object", name)
         }
   }
   return nil
}

// 驗證工具參數
func(s *MCPServer) validateToolArguments(tool Tool, args map[string]interface{})(error) {
   // 檢查必需參數
   for _, param := range tool.InputSchema.Required {
      if _, exists := args[param]; !exists {
         return fmt.Errorf("missing required parameter: %v", param)
      }
   }
   // 驗證參數類型
   for argName, argValue := range args {
      if propSchema, exists := tool.InputSchema.Properties[argName]; exists {
         if err := s.validateArgumentType(argName, argValue, propSchema); err != nil {
            return err
         }
      }
   }
   return nil
}
