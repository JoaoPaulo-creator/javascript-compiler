use crate::ast::*;
use crate::error::CompileError;
use crate::lexer::Token;

pub struct Parser {
    tokens: Vec<Token>,
    pos: usize, // position
}

impl Parser {
    pub fn new(tokens: Vec<Token>) -> Self {
        Parser { tokens, pos: 0 }
    }

    fn peek(&self) -> &Token {
        &self.tokens[self.pos]
    }

    fn eat(&mut self) {
        if self.pos < self.tokens.len() {
            self.pos += 1;
        }
    }

    fn expect(&mut self, expected: Token) -> Result<(), CompileError> {
        if *self.peek() == expected {
            self.eat();
            Ok(())
        } else {
            Err(CompileError::Parse(format!(
                "Expected {:?}, got {:?}",
                expected,
                self.peek()
            )))
        }
    }

    pub fn parse_program(&mut self) -> Result<Program, CompileError> {
        let mut funcs = Vec::new();
        let mut stmts = Vec::new();

        while *self.peek() != Token::EOF {
            if *self.peek() == Token::Fn {
                funcs.push(self.parse_function()?);
            } else {
                let stmt = self.parse_statement()?;
                stmts.push(stmt);
            }
        }
        Ok(Program {
            functions: funcs,
            statements: stmts,
        })
    }

    fn parse_function(&mut self) -> Result<Function, CompileError> {
        self.expect(Token::Fn)?;
        let name = match self.peek() {
            Token::Ident(n) => n.clone(),
            _ => return Err(CompileError::Parse("Expected function name".into())),
        };

        self.eat();
        self.expect(Token::LParen)?;

        let mut params = Vec::new();
        if *self.peek() == Token::RParen {
            loop {
                if let Token::Ident(n) = self.peek() {
                    params.push(n.clone());
                    self.eat();
                } else {
                    return Err(CompileError::Parse("Expected parameter name".into()));
                }

                if *self.peek() == Token::Comma {
                    self.eat();
                    continue;
                }

                break;
            }
        }

        self.expect(Token::RParen)?;
        let body = self.parse_block()?;
        Ok(Function { name, params, body })
    }

    fn parse_block(&mut self) -> Result<Vec<Statement>, CompileError> {
        self.expect(Token::LBrace)?;
        let mut v = Vec::new();
        while *self.peek() != Token::RBrace {
            let stmt = self.parse_statement()?;
            v.push(stmt);
        }

        self.expect(Token::RBrace)?;
        Ok(v)
    }

    fn parse_statement(&mut self) -> Result<Statement, CompileError> {
        match self.peek() {
            Token::Print => {
                self.eat();
                let expr = self.parse_expr()?;
                Ok(Statement::Print { expr })
            }
            Token::Ident(_) => {
                let expr = self.parse_expr()?;
                if *self.peek() == Token::Eq {
                    self.eat();
                    let value = self.parse_expr()?;
                    match expr {
                        Expr::Variable(name) => Ok(Statement::Assign { name, expr: value }),
                        _ => Err(CompileError::Parse(
                            "Left-hand side of assignment must be a variable or array index".into(),
                        )),
                    }
                } else {
                    Ok(Statement::ExprStmt(expr))
                }
            }
            _ => {
                let expr = self.parse_expr()?;
                Ok(Statement::ExprStmt(expr))
            }
        }
    }

    fn parse_expr(&mut self) -> Result<Expr, CompileError> {
        self.parse_equality()
    }

    fn parse_equality(&mut self) -> Result<Expr, CompileError> {
        let mut lhs = self.parse_comparison()?;
        while matches!(self.peek(), Token::EqEq | Token::Ne) {
            let op = match self.peek() {
                Token::EqEq => BinOp::Eq,
                Token::Ne => BinOp::Ne,
                _ => unreachable!(),
            };
            self.eat();
            let rhs = self.parse_comparison()?;
            lhs = Expr::Binary {
                op,
                left: Box::new(lhs),
                right: Box::new(rhs),
            };
        }
        Ok(lhs)
    }

    fn parse_comparison(&mut self) -> Result<Expr, CompileError> {
        let mut lhs = self.parse_addition()?;
        while matches!(self.peek(), Token::Lt | Token::Le | Token::Gt | Token::Ge) {
            let op = match self.peek() {
                Token::Lt => BinOp::Lt,
                Token::Le => BinOp::Le,
                Token::Gt => BinOp::Gt,
                Token::Ge => BinOp::Ge,
                _ => unreachable!(),
            };

            self.eat();
            let rhs = self.parse_addition()?;
            lhs = Expr::Binary {
                op,
                left: Box::new(lhs),
                right: Box::new(rhs),
            }
        }

        Ok(lhs)
    }

    fn parse_addition(&mut self) -> Result<Expr, CompileError> {
        let mut lhs = self.parse_term()?;
        while matches!(self.peek(), Token::Plus | Token::Minus) {
            let op = if *self.peek() == Token::Plus {
                BinOp::Add
            } else {
                BinOp::Sub
            };

            self.eat();
            let rhs = self.parse_term()?;
            lhs = Expr::Binary {
                op,
                left: Box::new(lhs),
                right: Box::new(rhs),
            }
        }

        Ok(lhs)
    }

    fn parse_term(&mut self) -> Result<Expr, CompileError> {
        let mut lhs = self.parse_factor()?;
        while matches!(self.peek(), Token::Star | Token::Slash) {
            let op = match self.peek() {
                Token::Star => BinOp::Mul,
                Token::Slash => BinOp::Div,
                _ => unreachable!(),
            };

            self.eat();
            let rhs = self.parse_factor()?;
            lhs = Expr::Binary {
                op,
                left: Box::new(lhs),
                right: Box::new(rhs),
            }
        }

        Ok(lhs)
    }

    fn parse_factor(&mut self) -> Result<Expr, CompileError> {
        let mut node = match self.peek() {
            Token::StrLiteral(s) => {
                let v = s.clone();
                self.eat();
                Expr::StrLiteral(v)
            }
            Token::Ident(name) => {
                let name = name.clone();
                self.eat();
                if *self.peek() == Token::LParen {
                    // function call
                    self.eat(); // '('
                    let mut args = Vec::new();
                    if *self.peek() != Token::RParen {
                        loop {
                            args.push(self.parse_expr()?);
                            if *self.peek() == Token::Comma {
                                self.eat();
                                continue;
                            }
                            break;
                        }
                    }
                    self.expect(Token::RParen)?;
                    Expr::Call { name, args }
                } else {
                    Expr::Variable(name)
                }
            }
            Token::LParen => {
                self.eat();
                let e = self.parse_expr()?;
                self.expect(Token::RParen)?;
                e
            }
            other => {
                return Err(CompileError::Parse(format!(
                    "Unexpected token in factor: {:?}",
                    other
                )));
            }
        };

        loop {
            match self.peek() {
                Token::Dot => {
                    self.eat();
                    match self.peek() {
                        Token::Ident(method_name) if method_name == "length" => {
                            self.eat();
                            self.expect(Token::LParen)?;
                            self.expect(Token::RParen)?;
                            node = Expr::Length {
                                array: Box::new(node),
                            };
                        }
                        other => {
                            return Err(CompileError::Parse(format!(
                                "Expected 'length' after '.', found {:?}",
                                other
                            )));
                        }
                    }
                }
                _ => break,
            }
        }

        Ok(node)
    }
}
