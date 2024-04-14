grammar morpheus;
import neo;

@parser::header {
    import "github.com/adam-bunce/morpheus/backend"

}

program returns [backend.Block statements]
    :
      { var listOfExpressions []backend.Expression }
      (statement { listOfExpressions = append(listOfExpressions, $statement.expression); })*
      { $statements = backend.Block{Exprs: listOfExpressions}; }
    ;


statement returns [backend.Expression expression]
    : assignment SEMICOLON? { $expression = $assignment.expression }
    | expr SEMICOLON? { $expression = $expr.expression; }
    | loop { $expression = $loop.expression }
    | funDef { $expression = $funDef.expression }
    | ifElse { $expression = $ifElse.expression}
    | builtIn SEMICOLON? { $expression = $builtIn.expression }
    ;

assignment returns [backend.Assign expression]
    : 'let'? ID ASSIGN expr { $expression = backend.Assign{Name: $ID.text, Expr: $expr.expression} } // bind expr to ID
    ;

expr returns [backend.Expression expression]
    : LPAREN expr RPAREN { $expression = $expr.expression }
    | e1=expr PLUS e2=expr { $expression = backend.Arithmetic{Left: $e1.expression, Right: $e2.expression, Op: backend.ADD} }
    | e1=expr SUBTRACT e2=expr { $expression = backend.Arithmetic{Left: $e1.expression, Right: $e2.expression, Op: backend.SUB} }
    | e1=expr ASTERISK e2=expr { $expression = backend.Arithmetic{Left: $e1.expression, Right: $e2.expression, Op: backend.MUL} }
    | e1=expr SLASH e2=expr { $expression = backend.Arithmetic{Left: $e1.expression, Right: $e2.expression, Op: backend.DIV} }
    | e1=expr PLUS PLUS e2=expr{ $expression = backend.Concat{Left: $e1.expression, Right: $e2.expression} } // str concat
    | ID LPAREN al=argList RPAREN { $expression = backend.FunctionCall{Name: $ID.text, Args: $al.expressionList} } // func call
    | LPAREN compare RPAREN { $expression = $compare.expression } // evals to boolean
    | ID { $expression = backend.Dereference{ Name: $ID.text } } // derefrence var
    | list { $expression = $list.expression }
    | listOp { $expression = $listOp.expression }
    | NUMBER { $expression = backend.NewIntLiteral($NUMBER.text) }
    | STRING { $expression = backend.NewStringLiteral($STRING.text) }
    | BOOLEAN { $expression = backend.NewBooleanLiteral($BOOLEAN.text) }

    | 'Box' LPAREN STRING RPAREN{ $expression = backend.BoxExpr{Id: $STRING.text} }
    | 'Group' LPAREN list COLON LSQBRACE cl=constraintList RSQBRACE RPAREN
         { $expression = backend.GroupExpr{Items: $list.expression, Constraints: $cl.ret} }
    ;

block returns [backend.Block expression]
    :  {var blockExprs []backend.Expression}
       ( statement { blockExprs = append(blockExprs, $statement.expression); } )*
       { $expression = backend.Block{Exprs: blockExprs}; }
    ;

list returns [backend.Expression expression]
    : { var exprList []backend.Expression }
      LSQBRACE (e1=expr { exprList = append(exprList, $e1.expression)} (COMMA e2=expr { exprList = append(exprList, $e2.expression) })*)? RSQBRACE
      { $expression = backend.List{Values: exprList} }
    ;

listOrId returns [backend.Expression expression]
    : list { $expression = $list.expression }
    | ID { $expression = backend.Dereference{Name: $ID.text } }
    ;

// NOTE: you can't chain these together, so no [1,2,3].del(3).add(3).len
listOp returns [backend.Expression expression]
    : listOrId '.get' LPAREN e1=expr RPAREN{ $expression = backend.ListIndex{List: $listOrId.expression, Position: $e1.expression} }
    | listOrId '.add' LPAREN e1=expr RPAREN{ $expression = backend.ListAdd{List: $listOrId.expression, Value: $e1.expression} }
    | listOrId '.del' LPAREN e1=expr RPAREN{ $expression = backend.ListDelete{List: $listOrId.expression, Position: $e1.expression} }
    | listOrId '.len'  { $expression = backend.ListLength{List: $listOrId.expression} }
    ;

loop returns [backend.Expression expression]
    : 'for' ID 'in'  LPAREN e1=expr COMMA e2=expr COMMA e3=expr RPAREN LBRACE
        body=block
      RBRACE
      { $expression = backend.Loop{ Iterator: $ID.text,
                                    Start:    $e1.expression,
                                    Stop:     $e2.expression,
                                    Step:     $e3.expression,
                                    Body:     $body.expression } }
    ;


funDef returns [backend.Expression expression]
    : 'function' ID LPAREN paramList RPAREN LBRACE
        block
      RBRACE
      { $expression = backend.Declare{Name: $ID.text, Args: $paramList.params, Body: $block.expression} }
    ;

paramList returns [[]string params]
    :
    { var parameterList []string }
    (id1=ID { parameterList = append(parameterList, $id1.text)}
      (COMMA idn=ID {parameterList = append(parameterList, $idn.text)} )*)?
    { $params = parameterList}
    ;

argList returns [[]backend.Expression expressionList]
    :
    { var exprs []backend.Expression }
      (e1=expr        { exprs = append(exprs, $e1.expression) }
      (COMMA e2=expr { exprs = append(exprs, $e2.expression) } ))*
     { $expressionList = exprs }
    ;


compare returns [backend.Expression expression]
    : e1=expr LT e2=expr { $expression = backend.Compare{Left: $e1.expression, Right: $e2.expression, Op: backend.LT} }
    | e1=expr GT e2=expr { $expression = backend.Compare{Left: $e1.expression, Right: $e2.expression, Op: backend.GT} }
    | e1=expr ASSIGN ASSIGN e2=expr { $expression = backend.Compare{Left: $e1.expression, Right: $e2.expression, Op: backend.EQ} }
    | e1=expr 'and' e2=expr { $expression = backend.Compare{Left: $e1.expression, Right: $e2.expression, Op: backend.AND} }
    | e1=expr 'or' e2=expr { $expression = backend.Compare{Left: $e1.expression, Right: $e2.expression, Op: backend.OR} }
    ;

ifElse returns [backend.Expression expression]
    :
         { var elifConds []backend.Conditional }
         { var elseExpr backend.Expression }
        'if' LPAREN ifComparison=compare RPAREN LBRACE
            ifBlock=block
        RBRACE

        ( 'elif' LPAREN elifComparison=compare RPAREN LBRACE
            elifBlock=block
        RBRACE {
         elifConds = append(elifConds, backend.Conditional{Condition: $elifComparison.expression, Body: $elifBlock.expression})
         } )*


        ( 'else' LBRACE
            elseBlock=block
        RBRACE )?
        {
            if $elseBlock.text != "" {
                elseExpr = $elseBlock.expression
            }
        }

        { $expression = backend.IfElifElse{
            If: backend.Conditional{Condition: $ifComparison.expression, Body: $ifBlock.expression },
            ElseIf: elifConds,
            Else: elseExpr,
        } }

    ;

builtIn returns [backend.Expression expression]
    : 'print' LPAREN expr RPAREN { $expression = backend.Print{ToPrint: $expr.expression} }
    | expr '.htmlify'LPAREN STRING RPAREN { $expression = backend.Htmlify{Layout: $expr.expression, File: $STRING.text}}
    ;



// Lexer Rules
fragment DIGIT: '0' .. '9' ;
fragment LETTER: 'a'..'z' | 'A'..'Z' ;

COMMENT: '//' ~('\r' | '\n')* -> skip;

LPAREN: '(' ;
RPAREN: ')' ;
SLASH: '/' ;
PLUS: '+' ;
SUBTRACT: '-' ;
ASTERISK: '*' ;
COMMA: ',' ;
ASSIGN: '=' ;
LBRACE: '{' ;
RBRACE: '}' ;
LSQBRACE: '[' ;
RSQBRACE: ']' ;
LT: '<' ;
GT: '>' ;

NUMBER: '-'?DIGIT+('.'DIGIT+)? ;
STRING: '"' ~('"')+ '"' ;
BOOLEAN: 'true' | 'false' ;

SEMICOLON: ';' ;
COLON: ':';
WHITESPACE : [ \t\r\n] -> skip;
ID: (LETTER | '_')+ ;
