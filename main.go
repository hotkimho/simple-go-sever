package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"
	"text/template/parse"

	sprig "github.com/go-task/slim-sprig/v3"
)

// extractFields 재귀적으로 템플릿 파스 트리를 순회하면서
// FieldNode에 나타난 필드 경로들을 추출합니다.
// 예를 들어, .values.type.dockerfileContents 라는 경로를 "values.type.dockerfileContents"로 저장합니다.
func extractFields(node parse.Node, paths *[]string) {
	fmt.Println("node : ", node.String(), " ", "paths : ", *paths)
	switch n := node.(type) {
	case *parse.ListNode:
		fmt.Println("listNode : ", n.String())
		if n != nil {
			for _, node := range n.Nodes {
				extractFields(node, paths)
			}
		}
	case *parse.ActionNode:
		fmt.Println("actionNode : ", n.String())
		if n.Pipe != nil {
			for _, cmd := range n.Pipe.Cmds {
				for _, arg := range cmd.Args {
					extractFields(arg, paths)
				}
			}
		}
	case *parse.FieldNode:
		fmt.Println("fieldNode : ", n.String())
		// FieldNode.Ident는 []string로 구성되어 있음.
		*paths = append(*paths, strings.Join(n.Ident, "."))
	case *parse.PipeNode:
		fmt.Println("pipeNode : ", n.String())
		for _, cmd := range n.Cmds {
			extractFields(cmd, paths)
		}
	case *parse.CommandNode:
		fmt.Println("commandNode : ", n.String())
		for _, arg := range n.Args {
			extractFields(arg, paths)
		}
		// 다른 노드 타입(TextNode, VariableNode 등)은 무시.
	}
}

// buildSafeMapFromPaths는 추출된 필드 경로들을 기반으로 중첩된 맵 구조를 생성합니다.
// 최종 값은 빈 문자열("")로 채워넣습니다.
func buildSafeMapFromPaths(paths []string) map[string]interface{} {
	safeMap := make(map[string]interface{})
	for _, path := range paths {
		parts := strings.Split(path, ".")
		current := safeMap
		for i, part := range parts {
			if i == len(parts)-1 {
				// 마지막 키에는 빈 문자열 할당
				if _, exists := current[part]; !exists {
					current[part] = ""
				}
			} else {
				// 중첩된 맵 생성
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				// 다음 단계로 이동
				if next, ok := current[part].(map[string]interface{}); ok {
					current = next
				} else {
					// 만약 이미 값이 존재하는데 map이 아니라면 무시
					break
				}
			}
		}
	}
	return safeMap
}

func main() {
	// 템플릿 문자열 (수정할 수 없는 상태)
	tmplStr := "echo .PIPELINE.NAME⦀{{{.PIPELINE.NAME}}}\" >> /dev/termination-log; echo \".PIPELINE.NAMESPACE⦀{{{.PIPELINE.NAMESPACE}}}\" >> /dev/termination-log; echo \".PIPELINE.UID⦀{{{.PIPELINE.UID}}}\" >> /dev/termination-log; echo \".PIPELINE.INSTANCE⦀{{{.PIPELINE.INSTANCE}}}\" >> /dev/termination-log; echo \".BUILD.NAME⦀{{{.BUILD.NAME}}}\" >> /dev/termination-log; echo \".BUILD.VERSION⦀{{{.BUILD.VERSION}}}\" >> /dev/termination-log; echo \".BUILD.CREATOR.USERNAME⦀{{{.BUILD.CREATOR.USERNAME}}}\" >> /dev/termination-log;"

	// 템플릿 파싱
	tmpl := template.New("example")
	tmpl.Option("missingkey=default").Delims("{{{", "}}}")
	tmpl = tmpl.Funcs(sprig.TxtFuncMap())
	tmpl, err := tmpl.Parse(tmplStr)
	if err != nil {
		log.Fatalf("template parsing error: %v", err)
	}

	// 템플릿 파스 트리에서 필드 경로 추출
	var paths []string
	extractFields(tmpl.Tree.Root, &paths)
	log.Println("Extracted field paths:", paths)

	safeData := buildSafeMapFromPaths(paths)
	fmt.Println("safeData:", safeData)

	bb := &bytes.Buffer{}
	// 템플릿 실행
	err = tmpl.Execute(bb, safeData)
	if err != nil {
		log.Fatalf("template execution error: %v", err)
	}

	fmt.Println("------------------------------------------------------------")
	fmt.Println("bb : ", bb.String())
}

// func main() {
// 	server := "427800856788.dkr.ecr.ap-northeast-2.amazonaws.com"
// 	domain := "427800856788.dkr.ecr.ap-northeast-2.amazonaws.com"
// 	r := strings.TrimPrefix(server, domain+"/")
// 	fmt.Println("r: ", r)
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}
// 	a := ""
// 	b, err := strconv.Atoi(a)
// 	if err != nil {
// 		fmt.Println("error ", err)
// 	}
// 	fmt.Println("b : ", b)
// 	corsOrigin := os.Getenv("CORS_ORIGIN")
// 	debugMode := os.Getenv("DEBUG_MODE")
// 	logLevel := os.Getenv("LOG_LEVEL")

// 	if corsOrigin == "" {
// 		corsOrigin = "*(default value)"
// 	}
// 	if debugMode == "" {
// 		debugMode = "false(default value)"
// 	}
// 	if logLevel == "" {
// 		logLevel = "info(default value)"
// 	}

// 	fmt.Println("trigger test-5")
// 	fmt.Printf("PORT: %s\n", port)
// 	fmt.Printf("CORS_ORIGIN: %s\n", corsOrigin)
// 	fmt.Printf("DEBUG_MODE: %s\n", debugMode)
// 	fmt.Printf("LOG_LEVEL: %s\n", logLevel)

// 	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Test 11")
// 		// content-type값이 application/json으로 설정하고 chartset=utf-8 설정 삭제
// 		w.Header().Del("Content-Type")
// 		w.Header().Set("Content-Type", "application/json")

// 		w.WriteHeader(http.StatusOK)
// 	})

// 	log.Printf("Starting server on port %s...", port)
// 	if err := http.ListenAndServe(":"+port, nil); err != nil {
// 		log.Fatal(err)
// 	}
// }
