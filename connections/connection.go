package connections

import (
	"database/contracts"
	"database/kernel/config"
	"database/query"
	"database/sql"
	"fmt"
)

type Connection struct {
	pdo          *sql.DB
	config       *config.DatabaseDriver
	queryGrammar contracts.Grammar
}

func NewConnection(pdo *sql.DB, config *config.DatabaseDriver, grammar contracts.Grammar) *Connection {

	return &Connection{
		pdo:          pdo,
		config:       config,
		queryGrammar: grammar,
	}
}

func (c *Connection) Query() contracts.QueryBuilder {

	return query.NewBuilder(c, c.queryGrammar)
}

func (c *Connection) Select(query string, bindings []interface{}) (*sql.Rows, error) {

	fmt.Println("== SELECT ==", query)

	statement, err := c.pdo.Prepare(query)
	prepareError(err)

	defer statement.Close()

	return statement.Query(bindings...)
}

func (c *Connection) Insert(query string, bindings []interface{}) sql.Result {

	return c.statement(query, bindings)
}

func (c *Connection) Update(query string, bindings []interface{}) int64 {

	return c.affectingStatement(query, bindings)
}

func (c *Connection) Delete(query string, bindings []interface{}) int64 {

	return c.affectingStatement(query, bindings)
}

func (c *Connection) statement(query string, bindings []interface{}) sql.Result {

	statement, err := c.pdo.Prepare(query)
	prepareError(err)

	defer statement.Close()

	res, err := statement.Exec(bindings...)
	prepareError(err)

	return res
}

func (c *Connection) affectingStatement(query string, bindings []interface{}) int64 {

	statement, err := c.pdo.Prepare(query)
	prepareError(err)

	defer statement.Close()

	res, err := statement.Exec(bindings...)
	prepareError(err)

	cont, err := res.RowsAffected()
	if err != nil {
		return 0
	}

	return cont
}

func prepareError(err error) {

	if err != nil {
		panic(err)
	}
}
