package schemainit

import "github.com/graphql-go/graphql"

var caseUpdateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CaseUpdate",
	Fields: graphql.Fields{
		"update_date": &graphql.Field{
			Type: graphql.DateTime,
		},
		"update_details": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var courtCaseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CourtCase",
	Fields: graphql.Fields{
        "cnj": &graphql.Field{
            Type: graphql.String,
        },
        "plaintiff": &graphql.Field{
            Type: graphql.String,
        },
        "defendant": &graphql.Field{
            Type: graphql.String,
        },
        "court_of_origin": &graphql.Field{
            Type: graphql.String,
        },
        "start_date": &graphql.Field{
            Type: graphql.DateTime,
        },
		"updates": &graphql.Field{
			Type: graphql.NewList(caseUpdateType),
		},
    },
})

func SchemaInit(
	resolver func(p graphql.ResolveParams) (interface{}, error),
) (*graphql.Schema) {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"court_case": &graphql.Field{
				Type: courtCaseType,
				Args: graphql.FieldConfigArgument{
					"cnj": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: resolver,
			},
		},
	})

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	return &schema
}
