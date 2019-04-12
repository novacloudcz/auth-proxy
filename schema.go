package main

import (
	"fmt"

	// parser "github.com/bucketd/go-graphqlparser/graphql"
	// "github.com/bucketd/go-graphqlparser/graphql/types"
	// "github.com/graphql-go/graphql"
	"github.com/davecgh/go-spew/spew"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

// var rawSchema = []byte(`
// type Query {
//   bar: String
//   baz: String
// }
// `)

var rawQuery = []byte(`
query {
  bar
  baz
}
`)

func testGQL() {
	astDoc, _ := parser.Parse(parser.ParseParams{
		Source: string(rawQuery),
		Options: parser.ParseOptions{
			NoLocation: true,
		},
	})
	// spew.Dump(astDoc, err)

	// printed := printer.Print(astDoc)
	// fmt.Println("??", printed)

	// fmt.Println("?!!", reflect.DeepEqual(string(rawQuery), printed))

	for _, def := range astDoc.Definitions {
		spew.Dump(def)
		op, ok := def.(*ast.OperationDefinition)
		fmt.Println("?", op, ok)
		if ok {
			for _, sel := range op.SelectionSet.Selections {
				fmt.Println("?", sel)
			}
		}
	}

	// start := time.Now()
	// schema, errs, err := parser.BuildSchema(nil, rawSchema)
	// if err != nil {
	// 	panic(err)
	// }

	// if errs.Len() > 0 {
	// 	fmt.Println("Failed to validate schema")
	// 	errs.ForEach(func(e types.Error, i int) {
	// 		fmt.Println(e.Message)
	// 	})

	// 	os.Exit(1)
	// }
	// // Schema
	// fields := graphql.Fields{
	// 	"hello": &graphql.Field{
	// 		Type: graphql.String,
	// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
	// 			return "world", nil
	// 		},
	// 	},
	// }
	// rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	// schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	// schema, err := graphql.NewSchema(schemaConfig)
	// if err != nil {
	// 	log.Fatalf("failed to create new schema, error: %v", err)
	// }

	// Query
	// query := `
	// 	{
	// 		hello
	// 	}
	// `
	// params := graphql.Params{Schema: schema, RequestString: query}
	// r := graphql.Do(params)
	// if len(r.Errors) > 0 {
	// 	log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	// }
	// rJSON, _ := json.Marshal(r)
	// fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}

	// elapsed := time.Since(start)
	// log.Printf("took %s", elapsed)
}
