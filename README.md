# Abgabe_MSE
Mini-Compiler Projekt für Modelbasierte Software Entwicklung

Kurzbeschreibung

  Ein gegebener Input wird anhand der Regeln einer einfachen imperativen Programmiersprache geparsed, überprüft und anschließend interpretiert.
  
  Zuerst wird versucht den Input zu parsen, anschließend wird auf dem Ergebnis ein Typ-Check durchgeführt, falls dieser erfolgreich ist, wird das Ergebnis des Parsens evaluiert.

Einfache imperative Programmiersprache / IMP [^1]
  
  Syntax
    
    vars      Variable names, start with lower-case letter

    prog      ::= block
    block     ::= "{" statement "}"
    statement ::=  statement ";" statement           -- Command sequence
                |  vars ":=" exp                     -- Variable declaration
                |  vars "=" exp                      -- Variable assignment
                |  "while" exp block                 -- While
                |  "if" exp block "else" block       -- If-then-else
                |  "print" exp                       -- Print

    exp       ::= 0 | 1 | -1 | ...     -- Integers
                | "true" | "false"      -- Booleans
                | exp "+" exp           -- Addition
                | exp "*" exp           -- Multiplication
                | exp "||" exp          -- Disjunction
                | exp "&&" exp          -- Conjunction
                | "!" exp               -- Negation
                | exp "==" exp          -- Equality test
                | exp "<" exp           -- Lesser test
                | "(" exp ")"           -- Grouping of expressions
                | vars                  -- Variables
                
  Static Semantics used for type checker
  
    Types
    
      T ::= int | bool
    
    Variable environment
    
      G ::= [] | [x : T] | G ++ G
  
    Typing rules
    
      Expressions G |- e : T
    
        i some number
        ------------
        G |- i : int

        G |- true : bool

        G |- false : bool
      
        lookup(G,x) = T
        ------------------
        G |- x : T
      
        G |- e1 : int    G |- e2 : int
        ----------------------------------------
        G |- e2 + e2 : int

        G |- e1 : int    G |- e2 : int
        ----------------------------------------
        G |- e2 * e2 : int

        G |- e1 : bool    G |- e2 : bool
        ----------------------------------------
        G |- e2 || e2 : bool

        G |- e1 : bool    G |- e2 : bool
        ----------------------------------------
        G |- e2 && e2 : bool

        G |- e : bool
        ----------------------------------------
        G |- ! e : bool

        G |- e1 : T   G |- e2 : T
        ----------------------------------------
        G |- e1 == e2 : bool

        G |- e1 : int   G |- e2 : int
        ----------------------------------------
        G |- e1 < e2 : int
       
      Statements G |- (s,G2)
    
        G1 |- (s1, G2)   G2 |- (s2,G3)
        ----------------------------------------
        G1 |- (s1 ; s2, G3)

        G |- e : T  G2 = G ++ [x : T]
        ----------------------------------------
        G |- (x := e, G2)

        G |- x : T   G |- e : T
        ----------------------------------------
        G |- (x = e, G)

        G |- e : bool  G |- (s,_)
        ----------------------------------------
        G |- (while e s, G)

        G |- e : bool  G |- (s1,_)  G |- (s2,_)
        ----------------------------------------
        G |- (if e s1 else s2, G)
  
        G |- e : T
        ----------------------------------------
        G |- (print e, G)
      
    Example of block seen valid by the type checker
    
      x := false;
      while x {
        x := 1
      };
      x := true
     
  Dynamic semantics (interpreter)
  
    Symplifying assumption

      Es wird angenommen, dass jede neu erstelle Variable distinkt ist.

    Values and state
      
      V ::= i | true | false

      S ::= [] | [x : V] | S ++ S
      
    Evaluation rules
    
      Expressions: S | e => V
        
        S |- i => i

        S |- true => true

        S |- false => false

        lookup(S,x) = V
        ------------------
        S |- x => V

        S |- e1 => i1    G |- e2 => i2
        i = i1 + i2
        ----------------------------------------
        S |- e2 + e2 => i

        S |- e1 => i1    G |- e2 => i2
        i = i1 * i2
        ----------------------------------------
        S |- e2 * e2 => i

        G |- e1 => true
        ----------------------------------------
        G |- e2 || e2 => true
       
        G |- e1 => false  G |- e1 => V
        where V is a Boolean value
        ----------------------------------------
        G |- e2 || e2 => V
        
        G |- e1 => false
        ----------------------------------------
        G |- e2 && e2 => false       
        
        G |- e1 => true  G |- e1 => V
        where V is a Boolean value
        ----------------------------------------
        G |- e2 && e2 => V       
        
        G |- e => true
        ----------------------------------------
        G |- ! e => false
               
        G |- e => false
        ----------------------------------------
        G |- ! e => true        
        
        G |- e1 => V   G |- e2 => V
        ----------------------------------------
        G |- e1 == e2 => true        
        
        G |- e1 => V1   G |- e2 => V2
        where V1 and V2 are different values
        ----------------------------------------
        G |- e1 == e2 => false              
        
        G |- e1 => V1   G |- e2 => V2
        where V1 and V2 are numbers and
        V1 is smaller than V2
        ----------------------------------------
        G |- e1 < e2 : true        
        
        G |- e1 => V1   G |- e2 => V2
        where V1 and V2 are numbers and
        V1 is not smaller than V2
        ----------------------------------------
        G |- e1 < e2 : false
        
      Statements: S | s => S2
        
        S1 |- s1 => S2   S2 |- s2 => S3
        ----------------------------------------
        S1 |- s1 ; s2 => S3

        S |- e => V  S2 = S ++ [x : V]
        ----------------------------------------
        S |- x := e => S2
        
        S |- e => V  S2 = S ++ [x : V]
        ----------------------------------------
        S |- x = e => S2
        
        S |- e => false
        ----------------------------------------
        S |- while e s => S
        
        S |- e => true
        S |- s => S2
        S2 |- while e s => S3
        ----------------------------------------
        S |- while e s => S3
        
        S |- e => true   S |- s1 => S2
        ----------------------------------------
        S |- if e s1 else s2 => S2
        
        S |- e => false   S |- s2 => S2
        ----------------------------------------
        S |- if e s1 else s2 => S2
        
        S |- e => V
        "print V on console"
        ----------------------------------------
        S |- print e => S
             
[^1]: Source:  [Lecture-Semantics](https://sulzmann.github.io/ModelBasedSW/lec-semantics.html#(6))

Used [Interface](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L11-L15) for Expression

    type Exp interface {
	    pretty() string
	    eval(s ValState) Val
	    infer(t TyState) (Type, ErrorCodeExpression)
    }
    
  Methods of Exp interface
  
    pretty() string
    
      T
    
    eval(s ValState) Val
    
      T
    
    infer(t TyState) (Type, ErrorCodeExpression)
    
      T
    
Used [Interface](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L19-L23) for Statement

    type Stmt interface {
	    pretty() string
	    eval(s ValState)
	    check(t TyState) (bool, ErrorCodeStatement, ErrorCodeExpression)
    }
    
  Methods of Stmt interface
    
    pretty() string
    
      T
    
    eval(s ValState)
    
      T
    
    check(t TyState) (bool, ErrorCodeStatement, ErrorCodeExpression)
    
      T
    
Tests for different possibilities

  [Test 1 Declaration Statement](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1401-L1404)
  
    Test 1.1 - Declaration
    
	Input: {varX:=3}
 	Output Parse: varX := 3
	Check: true 
	Evalutaion: 

    Test 1.2- False Declaration Example - Error at char 8
    
	Input: {varX:==3}
 	ERROR ON PARSE 
 	AT CHARACTER 8 
  
  [Test 2 Command Sequence Statement](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1406-L1411)
  
    Test 2.1 - Command Sequence - two expressions
  
	Input: {varX:=3;varY:=4}
 	Output Parse: varX := 3 ; varY := 4
 	Check: true 
 	Evalutaion: 

    Test 2.2 - Command Sequence - three expressions
  
	Input: {varX:=3;varY:=4;varZ:=7}
 	Output Parse: varX := 3 ; varY := 4 ; varZ := 7
 	Check: true 
 	Evalutaion: 

    Test 2.3 - False Command Sequence - Error at char 18]()
    
	Input: {varX:=3;varY:=4;;varZ:=7}
 	ERROR ON PARSE 
 	AT CHARACTER 18 

  [Test 3 Print Statement](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1413-L1416)
  
    Test 3.1 - Print - print expression
  
	Input: {print false; print true}
 	Output Parse: print: false ; print: true
 	Check: true 
 	Evalutaion: 
 	false
 	true

    Test 3.2 - False Print - print expression
    
	Input: {print fasle; print true}
 	Output Parse: print: fasle ; print: true
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = PRINT, Reason = Variable not declarated
  
  [Test 4 Assignment Statement](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1418-L1421)
  
    Test 4.1 - Assignment - set varX 4
    
	Input: {varX:=3;print varX; varX=4; print varX}
 	Output Parse: varX := 3 ; print: varX ; varX = 4 ; print: varX
 	Check: true 
 	Evalutaion: 
 	3
 	4
  
    Test 4.2 - False Assignment - not declarated
    
	Input: {varX=4; print varX}
 	Output Parse: varX = 4 ; print: varX
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = ASSIGN, Reason = Variable not declarated
  
  [Test 5 Plus Expression](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1423-L1426)
  
    Test 5.1 - Plus - print varX + varY
    
	Input: {varX:=3;varY:=4;print varX+varY}
 	Output Parse: varX := 3 ; varY := 4 ; print: (varX+varY)
 	Check: true 
 	Evalutaion: 
 	7

    Test 5.2 - False Plus - bool + int
    
	Input: {varX:=true;varY:=4;print varX+varY}
 	Output Parse: varX := true ; varY := 4 ; print: (varX+varY)
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = PRINT, Reason = IllTyped Addition

  [Test 6 Multiplication Expression](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1428-L1431)
  
    Test 6.1 - Mult - print varX * varY
    
	Input: {varX:=3;varY:=4;print varX*varY}
 	Output Parse: varX := 3 ; varY := 4 ; print: (varX*varY)
 	Check: true 
 	Evalutaion: 
 	12

    Test 6.2 - False Mult - bool * int
    
	Input: {varX:=true;varY:=4;print varX*varY}
 	Output Parse: varX := true ; varY := 4 ; print: (varX*varY)
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = PRINT, Reason = IllTyped Multiplication
  
  [Test 7 Lesser Expression](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1433-L1436)
  
    Test 7.1 - Less - print varX < varY
    
	Input: {varX:=3;varY:=4;print varX<varY}
 	Output Parse: varX := 3 ; varY := 4 ; print: (varX<varY)
 	Check: true 
 	Evalutaion: 
 	true
  
    Test 7.2 - False Less - bool < int
    
	Input: {varX:=true;varY:=4;print varX<varY}
 	Output Parse: varX := true ; varY := 4 ; print: (varX<varY)
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = PRINT, Reason = IllTyped Lesser

  [Test 8 And Expression](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1438-L1442)
  
    Test 8.1 - And - print varX && varY
    
	Input: {varX:=true;varY:=true;print varX&&varY}
 	Output Parse: varX := true ; varY := true ; print: (varX&&varY)
 	Check: true 
 	Evalutaion: 
 	true
 
	Input: {varX:=false;varY:=true;print varX&&varY}
	Output Parse: varX := false ; varY := true ; print: (varX&&varY)
 	Check: true 
 	Evalutaion: 
 	false
  
    Test 8.2 - False And - bool && int
    
	Input: {varX:=true;varY:=4;print varX&&varY}
 	Output Parse: varX := true ; varY := 4 ; print: (varX&&varY)
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = PRINT, Reason = IllTyped Conjuction

  [Test 9 Or Expression](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1444-L1448)
  
    Test 9.1 - Or - print varX || varY
    
	Input: {varX:=true;varY:=false;print varX||varY}
 	Output Parse: varX := true ; varY := false ; print: (varX||varY)
 	Check: true 
 	Evalutaion: 
 	true

	Input: {varX:=false;varY:=false;print varX||varY}
 	Output Parse: varX := false ; varY := false ; print: (varX||varY)
 	Check: true 
 	Evalutaion: 
 	false
  
    Test 9.2 - False Or - bool || int
    
	Input: {varX:=true;varY:=4;print varX||varY}
 	Output Parse: varX := true ; varY := 4 ; print: (varX||varY)
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = PRINT, Reason = IllTyped Disjunction

  [Test 10 Equality Expression](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1450-L1456)
  
    Test 10.1 - Equality - print varX == varY
    	
	Input: {varX:=true;varY:=false;print varX==varY}
 	Output Parse: varX := true ; varY := false ; print: (varX==varY)
 	Check: true 
 	Evalutaion: 
 	false

	Input: {varX:=false;varY:=false;print varX==varY}
 	Output Parse: varX := false ; varY := false ; print: (varX==varY)
 	Check: true 
 	Evalutaion: 
 	true

	Input: {varX:=1;varY:=1;print varX==varY}
 	Output Parse: varX := 1 ; varY := 1 ; print: (varX==varY)
 	Check: true 
 	Evalutaion: 
 	true

	Input: {varX:=1;varY:=2;print varX==varY}
 	Output Parse: varX := 1 ; varY := 2 ; print: (varX==varY)
 	Check: true 
 	Evalutaion: 
 	false
  
    Test 10.2 - False Equality - bool == int
    
	Input: {varX:=true;varY:=4;print varX==varY}
 	Output Parse: varX := true ; varY := 4 ; print: (varX==varY)
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = PRINT, Reason = IllTyped Equality

  [Test 11 Negation Expression](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1458-L1462)
  
    Test 11.1 - Negation - print !varX
    
	Input: {varX:=true;print !varX}
 	Output Parse: varX := true ; print: !varX
 	Check: true 
 	Evalutaion: 
 	false

	Input: {varX:=false;print !varX}
 	Output Parse: varX := false ; print: !varX
 	Check: true 
 	Evalutaion: 
 	true
  
    Test 11.2 - False Negation - IllTyped
    
	Input: {varX:=1;print !varX}
 	Output Parse: varX := 1 ; print: !varX
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = PRINT, Reason = IllTyped Negation

  [Test 12  If Else Statement](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1464-L1471)
  
    Test 12.1 - If Else - if varX { print true } else { print false }
    
	Input: {varX:=true;if varX {print true} else {print false}}
	Output Parse: varX := true ; if varX then print: true else print: false
 	Check: true 
 	Evalutaion: 
 	true

    Test 12.2 - If Else - if varX == 1 { print true } else { print false }
    
	Input: {varX:=false;if varX {print true} else {print false}}
 	Output Parse: varX := false ; if varX then print: true else print: false
 	Check: true 
 	Evalutaion: 
 	false

	Input: {varX:=2;if varX == 2 {print true} else {print false}}
 	Output Parse: varX := 2 ; if (varX==2) then print: true else print: false
 	Check: true 
 	Evalutaion: 
 	true
    
     Test 12.3 - If Else - IllTyped 
 
	Input: {varX:=1;if varX {print true} else {print false}}
 	Output Parse: varX := 1 ; if varX then print: true else print: false
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = IF, Reason = Condition IllTyped

  [Test 13 While Statement](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1473-L1478)
  
    Test 13.1 - While - varX:=1; while varX<4 { print varX; varX+1}
    
	Input: {varX:=1;while varX<4 {print varX; varX = varX+1}}
 	Output Parse: varX := 1 ;  while (varX<4) { print: varX ; varX = (varX+1) } 
 	Check: true 
 	Evalutaion: 
 	1
 	2
 	3
  
    Test 13.2 - While - varX:=true; while varX { print varX; varX=false}
    
	Input: {varX:=true;while varX {print varX; varX = false}}
 	Output Parse: varX := true ;  while varX { print: varX ; varX = false } 
 	Check: true 
 	Evalutaion: 
 	true
  
    Test 13.3 - While - IllTyped
    
	Input: {varX:=1;while varX {print varX}}
 	Output Parse: varX := 1 ;  while varX { print: varX } 
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = WHILE, Reason = Condition IllTyped
 
  [Test 14 ExpressionErrorCode Values](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/f601a102fff12f22166031231dc74ccdf336769c/abgabe.go#L1484-L1488)
  
    [Test 14.1 - return value infer/check - Plus](Link)
    
	Input: {varX:=1;varY:=1;varZ:=true;while 1<4 {print varX; if varX<3 {varX = varX+varY}else{varX = varX+varZ}}}
 	Output Parse: varX := 1 ; varY := 1 ; varZ := true ;  while (1<4) { print: varX ; if (varX<3) then varX = (varX+varY) else varX = (varX+varZ) } 
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = ASSIGN, Reason = IllTyped Addition
  
    [Test 14.2 - return value infer/check - Mult](Link)
    
	Input: {varX:=1;varY:=1;varZ:=true;while 1<4 {print varX; if varX<3 {varX = varX+varY}else{varX = varX*varZ}}}
	Output Parse: varX := 1 ; varY := 1 ; varZ := true ;  while (1<4) { print: varX ; if (varX<3) then varX = (varX+varY) else varX = (varX*varZ) } 
 	Check: false 
 	ERROR ON EVALUATION 
 	Illtyped Statement found, StatementType = ASSIGN, Reason = IllTyped Multiplication
