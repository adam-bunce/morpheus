grammar neo;

@parser::header {
}

constraintList returns [[]backend.Constraint ret]:
{ var listOfConstraints []backend.Constraint}
(constraint { listOfConstraints = append(listOfConstraints, $constraint.ret)}
(',' constraint { listOfConstraints = append(listOfConstraints, $constraint.ret)})*)?
{ $ret = listOfConstraints }
;

constraint returns [backend.Constraint ret]
    : li=ITEM 'is left of' ri=ITEM { $ret = backend.Constraint{LeftItemName: $li.text, RightItemName: $ri.text, ConstraintType: backend.Left}}
    | li=ITEM 'is right of' ri=ITEM { $ret = backend.Constraint{LeftItemName: $li.text, RightItemName: $ri.text, ConstraintType: backend.Right}}
    | li=ITEM 'is below' ri=ITEM { $ret = backend.Constraint{LeftItemName: $li.text, RightItemName: $ri.text, ConstraintType: backend.Below}}
    | li=ITEM 'is above' ri=ITEM{ $ret = backend.Constraint{LeftItemName: $li.text, RightItemName: $ri.text, ConstraintType: backend.Above}}
    ;

fragment LETTER: 'a'..'z' | 'A'..'Z' ;
ITEM: '*'(LETTER | '_')+ ;

WHITESPACE : [ \t\r\n] -> skip;
