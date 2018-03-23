// Code generated by gnorm. Source: ../gnorm/templates/store.gotmpl. DO NOT EDIT!

package armoury

import (
	"fmt"
	"strconv"
	"strings"
)

// QueryClause builds a query clause from a where clause and a set of sorts.
func QueryClause(whereClause WhereClause, order Sort) Clause {
	return Clause{
		whereClause: whereClause,
		order:       order,
	}
}

// Clause is a where clause and a set of sorts.
type Clause struct {
	whereClause WhereClause
	order       Sort
}

func (c Clause) String(idx *int) string {
	var ret string

	if where := c.whereClause.String(idx); where != "" {
		ret = fmt.Sprintf(" WHERE %s", where)
	}

	if order := c.order.String(); order != "" {
		ret += fmt.Sprintf(" ORDER BY %s", order)
	}

	return ret
}

func (c Clause) Values() []interface{} { // nolint: golint
	return c.whereClause.Values()
}

// OrderBy receives a variadic list of orderings and returns an Orderings type.
func OrderBy(orderings ...Ordering) Orderings {
	return Orderings(orderings)
}

// Sort is the interface that a printable sort/order by statement should satisfy.
type Sort interface {
	isSortable()
	String() string
}

// Orderings is a slice of orderings.
type Orderings []Ordering

func (o Orderings) isSortable() {}

func (o Orderings) String() string {
	orderings := make([]string, 0, len(o))
	for i := 0; i < len(o); i++ {
		orderings = append(orderings, o[i].String())
	}

	if len(orderings) == 0 {
		return ""
	}
	return strings.Join(orderings, ", ")
}

// UnOrdered is a convenience value to make it clear you're not sorting a query.
var UnOrdered = Ordering{}

// OrderByDesc returns a descending sort on the given field.
func OrderByDesc(field string) Ordering {
	return Ordering{
		field: field,
		order: orderDesc,
	}
}

// OrderByAsc returns an ascending sort on the given field.
func OrderByAsc(field string) Ordering {
	return Ordering{
		field: field,
		order: orderAsc,
	}
}

// Ordering indicates how rows should be sorted.
type Ordering struct {
	field string
	order sortOrder
}

func (o Ordering) isSortable() {}

func (o Ordering) String() string {
	if o.order == orderNone {
		return ""
	}
	return fmt.Sprintf("%s %s", o.field, o.order.String())
}

// sortOrder defines how to order rows returned.
type sortOrder int

// Defined sort orders for not sorted, descending and ascending.
const (
	orderNone sortOrder = iota
	orderDesc
	orderAsc
)

// String returns the sql string representation of this sort order.
func (s sortOrder) String() string {
	switch s {
	case orderDesc:
		return "DESC"
	case orderAsc:
		return "ASC"
	}
	return ""
}

// WhereClause has a String function should return a properly formatted where
// clause (not including the WHERE) for positional arguments starting at idx.
type WhereClause interface {
	String(idx *int) string
	Values() []interface{}
	isWhere()
}

// NoWhere is a convenience value to make it clear you're not utilising a where in a query.
var NoWhere = NoopWhere{}

// NoopWhere is a convenience struct to make allow a no op where.
type NoopWhere struct{}

func (n NoopWhere) isWhere() {}

func (n NoopWhere) String(idx *int) string {
	return ""
}

func (n NoopWhere) Values() []interface{} { // nolint: golint
	return []interface{}{}
}

// Comparison is used by WhereClauses to create valid sql.
type Comparison string

// Comparison types.
const (
	CompEqual   Comparison = "="
	CompGreater Comparison = ">"
	CompLess    Comparison = "<"
	CompGTE     Comparison = ">="
	CompLTE     Comparison = "<="
	CompNE      Comparison = "<>"
)

// Where is a clause that checks a column against a comparison value.
type Where struct {
	Field string
	Comp  Comparison
	Value interface{}
}

func (w Where) isWhere() {}

func (w Where) String(idx *int) string {
	str := fmt.Sprintf("%s %s $%s", w.Field, string(w.Comp), strconv.Itoa(*idx))
	(*idx)++
	return str
}

func (w Where) Values() []interface{} { // nolint: golint
	return []interface{}{w.Value}
}

// NullClause is a clause that checks for a column being null or not.
type NullClause struct {
	Field string
	Null  bool
}

func (n NullClause) isWhere() {}

func (n NullClause) String(idx *int) string {
	if n.Null {
		return fmt.Sprintf("%s IS NULL", n.Field)
	}
	return fmt.Sprintf("%s IS NOT NULL", n.Field)
}

func (n NullClause) Values() []interface{} { // nolint: golint
	return []interface{}{}
}

// AndClause returns a WhereClause that serializes to the AND
// of all the given where clauses.
func AndClause(wheres ...WhereClause) WhereClause {
	return andClause(wheres)
}

type andClause []WhereClause

func (a andClause) isWhere() {}

func (a andClause) String(idx *int) string {
	wheres := make([]string, len(a))
	for x := 0; x < len(a); x++ {
		wheres[x] = fmt.Sprintf("(%s)", a[x].String(idx))
	}
	return strings.Join(wheres, " AND ")
}

func (a andClause) Values() []interface{} { // nolint: golint
	vals := make([]interface{}, 0, len(a))
	for x := 0; x < len(a); x++ {
		vals = append(vals, a[x].Values()...)
	}
	return vals
}

// OrClause returns a WhereClause that serializes to the OR
// of all the given where clauses.
func OrClause(wheres ...WhereClause) WhereClause {
	return orClause(wheres)
}

type orClause []WhereClause

func (o orClause) isWhere() {}

func (o orClause) String(idx *int) string {
	wheres := make([]string, len(o))
	for x := 0; x < len(wheres); x++ {
		wheres[x] = fmt.Sprintf("(%s)", o[x].String(idx))
	}
	return strings.Join(wheres, " OR ")
}

func (o orClause) Values() []interface{} { // nolint: golint
	vals := make([]interface{}, len(o))
	for x := 0; x < len(o); x++ {
		vals = append(vals, o[x].Values()...)
	}
	return vals
}

// InClause takes a slice of values that it matches against.
type InClause struct {
	Field string
	Vals  []interface{}
}

func (in InClause) isWhere() {}

func (in InClause) String(idx *int) string {
	str := in.Field + " IN("
	for x := range in.Vals {
		if x != 0 {
			str += ", "
		}
		str += "$" + strconv.Itoa(*idx)
		(*idx)++
	}
	str += ")"
	return str
}

func (in InClause) Values() []interface{} { // nolint: golint
	return in.Vals
}
