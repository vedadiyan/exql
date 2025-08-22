%{
package lang

%}

%union {
    expr     ExprNode
    exprList []ExprNode
    str      string
    num      float64
    boolean  bool
}

%token <str> IDENTIFIER STRING DSTRING
%token <num> NUMBER
%token <boolean> BOOLEAN

%token AND OR NOT IN
%token EQ NE LT LE GT GE
%token LPAREN RPAREN LBRACKET RBRACKET
%token DOT COMMA QUOTE DQUOTE

%type <expr> expr logical_expr equality_expr relational_expr additive_expr multiplicative_expr unary_expr primary_expr
%type <expr> field_access function_call list_literal
%type <exprList> argument_list expression_list

%left OR
%left AND
%left IN
%left EQ NE
%left LT LE GT GE
%left '+' '-'
%left '*' '/'
%right NOT UMINUS
%left DOT LBRACKET

%%

program: expr { yylex.(*yyLex).result = $1 }

expr: logical_expr { $$ = $1 }

logical_expr: logical_expr AND equality_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "and"}
    }
    | logical_expr OR equality_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "or"}
    }
    | equality_expr { $$ = $1 }

equality_expr: equality_expr EQ relational_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "="}
    }
    | equality_expr NE relational_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "!="}
    }
    | relational_expr { $$ = $1 }

relational_expr: relational_expr LT additive_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "<"}
    }
    | relational_expr LE additive_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "<="}
    }
    | relational_expr GT additive_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: ">"}
    }
    | relational_expr GE additive_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: ">="}
    }
    | relational_expr IN additive_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "in"}
    }
    | relational_expr NOT IN additive_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $4, Operator: "not in"}
    }
    | additive_expr { $$ = $1 }

additive_expr: additive_expr '+' multiplicative_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "+"}
    }
    | additive_expr '-' multiplicative_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "-"}
    }
    | multiplicative_expr { $$ = $1 }

multiplicative_expr: multiplicative_expr '*' unary_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "*"}
    }
    | multiplicative_expr '/' unary_expr {
        $$ = &BinaryOpNode{Left: $1, Right: $3, Operator: "/"}
    }
    | unary_expr { $$ = $1 }

unary_expr: NOT unary_expr {
        $$ = &UnaryOpNode{Operand: $2, Operator: "not"}
    }
    | '-' unary_expr %prec UMINUS {
        $$ = &UnaryOpNode{Operand: $2, Operator: "-"}
    }
    | primary_expr { $$ = $1 }

primary_expr: IDENTIFIER { 
        $$ = &VariableNode{Name: $1}
    }
    | NUMBER {
        $$ = &LiteralNode{Value: NumberValue($1)}
    }
    | STRING {
        $$ = &LiteralNode{Value: StringValue($1)}
    }
    | DSTRING {
        $$ = &LiteralNode{Value: StringValue($1)}
    }
    | BOOLEAN {
        $$ = &LiteralNode{Value: BoolValue($1)}
    }
    | LPAREN expr RPAREN {
        $$ = $2
    }
    | field_access { $$ = $1 }
    | function_call { $$ = $1 }
    | list_literal { $$ = $1 }

field_access: primary_expr DOT IDENTIFIER {
        $$ = &FieldAccessNode{Object: $1, Field: $3}
    }
    | primary_expr LBRACKET expr RBRACKET {
        $$ = &IndexAccessNode{Object: $1, Index: $3}
    }

function_call: IDENTIFIER LPAREN argument_list RPAREN {
        $$ = &FunctionCallNode{Name: $1, Args: $3}
    }
    | IDENTIFIER LPAREN RPAREN {
        $$ = &FunctionCallNode{Name: $1, Args: []ExprNode{}}
    }

list_literal: LBRACKET expression_list RBRACKET {
        $$ = &ListNode{Elements: $2}
    }
    | LBRACKET RBRACKET {
        $$ = &ListNode{Elements: []ExprNode{}}
    }

argument_list: expr {
        $$ = []ExprNode{$1}
    }
    | argument_list COMMA expr {
        $$ = append($1, $3)
    }

expression_list: expr {
        $$ = []ExprNode{$1}
    }
    | expression_list COMMA expr {
        $$ = append($1, $3)
    }

%%