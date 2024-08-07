package md

// example for https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html

import (
	// "os"

	"regexp"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"

	"fmt"
)

// get a short name of the type of v which excludes package name
// and strips "()" from the end
func getNodeType(node ast.Node) string {
	s := fmt.Sprintf("%T", node)
	s = strings.TrimSuffix(s, "()")
	if idx := strings.Index(s, "."); idx != -1 {
		return s[idx+1:]
	}
	return s
}

func contentToString(d1 []byte, d2 []byte) string {
	if d1 != nil {
		return string(d1)
	}
	if d2 != nil {
		return string(d2)
	}
	return ""
}

func getContentOrig(node ast.Node) string {
	if c := node.AsContainer(); c != nil {
		return contentToString(c.Literal, c.Content)
	}
	leaf := node.AsLeaf()
	return contentToString(leaf.Literal, leaf.Content)
}

func getContent(node ast.Node) string {
	// nt := getNodeType(node)
	if c := node.AsContainer(); c != nil {
		return contentToString(c.Literal, c.Content)
	}
	leaf := node.AsLeaf()
	return contentToString(leaf.Literal, leaf.Content)
}

func parseBlock(node ast.Node) []SyntaxNode {
	// nt := getNodeType(node)
	// fmt.Println("parseBlock", nt, node)
	if c := node.AsContainer(); c != nil {
		switch v := node.(type) {
		case *ast.Heading:
			return []SyntaxNode{
				&Heading{
					Level:    fmt.Sprintf("%d", v.Level),
					Contents: getContentOrig(c.Children[0]),
				},
			}
		case *ast.Paragraph:
			out_c := ""
			blocks := make([]SyntaxNode, 0)

			// If there's a link by itself, treat it a a 'fancy' link
			// It's similar to what we do with images, but there we let them interrupt the text flow
			// I'm not sure we can distinguish in-line links and standalone links, so we'll just
			// parse out those on their own.
			if len(c.Children) == 2 {
				if _, ok := c.Children[0].(*ast.Text); ok {
					if len(c.Children[0].AsLeaf().Content) == 0 {
						if lz, ok := c.Children[1].(*ast.Link); ok {

							link_text := getContentOrig(lz.Children[0].AsLeaf())
							link_url := lz.Destination

							block := &Link{
								Contents: link_text,
								Url:      string(link_url),
							}
							blocks = append(blocks, block)
							return blocks
						}
					}
				}
			}

			for _, n := range c.Children {
				switch vv := n.(type) {
				case *ast.Text:
					// if it's an empty text node, ignore it.
					if vv.Literal == nil && vv.Content == nil {
						continue
					} else {
						out_c += getContentOrig(n)
					}
				case *ast.Strong:
					out_c += fmt.Sprintf("**%s**", getContentOrig(n.AsContainer().Children[0]))
				case *ast.Emph:
					out_c += fmt.Sprintf("_%s_", getContentOrig(n.AsContainer().Children[0]))
				case *ast.Link:
					out_c += fmt.Sprintf("%s", vv.Destination) //fmt.Sprintf("[%s](%s)", getContentOrig(n.AsContainer().Children[0]), vv.Destination)
				case *ast.Image:
					// If there's existing contents, flush it to a block
					if out_c != "" {
						blocks = append(blocks, &Paragraph{
							Contents: strings.TrimSpace(out_c),
						})
						// start afresh
						out_c = ""
					}
					// Add the image block
					blocks = append(blocks, &Image{
						Url:     string(vv.Destination),
						AltText: getContentOrig(n.AsContainer().Children[0]),
					})
					// out_c += fmt.Sprintf("![%s](%s)", getContentOrig(n.AsContainer().Children[0]), vv.Destination)
				default:
					panic(fmt.Sprintf("Node type not handled, error. %T\n", vv))
					out_c += getContentOrig(n)
					// pass
				}
			}
			if out_c != "" {
				blocks = append(blocks, &Paragraph{
					Contents: strings.TrimSpace(out_c),
				})
			}
			return blocks
			// return []*SyntaxNode{
			// 	&SyntaxNode{
			// 		Type:     P,
			// 		Contents: out_c,
			// 	},
			// }
		case *ast.List:
			// Get list type
			list_contents := []string{}
			for _, n := range c.Children {
				kid := n.AsContainer().Children[0].AsContainer().Children[0]
				list_contents = append(list_contents, getContentOrig(kid))
			}
			return []SyntaxNode{
				&List{
					Contents: list_contents,
					Ordered:  v.ListFlags&ast.ListTypeOrdered != 0,
				},
			}
		case *ast.Table:
			// need header, body
			// header
			header := []string{}
			for _, row := range c.Children[0].AsContainer().Children {
				for _, cell := range row.AsContainer().Children {
					header_actual := cell.AsContainer().Children[0]
					header = append(header, getContentOrig(header_actual))
				}
			}
			// body
			body := [][]string{}
			for _, row := range c.Children[1].AsContainer().Children {
				row_contents := []string{}
				for _, cell := range row.AsContainer().Children {
					cell_actual := cell.AsContainer().Children[0]
					fmt.Println("cell_actual", cell_actual, getContentOrig(cell_actual))
					row_contents = append(row_contents, getContentOrig(cell_actual))
				}
				body = append(body, row_contents)
			}
			return []SyntaxNode{
				&Table{
					Header: header,
					Body:   body,
				},
			}
		default:
			panic(fmt.Sprintf("Unhandled container type %T", v))
		}
	} else {
		switch v := node.(type) {
		case *ast.CodeBlock:
			_ = v
			table_marker := regexp.MustCompile(`^table_view\|?(.*)$`)
			if table_marker.MatchString(string(v.Info)) {
				return []SyntaxNode{
					&TableView{
						Display: table_marker.FindStringSubmatch(string(v.Info))[1],
						Query:   strings.TrimSpace(getContentOrig(node)),
					},
				}
			} else {
				return []SyntaxNode{
					&Code{
						Contents: strings.TrimSpace(getContentOrig(node)),
						Lang:     string(v.Info),
					},
				}
			}
		case *ast.HorizontalRule:
			return []SyntaxNode{
				&HorizontalRule{},
			}
		default:
			panic(fmt.Sprintf("Unhandled leaf type %T", v))
		}
	}
}

func MdToBlocks(md []byte) []SyntaxNode {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)
	//
	// fmt.Println("Official AST")
	// fmt.Println("vvvvvvvvvv")
	// ast.Print(os.Stdout, doc)
	// fmt.Println("^^^^^^^^^^")

	out := []SyntaxNode{}

	for _, node := range doc.GetChildren() {

		// content := getContent(node)
		// typeName := getNodeType(node)
		blocks := parseBlock(node)
		out = append(out, blocks...)
		// fmt.Println(typeName, content)
	}
	return out
}

// func main() {
// 	md := []byte(mds)
// 	mdToHTML(md)
// }
