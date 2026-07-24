//go:build api_parity

package tgbotapi

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type parityDoc struct {
	Methods map[string]parityMethod `json:"methods"`
	Types   map[string]parityType   `json:"types"`
}

type parityMethod struct {
	Parameters []parityParameter `json:"parameters"`
	Fields     []parityParameter `json:"fields"`
}

type parityParameter struct {
	Name string `json:"name"`
}

type parityType struct {
	Fields []parityField `json:"fields"`
}

type parityField struct {
	Name string `json:"name"`
}

type parityTypeDecl struct {
	Expr ast.Expr
}

func (method parityMethod) parameterFields() []parityParameter {
	if len(method.Parameters) > 0 {
		return method.Parameters
	}
	return method.Fields
}

func TestAPIParityMethods(t *testing.T) {
	doc := loadParityDoc(t)
	index := loadPackageIndex(t)

	implementedMethods := index.methodNames
	for _, name := range []string{"getMe", "getWebhookInfo"} {
		implementedMethods[name] = struct{}{}
	}

	missingMethods := make([]string, 0)
	for name := range doc.Methods {
		if _, ok := implementedMethods[name]; !ok {
			missingMethods = append(missingMethods, name)
		}
	}

	extraMethods := make([]string, 0)
	for name := range implementedMethods {
		if _, ok := doc.Methods[name]; !ok {
			extraMethods = append(extraMethods, name)
		}
	}

	sort.Strings(missingMethods)
	sort.Strings(extraMethods)

	if len(missingMethods) > 0 || len(extraMethods) > 0 {
		t.Fatalf("method parity failed, missing=%v extra=%v", missingMethods, extraMethods)
	}
}

func TestAPIParityMethodParameters(t *testing.T) {
	doc := loadParityDoc(t)
	index := loadPackageIndex(t)

	allowedExtraParams := map[string]map[string]string{
		"approveSuggestedPost": {
			"business_connection_id": "implemented before current docs listed it for this method",
		},
		"banChatSenderChat": {
			"until_date": "legacy compatibility field",
		},
		"closeGeneralForumTopic": {
			"message_thread_id": "promoted legacy forum-topic field",
		},
		"copyMessage": {
			"business_connection_id": "implemented before current docs listed it for this method",
		},
		"copyMessages": {
			"allow_paid_broadcast":      "shared forward/copy compatibility field",
			"business_connection_id":    "shared forward/copy compatibility field",
			"message_effect_id":         "shared forward/copy compatibility field",
			"reply_markup":              "shared forward/copy compatibility field",
			"reply_parameters":          "shared forward/copy compatibility field",
			"suggested_post_parameters": "shared forward/copy compatibility field",
		},
		"declineSuggestedPost": {
			"business_connection_id": "implemented before current docs listed it for this method",
		},
		"deleteBusinessMessages": {
			"chat_id": "promoted legacy chat field",
		},
		"deleteMessage": {
			"business_connection_id": "implemented before current docs listed it for this method",
		},
		"deleteMessageReaction": {
			"business_connection_id": "implemented before current docs listed it for this method",
		},
		"deleteMessages": {
			"business_connection_id": "implemented before current docs listed it for this method",
		},
		"editGeneralForumTopic": {
			"message_thread_id": "promoted legacy forum-topic field",
		},
		"forwardMessage": {
			"allow_paid_broadcast":   "shared forward/copy compatibility field",
			"business_connection_id": "shared forward/copy compatibility field",
			"reply_markup":           "shared forward/copy compatibility field",
			"reply_parameters":       "shared forward/copy compatibility field",
		},
		"forwardMessages": {
			"allow_paid_broadcast":      "shared forward/copy compatibility field",
			"business_connection_id":    "shared forward/copy compatibility field",
			"message_effect_id":         "shared forward/copy compatibility field",
			"reply_markup":              "shared forward/copy compatibility field",
			"reply_parameters":          "shared forward/copy compatibility field",
			"suggested_post_parameters": "shared forward/copy compatibility field",
		},
		"getGameHighScores": {
			"business_connection_id": "promoted edit/game compatibility field",
		},
		"hideGeneralForumTopic": {
			"message_thread_id": "promoted legacy forum-topic field",
		},
		"reopenGeneralForumTopic": {
			"message_thread_id": "promoted legacy forum-topic field",
		},
		"sendChatAction": {
			"allow_paid_broadcast":      "promoted BaseChat field",
			"direct_messages_topic_id":  "promoted BaseChat field",
			"disable_notification":      "promoted BaseChat field",
			"message_effect_id":         "promoted BaseChat field",
			"protect_content":           "promoted BaseChat field",
			"reply_markup":              "promoted BaseChat field",
			"reply_parameters":          "promoted BaseChat field",
			"suggested_post_parameters": "promoted BaseChat field",
		},
		"sendChecklist": {
			"allow_paid_broadcast":      "promoted BaseChat field",
			"direct_messages_topic_id":  "promoted BaseChat field",
			"message_thread_id":         "promoted BaseChat field",
			"suggested_post_parameters": "promoted BaseChat field",
		},
		"sendGame": {
			"direct_messages_topic_id":  "promoted BaseChat field",
			"suggested_post_parameters": "promoted BaseChat field",
		},
		"sendInvoice": {
			"business_connection_id": "promoted BaseChat field",
		},
		"sendMediaGroup": {
			"reply_markup":              "promoted BaseChat field",
			"suggested_post_parameters": "promoted BaseChat field",
		},
		"sendPaidMedia": {
			"message_effect_id": "promoted BaseChat field",
		},
		"sendPhoto": {
			"thumbnail": "legacy compatibility field",
		},
		"sendPoll": {
			"direct_messages_topic_id":  "promoted BaseChat field",
			"suggested_post_parameters": "promoted BaseChat field",
		},
		"sendVoice": {
			"thumbnail": "legacy compatibility alias for older Bot API docs",
		},
		"setChatPhoto": {
			"allow_paid_broadcast":      "promoted BaseChat field",
			"business_connection_id":    "promoted BaseChat field",
			"direct_messages_topic_id":  "promoted BaseChat field",
			"disable_notification":      "promoted BaseChat field",
			"message_effect_id":         "promoted BaseChat field",
			"message_thread_id":         "promoted BaseChat field",
			"protect_content":           "promoted BaseChat field",
			"reply_markup":              "promoted BaseChat field",
			"reply_parameters":          "promoted BaseChat field",
			"suggested_post_parameters": "promoted BaseChat field",
		},
		"setGameScore": {
			"business_connection_id": "promoted edit/game compatibility field",
		},
		"setMessageReaction": {
			"business_connection_id": "implemented before current docs listed it for this method",
		},
		"stopPoll": {
			"inline_message_id": "legacy edit-style compatibility field",
		},
		"unhideGeneralForumTopic": {
			"message_thread_id": "promoted legacy forum-topic field",
		},
		"unpinAllGeneralForumTopicMessages": {
			"message_thread_id": "promoted legacy forum-topic field",
		},
	}

	missingParams := make(map[string][]string)
	extraParams := make(map[string][]string)

	for methodName, method := range doc.Methods {
		implemented, ok := index.collectMethodParams(methodName)
		if !ok {
			continue
		}

		expected := make(map[string]struct{})
		for _, parameter := range method.parameterFields() {
			expected[parameter.Name] = struct{}{}
			if _, exists := implemented[parameter.Name]; !exists {
				missingParams[methodName] = append(missingParams[methodName], parameter.Name)
			}
		}

		for parameter := range implemented {
			if _, exists := expected[parameter]; exists {
				continue
			}
			if _, allowed := allowedExtraParams[methodName][parameter]; allowed {
				continue
			}
			extraParams[methodName] = append(extraParams[methodName], parameter)
		}
	}

	unusedAllowed := unusedAllowedMethodParams(allowedExtraParams, index, doc)
	for _, params := range missingParams {
		sort.Strings(params)
	}
	for _, params := range extraParams {
		sort.Strings(params)
	}

	missingMethods := sortedMapKeys(missingParams)
	extraMethods := sortedMapKeys(extraParams)
	sort.Strings(unusedAllowed)

	if len(missingMethods) > 0 || len(extraMethods) > 0 || len(unusedAllowed) > 0 {
		builder := strings.Builder{}
		if len(missingMethods) > 0 {
			builder.WriteString("methods with missing params:\n")
			for _, methodName := range missingMethods {
				builder.WriteString(fmt.Sprintf("- %s: %v\n", methodName, missingParams[methodName]))
			}
		}
		if len(extraMethods) > 0 {
			builder.WriteString("methods with unexpected params:\n")
			for _, methodName := range extraMethods {
				builder.WriteString(fmt.Sprintf("- %s: %v\n", methodName, extraParams[methodName]))
			}
		}
		if len(unusedAllowed) > 0 {
			builder.WriteString(fmt.Sprintf("allowed extra params no longer needed: %v\n", unusedAllowed))
		}
		t.Fatalf("method parameter parity failed\n%s", builder.String())
	}
}

func TestAPIParityTypesAndFields(t *testing.T) {
	doc := loadParityDoc(t)
	index := loadPackageIndex(t)

	allowedMissingTypeNames := map[string]string{
		"ChatBoostSourcePremium":  "type name conflicts with exported constant",
		"ChatBoostSourceGiftCode": "type name conflicts with exported constant",
		"ChatBoostSourceGiveaway": "type name conflicts with exported constant",
		"MessageOriginUser":       "type name conflicts with exported constant",
		"MessageOriginHiddenUser": "type name conflicts with exported constant",
		"MessageOriginChat":       "type name conflicts with exported constant",
		"MessageOriginChannel":    "type name conflicts with exported constant",
		"ReactionTypeEmoji":       "type name conflicts with exported constant",
		"ReactionTypeCustomEmoji": "type name conflicts with exported constant",
		"ReactionTypePaid":        "type name conflicts with exported constant",
	}

	missingTypes := make([]string, 0)
	missingFields := make(map[string][]string)

	for typeName, apiType := range doc.Types {
		if _, ok := index.types[typeName]; !ok {
			if _, allowed := allowedMissingTypeNames[typeName]; allowed {
				continue
			}
			missingTypes = append(missingTypes, typeName)
			continue
		}

		if len(apiType.Fields) == 0 {
			continue
		}

		fields, ok := index.collectJSONFields(typeName)
		if !ok {
			missingFields[typeName] = []string{"<non-struct>"}
			continue
		}

		for _, field := range apiType.Fields {
			if _, exists := fields[field.Name]; !exists {
				missingFields[typeName] = append(missingFields[typeName], field.Name)
			}
		}
	}

	unusedAllowed := make([]string, 0)
	for typeName := range allowedMissingTypeNames {
		if _, ok := index.types[typeName]; ok {
			unusedAllowed = append(unusedAllowed, typeName)
		}
	}

	sort.Strings(missingTypes)
	sort.Strings(unusedAllowed)

	missingFieldTypes := make([]string, 0, len(missingFields))
	for typeName, fields := range missingFields {
		sort.Strings(fields)
		missingFieldTypes = append(missingFieldTypes, typeName)
	}
	sort.Strings(missingFieldTypes)

	if len(missingTypes) > 0 || len(missingFieldTypes) > 0 || len(unusedAllowed) > 0 {
		builder := strings.Builder{}
		if len(missingTypes) > 0 {
			builder.WriteString(fmt.Sprintf("missing types: %v\n", missingTypes))
		}
		if len(missingFieldTypes) > 0 {
			builder.WriteString("types with missing fields:\n")
			for _, typeName := range missingFieldTypes {
				builder.WriteString(fmt.Sprintf("- %s: %v\n", typeName, missingFields[typeName]))
			}
		}
		if len(unusedAllowed) > 0 {
			builder.WriteString(fmt.Sprintf("allowed missing types no longer needed: %v\n", unusedAllowed))
		}
		t.Fatalf("type parity failed\n%s", builder.String())
	}
}

type packageIndex struct {
	types           map[string]parityTypeDecl
	methodNames     map[string]struct{}
	methodReceivers map[string]string
	paramsFuncs     map[string]*ast.FuncDecl
	filesFuncs      map[string]*ast.FuncDecl
}

func loadPackageIndex(t *testing.T) *packageIndex {
	t.Helper()

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}

	goFiles := make([]string, 0)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		goFiles = append(goFiles, name)
	}
	sort.Strings(goFiles)

	fset := token.NewFileSet()
	index := &packageIndex{
		types:           make(map[string]parityTypeDecl),
		methodNames:     make(map[string]struct{}),
		methodReceivers: make(map[string]string),
		paramsFuncs:     make(map[string]*ast.FuncDecl),
		filesFuncs:      make(map[string]*ast.FuncDecl),
	}

	for _, file := range goFiles {
		parsed, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}

		for _, decl := range parsed.Decls {
			switch current := decl.(type) {
			case *ast.GenDecl:
				if current.Tok != token.TYPE {
					continue
				}
				for _, spec := range current.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					index.types[typeSpec.Name.Name] = parityTypeDecl{Expr: typeSpec.Type}
				}
			case *ast.FuncDecl:
				if current.Recv == nil || current.Name == nil {
					continue
				}
				receiver, ok := receiverTypeName(current.Recv)
				if !ok {
					continue
				}
				switch current.Name.Name {
				case "method":
					if current.Body == nil || len(current.Body.List) == 0 {
						continue
					}
					for _, statement := range current.Body.List {
						returnStatement, ok := statement.(*ast.ReturnStmt)
						if !ok || len(returnStatement.Results) != 1 {
							continue
						}

						literal, ok := returnStatement.Results[0].(*ast.BasicLit)
						if !ok || literal.Kind != token.STRING {
							continue
						}

						methodName, err := strconv.Unquote(literal.Value)
						if err != nil {
							t.Fatalf("unquote method literal in %s: %v", file, err)
						}
						index.methodNames[methodName] = struct{}{}
						index.methodReceivers[methodName] = receiver
						break
					}
				case "params":
					index.paramsFuncs[receiver] = current
				case "files":
					index.filesFuncs[receiver] = current
				}
			}
		}
	}

	return index
}

func loadParityDoc(t *testing.T) parityDoc {
	t.Helper()

	content, err := os.ReadFile(filepath.Join(".", "api.json"))
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("api.json not found, skipping api parity tests")
			return parityDoc{}
		}
		t.Fatalf("read api.json: %v", err)
	}

	trimmed := strings.TrimSpace(string(content))
	if trimmed == "" {
		t.Fatalf("api.json is empty")
	}
	if trimmed[0] == '<' {
		preview := trimmed
		if newline := strings.IndexByte(preview, '\n'); newline >= 0 {
			preview = preview[:newline]
		}
		if len(preview) > 120 {
			preview = preview[:120] + "..."
		}
		t.Fatalf("api.json contains HTML instead of JSON; fetch the raw file, not the GitHub blob page (starts with %q)", preview)
	}

	var doc parityDoc
	if err := json.Unmarshal(content, &doc); err != nil {
		t.Fatalf("unmarshal api.json: %v", err)
	}

	return doc
}

func receiverTypeName(recv *ast.FieldList) (string, bool) {
	if recv == nil || len(recv.List) == 0 {
		return "", false
	}
	return exprTypeName(recv.List[0].Type)
}

func exprTypeName(expr ast.Expr) (string, bool) {
	switch current := expr.(type) {
	case *ast.Ident:
		return current.Name, true
	case *ast.StarExpr:
		return exprTypeName(current.X)
	case *ast.IndexExpr:
		return exprTypeName(current.X)
	case *ast.IndexListExpr:
		return exprTypeName(current.X)
	case *ast.SelectorExpr:
		return current.Sel.Name, true
	default:
		return "", false
	}
}

func sortedMapKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func unusedAllowedMethodParams(allowed map[string]map[string]string, index *packageIndex, doc parityDoc) []string {
	unused := make([]string, 0)
	for methodName, params := range allowed {
		implemented, ok := index.collectMethodParams(methodName)
		if !ok {
			for param := range params {
				unused = append(unused, methodName+"."+param)
			}
			continue
		}
		expected := make(map[string]struct{})
		for _, parameter := range doc.Methods[methodName].parameterFields() {
			expected[parameter.Name] = struct{}{}
		}
		for param := range params {
			if _, exists := implemented[param]; !exists {
				unused = append(unused, methodName+"."+param)
				continue
			}
			if _, exists := expected[param]; exists {
				unused = append(unused, methodName+"."+param)
			}
		}
	}
	return unused
}

func (index *packageIndex) collectMethodParams(methodName string) (map[string]struct{}, bool) {
	receiver, ok := index.methodReceivers[methodName]
	if !ok {
		switch methodName {
		case "getMe", "getWebhookInfo":
			return map[string]struct{}{}, true
		default:
			return nil, false
		}
	}

	params, ok := index.collectParamsByType(receiver, map[string]bool{})
	if !ok {
		params = make(map[string]struct{})
	}
	files := index.collectFileParamsByType(receiver)
	for name := range files {
		params[name] = struct{}{}
	}
	return params, true
}

func (index *packageIndex) collectParamsByType(typeName string, path map[string]bool) (map[string]struct{}, bool) {
	if path[typeName] {
		return nil, false
	}
	path[typeName] = true
	defer delete(path, typeName)

	if paramsFunc, ok := index.paramsFuncs[typeName]; ok {
		params := make(map[string]struct{})
		ast.Inspect(paramsFunc.Body, func(node ast.Node) bool {
			switch current := node.(type) {
			case *ast.CallExpr:
				index.collectCallParam(typeName, current, params, path)
			case *ast.IndexExpr:
				if key, ok := stringIndexKey(current); ok {
					params[key] = struct{}{}
				}
			}
			return true
		})
		return params, true
	}

	embedded := index.embeddedTypeNames(typeName)
	if len(embedded) == 0 {
		return nil, false
	}

	params := make(map[string]struct{})
	for _, embeddedType := range embedded {
		embeddedParams, ok := index.collectParamsByType(embeddedType, path)
		if !ok {
			continue
		}
		for name := range embeddedParams {
			params[name] = struct{}{}
		}
	}
	return params, true
}

func (index *packageIndex) collectCallParam(receiverType string, call *ast.CallExpr, params map[string]struct{}, path map[string]bool) {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	switch selector.Sel.Name {
	case "AddNonEmpty", "AddNonZero", "AddNonZero64", "AddNonZeroFloat", "AddBool", "AddBoolValue", "AddBoolPtr", "AddInterface", "AddInterfaceNonZero", "AddFirstValid", "paramsWithKey":
		if key, ok := firstStringArg(call); ok {
			params[key] = struct{}{}
		}
	case "params":
		calledType, ok := index.paramReceiverExprType(receiverType, selector.X)
		if !ok {
			return
		}
		nested, ok := index.collectParamsByType(calledType, path)
		if !ok {
			return
		}
		for name := range nested {
			params[name] = struct{}{}
		}
	}
}

func (index *packageIndex) collectFileParamsByType(typeName string) map[string]struct{} {
	params := make(map[string]struct{})
	filesFunc, ok := index.filesFuncs[typeName]
	if !ok {
		return params
	}

	ast.Inspect(filesFunc.Body, func(node ast.Node) bool {
		switch current := node.(type) {
		case *ast.CallExpr:
			selector, ok := current.Fun.(*ast.Ident)
			if !ok || selector.Name != "requestFile" {
				return true
			}
			if key, ok := firstStringArg(current); ok {
				params[key] = struct{}{}
			}
		case *ast.KeyValueExpr:
			ident, ok := current.Key.(*ast.Ident)
			if !ok || ident.Name != "Name" {
				return true
			}
			if key, ok := stringLiteralValue(current.Value); ok {
				params[key] = struct{}{}
			}
		}
		return true
	})

	return params
}

func (index *packageIndex) paramReceiverExprType(receiverType string, expr ast.Expr) (string, bool) {
	switch current := expr.(type) {
	case *ast.Ident:
		return receiverType, true
	case *ast.SelectorExpr:
		baseType, ok := index.paramReceiverExprType(receiverType, current.X)
		if !ok {
			return "", false
		}
		return index.structFieldType(baseType, current.Sel.Name)
	default:
		return "", false
	}
}

func (index *packageIndex) structFieldType(typeName, fieldName string) (string, bool) {
	decl, ok := index.types[typeName]
	if !ok {
		return "", false
	}

	structType, ok := index.structTypeByExpr(decl.Expr)
	if !ok {
		return "", false
	}

	for _, field := range structType.Fields.List {
		fieldType, ok := exprTypeName(field.Type)
		if !ok {
			continue
		}
		if len(field.Names) == 0 {
			if fieldType == fieldName {
				return fieldType, true
			}
			if nestedType, ok := index.structFieldType(fieldType, fieldName); ok {
				return nestedType, true
			}
			continue
		}
		for _, name := range field.Names {
			if name.Name == fieldName {
				return fieldType, true
			}
		}
	}

	return "", false
}

func (index *packageIndex) embeddedTypeNames(typeName string) []string {
	decl, ok := index.types[typeName]
	if !ok {
		return nil
	}

	structType, ok := index.structTypeByExpr(decl.Expr)
	if !ok {
		return nil
	}

	names := make([]string, 0)
	for _, field := range structType.Fields.List {
		if len(field.Names) != 0 {
			continue
		}
		if name, ok := exprTypeName(field.Type); ok {
			names = append(names, name)
		}
	}
	return names
}

func (index *packageIndex) structTypeByExpr(expr ast.Expr) (*ast.StructType, bool) {
	switch current := expr.(type) {
	case *ast.StructType:
		return current, true
	case *ast.Ident:
		decl, ok := index.types[current.Name]
		if !ok {
			return nil, false
		}
		return index.structTypeByExpr(decl.Expr)
	default:
		return nil, false
	}
}

func firstStringArg(call *ast.CallExpr) (string, bool) {
	if len(call.Args) == 0 {
		return "", false
	}
	return stringLiteralValue(call.Args[0])
}

func stringIndexKey(index *ast.IndexExpr) (string, bool) {
	return stringLiteralValue(index.Index)
}

func stringLiteralValue(expr ast.Expr) (string, bool) {
	literal, ok := expr.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return "", false
	}
	value, err := strconv.Unquote(literal.Value)
	if err != nil {
		return "", false
	}
	return value, true
}

func (index *packageIndex) collectJSONFields(typeName string) (map[string]struct{}, bool) {
	return index.collectJSONFieldsByExpr(&ast.Ident{Name: typeName}, map[string]bool{})
}

func (index *packageIndex) collectJSONFieldsByExpr(expr ast.Expr, path map[string]bool) (map[string]struct{}, bool) {
	switch current := expr.(type) {
	case *ast.ParenExpr:
		return index.collectJSONFieldsByExpr(current.X, path)
	case *ast.StarExpr:
		return index.collectJSONFieldsByExpr(current.X, path)
	case *ast.Ident:
		if path[current.Name] {
			return nil, false
		}
		path[current.Name] = true
		defer delete(path, current.Name)

		decl, ok := index.types[current.Name]
		if !ok {
			return nil, false
		}

		switch typed := decl.Expr.(type) {
		case *ast.StructType:
			fields := make(map[string]struct{})
			for _, field := range typed.Fields.List {
				tagName, hasTag := parseJSONTag(field.Tag)

				if len(field.Names) > 0 {
					for _, fieldName := range field.Names {
						jsonName := fallbackJSONFieldName(tagName, hasTag, fieldName.Name)
						if jsonName != "" {
							fields[jsonName] = struct{}{}
						}
					}
					continue
				}

				if hasTag {
					if tagName != "" {
						fields[tagName] = struct{}{}
					}
					continue
				}

				embeddedFields, ok := index.collectJSONFieldsByExpr(field.Type, path)
				if !ok {
					continue
				}
				for name := range embeddedFields {
					fields[name] = struct{}{}
				}
			}
			return fields, true
		default:
			return index.collectJSONFieldsByExpr(decl.Expr, path)
		}
	default:
		return nil, false
	}
}

func parseJSONTag(tag *ast.BasicLit) (string, bool) {
	if tag == nil {
		return "", false
	}

	value, err := strconv.Unquote(tag.Value)
	if err != nil {
		return "", false
	}

	parsed := reflect.StructTag(value).Get("json")
	if parsed == "" {
		return "", false
	}

	name := strings.Split(parsed, ",")[0]
	if name == "-" {
		return "", true
	}

	return name, true
}

func fallbackJSONFieldName(tagName string, hasTag bool, fieldName string) string {
	if hasTag {
		if tagName == "" {
			return fieldName
		}
		return tagName
	}
	return fieldName
}
