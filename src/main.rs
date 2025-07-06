use crate::error::CompileError;
use std::env;
use std::fs;

mod ast;
mod error;
mod lexer;
mod parser;

use ast::Program;

fn main() {
    if let Err(e) = run() {
        eprintln!("{}", e);
        std::process::exit(1);
    }
}

fn run() -> Result<(), CompileError> {
    let path = env::args()
        .nth(1)
        .ok_or_else(|| CompileError::Io("No input file specified".into()))?;
    let src = fs::read_to_string(path).map_err(|e| CompileError::Io(e.to_string()))?;

    let tokens = lexer::lex(&src)?;
    let mut parser = parser::Parser::new(tokens);
    let prog: Program = parser.parse_program()?;
    println!("{:#?}", prog);

    Ok(())
}
