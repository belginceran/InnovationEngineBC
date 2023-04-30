package parsers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

var markdownParser = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithBlockParsers(),
	),
	goldmark.WithRendererOptions(
		html.WithXHTML(),
	),
)

type MarkdownElementType string

const (
	ElementHeading    MarkdownElementType = "heading"
	ElementCodeBlock  MarkdownElementType = "code_block"
	ElementList       MarkdownElementType = "list"
	ElementBlockQuote MarkdownElementType = "block_quote"
	ElementParagraph  MarkdownElementType = "paragraph"
)

// Represents a markdown element.
type MarkdownElement struct {
	Type    MarkdownElementType
	Content string
	Result  string
}

// Parses a markdown file into an AST representing the markdown document.
func ParseMarkdownIntoAst(source []byte) ast.Node {
	document := markdownParser.Parser().Parse(text.NewReader(source))
	return document
}

// Extracts the code blocks from a provided markdown AST that match the
// languagesToExtract.
func ExtractCodeBlocksFromAst(node ast.Node, source []byte, languagesToExtract []string) []string {
	var commands []string
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == ast.KindFencedCodeBlock {
			codeBlock := node.(*ast.FencedCodeBlock)
			for _, language := range languagesToExtract {
				if string(codeBlock.Language(source)) == language {
					commands = append(commands, extractCommandFromCodeBlock(codeBlock, source))
				}
			}
		}
		return ast.WalkContinue, nil
	})

	return commands
}

// This regex matches HTML comments within markdown blocks that contain
// variables to use within the scenario.
var variableCommentBlockRegex = regexp.MustCompile("(?s)<!--.*?```variables(.*?)```.*?")

// Extracts the variables from a provided markdown AST.
func ExtractScenarioVariablesFromAst(node ast.Node, source []byte) map[string]string {
	scenarioVariables := make(map[string]string)

	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == ast.KindHTMLBlock {
			htmlNode := node.(*ast.HTMLBlock)
			blockContent := extractVariablesFromHTMLBlock(htmlNode, source)
			fmt.Printf("Found HTML block with the content: %s\n", blockContent)
			match := variableCommentBlockRegex.FindStringSubmatch(blockContent)
			fmt.Printf("Found %d matches\n", len(match))

			// Extract the variables from the comment block.
			if len(match) > 1 {
				fmt.Println("Found: ", match[1])
				variables := convertScenarioVariablesToMap(match[1])
				for key, value := range variables {
					scenarioVariables[key] = value
				}
			}
		}
		return ast.WalkContinue, nil
	})

	return scenarioVariables
}

func convertScenarioVariablesToMap(variableBlock string) map[string]string {
	variableMap := make(map[string]string)

	// Only process statements that begin with export.
	for _, variable := range strings.Split(variableBlock, "\n") {
		if strings.HasPrefix(variable, "export") {
			parts := strings.SplitN(variable, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimPrefix(parts[0], "export ")
				value := parts[1]
				fmt.Printf("Found variable: %s=%s\n", key, value)
				variableMap[key] = value
			}
		}

	}

	return variableMap
}

// Extracts the command text from an already parsed markdown code block.
func extractCommandFromCodeBlock(codeBlock *ast.FencedCodeBlock, source []byte) string {
	lines := codeBlock.Lines()
	var command strings.Builder

	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		command.WriteString(string(line.Value(source)))
	}

	return command.String()
}

// TODO: Merge this with the above function.
func extractVariablesFromHTMLBlock(htmlBlock *ast.HTMLBlock, source []byte) string {
	lines := htmlBlock.Lines()
	var command strings.Builder

	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		command.WriteString(string(line.Value(source)))
	}

	return command.String()
}