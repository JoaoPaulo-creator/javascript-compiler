#[derive(Debug)]
pub struct Program {
    pub functions: Vec<Function>,
    pub statements: Vec<Statement>,
}

#[derive(Debug)]
pub struct Function {
    pub name: String,
    pub params: Vec<String>,
    pub body: Vec<Statement>,
}

#[derive(Debug)]
pub enum Statement {
    Assign {
        name: String,
        expr: Expr,
    },
    Print {
        expr: Expr,
    },
    If {
        cond: Expr,
        then_branch: Vec<Statement>,
        else_branch: Vec<Statement>,
    },

    ExprStmt(Expr),
}

#[derive(Debug)]
pub enum Expr {
    StrLiteral(String),
    Variable(String),
    Binary {
        op: BinOp,
        left: Box<Expr>,
        right: Box<Expr>,
    },
    Call {
        name: String,
        args: Vec<Expr>,
    },
    Length {
        array: Box<Expr>,
    }
}

#[derive(Debug, Clone, Copy)]
pub enum BinOp {
    Add,
    Sub,
    Eq,
    Ne,
    Lt,
    Le,
    Ge,
    Gt,
    Plus,
    Minus,
    Mul,
    Div,
}
