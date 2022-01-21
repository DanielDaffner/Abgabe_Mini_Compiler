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

Used Interfaces for Typechecker

Used Interfaces for Interpreter

Tests for different opportunities

