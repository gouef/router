package router

type Route struct {
	mask     string
	sequence []interface{}

	/** regular expression pattern */
	re string

	/** @var string[]  parameter aliases in regular expression */
	aliases []string

	/** @var array of [value & fixity, filterIn, filterOut] */
	metadata []interface{}
	xlat     []interface{}

	/** Host, Path, Relative */
	routeType int

	/** http | https */
	scheme string
}
