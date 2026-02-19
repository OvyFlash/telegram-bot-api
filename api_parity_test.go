//go:build api_parity
// +build api_parity

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
	Methods map[string]json.RawMessage `json:"methods"`
	Types   map[string]parityType      `json:"types"`
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
	types       map[string]parityTypeDecl
	methodNames map[string]struct{}
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
		types:       make(map[string]parityTypeDecl),
		methodNames: make(map[string]struct{}),
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
				if current.Recv == nil || current.Name == nil || current.Name.Name != "method" {
					continue
				}
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
					break
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

	var doc parityDoc
	if err := json.Unmarshal(content, &doc); err != nil {
		t.Fatalf("unmarshal api.json: %v", err)
	}

	return doc
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
