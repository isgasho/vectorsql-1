// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package planners

import (
	"testing"

	"parsers"

	"parsers/sqlparser"

	"github.com/stretchr/testify/assert"
)

func TestSelectPlan(t *testing.T) {
	tests := []struct {
		err    error
		name   string
		query  string
		expect string
	}{
		{
			name:  "simple",
			query: "select * from t1",
			expect: `{
    "Name": "SelectPlan",
    "SubPlan": {
        "Name": "MapPlan",
        "SubPlans": [
            {
                "Name": "ScanPlan",
                "Table": "t1",
                "Schema": ""
            },
            {
                "Name": "ProjectPlan",
                "SubPlan": {
                    "Name": "MapPlan"
                }
            },
            {
                "Name": "SinkPlan"
            }
        ]
    }
}`,
		},
		{
			name:  "filter",
			query: "select * from t1 where (id>1 and id<8) or id=9 and c>2019.12 or d='name'",
			expect: `{
    "Name": "SelectPlan",
    "SubPlan": {
        "Name": "MapPlan",
        "SubPlans": [
            {
                "Name": "ScanPlan",
                "Table": "t1",
                "Schema": ""
            },
            {
                "Name": "ProjectPlan",
                "SubPlan": {
                    "Name": "MapPlan"
                }
            },
            {
                "Name": "FilterPlan",
                "SubPlan": {
                    "Name": "BinaryExpressionPlan",
                    "FuncName": "OR",
                    "Left": {
                        "Name": "BinaryExpressionPlan",
                        "FuncName": "OR",
                        "Left": {
                            "Name": "BinaryExpressionPlan",
                            "FuncName": "AND",
                            "Left": {
                                "Name": "BinaryExpressionPlan",
                                "FuncName": "\u003e",
                                "Left": {
                                    "Name": "VariablePlan",
                                    "Value": "id"
                                },
                                "Right": {
                                    "Name": "ConstantPlan",
                                    "Value": 1
                                }
                            },
                            "Right": {
                                "Name": "BinaryExpressionPlan",
                                "FuncName": "\u003c",
                                "Left": {
                                    "Name": "VariablePlan",
                                    "Value": "id"
                                },
                                "Right": {
                                    "Name": "ConstantPlan",
                                    "Value": 8
                                }
                            }
                        },
                        "Right": {
                            "Name": "BinaryExpressionPlan",
                            "FuncName": "AND",
                            "Left": {
                                "Name": "BinaryExpressionPlan",
                                "FuncName": "=",
                                "Left": {
                                    "Name": "VariablePlan",
                                    "Value": "id"
                                },
                                "Right": {
                                    "Name": "ConstantPlan",
                                    "Value": 9
                                }
                            },
                            "Right": {
                                "Name": "BinaryExpressionPlan",
                                "FuncName": "\u003e",
                                "Left": {
                                    "Name": "VariablePlan",
                                    "Value": "c"
                                },
                                "Right": {
                                    "Name": "ConstantPlan",
                                    "Value": 2019.12
                                }
                            }
                        }
                    },
                    "Right": {
                        "Name": "BinaryExpressionPlan",
                        "FuncName": "=",
                        "Left": {
                            "Name": "VariablePlan",
                            "Value": "d"
                        },
                        "Right": {
                            "Name": "ConstantPlan",
                            "Value": "name"
                        }
                    }
                }
            },
            {
                "Name": "SinkPlan"
            }
        ]
    }
}`,
		},
		{
			name:  "orderby",
			query: "select * from t1 order by c1 desc, c2 asc, c3 desc",
			expect: `{
    "Name": "SelectPlan",
    "SubPlan": {
        "Name": "MapPlan",
        "SubPlans": [
            {
                "Name": "ScanPlan",
                "Table": "t1",
                "Schema": ""
            },
            {
                "Name": "ProjectPlan",
                "SubPlan": {
                    "Name": "MapPlan"
                }
            },
            {
                "Name": "OrderByPlan",
                "Orders": [
                    {
                        "Expression": {
                            "Name": "VariablePlan",
                            "Value": "c1"
                        },
                        "Direction": "desc"
                    },
                    {
                        "Expression": {
                            "Name": "VariablePlan",
                            "Value": "c2"
                        },
                        "Direction": "asc"
                    },
                    {
                        "Expression": {
                            "Name": "VariablePlan",
                            "Value": "c3"
                        },
                        "Direction": "desc"
                    }
                ]
            },
            {
                "Name": "SinkPlan"
            }
        ]
    }
}`,
		},
		{
			name:  "tvf-range",
			query: "SELECT * FROM range(range_start -> 1, range_end -> 5) r",
			expect: `{
    "Name": "SelectPlan",
    "SubPlan": {
        "Name": "MapPlan",
        "SubPlans": [
            {
                "Name": "TableValuedFunctionPlan",
                "As": "",
                "FuncName": "range",
                "SubPlan": {
                    "Name": "MapPlan",
                    "SubPlans": [
                        {
                            "Name": "TableValuedFunctionExpressionPlan",
                            "FuncName": "",
                            "SubPlan": {
                                "Name": "BinaryExpressionPlan",
                                "FuncName": "-\u003e",
                                "Left": {
                                    "Name": "VariablePlan",
                                    "Value": "range_start"
                                },
                                "Right": {
                                    "Name": "ConstantPlan",
                                    "Value": 1
                                }
                            }
                        },
                        {
                            "Name": "TableValuedFunctionExpressionPlan",
                            "FuncName": "",
                            "SubPlan": {
                                "Name": "BinaryExpressionPlan",
                                "FuncName": "-\u003e",
                                "Left": {
                                    "Name": "VariablePlan",
                                    "Value": "range_end"
                                },
                                "Right": {
                                    "Name": "ConstantPlan",
                                    "Value": 5
                                }
                            }
                        }
                    ]
                }
            },
            {
                "Name": "ProjectPlan",
                "SubPlan": {
                    "Name": "MapPlan"
                }
            },
            {
                "Name": "SinkPlan"
            }
        ]
    }
}`,
		},
		{
			name:  "select-aggregate",
			query: "SELECT max(a), sum(b), c, (id+1) FROM t1 where (id+1)!=2 group by d,e order by c desc",
			expect: `{
    "Name": "SelectPlan",
    "SubPlan": {
        "Name": "MapPlan",
        "SubPlans": [
            {
                "Name": "ScanPlan",
                "Table": "t1",
                "Schema": ""
            },
            {
                "Name": "ProjectPlan",
                "SubPlan": {
                    "Name": "MapPlan",
                    "SubPlans": [
                        {
                            "Name": "UnaryExpressionPlan",
                            "FuncName": "MAX",
                            "Expr": {
                                "Name": "VariablePlan",
                                "Value": "a"
                            }
                        },
                        {
                            "Name": "UnaryExpressionPlan",
                            "FuncName": "SUM",
                            "Expr": {
                                "Name": "VariablePlan",
                                "Value": "b"
                            }
                        },
                        {
                            "Name": "VariablePlan",
                            "Value": "c"
                        },
                        {
                            "Name": "BinaryExpressionPlan",
                            "FuncName": "+",
                            "Left": {
                                "Name": "VariablePlan",
                                "Value": "id"
                            },
                            "Right": {
                                "Name": "ConstantPlan",
                                "Value": 1
                            }
                        }
                    ]
                }
            },
            {
                "Name": "FilterPlan",
                "SubPlan": {
                    "Name": "BinaryExpressionPlan",
                    "FuncName": "!=",
                    "Left": {
                        "Name": "BinaryExpressionPlan",
                        "FuncName": "+",
                        "Left": {
                            "Name": "VariablePlan",
                            "Value": "id"
                        },
                        "Right": {
                            "Name": "ConstantPlan",
                            "Value": 1
                        }
                    },
                    "Right": {
                        "Name": "ConstantPlan",
                        "Value": 2
                    }
                }
            },
            {
                "Name": "OrderByPlan",
                "Orders": [
                    {
                        "Expression": {
                            "Name": "VariablePlan",
                            "Value": "c"
                        },
                        "Direction": "desc"
                    }
                ]
            },
            {
                "Name": "SinkPlan"
            }
        ]
    }
}`,
		},
		{
			name:  "aliased-select-test",
			query: "SELECT max(a+1), (id+1) as b, c as c1 FROM t1 where b>5",
			expect: `{
    "Name": "SelectPlan",
    "SubPlan": {
        "Name": "MapPlan",
        "SubPlans": [
            {
                "Name": "ScanPlan",
                "Table": "t1",
                "Schema": ""
            },
            {
                "Name": "ProjectPlan",
                "SubPlan": {
                    "Name": "MapPlan",
                    "SubPlans": [
                        {
                            "Name": "UnaryExpressionPlan",
                            "FuncName": "MAX",
                            "Expr": {
                                "Name": "BinaryExpressionPlan",
                                "FuncName": "+",
                                "Left": {
                                    "Name": "VariablePlan",
                                    "Value": "a"
                                },
                                "Right": {
                                    "Name": "ConstantPlan",
                                    "Value": 1
                                }
                            }
                        },
                        {
                            "Name": "AliasedExpressionPlan",
                            "As": "b",
                            "Expr": {
                                "Name": "BinaryExpressionPlan",
                                "FuncName": "+",
                                "Left": {
                                    "Name": "VariablePlan",
                                    "Value": "id"
                                },
                                "Right": {
                                    "Name": "ConstantPlan",
                                    "Value": 1
                                }
                            }
                        },
                        {
                            "Name": "AliasedExpressionPlan",
                            "As": "c1",
                            "Expr": {
                                "Name": "VariablePlan",
                                "Value": "c"
                            }
                        }
                    ]
                }
            },
            {
                "Name": "FilterPlan",
                "SubPlan": {
                    "Name": "BinaryExpressionPlan",
                    "FuncName": "\u003e",
                    "Left": {
                        "Name": "BinaryExpressionPlan",
                        "FuncName": "+",
                        "Left": {
                            "Name": "VariablePlan",
                            "Value": "id"
                        },
                        "Right": {
                            "Name": "ConstantPlan",
                            "Value": 1
                        }
                    },
                    "Right": {
                        "Name": "ConstantPlan",
                        "Value": 5
                    }
                }
            },
            {
                "Name": "SinkPlan"
            }
        ]
    }
}`,
		},
	}

	for _, test := range tests {
		statement, err := parsers.Parse(test.query)
		assert.Nil(t, err)

		plan := NewSelectPlan(statement.(*sqlparser.Select))
		err = plan.Build()
		assert.Nil(t, err)

		expect := test.expect
		actual := plan.String()
		assert.Equal(t, expect, actual)
	}
}
