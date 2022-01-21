# Abgabe_MSE
Mini-Compiler Projekt für Modelbasierte Software Entwicklung

Kurzbeschreibung

  Ein gegebener Input wird anhand der Regeln einer einfachen imperativen Programmiersprache geparsed, überprüft und anschließend interpretiert.

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

Used [Interface](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L11-L15) for Typechecker

    type Exp interface {
	    pretty() string
	    eval(s ValState) Val
	    infer(t TyState) (Type, ErrorCodeExpression)
    }
    
  Methods of [Exp](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L19-L23) interface
  
    pretty() string
    
      T
    
    eval(s ValState) Val
    
      T
    
    infer(t TyState) (Type, ErrorCodeExpression)
    
      T
    
Used Interface for Interpreter

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

  [Test 1](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1401-L1404)
  
    Test 1.1 - Declaration
  
    Test 1.2- False Declaration Example - Error at char 8
  
  [Test 2](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1406-L1411)
  
    Test 2.1 - Command Sequence - two expressions
  
    Test 2.2 - Command Sequence - three expressions
  
    Test 2.3 - False Command Sequence - Error at char 18]()
    
  [Test 3](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1413-L1416)
  
    Test 3.1 - Print - print expression
  
    Test 3.2 - False Print - print expression
  
  [Test 4](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1418-L1421)
  
    Test 4.1 - Assignment - set varX 4
  
    Test 4.2 - False Assignment - not declarated
  
  [Test 5](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1423-L1426)
  
    Test 5.1 - Plus - print varX + varY
  
    Test 5.2 - False Plus - bool + int
  
  [Test 6](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1428-L1431)
  
    Test 6.1 - Mult - print varX * varY
  
    Test 6.2 - False Mult - bool * int
  
  [Test 7](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1433-L1436)
  
    Test 7.1 - Less - print varX < varY
  
    Test 7.2 - False Less - bool < int
  
  [Test 8](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1438-L1442)
  
    Test 8.1 - And - print varX && varY
  
    Test 8.2 - False And - bool && int
  
  [Test 9](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1444-L1448)
  
    Test 9.1 - Or - print varX || varY
  
    [Test 9.2 - False Or - bool || int
  
  [Test 10](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1450-L1456)
  
    Test 10.1 - Equality - print varX == varY
  
    Test 10.2 - False Equality - bool == int
  
  [Test 11](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1458-L1462)
  
    Test 11.1 - Negation - print !varX
  
    Test 11.2 - False Negation - IllTyped
  
  [Test 12](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1464-L1471)
  
    Test 12.1 - If Else - if varX { print true } else { print false }
  
    Test 12.2 - If Else - if varX == 1 { print true } else { print false }
  
  [Test 13](https://github.com/DanielDaffner/Abgabe_Mini_Compiler/blob/e0b7b71241ad48128fd4002078c1d8dde65b08dd/abgabe.go#L1473-L1478)
  
    Test 13.1 - While - varX:=1; while varX<4 { print varX; varX+1}
  
    Test 13.2 - While - varX:=true; while varX { print varX; varX=false}
  
    Test 13.3 - While - IllTyped
     
  [Test 14](Link)
  
    [Test 14.1 - return value infer/check - Plus](Link)
  
    [Test 14.2 - return value infer/check - Mult](Link)
