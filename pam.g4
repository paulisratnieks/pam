grammar pam;
progr 		:   series;
series		: 	stmt ';' NEWLINE* (stmt ';' NEWLINE*)*;
stmt		:	input_stmt | output_stmt | assign_stmt | cond_stmt | loop_stmt;

input_stmt	:	'read' varlist;
output_stmt	:	'write' varlist;
assign_stmt	:	VARNAME ':=' expr;
cond_stmt	:	'if' logical 'then' series ('else' series)? 'fi';
loop_stmt	:	'while' logical 'do' series 'end';
varlist		:	VARNAME (',' VARNAME)*;

logical     : logical_or;
logical_or	: logical_and ('or' logical_and)*;
logical_and	: logical_not ('and' logical_not)*;
logical_not : 'not' logical | logical_relation;
logical_relation: expr RELATION expr | logical_elem;
logical_elem: LCONST | '(' logical ')';

expr		:	term (WEAKOP term)*;
term		:	elem (STRONGOP elem)*;
elem		:	NUMBER | VARNAME | '(' expr ')';
LCONST      :   'true' | 'false';
NEWLINE	    :	'\r' ? '\n';
WEAKOP		:	'+' | '-';
STRONGOP	:	'*' | '/';
RELATION	:	'<>' | '<=' | '>='| '=' |  '<' | '>';
NUMBER		:	[0-9][0-9]* ;
VARNAME 	:	([a-z]|[A-Z])([a-z]|[A-Z]|[0-9])*;
WS          :   [ \n\t\r]+ -> skip;
