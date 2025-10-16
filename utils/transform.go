package utils

import (
	"data_mapping/models"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/antonmedv/expr"
)

// Utility: GetNestedValue retrieves a value from a nested map by path
func GetNestedValue(data map[string]interface{}, path []string) (interface{}, bool) {
	var current interface{} = data
	for i, key := range path {
		if arr, ok := current.([]interface{}); ok {
			if idx, err := strconv.Atoi(key); err == nil && idx >= 0 && idx < len(arr) {
				if i == len(path)-1 {
					return arr[idx], true
				}
				current = arr[idx]
				continue
			}
			return nil, false
		}
		if currentMap, ok := current.(map[string]interface{}); ok {
			val, exists := currentMap[key]
			if !exists {
				return nil, false
			}
			if i == len(path)-1 {
				return val, true
			}
			current = val
		} else {
			return nil, false
		}
	}
	return nil, false
}

// Utility: SetNestedValue sets a value in a nested map by path
func SetNestedValue(data map[string]interface{}, path []string, value interface{}) {
	current := data
	for i, key := range path {
		if i == len(path)-1 {
			current[key] = value
			return
		}
		if _, exists := current[key]; !exists {
			current[key] = make(map[string]interface{})
		}
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			newMap := make(map[string]interface{})
			current[key] = newMap
			current = newMap
		}
	}
}

func Transform(input map[string]interface{}, rules []models.MappingRule) (map[string]interface{}, error) {
	return ApplyRules(input, rules), nil
}

func ApplyRules(input map[string]interface{}, rules []models.MappingRule) map[string]interface{} {
	output := make(map[string]interface{})
	for _, rule := range rules {
		val, exists := GetNestedValue(input, rule.SourcePath)
		var transformedVal interface{}
		var err error

		if !exists {
			// Handle case where source field doesn't exist
			if rule.Required {
				// For required fields, use default value if provided
				if rule.DefaultValue != "" {
					// Try to parse default value based on expected type
					var defaultVal interface{}
					defaultVal = rule.DefaultValue

					// Check if it's a boolean
					if rule.DefaultValue == "true" || rule.DefaultValue == "false" {
						defaultVal = (rule.DefaultValue == "true")
					} else if val, err := strconv.Atoi(rule.DefaultValue); err == nil {
						// Check if it's an integer
						defaultVal = val
					} else if val, err := strconv.ParseFloat(rule.DefaultValue, 64); err == nil {
						// Check if it's a float
						defaultVal = val
					}

					SetNestedValue(output, rule.DestinationPath, defaultVal)
				} else {
					// No default value provided, but field is required
					// Set an empty value based on destination field name hints
					destField := rule.DestinationPath[len(rule.DestinationPath)-1]

					// Try to infer type from field name
					if strings.Contains(strings.ToLower(destField), "count") ||
						strings.Contains(strings.ToLower(destField), "number") ||
						strings.Contains(strings.ToLower(destField), "id") {
						SetNestedValue(output, rule.DestinationPath, 0)
					} else if strings.Contains(strings.ToLower(destField), "is") ||
						strings.Contains(strings.ToLower(destField), "has") {
						SetNestedValue(output, rule.DestinationPath, false)
					} else {
						SetNestedValue(output, rule.DestinationPath, "")
					}
				}
			}
			continue
		}

		// Normal transformation for existing fields
		if rule.TransformType == "expression" || rule.TransformLogic != "" {
			// Create a rich context with various helper functions and input data
			params := map[string]interface{}{
				"value":      val,
				"input":      input,
				"output":     output,
				"sourcePath": rule.SourcePath,
				"destPath":   rule.DestinationPath,
				"rule":       rule,
			}

			// Use transform logic if available, otherwise create a simple expression that just returns the value
			exprToEval := rule.TransformLogic
			if exprToEval == "" {
				exprToEval = "value"
			}

			transformedVal, err = EvaluateExpression(exprToEval, params)

			// If result is a JSON string, try to parse it
			if err == nil {
				if jsonStr, ok := transformedVal.(string); ok {
					if strings.HasPrefix(jsonStr, "[") || strings.HasPrefix(jsonStr, "{") {
						var jsonObj interface{}
						if jsonErr := json.Unmarshal([]byte(jsonStr), &jsonObj); jsonErr == nil {
							transformedVal = jsonObj
						}
					}
				}
			}
		} else {
			transformedVal, err = ApplyTransform(val, rule.TransformType)
		}
		if err == nil {
			SetNestedValue(output, rule.DestinationPath, transformedVal)
		}
	}
	return output
}
func ApplyTransform(value interface{}, transformType string) (interface{}, error) {
	return value, nil
}

// StreamTransformJSON streams and transforms large client JSONs in real-time.
func StreamTransformJSON(r io.Reader, w io.Writer, transform func(key string, value interface{}) (string, interface{})) error {
	dec := json.NewDecoder(r)

	t, err := dec.Token()
	if err != nil || t != json.Delim('{') {
		return fmt.Errorf("expected start of object: %v", err)
	}
	w.Write([]byte("{"))
	first := true
	for dec.More() {
		keyToken, err := dec.Token()
		if err != nil {
			return err
		}
		key := keyToken.(string)
		var value interface{}
		if err := dec.Decode(&value); err != nil {
			return err
		}
		newKey, newValue := transform(key, value)
		if !first {
			w.Write([]byte(","))
		}
		first = false
		keyBytes, _ := json.Marshal(newKey)
		valueBytes, _ := json.Marshal(newValue)
		w.Write(keyBytes)
		w.Write([]byte(":"))
		w.Write(valueBytes)
	}
	t, err = dec.Token()
	if err != nil || t != json.Delim('}') {
		return fmt.Errorf("expected end of object: %v", err)
	}
	w.Write([]byte("}"))
	return nil
}

// StreamTransformJSONWithRules streams and transforms large JSONs using the same rules as the standard transform logic.
func StreamTransformJSONWithRules(r io.Reader, w io.Writer, rules []models.MappingRule) error {
	dec := json.NewDecoder(r)
	t, err := dec.Token()
	if err != nil || t != json.Delim('{') {
		return fmt.Errorf("expected start of object: %v", err)
	}
	w.Write([]byte("{"))
	first := true
	for dec.More() {
		keyToken, err := dec.Token()
		if err != nil {
			return err
		}
		key := keyToken.(string)
		var value interface{}
		if err := dec.Decode(&value); err != nil {
			return err
		}
		// Use ApplyRules for each top-level object
		var transformed interface{}
		if vMap, ok := value.(map[string]interface{}); ok {
			transformed = ApplyRules(vMap, rules)
		} else {
			transformed = value
		}
		if !first {
			w.Write([]byte(","))
		}
		first = false
		keyBytes, _ := json.Marshal(key)
		valueBytes, _ := json.Marshal(transformed)
		w.Write(keyBytes)
		w.Write([]byte(":"))
		w.Write(valueBytes)
	}
	t, err = dec.Token()
	if err != nil || t != json.Delim('}') {
		return fmt.Errorf("expected end of object: %v", err)
	}
	w.Write([]byte("}"))
	return nil
}

// expressionFuncs contains reusable functions for expression evaluation

// EvaluateExpression evaluates an expression with rich context and helper functions
func EvaluateExpression(expression string, context map[string]interface{}) (interface{}, error) {
	// Create a set of functions for the expression environment
	env := map[string]interface{}{
		// Pass through all existing context
		"value":  context["value"],
		"input":  context["input"],
		"output": context["output"],

		// Date/time functions
		"formatDate": func(dateStr string, format string) string {
			formats := []string{
				"02-January-2006",
				"02-Jan-2006",
				"02/January/2006",
				"02-January-06",
				"2006-01-02",
				time.RFC3339,
			}
			for _, f := range formats {
				if t, err := time.Parse(f, dateStr); err == nil {
					return t.Format(format)
				}
			}
			return dateStr
		},

		// String manipulation
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
		"trim":    strings.TrimSpace,
		"replace": strings.Replace,
		"contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},
		"startsWith": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"endsWith": func(s, suffix string) bool {
			return strings.HasSuffix(s, suffix)
		},
		"capitalize": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
		},

		// Type conversion
		"toInt": func(v interface{}) int {
			switch val := v.(type) {
			case string:
				i, _ := strconv.Atoi(val)
				return i
			case float64:
				return int(val)
			case int:
				return val
			default:
				return 0
			}
		},
		"toFloat": func(v interface{}) float64 {
			switch val := v.(type) {
			case string:
				f, _ := strconv.ParseFloat(val, 64)
				return f
			case float64:
				return val
			case int:
				return float64(val)
			default:
				return 0
			}
		},
		"toString": func(v interface{}) string {
			return fmt.Sprintf("%v", v)
		},
		"toBool": func(v interface{}) bool {
			switch val := v.(type) {
			case bool:
				return val
			case string:
				s := strings.ToLower(val)
				return s == "true" || s == "yes" || s == "1" || s == "y"
			case int:
				return val != 0
			case float64:
				return val != 0
			default:
				return false
			}
		},

		// Array/slice helpers
		"join":  strings.Join,
		"split": strings.Split,
		"length": func(v interface{}) int {
			switch val := v.(type) {
			case string:
				return len(val)
			case []interface{}:
				return len(val)
			case map[string]interface{}:
				return len(val)
			default:
				return 0
			}
		},

		// Path helpers
		"getPath": func(data map[string]interface{}, path ...string) interface{} {
			val, exists := GetNestedValue(data, path)
			if !exists {
				return nil
			}
			return val
		},

		// Conditional helpers
		"ifThen": func(condition bool, trueVal, falseVal interface{}) interface{} {
			if condition {
				return trueVal
			}
			return falseVal
		},
		"coalesce": func(val1, val2 interface{}) interface{} {
			if val1 == nil {
				return val2
			}
			return val1
		},

		// Math operations
		"add":      func(a, b float64) float64 { return a + b },
		"subtract": func(a, b float64) float64 { return a - b },
		"multiply": func(a, b float64) float64 { return a * b },
		"divide": func(a, b float64) float64 {
			if b == 0 {
				return 0 // Prevent division by zero
			}
			return a / b
		},
		"round": func(val float64, precision int) float64 {
			p := math.Pow10(precision)
			return math.Round(val*p) / p
		},

		// Current time
		"now":     time.Now(),
		"today":   time.Now().Format("2006-01-02"),
		"isoDate": time.Now().Format(time.RFC3339),
	}

	// Add any other context variables
	for k, v := range context {
		if _, exists := env[k]; !exists {
			env[k] = v
		}
	}

	// Evaluate the expression with the enriched context
	return expr.Eval(expression, env)
}
