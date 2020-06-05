// Copyright 2015 Osipov Konstantin <k.osipov.msk@gmail.com>. All rights reserved.
// license that can be found in the LICENSE file.

// This file is part of the application source code leveldb-cli
// This software provides a console interface to leveldb.

package commands

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/liderman/leveldb-cli/cliutil"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type params struct {
	filter    string
	keyformat string
	format    string
	start     string
	end       string
	prefix    string
	contain   string
	tofile    string
}

func ParseSubCommand(args []string, para *params) bool {
	flags := flag.NewFlagSet("show sub commands", flag.ContinueOnError)
	fmt.Println("parsing:", args)
	flags.StringVar(&para.filter, "filter", "all", "choose filter. avaliable arguments- all,range,prefix,contain")
	flags.StringVar(&para.keyformat, "keyformat", "raw", "choose keys display format")
	flags.StringVar(&para.format, "format", "raw", "choose values display format")
	flags.StringVar(&para.start, "start", "", "key prefix indicating the beginning of the range")
	flags.StringVar(&para.end, "end", "", "key prefix indicating the end of the range")
	flags.StringVar(&para.prefix, "prefix", "", "prefix")
	flags.StringVar(&para.contain, "contain", "", "contain")
	flags.StringVar(&para.tofile, "tofile", "", "output to a file")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return false
	}
	return true
}

func file_write(context, file string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(context)
}

func Show(args []string) string {
	content := ""
	if !isConnected {
		content = AppError(ErrDbDoesNotOpen)
	}
	para := &params{}
	if ParseSubCommand(args, para) {
		if para.filter == "range" {
			content = showByIterator(para, dbh.NewIterator(&util.Range{Start: []byte(para.start), Limit: []byte(para.end)}, nil))
		} else if para.filter == "prefix" {
			content = showByIterator(para, dbh.NewIterator(util.BytesPrefix([]byte(para.prefix)), nil))
		} else {
			content = showByIterator(para, dbh.NewIterator(nil, nil))
		}
	}
	if len(para.tofile) > 0 {
		file_write(content, para.tofile)
		return ""
	}
	return content
}

// Show by iterator
//
// Returns a string containing information about the result of the operation.
func showByIterator(para *params, iter iterator.Iterator) string {
	if iter.Error() != nil {
		return "Empty result!"
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	w := new(tabwriter.Writer)

	w.Init(writer, 0, 8, 0, '\t', 0)
	fmt.Println(para)
	fmt.Fprintln(w, "<Key\t| Value>")
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		keystr := cliutil.ToString(para.keyformat, key)
		if strings.Contains(keystr, para.contain) {
			// fmt.Fprintf(w, "%s\t| %s\n", string(key), cliutil.ToString(format, value))
			fmt.Fprintf(w, "%s\t| %s\n", keystr, cliutil.ToString(para.format, value))
		}
	}

	w.Flush()

	iter.Release()
	err := iter.Error()
	if err != nil {
		return "Error iterator!"
	}

	writer.Flush()
	return string(b.Bytes())
}
