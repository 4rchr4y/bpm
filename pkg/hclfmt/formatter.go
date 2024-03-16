package hclfmt

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Context struct {
	Level uint
}

func FormatFile(f *hclwrite.File) error {
	fmtctx := Context{
		Level: 0,
	}

	for _, block := range f.Body().Blocks() {
		FormatBlock(fmtctx, block)
	}
	return nil
}

func FormatBlock(fmtctx Context, block *hclwrite.Block) {
	for _, b := range block.Body().Blocks() {
		fmt.Println(b.Labels())
	}
}

// func FormatBody(body *hclwrite.Body) {
// 	body.
// }
