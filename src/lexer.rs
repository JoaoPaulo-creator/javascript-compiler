use crate::error::CompileError;

#[derive(Debug, PartialEq)]
pub enum Token {
    // identifiers and literals
    Ident(String),
    Number(i64),
    StrLiteral(String),
    BoolLiteral(bool),
    Fn,
    // keywords
    Print,
    // delimiters
    LParen,
    RParen,
    LBrace,
    RBrace,
    LBracket,
    RBracket,
    Comma,
    Dot,
    EOF,

    //tokens
    Gt,
    Ge,
    Eq,
    EqEq,
    Ne,
    NeEq,
    Lt,
    Le,
    Plus,
    Minus,
    Star,
    Slash,
}

fn is_ident_start(c: char) -> bool {
    return c.is_ascii_alphabetic() || c == '_';
}

fn is_ident_continue(c: char) -> bool {
    return c.is_ascii_alphabetic() || c == '_';
}

pub fn lex(input: &str) -> Result<Vec<Token>, CompileError> {
    let mut tokens = Vec::new();
    let mut chars = input.chars().peekable();

    while let Some(&ch) = chars.peek() {
        match ch {
            c if c.is_whitespace() => {
                chars.next();
            }

            '(' => {
                chars.next();
                tokens.push(Token::LParen);
            }
            ')' => {
                chars.next();
                tokens.push(Token::RParen);
            }
            '"' => {
                chars.next();
                let mut s = String::new();
                while let Some(&c2) = chars.peek() {
                    if c2 == '"' {
                        chars.next();
                        break;
                    }

                    s.push(c2);
                    chars.next();
                }
                tokens.push(Token::StrLiteral(s));
            }

            c if is_ident_start(c) => {
                let mut ident = String::new();
                ident.push(c);
                chars.next();
                while let Some(&c2) = chars.peek() {
                    if is_ident_continue(c2) {
                        ident.push(c2);
                        chars.next();
                    } else {
                        break;
                    }
                }

                let tok = match ident.as_str() {
                    "print" => Token::Print,
                    _ => Token::Ident(ident),
                };
            }
            //
            other => {
                return Err(CompileError::Lex(format!(
                    "unexpected character '{}'",
                    other
                )));
            }
        }
    }

    tokens.push(Token::EOF);
    Ok(tokens)
}
