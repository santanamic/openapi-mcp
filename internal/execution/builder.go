package execution

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"openapi-mcp/internal/tools"
)

// BuildRequest constructs an http.Request from an operation definition, base URL, and tool arguments.
func BuildRequest(opDef tools.OperationDefinition, baseURL string, args map[string]any) (*http.Request, error) {
	baseURL = strings.TrimSuffix(baseURL, "/")
	path := opDef.Path
	
	var allParams []*openapi3.ParameterRef
	if opDef.Operation != nil {
		allParams = append([]*openapi3.ParameterRef{}, opDef.Operation.Parameters...)
	}
	
	queryParams := url.Values{}
	
	for _, paramRef := range allParams {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}
		param := paramRef.Value
		val, exists := args[param.Name]
		if !exists {
			if param.Required {
				return nil, fmt.Errorf("parâmetro obrigatório ausente: '%s'", param.Name)
			}
			continue
		}
		
		valStr := fmt.Sprintf("%v", val)
		if param.In == "path" {
			placeholder := fmt.Sprintf("{%s}", param.Name)
			path = strings.ReplaceAll(path, placeholder, url.PathEscape(valStr))
		} else if param.In == "query" {
			queryParams.Add(param.Name, valStr)
		}
	}

	targetURL := baseURL + path
	if len(queryParams) > 0 {
		targetURL += "?" + queryParams.Encode()
	}

	var bodyReader *bytes.Reader
	if opDef.Operation != nil && opDef.Operation.RequestBody != nil && opDef.Operation.RequestBody.Value != nil {
		rbVal, ok := args["request_body"]
		if !ok && opDef.Operation.RequestBody.Value.Required {
			return nil, fmt.Errorf("request_body é obrigatório para esta operação")
		}
		if ok {
			bodyJSON, err := json.Marshal(rbVal)
			if err != nil {
				return nil, fmt.Errorf("falha ao serializar request_body para JSON: %w", err)
			}
			bodyReader = bytes.NewReader(bodyJSON)
		}
	}

	var req *http.Request
	var err error
	if bodyReader != nil {
		req, err = http.NewRequest(opDef.Method, targetURL, bodyReader)
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		req, err = http.NewRequest(opDef.Method, targetURL, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição HTTP: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}
