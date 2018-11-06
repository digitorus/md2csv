package main

import (
	"bufio"
	"os"
	"regexp"
	"path"
	"sync"
	"fmt"
	"io"
	"log"
	"io/ioutil"
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/gomarkdown/markdown/parser"
	"github.com/gomarkdown/markdown/ast"
)

// Service implements the MarkDown (MD) to CSV conversion
type Service struct {
	doc ast.Node

	wg     sync.WaitGroup
	writer *csv.Writer
	save   chan Row
}

var sectionFormat = regexp.MustCompile(`^([0-9]+\.?)+$`)
var sectionSeperator = regexp.MustCompile(`\s`)
var doubleNewlines = regexp.MustCompile(`(?m)^[\s\n$]+`)

var skipSections = []string{"Introduction", "Scope", "Definitions", "Acronyms",
	"Revisions", "PUBLICATION AND REPOSITORY RESPONSIBILITIES", 
	"Acknowledgements", "References"}

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: " + os.Args[0] + " https://example.com/document.md ...")
		os.Exit(1)
	}

	for _, arg := range os.Args[1:] {
		fmt.Println("Processing:", arg)

		if strings.HasPrefix(arg, "http") {
			err := fromURL(arg)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			file, err := os.Open(arg)
			if err != nil {
				log.Fatal(err)
			}
			parse(arg, bufio.NewReader(file))
		}

	}
}

func fromURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}
	// text/markdown - text/x-markdown - text/plain
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/") {
		return fmt.Errorf("content type must be text/")
	}

	return parse(path.Base(url), resp.Body)
}

func parse(name string, data io.Reader) error {

	s := &Service{}
	go s.newWriter(name)

	markdown, err := ioutil.ReadAll(data)
	if err != nil {
		log.Fatal(err)
	}

	p := parser.New()
	s.doc = p.Parse(markdown)

	fmt.Println("Walking", name)
	ast.Walk(s.doc, &NodeVisitor{s: s})
	fmt.Println("Done walking", name)

	close(s.save)

	fmt.Println("Almost done with", name)
	s.wg.Wait()
	fmt.Println("Done with", name)

	return nil
}


// NodeVisitor implements the ast.NodeVisitor interface
type NodeVisitor struct {
	Row
	index sync.Map
	s *Service
}

// Visit defines what to do for a specific node type
func (n *NodeVisitor) Visit(node ast.Node, entering bool) ast.WalkStatus {
	fmt.Print(".")
	//fmt.Printf("---------------------------- %T, %t\n", node, entering)
	switch node := node.(type) {
	case *ast.Text:
		if n.Category != "" {
			n.Description = n.Description + string(node.Literal)
		}
	case *ast.Softbreak:
		if n.Category != "" {
			n.Description = n.Description + "\n"
		}

	case *ast.Hardbreak:
		if n.Category != "" {
			n.Description = n.Description + "\n\n"
		}

	case *ast.Emph:
	case *ast.Strong:
	case *ast.Del:
	case *ast.BlockQuote:
	case *ast.Aside:
	case *ast.Link:
	case *ast.CrossReference:
	case *ast.Citation:
	case *ast.Image:
	case *ast.Code:
	case *ast.CodeBlock:
	case *ast.Caption:
	case *ast.CaptionFigure:
	case *ast.Document:
		// save remaining data when leaving document
		if !entering {
			n.save()
			fmt.Println()
		}

	case *ast.Paragraph:
		if n.Category != "" {
			n.Description = n.Description + "\n"
		}

	case *ast.HTMLSpan:
	case *ast.HTMLBlock:
	case *ast.Heading:
		if entering {
			// Check if heading has a text child
			var title []string
			titleNode := ast.GetFirstChild(node)
			if txt, ok := titleNode.(*ast.Text); ok {
				title = sectionSeperator.Split(string(txt.Literal), 2)
			} else {
				return ast.GoToNext
			}
			
			if len(title) != 2 {
				return ast.GoToNext
			}
			// first part of title is expected to be a pragraph/section number
			title[0] = strings.TrimSpace(title[0])
			if !sectionFormat.MatchString(title[0]) {
				return ast.GoToNext
			}
			
			n.save()

			n.Name = strings.Trim(title[0], ".") // trimmed above for regex
			n.Title = strings.Trim(strings.TrimSpace(title[1]), ".")
			n.Description = ""

			if strings.HasPrefix(n.Category, n.Name) {
				n.Category = n.Title
			}

			n.index.Store(n.Name, n.Title)
			if strings.Count(n.Name, ".") == 0 {
				n.Category = n.Title
			} else if cat, ok := n.index.Load(n.Name[:strings.LastIndex(n.Name, ".")]); ok {
				n.Category = cat.(string)
			}

			return ast.SkipChildren
		}
		
	case *ast.HorizontalRule:
	case *ast.List:
	case *ast.ListItem:
		if entering {
			n.Description = n.Description + "\t"
		}
	case *ast.Table:
	case *ast.TableCell:
	case *ast.TableHeader:
	case *ast.TableBody:
	case *ast.TableRow:
	case *ast.TableFooter:
	case *ast.Math:
	case *ast.MathBlock:
	case *ast.DocumentMatter:
	case *ast.Callout:
	case *ast.Index:
	case *ast.Subscript:
	case *ast.Superscript:
	case *ast.Footnotes:

	default:
		panic(fmt.Sprintf("Unknown node %T", node))
	}
	return ast.GoToNext
}

func (n *NodeVisitor) save() {
	if n.Description != "" && !inSlice(n.Category, skipSections) &&
		!inSlice(n.Title, skipSections) {

		n.Description = strings.TrimSpace(n.Description)
		n.Description = doubleNewlines.ReplaceAllString(n.Description, "\n")

		n.s.save <- n.Row
	}
}